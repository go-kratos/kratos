package conf

import (
	"errors"
	"flag"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc/warden"
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
	Tracer     *trace.Config
	Memcache   *memcache.Config
	MySQL      *sql.Config
	Ecode      *ecode.Config
	Auth       *permit.Config
	BFS        *BFS
	Host       *Host
	GRPCClient *GRPCClient
	// http client
	HTTPClient *bm.ClientConfig
	// manager log config
	ManagerLog *databus.Config
	Property   *Property
}

// Property .
type Property struct {
	MSG []*model.MSG
}

// GRPCClient .
type GRPCClient struct {
	Account *warden.ClientConfig
	Member  *warden.ClientConfig
	Archive *warden.ClientConfig
}

// Host host config .
type Host struct {
	Bfs     string
	Msg     string
	Videoup string
	API     string
}

// BFS bfs config
type BFS struct {
	Bucket string
	Key    string
	Secret string
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
