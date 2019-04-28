package conf

import (
	"errors"
	"flag"

	"go-common/app/job/main/member/model/block"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	// base
	// app
	App *bm.App
	Biz *BizConfig
	// Env
	Env string
	// goroutine sleep
	Tick time.Duration
	// log
	Xlog *log.Config
	// databus
	DataBus *databus.Config
	// databus
	AccDataBus *databus.Config
	// passport
	PassortDataBus *databus.Config
	// operate log.
	LogDatabus *databus.Config
	// operate log for publish
	PLogDatabus *databus.Config
	// add exp,pub by archive,history.
	ExpDatabus *databus.Config
	// login award ,pub by passport.
	LoginDatabus    *databus.Config
	AwardDatabus    *databus.Config
	RealnameDatabus *databus.Config
	// account notify to purge cache
	AccountNotify   *databus.Config
	ShareMidDatabus *databus.Config
	// mc
	Memcache *Memcache
	// httpClinet
	HTTPClient    *bm.ClientConfig
	Mysql         *sql.Config
	AccCheckMysql *sql.Config
	AccMysql      *sql.Config
	AsoMysql      *sql.Config
	PasslogMysql  *sql.Config
	// hbase
	// HBase *conf.HBase
	// redis
	Redis        *redis.Config
	Databusutil  *databusutil.Config
	On           bool
	FeatureGates *FeatureGates
	SyncRange    *SyncRange
	// bm
	BM *bm.ServerConfig
	// Report
	UserReport    *databus.Config
	ManagerReport *databus.Config

	//block config.
	BlockMemcache      *memcache.Config
	BlockDB            *sql.Config
	BlockCreditDatabus *databus.Config
	BlockProperty      *Property

	//realname
	RealnameRsaPriv        []byte
	RealnameAlipayPub      []byte
	RealnameAlipayBiliPriv []byte

	// Parsed Realname Infoc
	ParsedRealnameInfoc *infoc.Config
}

// Property .
type Property struct {
	LimitExpireCheckLimit  int
	LimitExpireCheckTick   time.Duration
	CreditExpireCheckLimit int
	CreditExpireCheckTick  time.Duration
	MSGURL                 string
	MSG                    *MSG
	Flag                   *struct {
		ExpireCheck bool
		CreditSub   bool
	}
}

// MSG .
type MSG struct {
	BlockRemove block.MSG
}

// Memcache memcache.
type Memcache struct {
	*memcache.Config
	Expire time.Duration
}

// SyncRange syncRange.
type SyncRange struct {
	Start int64
	End   int64
}

// BizConfig biz common config
type BizConfig struct {
	ExpiredBase   int32
	ExpiredDetail int32
	IsFree        bool
	AccprocCount  int32
	ExpprocCount  int32

	// realname alipay
	RealnameAlipayCheckTick  time.Duration
	RealnameAlipayCheckLimit int
	RealnameAlipayAppID      string
	RealnameAlipayGateway    string
}

// FeatureGates is.
type FeatureGates struct {
	DataFixer bool
	FaceCheck bool
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init create config instance.
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
	// go func() {
	// 	for range client.Event() {
	// 		log.Info("config reload")
	// 		if load() != nil {
	// 			log.Error("config reload error (%v)", err)
	// 		}
	// 	}
	// }()
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
	tmpConf.RealnameRsaPriv, tmpConf.RealnameAlipayPub, tmpConf.RealnameAlipayBiliPriv = loadRealnameKey()
	*Conf = *tmpConf
	return
}

func loadRealnameKey() (rasPriv, alipayPub, alipayBiliPriv []byte) {
	var (
		emptyBytes = []byte("")
	)
	rasPriv, alipayPub, alipayBiliPriv = emptyBytes, emptyBytes, emptyBytes
	if client == nil {
		return
	}
	if str, ok := client.Value("realname.rsa.priv"); ok {
		rasPriv = []byte(str)
	}
	if str, ok := client.Value("realname.alipay.pub"); ok {
		alipayPub = []byte(str)
	}
	if str, ok := client.Value("realname.alipay.bili.priv"); ok {
		alipayBiliPriv = []byte(str)
	}
	return
}
