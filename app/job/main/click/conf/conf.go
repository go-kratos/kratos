package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

type Config struct {
	// Env
	Env            string
	LastChangeTime int64
	ReleaseTime    int64
	BnjMainAid     int64
	BnjListAids    []int64
	// interface XLog
	XLog *log.Config
	// tracer
	Tracer *trace.Config
	// stat databus pub
	StatPub            *databus.Config
	StatViewPub        *databus.Config
	ReportMergeDatabus *databus.Config
	// click databus pub
	ClickPub   *databus.Config
	ArchiveRPC *rpc.ClientConfig
	// http
	BM *bm.ServerConfig
	// redis
	Redis      *redis.Config
	HTTPClient *bm.ClientConfig
	// cache time conf
	CacheConf struct {
		PGCReplayTime           int64
		ArcUpCacheTime          int64
		NewAnonymousCacheTime   int64
		NewAnonymousBvCacheTime int64
	}
	// db
	DB *sql.Config
	// hash number
	HashNum int64
	// chan number
	ChanNum int64
	// consumer num
	ConsumeNum int
	// need Init
	NeedInit bool
	// infoc
	Infoc2 *infoc.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
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
	client.Watch("click-job.toml")
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
