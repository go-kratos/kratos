package conf

import (
	"errors"
	"github.com/BurntSushi/toml"
	"go-common/library/conf"
	"go-common/library/log"
	xtime "go-common/library/time"
)

type BroadcastProxyConfig struct {
	Perf      string `toml:"perf"`
	Log       *log.Config
	Http      *HttpConfig
	Backend   *BackendConfig
	ZooKeeper *ZooKeeperConfig
	Ipip      *IpipConfig
	Dispatch  *DispatchConfig
	Sven      *SvenConfig
}

type HttpConfig struct {
	Address string
}

type BackendConfig struct {
	MaxIdleConnsPerHost int
	ProbePath           string
	BackendServer       []string
	ProbeSample         int
}

type ZooKeeperConfig struct {
	Address    []string
	Timeout    xtime.Duration
	ConfigPath string
}

type SinaIPConfig struct {
	Data string
}

type IpipConfig struct {
	V4 string
	V6 string
}
type DispatchConfig struct {
	MaxLimit             int
	DefaultDomain        string
	WildcardDomainSuffix string
	FileName             string
}

type SvenConfig struct {
	TreeID string
	Zone   string
	Env    string
	Build  string
	Token  string
}

func NewBroadcastProxyConfig(file string) (*BroadcastProxyConfig, error) {
	config := new(BroadcastProxyConfig)
	if file != "" {
		if err := config.local(file); err != nil {
			return nil, err
		}
	} else {
		if err := config.remote(); err != nil {
			return nil, err
		}
	}
	return config, nil
}

func (config *BroadcastProxyConfig) local(filename string) (err error) {
	_, err = toml.DecodeFile(filename, config)
	return
}

func (config *BroadcastProxyConfig) remote() error {
	client, err := conf.New()
	if err != nil {
		return err
	}
	if err = config.load(client); err != nil {
		return err
	}
	go func() {
		for range client.Event() {
			log.Info("config event")
		}
	}()
	return nil
}

func (config *BroadcastProxyConfig) load(c *conf.Client) error {
	s, ok := c.Value("live-broadcast-proxy.toml")
	if !ok {
		return errors.New("load config center error")
	}
	if _, err := toml.Decode(s, config); err != nil {
		return errors.New("could not decode config")
	}
	return nil
}
