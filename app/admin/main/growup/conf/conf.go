package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

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
	// identify
	App *bm.App
	// tracer
	Tracer *trace.Config
	// tick load pgc
	Tick time.Duration
	// orm
	ORM *ORM
	// host
	Host *Host
	// Bfs
	Bfs *Bfs
	// http client
	HTTPClient *bm.ClientConfig
	// image client
	ImageClient *bm.ClientConfig
	// BM HTTPServers
	BM *bm.ServerConfig
	// budget
	Budget *Budget
	// rpc client
	VipRPC *rpc.ClientConfig
	// grpc client
	Account *warden.ClientConfig
	// shell config
	ShellConf *ShellConfig
	OtherConf *OtherConfig
}

// ORM is orm config
type ORM struct {
	Allowance *sql.Config
	Growup    *orm.Config
}

// Host is hosts
type Host struct {
	Message    string
	Common     string
	VideoType  string
	ColumnType string
	Creative   string
	API        string
}

// Bfs struct.
type Bfs struct {
	Addr        string
	Bucket      string
	Key         string
	Secret      string
	MaxFileSize int
}

// Budget config.
type Budget struct {
	Video  *BBudget
	Column *BBudget
	Bgm    *BBudget
}

// BBudget config.
type BBudget struct {
	Year         int64
	AnnualBudget int64
	DayBudget    int64
}

//ShellConfig 贝壳系统配置
type ShellConfig struct {
	CustomID    string
	Token       string
	PayHost     string
	CallbackURL string
}

//OtherConfig 其他配置
type OtherConfig struct {
	// true需要consume数据，false不consume数据
	OfflineOrderConsume bool
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
