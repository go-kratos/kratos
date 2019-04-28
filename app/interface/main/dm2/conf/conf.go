package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/bfs"
	"go-common/library/database/elastic"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
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

var (
	confPath string
	client   *conf.Client
	// Conf config export var
	Conf = &Config{}
)

// Config dm config struct
type Config struct {
	// ecode
	Ecode *ecode.Config
	// log
	Xlog   *log.Config
	Infoc2 *infoc.Config
	// tracer
	Tracer *trace.Config
	Auth   *auth.Config
	Verify *verify.Config
	// database
	DB *DB
	// redis
	Redis *Redis
	// memcache
	Memcache *Memcache
	// localcache
	Localcache *Localcache
	// host
	Host *Host
	// http
	HTTPServer *bm.ServerConfig
	// http client
	HTTPCli *bm.ClientConfig
	// rpc client
	ArchiveRPC *rpc.ClientConfig
	// account rpc client
	AccountRPC *warden.ClientConfig
	// assist rpc client
	AssistRPC *rpc.ClientConfig
	//coin rpc client
	CoinRPC *rpc.ClientConfig
	// thumbup rpc client
	ThumbupRPC *warden.ClientConfig
	// relation rpc client
	RelationRPC *rpc.ClientConfig
	// seq-server rpc client
	SeqRPC *rpc.ClientConfig
	Seq    *Seq
	// Location rpc client
	LocationRPC *rpc.ClientConfig
	// MemberRpc Rpc
	MemberRPC *rpc.ClientConfig
	// Filter Rpc
	FilterRPC *warden.ClientConfig
	// Figure RPC
	FigureRPC *rpc.ClientConfig
	// UgcPayRPC RPC
	UgcPayRPC *warden.ClientConfig
	// SpyRPC RPC
	SpyRPC *rpc.ClientConfig
	// Season RPC
	SeasonRPC *warden.ClientConfig
	// databus
	Databus   *databus.Config
	ActionPub *databus.Config
	// dm list realname
	Realname *Realname
	// dm rpc server
	RPCServer *rpc.ServerConfig
	// Antispam
	Antispam *antispam.Config
	// supervision conf
	Supervision *Supervision
	// Elastic config
	Elastic *elastic.Config
	// Bfs
	Bfs *Bfs
	// Subtitle Check Databus
	SubtitleCheckPub *databus.Config
	// Garbage danmu Switch
	Switch *Switch
	// BroadcastLimit
	BroadcastLimit *BroadcastLimit
	// DmFlag
	DmFlag *DmFlag
}

// DmFlag .
type DmFlag struct {
	RecFlag   int
	RecText   string
	RecSwitch int
}

// BroadcastLimit broadcast limit
type BroadcastLimit struct {
	Limit    int
	Interval int
}

// Switch .
type Switch struct {
	GarbageDanmu bool
}

// Bfs .
type Bfs struct {
	Client         *bfs.Config
	BucketSubtitle string
}

//Seq Conf
type Seq struct {
	DM *struct {
		BusinessID int64
		Token      string
	}
	Subtitle *struct {
		BusinessID int64
		Token      string
	}
}

// DB mysql
type DB struct {
	DMReader *sql.Config
	DMWriter *sql.Config
	DM       *sql.Config
}

// Redis dm redis
type Redis struct {
	DM *struct {
		*redis.Config
		Expire time.Duration
	}
	DMRct *struct {
		*redis.Config
		Expire time.Duration
	}
	DMSeg *struct {
		*redis.Config
		Expire time.Duration
	}
}

// Memcache dm memcache
type Memcache struct {
	DM *struct {
		*memcache.Config
		DMExpire      time.Duration
		SubjectExpire time.Duration
		HistoryExpire time.Duration
		AjaxExpire    time.Duration
		DMMaskExpire  time.Duration
	}
	Filter *struct {
		*memcache.Config
		Expire time.Duration
	}
	Subtitle *struct {
		*memcache.Config
		Expire time.Duration
	}
	DMSeg *struct {
		*memcache.Config
		DMExpire        time.Duration
		DMLimiterExpire time.Duration
	}
}

// Localcache cache stored in local
type Localcache struct {
	Oids       []int64
	Expire     time.Duration
	ViewAids   []int64
	ViewExpire time.Duration
}

// Realname realname switch and config
type Realname struct {
	SwitchOn  bool
	Whitelist []int64
	Threshold map[string]int64
}

// Host http host
type Host struct {
	AI        string
	API       string
	Archive   string
	Message   string
	Search    string
	MaskCloud string
	Advert    string
	Upos      string
	Self      string
}

// Supervision supervision .
type Supervision struct {
	Completed bool
	StartTime string
	EndTime   string
	Location  string
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

//Init int config
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
