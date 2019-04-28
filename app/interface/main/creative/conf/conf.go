package conf

import (
	"encoding/json"
	"errors"
	"flag"
	"strconv"

	"go-common/app/interface/main/creative/model/academy"
	appMdl "go-common/app/interface/main/creative/model/app"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

	"go-common/library/net/rpc/warden"

	"github.com/BurntSushi/toml"
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
	// app
	App *bm.App
	// HTTPClient
	HTTPClient *HTTPClient
	// BM
	BM *HTTPServers
	// rpc client2
	ArchiveRPC  *rpc.ClientConfig
	ArticleRPC  *rpc.ClientConfig
	ResourceRPC *rpc.ClientConfig
	RelationRPC *rpc.ClientConfig
	UPRPC       *rpc.ClientConfig
	SubRPC      *rpc.ClientConfig
	// mc
	Memcache *Memcache

	// redis
	Redis *Redis
	// tracer
	Tracer *trace.Config
	//nhbase
	HBase    *HBaseConfig
	HBaseOld *HBaseConfig
	// white list
	WhiteAccessKey string
	WhiteMid       int64
	// host
	Host   *Host
	H5Page *H5Page
	// geetest
	Geetest *Geetest
	// whitelist
	Whitelist *Whitelist
	// ArchStatus
	ArchStatus     map[string]string
	RouterAntispam *antispam.Config
	DmAntispam     *antispam.Config
	// BFS
	BFS       *BFS
	AppealTag int64
	// databus sub
	UserInfoSub *databus.Config
	TaskPub     *databus.Config
	// WaterMark
	WaterMark   *WaterMark
	Game        *Game
	Growup      *Growup
	StatCacheOn bool
	AppIcon     *AppIcon
	UgcPay      *UgcPay
	//academy
	Coefficient *Coefficient
	// academy
	Academy      *Academy
	AcaRecommend *AcaRecommend
	//ManagerReport 行为日志平台
	ManagerReport *databus.Config
	// rpc server
	WardenServer *warden.ServerConfig
	WardenClient *warden.ClientConfig
	CoinClient   *warden.ClientConfig
	AccClient    *warden.ClientConfig
	UpClient     *warden.ClientConfig
	// task condition
	TaskCondition *TaskCondition

	//联合投稿配置
	StaffConf *StaffConf
	// honor weekly degrade switch
	HonorDegradeSwitch bool
}

// TaskCondition task condition
type TaskCondition struct {
	Fans              int64
	ReceiveMsg        string
	ReceiveMsgPendant string
	WhiteSwitch       bool
	AppIndexSwitch    bool
}

// StaffConf 联合投稿配置
type StaffConf struct {
	IsGray   bool             `json:"is_gray"`
	TypeList []*StaffTypeConf `json:"typelist"`
}

// StaffTypeConf 联合投稿的分区配置
type StaffTypeConf struct {
	TypeID   int16 `json:"typeid"`
	MaxStaff int   `json:"max_staff"`
}

// UgcPay str
type UgcPay struct {
	ProtocolID      string
	AllowDeleteDays int
	AllowEditDays   int
}

// Coefficient str
type Coefficient struct {
	ActHeat float64
}

// AppIcon str
type AppIcon struct {
	CameraInput *appMdl.Icon `json:"camera_input"`
	CameraCoo   *appMdl.Icon `json:"camera_coo"`
}

// Game str Conf
type Game struct {
	OpenHost string
	App      *bm.App
}

// Growup str
type Growup struct {
	LimitFanCnt     int64 // LimitFanCnt 一万粉
	LimitTotalClick int64 // LimitTotalClick 五十万点击量
}

// DB conf.
type DB struct {
	// archive db
	Creative *sql.Config
	Archive  *sql.Config
}

// Thrift conf.
type Thrift struct {
	Addr                     string
	Idle                     int
	DialTimeout, ReadTimeout time.Duration
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
	Local *bm.ServerConfig
}

// HTTPClient conf.
type HTTPClient struct {
	Normal   *bm.ClientConfig
	Slow     *bm.ClientConfig
	UpMng    *bm.ClientConfig
	Fast     *bm.ClientConfig
	Chaodian *bm.ClientConfig
}

// Memcache conf.
type Memcache struct {
	Data struct {
		*memcache.Config
		DataExpire  time.Duration
		IndexExpire time.Duration
	}
	Archive struct {
		*memcache.Config
		TplExpire time.Duration
	}
	Honor struct {
		*memcache.Config
		HonorExpire time.Duration
		ClickExpire time.Duration
	}
}

// Redis conf.
type Redis struct {
	Antispam *struct {
		*redis.Config
		Expire time.Duration
	}
	Cover *struct {
		*redis.Config
		Expire time.Duration
	}
}

// Host conf.
type Host struct {
	Passport   string
	Archive    string
	Search     string
	API        string
	Data       string
	Member     string
	Act        string
	Activity   string
	Videoup    string
	Tag        string
	Geetest    string
	Account    string
	UpMng      string
	Elec       string
	Live       string
	Monitor    string
	Coverrec   string
	Growup     string
	Matsuri    string
	ArcTip     string
	Message    string
	HelpAPI    string
	MainSearch string
	Dynamic    string
	Mall       string //会员购
	BPay       string //B币券
	Pendant    string //挂件
	BigMember  string //大会员
	Profit     string //激励计划
	Notify     string //消息通知
	Chaodian   string //超电
}

// H5Page conf.
type H5Page struct {
	FAQVideoEditor  string
	CreativeCollege string
	HotAct          string
	Draft           string
	Passport        string
	Mission         string
	Cooperate       string
}

// Geetest geetest id & key
type Geetest struct {
	CaptchaID   string
	MCaptchaID  string
	PrivateKEY  string
	MPrivateKEY string
}

// Whitelist str
type Whitelist struct {
	DataMids          []int64
	ArcMids           []int64
	ForbidVideoupMids []int64
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

// WaterMark config
type WaterMark struct {
	UnameMark string
	UIDMark   string
	SaveImg   string
	FontFile  string
	FontSize  int
	Consume   bool
}

// HBaseConfig for new hbase client.
type HBaseConfig struct {
	*hbase.Config
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}

//Academy for academy h5 conf
type Academy struct {
	academy.H5Conf
}

//AcaRecommend for h5 rec conf
type AcaRecommend struct {
	academy.Recommend
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

	//bs, err := ioutil.ReadFile("academy.json")
	//if err != nil {
	//	return err
	//}
	//if err = json.Unmarshal([]byte(bs), &Conf.AcaRecommend); err != nil {
	//	return errors.New("could not decode json config")
	//}
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
		tomlStr         string
		jsonStr, acaStr string
		ok              bool
		tmpConf         *Config
		archStatus      map[string]string
	)
	if tomlStr, ok = client.Toml2(); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(tomlStr, &tmpConf); err != nil {
		return errors.New("could not decode toml config")
	}
	if jsonStr, ok = client.Value("archStatus.json"); !ok {
		return errors.New("load config center error")
	}
	if err = json.Unmarshal([]byte(jsonStr), &archStatus); err != nil {
		return errors.New("could not decode json config")
	}

	if acaStr, ok = client.Value("academy.json"); !ok {
		return errors.New("load config center error")
	}
	if err = json.Unmarshal([]byte(acaStr), &tmpConf.AcaRecommend); err != nil {
		return errors.New("could not decode json config")
	}
	tmpConf.ArchStatus = archStatus
	*Conf = *tmpConf
	return
}

// StatDesc define
func (c *Config) StatDesc(state int) (desc string) {
	statusStr := strconv.Itoa(state)
	if v, ok := c.ArchStatus[statusStr]; !ok {
		desc = statusStr
	} else {
		desc = v
	}
	return
}
