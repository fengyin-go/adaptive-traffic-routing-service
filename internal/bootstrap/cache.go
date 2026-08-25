package bootstrap

import "sync"

type Runtime struct {
	Name     string
	Backends map[string]string
	Ready    bool
}

func cloneRuntime(runtime *Runtime) *Runtime {
	if runtime == nil {
		return nil
	}
	copy := &Runtime{Name: runtime.Name, Ready: runtime.Ready, Backends: make(map[string]string, len(runtime.Backends))}
	for id, address := range runtime.Backends {
		copy.Backends[id] = address
	}
	return copy
}

type Cache struct {
	mu    sync.Mutex
	items map[string]*Runtime
}

func NewCache() *Cache { return &Cache{items: make(map[string]*Runtime)} }

func (c *Cache) Store(key string, runtime *Runtime) bool {
	if runtime == nil {
		return false
	}
	c.mu.Lock()
	c.items[key] = cloneRuntime(runtime)
	c.mu.Unlock()
	return true
}

func (c *Cache) Load(key string) (*Runtime, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	runtime, ok := c.items[key]
	return cloneRuntime(runtime), ok
}
