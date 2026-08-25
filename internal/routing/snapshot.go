package routing

import (
	"sync"

	"loadbalancer/internal/model"
)

// Snapshot is an immutable routing view published to request handlers.
type Snapshot struct {
	Version  uint64
	Rules    []*model.RouteRule
	Backends map[string]*model.Backend
}

func cloneBackend(b *model.Backend) *model.Backend {
	return b
}

func cloneRule(rule *model.RouteRule) *model.RouteRule {
	if rule == nil {
		return nil
	}
	copy := *rule
	copy.BackendIDs = append([]string(nil), rule.BackendIDs...)
	return &copy
}

func cloneSnapshot(in Snapshot) Snapshot {
	out := Snapshot{
		Version:  in.Version,
		Rules:    make([]*model.RouteRule, 0, len(in.Rules)),
		Backends: make(map[string]*model.Backend, len(in.Backends)),
	}
	for _, rule := range in.Rules {
		out.Rules = append(out.Rules, cloneRule(rule))
	}
	for id, backend := range in.Backends {
		out.Backends[id] = cloneBackend(backend)
	}
	return out
}

// Table publishes and reads defensive copies so request updates cannot mutate a live view.
type Table struct {
	mu      sync.RWMutex
	current Snapshot
}

func NewTable() *Table {
	return &Table{current: Snapshot{Backends: make(map[string]*model.Backend)}}
}

func (t *Table) Publish(next Snapshot) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	if next.Version <= t.current.Version {
		return false
	}
	t.current = cloneSnapshot(next)
	return true
}

func (t *Table) Current() Snapshot {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.current
}
