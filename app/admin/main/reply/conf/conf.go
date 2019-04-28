package conf

import (
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	bauth "go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	// ConfPath config toml file path for reply admin project
	ConfPath string
	// Conf global config instance
	Conf *Config
)

// Config config.
type Config struct {
	// base
	// ecode
	Ecode *ecode.Config
	// log
	Log *log.Config
	// tracer
	Tracer *trace.Config
	// verify
	Verify *verify.Config
	// http
	HTTPServer *bm.ServerConfig
	// db
	DB *DB
	// memcache
	Memcache *Memcache
	// Redis
	Redis *Redis
	// http client
	HTTPClient        *bm.ClientConfig
	DrawyooHTTPClient *bm.ClientConfig
	// reply
	Reply *Reply
	// host
	Host *Host
	// Stats
	StatTypes     map[string]int32
	Weight        *Weight
	RPCClient2    *RPCClient2
	Databus       *Databus
	ManagerAuth   *bauth.Config
	ManagerReport *databus.Config
	Es            *elastic.Config
	ThumbupClient *warden.ClientConfig
	AccountClient *warden.ClientConfig
}

// Databus databus.
type Databus struct {
	Event *databus.Config
	Stats *databus.Config
}

// RPCClient2 rpc client.
type RPCClient2 struct {
	Account  *rpc.ClientConfig
	Archive  *rpc.ClientConfig
	Article  *rpc.ClientConfig
	Assist   *rpc.ClientConfig
	Thumbup  *rpc.ClientConfig
	Relation *rpc.ClientConfig
}

// Weight weight.
type Weight struct {
	Like int32
	Hate int32
}

// Stats stats.
type Stats struct {
	*databus.Config
	Type  int32
	Field string
}

// Redis redis.
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Reply reply.
type Reply struct {
	PageNum    int
	PageSize   int
	LikeWeight int32
	HateWeight int32
	// 大忽悠事件删除评论的管理员
	AdminName []string
	// 针对大忽悠的跳转链接
	Link string
	// 针对大忽悠时间的特殊稿件
	Oids []int64
	Tps  []int32
}

// Host host.
type Host struct {
	Search string
}

// DB db.
type DB struct {
	Reply      *sql.Config
	ReplySlave *sql.Config
}

// Memcache memcache.
type Memcache struct {
	*memcache.Config
	Expire time.Duration
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "config path")
}

// Init inti config.
func Init() (err error) {
	if ConfPath == "" {
		return configCenter()
	}
	_, err = toml.DecodeFile(ConfPath, &Conf)
	return
}

func configCenter() (err error) {
	var (
		ok     bool
		value  string
		client *conf.Client
	)
	if client, err = conf.New(); err != nil {
		return
	}
	if value, ok = client.Toml2(); !ok {
		panic(err)
	}
	_, err = toml.Decode(value, &Conf)
	return
}
