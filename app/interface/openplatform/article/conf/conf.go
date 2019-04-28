package conf

import (
	"errors"
	"flag"

	"go-common/app/interface/openplatform/article/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
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
	hbase "go-common/library/database/hbase.v2"
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
	// base
	// elk
	Log *log.Config
	// BM
	BM *bm.ServerConfig
	// HTTPClient .
	HTTPClient        *bm.ClientConfig
	MessageHTTPClient *bm.ClientConfig
	// tracer
	Tracer *trace.Config
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// redis
	Redis *redis.Config
	// memcache
	Memcache *Memcache
	// MySQL
	MySQL MySQL
	// rpc
	AccountRPC *rpc.ClientConfig
	TagRPC     *rpc.ClientConfig
	FavRPC     *rpc.ClientConfig
	ArcRPC     *rpc.ClientConfig
	CoinRPC    *rpc.ClientConfig
	ResRPC     *rpc.ClientConfig
	ThumbupRPC *rpc.ClientConfig
	FilterRPC  *rpc.ClientConfig
	HistoryRPC *rpc.ClientConfig
	SearchRPC  *warden.ClientConfig
	// databus
	StatDatabus *databus.Config
	// infoc log2
	DisplayInfoc *infoc.Config
	ClickInfoc   *infoc.Config
	AIClickInfoc *infoc.Config
	ShowInfoc    *infoc.Config
	CheatInfoc   *infoc.Config
	// ecode
	Ecode *ecode.Config
	//RankCategories .
	RankCategories []*model.RankCategory
	//Message
	Message Message
	Cards   Cards
	//hbase
	HBase *hbase.Config
	// BFS
	BFS *BFS
	// Antispam
	Antispam *antispam.Config
	// DegradeConfig
	DegradeConfig *DegradeConfig
	// artcile
	Article Article
	//Berserker
	Berserker Berserker
	//Sentinel
	Sentinel *Sentinel
}

//Sentinel .
type Sentinel struct {
	EnableSentinel     int `json:"enableSentinel"`
	DurationSample     int `json:"durationSample"`
	MonitorCountSample int `json:"monitorCountSample"`
	MonitorRateSample  int `json:"monitorRateSample"`
	DebugSample        int `json:"debugSample"`
}

// Berserker .
type Berserker struct {
	AppKey    string
	AppSecret string
	URL       string
}

// Memcache config
type Memcache struct {
	*memcache.Config
	ArticleExpire     time.Duration
	StatsExpire       time.Duration
	LikeExpire        time.Duration
	CardsExpire       time.Duration
	SubmitExpire      time.Duration
	ListArtsExpire    time.Duration
	ListExpire        time.Duration
	ArtListExpire     time.Duration
	UpListsExpire     time.Duration
	ListReadExpire    time.Duration
	HotspotExpire     time.Duration
	AuthorExpire      time.Duration
	ArticlesIDExpire  time.Duration
	ArticleTagExpire  time.Duration
	UpStatDailyExpire time.Duration
}

// MySQL config
type MySQL struct {
	Article *sql.Config
}

// Cards config
type Cards struct {
	TicketURL  string
	MallURL    string
	AudioURL   string
	BangumiURL string
}

// BFS bfs config
type BFS struct {
	Timeout     time.Duration
	MaxFileSize int
	Bucket      string
	URL         string
	Method      string
	Key         string
	Secret      string
}

// DegradeConfig .
type DegradeConfig struct {
	Expire   int32
	Memcache *memcache.Config
}

// Article article config
type Article struct {
	ExpireUpper                time.Duration
	ExpireArtLikes             time.Duration
	ExpireSortArts             time.Duration
	TTLSortArts                time.Duration
	ExpireRank                 time.Duration
	TTLRank                    time.Duration
	ExpireMaxLike              time.Duration
	ExpireHotspot              time.Duration
	CreationDefaultSize        int
	CreationMaxSize            int
	UpperDraftLimit            int
	UpperArticleLimit          int
	UpdateRecommendsInteval    time.Duration
	MaxRecommendPnSize         int64
	MaxRecommendPsSize         int64
	MaxUpperListPsSize         int64
	MaxArchives                int
	MaxComplaintReasonLimit    int64
	MaxArticleMetas            int
	MaxApplyContentLimit       int64
	MaxApplyCategoryLimit      int64
	MaxLikeMidLen              int
	RecommendAidLen            int
	SortLimitTime              time.Duration
	UpdateBannersInteval       time.Duration
	BannerIDs                  []int
	ActBannerIDs               []int
	RecommendRegionLen         int
	SkyHorseRecommendRegionLen int
	RankHost                   string
	MessageMids                []int64
	MaxContentSize             int
	MaxContentLength           int
	MinContentLength           int
	ActAddURI                  string
	ActDelURI                  string
	ActURI                     string
	ListLimit                  int
	ListArtsLimit              int
	AppCategoryName            string
	AppCategoryURL             string
	SkyHorseURL                string
	SkyHorseGray               []int64
	SkyHorseGrayUsers          []int64
	ListDefaultImage           string
	RecommendAuthors           int
	ExpireReadPing             time.Duration
	ExpireReadSet              time.Duration
	Media                      []int64
	EditTimes                  int
	RecommendAuthorsURL        string
}

// Message .
type Message struct {
	URL string
	MC  string
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
	err = load()
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
