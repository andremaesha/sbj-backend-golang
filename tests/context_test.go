package tests

import (
	"context"
	"testing"
	"time"
)

func TestFirstContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(6 * time.Second)
		cancel()
	}()

	select {
	case <-ctx.Done():
		println("operation cancelled: ", ctx.Err())
	case <-time.After(5 * time.Second):
		println("operation complated")
	}
}
