package tags

import (
	"reflect"
	"unsafe"
)

type StringMaker interface {
	MakeString(string) (bool, string, error)
}

type BoolMaker interface {
	MakeBool(string) (bool, bool, error)
}

type IntMaker interface {
	MakeInt(s string, bitSize int) (bool, int64, error)
}

type UintMaker interface {
	MakeUint(s string, bitSize int) (bool, uint64, error)
}

type UintPtrMaker interface {
	MakeUintPtr(s string, bitSize int) (bool, uintptr, error)
}

type FloatMaker interface {
	MakeFloat(s string, bitSize int) (bool, float64, error)
}

type ComplexMaker interface {
	MakeComplex(s string, bitSize int) (bool, complex128, error)
}

type ArrayMaker interface {
	MakeArray(string, reflect.Type) (bool, reflect.Value, error)
}

type ChanMaker interface {
	MakeChan(string, reflect.Type) (bool, reflect.Value, error)
}

type FuncMaker interface {
	MakeFunc(string, reflect.Type) (bool, reflect.Value, error)
}

type InterfaceMaker interface {
	MakeInterface(string) (bool, interface{}, error)
}

type MapMaker interface {
	MakeMap(string) (bool, reflect.Value, error)
}

type PtrMaker interface {
	MakePtr(string) (bool, reflect.Value, error)
}

type SliceMaker interface {
	MakeSlice(string, reflect.Type) (bool, reflect.Value, error)
}

type StructMaker interface {
	MakeStruct(string, reflect.Type) (bool, reflect.Value, error)
}

type UnsafePointerMaker interface {
	MakeUnsafePointer(string) (bool, unsafe.Pointer, error)
}