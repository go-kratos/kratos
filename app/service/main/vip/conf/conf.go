package conf

import (
	"errors"
	"flag"

	eleclient "go-common/app/service/main/vip/dao/ele-api-client"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	xtime "go-common/library/time"

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
	MsgURI string
	PayURI string
	VipURI string
	// log
	Log *log.Config
	// gorpc server
	RPCServer *rpc.ServerConfig
	// db
	Mysql *sql.Config
	// ecodes FIXME
	Ecode *ecode.Config
	//old db
	OldMysql *sql.Config
	// http client
	HTTPClient *bm.ClientConfig
	// mc
	Memcache *Memcache
	// pay conf
	PayConf *PayConf
	// rpc clients
	RPCClient2 *RPC
	// property
	Property *Property
	// http
	BM *bm.ServerConfig
	// redis
	Redis *Redis
	// associate conf
	AssociateConf *AssociateConf
	// ele conf
	ELEConf *eleclient.Config
	Host    *Host
	// grpc server
	WardenServer *warden.ServerConfig
	// grpc client
	CouponClient *warden.ClientConfig
}

// Host host.
type Host struct {
	Ele  string
	Mail string
}

// Redis redis
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
}

//PayConf pay config
type PayConf struct {
	CustomerID     int64
	Token          string
	OrderNotifyURL string
	SignNotifyURL  string
	PlanID         int32
	ProductID      string
	Version        string
	ReturnURL      string
	OrderExpire    int
	SignType       string
}

// RPC clients config.
type RPC struct {
	Member *rpc.ClientConfig
	Point  *rpc.ClientConfig
	Coupon *rpc.ClientConfig
}

// Property config for biz logic.
type Property struct {
	NotifyURL                      string
	MsgURL                         string
	PayURL                         string
	PayCoURL                       string
	AccountURL                     string
	PassportURL                    string
	APIURL                         string
	APICoURL                       string
	VipURL                         string
	TokenBID                       string
	PGCURL                         string
	ActiveDate                     string
	ActiveTip                      string
	Expire                         string
	AnnualVipBcoinDay              int16
	AnnualVipBcoinCouponMoney      int
	AnnualVipBcoinCouponActivityID int
	GiveBpDay                      int8
	PointGetRule                   map[string]int
	PointActiveDate                map[string]string
	BubbleTicker                   xtime.Duration
	PayType                        map[string]string
	PayChannelMapping              map[string]string
	PointBalance                   int64
	ActiveStart                    string
	ActiveEnd                      string
	ConfigMap                      map[string]string
	PointExchangeTitle             map[string]string
	WillExpiredTitle               map[string]string
	ExpiredTitle                   map[string]string
	TipButtonName                  string
	TipButtonLink                  string
	AllowanceSwitch                int8
	CodeSwitch                     int8
	GiveSwitch                     int8
	PanelBgURL                     string
	CodeOpenedSearchSize           int
	WelfareBgHost                  string
}

// AssociateConf associate vip conf.
type AssociateConf struct {
	// user grant count limit
	GrantDurationMap         map[string]int64 //限制饿了么发放联合会员的次数
	BilibiliPrizeGrantKeyMap map[string]string
	MailCouponID1            string           //票务优惠券满99减5
	MailCouponID2            string           //电商优惠券满299减20
	BilibiliBuyDurationMap   map[string]int64 //限制bilibili购买联合会员的次数
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
