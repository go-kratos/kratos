package hostlogcollector

import (
	"errors"
	"time"

	xtime "go-common/library/time"
)

type Config struct {
	HostConfigPath string        `toml:"hostConfigPath"`
	ConfigSuffix   string        `toml:"configSuffix"`
	MetaPath       string        `toml:"metaPath"`
	ScanInterval   xtime.Duration `toml:"scanInterval"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of host log collector can't be nil")
	}

	if c.HostConfigPath == "" {
		return errors.New("hostConfigPath of host log collector config can't be nil")
	}

	if c.MetaPath == "" {
		c.MetaPath = "/data/log-agent/meta"
	}

	if c.ConfigSuffix == "" {
		c.ConfigSuffix = ".conf"
	}

	if c.ScanInterval == 0 {
		c.ScanInterval = xtime.Duration(time.Second * 10)
	}

	return nil
}
