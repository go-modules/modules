package inject

import (
	"reflect"
	"unsafe"
)

type StringMaker interface {
	// Makes a string value base on the input string.
	// Returning false indicates no value was made.
	MakeString(string) (bool, string, error)
}

type BoolMaker interface {
	// Makes a bool value based on the input string.
	// Returning false indicates no value was made.
	MakeBool(string) (bool, bool, error)
}

type IntMaker interface {
	// Makes an int value with the given bitSize based on the input string.
	// Returning false indicates no value was made.
	MakeInt(s string, bitSize int) (bool, int64, error)
}

type UintMaker interface {
	// Makes a uint value with the given bitSize based on the input string.
	// Returning false indicates no value was made.
	MakeUint(s string, bitSize int) (bool, uint64, error)
}

type UintPtrMaker interface {
	// Makes a uintptr value with the given bitSize based on the input string.
	// Returning false indicates no value was made.
	MakeUintPtr(s string, bitSize int) (bool, uintptr, error)
}

type FloatMaker interface {
	// Makes a float value with the given bitSize based on the input string.
	// Returning false indicates no value was made.
	MakeFloat(s string, bitSize int) (bool, float64, error)
}

type ComplexMaker interface {
	// Makes a complex value with the given bitSize based on the input string.
	// Returning false indicates no value was made.
	MakeComplex(s string, bitSize int) (bool, complex128, error)
}

type ArrayMaker interface {
	// Makes an array value of the given type based on the input string.
	// Returning false indicates no value was made.
	MakeArray(string, reflect.Type) (bool, reflect.Value, error)
}

type ChanMaker interface {
	// Makes a chan value of the given type based on the input string.
	// Returning false indicates no value was made.
	MakeChan(string, reflect.Type) (bool, reflect.Value, error)
}

type FuncMaker interface {
	// Makes a func value of the given type based on the input string.
	// Returning false indicates no value was made.
	MakeFunc(string, reflect.Type) (bool, reflect.Value, error)
}

type InterfaceMaker interface {
	// Makes an interface value based on the input string.
	// Returning false indicates no value was made.
	MakeInterface(string) (bool, interface{}, error)
}

type MapMaker interface {
	// Makes a map value based on the input string.
	// Returning false indicates no value was made.
	MakeMap(string) (bool, reflect.Value, error)
}

type PtrMaker interface {
	// Makes a pointer value based on the input string.
	// Returning false indicates no value was made.
	MakePtr(string) (bool, reflect.Value, error)
}

type SliceMaker interface {
	// Makes a slice value of the given type based on the input string.
	// Returning false indicates no value was made.
	MakeSlice(string, reflect.Type) (bool, reflect.Value, error)
}

type StructMaker interface {
	// Makes a struct value of the given type based on the input string.
	// Returning false indicates no value was made.
	MakeStruct(string, reflect.Type) (bool, reflect.Value, error)
}

type UnsafePointerMaker interface {
	// Makes an unsafe pointer value based on the input string.
	// Returning false indicates no value was made.
	MakeUnsafePointer(string) (bool, unsafe.Pointer, error)
}
