package main

import (
	"encoding/json"
	"fmt"
	"gitlab.com/firestart/ignition/pkg/config"
)

type Config struct {
	Database struct {
		Host     string
		Port     int
		Username string
	}
}

func main() {
	var cfg Config
	settings, _ := config.LoadSettings()
	_ = settings.Unpack(&cfg)

	//fmt.Println(settings.Get("Greeting:Hello"))

	marshal, _ := json.Marshal(cfg)
	fmt.Println(string(marshal))

	//cfg := Config{
	//	Database: struct{}{Host: "", Port: 0, Username: ""},
	//}
	//cfg, err := config.LoadConfig(cfg)
	//if err != nil {
	//	panic(err)
	//}
	//

	//
	//_ = cfg
}
