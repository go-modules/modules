package env

import (
	"testing"
	"os"
	"reflect"
)

func TestValueSetter(t *testing.T) {
	for _,testCase := range []struct {
		value reflect.Value
		tagValue string
		envVars map[string]string
		expected reflect.Value
	} {
		{
			reflect.New(reflect.TypeOf("")).Elem(),
			"envVarName",
			map[string]string{
				"envVarName": "value",
			},
			reflect.ValueOf("value"),
		},
		{
			reflect.New(reflect.TypeOf(0)).Elem(),
			"envVarName",
			map[string]string{
				"envVarName": "1234",
			},
			reflect.ValueOf(1234),
		},
	} {
		os.Clearenv()
		for k,v := range testCase.envVars {
			if err := os.Setenv(k, v); err != nil {
				t.Errorf("failed to set environment variable", err)
			}
		}
		if ok, err := valueSetterFunc(testCase.value, testCase.tagValue); err != nil {
			t.Error(err)
		} else if !ok {
			t.Error("expected value to be set")
		} else if testCase.value.Interface() != testCase.expected.Interface() {
			t.Errorf("expected %q but got %q", testCase.expected, testCase.value)
		}
	}
}
