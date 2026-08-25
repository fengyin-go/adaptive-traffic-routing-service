package routing

import "loadbalancer/internal/model"

// Publisher assembles store data into a single versioned routing snapshot.
type Publisher struct {
	table *Table
}

func NewPublisher(table *Table) *Publisher {
	return &Publisher{table: table}
}

func (p *Publisher) Refresh(version uint64, rules []*model.RouteRule, backends []*model.Backend) bool {
	byID := make(map[string]*model.Backend, len(backends))
	for _, backend := range backends {
		if backend != nil {
			byID[backend.ID] = backend
		}
	}
	ownedRules := make([]*model.RouteRule, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		ownedRules = append(ownedRules, rule)
	}
	return p.table.Publish(Snapshot{Version: version, Rules: ownedRules, Backends: byID})
}
