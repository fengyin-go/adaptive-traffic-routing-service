package requestscope

import (
	"context"
	"sync"
)

type contextKey string

const requestIDKey contextKey = "request-id"

func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, requestIDKey, id)
}

type Scope struct {
	RequestID string
	Labels    []string
}

type Snapshot struct {
	RequestID string
	Labels    []string
}

var scopePool = sync.Pool{New: func() any { return &Scope{} }}

func Acquire(ctx context.Context, labels []string) *Scope {
	scope := scopePool.Get().(*Scope)
	scope.RequestID, _ = ctx.Value(requestIDKey).(string)
	scope.Labels = append(scope.Labels[:0], labels...)
	return scope
}

func (s *Scope) Snapshot() Snapshot {
	return Snapshot{RequestID: s.RequestID, Labels: append([]string(nil), s.Labels...)}
}

func Release(scope *Scope) {
	scope.RequestID = ""
	clear(scope.Labels)
	scope.Labels = scope.Labels[:0]
	scopePool.Put(scope)
}
