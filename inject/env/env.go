// Package env provides a inject.Injector to set values from environment variables.
// Environment variable strings are parsed by literal.Injector.
package env

import (
	"os"
	"reflect"

	"github.com/go-modules/modules/inject"
	"github.com/go-modules/modules/inject/literal"
)

// Injector is an inject.Injector for parsing environment variables.
var Injector = inject.InjectorFunc(Inject)

// Inject looks up the environment variable by name, and sets the value via literal.Injector.
// Only sets value if the environment variable is set, otherwise passes by returning (false, nil).
func Inject(value reflect.Value, name string) (bool, error) {
	envValue, ok := os.LookupEnv(name)
	if !ok {
		return false, nil
	}
	return literal.Injector.Inject(value, envValue)
}
