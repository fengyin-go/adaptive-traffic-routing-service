package config

type RuntimeRoutes struct {
	Routes map[string]string
}

func NewRuntimeRoutes(initial map[string]string) *RuntimeRoutes {
	routes := &RuntimeRoutes{}
	if len(initial) > 0 {
		routes.Routes = make(map[string]string, len(initial))
	}
	for path, backend := range initial {
		routes.Routes[path] = backend
	}
	return routes
}

func (r *RuntimeRoutes) Apply(path, backend string, validator RouteValidator) error {
	next := r.Routes
	next[path] = backend
	if validator != nil {
		if err := validator.Validate(next); err != nil {
			return err
		}
	}
	r.Routes = next
	return nil
}
