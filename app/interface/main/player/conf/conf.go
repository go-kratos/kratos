package conf

import (
	"errors"
	"flag"
	xtime "time"

	"go-common/app/interface/main/player/model"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	Conf     = &Config{}
)

// Config is service conf.
type Config struct {
	// 广播
	Broadcast Broadcast
	// policy
	Policy *model.Policy
	// policy items
	Pitem []*model.Pitem
	// 拜年祭
	Matsuri Matsuri
	// XLog
	XLog *log.Config
	// ecode
	Ecode *ecode.Config
	// host
	Host *Host
	// tracer
	Tracer *trace.Config
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// bm
	BM *HTTPServers
	// mysql
	MySQL *MySQL
	// rpc
	ArchiveRPC  *rpc.ClientConfig
	AccountRPC  *rpc.ClientConfig
	HistoryRPC  *rpc.ClientConfig
	AssistRPC   *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	Dm2RPC      *rpc.ClientConfig
	LocRPC      *rpc.ClientConfig
	TagRPC      *rpc.ClientConfig
	// HTTPClient
	HTTPClient *bm.ClientConfig
	// Rule
	Rule *Rule
	Tick *Tick
	// Infoc2
	Infoc2 *infoc.Config
	// PlayURLToken
	PlayURLToken *PlayURLToken
	// grpc client
	AccClient    *warden.ClientConfig
	ArcClient    *warden.ClientConfig
	UGCPayClient *warden.ClientConfig
	// icon
	Icon *icon
	// bnj2019
	Bnj2019 *bnj2019
}

// Tick tick time.
type Tick struct {
	// tick time
	CarouselTick time.Duration
	ParamTick    time.Duration
	IconTick     time.Duration
}

// Rule rules
type Rule struct {
	// timeout
	VsTimeout   time.Duration
	NoAssistMid int64
	VipQn       []int
	LoginQn     int
	MaxFreeQn   int
	AutoQn      int
	PlayurlGray int64
}

// Host is host info
type Host struct {
	APICo       string
	AccCo       string
	PlayurlCo   string
	H5Playurl   string
	HighPlayurl string
}

// MySQL mysql.
type MySQL struct {
	Show *sql.Config
}

// Broadcast breadcast.
type Broadcast struct {
	TCPAddr string
	WsAddr  string
	WssAddr string
	Begin   string
	End     string
}

// Matsuri matsuri.
type Matsuri struct {
	PastID  int64
	MatID   int64
	MatTime string
	Tick    time.Duration
}

// PlayURLToken playurl auth token.
type PlayURLToken struct {
	Secret      string
	PlayerToken string
}

// HTTPServers bm servers config.
type HTTPServers struct {
	Outer *bm.ServerConfig
}

type bnj2019 struct {
	BnjMainAid  int64
	BnjListAids []int64
	BnjTick     time.Duration
}

type icon struct {
	Start xtime.Time
	End   xtime.Time
	URL1  string
	Hash1 string
	URL2  string
	Hash2 string
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
