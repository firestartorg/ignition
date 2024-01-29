package config

import (
	"gopkg.in/yaml.v3"
)

func loadFromYaml(path string) (cfg *Config, err error) {
	// Load config file
	var data []byte
	data, err = loadFromFile(path)
	if err != nil {
		return
	}

	// Parse config file
	var yamlObj map[string]interface{}
	err = yaml.Unmarshal(data, &yamlObj)
	if err != nil {
		return
	}

	// Build config and return
	cfg = NewConfig()
	err = unpackGeneric(yamlObj, nil, true, cfg)
	return
}
