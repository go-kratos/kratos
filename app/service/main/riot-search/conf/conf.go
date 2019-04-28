package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

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
	Log        *log.Config
	BM         *bm.ServerConfig
	Verify     *verify.Config
	Tracer     *trace.Config
	Ecode      *ecode.Config
	Riot       *RiotConfig
	HTTPClient *bm.ClientConfig
	Databus    *databus.Config
	// db
	Mysql *sql.Config
	UT    bool
}

// RiotConfig configs
type RiotConfig struct {
	Dict      string
	StopToken string
	LoadPath  string
	// timeout ms
	Timeout int
	// interval(second) to flush index
	FlushTime int64
	// 索引器shards数目
	NumShards int
	// 持久化文件数目，尽量保证每个文件小于100M
	StoreShards int
	StoreFolder string
	// 持久化存储引擎，可选项bg(badger),leveldb,bolt
	StoreEngine string
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
