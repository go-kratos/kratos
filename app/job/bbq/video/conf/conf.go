package conf

import (
	"errors"
	"flag"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"go-common/library/net/rpc/warden"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log          *log.Config
	BM           *HTTPGeneral
	Verify       *verify.Config
	Tracer       *trace.Config
	Redis        *redis.Config
	BfRedis      *redis.Config
	Memcache     *memcache.Config
	MySQL        *sql.Config
	MySQLCms     *sql.Config
	MySQLOffline *sql.Config
	Ecode        *ecode.Config
	Scheduler    *Scheduler
	GRPCClient   map[string]*GRPCConf
	//berserker
	Berserker     *Berserker
	TagMap        *TagMap
	Databus       map[string]*databus.Config
	Mail          *Mail
	URLs          map[string]string
	Download      *Download
	FTP           *FTP
	VST           *Vst
	Path          map[string]string
	SubBvcControl *Sub
}

//Sub ...
type Sub struct {
	Control int8
}

//Vst ...
type Vst struct {
	TmpStatus int64
}

// FTP FTP.
type FTP struct {
	Addr       string
	User       string
	Password   string
	RemotePath map[string]string
	Timeout    xtime.Duration
	LocalPath  map[string]string
}

//GRPCConf .
type GRPCConf struct {
	WardenConf *warden.ClientConfig
	Addr       string
}

// Download .
type Download struct {
	File string
}

//Mail ...
type Mail struct {
	Host     string
	Port     int
	From     string
	Password string
	To       []string
}

//TagMap ...
type TagMap struct {
	TagTidMap    map[string]string
	TagSubTidMap map[string]string
}

// HTTPGeneral conf
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
	Rankdaily        string
	Userdmg          string
	Upuserdmg        string
	Operaonce        string
	Userbasic        string
	Upmid            string
	VideoView        string
	UserProfile      string
	UserProfileBuvid string
}

// BerSerkerKeyList conf
type BerSerkerKeyList struct {
	YYC *BerSerkerKey
	HSC *BerSerkerKey
	LZQ *BerSerkerKey
	LJ  *BerSerkerKey
	DW  *BerSerkerKey
	HM  *BerSerkerKey
}

// BerSerkerKey conf
type BerSerkerKey struct {
	Appkey string
	Secret string
}

//Scheduler .
type Scheduler struct {
	CheckVideo2ES         string
	SyncUserDmg           string
	Test                  string
	SyncUpUserDmg         string
	CheckVideo            string
	CheckVideoSt          string
	CheckVideoStHv        string
	CheckVideoTag         string
	CheckTag              string
	SyncVideoOper         string
	DeliveryNewVideoToCms string
	SyncUsrSta            string
	SyncSearch            string
	VideoViewHistory      string
	SysMsgTask            string
	UserProfileBbq        string
	TransToReview         string
	TransToCheckBack      string
}

// Databus .
type Databus struct {
	Video    *databus.Config
	VideoRep *databus.Config
	BvcSub   *databus.Config
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
	*Conf = *tmpConf
	return
}
