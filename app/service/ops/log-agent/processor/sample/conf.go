package sample

import (
	"fmt"
	"errors"
	"time"

	"go-common/app/service/ops/log-agent/conf/configcenter"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

const (
	logSample = "sample.toml"
)

type Config struct {
	Local        bool             `toml:"local"`
	SampleConfig map[string]int64 `toml:"sampleConfig"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return fmt.Errorf("Error can't be nil")
	}
	if c.SampleConfig == nil {
		c.SampleConfig = make(map[string]int64)
	}
	return nil
}

func DecodeConfig(md toml.MetaData, primValue toml.Primitive) (c interface{}, err error) {
	config := new(Config)
	if err = md.PrimitiveDecode(primValue, config); err != nil {
		return nil, err
	}

	// read config from config center
	if !config.Local {
		if err = config.readConfig(); err != nil {
			return nil, err
		}
		// watch update and reload config
		go func() {
			currentVersion := configcenter.Version
			for {
				if currentVersion != configcenter.Version {
					log.Info("sample config reload")
					if err := config.readConfig(); err != nil {
						log.Error("sample config reload error (%v", err)
					}
					currentVersion = configcenter.Version
				}
				time.Sleep(time.Second)
			}
		}()
	}

	return config, nil
}

func (c *Config) readConfig() (err error) {
	var (
		ok        bool
		value     string
		tmpSample map[string]int64
	)

	// sample config
	if value, ok = configcenter.Client.Value(logSample); !ok {
		return errors.New("failed to get sample.toml")
	}
	if _, err = toml.Decode(value, &tmpSample); err != nil {
		return err
	}
	c.SampleConfig = tmpSample
	return nil
}
