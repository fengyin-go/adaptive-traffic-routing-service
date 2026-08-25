package service

import (
	"sort"
	"strings"

	"loadbalancer/internal/model"
)

// RouteResult 表示路由选择结果。
type RouteResult struct {
	RouteRule *model.RouteRule `json:"route_rule"`
	Strategy  *model.Strategy  `json:"strategy"`
	Backend   *model.Backend   `json:"backend"`
}

// SelectRoute 根据请求路径选择命中的路由规则、策略与后端节点。
func (s *Service) SelectRoute(path string) (*RouteResult, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, model.NewValidationError("path", "请求路径不能为空")
	}
	rules := s.store.ListRouteRules()
	var matched *model.RouteRule
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		if strings.HasPrefix(path, r.PathPrefix) {
			if matched == nil || len(r.PathPrefix) > len(matched.PathPrefix) {
				matched = r
			}
		}
	}
	if matched == nil {
		return nil, model.NewValidationError("path", "未匹配到任何路由规则")
	}

	strategy, err := s.store.GetStrategy(matched.StrategyID)
	if err != nil {
		return nil, err
	}
	if !strategy.Enabled {
		return nil, model.NewValidationError("strategy", "策略已禁用")
	}

	backends := make([]*model.Backend, 0, len(matched.BackendIDs))
	for _, bid := range matched.BackendIDs {
		b, err := s.store.GetBackend(bid)
		if err != nil {
			continue
		}
		backends = append(backends, b)
	}
	if len(backends) == 0 {
		return nil, model.NewValidationError("backends", "没有可用的后端节点")
	}

	picked := s.PickBackend(strategy.Algorithm, backends)
	if picked == nil {
		return nil, model.NewValidationError("backends", "没有健康的后端节点")
	}

	return &RouteResult{
		RouteRule: matched,
		Strategy:  strategy,
		Backend:   picked,
	}, nil
}

// RouteWithBackendID 用于流量统计的便捷方法：给定 path，返回选中的 backendID。
func (s *Service) RouteWithBackendID(path string) (string, error) {
	res, err := s.SelectRoute(path)
	if err != nil {
		return "", err
	}
	return res.Backend.ID, nil
}

// ListAvailableBackendsForPath 返回给定路径下所有可用的后端节点。
func (s *Service) ListAvailableBackendsForPath(path string) ([]*model.Backend, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil, model.NewValidationError("path", "请求路径不能为空")
	}
	rules := s.store.ListRouteRules()
	var matched *model.RouteRule
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		if strings.HasPrefix(path, r.PathPrefix) {
			if matched == nil || len(r.PathPrefix) > len(matched.PathPrefix) {
				matched = r
			}
		}
	}
	if matched == nil {
		return nil, model.NewValidationError("path", "未匹配到任何路由规则")
	}
	backends := make([]*model.Backend, 0, len(matched.BackendIDs))
	for _, bid := range matched.BackendIDs {
		b, err := s.store.GetBackend(bid)
		if err != nil {
			continue
		}
		if b.Status == model.BackendStatusUp && b.Healthy {
			backends = append(backends, b)
		}
	}
	sort.Slice(backends, func(i, j int) bool {
		return backends[i].Name < backends[j].Name
	})
	return backends, nil
}
