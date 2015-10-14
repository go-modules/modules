package flag

import (
	"testing"
	"reflect"
	"flag"
)

func TestValueSetter(t *testing.T) {
	for _,testCase := range []struct {
		value reflect.Value
		tagValue string
		args []string
		expected reflect.Value
	} {
		{
			reflect.New(reflect.TypeOf("")).Elem(),
			"flagName",
			[]string{"-flagName", "value"},
			reflect.ValueOf("value"),
		},
		{
			reflect.New(reflect.TypeOf(0)).Elem(),
			"flagName",
			[]string{"-flagName", "1234"},
			reflect.ValueOf(1234),
		},
	} {
		flags := flag.NewFlagSet("test set", flag.ContinueOnError)
		flags.String(testCase.tagValue, "", "")
		if err := flags.Parse(testCase.args); err != nil {
			t.Fatalf("failed to set command line flags: ", err)
		}
		
		if ok, err := (&valueSetter{flags}).SetValue(testCase.value, testCase.tagValue); err != nil {
			t.Error(err)
		} else if !ok {
			t.Error("expected value to be set")
		} else if testCase.value.Interface() != testCase.expected.Interface() {
			t.Errorf("expected %q but got %q", testCase.expected, testCase.value)
		}
	}
}
