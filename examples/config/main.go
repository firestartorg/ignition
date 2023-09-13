package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/firestart/ignition/pkg/config"
)

type Enum string

const (
	EnumTest Enum = "test"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
	}

	Test     []string
	TestMap  map[string]string
	TestMap2 map[string]struct {
		Test string
	}
	TestEnum map[Enum]string
}

func main() {
	var cfg Config
	configuration, _ := config.LoadConfig()

	configuration.Set("TestMap:Test", "test")
	configuration.Set("TestMap2:Test:Test", "test")
	configuration.Set("TestEnum:test", "test")

	_ = configuration.Unpack(&cfg)

	fmt.Println(configuration.Get("Greeting:Hello"))

	marshal, _ := json.Marshal(cfg)
	fmt.Println(string(marshal))
}
