package snippet

import "sync"

func MergeChans[T any](ln int, cs ...<-chan T) <-chan T {
	wg := sync.WaitGroup{}
	out := make(chan T, ln)

	wg.Add(len(cs))

	for _, in := range cs {
		go func(in <-chan T) {
			defer wg.Done()

			for v := range in {
				out <- v
			}
		}(in)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
