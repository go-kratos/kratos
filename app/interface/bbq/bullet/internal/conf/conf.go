package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	antispam "go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"

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
	Log          *log.Config
	BM           *bm.ServerConfig
	Verify       *verify.Config
	Auth         *auth.Config
	Tracer       *trace.Config
	Redis        *redis.Config
	MySQL        *sql.Config
	OnlineMySQL  *sql.Config
	Ecode        *ecode.Config
	GRPCClient   map[string]*GRPCConf
	AntiSpam     map[string]*antispam.Config
	BulletConfig BulletConfig
	Infoc        *infoc.Config
}

// BulletConfig 弹幕的一些配置项
type BulletConfig struct {
	CloseWrite bool
	CloseRead  bool
}

// GRPCConf .
type GRPCConf struct {
	WardenConf *warden.ClientConfig
	Addr       string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() (err error) {
	if confPath != "" {
		err = local()
	} else {
		err = remote()
	}

	if Conf.Redis != nil {
		for _, anti := range Conf.AntiSpam {
			anti.Redis = Conf.Redis
		}
	}

	return
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
