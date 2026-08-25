package bootstrap

import "fmt"

type Source interface {
	Load(string) (map[string]string, error)
}

type Builder struct {
	source Source
	cache  *Cache
}

func NewBuilder(source Source, cache *Cache) *Builder {
	return &Builder{source: source, cache: cache}
}

func (b *Builder) Build(key string) (runtime *Runtime, err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			// A failed build must never surface or cache a partial runtime:
			// discard the half-built value so callers only observe the error,
			// and never seed the cache with an incomplete configuration that a
			// later lookup would read as if it were valid.
			runtime = nil
			err = fmt.Errorf("build runtime %q: %v", key, recovered)
		}
	}()
	backends, err := b.source.Load(key)
	if err != nil {
		return nil, err
	}
	runtime = &Runtime{Name: key, Backends: make(map[string]string, len(backends))}
	for id, address := range backends {
		if address == "panic" {
			panic("invalid backend address")
		}
		runtime.Backends[id] = address
	}
	runtime.Ready = true
	b.cache.Store(key, runtime)
	return cloneRuntime(runtime), nil
}
