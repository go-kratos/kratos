package classify

import (
	"fmt"
	"errors"
	"time"

	"go-common/app/service/ops/log-agent/conf/configcenter"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
)

const (
	logLevelMap = "logLevelMap.toml"
	logIdMap    = "logIdMap.toml"
)

type Config struct {
	Local             bool              `toml:"local"`
	LogLevelMapConfig map[string]string `toml:"logLevelMapConfig"`
	LogIdMapConfig    map[string]string `toml:"logIdMapConfig"`
	PriorityBlackList map[string]string `toml:"priorityBlackList"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return fmt.Errorf("Error can't be nil")
	}
	if c.LogLevelMapConfig == nil {
		c.LogLevelMapConfig = make(map[string]string)
	}

	if c.LogIdMapConfig == nil {
		return fmt.Errorf("LogIdMapConfig of classify can't be nil")
	}

	if c.PriorityBlackList == nil {
		c.PriorityBlackList = make(map[string]string)
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
					log.Info("classify config reload")
					if err := config.readConfig(); err != nil {
						log.Error("classify config reload error (%v)", err)
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
		ok             bool
		value          string
		tmplogLevelMap map[string]string
		tmplogIdMap    map[string]string
	)

	// logLevel config
	if value, ok = configcenter.Client.Value(logLevelMap); !ok {
		return errors.New("failed to get logLevelMap.toml")
	}
	if _, err = toml.Decode(value, &tmplogLevelMap); err != nil {
		return err
	}
	c.LogLevelMapConfig = tmplogLevelMap

	// logIdMap config
	if value, ok = configcenter.Client.Value(logIdMap); !ok {
		return errors.New("failed to get logIdMap.toml")
	}
	if _, err = toml.Decode(value, &tmplogIdMap); err != nil {
		return err
	}
	c.LogIdMapConfig = tmplogIdMap

	return nil
}
