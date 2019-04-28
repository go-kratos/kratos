package conf

import (
	"errors"
	"flag"

	"go-common/app/admin/main/block/model"
	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	Log           *log.Config
	Tracer        *trace.Config
	Auth          *permit.Config
	Ecode         *ecode.Config
	BM            *bm.ServerConfig
	RPCClients    *RPC
	Memcache      *memcache.Config
	MySQL         *sql.Config
	HTTPClient    *bm.ClientConfig
	Perms         *Perms
	AccountNotify *databus.Config
	Property      *Property
	// manager log config
	ManagerLog *databus.Config
}

// Property .
type Property struct {
	BlackHouseURL string
	MSGURL        string
	TelURL        string
	MailURL       string
	MSG           *MSG
}

// MSG .
type MSG struct {
	BlackHouseLimit   model.MSG
	BlackHouseForever model.MSG
	SysLimit          model.MSG
	SysForever        model.MSG
	BlockRemove       model.MSG
}

// RPC .
type RPC struct {
	Account *rpc.ClientConfig
	Figure  *rpc.ClientConfig
	Spy     *rpc.ClientConfig
}

// Perms .
type Perms struct {
	Perm map[string]string
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
