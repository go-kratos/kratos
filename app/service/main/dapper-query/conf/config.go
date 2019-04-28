package conf

import (
	"errors"
	"flag"

	"github.com/BurntSushi/toml"

	"go-common/library/conf"
	"go-common/library/log"
	xtime "go-common/library/time"
)

func init() {
	flag.StringVar(&confPath, "conf", "", "config file")
}

var (
	confPath string
	// Conf conf
	Conf   = &Config{}
	client *conf.Client
)

// Config config.
type Config struct {
	Log        *log.Config     `toml:"log"`
	HBase      *HBaseConfig    `toml:"hbase"`
	InfluxDB   *InfluxDBConfig `toml:"influx_db"`
	OpsLog     *OpsLog         `toml:"ops_log"`
	Collectors *Collectors     `toml:"collectors"`
}

// InfluxDBConfig InfluxDBConfig
type InfluxDBConfig struct {
	Addr     string `toml:"addr"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
}

// HBaseConfig hbase config
type HBaseConfig struct {
	Namespace           string         `toml:"namespace"`
	Addrs               string         `toml:"addrs"`
	RPCQueueSize        int            `toml:"rpc_queue_size"`
	FlushInterval       xtime.Duration `toml:"flush_interval"`
	EffectiveUser       string         `toml:"effective_user"`
	RegionLookupTimeout xtime.Duration `toml:"region_lookup_timeout"`
	RegionReadTimeout   xtime.Duration `toml:"region_read_timeout"`
}

// OpsLog .
type OpsLog struct {
	API string `toml:"api"`
}

// Collectors collector config
type Collectors struct {
	Nodes []string `toml:"nodes"`
}

// Init config
func Init() (err error) {
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
			if err := load(); err != nil {
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
