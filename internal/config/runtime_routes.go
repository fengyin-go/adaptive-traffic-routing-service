package config

type RuntimeRoutes struct {
	Routes map[string]string
}

func NewRuntimeRoutes(initial map[string]string) *RuntimeRoutes {
	routes := &RuntimeRoutes{Routes: make(map[string]string, len(initial))}
	for path, backend := range initial {
		routes.Routes[path] = backend
	}
	return routes
}

func (r *RuntimeRoutes) Apply(path, backend string, validator RouteValidator) error {
	next := make(map[string]string, len(r.Routes)+1)
	for key, value := range r.Routes {
		next[key] = value
	}
	next[path] = backend
	if validator != nil {
		if err := validator.Validate(next); err != nil {
			return err
		}
	}
	r.Routes = next
	return nil
}
