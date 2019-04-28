package conf

import (
	"go-common/app/common/openplatform/encoding"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// Conf common conf
	Conf = &Config{}
)

//Config config struct
type Config struct {
	//数据库配置
	DB *DB
	// redis
	Redis *Redis
	// http client
	HTTPClient HTTPClient
	// http
	BM *blademaster.ServerConfig
	// tracer
	Tracer *trace.Config
	// log
	Log *log.Config
	// UT
	UT         *UT
	GRPCClient map[string]*warden.ClientConfig
	Encrypt    *encoding.EncryptConfig
	URLs       map[string]string
	//basecenter配置
	BaseCenter *BaseCenter
	Databus    map[string]*databus.Config

	TestProject *TestProject
}

// HTTPClient config
type HTTPClient struct {
	Read  *blademaster.ClientConfig
	Write *blademaster.ClientConfig
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *blademaster.ServerConfig
	Local *blademaster.ServerConfig
}

// Redis config
type Redis struct {
	Master *redis.Config
	Expire time.Duration
}

// DB config
type DB struct {
	Master *sql.Config
}

// UT config
type UT struct {
	DistPrefix string
}

//BaseCenter 的配置
type BaseCenter struct {
	AppID string
	Token string
}

// TestProject 测试项目配置
type TestProject struct {
	IDs        []int64
	CheckQuery string
}

// Set set config and decode.
func (c *Config) Set(text string) error {
	var tmp Config
	if _, err := toml.Decode(text, &tmp); err != nil {
		return err
	}
	*c = tmp
	return nil
}
