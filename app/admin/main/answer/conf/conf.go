package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	confPath string
	client   *conf.Client
)

// Config struct of conf.
type Config struct {
	// log
	Log *log.Config
	// db
	Mysql *orm.Config
	// redis
	Redis *Redis
	// tracer
	Tracer *trace.Config
	// Answer
	Answer *Answer
	// bm
	BM   *bm.ServerConfig
	Auth *permit.Config
}

// Redis .
type Redis struct {
	*redis.Config
	Expire                time.Duration
	AnsCountExpire        time.Duration
	AnsAddFlagCountExpire time.Duration
}

// Answer conf.
type Answer struct {
	Debug        bool
	FontFilePath string
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
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

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init create config instance.
func Init() (err error) {
	if confPath == "" {
		err = configCenter()
	} else {
		_, err = toml.DecodeFile(confPath, &Conf)
	}
	return
}
