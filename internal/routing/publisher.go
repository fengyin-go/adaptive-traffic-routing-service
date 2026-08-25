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
			copy := *backend
			byID[backend.ID] = &copy
		}
	}
	ownedRules := make([]*model.RouteRule, 0, len(rules))
	for _, rule := range rules {
		if rule == nil {
			continue
		}
		copy := *rule
		ownedRules = append(ownedRules, &copy)
	}
	return p.table.Publish(Snapshot{Version: version, Rules: ownedRules, Backends: byID})
}
