package httpstream

import (
	"errors"
)

type Config struct {
	Addr string `toml:"addr"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of Sock Input is nil")
	}

	if c.Addr == "" {
		c.Addr = ":18123"
	}
	return nil
}
