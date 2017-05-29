package modules

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/go-modules/modules/tags"
)

// newBinding returns a new binding configured with b
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

// Inject injectss the value bound to bindName into value.
func (b *binding) inject(bindName string, value reflect.Value) {
	key := bindKey{value.Type(), bindName}
	// Wait to inject this field after it has been provided, or binding cancelled.
	select {
	case <-b.cancel:
		b.logf("nothing bound to %s\n", key.String())
	case <-b.gates.get(key):
		if bound, ok := b.fields.get(key); ok {
			value.Set(bound)
			b.logf("%v <- %s\n", bound, key.String())
		} else {
			b.logf("nothing bound to %s\n", key.String())
		}
	}
}

// provide binds value to bindName.
// Each recognized tag key's inject.Injector will be executed until one sets the value.
func (b *binding) provide(bindName string, singleton bool, tag tags.StructTag, value reflect.Value) error {
	key := bindKey{value.Type(), bindName}
	// Range over tag fields until a known tag key's inject.Injector sets the value.
	tag.ForEach(tags.Handler(func(tagKey, v string) (bool, error) {
		if tagKey == "provide" {
			return false, nil
		}
		if tagKey == "inject" {
			return false, errors.New(fmt.Sprintf("failed to parse tags for value %s ;a module field tagged with 'provide' cannot also be tagged with 'inject'", key))
		}
		if injector, ok := b.injectors[tagKey]; ok {
			if ok, err := injector.Inject(value, v); err != nil {
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

// get retrieves the value bound to key.
func (f *fields) get(key bindKey) (reflect.Value, bool) {
	f.RLock()
	value, ok := f.m[key]
	f.RUnlock()
	return value, ok
}

// bind binds value to key.
func (f *fields) bind(key bindKey, value reflect.Value) {
	f.Lock()
	f.m[key] = value
	f.Unlock()
}

// A gate is a channel intended to be closed to broadcast a signal to receivers.
type gate chan struct{}

// newGate returns a new gate instance.
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
