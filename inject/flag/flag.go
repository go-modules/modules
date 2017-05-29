// Package flag provides an inject.Injector to set values from command line flags.
// Depends on normal use of the stdlib flag package.
// String flag values are passed through literal.Injector.
package flag

import (
	stdFlag "flag"
	"reflect"

	"github.com/go-modules/modules/inject/literal"
)

// Injector is an inject.Injector for parsing command line flags.
var Injector = &injector{stdFlag.CommandLine}

// An injector wraps a FlagSet and implements inject.Injector
type injector struct {
	*stdFlag.FlagSet
}

// Inject looks up the flag by name and sets the value via literal.Injector.
// Only sets value if the flag is set, otherwise passes by returning (false, nil).
func (v injector) Inject(value reflect.Value, name string) (bool, error) {
	f := v.Lookup(name)
	if f == nil || f.Value.String() == "" {
		return false, nil
	}
	return literal.Injector.Inject(value, f.Value.String())
}
