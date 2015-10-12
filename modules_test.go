package modules

import (
	"testing"
)

// TestSimpleBind tests a one-way single-field binding.
// moduleA provides 'value' via Field.
// moduleB injects 'name' into Field.
func TestSimpleBind(t *testing.T) {
	moduleA := &struct {
		Field string `provide:""`
	}{
		Field: "value",
	}
	moduleB := &struct {
		Field string `inject:""`
	}{}

	if err := NewBinder().Bind(moduleA, moduleB); err != nil {
		t.Fatal(err)
	}

	assertString(t, "value", moduleA.Field)
	assertString(t, "value", moduleB.Field)
}

// TestSimpleBind tests a one-way, single-field, literal tag binding.
// moduleA provides 'name' as 'value' via Field.
// moduleB injects 'value' for 'name' into Field.
func TestSimpleTagBind(t *testing.T) {
	moduleA := &struct {
		Field string `provide:"name" literal:"value"`
	}{}
	moduleB := &struct {
		Field string `inject:"name"`
	}{}

	if err := NewBinder().Bind(moduleA, moduleB); err != nil {
		t.Fatal(err)
	}

	assertString(t, "value", moduleA.Field)
	assertString(t, "value", moduleB.Field)
}

// TestTwoWayBind tests a two-way, multi-field, literal tag binding.
// moduleA injects 'value1' for 'name1' into Field1 and provides 'name2' as 'value2' via Field2.
// moduleB provides 'name1' as 'value1' via Field1 and injects 'value2' for 'name2' into Field2.
func TestTwoWayBind(t *testing.T) {
	moduleA := &struct {
		Field1 string `inject:"name1"`
		Field2 string `provide:"name2" literal:"value2"`
	}{}
	moduleB := &struct {
		Field1 string `provide:"name1" literal:"value1"`
		Field2 string `inject:"name2"`
	}{}

	if err := NewBinder().Bind(moduleA, moduleB); err != nil {
		t.Fatal(err)
	}

	assertString(t, "value1", moduleA.Field1)
	assertString(t, "value2", moduleA.Field2)

	assertString(t, "value1", moduleB.Field1)
	assertString(t, "value2", moduleB.Field2)
}


// TestSingleton tests a simple singleton binding.
// Both moduleA and moduleB inject the singleton string function provided by moduleC.
func TestSingleton(t *testing.T) {
	moduleA := &struct {
		TestProvider func() string `inject:"test"`
	}{}
	moduleB := &struct {
		TestProvider func() string `inject:"test"`
	}{}
	moduleC := &struct {
		TestProvider func() string `provide:"test,singleton"`
	} {
		TestProvider: func() string {
			return "testValue"
		},
	}

	if err := NewBinder().Bind(moduleA, moduleB, moduleC); err != nil {
		t.Fatal(err)
	}

	assertNotNil(t, moduleA.TestProvider)
	assertString(t, "testValue", moduleA.TestProvider())
	assertNotNil(t, moduleB.TestProvider)
	assertString(t, "testValue", moduleB.TestProvider())
}

func assertNotNil(t *testing.T, value interface{}) {
	if value == nil {
		t.Errorf("expected non-nil value")
	}
}

func assertString(t *testing.T, expected, got string) {
	if expected != got {
		t.Errorf("expected %q got %q", expected, got)
	}
}
