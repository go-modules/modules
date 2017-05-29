package tags

import "testing"

func TestForEach(t *testing.T) {
	var handled bool
	handler := Handler(func(k, v string) (bool, error) {
		if k == "handle" {
			handled = true
			return true, nil
		} else {
			return false, nil
		}
	})

	for _, testTag := range []StructTag{
		`handle:"value"`,
		`skip:"value" handle:"value"`,
		`skip:"value" handle:"value" skip2:"value"`,
	} {
		handled = false
		if err := testTag.ForEach(handler); err != nil {
			t.Errorf("failed to handle tag: %s", err)
		} else if !handled {
			t.Error("expected tag to be handled")
		}
	}

	for _, testTag := range []StructTag{
		`skip:"value"`,
		`skip:"value" skip2:"value"`,
		`skip:"value" skip2:"value" skip3:"value"`,
	} {
		handled = false
		if err := testTag.ForEach(handler); err != nil {
			t.Errorf("failed to handle tag: %s", err)
		} else if handled {
			t.Error("expected tag to be skipped")
		}
	}
}

func TestGet(t *testing.T) {
	for _, testCase := range []StructTag{
		`key1:"value1"`,
		`key2:"value2" key1:"value1"`,
		`key2:"value2" key1:"value1" key3:"value3"`,
	} {
		if value, ok := testCase.Get("key1"); !ok {
			t.Error("expected to find key1")
		} else if value != "value1" {
			t.Errorf(`expected "value1" got %q`, value)
		}
	}
}
