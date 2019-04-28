package conf

import (
	"errors"
	"flag"
	xtime "time"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf conf
	Conf   = &Config{}
	client *conf.Client
)

// Config so config
type Config struct {
	// interface Log
	Log *log.Config
	//HTTPClient
	HTTPClient *bm.ClientConfig
	// BM
	BM *bm.ServerConfig
	// rpc
	ArchiveRPC *rpc.ClientConfig
	ArticleRPC *rpc.ClientConfig
	CoinRPC    *rpc.ClientConfig
	ActRPC     *rpc.ClientConfig
	// grpc
	AccClient *warden.ClientConfig
	// DB
	MySQL *MySQL
	// mc
	Memcache *Memcache
	// redis
	Redis *Redis
	// databus
	ActSub *databus.Config
	BnjSub *databus.Config
	// vip binlog databus
	//VipSub *databus.Config
	KfcSub *databus.Config
	// Interval
	Interval *interval
	// Rule
	Rule *rule
	// Host
	Host *host
	// Elastic
	Elastic *elastic.Config
	// bnj
	Bnj2019 *bnj2019
}

type bnj2019 struct {
	GameCancel    int
	LID           int64
	StartTime     xtime.Time
	TimelinePic   string
	H5TimelinePic string
	MsgSpec       string
	MidLimit      int64
	WxKey         string
	WxTitle       string
	WxUser        string
	Time          []*struct {
		Score    int64
		Second   int64
		Step     int
		WxMsg    string
		MsgTitle string
		MsgMc    string
		Msg      string
	}
	Message []*struct {
		Start   xtime.Time
		Title   string
		Content string
		Mc      string
	}
}

type interval struct {
	CoinInterval      time.Duration
	QueryInterval     time.Duration
	ObjStatInterval   time.Duration
	ViewRankInterval  time.Duration
	KingStoryInterval time.Duration
}

// MySQL is db config.
type MySQL struct {
	Like *sql.Config
}

// Redis config
type Redis struct {
	*redis.Config
	Expire time.Duration
}

// Memcache config
type Memcache struct {
	Like             *memcache.Config
	LikeExpire       time.Duration
	TimeFinishExpire time.Duration
	LessTimeExpire   time.Duration
}

type rule struct {
	BroadcastCid  int64
	BroadcastSid  int64
	ArcObjStatSid int64
	ArtObjStatSid int64
	KingStorySid  int64
	EleLotteryID  int64
}

type host struct {
	APICo    string
	Activity string
	MsgCo    string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
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
	client.Watch("activity-job.toml")
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
		tmpConf = &Config{}
	)
	if s, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
