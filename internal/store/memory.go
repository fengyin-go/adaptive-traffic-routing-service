package store

import (
	"sync"

	"loadbalancer/internal/model"
)

type MemoryStore struct {
	mu            sync.RWMutex
	backends      map[string]*model.Backend
	strategies    map[string]*model.Strategy
	routeRules    map[string]*model.RouteRule
	trafficStats  map[string]*model.TrafficStat
	backendStatID map[string]string // backendID -> trafficStat ID
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		backends:      make(map[string]*model.Backend),
		strategies:    make(map[string]*model.Strategy),
		routeRules:    make(map[string]*model.RouteRule),
		trafficStats:  make(map[string]*model.TrafficStat),
		backendStatID: make(map[string]string),
	}
}

var _ Store = (*MemoryStore)(nil)
