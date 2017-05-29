// The modules package implements binders for dependency injection of tagged struct fields.
package modules

import (
	"errors"
	"fmt"
	"io"
	"log"
	"reflect"
	"sync"

	"github.com/go-modules/modules/inject"
	"github.com/go-modules/modules/inject/env"
	"github.com/go-modules/modules/inject/file"
	"github.com/go-modules/modules/inject/flag"
	"github.com/go-modules/modules/inject/literal"
	"github.com/go-modules/modules/tags"
)

// A Provider is a binding module that implements the Provide() method.
// When a Provider is bound, Provide() will be called (prior to having fields injected).
type Provider interface {
	// Provide is called once to set provided fields.
	// Returns nil for success, or an error in the case of failed binding.
	// This method is called prior to field injection, so injected fields may not be directly referenced, but may be closed over.
	Provide() error
}

// A Binder holds a configuration for module binding.
type Binder struct {
	// The logger (if present) will receive informational binding messages.
	logger *log.Logger
	// Injectors by tag key.
	injectors map[string]inject.Injector
}

// NewBinder initializes a new Binder instance, and applies options.
func NewBinder(options ...BinderOption) *Binder {
	b := &Binder{
		injectors: map[string]inject.Injector{
			"literal": literal.Injector,
			"env":     env.Injector,
			"flag":    flag.Injector,
			"file":    file.Injector,
		},
	}

	for _, option := range options {
		option.configure(b)
	}
	return b
}

// A functional option for configuring a Binder.
type BinderOption interface {
	configure(*Binder)
}

// Logger is a functional option that sets a Binder's logger.
type Logger struct {
	io.Writer
}

func (l Logger) configure(b *Binder) {
	b.logger = log.New(l, "modules: ", log.LstdFlags)
}

// Injectors is a functional option that adds mapped Injectors to a Binder.
type Injectors map[string]inject.Injector

func (v Injectors) configure(b *Binder) {
	for k, v := range v {
		b.injectors[k] = v
	}
}

// Bind binds modules. Calls Provide() on modules implementing Provider, calls
// inject.Injectors for tagged fields, and injects provided fields.
func (b *Binder) Bind(modules ...interface{}) error {
	binding := newBinding(b)
	// Holds errors during binding.
	errs := make([]error, 0)

	// Validate configuration
	if _, ok := binding.injectors["provide"]; ok {
		return errors.New("the 'provide' tag key may not be overridden")
	}
	if _, ok := binding.injectors["inject"]; ok {
		return errors.New("the 'inject' tag key may not be overridden")
	}

	// Collect errors in a goroutine.
	go func() {
		for err := range binding.errors {
			errs = append(errs, err)
		}
	}()

	// Injection goroutines signal here when complete.
	var injections sync.WaitGroup

	// Bind each module.
	for _, module := range modules {

		// If this module is a Provider then call Provide().
		if provider, ok := module.(Provider); ok {
			if err := provider.Provide(); err != nil {
				return &AnnotatedError{msg: "error during call to Provide()", cause: err}
			}
		}

		// Bind each field in this module.
		moduleType := reflect.TypeOf(module).Elem()
		for i := 0; i < moduleType.NumField(); i++ {
			field := moduleType.Field(i)
			value := reflect.ValueOf(module).Elem().Field(i)
			tag := tags.StructTag(string(field.Tag))
			if bindName, ok := tag.Get("inject"); ok {
				if !value.CanSet() {
					binding.errors <- fmt.Errorf("cannot inject unexported field: %s", field.Name)
					continue
				}
				injections.Add(1)
				go func() {
					// Blocks until a provider binds key, or cancelled.
					binding.inject(bindName, value)
					injections.Done()
				}()
			} else if tagValue, ok := tag.Get("provide"); ok {
				bindName, options := tags.ParseTag(tagValue)
				// Releases blocking injections for key.
				if err := binding.provide(bindName, options.Contains("singleton"), tag, value); err != nil {
					binding.errors <- err
				}
			}
		}
	}

	// Wait for all injection goroutines to complete.
	injections.Wait()

	// Signal error processing goroutine to complete.
	close(binding.errors)

	if len(errs) > 0 {
		return &BindingError{errs}
	}

	return nil
}

// logf logs to b's Logger, if present.
func (b *Binder) logf(fmt string, a ...interface{}) {
	if b.logger != nil {
		b.logger.Printf(fmt, a...)
	}
}

// Bound fields are keyed by type, and (optionally) name.
type bindKey struct {
	reflect.Type
	// Optional distinguishing name for types with multiple fields bound (e.g. string).
	name string
}

// Example: {string|database.host} or {DatabaseClient}
func (k *bindKey) String() string {
	if k.name != "" {
		return "{" + k.Type.String() + "|" + k.name + "}"
	}
	return "{" + k.Type.String() + "}"
}

// asSingleton wraps funcValue (which must of type.Kind Func, and take no parameters) in a function which only
// calls funcValue once, and caches the returned value for future calls.
func asSingleton(funcValue reflect.Value) reflect.Value {
	var once sync.Once
	var cache []reflect.Value
	return reflect.MakeFunc(funcValue.Type(), func(args []reflect.Value) []reflect.Value {
		once.Do(func() {
			cache = funcValue.Call([]reflect.Value{})
		})
		return cache
	})
}
