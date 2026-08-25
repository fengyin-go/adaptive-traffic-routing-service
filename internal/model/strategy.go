package model

import (
	"strings"
	"time"
)

// Strategy 表示负载均衡策略。
type Strategy struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Algorithm string    `json:"algorithm"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Strategy) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	if s.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if s.Algorithm == "" {
		return NewValidationError("algorithm", "负载均衡算法不能为空")
	}
	if !IsValidAlgorithm(s.Algorithm) {
		return NewValidationError("algorithm", "不支持的负载均衡算法")
	}
	return nil
}

type StrategyFilter struct {
	Algorithm string
	Keyword   string
}

func (f StrategyFilter) Match(s *Strategy) bool {
	if f.Algorithm != "" && s.Algorithm != f.Algorithm {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) {
			return false
		}
	}
	return true
}
