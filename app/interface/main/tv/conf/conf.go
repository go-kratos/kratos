package conf

import (
	"errors"
	"flag"

	"go-common/app/interface/main/tv/model"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	confPath string
	client   *conf.Client
)

// Config struct of conf.
type Config struct {
	// zone configure
	Newzone map[string]*PageCfg
	// log
	Log *log.Config
	// tracer
	Tracer *trace.Config
	// http server config
	HTTPServer *bm.ServerConfig
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// mysql
	Mysql *sql.Config
	// memcache
	Memcache *Memcache
	// app
	TVApp *TVApp
	// homepage settings
	Homepage *PageConf
	// HTTPClient .
	HTTPClient    *bm.ClientConfig
	PlayurlClient *bm.ClientConfig
	SearchClient  *bm.ClientConfig
	// Redis
	Redis *Redis
	// Cfg common configuration
	Cfg *Cfg
	// Search Config
	Search *Search
	// RPC config
	ArcClient   *warden.ClientConfig
	AccClient   *warden.ClientConfig
	HisRPC      *rpc.ClientConfig
	FavoriteRPC *rpc.ClientConfig
	TvVipClient *warden.ClientConfig
	// Ip Whitelist
	IP *IP
	// ecode
	Ecode *ecode.Config
	// api url
	Host   *Host
	Region *Region
	Style  *Style
	Wild   *Wild
}

// IPWhite .
type IPWhite struct {
	TvVip []string
}

// IP .
type IP struct {
	White *IPWhite
}

// Style label .
type Style struct {
	LabelSpan xtime.Duration
}

// Region .
type Region struct {
	StopSpan xtime.Duration // get region time span
}

// IndexLabel def.
type IndexLabel struct {
	Fre       xtime.Duration
	PGCOrder  []string             // pgc order
	UGCOrder  []string             // ugc order
	YearV     map[string]*YearVDur // year value pair
	YearParam []string             // year params = pub_date, year
}

// YearVDur def
type YearVDur struct {
	Dur string `json:"dur"`
}

// IsYear distinguishes whether the param is year type param
func (u *IndexLabel) IsYear(param string) bool {
	for _, v := range u.YearParam {
		if v == param {
			return true
		}
	}
	return false
}

// Host api urls
type Host struct {
	Data        string // data.bilibili.co
	APIIndex    string // homepage pgc data source
	APIZone     string // zonepage pgc data source
	APIFollow   string // pgc follow
	APIMedia    string // pgc media detail
	APIMediaV2  string // pgc media detail v2
	APIRecom    string // pgc recom
	APINewindex string // pgc index_show
	UgcPlayURL  string // ugc play url
	AIUgcType   string // ai ugc type data
	APICo       string
	FavAdd      string // favorite add url
	FavDel      string // favorite del url
	ReqURL      string // version update request url
	ESHost      string // manager url
}

// Wild .
type Wild struct {
	WildSearch *WildSearch
}

// WildSearch wild search .
type WildSearch struct {
	UserNum        int
	UserVideoLimit int
	BiliUserNum    int
	BiliUserVl     int
	SeasonNum      int
	MovieNum       int
	SeasonMore     int
	MovieMore      int
}

// Cfg def.
type Cfg struct {
	ZonePs          int                // Zone index page size
	AuthMsg         *AuthMsg           // auth error message config
	ZonesInfo       *ZonesInfo         // all the zones info
	Dangbei         *Dangbei           // dangbei configuration
	PageReload      xtime.Duration     // all page reload duration
	IndexShowReload xtime.Duration     // index show reload duration
	EsIntervReload  xtime.Duration     // es intervention reload duration
	DefaultSplash   string             // default splash url
	FavPs           int                // favorite cfg
	PGCFilterBuild  int                // the build number, under which we export only pgc modules and data
	VipQns          []string           // the qualities dedicated for vips
	HisCfg          *HisCfg            // history related cfg
	EsIdx           *EsIdx             // elastic search index page cfg
	IndexLabel      *IndexLabel        // index label cfg
	EmptyArc        *EmptyArc          // chan size
	VipMark         *VipMark           // vip mark
	SnVipCorner     *model.SnVipCorner // season vip corner mark cfg
	AuditSign       *AuditSign
}

// AuditSign cfg is used to check license owner requests
type AuditSign struct {
	Key    string
	Secret string
}

// TvVip def.
type TvVip struct {
	Build int64
	Msg   string
}

// VipMark def.
type VipMark struct {
	V1HideChargeable bool // whether we hide chargeable episode in pgc view V1
	EpFree           int  // ep's pay status which means free
	EP               *model.CornerMark
	LoadepMsg        *TvVip // tv vip cfg
}

// EmptyArc def.
type EmptyArc struct {
	ChanSize   int64
	UnshelvePS int
}

// EsIdx def.
type EsIdx struct {
	UgcIdx, PgcIdx *EsCfg
}

// EsCfg def.
type EsCfg struct {
	Business string
	Index    string
}

// HisCfg def.
type HisCfg struct {
	Businesses []string
	Pagesize   int
}

// Dangbei cfg def.
type Dangbei struct {
	Pagesize int64          // dangbei api page size
	MangoPS  int            // mango page size
	Expire   xtime.Duration // dangbei page ID expiration
}

// AuthMsg configures the auth error messages
type AuthMsg struct {
	PGCOffline    string // offline pgc
	CMSInvalid    string // cms not valid
	LicenseReject string // license owner rejected
}

// App config
type App struct {
	*bm.App
}

func configCenter() (err error) {
	if client, err = conf.New(); err != nil {
		panic(err)
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

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf.
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
