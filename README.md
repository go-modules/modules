# Go Modules [![GoDoc](https://godoc.org/github.com/go-modules/modules?status.svg)](https://godoc.org/github.com/go-modules/modules) [![Build Status](https://travis-ci.org/go-modules/modules.svg)](https://travis-ci.org/go-modules/modules)
A dependency injection library using struct tags.

This project is currently alpha and subject to experimental change.

## Overview
This library simplifies the wiring of an application by injecting dependencies between modules.

A *module* is go struct containing fields tagged with 'inject' or 'provide' keys. When a set of modules are
*bound*, fields tagged 'inject' are set with values from corresponding fields tagged 'provide', respecting type and
(optionally) name. Provided fields may either be set normally prior to binding, during binding from a module's *Provide*
method, or from function calls tied to additional tag keys.

## How to Use

### Modules
A *module* is any tagged go struct. The 'inject' and 'provide' tag keys are treated specially. Other tag keys may be
registered with a Binder and trigger special processing. Unexported fields, and fields without recognized tags will be
ignored during binding.
```go
type simpleModule struct {
  // Requires a string value named 'injectMe' to be injected.
  FieldA string 'inject:"injectMe"'
  // Provides a string value named 'provideMe'.
  FieldB int 'provide:"provideMe"'
  // These fields are ignored by the Binder.
  FieldC bool
  fieldD string
}
```
When simpleModule is bound, it provides the string value named 'provideMe' via FieldB to the Binder, and expects the
string dependency named 'injectMe' to be provided by another module and injected into FieldA. FieldC and fieldD will be
ignored by the Binder.

### Providers
There are a few different ways for a module to provide values.

Fields may be set normally prior to binding.
```go
module := struct {
  FieldA string 'provide:"provideMe"'
} {
  FieldA: "providedValue"
}
```

Modules implementing the *Provider* interface may set fields from the *Provide* method.
```go
type module struct {
  FieldA string 'inject:"injectMe"'
  FieldB func() string 'provide:"provideMe"'
}
// Implements modules.Provider
func (m *Module) Provide() {
  // Note that injected fields have not yet necessarily been set at this point, so
  // they may not be accessed directly, but they may be closed over.
  m.FieldB = func() string {
    return = m.FieldA
  }
}
```

Additionally, a *Binder* may be configured to recognize certain tag keys and call a *ValueSetter* to set a value.
The 'literal' tag key is built-in, and parses string tag values into standard supported types.
```go
type module struct {
  FieldA string 'provide:"stringField" literal:"someString"'
  FieldB int    'provide:"intField" literal:"10"'
  FieldC complex128 'provide:"complexField" literal:"-1,1"'
}
```

### Binders
Modules are bound using a *Binder*. Binders are created with the *NewBinder* function, which optionally
accepts functional option arguments.
```go
binder := modules.NewBinder(modules.LogWriter(os.Stdout))
```
This binder logs information to stdout.

The *Bind* method binds a set of modules. All binding and injection occurs during this call. Modules implementing
*Provider* will have their *Provide* method called as well. Exported module fields are scanned for 'provide',
'inject' or other recognized tag keys.
```go
_ := binder.Bind(appModule, dataModule, serviceModule)
```
This call binds 3 modules. Each module's provided fields are available for injection into any module.

The functional option *ValueSetters* can be used to map custom or third party tags to *ValueSetter*s.
```go
valueSetters := modules.ValueSetters(map[string]ValueSetter{
  "customTag": customTag.ValueSetter,
})
binder := modules.NewBinder(valueSetters)
module := struct{
  FieldA CustomType 'provide:"someField" customTag:"tagValueArgument"'
}
_ := binder.Bind(module)
```
When this module is bound, *customTag.ValueSetter* may set the value of FieldA based on the tag value "tagValueArgument".