package config

import (
	"os"
	"strings"
)

// LoadEnviron loads configuration from environment variables.
func LoadEnviron(prefix string) (cfg *Config, err error) {
	prefix = prefix + "_"
	cfg = NewConfig()

	kvs := os.Environ()
	for _, kv := range kvs {
		// Check if key starts with IGNITION_ prefix and ignore if not found (not a config key)
		if !strings.HasPrefix(kv, prefix) {
			continue
		}
		kv = strings.TrimPrefix(kv, prefix)

		// Split key and value
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}

		// Replace double underscore with colon in key
		//key := formatEnvironmentKey(parts[0])
		key := strings.ReplaceAll(parts[0], "__", ":")
		val := parts[1]

		cfg.fields[key] = val
	}

	return cfg, nil
}

//func toTitleCase(str string) string {
//	str = strings.ToLower(strings.TrimSpace(str))
//	parts := strings.Split(str, "_")
//	for i, part := range parts {
//		parts[i] = strings.Title(part)
//	}
//	return strings.Join(parts, "")
//}
//
//func formatEnvironmentKey(key string) string {
//	parts := strings.Split(key, "__")
//
//	for i, part := range parts {
//		parts[i] = toTitleCase(part)
//	}
//
//	return strings.Join(parts, ":")
//}
