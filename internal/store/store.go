// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"loadbalancer/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Backend
	CreateBackend(b *model.Backend) error
	GetBackend(id string) (*model.Backend, error)
	ListBackends() []*model.Backend
	UpdateBackend(b *model.Backend) error
	DeleteBackend(id string) error

	// Strategy
	CreateStrategy(s *model.Strategy) error
	GetStrategy(id string) (*model.Strategy, error)
	ListStrategies() []*model.Strategy
	UpdateStrategy(s *model.Strategy) error
	DeleteStrategy(id string) error

	// RouteRule
	CreateRouteRule(r *model.RouteRule) error
	GetRouteRule(id string) (*model.RouteRule, error)
	ListRouteRules() []*model.RouteRule
	UpdateRouteRule(r *model.RouteRule) error
	DeleteRouteRule(id string) error

	// TrafficStat
	CreateTrafficStat(t *model.TrafficStat) error
	GetTrafficStat(id string) (*model.TrafficStat, error)
	ListTrafficStats() []*model.TrafficStat
	UpdateTrafficStat(t *model.TrafficStat) error
	DeleteTrafficStat(id string) error
	GetTrafficStatByBackend(backendID string) (*model.TrafficStat, error)
}
