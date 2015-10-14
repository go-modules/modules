// Package env provides a tags.ValueSetter to set values from environment variables.
// Environment variable strings are parsed by literal.ValueSetter.
package env

import (
	"reflect"
	"os"
	"github.com/go-modules/modules/tags"
	"github.com/go-modules/modules/tags/literal"
)

var ValueSetter = tags.ValueSetterFunc(valueSetterFunc)

// Looks up the environment variable tagValue. Returned string is parsed and used to set value via literal.ValueSetter.
// Only sets value if the environment variable is set, otherwise passes by returning (false, nil).
func valueSetterFunc(value reflect.Value, tagValue string) (bool, error) {
	envValue, ok := os.LookupEnv(tagValue)
	if !ok {
		return false, nil
	}
	return literal.ValueSetter.SetValue(value, envValue)
}