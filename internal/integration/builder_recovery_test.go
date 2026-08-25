package integration

import (
	"testing"

	"loadbalancer/internal/bootstrap"
)

type panicSource struct{}

func (panicSource) Load(string) (map[string]string, error) {
	return map[string]string{"edge-a": "127.0.0.1:8080", "edge-b": "panic"}, nil
}

func TestRecoveredBuildNeverPublishesPartialRuntime(t *testing.T) {
	cache := bootstrap.NewCache()
	builder := bootstrap.NewBuilder(panicSource{}, cache)
	runtime, err := builder.Build("media")
	cached, found := cache.Load("media")
	if err == nil || runtime != nil || found || cached != nil {
		t.Fatalf("recovered runtime build exposed a partial configuration: err=%v runtime=%+v found=%v cached=%+v", err, runtime, found, cached)
	}
}
