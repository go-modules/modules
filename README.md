# Go Modules [![GoDoc](https://godoc.org/github.com/go-modules/modules?status.svg)](https://godoc.org/github.com/go-modules/modules) [![Build Status](https://travis-ci.org/go-modules/modules.svg)](https://travis-ci.org/go-modules/modules)
A dependency injection library using struct tags.

- [Slack channel](https://gophers.slack.com/messages/go-modules/)

## Overview
This library simplifies the wiring of an application by injecting dependencies between modules.

A *module* is a go struct containing fields tagged with 'inject' or 'provide' keys. When a set of modules are
*bound*, fields tagged 'inject' are set with values from corresponding fields tagged 'provide', respecting type and
 name. Provided fields may be set normally prior to binding, during binding from a module's *Provide* method, or from
 *Injector*s triggered by additional tag keys.

## How to Use

### Modules
A *module* is any tagged struct. The 'inject' and 'provide' tag keys are treated specially. Other tag keys may be
registered with a *Binder* to trigger *Injector*s. Unexported fields, and fields without recognized tags will be
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
When simpleModule is bound, it provides the string value named 'provideMe' via FieldB to the *Binder*, and expects the
string dependency named 'injectMe' to be provided by another module and injected into FieldA. FieldC and fieldD will be
ignored by the *Binder*.


### Providers
There are a few different ways for a *module* to provide values.

Fields may be set normally prior to binding.
```go
module := struct {
  FieldA string 'provide:"provideMe"'
} {
  FieldA: "value"
}
// or
module.FieldA = "value"
```

Modules implementing the *Provider* interface may set fields from the *Provide* method.
```go
type module struct {
  Field string 'inject:"injectMe"'
  Func func() string 'provide:"provideMe"'
}
// Implements modules.Provider
func (m *Module) Provide() error {
  m.Func = func() string {
    return m.Field
  }
  return nil
}
```
The *Provide* method is called during binding. Injected fields have not yet necessarily been set when *Provide* is
called, so they may not be accessed directly, but they may be closed over.

Additionally, a *Binder* may be configured to recognize certain tag keys and call an *Injector* to set a value.
The 'literal' tag key is built-in, and parses string tag values into standard supported types.
```go
type module struct {
  FieldA string 'provide:"stringField" literal:"someString"'
  FieldB int    'provide:"intField" literal:"10"'
  FieldC complex128 'provide:"complexField" literal:"-1,1"'
}
```
Other built-in tag keys include:
- 'env' for environment variables
- 'file' for os.File handles, and decoding of txt, json, xml, and gob
- 'flag' for command line arguments

### Binders
Modules are bound using a *Binder*. Binders are created with the *NewBinder* function, which optionally
accepts functional option arguments.
```go
binder := modules.NewBinder(modules.LogWriter(os.Stdout))
```
This binder logs information to stdout.

The *Bind* method binds a set of modules. All binding and injection occurs during this call. Modules implementing
*Provider* will have their *Provide* method called. Exported fields are inspected for 'provide', 'inject' or other
registered tag keys.
```go
_ := binder.Bind(appModule, dataModule, serviceModule)
```
This call binds 3 modules. Each module's provided fields are available for injection into any module.


### Tags and Injectors
The functional option *Injectors* can be used to map tag keys (anything besides "provide" and "inject") to custom or
third party *Injector*s.
```go
injectors := modules.Injectors(map[string]Injector{
  "customTag": customInjector,
})
binder := modules.NewBinder(injectors)
module := struct{
  FieldA CustomType 'provide:"someField" customTag:"tagValueArgument"'
}
_ := binder.Bind(module)
```
When this module is bound, *customInjector* may set the value of FieldA based on the tag value "tagValueArgument".

The *Injector* interface is defined in the inject package.
```go
// An Injector sets a value based on a string.
type Injector interface {
	// May set a value based on string.
	// Returns (true, nil) when a value has been set, or (false, nil) when a value has not been set (e.g. environment
	// variable not set, file not found, etc.).
	Inject(reflect.Value, string) (bool, error)
}
```

If a field is tagged with multiple keys, *Inject* will be called for each *Injector* until one sets the value.
```go
module := struct{
  Field string `provide:"setting" flag:"setting" env:"SETTING" literal:"defaultValue"`
}
```
This module provides a string value named 'setting', which may be set via a command-line flag or environment variable,
and which falls back to the default literal 'defaultValue'.

See the [GoDoc](https://godoc.org/github.com/go-modules/modules) for more api documentation, and a working example.