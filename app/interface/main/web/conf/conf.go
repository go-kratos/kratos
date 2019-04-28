package conf

import (
	"errors"
	"flag"
	xtime "time"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// log
	Log *log.Config
	// ecode
	Ecode *ecode.Config
	// http client
	HTTPClient *httpClient
	// HTTPServer
	HTTPServer *blademaster.ServerConfig
	// tracer
	Tracer *trace.Config
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// CouponRPC
	CouponRPC *rpc.ClientConfig
	// ArchiveRPC
	ArchiveRPC *rpc.ClientConfig
	// DynamicRPC
	DynamicRPC *rpc.ClientConfig
	// LocationRPC
	LocationRPC *rpc.ClientConfig
	// CoinRPC
	CoinRPC *rpc.ClientConfig
	// TagRPC
	TagRPC *rpc.ClientConfig
	// ArticleRPC
	ArticleRPC *rpc.ClientConfig
	// ResourceRPC
	ResourceRPC *rpc.ClientConfig
	// RelationRPC
	RelationRPC *rpc.ClientConfig
	// ThumbupRPC
	ThumbupRPC *rpc.ClientConfig
	// DM2RPC
	Dm2RPC *rpc.ClientConfig
	// FavRPC
	FavRPC *rpc.ClientConfig
	// Host
	Host host
	// redis
	Redis *redisConf
	// degrade
	DegradeConfig *degradeConfig
	// WEB
	WEB *web
	// Tag
	Tag *tag
	// DefaultTop
	DefaultTop *defaultTop
	// Bfs
	Bfs *bfs
	// Infoc2
	Infoc2 *infoc.Config
	// Rule
	Rule *rule
	// Warden Client
	BroadcastClient *warden.ClientConfig
	CoinClient      *warden.ClientConfig
	ArcClient       *warden.ClientConfig
	AccClient       *warden.ClientConfig
	ShareClient     *warden.ClientConfig
	UGCClient       *warden.ClientConfig
	// bnj
	Bnj2019 *bnj2019
}

type rule struct {
	// min cache rank count
	MinRankCount int
	// min cache rank index count
	MinRankIndexCount int
	// min cache rank region count
	MinRankRegionCount int
	// min cache rank recommend count
	MinRankRecCount int
	// min cache rank tag count
	MinRankTagCount int
	// min cache dynamic count
	MinDyCount int
	// min newlist tid arc count
	MinNewListCnt int
	// Elec
	ElecShowTypeIDs []int32
	// AuthorRecCnt author recommend count
	AuthorRecCnt int
	// RelatedArcCnt related archive limit count
	RelatedArcCnt int
	// MaxHelpPageSize help detail search max page  count
	MaxHelpPageSize int
	// newlist
	MaxArcsPageSize int
	// max size of second region newlist.
	MaxSecondCacheSize int
	// max size of first region newlist.
	MaxFirstCacheSize int
	// default num of dynamic archives
	DynamicNumArcs int
	// regions count
	RegionsCount int
	// bangumi count
	BangumiCount int
	// MaxArtPageSize max article page size
	MaxArtPageSize int
	// article up list get count
	ArtUpListGetCnt int
	// article up list count
	ArtUpListCnt int
	// min wechat hot count
	MinWxHotCount int
	// Rids first region ids
	Rids []int32
	// no related aids
	NoRelAids []int64
}

type httpClient struct {
	Read    *blademaster.ClientConfig
	Write   *blademaster.ClientConfig
	BigData *blademaster.ClientConfig
	Help    *blademaster.ClientConfig
	Search  *blademaster.ClientConfig
	Pay     *blademaster.ClientConfig
}

type host struct {
	Rank     string
	API      string
	Data     string
	Space    string
	Elec     string
	ArcAPI   string
	LiveAPI  string
	HelpAPI  string
	Mall     string
	Search   string
	Manager  string
	Pay      string
	AbServer string
}

type tag struct {
	PageSize int
	MaxSize  int
}

type redisConf struct {
	LocalRedis *localRedis
	BakRedis   *bakRedis
}

type localRedis struct {
	*redis.Config
	RankingExpire   time.Duration
	NewlistExpire   time.Duration
	RcExpire        time.Duration
	IndexIconExpire time.Duration
	WxHotExpire     time.Duration
}

type bakRedis struct {
	*redis.Config
	RankingExpire     time.Duration
	NewlistExpire     time.Duration
	RegionExpire      time.Duration
	ArchiveExpire     time.Duration
	TagExpire         time.Duration
	CardExpire        time.Duration
	RcExpire          time.Duration
	ArtUpExpire       time.Duration
	IndexIconExpire   time.Duration
	HelpExpire        time.Duration
	OlListExpire      time.Duration
	WxHotExpire       time.Duration
	AppealLimitExpire time.Duration
}

type web struct {
	PullRegionInterval    time.Duration
	PullOnlineInterval    time.Duration
	PullIndexIconInterval time.Duration
	SearchEggInterval     time.Duration
	OnlineCount           int
	SpecailInterval       time.Duration
}

type defaultTop struct {
	SImg string
	LImg string
}

type bfs struct {
	Addr        string
	Bucket      string
	Key         string
	Secret      string
	MaxFileSize int
	Timeout     time.Duration
}

type degradeConfig struct {
	Expire   int32
	Memcache *memcache.Config
}

type bnj2019 struct {
	Open        bool
	LiveAid     int64
	BnjMainAid  int64
	FakeElec    int64
	BnjListAids []int64
	BnjTick     time.Duration
	Timeline    []*struct {
		Name    string
		Start   xtime.Time
		End     xtime.Time
		Cover   string
		H5Cover string
	}
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
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
	client.Watch("web-interface.toml")
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
