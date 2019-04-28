package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf of config
	Conf = &Config{}
)

// Memcache memcache.
type Memcache struct {
	*memcache.Config
	CmsExpire xtime.Duration
}

// Config def.
type Config struct {
	// base
	// http
	HTTPServer *bm.ServerConfig
	// auth
	Auth *permit.Config
	// db
	ORM *orm.Config
	// dbshow
	ORMShow *orm.Config
	// log
	Log *log.Config
	// tracer
	Tracer *trace.Config
	// httpsearch
	HTTPSearch *HTTPSearch
	// Cfg
	Cfg *Cfg
	// HTTPClient .
	HTTPClient *bm.ClientConfig
	// URLConf
	URLConf *URLConf
	// YSTParam
	YSTParam *YSTParam
	// Bfs
	Bfs *Bfs
	// grpc
	ArcClient *warden.ClientConfig
	AccClient *warden.ClientConfig
	// Ecode Cfg
	Ecode *ecode.Config
	//mc
	Memcache *Memcache
}

// Cfg def
type Cfg struct {
	Playpath      string // playurl
	AuditRSize    int    // page size for audit result checking
	PlayurlAPI    string // pgc playurl api
	SearInterMax  int
	IntervLimit   int
	PGCTypes      []string          // pgc types that need to filter archives
	PgcNames      map[string]string // pgc category name in CN
	ModIntMaxSize int               // ModIntMaxSize module intervene max size
	TypesLoad     xtime.Duration    // reloading type duratio
	UPlayurlAPI   string            // ugc playurl api
	SupportCat    *SupportCat
	MangoErr      string // mango error indication message
	LoadSnFre     xtime.Duration
	RefLabel      *RefLabel     // refresh label original data frequency
	AuditConsult  *AuditConsult // audit consult cfg
	Hosts         *Hosts
	Abnormal      *Abnormal // abnormal cid export related cfg
	EsIdx         *EsIdx    // es index cfg
}

// EsIdx def.
type EsIdx struct {
	UgcIdx *EsCfg
}

// EsCfg def.
type EsCfg struct {
	Business string
	Index    string
}

// RefLabel def.
type RefLabel struct {
	Fre      xtime.Duration
	PgcAPI   string // pgc api host
	UgcType  string
	UgcTime  string
	AllValue string
	AllName  string
}

// Hosts def.
type Hosts struct {
	ESUgc   string // ESUgc api
	Manager string // manager host
}

// AuditConsult related cfg
type AuditConsult struct {
	LikeLimit  int
	UnshelveNb int
	MatchPS    int64
}

// Abnormal cid export def
type Abnormal struct {
	CriticalCid  int64 // 12780000, critical cid for transcoding
	AbnormHours  int   // ugc abnormal cid interval hour
	ReloadFre    xtime.Duration
	ExportTitles []string // export titles
}

// SupportCat means the pgc&ugc types that we support to fill the modules
type SupportCat struct {
	PGCTypes  []int32
	UGCTypes  []int32
	ReloadFre xtime.Duration
}

// URLConf url conf
type URLConf struct {
	GetRemotePanelUrl string
	SyncPanelUrl      string
}

// YSTParam yst config param
type YSTParam struct {
	QueryPanelType  string
	InsertPanelType string
	Source          string
	Insert          string
	Update          string
}

// Bfs struct
type Bfs struct {
	Key     string
	Secret  string
	Host    string
	Timeout int
	Bucket  string
}

// HTTPSearch http client of search
type HTTPSearch struct {
	*bm.ClientConfig
	FullURL string
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

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}
