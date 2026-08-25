package model

import (
	"strings"
	"time"
)

const (
	BackendStatusUp       = "up"
	BackendStatusDown     = "down"
	BackendStatusDraining = "draining"
)

const (
	AlgoRoundRobin         = "round_robin"
	AlgoRandom             = "random"
	AlgoLeastConn          = "least_conn"
	AlgoWeightedRoundRobin = "weighted_round_robin"
)

var validAlgorithms = map[string]bool{
	AlgoRoundRobin:         true,
	AlgoRandom:             true,
	AlgoLeastConn:          true,
	AlgoWeightedRoundRobin: true,
}

func IsValidAlgorithm(a string) bool {
	return validAlgorithms[a]
}

// Backend 表示负载均衡后端节点。
type Backend struct {
	ID                  string    `json:"id"`
	Name                string    `json:"name"`
	Host                string    `json:"host"`
	Port                int       `json:"port"`
	Weight              int       `json:"weight"`
	Status              string    `json:"status"`
	Healthy             bool      `json:"healthy"`
	ConsecutiveFailures int       `json:"consecutive_failures"`
	ActiveConnections   int       `json:"active_connections"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (b *Backend) Validate() error {
	b.Name = strings.TrimSpace(b.Name)
	b.Host = strings.TrimSpace(b.Host)
	if b.Name == "" {
		return NewValidationError("name", "后端节点名称不能为空")
	}
	if b.Host == "" {
		return NewValidationError("host", "主机地址不能为空")
	}
	if b.Port <= 0 || b.Port > 65535 {
		return NewValidationError("port", "端口号必须在 1-65535 之间")
	}
	if b.Weight < 0 {
		return NewValidationError("weight", "权重不能为负数")
	}
	if b.Status == "" {
		b.Status = BackendStatusUp
	}
	if b.Status != BackendStatusUp && b.Status != BackendStatusDown && b.Status != BackendStatusDraining {
		return NewValidationError("status", "后端节点状态不合法")
	}
	return nil
}

var backendTransitions = map[string]map[string]bool{
	BackendStatusUp:       {BackendStatusDown: true, BackendStatusDraining: true},
	BackendStatusDown:     {BackendStatusUp: true},
	BackendStatusDraining: {BackendStatusDown: true, BackendStatusUp: true},
}

// CanTransitionBackend 校验后端节点状态流转是否合法。
func CanTransitionBackend(from, to string) bool {
	if m, ok := backendTransitions[from]; ok {
		return m[to]
	}
	return false
}

type BackendFilter struct {
	Status  string
	Keyword string
}

func (f BackendFilter) Match(b *Backend) bool {
	if f.Status != "" && b.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(b.Name), k) &&
			!strings.Contains(strings.ToLower(b.Host), k) {
			return false
		}
	}
	return true
}
