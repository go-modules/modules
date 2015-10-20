// Package flag provides an inject.Injector to set values from command line flags.
// Depends on normal use of the stdlib flag package.
// String flag values are passed through literal.Injector.
package flag

import (
	stdFlag "flag"
	"reflect"

	"github.com/go-modules/modules/inject/literal"
)

var Injector = &injector{stdFlag.CommandLine}

// An injector wraps a FlagSet and implements inject.Injector
type injector struct {
	*stdFlag.FlagSet
}

// Looks up the environment variable tagValue. Returned string is parsed and used to set value via literal.Injector.
// Only sets value if the environment variable is set, otherwise passes by returning (false, nil).
func (v injector) Inject(value reflect.Value, tagValue string) (bool, error) {
	f := v.Lookup(tagValue)
	if f == nil || f.Value.String() == "" {
		return false, nil
	}
	return literal.Injector.Inject(value, f.Value.String())
}
