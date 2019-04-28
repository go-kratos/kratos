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
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf .
	Conf = &Config{}
)

// Config .
type Config struct {
	Log         *log.Config
	Tracer      *trace.Config
	Ecode       *ecode.Config
	Tag         *Tag
	Supervision *Supervision
	URL         *URL
	BM          *bm.ServerConfig
	RPCServer   *rpc.ServerConfig
	MySQL       *sql.Config
	Memcache    *Memcache
	Redis       *Redis
	RPC         *RPC
	GRPC        *warden.ServerConfig
	Databus     *databus.Config
}

// URL .
type URL struct {
	APICo   string
	Account string
}

// RPC rpc.
type RPC struct {
	Archive *rpc.ClientConfig
}

// Redis redis.
type Redis struct {
	*redis.Config
	Expire        time.Duration
	SubExpire     time.Duration
	ActionExpire  time.Duration
	OperateExpire time.Duration
	NewResExpire  time.Duration
}

// Memcache .
type Memcache struct {
	*memcache.Config
	TagExpire          time.Duration
	ResExpire          time.Duration
	ResAllTidsExpire   time.Duration
	ChannelGroupExpire time.Duration
}

// Tag .
type Tag struct {
	// sub number
	SubTagMaxNum int // 用户订阅tag最大数量限制
	SubArcMaxNum int // 用户订阅视频最大数量限制
	// arctag
	LikeLimitToLock int64 // tag点赞数量到达后锁定tag
	MaxResPageSize  int   // 每页最大限制
	MaxResLimit     int64 // tag下的最新资源数量限制
	MaxSelTagNum    int   // 接口批量查询tag信息的最大数量
	// arctag
	ResTagMaxNum int // 每个视频绑定tag的最大数量限制
	// user level
	ResTagAddLevel  int32 // 增加tag操作的用户登录最低限制
	ResTagDelLevel  int32 // 删除tag操作的用户等级最低限制
	ResTagRptLevel  int32 // 举报tag操作的用户等级最低限制
	ResTagLikeLevel int32 // 点赞tag操作的用户等级最低限制
	ResTagHateLevel int32 // 点踩tag操作的用户等级最低限制
	// operation numbn
	ResTagAddMaxNum  int   // 用户每日添加tag的最大数量
	ResTagDelMaxNum  int   // 用户每日删除tag的总数量限制
	ResTagDelSomeNum int   // 用户每日删除同一个视频下的tag数量限制
	ResTagLikeMaxNum int   // 用户每日点顶tag的最大数量
	ResTagHateMaxNum int   // 用户每日点踩tag的最大数量
	ResTagRptMaxNum  int   // 用户每日举报tag的最大数量
	ResOidLimit      int64 //用户批量查询资源tag的最大数量限制
}

// Supervision TODO 增加实名制
type Supervision struct {
	SixFour *struct {
		Button bool
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init .
func Init() (err error) {
	if confPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

func configCenter() (err error) {
	var (
		client *conf.Client
		value  string
		ok     bool
	)
	if client, err = conf.New(); err != nil {
		return
	}
	if value, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	_, err = toml.Decode(value, &Conf)
	return
}
