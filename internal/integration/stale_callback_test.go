package integration

import (
	"testing"

	"loadbalancer/internal/model"
	"loadbalancer/internal/resilience"
	"loadbalancer/internal/routing"
)

func TestLateAttemptCannotRollBackReadyRoutingState(t *testing.T) {
	tracker := resilience.NewAttemptTracker()
	if !tracker.Apply("route-a", resilience.AttemptState{Version: 2, Status: resilience.AttemptReady, Operation: "publish-a", SideEffect: true}) {
		panic("ready attempt was rejected")
	}
	staleAccepted := tracker.Apply("route-a", resilience.AttemptState{Version: 1, Status: resilience.AttemptRunning, Operation: "publish-a", SideEffect: true})
	state, _ := tracker.Get("route-a")
	table := routing.NewTable()
	table.Publish(routing.Snapshot{Version: 2, Backends: map[string]*model.Backend{"edge-a": {ID: "edge-a"}}})
	oldSnapshotAccepted := table.Publish(routing.Snapshot{Version: 1, Backends: map[string]*model.Backend{"edge-b": {ID: "edge-b"}}})
	current := table.Current()
	if staleAccepted || oldSnapshotAccepted || state.Status != resilience.AttemptReady || state.Version != 2 || tracker.EffectCount() != 1 || current.Version != 2 || current.Backends["edge-a"] == nil {
		t.Fatalf("late attempt rolled routing state backward or repeated its publish effect: state=%+v effects=%d snapshot=%+v", state, tracker.EffectCount(), current)
	}
}
