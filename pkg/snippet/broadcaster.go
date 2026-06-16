package snippet

import (
	"context"
	"sync"
	"time"
)

// Broadcaster описывает объект, способный дублировать поток событий для множества подписчиков.
type Broadcaster[T any] interface {
	// BestEffort — пропустить, если занято
	Subscribe(ctx context.Context, buffer int) <-chan T
	// WithTimeout — ждать готовности N времени
	SubscribeWait(ctx context.Context, buffer int, timeout time.Duration) <-chan T
	// Strict — блокировать до победного (или отмены контекста)
	SubscribeStrict(ctx context.Context, buffer int) <-chan T

	// Публикация данных для всех подписчиков
	Publish(val T)

	Close()
}

// broadcaster — скрытая реализация интерфейса.
type broadcaster[T any] struct {
	mu   sync.RWMutex
	subs map[chan T]time.Duration // 0 - BestEffort, -1 - Strict, >0 - Timeout

	source chan T
}

func NewBroadcaster[T any]() Broadcaster[T] {
	return NewBroadcasterWithLen[T](0)
}

// NewBroadcaster создает и запускает новый транслятор событий.
func NewBroadcasterWithLen[T any](bufLen int) Broadcaster[T] {
	b := &broadcaster[T]{
		subs:   make(map[chan T]time.Duration),
		source: make(chan T, bufLen),
	}

	go b.run()
	return b
}

func (b *broadcaster[T]) run() {
	defer func() {
		b.mu.Lock()
		for ch := range b.subs {
			close(ch)
		}
		b.mu.Unlock()
	}()

	for val := range b.source {
		b.mu.RLock()
		var wg sync.WaitGroup

		// Рассылаем всем параллельно, чтобы один тип очереди не ждал другой
		for ch, timeout := range b.subs {
			wg.Add(1)
			go func(ch chan T, timeout time.Duration, val T) {
				defer wg.Done()
				defer func() { recover() }()

				switch {
				case timeout == 0: // BestEffort
					select {
					case ch <- val:
					default:
					}
				case timeout < 0: // Strict
					ch <- val
				default: // With Timeout
					t := time.NewTimer(timeout)
					defer t.Stop()
					select {
					case ch <- val:
					case <-t.C:
					}
				}
			}(ch, timeout, val)
		}
		wg.Wait()
		b.mu.RUnlock()
	}
}

func (b *broadcaster[T]) SubscribeWait(ctx context.Context, buffer int, timeout time.Duration) <-chan T {
	ch := make(chan T, buffer)
	b.mu.Lock()
	defer b.mu.Unlock()
	b.subs[ch] = timeout

	go func() {
		<-ctx.Done()
		close(ch)
		b.mu.Lock()
		defer b.mu.Unlock()
		delete(b.subs, ch)
	}()

	return ch
}

func (b *broadcaster[T]) Subscribe(ctx context.Context, buffer int) <-chan T {
	return b.SubscribeWait(ctx, buffer, time.Duration(0*time.Second))
}

func (b *broadcaster[T]) SubscribeStrict(ctx context.Context, buffer int) <-chan T {
	return b.SubscribeWait(ctx, buffer, time.Duration(-1*time.Second))
}

func (b *broadcaster[T]) Publish(val T) {
	b.source <- val
}

func (b *broadcaster[T]) Close() {
	close(b.source)
}
