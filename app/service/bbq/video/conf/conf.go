package conf

import (
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf/paladin"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	localConf string
	confName  string
	// Conf config
	Conf = &Config{}
	// ArchiveRules 稿件导入规则
	ArchiveRules = &Rules{}
)

// Config .
type Config struct {
	Log        *log.Config
	BM         *HTTPGeneral
	Verify     *verify.Config
	Tracer     *trace.Config
	Redis      *redis.Config
	MySQL      *sql.Config
	CMSMySQL   *sql.Config
	Ecode      *ecode.Config
	Berserker  *Berserker
	GRPCClient map[string]*GRPCConf
	URLs       map[string]string
	Databus    map[string]*databus.Config
	BPSCode    map[string]map[string]int64
	Upload     *Upload
}

//Upload .
type Upload struct {
	File     *UploadFile
	Endpoint *UploadEndPoint
	Auth     *UploadAuth
}

//UploadFile .
type UploadFile struct {
	Prefix string
	Line   string
}

//UploadEndPoint .
type UploadEndPoint struct {
	Main   string
	BackUp string
}

//UploadAuth .
type UploadAuth struct {
	AK string
	SK string
}

//GRPCConf .
type GRPCConf struct {
	WardenConf *warden.ClientConfig
	Addr       string
}

//HTTPGeneral ...
type HTTPGeneral struct {
	Server *bm.ServerConfig
	Client *bm.ClientConfig
}

// Berserker conf
type Berserker struct {
	Key *BerSerkerKeyList
	API *BerserkerAPI
}

// BerserkerAPI conf
type BerserkerAPI struct {
	Rankdaily string
	Userdmg   string
	Operaonce string
}

// BerSerkerKeyList conf
type BerSerkerKeyList struct {
	YYC *BerSerkerKey
	HSC *BerSerkerKey
	LZQ *BerSerkerKey
}

// BerSerkerKey conf
type BerSerkerKey struct {
	Appkey string
	Secret string
}

// Set .
func (c *Config) Set(text string) error {
	if _, err := toml.Decode(text, c); err != nil {
		panic(err)
	}
	return nil
}

// Set .
func (r *Rules) Set(text string) error {
	if _, err := toml.Decode(text, r); err != nil {
		panic(err)
	}
	return nil
}

func init() {
	//线下使用
	flag.StringVar(&localConf, "localconf", "", "default config path")
	flag.StringVar(&confName, "conf_name", "video-service.toml", "default config filename")
}

// Init init conf
func Init() error {
	if localConf != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(localConf, &Conf)
	return
}

func remote() (err error) {
	if err := paladin.Init(); err != nil {
		panic(err)
	}
	// var setter
	if err := paladin.Watch(confName, Conf); err != nil {
		panic(err)
	}

	if err := paladin.Watch("rule.toml", ArchiveRules); err != nil {
		panic(err)
	}
	return
}
