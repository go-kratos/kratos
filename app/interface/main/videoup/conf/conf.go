package conf

import (
	"errors"
	"flag"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	xlog "go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath str
	ConfPath string
	// Conf cfg
	Conf   = &Config{}
	client *conf.Client
)

// Config struct
type Config struct {
	AppEditorInfoc *infoc.Config
	// base
	Bfs   *Bfs
	Limit *Limit
	// db
	DB *DB
	// tick
	Tick time.Duration
	// app
	App *bm.App
	// host
	Host *Host
	// http
	HTTPClient *HTTPClient
	// xlog
	Xlog *xlog.Config
	// tracer
	Tracer *trace.Config
	// ecode
	Ecode *ecode.Config
	// rpc client
	RelationRPC *rpc.ClientConfig
	SubRPC      *rpc.ClientConfig
	// mc
	Memcache *Memcache
	// redis
	Redis       *Redis
	UpCoverAnti *antispam.Config
	Game        *Game
	BM          *bm.ServerConfig
	// geetest
	Geetest             *Geetest
	MaxAllVsCnt         int
	MaxAddVsCnt         int
	UgcPayAllowEditDays int
	AccClient           *warden.ClientConfig
	UpClient            *warden.ClientConfig
}

// Geetest geetest id & key
type Geetest struct {
	CaptchaID   string
	MCaptchaID  string
	PrivateKEY  string
	MPrivateKEY string
}

// Game str Conf
type Game struct {
	OpenHost string
	App      *bm.App
}

// DB conf.
type DB struct {
	// archive db
	Manager *sql.Config
}

// Limit config
type Limit struct {
	AddBasicExp time.Duration
}

// Bfs bfs config
type Bfs struct {
	Timeout     time.Duration
	MaxFileSize int
}

// Host conf
type Host struct {
	Account    string
	Archive    string
	APICo      string
	WWW        string
	Member     string
	UpMng      string
	Tag        string
	Elec       string
	Geetest    string
	Dynamic    string
	MainSearch string
	Chaodian   string //超电
}

// HTTPClient str
type HTTPClient struct {
	Read     *bm.ClientConfig
	Write    *bm.ClientConfig
	UpMng    *bm.ClientConfig
	FastRead *bm.ClientConfig
	Chaodian *bm.ClientConfig
}

// Memcache str
type Memcache struct {
	Account *struct {
		*memcache.Config
		SubmitExpire time.Duration
	}
}

// Redis str
type Redis struct {
	Videoup *struct {
		*redis.Config
		Expire time.Duration
	}
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init
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
			xlog.Info("config reload")
			if load() != nil {
				xlog.Error("config reload error (%v)", err)
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
