package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

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
	Log         *log.Config
	Tracer      *trace.Config
	BM          *bm.ServerConfig
	Verify      *verify.Config
	Memcache    *memcache.Config
	MySQL       *sql.Config
	Ecode       *ecode.Config
	LocalCache  *LocalCache
	ArchiveGRPC *warden.ClientConfig
	AccountGRPC *warden.ClientConfig
	CacheTTL    *CacheTTL
	Biz         *Biz
}

// Biz .
type Biz struct {
	ElecAVRankSize int
	ElecUPRankSize int
	RAMAVIDs       []int64
	RAMUPIDs       []int64
	ReloadDuration xtime.Duration
}

// LocalCache prop.
type LocalCache struct {
	ElecAVRankSize int
	ElecAVRankTTL  xtime.Duration
	ElecUPRankSize int
	ElecUPRankTTL  xtime.Duration
}

// CacheTTL .
type CacheTTL struct {
	ElecUPRankTTL      int32
	ElecAVRankTTL      int32
	ElecPrepUPRankTTL  int32
	ElecPrepAVRankTTL  int32
	ElecUserSettingTTL int32
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
