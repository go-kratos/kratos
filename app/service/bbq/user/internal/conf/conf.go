package conf

import (
	"context"
	"encoding/json"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	// Conf config
	Conf = &Config{}
	// UnameConf .
	UnameConf = &UnameConfig{}
)

// UnameConfig .
type UnameConfig struct {
	ForbiddenUname []string `json:"forbidden_uname"`
	unameSet       map[string]bool
}

// UnameForbidden 判断uname是否被禁用
func (c *UnameConfig) UnameForbidden(uname string) bool {
	_, exists := c.unameSet[uname]
	return exists
}

// Set .
func (c *UnameConfig) Set(text string) error {
	log.Infow(context.Background(), "log", "reload uname config")
	if err := json.Unmarshal([]byte(text), &UnameConf); err != nil {
		panic(err)
	}
	c.unameSet = make(map[string]bool)
	for _, uname := range c.ForbiddenUname {
		c.unameSet[uname] = true
	}
	log.Infow(context.Background(), "log", "reload uname config succ", "uname_size", len(c.unameSet))
	return nil
}

// Config .
type Config struct {
	Log        *log.Config
	BM         *bm.ServerConfig
	Verify     *verify.Config
	Tracer     *trace.Config
	Redis      *redis.Config
	MySQL      *sql.Config
	Ecode      *ecode.Config
	GRPC       *warden.ServerConfig
	GRPCClient map[string]*GRPCConf
	Databus    map[string]*databus.Config
}

//GRPCConf .
type GRPCConf struct {
	WardenConf *warden.ClientConfig
	Addr       string
}

// Set .
func (c *Config) Set(text string) error {
	log.Infow(context.Background(), "log", "reload config")
	if _, err := toml.Decode(text, c); err != nil {
		panic(err)
	}
	return nil
}
