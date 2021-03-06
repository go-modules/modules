// Package literal provides an inject.Injector that parses string literals into values.
//
// Supported Kinds: Bool, Int, Uint, Float, Complex, Chan; Slice and Array (of supported Kinds);
// Func (parameterless, single return, with string assignable/convertible to the return type);
// Interface, Struct (if string is assignable/convertible);
// Ptr, Uintptr, UnsafePointer
//
// Maps are not supported
package literal

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-modules/modules/inject"
)

// Injector is an inject.Injector for parsing string literals.
var Injector = inject.TypedInjector(&valueMaker{})

// valueMaker implements a subset of tags.*Maker interfaces.
type valueMaker struct{}

// MakeString passes str through as is.
func (valueMaker) MakeString(str string) (bool, string, error) {
	return true, str, nil
}

// MakeBool parses str into a bool value.
func (valueMaker) MakeBool(str string) (bool, bool, error) {
	val, err := strconv.ParseBool(str)
	return true, val, err
}

// MakeInt parses str into an int value.
func (valueMaker) MakeInt(str string, bitSize int) (bool, int64, error) {
	val, err := strconv.ParseInt(str, 10, bitSize)
	return true, val, err
}

// MakeUint parses str into a uint value.
func (valueMaker) MakeUint(str string, bitSize int) (bool, uint64, error) {
	val, err := strconv.ParseUint(str, 10, bitSize)
	return true, val, err
}

// MakeFloat parses str into a float value.
func (valueMaker) MakeFloat(str string, bitSize int) (bool, float64, error) {
	val, err := strconv.ParseFloat(str, bitSize)
	return true, val, err
}

// MakeComplex parses str into a complex value.
func (valueMaker) MakeComplex(str string, bits int) (bool, complex128, error) {
	values := strings.Split(str, ",")
	if values == nil || len(values) != 2 {
		return false, 0 + 0i, errors.New("illegal complex literal. expected 2 comma separated values")
	}
	var floatBitSize int
	if bits == 64 {
		floatBitSize = 32
	} else {
		floatBitSize = 64
	}
	real, err := strconv.ParseFloat(values[0], floatBitSize)
	if err != nil {
		return false, 0 + 0i, err
	}
	imaginary, err := strconv.ParseFloat(values[1], floatBitSize)
	if err != nil {
		return false, 0 + 0i, err
	}
	return true, complex(real, imaginary), nil
}

// MakeSlice parses str into a slice of typeOfElem.
// Returns a pointer to a slice populated with comma separated values parsed from str.
func (valueMaker) MakeSlice(str string, typeOfSlice reflect.Type) (bool, uintptr, error) {
	elements := strings.Split(str, ",")

	slice := reflect.MakeSlice(typeOfSlice, 0, len(elements))
	errs := make([]string, 0, 0)

	typeOfElem := typeOfSlice.Elem()
	for i, str := range elements {
		elem := reflect.New(typeOfElem).Elem()
		if _, err := Injector.Inject(elem, str); err != nil {
			errs = append(errs, fmt.Sprintf("element %d - %s", i, err.Error()))
		} else {
			// Note that if Inject did not set the value of elem then it will be zero valued.
			slice = reflect.Append(slice, elem)
		}
	}

	if len(errs) > 0 {
		return false, 0, errors.New(fmt.Sprintf("failed to parse list: %s", strings.Join(errs, "; ")))
	}

	return true, slice.Pointer(), nil
}

// MakeArray parses str into an array of typeOfElem.
// Returns a pointer to an array populated with comma separated values parsed from str.
func (valueMaker) MakeArray(str string, typeOfElem reflect.Type) (bool, uintptr, error) {
	elements := strings.Split(str, ",")

	array := reflect.New(reflect.ArrayOf(len(elements), typeOfElem))
	errs := make([]string, 0, 0)

	for i, str := range elements {
		elem := reflect.New(typeOfElem)
		if _, err := Injector.Inject(elem, str); err != nil {
			errs = append(errs, fmt.Sprintf("element %d - %s", i, err.Error()))
		} else {
			// Note that if Inject did not set the value of elem then it will be zero valued.
			array.Index(i).Set(elem)
		}
	}

	if len(errs) > 0 {
		return false, 0, errors.New(fmt.Sprintf("failed to parse list: %s", strings.Join(errs, "; ")))
	}
	return true, array.Pointer(), nil
}

// MakeChan parses str into a chan of typeOfElem.
// Returns a pointer to a channel of type typeOfElem with a buffer size parsed from str.
func (valueMaker) MakeChan(str string, typeOfElem reflect.Type) (bool, uintptr, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return false, 0, errors.New(fmt.Sprintf("failed to parse channel buffer capacity: %s", err.Error()))
	}

	c := reflect.MakeChan(typeOfElem, i)

	return true, c.Pointer(), nil
}

// MakeFunc parses str into a typeOfFn function.
// Returns a pointer to a function which always returns str.
func (valueMaker) MakeFunc(str string, typeOfFn reflect.Type) (bool, uintptr, error) {
	typeOfFnRet := typeOfFn.Out(0)
	ret := make([]reflect.Value, 1)
	if reflect.TypeOf(str).AssignableTo(typeOfFnRet) {
		ret[0] = reflect.ValueOf(str)
	} else if reflect.TypeOf(str).ConvertibleTo(typeOfFnRet) {
		ret[0] = reflect.ValueOf(str).Convert(typeOfFnRet)
	} else {
		return false, 0, errors.New(fmt.Sprintf("string is not assignable or convertible to function return type %s", typeOfFnRet))
	}
	f := reflect.MakeFunc(typeOfFn, func(args []reflect.Value) []reflect.Value {
		return ret
	})
	return true, f.Pointer(), nil
}

// MakeUintptr returns a uintPtr to str.
func (valueMaker) MakeUintptr(str string, base int, bitSize int) (bool, uintptr, error) {
	return true, reflect.ValueOf(str).Pointer(), nil
}
