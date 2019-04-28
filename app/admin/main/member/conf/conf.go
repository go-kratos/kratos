package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/hbase.v2"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

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
	Log           *log.Config
	BM            *bm.ServerConfig
	Tracer        *trace.Config
	ORM           *ORM
	ORMRead       *ORM
	Ecode         *ecode.Config
	Host          *Host
	HTTPClient    *HTTPClient
	HBase         *hbase.Config
	Auth          *permit.Config
	FHByOPHBase   *hbase.Config
	FHByMidHBase  *hbase.Config
	FacePriBFS    *BFS
	FaceBFS       *BFS
	ManagerReport *databus.Config
	ExpMsgDatabus *databus.Config
	RPCClient     *RPC
	ES            *elastic.Config
	Realname      *Realname
	ReviewNotify  *ReviewNotify
	Redis         *redis.Config
	Memcache      *memcache.Config
	// block
	AccountNotify *databus.Config
	BlockProperty *Property
	BlockMemcache *memcache.Config
	BlockMySQL    *sql.Config
}

// Host is Host config
type Host struct {
	Message  string
	Passport string
	Merak    string
}

// Property .
type Property struct {
	BlackHouseURL string
	MSGURL        string
	TelURL        string
	MailURL       string
}

// RPC config
type RPC struct {
	Coin     *rpc.ClientConfig
	Account  *warden.ClientConfig
	Figure   *rpc.ClientConfig
	Member   *rpc.ClientConfig
	Spy      *rpc.ClientConfig
	Relation *rpc.ClientConfig
}

// ORM is database config
type ORM struct {
	Member     *orm.Config
	MemberRead *orm.Config
	Account    *orm.Config
}

// BFS bfs config
type BFS struct {
	Timeout     time.Duration
	MaxFileSize int
	Bucket      string
	URL         string
	Key         string
	Secret      string
}

// Realname conf
type Realname struct {
	ImageURLTemplate string
	DataDir          string
	RsaPub           []byte
	RsaPriv          []byte
}

// HTTPClient http client
type HTTPClient struct {
	Read     *bm.ClientConfig
	Passport *bm.ClientConfig
}

// ReviewNotify notify users
type ReviewNotify struct {
	Users []string
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
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	pub, priv := loadRealnameKey()
	tmpConf.Realname.RsaPriv = []byte(priv)
	tmpConf.Realname.RsaPub = []byte(pub)
	*Conf = *tmpConf
	return
}

func loadRealnameKey() (string, string) {
	priv, ok := client.Value("realname.rsa.priv")
	if !ok {
		panic("Failed to load realname private key")
	}
	pub, ok := client.Value("realname.rsa.pub")
	if !ok {
		panic("Failed to load realname pubic key")
	}
	return pub, priv
}
