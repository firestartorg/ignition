package ignition

import "os"

func IsProduction() bool {
	env, ok := os.LookupEnv("GO_ENV")
	if ok && env == "production" {
		return true
	}
	return false
}
