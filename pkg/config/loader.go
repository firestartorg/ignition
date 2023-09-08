package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

var SupportedFileExtensions = [3]string{".yaml", ".yml", ".json"}

var NotFoundError = errors.New("config file not found")
var UnsupportedFormatError = errors.New("unsupported config file format")

// GetConfigLocation returns the path to the directory containing the configuration files.
func GetConfigLocation() (path string, err error) {
	var ok bool
	path, ok = os.LookupEnv("IGNITION_CONFIG_PATH")
	if !ok {
		path, err = os.Executable()
		if err != nil {
			return
		}
		path = filepath.Dir(path)
	}
	return
}

// findConfigFile searches for a configuration file in the specified path and name.
func findConfigFile(path, name string) (string, error) {
	file := filepath.Join(path, name)

	for _, ext := range SupportedFileExtensions {
		fileExt := file + ext

		info, err := os.Stat(fileExt)
		if err != nil {
			continue
		}
		if info.IsDir() {
			continue
		}

		return fileExt, nil
	}

	return "", NotFoundError
}

// loadConfigFile loads a configuration file from the specified path and name.
func loadConfigFile(path, name string) (cfg *Config, err error) {
	file, err := findConfigFile(path, name)
	if errors.Is(err, NotFoundError) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if filepath.Ext(file) == ".json" {
		return loadFromJson(file)
	} else if filepath.Ext(file) == ".yaml" || filepath.Ext(file) == ".yml" {
		return loadFromYaml(file)
	}
	return nil, UnsupportedFormatError
}

// Load loads a configuration file from the specified path and returns a Config object.
func Load(name string) (cfg *Config, err error) {
	var path string
	path, err = GetConfigLocation()
	if err != nil {
		return
	}

	var cfg0, cfg1 *Config

	// Load default config
	cfg0, err = loadConfigFile(path, name)
	if err != nil {
		return
	}

	// Load config for current environment
	env, ok := os.LookupEnv("GO_ENV")
	if ok {
		cfg1, err = loadConfigFile(path, name+"."+strings.ToLower(env))
		if err != nil {
			return
		}
	}

	// Merge configs
	cfg = NewConfig()
	cfg.Merge(cfg0)
	cfg.Merge(cfg1)

	return
}

// LoadConfig loads the default configuration file and environment variables.
func LoadConfig() (cfg *Config, err error) {
	// Load default config
	cfg, err = Load("appsettings")
	if err != nil {
		return
	}

	// Load environment variables
	var cfg1 *Config
	cfg1, err = LoadEnviron("IGNITION")
	if err != nil {
		return
	}
	cfg.Merge(cfg1)

	return
}
