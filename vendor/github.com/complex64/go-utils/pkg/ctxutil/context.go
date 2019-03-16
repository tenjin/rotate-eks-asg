package ctxutil

import (
	"context"
	"os"
	"os/signal"
	"time"
)

func ContextWithCancelSignals(sig ...os.Signal) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, sig...)
	go func() {
		<-exit
		cancel()
	}()
	return ctx, cancel
}

func CallWithTimeout(parentCtx context.Context, timeout time.Duration, function func(context.Context) error) error {
	ctx, cancel := context.WithTimeout(parentCtx, timeout)
	defer cancel()
	errors := make(chan error)
	go func() { errors <- function(ctx) }()
	select {
	case err := <-errors:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
