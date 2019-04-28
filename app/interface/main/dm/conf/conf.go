package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf export config variable
	Conf = &Config{}
)

// Config dm config file
type Config struct {
	// base
	// ecode
	Ecode *ecode.Config
	// xlog
	Xlog   *log.Config
	Auth   *auth.Config
	Verify *verify.Config
	// http client
	HTTPClient *bm.ClientConfig
	ArchiveRPC *rpc.ClientConfig
	AccountRPC *warden.ClientConfig
	AssistRPC  *rpc.ClientConfig
	DMRPC      *rpc.ClientConfig
	// http server
	HTTPServer *bm.ServerConfig
	// db
	DB *DB
	// redis
	Redis *Redis
	// tracer
	Tracer *trace.Config
	// databus
	Databus *databus.Config
	// Antispam
	Antispam *antispam.Config
	Host     Host
	ES       *elastic.Config
}

// Host hosts used in dm admin
type Host struct {
	API     string
	Archive string
	Message string
}

// DB mysql database instance
type DB struct {
	DM           *sql.Config
	DMMetaReader *sql.Config
	DMWriter     *sql.Config
}

// Redis redis instance
type Redis struct {
	DM *DMRedis
}

// DMRedis redis instance of dm
type DMRedis struct {
	*redis.Config
	DMIDExpire  time.Duration
	LockExpire  time.Duration
	IndexExpire time.Duration
	VideoExpire time.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

//Init int config
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
