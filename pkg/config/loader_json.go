package config

import (
	"encoding/json"
	"gitlab.com/firestart/ignition"
	"io"
	"os"
	"strconv"
	"strings"
)

func loadFromFile(path string) (data []byte, err error) {
	var file *os.File

	// Open file for reading
	file, err = os.Open(path)
	if err != nil {
		return
	}
	// Close file handle
	defer file.Close()

	// Read file contents into memory buffer
	data, err = io.ReadAll(file)
	return
}

func loadFromJson(path string) (cfg *Config, err error) {
	// Load config file
	var data []byte
	data, err = loadFromFile(path)
	if err != nil {
		return
	}

	// Parse config file
	var jsonObj map[string]interface{}
	err = json.Unmarshal(data, &jsonObj)
	if err != nil {
		return
	}

	// Build config and return
	cfg = NewConfig()
	err = unpackGeneric(jsonObj, nil, false, cfg)
	return
}

func unpackGeneric(jsonObj map[string]interface{}, keySpace *string, cap bool, config *Config) error {
	for key, val := range jsonObj {
		if cap {
			parts := strings.Split(key, "-")
			for i, p := range parts {
				parts[i] = ignition.CapitalizeString(p)
			}
			key = strings.Join(parts, "")
		}
		// Check if key is a nested object
		if keySpace != nil {
			key = *keySpace + ":" + key
		}

		switch val.(type) {
		case map[string]interface{}:
			err := unpackGeneric(val.(map[string]interface{}), &key, cap, config)
			if err != nil {
				return err
			}
		case string:
			config.Set(key, val.(string))
		case float64:
			config.Set(key, strconv.FormatFloat(val.(float64), 'f', -1, 64))
		case bool:
			config.Set(key, strconv.FormatBool(val.(bool)))
		case int:
			config.Set(key, strconv.FormatInt(int64(val.(int)), 10))
		case uint:
			config.Set(key, strconv.FormatUint(uint64(val.(uint)), 10))
		case []interface{}:
			m := map[string]interface{}{}
			for i, v := range val.([]interface{}) {
				m[strconv.Itoa(i)] = v
			}
			err := unpackGeneric(m, &key, cap, config)
			if err != nil {
				return err
			}
		default:
			config.Set(key, val.(string))
		}
	}
	return nil
}
