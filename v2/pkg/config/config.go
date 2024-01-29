package config

import (
	"reflect"
	"strconv"
	"strings"
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

func (cfg *Config) Sub(key string) *Config {
	key = key + ":"
	sub := NewConfig()
	for k, v := range cfg.fields {
		if len(k) > len(key) && k[:len(key)] == key {
			sub.fields[k[len(key):]] = v
		}
	}
	return sub
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
	// GetInjectable field information
	fields := reflect.TypeOf(v).Elem()
	for i := 0; i < fields.NumField(); i++ {
		field := fields.Field(i)

		// GetInjectable field value
		value := reflect.ValueOf(v).Elem().Field(i)

		// GetInjectable field name
		name := field.Tag.Get("config")
		if name == "" {
			name = field.Name
		}

		// GetInjectable key
		key := name
		if keySpace != nil {
			key = *keySpace + ":" + key
		}

		// GetInjectable field type
		kind := field.Type.Kind()
		if kind == reflect.Struct {
			err := unpack(value.Addr().Interface(), &key, cfg)
			if err != nil {
				return err
			}
			continue
		} else if kind == reflect.Map {
			cfg1 := cfg.Sub(key)
			set := map[string]bool{}

			// Determine all keys in the map
			for _, key := range cfg1.Keys() {
				key = strings.SplitN(key, ":", 2)[0]
				set[key] = true
			}

			// Create map
			value.Set(reflect.MakeMap(value.Type()))

			// SetInjectable map values
			for key := range set {
				// This is required to because of enum types
				keyValue := reflect.New(value.Type().Key()).Elem()
				keyValue.SetString(key)

				switch value.Type().Elem().Kind() {
				case reflect.String:
					value.SetMapIndex(keyValue, reflect.ValueOf(cfg1.Get(key)))
				case reflect.Struct:
					mapValue := reflect.New(value.Type().Elem()).Elem()

					err := unpack(mapValue.Addr().Interface(), &key, cfg1)
					if err != nil {
						return err
					}

					value.SetMapIndex(keyValue, mapValue)
				}
			}

			continue
		} else if kind == reflect.Array || kind == reflect.Slice {
			var counter int
			var keys []string
			for {
				key := key + ":" + strconv.Itoa(counter)

				_, ok := cfg.Lookup(key)
				if !ok {
					break
				}

				keys = append(keys, key)
				counter++
			}

			value.Set(reflect.MakeSlice(value.Type(), len(keys), len(keys)))

			for i, key := range keys {
				val := cfg.Get(key)

				value := value.Index(i)
				switch value.Kind() {
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
				case reflect.Struct:
					err := unpack(value.Addr().Interface(), &key, cfg)
					if err != nil {
						return err
					}
				}
			}

			continue
		}

		// GetInjectable field value
		val, ok := cfg.Lookup(key)
		if !ok {
			continue
		}

		// SetInjectable field value
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
