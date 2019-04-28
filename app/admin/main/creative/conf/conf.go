package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// ConfPath str
var (
	ConfPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config str
type Config struct {
	// base
	// channal len
	ChanSize int64
	// log
	Log *log.Config
	// identify
	App *bm.App
	// auth
	Auth *permit.Config
	// tracer
	Tracer *trace.Config
	// tick load pgc
	Tick time.Duration
	// orm
	ORM *orm.Config
	//orm
	ORMArchive *orm.Config
	// redis
	Redis *Redis
	// HTTPClient client
	HTTPClient *bm.ClientConfig
	// rpc
	ArchiveRPC *rpc.ClientConfig
	ArticleRPC *rpc.ClientConfig
	// http
	BM *bm.ServerConfig
	// bfs
	Bfs *Bfs
	//ManagerReport
	ManagerReport *databus.Config
	//Ecode
	Ecode *ecode.Config
	// host
	Host *Host
	// grpc
	AccClient *warden.ClientConfig
}

// Host host config .
type Host struct {
	Msg string
}

// Bfs reprensents the bfs config
type Bfs struct {
	Key         string
	Secret      string
	Host        string
	Timeout     int
	MaxFileSize int
}

// Redis str
type Redis struct {
	Track *struct {
		*redis.Config
		Expire time.Duration
	}
}

// HTTPClient str
type HTTPClient struct {
	Read  *bm.ClientConfig
	Write *bm.ClientConfig
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init fn
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
