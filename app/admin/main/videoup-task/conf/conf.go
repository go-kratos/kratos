package conf

import (
	"flag"
	"fmt"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

//Config config.
type Config struct {
	Xlog          *log.Config
	ManagerReport *databus.Config
	HTTPClient    *bm.ClientConfig
	BM            *bm.ServerConfig
	DB            *db
	Auth          *permit.Config
	Host          host
	HBase         *HBaseConfig
	// redis
	Redis *Redis

	GRPC *GRPC
}

// HBaseConfig extra hbase config
type HBaseConfig struct {
	*hbase.Config
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

//GRPC .
type GRPC struct {
	AccRPC *warden.ClientConfig
	UpsRPC *warden.ClientConfig
}

// Redis .
type Redis struct {
	Weight *struct {
		*redis.Config
		Expire time.Duration
	}
}

type host struct {
	API     string
	Manager string
	Search  string
}

type db struct {
	Archive     *sql.Config
	ArchiveRead *sql.Config
	Manager     *sql.Config
}

//common + xlog(agent) + trace(better) + http + perf(web代码性能分析) + os.signal监听（stop/reload服务）
var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

//Init config init.
func Init() error {
	if confPath != "" {
		return local()
	}

	return remote()
}

func local() (err error) {
	tmpConf := &Config{}
	if _, err = toml.DecodeFile(confPath, tmpConf); err == nil {
		Conf = tmpConf
	}

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
			if err = load(); err != nil {
				log.Error("config reload error (%v)", err)
			}
		}
	}()
	return
}

func load() (err error) {
	var (
		ok      bool
		tmpConf = &Config{}
	)

	if confPath, ok = client.Toml2(); !ok {
		err = fmt.Errorf("config load error")
		return
	}

	if _, err = toml.Decode(confPath, tmpConf); err != nil {
		err = fmt.Errorf("couldn't decode config, error (%v)", err)
		return
	}

	Conf = tmpConf
	return
}
