package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"go-common/library/database/hbase.v2"

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
	Log        *log.Config
	HTTPClient *blademaster.ClientConfig
	Tracer     *trace.Config
	Auth       *auth.Config
	Verify     *verify.Config
	// Wechat wechat config
	Wechat *wechat
	// HBase for fans
	HBase *hbaseConf
	// FansHBase for attention groups + active time
	FansHBase *hbaseConf
	Redis     *redis.Config
	// MySQL
	MySQL      *sql.Config
	AccountRPC *rpc.ClientConfig
	// ArchiveSub archive_result databus consumer
	ArchiveSub *databus.Config
	// RelationSub relation_xxx databus consumer
	RelationSub *databus.Config
	Push        *push
	// ArcPush archive push settings
	ArcPush *arcPush
	PushRPC *warden.ClientConfig
	// Anti antispam
	Anti   *antispam.Config
	Bm     *blademaster.ServerConfig
	Abtest *abtest
}

type abtest struct {
	HbaseBlacklistTable  string
	HbaseBlacklistFamily []string
	HbaseeWhitelistTable string
	HbaseWhitelistFamily []string
	TestGroup            []int
	ComparisonGroup      []int
	TestMids             []int64
}

type hbaseConf struct {
	hbase.Config
	ReadTimeout   xtime.Duration
	ReadsTimeout  xtime.Duration
	WriteTimeout  xtime.Duration
	WritesTimeout xtime.Duration
}

/**
* 配置规则：
PushStatisticsKeepDays 推送数据保留天数
PushStatisticsClearTim 每日推送数据删除的时间点
Order 分组优先级，元素=类型#组名，优先级只针对同一类型下有效,没配置优先级的分组不可用
ActiveTime 默认活跃时间(24小时制),过滤粉丝是否在活跃时间段内；未配置则不过滤；若希望过滤活跃时间但不提供默认活跃时间，配置成[0]
ForbidTimes 固定免推送时间段组
Proportions 灰度策略，粉丝尾号100内, 起始点+step
FanGroup 分组具体信息
*/
// arcPush 稿件更新的推送设置
type arcPush struct {
	PushStatisticsKeepDays  int
	PushStatisticsClearTime string
	Order                   []string
	ActiveTime              []int
	ForbidTimes             []ForbidTime
	Proportions             []Proportion
	FanGroup                []*fanGroup
	UpperLimitExpire        xtime.Duration
}

// ForbidTime 禁止时间范围
type ForbidTime struct {
	PushForbidStartTime string
	PushForbidEndTime   string
}

// Proportion 灰度uid范围
type Proportion struct {
	ProportionStartFrom string
	Proportion          string //必须是2位小数，比如:1.00, 0.05
}

// fanGroup 关注up主的粉丝分组
type fanGroup struct {
	Name            string
	Desc            string
	RelationType    int
	Hitby           string //命中分组规则,default=全部命中，hbase=hbase表过滤
	Limit           int
	PerUpperLimit   int
	LimitExpire     xtime.Duration
	HBaseTable      string
	HBaseFamily     []string
	MsgTemplateDesc string
	MsgTemplate     string
}

type wechat struct {
	UserName, Token, Secret string
}

type push struct {
	ProdSwitch           bool
	AddAPI               string
	MultiAPI             string
	BusinessID           int
	BusinessToken        string
	BusinessSpecialID    int
	BusinessSpecialToken string
	LoadSettingsInterval xtime.Duration
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
