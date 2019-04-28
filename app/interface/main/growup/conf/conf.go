package conf

import (
	"errors"
	"flag"
	"path"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"go-common/library/database/hbase.v2"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath local config path
	ConfPath string
	// Conf config
	Conf   = &Config{}
	client *conf.Client
)

// Config str
type Config struct {
	// base
	// channal len
	ChanSize int64
	// log
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// identify
	AppConf *bm.App
	// tracer
	Tracer *trace.Config
	// tick load pgc
	Tick time.Duration
	// orm
	DB *DB
	//redis
	Redis *Redis
	// http client of search
	HTTPClient *HTTPClient
	// host
	Host *Host
	// rpc
	ArticleRPC *rpc.ClientConfig
	AccountRPC *rpc.ClientConfig
	VipRPC     *rpc.ClientConfig
	AccCliConf *warden.ClientConfig
	// threshold
	Threshold *Threshold
	// hbase
	HBase *HBaseConfig
	// newbie
	Newbie Newbie
}

// DB def db struct
type DB struct {
	Allowance *sql.Config
	Growup    *sql.Config
}

// Redis define redis conf.
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// HTTPClient http client
type HTTPClient struct {
	Read *bm.ClientConfig
}

// Host http host
type Host struct {
	AccountURI    string
	ArchiveURI    string
	UperURI       string
	VipURI        string
	ActivitiesURI string
	RelationsURI  string
	VideoUpURI    string
	CategoriesURI string
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// Threshold up cretive threshold
type Threshold struct {
	LimitFanCnt      int64
	LimitTotalClick  int64
	LimitArticleView int64
}

// Newbie newbie config
type Newbie struct {
	Talents              map[string]string
	DefaultCover         string
	DefaultTalent        string
	RecommendUpCount     int
	RecommendUpPoolCount int
	ActivityCount        int
	ActivityShotType     int32
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init init conf
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
	if err != nil {
		return
	}
	var templateConfPath = path.Join(path.Dir(ConfPath), "newbie.toml")
	_, err = toml.DecodeFile(templateConfPath, &Conf)
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
