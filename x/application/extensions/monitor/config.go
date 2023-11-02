package monitor

import (
	"strconv"
	"strings"
	"time"
)

type Config struct {
	// Either a port number or auto
	Port string
	// Timeout is the timeout for the readiness and liveness checks in seconds
	Timeout int
}

func (c Config) toOptions() []Option {
	var opts []Option
	if c.Port != "" {
		if strings.ToUpper(c.Port) == "AUTO" {
			opts = append(opts, WithPortAntiCollision())
		} else if port, err := strconv.ParseUint(c.Port, 10, 16); err == nil {
			opts = append(opts, WithPort(int(port)))
		}
	}
	if c.Timeout != 0 {
		opts = append(opts, WithTimeout(time.Duration(c.Timeout)*time.Second))
	}
	return opts
}

func (c Config) apply(o *Options) {
	opts := c.toOptions()
	for _, opt := range opts {
		opt(o)
	}
}
