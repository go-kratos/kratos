package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	xtime "go-common/library/time"

	"github.com/BurntSushi/toml"
)

// Conf global variable.
var (
	Conf     = &Config{}
	client   *conf.Client
	confPath string
)

// Config struct of conf.
type Config struct {
	// base
	// log
	Log *log.Config
	// Databus cfg
	ContentSub       *databus.Config
	ArchiveNotifySub *databus.Config
	UgcSub           *databus.Config
	// tracer
	Tracer *trace.Config
	// http
	HTTPServer *bm.ServerConfig
	// db
	Mysql *sql.Config
	// sync params
	Sync *Sync
	// memcache
	Memcache *Memcache
	// redis
	Redis *Redis
	// playControl related config
	PlayControl *PlayControl
	// HTTPClient .
	HTTPClient *bm.ClientConfig
	// Search Cfg
	Search *Search
	// UgcSync cfg
	UgcSync *UgcSync
	// grpc
	ArcClient  *warden.ClientConfig
	ArchiveRPC *rpc.ClientConfig
	AccClient  *warden.ClientConfig
	Cfg        *Cfg
	Report     *Report
	DpClient   *bm.ClientConfig
	Style      *Style
}

// Style .
type Style struct {
	LabelSpan xtime.Duration
	StyleSpan xtime.Duration
}

// Report data .
type Report struct {
	ReportURI     string
	UpDataURI     string
	TimeDelay     string
	SendDataDelay xtime.Duration
	Env           string
	RoutineCount  int
	ReadSize      int
	Expire        xtime.Duration
	SeTimeSpan    xtime.Duration
	CronAc        string
	CronAd        string
	CronPd        string
	CronVe        string
}

// Cfg contains various of configuration
type Cfg struct {
	TitleFilter  []string
	LessStrategy int
	PgcTypes     []string // pgc types name, need to filter these ugc archives
	PGCZonesID   []int    // all the zones' ID that need to be loaded
	UgcZones     map[string]*UgcType
	SyncRetry    *SyncRetry // sync retry
	Merak        *Merak
}

// Merak cfg
type Merak struct {
	Host     string
	Key      string
	Secret   string
	Names    []string
	Template string
	Title    string
	Cron     string
	Onlyfree bool // if true, we consider free+audited episodes, otherwise we consider only audited episodes
}

// SyncRetry def.
type SyncRetry struct {
	MaxRetry int // max retry times for pgc already-passed sn & ep
	RetryFre xtime.Duration
}

// UgcType def.
type UgcType struct {
	TID  int32
	Name string
}

// Search represents the config for the search suggestion module
type Search struct {
	UgcSwitch      string // the Ugc search suggest Switch
	SugPath        string // the tvsug file local path
	Md5Path        string // the tvsug md5 file local path
	FTP            *FTP   // the ftp info
	PgcContPath    string // the pgc content file local path
	PgcContMd5Path string // the pgc content md5 file local path
	UgcContPath    string // the ugc content file local path
	UgcContMd5Path string // the ugc content md5 file local path
	Cfg            *SearchCfg
}

//SearchCfg synchronize files time
type SearchCfg struct {
	UploadFre xtime.Duration
}

// FTP represents the ftp login info
type FTP struct {
	Pass             string
	User             string
	Host             string
	URL              string
	Timeout          xtime.Duration // timeout in seconds
	UseEPSV          bool
	RemoteFName      string // file name in remote ftp server
	RemoteMd5        string // md5 file name in remote ftp server
	RemotePgcCont    string // pgc file name in remote ftp server
	RemotePgcURL     string // RemotePgcURL remote search pgc url dir
	RemotePgcContMd5 string // pgc md5 file name in remote ftp server
	RemoteUgcCont    string // ugc file name in remote ftp server
	RemoteUgcURL     string // RemotePgcCont remote search ugc url dir
	RemoteUgcContMd5 string // ugc md5 file name in remote ftp server
}

// Redis redis
type Redis struct {
	*redis.Config
	Expire  xtime.Duration
	CronPGC string
	CronUGC string
}

// PlayControl is the configuration for the play control interface, related to MC
type PlayControl struct {
	ProducerCron string
	PieceSize    int
}

// Memcache config
type Memcache struct {
	*memcache.Config
	Expire      xtime.Duration
	ExpireMedia xtime.Duration
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

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init int config
func Init() error {
	if confPath != "" {
		return local()
	}
	return remote()
}
