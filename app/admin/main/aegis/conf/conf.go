package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log      *log.Config
	BM       *bm.ServerConfig
	Verify   *verify.Config
	Auth     *permit.Config
	Tracer   *trace.Config
	Redis    *Redis
	Memcache *memcache.Config

	Ecode *ecode.Config
	RPC   *RPC
	Host  *Host
	// http client test
	HTTPClient HTTPClient
	// ORM
	ORM      *orm.Config
	MySQL    *sql.Config
	AegisPub *databus.Config
	Bfs      *Bfs

	Bucket   string
	Debug    string
	Admin    string // 所有业务的管理员
	Consumer *Consumer
	Gray     *Gray

	Auditstate map[string]string

	GRPC *GRPC
}

//GRPC .
type GRPC struct {
	AccRPC *warden.ClientConfig
	UpRPC  *warden.ClientConfig
}

// Gray .
type Gray struct {
	Biz []graybiz
}

type graybiz struct {
	BusinessID int64
	Options    []grayoption
}

type grayoption struct {
	Fields []struct {
		Name  string
		Value string
	}
}

// Consumer 在线过期时间，角色过期时间
type Consumer struct {
	OnExp   int32
	RoleExp int32
}

// Bfs reprensents the bfs config
type Bfs struct {
	Key         string
	Secret      string
	Host        string
	Timeout     int
	MaxFileSize int
}

// Host host config .
type Host struct {
	API        string
	Manager    string
	MainSearch string
}

// Redis .
type Redis struct {
	NetExpire xtime.Duration
	Cluster   *redis.Config
}

// HTTPClient str
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
	Es    *bm.ClientConfig
}

//DB .
type DB struct {
	Aegis *sql.Config
	MySQL *sql.Config
}

// RPC .
type RPC struct {
	Acc *rpc.ClientConfig
	Rel *rpc.ClientConfig
	Up  *rpc.ClientConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	if err = load(); err != nil {
		return
	}
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
