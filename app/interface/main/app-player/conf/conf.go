package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf all conf
	Conf = &Config{}
)

// Config is
type Config struct {
	HTTPClient  *bm.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	Ecode       *ecode.Config
	Log         *log.Config
	Host        *Host
	AidGray     int64
	PadAid      int64
	PadCid      int64
	PhoneAid    int64
	PhoneCid    int64
	PadHDAid    int64
	PadHDCid    int64
	Bnj         *Bnj
	// mc
	Memcache *memcache.Config
	// Warden Client
	ArchiveClient *warden.ClientConfig
	UGCpayClient  *warden.ClientConfig
	AccountClient *warden.ClientConfig
}

// Bnj is
type Bnj struct {
	Tick xtime.Duration
	Aids []int64
}

// Host struct
type Host struct {
	Playurl   string
	PlayurlBk string
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
	client.Watch("app-player.toml")
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
