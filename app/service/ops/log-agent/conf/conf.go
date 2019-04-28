package conf

import (
	"errors"
	"time"
	"go-common/app/service/ops/log-agent/pkg/flowmonitor"
	"go-common/library/log"
	"go-common/app/service/ops/log-agent/pkg/limit"
	"go-common/app/service/ops/log-agent/conf/configcenter"
	"go-common/app/service/ops/log-agent/pkg/httpstream"
	"go-common/app/service/ops/log-agent/pkg/lancermonitor"
	"go-common/app/service/ops/log-agent/pipeline/hostlogcollector"
	"go-common/app/service/ops/log-agent/pipeline/dockerlogcollector"
	"go-common/library/naming/discovery"

	"github.com/BurntSushi/toml"
)

const (
	config = "agent.toml"
)

var (
	// Conf conf
	Conf = &Config{}
)

type Config struct {
	// discovery
	Discovery *discovery.Config `toml:"discovery"`
	// log
	Log *log.Config `toml:"log"`
	// flow monitor
	Flowmonitor *flowmonitor.Config `toml:"flowmonitor"`
	// limit
	Limit *limit.LimitConf `toml:"limit"`
	// debug
	DebugAddr string `toml:"debugAddr"`
	// httpstream
	HttpStream *httpstream.Config `toml:"httpstream"`
	// lancermonitor
	LancerMonitor *lancermonitor.Config `toml:"lancermonitor"`
	// hostlogcollector
	HostLogCollector *hostlogcollector.Config `toml:"hostlogcollector"`
	// docker log collector
	DockerLogCollector *dockerlogcollector.Config `toml:"dockerLogCollector"`
}

func (c *Config) ConfigValidate() (error) {
	if c == nil {
		return errors.New("config of log agent can't be nil")
	}

	if c.DockerLogCollector == nil {
		c.DockerLogCollector = new(dockerlogcollector.Config)
	}

	if c.HostLogCollector == nil {
		c.HostLogCollector = new(hostlogcollector.Config)
	}

	return nil
}

// initConfig init config
func Init() (err error) {
	configcenter.InitConfigCenter()

	if err = readConfig(); err != nil {
		return
	}

	go func() {
		currentVersion := configcenter.Version
		for {
			if currentVersion != configcenter.Version {
				log.Info("lancer route config reload")
				if err := readConfig(); err != nil {
					log.Error("lancer route config reload error (%v", err)
				}
				currentVersion = configcenter.Version
			}
			time.Sleep(time.Second)
		}
	}()
	return Conf.ConfigValidate()
}

//// readConfig read config from config center
func readConfig() (err error) {
	var (
		ok        bool
		value     string
		tmpConfig *Config
	)
	//config
	if value, ok = configcenter.Client.Value(config); !ok {
		return errors.New("failed to get agent.toml")
	}
	if _, err = toml.Decode(value, &tmpConfig); err != nil {
		return err
	}
	Conf = tmpConfig
	return
}
