package dockerlogcollector

import (
	"errors"
	"time"

	xtime "go-common/library/time"
)

type Config struct {
	ConfigEnv    string        `toml:"configEnv"`
	ConfigSuffix string        `toml:"configSuffix"`
	MetaPath     string        `toml:"metaPath"`
	ScanInterval xtime.Duration `toml:"scanInterval"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of docker log collector can't be nil")
	}

	if c.ConfigEnv == "" {
		c.ConfigEnv = "LogCollectorConf"
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
