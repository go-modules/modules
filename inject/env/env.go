// Package env provides a inject.Injector to set values from environment variables.
// Environment variable strings are parsed by literal.Injector.
package env

import (
	"os"
	"reflect"

	"github.com/go-modules/modules/inject"
	"github.com/go-modules/modules/inject/literal"
)

var Injector = inject.InjectorFunc(Inject)

// Looks up the environment variable tagValue. Returned string is parsed and used to set value via literal.Injector.
// Only sets value if the environment variable is set, otherwise passes by returning (false, nil).
func Inject(value reflect.Value, tagValue string) (bool, error) {
	envValue, ok := os.LookupEnv(tagValue)
	if !ok {
		return false, nil
	}
	return literal.Injector.Inject(value, envValue)
}
