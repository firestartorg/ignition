package main

import (
	"fmt"
	"gitlab.com/firestart/ignition/pkg/injector"
)

type Config struct {
	// Name is the name of the application
	Name string
	// Version is the version of the application
	Version string
}

type MyService struct {
	// Config is the configuration of the application
	config Config
}

func NewMyService(inj *injector.Injector) (srv MyService, err error) {
	var config Config
	config, err = injector.Get[Config](inj)
	if err != nil {
		return
	}

	srv = MyService{
		config: config,
	}
	return
}

func main() {
	inj := injector.NewInjector()
	injector.MustInjectConfig[Config](inj)
	injector.Provide(inj, NewMyService)

	// Do something with the config
	srv := injector.MustGet[MyService](inj)
	fmt.Println(srv)
}
