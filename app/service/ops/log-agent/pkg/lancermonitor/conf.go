package lancermonitor

import (
	"errors"
)

type Config struct {
	Addr string `toml:"addr"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of LancerMonitor is nil")
	}

	if c.Addr == "" {
		return errors.New("addr of LancerMonitor can't be nil")
	}
	return nil
}
