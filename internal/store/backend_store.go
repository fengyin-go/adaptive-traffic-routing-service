package store

import (
	"loadbalancer/internal/model"
)

func cloneBackend(backend *model.Backend) *model.Backend {
	if backend == nil {
		return nil
	}
	copy := *backend
	return &copy
}

func (s *MemoryStore) CreateBackend(b *model.Backend) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.backends {
		if exist.Host == b.Host && exist.Port == b.Port {
			return ErrConflict
		}
	}
	s.backends[b.ID] = cloneBackend(b)
	return nil
}

func (s *MemoryStore) GetBackend(id string) (*model.Backend, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	b, ok := s.backends[id]
	if !ok {
		return nil, ErrNotFound
	}
	return cloneBackend(b), nil
}

func (s *MemoryStore) ListBackends() []*model.Backend {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Backend, 0, len(s.backends))
	for _, b := range s.backends {
		list = append(list, cloneBackend(b))
	}
	return list
}

func (s *MemoryStore) UpdateBackend(b *model.Backend) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.backends[b.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.backends {
		if exist.ID != b.ID && exist.Host == b.Host && exist.Port == b.Port {
			return ErrConflict
		}
	}
	s.backends[b.ID] = cloneBackend(b)
	return nil
}

func (s *MemoryStore) DeleteBackend(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.backends[id]; !ok {
		return ErrNotFound
	}
	delete(s.backends, id)
	return nil
}
