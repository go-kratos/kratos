package conf

import (
	"errors"
	"flag"
	"go-common/library/net/rpc/warden"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/queue/databus"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf is global config object.
	Conf = &Config{}
)

// Config is project all config.
type Config struct {
	Host *Host
	// interface XLog
	XLog  *log.Config
	Redis *redis.Config
	// db
	DB *DB
	// http client
	HTTPClient *blademaster.ClientConfig
	// databus sub
	Bvc2VuSub     *databus.Config
	VideoupSub    *databus.Config
	ArcResultSub  *databus.Config
	StatSub       *databus.Config
	VideoshotSub2 *databus.Config
	// databus pub
	VideoupPub *databus.Config
	BlogPub    *databus.Config
	CheckFrPub *databus.Config
	// rpc
	AccRPC *warden.ClientConfig
	// mail
	Mail          *Mail
	Bm            *blademaster.ServerConfig
	ManagerReport *databus.Config
	// others
	SpecialUp   []int64
	Tels        string
	ChangeDebug bool
	ChangeMid   int64
	//灰度控制
	Debug    bool
	DebugMid int64

	BvcConsumeTimeout int64
}

// Host is hosts
type Host struct {
	Message  string
	API      string
	Monitor  string
	Act      string
	Bvc      *BVC
	Push     *PushC
	RecCover string
}

//BVC key
type BVC struct {
	Bvc       string
	GapKey    string
	AppendKey string
}

//PushC 创作姬app推送消息配置
type PushC struct {
	AppID      string
	BusinessID string
	Token      string
}

// DB is db config.
type DB struct {
	Archive     *sql.Config
	ArchiveRead *sql.Config
	Manager     *sql.Config
}

//Mail 邮件配置
type Mail struct {
	Host, Checkout     string
	Port               int
	Username, Password string
	Addr               []*MailElemenet
	PrivateAddr        []*MailElemenet
}

//MailElemenet 邮件接收人配置
type MailElemenet struct {
	Type string
	Desc string
	Addr []string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
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
