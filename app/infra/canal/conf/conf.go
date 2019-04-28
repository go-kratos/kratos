package conf

import (
	"errors"
	"flag"
	"time"

	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/naming/discovery"
	bm "go-common/library/net/http/blademaster"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// ConfClient get config client
	ConfClient *conf.Client
	// Conf canal config variable
	Conf = &Config{}
)

// Config canal config struct
type Config struct {
	Monitor *Monitor
	// xlog
	Log *log.Config
	// http client
	HTTPClient *bm.ClientConfig
	// http server
	BM *bm.ServerConfig
	// master info
	MasterInfo *MasterInfoConfig
	// discovery
	Discovery *discovery.Config
	// db
	DB *sql.Config
}

// Monitor wechat monitor
type Monitor struct {
	User   string
	Token  string
	Secret string
}

// MasterInfoConfig save pos of binlog in file or db
type MasterInfoConfig struct {
	Addr     string        `toml:"addr"`
	DBName   string        `toml:"dbName"`
	User     string        `toml:"user"`
	Password string        `toml:"password"`
	Timeout  time.Duration `toml:"timeout"`
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
	flag.StringVar(&canalPath, "canal", "", "canal instance path")
}

//Init int config
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	return remote()
}

func remote() (err error) {
	if ConfClient, err = conf.New(); err != nil {
		return
	}
	ConfClient.WatchAll()
	err = LoadCanal()
	return
}

// LoadCanal canal config
func LoadCanal() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = ConfClient.Value("canal.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
