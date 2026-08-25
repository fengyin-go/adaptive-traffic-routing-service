package model

import (
	"strings"
	"time"
)

// RouteRule 表示路由规则。
type RouteRule struct {
	ID         string    `json:"id"`
	PathPrefix string    `json:"path_prefix"`
	StrategyID string    `json:"strategy_id"`
	BackendIDs []string  `json:"backend_ids"`
	Enabled    bool      `json:"enabled"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (r *RouteRule) Validate() error {
	r.PathPrefix = strings.TrimSpace(r.PathPrefix)
	if r.PathPrefix == "" {
		return NewValidationError("path_prefix", "路由路径前缀不能为空")
	}
	if !strings.HasPrefix(r.PathPrefix, "/") {
		return NewValidationError("path_prefix", "路由路径前缀必须以 / 开头")
	}
	if r.StrategyID == "" {
		return NewValidationError("strategy_id", "策略 ID 不能为空")
	}
	if len(r.BackendIDs) == 0 {
		return NewValidationError("backend_ids", "后端节点列表不能为空")
	}
	seen := make(map[string]bool)
	for _, id := range r.BackendIDs {
		if id == "" {
			return NewValidationError("backend_ids", "后端节点 ID 不能为空")
		}
		if seen[id] {
			return NewValidationError("backend_ids", "后端节点 ID 重复")
		}
		seen[id] = true
	}
	return nil
}

type RouteRuleFilter struct {
	PathPrefix string
	Keyword    string
}

func (f RouteRuleFilter) Match(r *RouteRule) bool {
	if f.PathPrefix != "" && !strings.HasPrefix(r.PathPrefix, f.PathPrefix) {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.PathPrefix), k) {
			return false
		}
	}
	return true
}
