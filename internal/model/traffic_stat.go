package model

import (
	"time"
)

// TrafficStat 表示流量统计。
type TrafficStat struct {
	ID              string    `json:"id"`
	BackendID       string    `json:"backend_id"`
	RequestCount    int64     `json:"request_count"`
	ErrorCount      int64     `json:"error_count"`
	TotalDurationMs int64     `json:"total_duration_ms"`
	LastRequestAt   time.Time `json:"last_request_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (t *TrafficStat) Validate() error {
	if t.BackendID == "" {
		return NewValidationError("backend_id", "后端节点 ID 不能为空")
	}
	if t.RequestCount < 0 {
		return NewValidationError("request_count", "请求数不能为负数")
	}
	if t.ErrorCount < 0 {
		return NewValidationError("error_count", "错误数不能为负数")
	}
	if t.ErrorCount > t.RequestCount {
		return NewValidationError("error_count", "错误数不能大于请求数")
	}
	if t.TotalDurationMs < 0 {
		return NewValidationError("total_duration_ms", "总耗时不能为负数")
	}
	return nil
}

type TrafficStatFilter struct {
	BackendID string
}

func (f TrafficStatFilter) Match(t *TrafficStat) bool {
	if f.BackendID != "" && t.BackendID != f.BackendID {
		return false
	}
	return true
}
