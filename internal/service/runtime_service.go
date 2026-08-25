package service

import (
	"context"
	"fmt"

	"loadbalancer/internal/health"
	"loadbalancer/internal/requestscope"
	"loadbalancer/internal/resilience"
	"loadbalancer/internal/routing"
)

// ProbeBackends checks a batch and applies every completed result to backend state.
func (s *Service) ProbeBackends(ctx context.Context, coordinator *health.Coordinator, ids []string) error {
	targets := make([]health.Target, 0, len(ids))
	for _, id := range ids {
		backend, err := s.store.GetBackend(id)
		if err != nil {
			return err
		}
		targets = append(targets, health.Target{ID: backend.ID, Host: backend.Host, Port: backend.Port})
	}
	results, err := coordinator.ProbeBatch(ctx, targets)
	if err != nil {
		return err
	}
	if len(results) != len(targets) {
		return fmt.Errorf("health batch incomplete: got %d of %d results", len(results), len(targets))
	}
	for _, result := range results {
		if _, err := s.RecordHealthCheck(result.Target.ID, result.Err == nil && result.Healthy); err != nil {
			return err
		}
	}
	return nil
}

// PublishRoutingSnapshot builds one immutable view from the current store state.
func (s *Service) PublishRoutingSnapshot(publisher *routing.Publisher, version uint64) bool {
	return publisher.Refresh(version, s.store.ListRouteRules(), s.store.ListBackends())
}

// ResolveWithRetry commits a selected backend only after a retry sequence succeeds.
func (s *Service) ResolveWithRetry(ctx context.Context, max int, attempt func(context.Context, int) (string, error), commit func(string) error) error {
	var selected string
	err := resilience.Retry(ctx, max, func(ctx context.Context, number int) error {
		backendID, err := attempt(ctx, number)
		if err == nil {
			selected = backendID
		}
		return err
	})
	if err != nil {
		return err
	}
	return commit(selected)
}

// RunHealthSweep owns one resource-bounded sequential probe pass.
func (s *Service) RunHealthSweep(ctx context.Context, pool *health.LeasePool, targets []health.Target, check func(health.Target) error) error {
	return health.RunLeasedChecks(ctx, pool, targets, check)
}

// AuditRequest snapshots request-owned data before handing it to asynchronous work.
func (s *Service) AuditRequest(ctx context.Context, labels []string, sink *requestscope.AuditSink, ready <-chan struct{}) {
	scope := requestscope.Acquire(ctx, labels)
	snapshot := scope.Snapshot()
	requestscope.Release(scope)
	go func() {
		if ready != nil {
			<-ready
		}
		sink.Record(snapshot)
	}()
}
