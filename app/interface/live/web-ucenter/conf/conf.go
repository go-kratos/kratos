package conf

import (
	"errors"
	"flag"

	"go-common/library/net/http/blademaster/middleware/auth"

	"go-common/library/net/rpc"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"

	"go-common/library/log/infoc"
	"go-common/library/net/rpc/warden"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

const (
	// APPKey to call main site api
	APPKey = "fb06a25c6338edbc"
)
const (
	// MainInnerHostHTTP api host
	MainInnerHostHTTP = "http://api.bilibili.co"
)

// Config .
type Config struct {
	Log        *log.Config
	BM         *bm.ServerConfig
	Verify     *verify.Config
	Tracer     *trace.Config
	Redis      *redis.Config
	Memcache   *memcache.Config
	MySQL      *sql.Config
	Ecode      *ecode.Config
	LiveRpc    map[string]*liverpc.ClientConfig
	HTTPClient *bm.ClientConfig
	HistoryRPC *rpc.ClientConfig
	Auth       *auth.Config
	Warden     *warden.ClientConfig
	Infoc      *Infoc
	Host       *Host
	AccountRPC *rpc.ClientConfig
}

// Infoc .
type Infoc struct {
	CapsuleInfoc *infoc.Config
}

// Host prc host
type Host struct {
	LiveRpc string
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
