package ctxutil

import (
	"context"
	"net/http"
	"time"
)

func ListenUntilCancelled(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	errors := make(chan error)
	go func() { errors <- server.ListenAndServe() }()
	select {
	case err := <-errors:
		return err
	case <-ctx.Done():
		return CallWithTimeout(context.Background(), shutdownTimeout, server.Shutdown)
	}
}
