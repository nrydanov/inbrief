package channels

import (
	"context"
	"sync"
)

func FanIn[T any](
	ctx context.Context,
	chs ...<-chan T,
) <-chan T {
	out := make(chan T)
	wg := sync.WaitGroup{}
	wg.Add(len(chs))
	for _, ch := range chs {
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case v, ok := <-ch:
					if !ok {
						return
					}
					out <- v
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
