package service

import (
	"sort"
	"time"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/idgen"
)

func (s *Service) CreateRouteRule(input model.RouteRule) (*model.RouteRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStrategy(input.StrategyID); err != nil {
		return nil, model.NewValidationError("strategy_id", "指定的策略不存在")
	}
	for _, bid := range input.BackendIDs {
		if _, err := s.store.GetBackend(bid); err != nil {
			return nil, model.NewValidationError("backend_ids", "后端节点不存在: "+bid)
		}
	}
	now := time.Now()
	r := &model.RouteRule{
		ID:         idgen.Hex(),
		PathPrefix: input.PathPrefix,
		StrategyID: input.StrategyID,
		BackendIDs: append([]string{}, input.BackendIDs...),
		Enabled:    input.Enabled,
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	if err := s.store.CreateRouteRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) GetRouteRule(id string) (*model.RouteRule, error) {
	return s.store.GetRouteRule(id)
}

func (s *Service) ListRouteRules(filter model.RouteRuleFilter, page, size int) ([]*model.RouteRule, int, error) {
	all := s.store.ListRouteRules()
	matched := make([]*model.RouteRule, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.RouteRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateRouteRule(id string, input model.RouteRule) (*model.RouteRule, error) {
	r, err := s.store.GetRouteRule(id)
	if err != nil {
		return nil, err
	}
	if input.PathPrefix != "" {
		r.PathPrefix = input.PathPrefix
	}
	if input.StrategyID != "" {
		if _, err := s.store.GetStrategy(input.StrategyID); err != nil {
			return nil, model.NewValidationError("strategy_id", "指定的策略不存在")
		}
		r.StrategyID = input.StrategyID
	}
	if len(input.BackendIDs) > 0 {
		for _, bid := range input.BackendIDs {
			if _, err := s.store.GetBackend(bid); err != nil {
				return nil, model.NewValidationError("backend_ids", "后端节点不存在: "+bid)
			}
		}
		r.BackendIDs = append([]string{}, input.BackendIDs...)
	}
	r.Enabled = input.Enabled
	r.UpdatedAt = time.Now()
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateRouteRule(r); err != nil {
		return nil, err
	}
	return r, nil
}

func (s *Service) DeleteRouteRule(id string) error {
	return s.store.DeleteRouteRule(id)
}
