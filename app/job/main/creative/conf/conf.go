package conf

import (
	"errors"
	"flag"

	"github.com/BurntSushi/toml"

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
	"go-common/library/time"

	"go-common/library/database/hbase.v2"
)

// Conf info.
var (
	ConfPath string
	Conf     = &Config{}
	client   *conf.Client
)

// Config struct.
type Config struct {
	// db
	DB *DB
	// base
	// ecode
	Ecode *ecode.Config
	// log
	Log *log.Config
	// httpClinet
	HTTPClient *HTTPClient
	// BM
	BM *HTTPServers
	// tracer
	Tracer *trace.Config
	// host
	Host *Host
	// Databus
	Env          string
	Consume      bool
	Pub          bool
	ArcSub       *databus.Config
	ArcNotifySub *databus.Config
	UpPub        *databus.Config
	//task
	TaskSub, ShareSub, RelationSub, StatLikeSub *databus.Config
	StatShareSub, StatCoinSub, StatFavSub       *databus.Config
	StatReplySub, StatDMSub, StatViewSub        *databus.Config
	NewUpSub                                    *databus.Config

	Task *Task
	// channal len
	ChanSize int64
	//moni
	Monitor *Monitor
	// rpc client2
	ArchiveRPC *rpc.ClientConfig
	ArticleRPC *rpc.ClientConfig
	// grpc Client
	CreativeGRPClient *warden.ClientConfig
	UpGRPCClient      *warden.ClientConfig
	//hot compute switch
	HotSwitch bool
	//hot compute switch
	HonorSwitch    bool
	HonorStep      int
	HonorMSGSpec   string
	HonorFlushSpec string
	SendEveryWeek  bool
	// hbase
	HBaseOld *HBaseConfig
	// infoc
	WeeklyHonorInfoc *infoc.Config
}

// DB conf.
type DB struct {
	// Creative db
	Creative *sql.Config
	// Archive db
	Archive *sql.Config
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
}

// HTTPClient conf.
type HTTPClient struct {
	Normal *bm.ClientConfig
	Slow   *bm.ClientConfig
}

// Host conf.
type Host struct {
	Monitor  string
	Passport string
	Account  string
	Message  string
	API      string
	Videoup  string
}

// Monitor conf.
type Monitor struct {
	Host      string
	Moni      string
	UserName  string
	AppSecret string
	AppToken  string
}

// Task conf.
type Task struct {
	//扫表
	RowLimit        int //每次从表中取的最大数据量
	TableJobNum     int //开启n张表的扫描任务的协程数量
	TableConsumeNum int //开启消费表数据的协程数量
	//databus 消费
	SwitchHighQPS    bool
	SwitchDatabus    bool
	DatabusQueueLen  int
	StatViewQueueLen int
	StatLikeQueueLen int
	ChanSize         int
	//task notify
	SwitchMsgNotify     bool
	TaskRowLimitNum     int
	TaskTableJobNum     int
	TaskTableConsumeNum int
	TaskExpireTime      int64
	TaskSendHour        int
	TaskSendMiniute     int
	TaskSendSecond      int
	TaskBatchMidNum     int

	TaskMsgCode    string
	TaskTitle      string
	TaskContent    string
	TestNotifyMids string
	//reward notify
	RewardRowLimitNum     int
	RewardTableJobNum     int
	RewardTableConsumeNum int
	RewardWeek            int
	RewardLastDay         int
	RewardLastHour        int
	RewardLastMiniute     int
	RewardLastSecond      int
	RewardNowHour         int
	RewardNowMiniute      int
	RewardNowSecond       int
	RewardSendHour        int
	RewardSendMiniute     int
	RewardSendSecond      int
	RewardBatchMidNum     int

	RewardMsgCode string
	RewardTitle   string
	RewardContent string
	BiliMID       int64

	//新手和进阶粉丝数
	NewFollower      int64
	AdvancedFollower int64

	//单个稿件计数
	StatView  int64
	StatLike  int64
	StatReply int64
	StatShare int64
	StatFav   int64
	StatCoin  int64
	StatDM    int64
	//单个稿件计数设定上限
	StatViewUp  int64
	StatLikeUp  int64
	StatReplyUp int64
	StatShareUp int64
	StatFavUp   int64
	StatCoinUp  int64
	StatDMUp    int64
}

func init() {
	flag.StringVar(&ConfPath, "conf", "", "default config path")
}

// Init conf.
func Init() (err error) {
	if ConfPath != "" {
		return local()
	}
	return remote()
}

func local() (err error) {
	_, err = toml.DecodeFile(ConfPath, &Conf)
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
		tomlStr string
		ok      bool
		tmpConf *Config
	)
	if tomlStr, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(tomlStr, &tmpConf); err != nil {
		return errors.New("could not decode toml config")
	}
	*Conf = *tmpConf
	return
}
