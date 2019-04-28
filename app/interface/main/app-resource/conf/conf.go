package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	Conf     = &Config{}
	client   *conf.Client
)

type Config struct {
	// Env
	Env string
	// show  XLog
	Log *log.Config
	// tick time
	Tick xtime.Duration
	// tracer
	Tracer *trace.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// bm http
	BM *HTTPServers
	// db
	Ecode *ecode.Config
	MySQL *MySQL
	// duration
	Duration *Duration
	// Splash
	Splash *Splash
	// interestJSONFile
	InterestJSONFile string
	// StaticJsonFile
	StaticJSONFile string
	// guide rand
	GuideRandom *GuideRandom
	// domain
	Domain *Domain
	ABTest *ABTest
	// host
	Host *Host
	// sideBar limit id
	SideBarLimit []int64
	// resource
	ResourceRPC *rpc.ClientConfig
	// infoc2
	InterestInfoc *infoc.Config
	// BroadcastRPC grpc
	BroadcastRPC *warden.ClientConfig
	// White
	White *White
	// 垃圾白名单
	ShowTabMids []int64
	// location rpc
	LocationRPC *rpc.ClientConfig
	// show hot all
	ShowHotAll bool
	// rpc server2
	RPCServer *rpc.ServerConfig
}

type HTTPServers struct {
	Outer *bm.ServerConfig
}

type Host struct {
	Ad   string
	Data string
	VC   string
}

type White struct {
	List map[string][]string
}

type ABTest struct {
	Range int
}

type GuideRandom struct {
	Random map[string]int
	Buvid  map[string]int
	Feed   uint32
}

type Duration struct {
	// splash
	Splash string
}

type Splash struct {
	Random map[string][]string
}

type MySQL struct {
	Show     *sql.Config
	Resource *sql.Config
}

type Domain struct {
	Addr      []string
	ImageAddr []string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

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
	client.Watch("app-resource.toml")
	go func() {
		for range client.Event() {
			log.Info("config reload")
			if load() != nil {
				log.Error("config reload error(%v)", err)
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
