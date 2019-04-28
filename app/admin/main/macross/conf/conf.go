package conf

import (
	"errors"
	"flag"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// config init.
var (
	confPath string
	Conf     *Config
)

// Config struct.
type Config struct {
	Version string `toml:"version"`
	// xlog
	Log *log.Config
	// HTTPServer
	BM *BM
	// tracer
	Tracer *trace.Config
	// oss
	Oss *Oss
	// db
	DB *DB
	// reload time
	Reload      time.Duration
	GitRelation map[string]string
	// identify
	Auth *permit.Config

	// other
	InvalidFrameworkFile []string
	TextSizeLimitList    []TextSizeLimit
	Property             *Property
}

// BM http.
type BM struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// TextSizeLimit struct.
type TextSizeLimit struct {
	Size  int64
	Limit float64
}

// Oss struct.
type Oss struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	Bucket          string
	OriginDir       string
	PublishDir      string
}

// DB struct.
type DB struct {
	Macross *sql.Config
}

// Property struct.
type Property struct {
	Mail *struct {
		Host    string
		Port    int
		Address string
		Pwd     string
		Name    string
	}

	Package *struct {
		URLPrefix string
		SavePath  string
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = configCenter()
	return
}

// configCenter ugc
func configCenter() (err error) {
	var (
		client *conf.Client
		c      string
		ok     bool
	)
	if client, err = conf.New(); err != nil {
		panic(err)
	}
	if c, ok = client.Toml2(); !ok {
		err = errors.New("load config center error")
		return
	}
	_, err = toml.Decode(c, &Conf)
	go func() {
		for e := range client.Event() {
			log.Error("get config from config center error(%v)", e)
		}
	}()
	return
}
