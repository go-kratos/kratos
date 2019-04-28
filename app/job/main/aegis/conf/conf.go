package conf

import (
	"errors"
	"flag"

	"go-common/app/job/main/aegis/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"

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
	Debug    bool
	Log      *log.Config
	BM       *bm.ServerConfig
	Verify   *verify.Config
	Tracer   *trace.Config
	Redis    *redis.Config
	Memcache *memcache.Config
	MySQL    *MySQL
	Ecode    *ecode.Config
	// ORM
	ORM *orm.Config
	// DataBus databus
	DataBus *DataBus
	// Databusutil
	Databusutil *Databusutil
	// RPC
	RPC *RPC
	//GRPC
	GRPC *GRPC
	// BizConfiger
	BizCfg BizConfiger
	HTTP   *HTTP

	Host *Host

	// mail
	Mail *Mail
}

//MySQL .
type MySQL struct {
	Slow *sql.Config
	Fast *sql.Config
}

//Host .
type Host struct {
	API     string
	Videoup string
}

//HTTP .
type HTTP struct {
	Fast *bm.ClientConfig
	Slow *bm.ClientConfig
}

//BizConfiger .
type BizConfiger struct {
	WeightOpt []*model.WeightOPT
}

//RPC .
type RPC struct {
	Rel *rpc.ClientConfig
	Up  *rpc.ClientConfig
}

//GRPC .
type GRPC struct {
	Up  *warden.ClientConfig
	Acc *warden.ClientConfig
}

// DataBus databus infomation
type DataBus struct {
	BinLogSub   *databus.Config
	ResourceSub *databus.Config
	TaskSub     *databus.Config
	ArchiveSub  *databus.Config
}

//Mail 邮件配置
type Mail struct {
	Host               string
	Port               int
	Username, Password string
}

// Databusutil databus group
type Databusutil struct {
	Task     *databusutil.Config
	Resource *databusutil.Config
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
