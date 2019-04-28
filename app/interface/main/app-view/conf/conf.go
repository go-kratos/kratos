package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf is
var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

//Config struct
type Config struct {
	// Env
	Env       string
	AutoLimit int64
	DMRegion  []int16
	// interface XLog
	XLog    *log.Config
	Bnj2019 *Bnj2019
	// infoc
	InfocCoin   *infoc.Config
	InfocView   *infoc.Config
	InfocRelate *infoc.Config
	UseractPub  *databus.Config
	DislikePub  *databus.Config
	// tick time
	Tick xtime.Duration
	// vip tick
	VipTick xtime.Duration
	// tracer
	Tracer *trace.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// httpWrite
	HTTPWrite *bm.ClientConfig
	// httpbangumi
	HTTPBangumi *bm.ClientConfig
	// httpaudio
	HTTPAudio *bm.ClientConfig
	// http
	BM *HTTPServers
	// httpAd
	HTTPAD *bm.ClientConfig
	// httpGame
	HTTPGame *bm.ClientConfig
	// HTTPAsync
	HTTPAsync *bm.ClientConfig
	// HTTPGameAsync
	HTTPGameAsync *bm.ClientConfig
	// httpClinet
	HTTPSearch *bm.ClientConfig
	// host
	Host *Host
	// rpc client
	AccountRPC  *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	TagRPC      *rpc.ClientConfig
	AssistRPC   *rpc.ClientConfig
	HisRPC      *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	RelationRPC *rpc.ClientConfig
	FavoriteRPC *rpc.ClientConfig
	CoinRPC     *rpc.ClientConfig
	DMRPC       *rpc.ClientConfig
	ActivityRPC *rpc.ClientConfig
	LocationRPC *rpc.ClientConfig
	// db
	MySQL *MySQL
	// ecode
	Ecode *ecode.Config
	// mc
	Memcache *Memcache
	// PlayURL
	PlayURL *PlayURL
	// 相关推荐秒开个数
	RelateCnt  int
	RelateGray int64
	// buildLimit
	BuildLimit *BuildLimit
	// Warden Client
	ThumbupClient *warden.ClientConfig
}

// Bnj2019 is
type Bnj2019 struct {
	Tick          xtime.Duration
	MainAid       int64
	AdAv          int64
	AidList       []int64
	ElecBigText   string
	ElecSmallText string
	WhiteMids     []int64
	FakeElec      int
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
}

// Host struct
type Host struct {
	Bangumi   string
	APICo     string
	Activity  string
	Elec      string
	AD        string
	Data      string
	Archive   string
	APILiveCo string
	Game      string
	VIP       string
	AI        string
	Search    string
	Bvcvod    string
}

// MySQL struct
type MySQL struct {
	Show    *sql.Config
	Manager *sql.Config
}

// Memcache struct
type Memcache struct {
	Archive *struct {
		*memcache.Config
		ArchiveExpire  xtime.Duration
		RelateExpire   xtime.Duration
		AddonExpire    xtime.Duration
		RecommedExpire xtime.Duration
	}
}

// PlayURL playurl token's secret.
type PlayURL struct {
	Secret string
}

// BuildLimt for build limit
type BuildLimit struct {
	CooperationIOS     int
	CooperationAndroid int
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
	client.Watch("app-view.toml")
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
