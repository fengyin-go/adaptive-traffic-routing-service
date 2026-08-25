package store

import (
	"loadbalancer/internal/model"
)

func cloneRouteRule(rule *model.RouteRule) *model.RouteRule {
	if rule == nil {
		return nil
	}
	copy := *rule
	return &copy
}

func (s *MemoryStore) CreateRouteRule(r *model.RouteRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.routeRules {
		if exist.PathPrefix == r.PathPrefix {
			return ErrConflict
		}
	}
	s.routeRules[r.ID] = cloneRouteRule(r)
	return nil
}

func (s *MemoryStore) GetRouteRule(id string) (*model.RouteRule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.routeRules[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneRouteRule(r), nil
}

func (s *MemoryStore) ListRouteRules() []*model.RouteRule {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RouteRule, 0, len(s.routeRules))
	for _, r := range s.routeRules {
		list = append(list, cloneRouteRule(r))
	}
	return list
}

func (s *MemoryStore) UpdateRouteRule(r *model.RouteRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.routeRules[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.routeRules {
		if exist.ID != r.ID && exist.PathPrefix == r.PathPrefix {
			return ErrConflict
		}
	}
	s.routeRules[r.ID] = cloneRouteRule(r)
	return nil
}

func (s *MemoryStore) DeleteRouteRule(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.routeRules[id]; !ok {
		return ErrNotFound
	}
	delete(s.routeRules, id)
	return nil
}
