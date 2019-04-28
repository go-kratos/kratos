package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath local config path
	confPath string
	client   *conf.Client
	// Conf is global config object.
	Conf = &Config{}
)

// Config is project all config
type Config struct {
	// base
	// log
	Log *log.Config
	// Mysql
	Mysql *Mysql
	// tracer
	Tracer *trace.Config
	// key secret
	KeySecret *KeySecret
	// mail
	Mail *Mail
	// avratio
	Ratio *TagConf
	// upincome
	Income *TagConf
	// http client
	HTTPClient *bm.ClientConfig
	// berserker client
	DPClient *bm.ClientConfig
	// Host
	Host *Host
	// bm
	BM *bm.ServerConfig
	// hbase
	HBase *HBaseConfig
	// databus
	ArchiveSub *databus.Config
	//
	Bubble *BubbleConfig
}

// KeySecret.
type KeySecret struct {
	Key    string
	Secret string
}

// Mysql mysql config
type Mysql struct {
	Growup    *sql.Config
	Allowance *sql.Config
}

// Mail config
type Mail struct {
	Host               string
	Port               int
	Username, Password string
	Send               []*MailAddr
}

// MailAddr mail send addr.
type MailAddr struct {
	Type int
	Addr []string
}

// Host is hosts
type Host struct {
	Archive      string
	VideoType    string
	ColumnType   string
	DataPlatform string
	ColumnAct    string
	Profit       string
	VC           string
	Porder       string
	Archives     string
	API          string
}

// TagConf tag config
type TagConf struct {
	Hour  int
	Sleep int64
	Num   int64
	Limit int64
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// BubbleConfig.
type BubbleConfig struct {
	BRatio []*BRatio
}

// BRatio .
type BRatio struct {
	BType int
	Ratio float64
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init config.
func Init() (err error) {
	if confPath != "" {
		_, err = toml.DecodeFile(confPath, &Conf)
		return
	}
	err = remote()
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
