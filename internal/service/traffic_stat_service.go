package service

import (
	"sort"
	"time"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/idgen"
)

func (s *Service) CreateTrafficStat(input model.TrafficStat) (*model.TrafficStat, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetBackend(input.BackendID); err != nil {
		return nil, model.NewValidationError("backend_id", "指定的后端节点不存在")
	}
	now := time.Now()
	t := &model.TrafficStat{
		ID:              idgen.Hex(),
		BackendID:       input.BackendID,
		RequestCount:    input.RequestCount,
		ErrorCount:      input.ErrorCount,
		TotalDurationMs: input.TotalDurationMs,
		LastRequestAt:   input.LastRequestAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}
	if t.LastRequestAt.IsZero() {
		t.LastRequestAt = now
	}
	if err := s.store.CreateTrafficStat(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) GetTrafficStat(id string) (*model.TrafficStat, error) {
	return s.store.GetTrafficStat(id)
}

func (s *Service) ListTrafficStats(filter model.TrafficStatFilter, page, size int) ([]*model.TrafficStat, int, error) {
	all := s.store.ListTrafficStats()
	matched := make([]*model.TrafficStat, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].UpdatedAt.After(matched[j].UpdatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.TrafficStat{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateTrafficStat(id string, input model.TrafficStat) (*model.TrafficStat, error) {
	t, err := s.store.GetTrafficStat(id)
	if err != nil {
		return nil, err
	}
	if input.RequestCount >= 0 {
		t.RequestCount = input.RequestCount
	}
	if input.ErrorCount >= 0 {
		t.ErrorCount = input.ErrorCount
	}
	if input.TotalDurationMs >= 0 {
		t.TotalDurationMs = input.TotalDurationMs
	}
	if !input.LastRequestAt.IsZero() {
		t.LastRequestAt = input.LastRequestAt
	}
	t.UpdatedAt = time.Now()
	if err := t.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateTrafficStat(t); err != nil {
		return nil, err
	}
	return t, nil
}

func (s *Service) DeleteTrafficStat(id string) error {
	return s.store.DeleteTrafficStat(id)
}

func (s *Service) GetTrafficStatByBackend(backendID string) (*model.TrafficStat, error) {
	return s.store.GetTrafficStatByBackend(backendID)
}

// BackendSummary 返回单个后端节点的统计摘要。
func (s *Service) BackendSummary(backendID string) (map[string]interface{}, error) {
	if _, err := s.store.GetBackend(backendID); err != nil {
		return nil, err
	}
	stat, err := s.store.GetTrafficStatByBackend(backendID)
	if err != nil {
		return map[string]interface{}{
			"backend_id":      backendID,
			"request_count":   int64(0),
			"error_count":     int64(0),
			"error_rate":      float64(0),
			"avg_duration_ms": float64(0),
		}, nil
	}
	avg := float64(0)
	if stat.RequestCount > 0 {
		avg = float64(stat.TotalDurationMs) / float64(stat.RequestCount)
	}
	errRate := float64(0)
	if stat.RequestCount > 0 {
		errRate = float64(stat.ErrorCount) / float64(stat.RequestCount)
	}
	return map[string]interface{}{
		"backend_id":      backendID,
		"request_count":   stat.RequestCount,
		"error_count":     stat.ErrorCount,
		"error_rate":      errRate,
		"avg_duration_ms": avg,
	}, nil
}

// GlobalSummary 返回全局流量汇总。
func (s *Service) GlobalSummary() (map[string]interface{}, error) {
	all := s.store.ListTrafficStats()
	var totalReq, totalErr, totalDur int64
	for _, t := range all {
		totalReq += t.RequestCount
		totalErr += t.ErrorCount
		totalDur += t.TotalDurationMs
	}
	avg := float64(0)
	if totalReq > 0 {
		avg = float64(totalDur) / float64(totalReq)
	}
	errRate := float64(0)
	if totalReq > 0 {
		errRate = float64(totalErr) / float64(totalReq)
	}
	return map[string]interface{}{
		"total_request_count": totalReq,
		"total_error_count":   totalErr,
		"total_duration_ms":   totalDur,
		"avg_duration_ms":     avg,
		"global_error_rate":   errRate,
		"backend_count":       len(all),
	}, nil
}

// SummaryByStrategy 按策略分组统计流量。
func (s *Service) SummaryByStrategy() (map[string]interface{}, error) {
	rules := s.store.ListRouteRules()
	result := make(map[string]interface{})
	for _, r := range rules {
		if !r.Enabled {
			continue
		}
		var totalReq, totalErr, totalDur int64
		for _, bid := range r.BackendIDs {
			stat, err := s.store.GetTrafficStatByBackend(bid)
			if err != nil {
				continue
			}
			totalReq += stat.RequestCount
			totalErr += stat.ErrorCount
			totalDur += stat.TotalDurationMs
		}
		avg := float64(0)
		if totalReq > 0 {
			avg = float64(totalDur) / float64(totalReq)
		}
		errRate := float64(0)
		if totalReq > 0 {
			errRate = float64(totalErr) / float64(totalReq)
		}
		result[r.ID] = map[string]interface{}{
			"route_rule_id":       r.ID,
			"path_prefix":         r.PathPrefix,
			"total_request_count": totalReq,
			"total_error_count":   totalErr,
			"avg_duration_ms":     avg,
			"error_rate":          errRate,
		}
	}
	return result, nil
}
