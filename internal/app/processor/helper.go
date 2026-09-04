package processor

import (
	"context"
	"sync"
)

func Wrap(ctx context.Context, wg *sync.WaitGroup, cb func(context.Context)) {
	if wg != nil {
		wg.Add(1)
	}

	go func() {
		defer func() {
			if wg != nil {
				wg.Done()
			}
		}()

		select {
		case <-ctx.Done():
			return
		default:
			cb(ctx)
		}
	}()
}
