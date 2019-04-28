package conf

import (
	"errors"
	"flag"
	"go-common/library/net/rpc/warden"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	//Conf .
	Conf   = &Config{}
	client *conf.Client
)

//Config .
type Config struct {
	Env string
	// base
	// host
	Host *Host
	// channal len
	ChanSize int64
	// log
	Xlog *log.Config
	// http
	BM *bm.ServerConfig
	// Auth
	Auth *permit.Config
	// tracer
	Tracer *trace.Config
	// tick load pgc
	Tick time.Duration
	// db
	DB *DB
	// db
	ORMArchive *orm.Config
	// databus
	VideoupPub  *databus.Config
	UpCreditPub *databus.Config
	// redis
	Redis *Redis
	// hbase
	HBase *hbaseConf
	// http client test
	HTTPClient HTTPClient
	// rpc
	AccountRPC    *warden.ClientConfig
	UpsRPC        *warden.ClientConfig
	TagDisRPC     *rpc.ClientConfig
	Ecode         *ecode.Config
	ManagerReport *databus.Config
}

type hbaseConf struct {
	hbase.Config
	ReadTimeout   time.Duration
	ReadsTimeout  time.Duration
	WriteTimeout  time.Duration
	WritesTimeout time.Duration
}

//Host .
type Host struct {
	API       string
	MngSearch string
	Manager   string
	Data      string
	Account   string
	Task      string
	Archive   string
}

//DB .
type DB struct {
	Archive     *sql.Config
	ArchiveRead *sql.Config
	Manager     *sql.Config
	Oversea     *orm.Config
	Creative    *sql.Config
}

//Redis .
type Redis struct {
	Track *struct {
		*redis.Config
		Expire time.Duration
	}
	Secondary *struct {
		*redis.Config
		Expire time.Duration
	}
}

// HTTPClient test
type HTTPClient struct {
	Read   *bm.ClientConfig
	Write  *bm.ClientConfig
	Search *bm.ClientConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

//Init .
func Init() (err error) {
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
