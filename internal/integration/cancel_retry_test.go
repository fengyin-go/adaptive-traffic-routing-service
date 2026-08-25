package integration

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"loadbalancer/internal/health"
	"loadbalancer/internal/resilience"
)

type cancelAwareChecker struct{ active atomic.Int32 }

func (c *cancelAwareChecker) Check(ctx context.Context, _ health.Target) (bool, error) {
	c.active.Add(1)
	defer c.active.Add(-1)
	<-ctx.Done()
	return false, ctx.Err()
}

func TestCancellationStopsRetryDelayAndProbeWorkers(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	started := time.Now()
	err := resilience.Retry(ctx, 3, func(context.Context, int) error {
		return &resilience.TemporaryError{After: 150 * time.Millisecond, Err: errors.New("warming")}
	})
	if !errors.Is(err, context.DeadlineExceeded) || time.Since(started) > 90*time.Millisecond {
		t.Fatalf("cancelled route retry kept waiting past its request deadline: err=%v elapsed=%v", err, time.Since(started))
	}
	checker := &cancelAwareChecker{}
	probeCtx, probeCancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	_, _ = health.NewCoordinator(checker).ProbeBatch(probeCtx, []health.Target{{ID: "edge-a"}})
	probeCancel()
	time.Sleep(20 * time.Millisecond)
	if checker.active.Load() != 0 {
		t.Fatalf("cancelled probe request left %d backend worker running", checker.active.Load())
	}
}
