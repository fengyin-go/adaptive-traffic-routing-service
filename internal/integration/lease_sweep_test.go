package integration

import (
	"context"
	"sync/atomic"
	"testing"

	"loadbalancer/internal/config"
	"loadbalancer/internal/health"
	"loadbalancer/internal/service"
	"loadbalancer/internal/store"
	"loadbalancer/pkg/logger"
)

func TestHealthSweepReleasesEachLeaseBeforeNextProbe(t *testing.T) {
	svc := service.New(store.NewMemoryStore(), logger.New(), &config.Config{})
	pool := health.NewLeasePool(2)
	targets := []health.Target{{ID: "edge-a"}, {ID: "edge-b"}, {ID: "edge-c"}}
	var checks atomic.Int32
	err := svc.RunHealthSweep(context.Background(), pool, targets, func(health.Target) error {
		checks.Add(1)
		return nil
	})
	if err != nil || checks.Load() != 3 || pool.Active() != 0 || pool.Peak() > 1 {
		t.Fatalf("health sweep exhausted leases or repeated probes: err=%v checks=%d active=%d peak=%d", err, checks.Load(), pool.Active(), pool.Peak())
	}
}
