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
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
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
	// log
	Xlog *log.Config
	// AuthN
	AuthN *auth.Config
	// Verify
	Verify *verify.Config
	// http
	BM *HTTPServers

	// db
	Mysql *sql.Config
	// mc
	Memcache *Memcache
	//redis
	Redis *Redis
	// databus
	DataBus *databus.Config
	// tracer
	Tracer *trace.Config

	// bm client
	HTTPClient *bm.ClientConfig
	// Judge conf
	Judge *Judge
	// ecodes
	Ecode *ecode.Config
	// Antispam
	Antispam *antispam.Config
	// TagID
	TagID *TagID
	// Property
	Property *Property
	// rpc client
	RPCClient2 *RPC
	//  GRPCClient
	GRPCClient *GRPC
}

// Redis define redis conf.
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Memcache define memcache conf.
type Memcache struct {
	*memcache.Config
	UserExpire      time.Duration
	MinCommonExpire time.Duration
	CommonExpire    time.Duration
}

// Host define host conf.
type Host struct {
	MessageURI  string
	BigDataURI  string
	APICoURI    string
	ManagersURI string
}

// Judge define judge conf.
type Judge struct {
	ConfTimer       time.Duration // 定时load数据时间间隔
	ReservedTime    time.Duration // 结案前N分钟停止获取case
	LoadManagerTime time.Duration // load manager user的时间间隔
	CaseGiveHours   int64         // 案件发放时长
	CaseCheckTime   int64         // 单案审核时长
	JuryRatio       int64         // 投准率下限
	JudgeRadio      int64         // 判决阙值
	CaseVoteMin     int64         // 案件投票数下限
	CaseObtainMax   int64         // 每日获取案件数
	CaseVoteMax     int64         // 结案投票数
	JuryApplyMax    int64         // 每日发放风纪委上限
	CaseLoadMax     int           // 案件发放最大队列数
	CaseLoadSwitch  int8          // 案件发放进入队列开关
	VoteNum
}

// VoteNum .
type VoteNum struct {
	RateS int8 `json:"rate_s"`
	RateA int8 `json:"rate_a"`
	RateB int8 `json:"rate_b"`
	RateC int8 `json:"rate_c"`
	RateD int8 `json:"rate_d"`
}

// RPC rpc client config.
type RPC struct {
	Archive *rpc.ClientConfig
	Member  *rpc.ClientConfig
}

// GRPC .
type GRPC struct {
	Filter  *warden.ClientConfig
	Account *warden.ClientConfig
}

// HTTPServers Http Servers
type HTTPServers struct {
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

// TagID is workflow id .
type TagID struct {
	Reply    int64
	DM       int64
	Msg      int64
	Tag      int64
	Member   int64
	Archive  int64
	Music    int64
	Article  int64
	SpaceTop int64
}

// Property .
type Property struct {
	QsNum    int
	PerScore int64
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
