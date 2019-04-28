package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// base
	// elk
	Log *log.Config
	// http
	BM *bm.ServerConfig
	// identify
	//Identify *identify.Config
	// tracer
	Tracer *trace.Config
	// MySQL
	MySQL *sql.Config
	// MySQL
	OldMySQL *sql.Config
	// Databus
	Databus *databus.Config
	// AuthDataBus authService binlog
	AuthDataBus *databus.Config
	// DataUtil config
	DatabusUtil *databusutil.Config
	// auth rpc
	AuthRPC *rpc.ClientConfig
	// HTTPClientConfig
	HTTPClientConfig *bm.ClientConfig
	// AuthJobConfig job config
	AuthJobConfig *AuthJobConfig
	// user defined
	SyncLines int64
	IDXFrom   int64
	IDXTo     int64
}

// AuthJobConfig auth job config
type AuthJobConfig struct {
	AsoCleanURL string
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
	if s, ok = client.Value("passport-auth-job.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
