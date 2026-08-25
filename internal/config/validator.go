package config

import "errors"

var ErrMissingRoute = errors.New("at least one default route is required")

type RouteValidator interface {
	Validate(map[string]string) error
}

type defaultRouteValidator struct{}

func (v *defaultRouteValidator) Validate(routes map[string]string) error {
	if v == nil || len(routes) == 0 {
		return ErrMissingRoute
	}
	return nil
}

func NewRouteValidator(enabled bool) RouteValidator {
	if !enabled {
		var validator *defaultRouteValidator
		return validator
	}
	return &defaultRouteValidator{}
}
