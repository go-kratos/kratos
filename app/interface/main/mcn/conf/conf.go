package conf

import (
	"errors"
	"flag"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/cache/memcache"
	"go-common/library/conf"
	"go-common/library/database/orm"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"
	"go-common/library/time"

	"go-common/app/interface/main/mcn/tool/datacenter"

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
	Log      *log.Config
	BM       *bm.ServerConfig
	Tracer   *trace.Config
	Memcache *MemcacheConfig
	MCNorm   *orm.Config
	Ecode    *ecode.Config
	BFS      *BFS
	Host     *Host
	// http client
	HTTPClient *bm.ClientConfig
	Property   *Property
	Other      *OtherConfig
	GRPCClient *GRPCClient
	RankCache  *GCacheConfig
	//upload Antispam
	UploadAntispam *antispam.Config
	DataClientConf *datacenter.ClientConfig
}

// AfterLoad .
func (s *Config) AfterLoad() {
	s.Other.WhiteListMidMap = make(map[int64]struct{}, len(s.Other.WhiteListMid))
	for _, v := range s.Other.WhiteListMid {
		s.Other.WhiteListMidMap[v] = struct{}{}
	}
}

// GRPCClient .
type GRPCClient struct {
	Tag     *warden.ClientConfig
	Account *warden.ClientConfig
	Member  *warden.ClientConfig
	Archive *warden.ClientConfig
}

// MemcacheConfig .
type MemcacheConfig struct {
	memcache.Config
	McnSignCacheExpire time.Duration
	McnDataCacheExpire time.Duration
}

// Property .
type Property struct {
	MSG []*model.MSG
}

// BFS bfs config
type BFS struct {
	Bucket string
	Key    string
	Secret string
}

// Host host config .
type Host struct {
	Bfs     string
	Msg     string
	Videoup string
	API     string
}

//OtherConfig some config.
type OtherConfig struct {
	Debug                       bool
	PublicationPriceChangeLimit time.Duration
	WhiteListMid                []int64 // 超级查看权限
	WhiteListMidMap             map[int64]struct{}
}

//IsWhiteList check is in white list
func (s *OtherConfig) IsWhiteList(mid int64) bool {
	_, ok := s.WhiteListMidMap[mid]
	return ok
}

//GCacheConfig gcache
type GCacheConfig struct {
	Size                    int           // gcache
	ExpireTime              time.Duration // key expire time
	RecommendPoolExpireTime time.Duration
}

func init() {
	flag.StringVar(&confPath, "conf", "", "default config path")
}

// Init init conf
func Init() error {
	defer Conf.AfterLoad()
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
