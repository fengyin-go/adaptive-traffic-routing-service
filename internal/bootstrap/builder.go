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
			runtime = nil
			err = fmt.Errorf("build runtime %q: %v", key, recovered)
		}
	}()
	backends, err := b.source.Load(key)
	if err != nil {
		return nil, err
	}
	next := &Runtime{Name: key, Backends: make(map[string]string, len(backends))}
	for id, address := range backends {
		if address == "panic" {
			panic("invalid backend address")
		}
		next.Backends[id] = address
	}
	next.Ready = true
	b.cache.Store(key, next)
	return cloneRuntime(next), nil
}
