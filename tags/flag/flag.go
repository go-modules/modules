// Package flag provides a tags.ValueSetter to set values from command line flags.
// Depends on normal use of the stdlib flag package.
// String flag values are passed through literal.ValueSetter.
package flag

import (
	stdFlag "flag"
	"github.com/go-modules/modules/tags/literal"
	"reflect"
)

var ValueSetter = &valueSetter{stdFlag.CommandLine}

// A valueSetter wraps a FlagSet and implements tags.ValueSetter
type valueSetter struct {
	*stdFlag.FlagSet
}

// Looks up the environment variable tagValue. Returned string is parsed and used to set value via literal.ValueSetter.
// Only sets value if the environment variable is set, otherwise passes by returning (false, nil).
func (v valueSetter) SetValue(value reflect.Value, tagValue string) (bool, error) {
	f := v.Lookup(tagValue)
	if f == nil || f.Value.String() == "" {
		return false, nil
	}
	return literal.ValueSetter.SetValue(value, f.Value.String())
}
