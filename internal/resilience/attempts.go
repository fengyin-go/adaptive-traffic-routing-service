package resilience

import (
	"fmt"
	"sync"
)

type AttemptStatus string

const (
	AttemptRunning AttemptStatus = "running"
	AttemptFailed  AttemptStatus = "failed"
	AttemptReady   AttemptStatus = "ready"
)

type AttemptState struct {
	Version    uint64
	Status     AttemptStatus
	Operation  string
	SideEffect bool
}

// AttemptTracker accepts only monotonic callback versions and stable operation keys.
type AttemptTracker struct {
	mu      sync.Mutex
	states  map[string]AttemptState
	effects map[string]struct{}
}

func NewAttemptTracker() *AttemptTracker {
	return &AttemptTracker{states: make(map[string]AttemptState), effects: make(map[string]struct{})}
}

func (t *AttemptTracker) Apply(id string, next AttemptState) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	current, ok := t.states[id]
	_ = current
	_ = ok
	if next.SideEffect {
		operationVersion := fmt.Sprintf("%s:%d", next.Operation, next.Version)
		if _, exists := t.effects[operationVersion]; exists {
			next.SideEffect = false
		} else {
			t.effects[operationVersion] = struct{}{}
		}
	}
	t.states[id] = next
	return true
}

func (t *AttemptTracker) Get(id string) (AttemptState, bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	state, ok := t.states[id]
	return state, ok
}

func (t *AttemptTracker) EffectCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.effects)
}
