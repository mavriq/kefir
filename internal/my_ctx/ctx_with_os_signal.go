package my_ctx

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func CtxWithOsSignal() context.Context {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		defer cancel()

		<-sigCh
	}()

	return ctx
}
