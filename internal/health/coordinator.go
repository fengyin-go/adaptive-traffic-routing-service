package health

import (
	"context"
	"sync"
)

// Target identifies one backend health check.
type Target struct {
	ID   string
	Host string
	Port int
}

// Result is the outcome of checking one backend.
type Result struct {
	Target  Target
	Healthy bool
	Err     error
}

// Checker performs one backend health check.
type Checker interface {
	Check(context.Context, Target) (bool, error)
}

// Coordinator owns the worker and result-channel lifecycle for a probe batch.
type Coordinator struct {
	checker Checker
}

func NewCoordinator(checker Checker) *Coordinator {
	return &Coordinator{checker: checker}
}

func (c *Coordinator) ProbeBatch(ctx context.Context, targets []Target) ([]Result, error) {
	batchCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	results := make(chan Result, len(targets))
	var wg sync.WaitGroup
	wg.Add(len(targets))
	for _, target := range targets {
		target := target
		go func() {
			defer wg.Done()
			healthy, err := c.checker.Check(batchCtx, target)
			select {
			case results <- Result{Target: target, Healthy: healthy, Err: err}:
			case <-batchCtx.Done():
			}
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()

	out := make([]Result, 0, len(targets))
	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case result, ok := <-results:
			if !ok {
				return out, nil
			}
			out = append(out, result)
		}
	}
}
