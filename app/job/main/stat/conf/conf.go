package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/queue/databus"
	"go-common/library/time"

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
	// Env
	Env string
	// monitor
	MonitorIdle time.Duration
	// interface XLog
	XLog *log.Config
	// BM
	BM *bm.ServerConfig
	// http client
	HTTPClient *bm.ClientConfig
	// databus
	StatSub   *databus.Config
	ViewSub   *databus.Config
	DmSub     *databus.Config
	ReplySub  *databus.Config
	FavSub    *databus.Config
	CoinSub   *databus.Config
	ShareSub  *databus.Config
	RankSub   *databus.Config
	LikeSub   *databus.Config
	Memcaches []*memcache.Config
	// DB
	DB      *sql.Config
	ClickDB *sql.Config
	// rpc
	ArchiveRPC  *rpc.ClientConfig
	ArchiveRPC2 *rpc.ClientConfig
	// Monitor
	Monitor *Monitor
}

// Monitor is
type Monitor struct {
	Users  string
	Token  string
	Secret string
	URL    string
}

// SMS config
type SMS struct {
	Phone string
	Token string
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
