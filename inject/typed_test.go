package inject

import (
	"reflect"
	"testing"
)

func TestInject(t *testing.T) {
	testCases := []struct {
		value    reflect.Value
		expected reflect.Value
	}{
		{
			value:    reflect.New(reflect.TypeOf("")).Elem(),
			expected: reflect.ValueOf(testString),
		},
		{
			value:    reflect.New(reflect.TypeOf(false)).Elem(),
			expected: reflect.ValueOf(testBool),
		},
		{
			value:    reflect.New(reflect.TypeOf(int(0))).Elem(),
			expected: reflect.ValueOf(testInt),
		},
		{
			value:    reflect.New(reflect.TypeOf(int8(0))).Elem(),
			expected: reflect.ValueOf(testInt8),
		},
		{
			value:    reflect.New(reflect.TypeOf(int16(0))).Elem(),
			expected: reflect.ValueOf(testInt16),
		},
		{
			value:    reflect.New(reflect.TypeOf(int32(0))).Elem(),
			expected: reflect.ValueOf(testInt32),
		},
		{
			value:    reflect.New(reflect.TypeOf(int64(0))).Elem(),
			expected: reflect.ValueOf(testInt64),
		},
		{
			value:    reflect.New(reflect.TypeOf(uint(0))).Elem(),
			expected: reflect.ValueOf(testUint),
		},
		{
			value:    reflect.New(reflect.TypeOf(uint8(0))).Elem(),
			expected: reflect.ValueOf(testUint8),
		},
		{
			value:    reflect.New(reflect.TypeOf(uint16(0))).Elem(),
			expected: reflect.ValueOf(testUint16),
		},
		{
			value:    reflect.New(reflect.TypeOf(uint32(0))).Elem(),
			expected: reflect.ValueOf(testUint32),
		},
		{
			value:    reflect.New(reflect.TypeOf(uint64(0))).Elem(),
			expected: reflect.ValueOf(testUint64),
		},
		{
			value:    reflect.New(reflect.TypeOf(float32(0))).Elem(),
			expected: reflect.ValueOf(testFloat32),
		},
		{
			value:    reflect.New(reflect.TypeOf(float64(0))).Elem(),
			expected: reflect.ValueOf(testFloat64),
		},
		{
			value:    reflect.New(reflect.TypeOf(complex64(0 + 0i))).Elem(),
			expected: reflect.ValueOf(testComplex64),
		},
		{
			value:    reflect.New(reflect.TypeOf(complex128(0 + 0i))).Elem(),
			expected: reflect.ValueOf(testComplex128),
		},
		{
			value:    reflect.New(reflect.TypeOf([]string{})).Elem(),
			expected: reflect.ValueOf(testSlice),
		},
		{
			value:    reflect.New(reflect.TypeOf([3]int{})).Elem(),
			expected: reflect.ValueOf(testArray),
		},
		{
			value:    reflect.New(reflect.TypeOf(make(chan int))).Elem(),
			expected: reflect.ValueOf(testChan),
		},
	}
	for _, testCase := range testCases {
		if b, err := constantInjector.Inject(testCase.value, ""); err != nil {
			t.Errorf("unexpected error: %s", err)
		} else if !b {
			t.Error("expected value to be set")
		} else if !reflect.DeepEqual(testCase.expected.Interface(), testCase.value.Interface()) {
			t.Errorf("expected %v but got %v", testCase.expected, testCase.value)
		}
	}

	for _, testCase := range testCases {
		if _, err := errInjector.Inject(testCase.value, ""); err == nil {
			t.Error("expected error")
		} else {
			if _, ok := err.(*UnsupportedKindError); !ok {
				t.Errorf("unexpected error type: %s", err)
			}
		}
	}
}

var errInjector = TypedInjector(&struct{}{})

const (
	testString     = "test"
	testBool       = true
	testInt        = int(10)
	testInt8       = int8(100)
	testInt16      = int16(500)
	testInt32      = int32(-2000)
	testInt64      = int64(-40000)
	testUint       = uint(50)
	testUint8      = uint8(200)
	testUint16     = uint16(2500)
	testUint32     = uint32(30000)
	testUint64     = uint64(50000)
	testFloat32    = float32(10.056)
	testFloat64    = float64(-5367.0235)
	testComplex64  = complex64(25 + 68i)
	testComplex128 = complex128(-649 - 2393i)
)

var (
	testSlice = []string{"elem1", "elem2"}
	testArray = [3]int{10, 56, 100}
	testChan  = make(chan int)
)

var constantInjector = TypedInjector(&constantMaker{})

type constantMaker struct{}

func (*constantMaker) MakeString(tagValue string) (bool, string, error) {
	return true, testString, nil
}

func (*constantMaker) MakeBool(tagValue string) (bool, bool, error) {
	return true, testBool, nil
}

func (*constantMaker) MakeInt(tagValue string, bitSize int) (bool, int64, error) {
	var ret int64
	switch bitSize {
	case 0:
		ret = int64(testInt)
	case 8:
		ret = int64(testInt8)
	case 16:
		ret = int64(testInt16)
	case 32:
		ret = int64(testInt32)
	case 64:
		ret = testInt64
	}
	return true, ret, nil
}

func (*constantMaker) MakeUint(tagValue string, bitSize int) (bool, uint64, error) {
	var ret uint64
	switch bitSize {
	case 0:
		ret = uint64(testUint)
	case 8:
		ret = uint64(testUint8)
	case 16:
		ret = uint64(testUint16)
	case 32:
		ret = uint64(testUint32)
	case 64:
		ret = testUint64
	}
	return true, ret, nil
}

func (*constantMaker) MakeFloat(tagValue string, bitSize int) (bool, float64, error) {
	var ret float64
	switch bitSize {
	case 32:
		ret = float64(testFloat32)
	case 64:
		ret = testFloat64
	}
	return true, ret, nil
}

func (*constantMaker) MakeComplex(tagValue string, bitSize int) (bool, complex128, error) {
	var ret complex128
	switch bitSize {
	case 64:
		ret = complex128(testComplex64)
	case 128:
		ret = testComplex128
	}
	return true, ret, nil
}

func (*constantMaker) MakeSlice(tagValue string, typeOfSlice reflect.Type) (bool, reflect.Value, error) {
	return true, reflect.ValueOf(testSlice), nil
}

func (*constantMaker) MakeArray(tagValue string, typeOfArray reflect.Type) (bool, reflect.Value, error) {
	return true, reflect.ValueOf(testArray), nil
}

func (*constantMaker) MakeChan(tagValue string, typeOfChan reflect.Type) (bool, reflect.Value, error) {
	return true, reflect.ValueOf(testChan), nil
}
