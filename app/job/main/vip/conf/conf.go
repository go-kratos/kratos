package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
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
	VipURI string
	// base
	App *bm.App
	// log
	Xlog *log.Config
	// http
	BM *bm.ServerConfig
	// db
	NewMysql *sql.Config
	//old db
	OldMysql *sql.Config
	// http client
	HTTPClient *bm.ClientConfig
	//Property
	Property *Property

	URLConf *URLConf
	//databus group config
	DatabusUtil *databusutil.Config
	//point databus
	Databus *DataSource
	// mc
	Memcache *Memcache
	// redis
	Redis   *Redis
	PayConf *PayConf
	// rpc clients
	RPCClient2 *RPC
	// grpc
	VipClient *warden.ClientConfig
}

//RPC rpc clients.
type RPC struct {
	Coupon *rpc.ClientConfig
}

//URLConf url conf
type URLConf struct {
	PayCoURL    string
	PayURL      string
	MsgURL      string
	MallURL     string
	AccountURL  string
	APICoURL    string
	OldVipCoURL string
}

// Property config for biz logic.
type Property struct {
	UpdateUserInfoCron               string
	AutoRenewCron                    string
	SendMessageCron                  string
	SendBcoinCron                    string
	WillExpireMsgCron                string
	HadExpiredMsgCron                string
	PushDataCron                     string
	EleEompensateCron                string
	HandlerThread                    int
	ReadThread                       int
	Retry                            int
	FrozenExpire                     xtime.Duration
	FrozenDate                       xtime.Duration
	FrozenLimit                      int64
	FrozenCron                       string
	PayMapping                       map[string]string
	MsgURL                           string
	ActivityID                       int64
	AnnualVipBcoinDay                int
	AnnualVipBcoinCouponMoney        int
	PayCoURL                         string
	SalaryDay                        int
	AnnualVipSalaryCount             int
	NormalVipSalaryCount             int
	SalaryVideoCouponnIterval        xtime.Duration
	SalaryVideoCouponCron            string
	MsgOpen                          bool
	BatchSize                        int
	SalaryCouponMaps                 map[string]map[string]int64 // map[coupontype]map[viptype]salarycount
	SalaryCouponTypes                []int8
	SalaryCouponBatchNoMaps          map[string]string // map[coupontype]batchnofmt
	SalaryCouponMsgTitleMaps         map[string]string // map[ coupontype]msgTitle
	SalaryCouponMsgContentMaps       map[string]string // map[coupontype]msgsContent
	SalaryCouponMsgSupplyContentMaps map[string]string // map[coupontype]msgsContent
	SalaryCouponURL                  string
	ActiveStartTime                  string
	SendMedalEndTime                 string
	SendVipbuyEndTime                string
	SummerActiveStartTime            string
	SummerActiveEndTime              string
	SendCodeStartTime                string
	SendCodeEndTime                  string
	CouponIDs                        []string
	MedalID                          int64
	CodeExchangeMap                  map[string][]int64
	CodeExchangeTimeMap              map[string]int
	CodeExchangePicMap               map[string]string
	VipbuyExchangeNameMap            map[string]string
	GrayScope                        int64
	PushToken                        string
	BusinessID                       int64
	SplitPush                        int
	UpdateDB                         bool
	NotGrantLimit                    int
}

//PayConf pay conf info
type PayConf struct {
	BasicURL       string
	CustomerID     string
	Token          string
	NotifyURL      string
	OrderNotifyURL string
	SignNotifyURL  string
	PlanID         int32
	ProductID      string
	Version        string
}

// Memcache memcache
type Memcache struct {
	*memcache.Config
	Expire xtime.Duration
}

// Redis redis
type Redis struct {
	*redis.Config
	Expire xtime.Duration
}

// DataSource databus config zone.
type DataSource struct {
	AccLogin      *databus.Config
	OldVipBinLog  *databus.Config
	SalaryCoupon  *databus.Config
	NewVipBinLog  *databus.Config
	AccountNotify *databus.Config
	CouponNotify  *databus.Config
	AutoRenew     *databus.Config
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
