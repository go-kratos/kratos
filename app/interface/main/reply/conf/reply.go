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
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

// global var
var (
	ConfPath string
	Conf     *Config
)

// Config represent service conf
type Config struct {
	BM *bm.ServerConfig
	//reply
	Reply *Reply
	// HTTPClinet
	HTTPClient        *bm.ClientConfig
	DrawyooHTTPClient *bm.ClientConfig
	FilterGRPCClient  *warden.ClientConfig
	FeedGRPCClient    *warden.ClientConfig
	AccountGRPCClient *warden.ClientConfig
	// rpc
	RPCClient2 *RPCClient2
	// mysql
	MySQL *MySQL
	// redis
	Redis *Redis
	// mc
	Memcache *Memcache
	// seq conf
	Seq *Seq
	// kafka
	Databus *databus.Config
	// tracer
	Tracer *trace.Config
	// XLog
	XLog *log.Config
	// auth
	Auth *auth.Config
	// verify
	Verify *verify.Config
	// ecode
	Ecode *ecode.Config

	Host *Host
	// appkey type
	AppkeyType map[string][]int8
	// supervision conf
	Supervision    *Supervision
	AssistConfig   *AssistConfig
	Identification *Identification
	ReportAgent    *log.AgentConfig
	UserReport     *databus.Config
	// es config
	Es *elastic.Config
	// info config
	Infoc *infoc.Config
}

//Seq Conf
type Seq struct {
	BusinessID int64
	Token      string
}

// Reply represents reply conf
type Reply struct {
	HotReply         int
	MaxPageSize      int
	MinConLen        int
	MaxConLen        int
	SecondDefSize    int
	SecondDefPageNum int
	EmojiExpire      time.Duration
	MaxEmoji         int
	BigdataFilter    bool
	// url
	BigdataURL          string
	AiTopicURL          string
	VipURL              string
	FansReceivedListURL string
	BlockStatusURL      string
	CaptchaTokenURL     string
	CaptchaVerifyURL    string
	CreditUserURL       string
	ReplyLogSearchURL   string

	AidWhiteList []int64
	ForbidList   []int64
	BnjAidList   []int64

	// 默认排序开关
	SortByHotOids  map[string]int8
	SortByTimeOids map[string]int8
	HideFloorOids  map[string]int8

	// 拜年祭的一些视频默认热评数目需要调整到N个
	HotReplyConfig map[string]map[string]int
}

// Host host.
type Host struct {
	API    string
	Search string
}

// MySQL represent mysql conf
type MySQL struct {
	Reply      *sql.Config
	ReplySlave *sql.Config
}

// Redis represent redis conf
type Redis struct {
	*redis.Config
	IndexExpire   time.Duration
	ReportExpire  time.Duration
	UserCntExpire time.Duration
	UserActExpire time.Duration
}

// Memcache represent mc conf
type Memcache struct {
	*memcache.Config
	Expire      time.Duration
	EmptyExpire time.Duration
}

// RPCClient2 represent rpc conf
type RPCClient2 struct {
	Account  *rpc.ClientConfig
	Filter   *rpc.ClientConfig
	Location *rpc.ClientConfig
	Assist   *rpc.ClientConfig
	Figure   *rpc.ClientConfig
	Seq      *rpc.ClientConfig
	Thumbup  *rpc.ClientConfig
	Archive  *rpc.ClientConfig
	Article  *rpc.ClientConfig
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "config path")
}

// Supervision supervision .
type Supervision struct {
	StartTime   string
	EndTime     string
	Completed   bool
	Location    string
	ReportAgent *log.AgentConfig
}

// AssistConfig Assist configurations .
type AssistConfig struct {
	StartTime string
}

// Identification identification configurations.
type Identification struct {
	SwitchOn bool
}

// Init init conf
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
