// Package inject contains code for working with Injectors.
package inject

import "reflect"

// An Injector sets a value based on a string.
type Injector interface {
	// May set a value based on string.
	// Returns (true, nil) when a value has been set, or (false, nil) when a value has not been set (e.g. environment
	// variable not set, file not found, etc.).
	Inject(reflect.Value, string) (bool, error)
}

// InjectorFunc implements Injector.
type InjectorFunc func(reflect.Value, string) (bool, error)

// Inject implements the Injector interface.
func (f InjectorFunc) Inject(value reflect.Value, tagValue string) (bool, error) {
	return f(value, tagValue)
}
