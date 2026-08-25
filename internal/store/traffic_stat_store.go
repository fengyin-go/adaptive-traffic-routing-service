package store

import (
	"loadbalancer/internal/model"
)

func (s *MemoryStore) CreateTrafficStat(t *model.TrafficStat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if existID, ok := s.backendStatID[t.BackendID]; ok {
		if exist, ok2 := s.trafficStats[existID]; ok2 {
			exist.RequestCount += t.RequestCount
			exist.ErrorCount += t.ErrorCount
			exist.TotalDurationMs += t.TotalDurationMs
			if t.LastRequestAt.After(exist.LastRequestAt) {
				exist.LastRequestAt = t.LastRequestAt
			}
			exist.UpdatedAt = t.UpdatedAt
			return nil
		}
	}
	s.trafficStats[t.ID] = t
	s.backendStatID[t.BackendID] = t.ID
	return nil
}

func (s *MemoryStore) GetTrafficStat(id string) (*model.TrafficStat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.trafficStats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *MemoryStore) ListTrafficStats() []*model.TrafficStat {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.TrafficStat, 0, len(s.trafficStats))
	for _, t := range s.trafficStats {
		list = append(list, t)
	}
	return list
}

func (s *MemoryStore) UpdateTrafficStat(t *model.TrafficStat) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.trafficStats[t.ID]; !ok {
		return ErrNotFound
	}
	s.trafficStats[t.ID] = t
	return nil
}

func (s *MemoryStore) DeleteTrafficStat(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	t, ok := s.trafficStats[id]
	if !ok {
		return ErrNotFound
	}
	delete(s.trafficStats, id)
	delete(s.backendStatID, t.BackendID)
	return nil
}

func (s *MemoryStore) GetTrafficStatByBackend(backendID string) (*model.TrafficStat, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id, ok := s.backendStatID[backendID]
	if !ok {
		return nil, ErrNotFound
	}
	t, ok := s.trafficStats[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}
