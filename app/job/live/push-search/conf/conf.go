package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
	"go-common/library/queue/databus"
	"go-common/library/database/hbase.v2"
	"go-common/library/net/rpc/liverpc"
	xtime "go-common/library/time"
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
	Tracer   *trace.Config
	Redis    *redis.Config
	Memcache *memcache.Config
	MySQL    *sql.Config
	Ecode    *ecode.Config
	LiveRpc  map[string]*liverpc.ClientConfig
	//DataBus
	DataBus     *DataBus
	Group       *Group
	SearchHBase *hbaseConf
	MigrateNum	int
}

type DataBus struct {
	RoomInfo *databus.Config
	Attention *databus.Config
	UserName *databus.Config
	PushSearch *databus.Config
}

// Group group.
type Group struct {
	RoomInfo *GroupConf
	Attention *GroupConf
	UserInfo *GroupConf
}

// GroupConf group conf.
type GroupConf struct {
	Num    int
	Chan   int
}

type hbaseConf struct {
	hbase.Config
	ReadTimeout   xtime.Duration
	ReadsTimeout  xtime.Duration
	WriteTimeout  xtime.Duration
	WritesTimeout xtime.Duration
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
