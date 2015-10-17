package modules

import (
	"fmt"
)

// The KVClient interface models a simple key/value store.
type KVClient interface {
	Get(key string) string
	Put(key, value string)
}

// A MapDBClient is a simple mock KVClient implementation backed by a map and configured with a default value for missing keys.
type MapDBClient struct {
	defaultValue string
	db           map[string]string
}

func (client *MapDBClient) Get(key string) string {
	if value, ok := client.db[key]; ok {
		return value
	} else {
		return client.defaultValue
	}
}

func (client *MapDBClient) Put(key, value string) {
	client.db[key] = value
}

// A service module has a 'GetData' service which utilizes an injected DBClient.
type ServiceModule struct {
	Client func() KVClient `inject:""`
}

func (service *ServiceModule) GetData(key string) string {
	return service.Client().Get(key)
}

func (service *ServiceModule) StoreData(key, value string) {
	service.Client().Put(key, value)
}

type defaultValue string

// This data module provides a Client function for retrieving a KVClient, which returns a DBClient configured with the
// injected default value.
type DataModule struct {
	DefaultValue defaultValue    `inject:""`
	Client       func() KVClient `provide:",singleton"`
}

func (data *DataModule) Provide() error {
	data.Client = func() KVClient {
		return &MapDBClient{defaultValue: string(data.DefaultValue), db: make(map[string]string)}
	}
	return nil
}

func Example() {
	serviceModule := &ServiceModule{}

	// This config module provides the default value required by the data module.
	configModule := &struct {
		DefaultValue defaultValue `provide:""`
	}{
		DefaultValue: "default",
	}

	binder := NewBinder()
	if err := binder.Bind(serviceModule, &DataModule{}, configModule); err != nil {
		panic(err)
	}

	fmt.Println(serviceModule.GetData("key"))

	serviceModule.StoreData("key", "value")
	fmt.Println(serviceModule.GetData("key"))

	// Output:
	// default
	// value
}
