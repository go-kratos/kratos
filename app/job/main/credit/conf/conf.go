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
	"go-common/library/queue/databus"
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
	// host
	Host *Host
	// Env
	Env string
	// goroutine sleep
	Tick time.Duration
	// log
	Xlog *log.Config
	// databus
	DataBus *DataBus
	// bm service
	BM *bm.ServerConfig
	// httpClinet
	HTTPClient *bm.ClientConfig
	Mysql      *sql.Config
	Judge      *Judge
	Sms        *Sms
	Redis      *Redis
	Memcache   *Memcache
	// rpc client
	RPCClient *RPC
}

// RPC rpc client config.
type RPC struct {
	Archive *rpc.ClientConfig
	Member  *rpc.ClientConfig
}

// DataBus databus config.
type DataBus struct {
	CreditDBSub *databus.Config
	ReplyAllSub *databus.Config
	LabourSub   *databus.Config
}

// Host config host.
type Host struct {
	APICoURI     string
	AccountCoURI string
	MsgCoURI     string
}

// Redis redis conf.
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Sms is sms monitor config.
type Sms struct {
	Phone string
	Token string
}

// Judge is judge config.
type Judge struct {
	ConfTimer          time.Duration // 定时load数据时间间隔
	ReservedTime       time.Duration // 结案前N分钟停止获取case
	CaseGiveHours      int64         // 案件发放时长
	CaseCheckTime      int64         // 单案审核时长
	JuryRatio          int64         // 投准率下限
	JudgeRadio         int64         // 判决阙值
	CaseVoteMin        int64         // 案件投票数下限
	CaseObtainMax      int64         // 每日获取案件数
	CaseVoteMax        int64         // 结案投票数
	JuryApplyMax       int64         // 每日发放风纪委上限
	CaseLoadMax        int           // 案件发放最大队列数
	CaseLoadSwitch     int8          // 案件发放进入队列开关
	CaseVoteMaxPercent int           // 结案投票数的百分比
	VoteNum
}

// VoteNum struct.
type VoteNum struct {
	RateS int8 `json:"rate_s"`
	RateA int8 `json:"rate_a"`
	RateB int8 `json:"rate_b"`
	RateC int8 `json:"rate_c"`
	RateD int8 `json:"rate_d"`
}

// Memcache memcache.
type Memcache struct {
	*memcache.Config
	Expire time.Duration
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
