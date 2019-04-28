package conf

import (
	"errors"
	"flag"

	"go-common/app/admin/ep/melloi/model"
	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/permit"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log *log.Config
	//BM *HTTPServers
	Tracer         *trace.Config
	Redis          *redis.Config
	Memcache       *memcache.Config
	Ecode          *ecode.Config
	ORM            *orm.Config
	Permit         *permit.Config2
	PermitGRPC     *warden.ClientConfig
	HTTPClient     *bm.ClientConfig
	ServiceTree    *model.TreeConf
	Dapper         *Dapper
	BfsConf        *model.BFSConf
	ServiceCluster *model.ClusterConf
	Melloi         *Melloi
	BM             *bm.ServerConfig
	Wechat         *Wechat
	Mail           *Mail
	Paas           *Paas
	Grpc           *Grpc
	Jmeter         *Jmeter
	DockerStatus   *DockerStatus
}

// Dapper conf
type Dapper struct {
	Host string
}

// DockerStatus conf
type DockerStatus struct {
	Host string // ip
	Port int    // port
}

//Paas conf
type Paas struct {
	APIToken       string // PaaS token获取
	PlatformID     string //PaaS token获取
	BusinessUnit   string //BU
	Project        string //项目
	App            string //应用
	Env            string //环境
	Image          string //镜像名称
	ImageVersion   string //镜像版本
	Volumes        string //PaaS 创建job
	ResourcePoolID string //PaaS 创建job
	Completions    int    //PaaS 创建job
	RetriesLimit   int    //PaaS 创建job
	NetworkID      int
	ClusterID      int
	TreeID         int    //服务树ID
	HostInfo       string //PaaS 创建job
	Action         string //paas 查询容器cpu
	PublicKey      string //key
	Signature      int    //paas 查询容器cpu
	DataSource     string //数据源
	Query          string //paas 查询容器cpu语句
	CPUCore        int    //cpu 核数
	CPUCoreDebug   int    //debug cpu 核数
}

//Melloi melloi config
type Melloi struct {
	AppkeyProd          string   //线上appkey
	SecretProd          string   //线上sceret
	AppkeyUat           string   //uat appkey
	SecretUat           string   //uat secret
	Executor            []string //白名单
	CheckTime           bool     //是否校验压测时间
	MaxFileSize         int64
	DefaultHost         string
	MaxDowloadSize      int64
	DefaultFusing       int //默认熔断成功率
	DefaultBusinessRate int //默认业务熔断阈值
	Recent              int //最近的qps取数
}

//Wechat wechat config
type Wechat struct {
	Host        string //微信通知id
	Chatid      string
	Msgtype     string
	Safe        int
	SendMessage bool //是否发送通知
}

// Mail mail
type Mail struct {
	Host        string
	Port        int
	Username    string
	Password    string
	NoticeOwner []string
}

// Grpc grpc
type Grpc struct {
	ProtoJavaPluginPath string
}

// Jmeter jmeter
type Jmeter struct {
	JmeterExtLibPath          string
	JmeterExtLibPathContainer string
	GRPCTemplatePath          string
	TestTimeLimit             int
	ThreadGroupPort           int //执行生成线程组的接口的端口号
	JmeterScUcodedTmp         string
	JmeterScTmp               string
	JmeterSampleTmp           string
	JmeterSamplePostTmp       string
	JmeterThGroupTmp          string
	JmeterThGroupPostTmp      string
	JmeterThGroupDuliTmp      string
	JmeterThGroupPostDuliTmp  string
	JmeterSceneTmp            string
	JSONExtractorTmp          string
}

// init init
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

// local local
func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}

// remote remote
func remote() (err error) {
	if client, err = conf.New(); err != nil {
		return
	}
	return load()
}

// load load
func load() (err error) {
	var (
		s       string
		ok      bool
		tmpConf *Config
	)
	if s, ok = client.Value("melloi.toml"); !ok {
		return errors.New("load config center error")
	}
	if _, err = toml.Decode(s, &tmpConf); err != nil {
		return errors.New("could not decode config")
	}
	*Conf = *tmpConf
	return
}
