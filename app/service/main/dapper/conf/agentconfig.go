package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"go-common/library/conf"
	"go-common/library/log"
)

const (
	_agentConfigKey = "dapper-agent.toml"
)

// LoadAgentConfig LoadAgentConfig
func LoadAgentConfig() (*AgentConfig, error) {
	if confPath != "" {
		cfg := new(AgentConfig)
		_, err := toml.DecodeFile(confPath, cfg)
		return cfg, err
	}
	return remoteAgentConfig()
}

// AgentConfig config for dapper agent
type AgentConfig struct {
	Servers    []string         `toml:"servers"`
	Log        *log.Config      `toml:"log"`
	Queue      *QueueConfig     `toml:"queue"`
	UDPCollect UDPCollectConfig `toml:"udp_collect"`
}

// QueueConfig internal queue config
type QueueConfig struct {
	// queue local stroage path
	MemBuckets  int    `toml:"mem_buckets"`
	BucketBytes int    `toml:"bucket_bytes"`
	CacheDir    string `toml:"cache_dir"`
}

// UDPCollectConfig collect config
type UDPCollectConfig struct {
	Workers int    `toml:"workers"`
	Addr    string `toml:"addr"`
}

func remoteAgentConfig() (*AgentConfig, error) {
	client, err := conf.New()
	if err != nil {
		return nil, fmt.Errorf("new config center client error: %s", err)
	}
	data, ok := client.Value2(_agentConfigKey)
	if !ok {
		return nil, fmt.Errorf("load config center error key %s not found", _agentConfigKey)
	}
	cfg := new(AgentConfig)
	_, err = toml.Decode(data, cfg)
	if err != nil {
		return nil, fmt.Errorf("could not decode config file %s, error: %s", _agentConfigKey, err)
	}
	go func() {
		for range client.Event() {
			// ignore config change event
		}
	}()
	return cfg, nil
}
