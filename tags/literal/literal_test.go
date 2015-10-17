package literal

import (
	"testing"
)

var fixture = &valueMaker{}

func TestMakeString(t *testing.T) {
	for _, literal := range []string{"test", ""} {
		if ok, got, err := fixture.MakeString(literal); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected string %q to be made", literal)
		} else if got != literal {
			t.Errorf("expected %q got %q", literal, got)
		}
	}
}

func TestMakeBool(t *testing.T) {
	for _, testCase := range []struct {
		literal  string
		expected bool
	}{
		{"true", true},
		{"false", false},
	} {
		if ok, got, err := fixture.MakeBool(testCase.literal); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected bool %q to be made", testCase.expected)
		} else if got != testCase.expected {
			t.Errorf("expected %q got %q", testCase.expected, got)
		}
	}
}

func TestMakeInt(t *testing.T) {
	for _, testCase := range []struct {
		literal  string
		bitSize  int
		expected int64
	}{
		{"1", 0, 1},
		{"-127", 8, -127},
		{"32767", 16, 32767},
		{"-2147483647", 32, -2147483647},
		{"-9223372036854775807", 64, -9223372036854775807},
	} {
		if ok, got, err := fixture.MakeInt(testCase.literal, testCase.bitSize); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected int %d to be made", testCase.expected)
		} else if got != testCase.expected {
			t.Errorf("expected %d got %d", testCase.expected, got)
		}
	}
}

func TestMakeUint(t *testing.T) {
	for _, testCase := range []struct {
		literal  string
		bitSize  int
		expected uint64
	}{
		{"1", 0, 1},
		{"255", 8, 255},
		{"65535", 16, 65535},
		{"4294967295", 32, 4294967295},
		{"9223372036854775806", 64, 9223372036854775806},
	} {
		if ok, got, err := fixture.MakeUint(testCase.literal, testCase.bitSize); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected uint %d to be made", testCase.expected)
		} else if got != testCase.expected {
			t.Errorf("expected %d got %d", testCase.expected, got)
		}
	}
}

func TestMakeFloat(t *testing.T) {
	for _, testCase := range []struct {
		literal  string
		bitSize  int
		expected float64
	}{
		{"1", 32, 1},
		{"10.0", 32, 10.0},
		{"200.345", 64, 200.345},
		{".12345E+5", 64, .12345E+5},
		{"6.67428e-11", 64, 6.67428e-11},
	} {
		if ok, got, err := fixture.MakeFloat(testCase.literal, testCase.bitSize); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected float %f to be made", testCase.expected)
		} else if got != testCase.expected {
			t.Errorf("expected %f got %f", testCase.expected, got)
		}
	}
}

func TestMakeComplex(t *testing.T) {
	for _, testCase := range []struct {
		literal  string
		bitSize  int
		expected complex64
	}{
		{"1,1", 64, 1 + 1i},
		{"10.1,-6", 64, 10.1 + -6i},
	} {
		if ok, got, err := fixture.MakeComplex(testCase.literal, testCase.bitSize); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected complex %f to be made", testCase.expected)
		} else if complex64(got) != testCase.expected {
			t.Errorf("expected %f got %f", testCase.expected, got)
		}
	}

	for _, testCase := range []struct {
		literal  string
		bitSize  int
		expected complex128
	}{
		{"-200,-600", 128, -200 + -600i},
		{"1234,56", 128, 1234 + 56i},
		{"654321,123", 128, 654321 + 123i},
	} {
		if ok, got, err := fixture.MakeComplex(testCase.literal, testCase.bitSize); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !ok {
			t.Errorf("expected complex %f to be made", testCase.expected)
		} else if got != testCase.expected {
			t.Errorf("expected %f got %f", testCase.expected, got)
		}
	}
}
