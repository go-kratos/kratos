package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/bfs"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

//Conf config var
var (
	ConfPath string
	client   *conf.Client
	Conf     = &Config{}
)

//Config Config
type Config struct {
	// base
	// ecode
	Ecode *ecode.Config
	// log
	Xlog     *log.Config
	Infoc2   *infoc.Config
	InfocBak *infoc.Config
	Verify   *verify.Config
	// rpc client
	AccountRPC *warden.ClientConfig
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	// http server config
	HTTPServer *bm.ServerConfig
	// db
	DB *DB
	// memcache
	Memcache *Memcache
	// http client
	HTTPClient *HTTPClient
	// http client of search
	HTTPSearch *HTTPSearch
	// http client of Infoc Data
	HTTPInfoc *HTTPInfoc
	// tracer
	Tracer *trace.Config
	// databus client
	ActionPub   *databus.Config
	Host        Host
	ManagerAuth *permit.Config
	// elastic config
	Elastic *elastic.Config
	// manager log config
	ManagerLog *databus.Config
	//bfs config
	BFS *bfs.Config
}

// DB mysql config struct
type DB struct {
	DM           *sql.Config
	DMMetaWriter *sql.Config
	DMMetaReader *sql.Config
}

// Memcache dm memcache
type Memcache struct {
	Filter *struct {
		*memcache.Config
	}
	Subtitle *struct {
		*memcache.Config
	}
}

// DMReport dm report
type DMReport struct {
	Count    int64
	MsgURL   string
	MoralURL string
	BlockURL string
}

// Host hosts used in dm admin
type Host struct {
	Videoup   string
	API       string
	Search    string
	Season    string
	Message   string
	Account   string
	Mask      string
	Berserker string
}

// HTTPClient http client
type HTTPClient struct {
	*bm.ClientConfig
	JudgeURL    string
	ArcInfoURL  string
	CidAidsURL  string
	TypeInfoURL string
}

// HTTPSearch http client of search
type HTTPSearch struct {
	*bm.ClientConfig
	SearchURL    string
	UpdateURL    string
	ReportURL    string
	UpdateOldURL string
	CountURL     string
}

// HTTPInfoc http client of Infoc query
type HTTPInfoc struct {
	*bm.ClientConfig
	InfocQueryURL string
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "config path")
}

//Init int config
func Init() error {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
