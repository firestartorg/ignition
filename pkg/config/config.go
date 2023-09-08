package config

import (
	"reflect"
	"strconv"
)

type Config struct {
	fields map[string]string
}

func NewConfig() *Config {
	return &Config{
		fields: make(map[string]string),
	}
}

func (cfg *Config) Get(key string) string {
	return cfg.fields[key]
}

func (cfg *Config) Lookup(key string) (string, bool) {
	val, ok := cfg.fields[key]
	return val, ok
}

func (cfg *Config) Set(key, value string) {
	cfg.fields[key] = value
}

func (cfg *Config) Unset(key string) {
	delete(cfg.fields, key)
}

func (cfg *Config) Keys() []string {
	keys := make([]string, len(cfg.fields))
	i := 0
	for key := range cfg.fields {
		keys[i] = key
		i++
	}
	return keys
}

func (cfg *Config) Merge(other *Config) {
	if other == nil {
		return
	}
	for key, val := range other.fields {
		cfg.fields[key] = val
	}
}

func (cfg *Config) Clone() *Config {
	clone := NewConfig()
	clone.Merge(cfg)
	return clone
}

// Unpack unpacks the configuration into the specified struct.
func (cfg *Config) Unpack(v interface{}) error {
	return unpack(v, nil, cfg)
}

//func Pack(v interface{}) (*Config, error) {
//	return nil, nil
//}

// Unpack unpacks the configuration into the specified struct.
func unpack(v interface{}, keySpace *string, cfg *Config) error {
	// Get field information
	fields := reflect.TypeOf(v).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		// Get field value
		value := reflect.ValueOf(v).Elem().Field(i)

		// Get field name
		name := field.Tag.Get("config")
		if name == "" {
			name = field.Name
		}

		// Get key
		key := name
		if keySpace != nil {
			key = *keySpace + ":" + key
		}

		// Get field type
		kind := field.Type.Kind()
		if kind == reflect.Struct {
			err := unpack(value.Addr().Interface(), &key, cfg)
			if err != nil {
				return err
			}
			continue
		}

		// Get field value
		val, ok := cfg.Lookup(key)
		if !ok {
			continue
		}

		// Set field value
		switch kind {
		case reflect.String:
			value.SetString(val)
		case reflect.Bool:
			boolVal, err := strconv.ParseBool(val)
			if err != nil {
				return err
			}
			value.SetBool(boolVal)
		case reflect.Float32:
			floatVal, err := strconv.ParseFloat(val, 32)
			if err != nil {
				return err
			}
			value.SetFloat(floatVal)
		case reflect.Float64:
			floatVal, err := strconv.ParseFloat(val, 64)
			if err != nil {
				return err
			}
			value.SetFloat(floatVal)
		case reflect.Uint:
			uintVal, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.Uint8:
			uintVal, err := strconv.ParseUint(val, 10, 8)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.Uint16:
			uintVal, err := strconv.ParseUint(val, 10, 16)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.Uint32:
			uintVal, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.Uint64:
			uintVal, err := strconv.ParseUint(val, 10, 64)
			if err != nil {
				return err
			}
			value.SetUint(uintVal)
		case reflect.Int:
			intVal, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			value.SetInt(intVal)
		case reflect.Int8:
			intVal, err := strconv.ParseInt(val, 10, 8)
			if err != nil {
				return err
			}
			value.SetInt(intVal)
		case reflect.Int16:
			intVal, err := strconv.ParseInt(val, 10, 16)
			if err != nil {
				return err
			}
			value.SetInt(intVal)
		case reflect.Int32:
			intVal, err := strconv.ParseInt(val, 10, 32)
			if err != nil {
				return err
			}
			value.SetInt(intVal)
		case reflect.Int64:
			intVal, err := strconv.ParseInt(val, 10, 64)
			if err != nil {
				return err
			}
			value.SetInt(intVal)
		}
	}

	return nil
}
