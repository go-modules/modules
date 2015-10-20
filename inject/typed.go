package inject

import "reflect"

// TypedInjector returns an Injector utilizing valueMaker, which should implement one or more *Maker interfaces.
func TypedInjector(valueMaker interface{}) Injector {
	return &typedInjector{valueMaker}
}

// A typedInjector adapts valueMaker to Injector.
type typedInjector struct {
	// valueMaker should implement one or more *Maker interfaces.
	valueMaker interface{}
}

// Injects value based on tagValue, if the valueMaker supports value's Kind.
// Implements Injector.
func (tvs typedInjector) Inject(value reflect.Value, tagValue string) (bool, error) {
	kind := value.Kind()
	switch kind {
	case reflect.String:
		stringValueMaker, ok := tvs.valueMaker.(StringMaker)
		if !ok {
			return false, &UnsupportedKindError{reflect.String}
		}
		set, s, err := stringValueMaker.MakeString(tagValue)
		if err != nil {
			return false, err
		} else if set {
			value.SetString(s)
		}
		return set, err
	case reflect.Bool:
		boolValueMaker, ok := tvs.valueMaker.(BoolMaker)
		if !ok {
			return false, &UnsupportedKindError{reflect.Bool}
		}
		set, b, err := boolValueMaker.MakeBool(tagValue)
		if err != nil {
			return false, err
		} else if set {
			value.SetBool(b)
		}
		return set, nil
	case reflect.Int:
		return intSetter(0, tvs.valueMaker, value, tagValue)
	case reflect.Int8:
		return intSetter(8, tvs.valueMaker, value, tagValue)
	case reflect.Int16:
		return intSetter(16, tvs.valueMaker, value, tagValue)
	case reflect.Int32:
		return intSetter(32, tvs.valueMaker, value, tagValue)
	case reflect.Int64:
		return intSetter(64, tvs.valueMaker, value, tagValue)
	case reflect.Uint:
		return uintSetter(0, tvs.valueMaker, value, tagValue)
	case reflect.Uint8:
		return uintSetter(8, tvs.valueMaker, value, tagValue)
	case reflect.Uint16:
		return uintSetter(16, tvs.valueMaker, value, tagValue)
	case reflect.Uint32:
		return uintSetter(32, tvs.valueMaker, value, tagValue)
	case reflect.Uint64:
		return uintSetter(64, tvs.valueMaker, value, tagValue)
	case reflect.Float32:
		return floatSetter(32, tvs.valueMaker, value, tagValue)
	case reflect.Float64:
		return floatSetter(64, tvs.valueMaker, value, tagValue)
	case reflect.Complex64:
		return complexSetter(64, tvs.valueMaker, value, tagValue)
	case reflect.Complex128:
		return complexSetter(128, tvs.valueMaker, value, tagValue)
	case reflect.Slice:
		sliceValueMaker, ok := tvs.valueMaker.(SliceMaker)
		if !ok {
			return false, &UnsupportedKindError{reflect.Slice}
		}
		set, made, err := sliceValueMaker.MakeSlice(tagValue, value.Type())
		if err != nil {
			return false, err
		} else if set {
			value.Set(made)
		}
		return set, nil
	case reflect.Array:
		arrayValueMaker, ok := tvs.valueMaker.(ArrayMaker)
		if !ok {
			return false, &UnsupportedKindError{reflect.Array}
		}
		set, made, err := arrayValueMaker.MakeArray(tagValue, value.Type())
		if err != nil {
			return false, err
		} else if set {
			value.Set(made)
		}
		return set, nil
	case reflect.Chan:
		chanValueMaker, ok := tvs.valueMaker.(ChanMaker)
		if !ok {
			return false, &UnsupportedKindError{reflect.Chan}
		}
		set, made, err := chanValueMaker.MakeChan(tagValue, reflect.TypeOf(value))
		if err != nil {
			return false, err
		} else if set {
			value.Set(made)
		}
		return set, nil
	default:
		return false, &UnsupportedKindError{kind}
	}
}

// intSetter sets value with the int returned from valueMaker, if it implements IntMaker. Otherwise it returns an error.
func intSetter(bits int, valueMaker interface{}, value reflect.Value, tagValue string) (bool, error) {
	intMaker, ok := valueMaker.(IntMaker)
	if !ok {
		return false, &UnsupportedKindError{reflect.Int}
	}
	set, i64, err := intMaker.MakeInt(tagValue, bits)
	if err != nil {
		return false, err
	} else if set {
		value.SetInt(i64)
	}
	return set, nil
}

// uintSetter sets value with the int returned from valueMaker, if it implements UintMaker. Otherwise it returns an error.
func uintSetter(bits int, valueMaker interface{}, value reflect.Value, tagValue string) (bool, error) {
	uintMaker, ok := valueMaker.(UintMaker)
	if !ok {
		return false, &UnsupportedKindError{reflect.Uint}
	}
	set, u64, err := uintMaker.MakeUint(tagValue, bits)
	if err != nil {
		return false, err
	} else if set {
		value.SetUint(u64)
	}
	return set, nil
}

// floatSetter sets value with the float returned from valueMaker, if it implements FloatMaker. Otherwise it returns an error.
func floatSetter(bitSize int, valueMaker interface{}, value reflect.Value, tagValue string) (bool, error) {
	floatMaker, ok := valueMaker.(FloatMaker)
	if !ok {
		var kind reflect.Kind
		if bitSize == 32 {
			kind = reflect.Float32
		} else {
			kind = reflect.Float64
		}
		return false, &UnsupportedKindError{kind}
	}
	set, f64, err := floatMaker.MakeFloat(tagValue, bitSize)
	if err != nil {
		return false, err
	} else if set {
		value.SetFloat(f64)
	}
	return set, nil
}

// complexSetter sets value with the complex returned from valueMaker, if it implements ComplexMaker. Otherwise it returns an error.
func complexSetter(bitSize int, valueMaker interface{}, value reflect.Value, tagValue string) (bool, error) {
	complexValueMaker, ok := valueMaker.(ComplexMaker)
	if !ok {
		var kind reflect.Kind
		if bitSize == 64 {
			kind = reflect.Complex64
		} else {
			kind = reflect.Complex128
		}
		return false, &UnsupportedKindError{kind}
	}
	set, c128, err := complexValueMaker.MakeComplex(tagValue, bitSize)
	if err != nil {
		return false, err
	} else if set {
		value.SetComplex(c128)
	}
	return set, nil
}

// An UnsupportedKindError indicates that a value maker does not support a certain reflect.Kind.
type UnsupportedKindError struct {
	reflect.Kind
}

func (e *UnsupportedKindError) Error() string {
	return "value maker does not support kind: " + e.Kind.String()
}
