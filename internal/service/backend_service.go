package service

import (
	"sort"
	"sync/atomic"
	"time"

	"loadbalancer/internal/model"
	"loadbalancer/pkg/idgen"
)

func (s *Service) CreateBackend(input model.Backend) (*model.Backend, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	b := &model.Backend{
		ID:                  idgen.Hex(),
		Name:                input.Name,
		Host:                input.Host,
		Port:                input.Port,
		Weight:              input.Weight,
		Status:              input.Status,
		Healthy:             true,
		ConsecutiveFailures: 0,
		ActiveConnections:   0,
		CreatedAt:           now,
		UpdatedAt:           now,
	}
	if b.Weight == 0 {
		b.Weight = 1
	}
	if err := s.store.CreateBackend(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) GetBackend(id string) (*model.Backend, error) {
	return s.store.GetBackend(id)
}

func (s *Service) ListBackends(filter model.BackendFilter, page, size int) ([]*model.Backend, int, error) {
	all := s.store.ListBackends()
	matched := make([]*model.Backend, 0, len(all))
	for _, b := range all {
		if filter.Match(b) {
			matched = append(matched, b)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Backend{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateBackend(id string, input model.Backend) (*model.Backend, error) {
	b, err := s.store.GetBackend(id)
	if err != nil {
		return nil, err
	}
	if input.Name != "" {
		b.Name = input.Name
	}
	if input.Host != "" {
		b.Host = input.Host
	}
	if input.Port > 0 {
		b.Port = input.Port
	}
	if input.Weight >= 0 {
		b.Weight = input.Weight
	}
	b.UpdatedAt = time.Now()
	if err := b.Validate(); err != nil {
		return nil, err
	}
	if err := s.store.UpdateBackend(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) DeleteBackend(id string) error {
	return s.store.DeleteBackend(id)
}

func (s *Service) TransitionBackendStatus(id string, toStatus string) (*model.Backend, error) {
	b, err := s.store.GetBackend(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionBackend(b.Status, toStatus) {
		return nil, model.NewValidationError("status", "非法的状态流转")
	}
	b.Status = toStatus
	if toStatus == model.BackendStatusUp {
		b.Healthy = true
		b.ConsecutiveFailures = 0
	}
	if toStatus == model.BackendStatusDown {
		b.Healthy = false
	}
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBackend(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) RecordHealthCheck(id string, success bool) (*model.Backend, error) {
	b, err := s.store.GetBackend(id)
	if err != nil {
		return nil, err
	}
	if success {
		if b.ConsecutiveFailures > 0 {
			b.ConsecutiveFailures = 0
		}
		if !b.Healthy && b.Status == model.BackendStatusDown {
			b.Status = model.BackendStatusUp
			b.Healthy = true
		}
	} else {
		b.ConsecutiveFailures++
		if b.ConsecutiveFailures >= 3 {
			b.Status = model.BackendStatusDown
			b.Healthy = false
		}
	}
	b.UpdatedAt = time.Now()
	if err := s.store.UpdateBackend(b); err != nil {
		return nil, err
	}
	return b, nil
}

func (s *Service) PickBackend(algo string, backends []*model.Backend) *model.Backend {
	available := make([]*model.Backend, 0, len(backends))
	for _, b := range backends {
		if b.Status == model.BackendStatusUp && b.Healthy {
			available = append(available, b)
		}
	}
	if len(available) == 0 {
		return nil
	}
	switch algo {
	case model.AlgoRoundRobin:
		return pickRoundRobin(available)
	case model.AlgoRandom:
		return pickRandom(available)
	case model.AlgoLeastConn:
		return pickLeastConn(available)
	case model.AlgoWeightedRoundRobin:
		return pickWeightedRoundRobin(available)
	default:
		return pickRoundRobin(available)
	}
}

var rrCounter uint64

func pickRoundRobin(backends []*model.Backend) *model.Backend {
	if len(backends) == 0 {
		return nil
	}
	n := atomic.AddUint64(&rrCounter, 1)
	return backends[int(n-1)%len(backends)]
}

func pickRandom(backends []*model.Backend) *model.Backend {
	if len(backends) == 0 {
		return nil
	}
	n := atomic.AddUint64(&rrCounter, 1)
	idx := int(n % uint64(len(backends)))
	return backends[idx]
}

func pickLeastConn(backends []*model.Backend) *model.Backend {
	if len(backends) == 0 {
		return nil
	}
	best := backends[0]
	for _, b := range backends[1:] {
		if b.ActiveConnections < best.ActiveConnections {
			best = b
		}
	}
	return best
}

var wrrMu uint64

func pickWeightedRoundRobin(backends []*model.Backend) *model.Backend {
	if len(backends) == 0 {
		return nil
	}
	totalWeight := 0
	for _, b := range backends {
		totalWeight += b.Weight
	}
	if totalWeight == 0 {
		return backends[0]
	}
	n := atomic.AddUint64(&wrrMu, 1)
	idx := int((n - 1) % uint64(totalWeight))
	cum := 0
	for _, b := range backends {
		cum += b.Weight
		if idx < cum {
			return b
		}
	}
	return backends[len(backends)-1]
}
