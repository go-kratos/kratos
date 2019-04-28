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
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// conf init.
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config struct
type Config struct {
	// Local cache
	LocalCache bool
	// interface XLog
	XLog *log.Config
	// tick time
	Tick xtime.Duration
	// tracer
	Tracer *trace.Config
	// databus
	UseractPub *databus.Config
	// httpClinet
	HTTPClient *bm.ClientConfig
	// httpIm9
	HTTPIm9 *bm.ClientConfig
	// httpSearch
	HTTPSearch *bm.ClientConfig
	// httpWrite
	HTTPWrite *bm.ClientConfig
	// httpLive
	HTTPLive *bm.ClientConfig
	// httpbangumi
	HTTPBangumi *bm.ClientConfig
	// httpbplus
	HTTPBPlus *bm.ClientConfig
	// httpgame
	HTTPGame *bm.ClientConfig
	// http
	BM *HTTPServers
	// host
	Host *Host
	// rpc client
	AccountRPC  *rpc.ClientConfig
	ArchiveRPC  *rpc.ClientConfig
	TagRPC      *rpc.ClientConfig
	ArticleRPC  *rpc.ClientConfig
	RelationRPC *rpc.ClientConfig
	ThumbupRPC  *rpc.ClientConfig
	HistoryRPC  *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	MemberRPC   *rpc.ClientConfig
	LocationRPC *rpc.ClientConfig
	// db
	MySQL *MySQL
	// ecode
	// ecode
	Ecode *ecode.Config
	// redis
	Redis *Redis
	// mc
	Memcache *Memcache
	// search
	Search *Search
	// space
	Space *Space
	// contribute
	ContributePub *databus.Config
	CoinClient    *warden.ClientConfig
	// BroadcastRPC grpc
	PGCRPC *warden.ClientConfig
	// build limit
	SearchBuildLimit *SearchBuildLimit
	// login build
	LoginBuild *LoginBuild
	// infoc
	Infoc *infoc.Config
	// fav Client
	FavClient *warden.ClientConfig
	// search dynamic
	SearchDynamicSwitch *SearchDynamicSwitch
}

// LoginBuild is
type LoginBuild struct {
	Iphone int
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
}

// Host struct
type Host struct {
	Account   string
	Bangumi   string
	APICo     string
	Im9       string
	Search    string
	Game      string
	Space     string
	Elec      string
	VC        string
	APILiveCo string
	WWW       string
	Show      string
	Pay       string
	Member    string
	Mall      string
}

// MySQL struct
type MySQL struct {
	Show *sql.Config
}

// Redis struct
type Redis struct {
	Contribute *struct {
		*redis.Config
		Expire xtime.Duration
	}
}

// Memcache struct
type Memcache struct {
	Archive *struct {
		*memcache.Config
		RecommedExpire xtime.Duration
		ArchiveExpire  xtime.Duration
	}
}

// Search struct
type Search struct {
	SeasonNum          int
	MovieNum           int
	SeasonMore         int
	MovieMore          int
	UpUserNum          int
	UVLimit            int
	UserNum            int
	UserVideoLimit     int
	BiliUserNum        int
	BiliUserVideoLimit int
	OperationNum       int
	IPadSearchBangumi  int
	IPadSearchFt       int
}

// Space struct
type Space struct {
	ForbidMid []int64
}

// SearchBuildLimit struct
type SearchBuildLimit struct {
	PGCHighLightIOS       int
	PGCHighLightAndroid   int
	PGCALLIOS             int
	PGCALLAndroid         int
	SpecialerGuideIOS     int
	SpecialerGuideAndroid int
	SearchArticleIOS      int
	SearchArticleAndroid  int
	ComicIOS              int
	ComicAndroid          int
	ChannelIOS            int
	ChannelAndroid        int
	CooperationIOS        int
	CooperationAndroid    int
	QueryCorIOS           int
	QueryCorAndroid       int
	SugDetailIOS          int
	SugDetailAndroid      int
	NewTwitterIOS         int
	NewTwitterAndroid     int
}

type SearchDynamicSwitch struct {
	IsUP    bool
	IsCount bool
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init init conf
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
	client.Watch("app-interface.toml")
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
