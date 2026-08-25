package health

import (
	"context"
	"errors"
	"sync"
)

var ErrLeaseLimit = errors.New("health probe lease limit reached")

// LeasePool bounds concurrently open probe connections.
type LeasePool struct {
	mu     sync.Mutex
	limit  int
	active int
	peak   int
}

func NewLeasePool(limit int) *LeasePool {
	return &LeasePool{limit: limit}
}

func (p *LeasePool) Acquire(ctx context.Context) (func(), error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	if p.active >= p.limit {
		return nil, ErrLeaseLimit
	}
	p.active++
	if p.active > p.peak {
		p.peak = p.active
	}
	var once sync.Once
	return func() {
		once.Do(func() {
			p.mu.Lock()
			p.active--
			p.mu.Unlock()
		})
	}, nil
}

func (p *LeasePool) Active() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.active
}

func (p *LeasePool) Peak() int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.peak
}

// RunLeasedChecks runs one probe at a time, releasing each lease before the
// next target begins. Sequential sweeps therefore never hold more than one
// lease at once, so a limit of N comfortably covers N targets checked in turn.
func RunLeasedChecks(ctx context.Context, pool *LeasePool, targets []Target, check func(Target) error) error {
	for _, target := range targets {
		release, err := pool.Acquire(ctx)
		if err != nil {
			return err
		}
		err = check(target)
		release()
		if err != nil {
			return err
		}
	}
	return nil
}
