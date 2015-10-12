package modules

import (
	"reflect"
	"github.com/jmank88/go-modules/tags"
	"errors"
	"sync"
	"fmt"
)


// newBinding returns a new binding configured with binder
func newBinding(binder *Binder) *binding {
	return &binding{
		binder,
		fields{m: make(map[bindKey]reflect.Value)},
		gates{m: make(map[bindKey]gate)},
		newGate(),
		make(chan error),
	}
}

// A binding holds fields during a call to Binder.Bind
type binding struct {
	// Configuration.
	*Binder
	// The bound fields.
	fields
	// The provider/injector gates.
	gates
	// Closing the cancel gate signals binding goroutines to complete.
	cancel gate
	// Cancelled injection goroutines send errors here
	errors chan error
}

// Injects the value bound to key into fieldValue.
func (b *binding) inject(key bindKey, fieldValue reflect.Value) {
	// Wait to inject this field after it has been provided, or binding cancelled.
	select {
	case <-b.cancel:
		b.logf("nothing bound to %s\n", key.String())
	case <-b.gates.get(key):
		if bound, ok := b.fields.get(key); ok {
			fieldValue.Set(bound)
			b.logf("%v <- %s\n", bound, key.String())
		} else {
			b.logf("nothing bound to %s\n", key.String())
		}
	}
}

// provide binds value to key.
// Each recognized tag's tags.ValueSetFn will be executed until one sets the value.
func (b *binding) provide(key bindKey, singleton bool, tag tags.StructTag, value reflect.Value) error {
	// Range over tag fields until a know tag key's tags.ValueSetFn sets the value.
	// Note: We can't detect if value has already been set, so a mis-configured module could result in a tag overriding
	// a value set during Provide().
	tag.ForEach(tags.Handler(func(tagKey, v string) (bool, error) {
		if tagKey == "provide" {
			return false, nil
		}
		if tagKey == "inject" {
			return false, errors.New(fmt.Sprintf("failed to parse tags for value %s ;a module field tagged with 'provide' cannot also be tagged with 'inject'", key))
		}
		if valueSetter, ok := b.valueSetters[tagKey]; ok {
			if ok, err := valueSetter.SetValue(value, v); err != nil {
				// Failed to set value.
				return false, &AnnotatedError{msg: fmt.Sprintf("failed to provide value for %s from tag key %s", key, tagKey), cause: err}
			} else if ok {
				// Value has been set. Done.
				return true, nil
			} else {
				// Value has not been set. Continue.
				return false, nil
			}
		} else {
			// Unrecognized tag. Continue.
			return false, nil
		}
	}))

	// The value to bind.
	var toBind reflect.Value
	if singleton && value.Kind() == reflect.Func && !value.IsNil() {
		// Inject a singleton by wrapping the provided function.
		toBind = asSingleton(value)
		b.logf("singleton(%v) -> %v -> %s\n", value, toBind, key)
	} else {
		// Inject the value that was provided, as is.
		toBind = value
		b.logf("%v -> %s\n", value, key)
	}

	// Provide this field.
	b.fields.bind(key, toBind)

	// Broadcast to waiting injectors.
	close(b.gates.get(key))
	return nil
}

// A fields instance holds bound field values mapped by bindKeys.
type fields struct {
	sync.RWMutex
	m map[bindKey]reflect.Value
}

// The get method retrieves the value bound to key.
func (f *fields) get(key bindKey) (reflect.Value, bool) {
	f.RLock()
	value, ok := f.m[key]
	f.RUnlock()
	return value, ok
}

// Bind binds value to key.
func (f *fields) bind(key bindKey, value reflect.Value) {
	f.Lock()
	f.m[key] = value
	f.Unlock()
}

// A gate is a channel intended to be closed to broadcast a signal to receivers.
type gate chan struct{}

// The newGate function returns a new gate instance.
func newGate() gate {
	return gate(make(chan struct{}))
}

// A gates instance holds a lazily created singleton gate per bindKey.
type gates struct {
	sync.Mutex
	m map[bindKey]gate
}

// get returns the gate for key, or creates a new one if none exists.
// Safe to call from multiple goroutines.
func (g *gates) get(key bindKey) gate {
	g.Lock()
	gate, ok := g.m[key]
	if !ok {
		gate = newGate()
		g.m[key] = gate
	}
	g.Unlock()
	return gate
}