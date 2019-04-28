package conf

import (
	"go-common/app/interface/bbq/app-bbq/api/http/v1"
	"go-common/library/cache/redis"
	"go-common/library/conf/paladin"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	"go-common/library/log/infoc"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	"go-common/library/net/http/blademaster/middleware/auth"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/warden"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

var (
	// Conf config
	Conf = &Config{}
	// App setting
	App = &AppSetting{}
	// Filter .
	Filter = &UploadFilter{}
)

// Config .
type Config struct {
	Log        *log.Config
	BM         *bm.ServerConfig
	Verify     *verify.Config
	Auth       *auth.Config
	Tracer     *trace.Config
	Redis      *redis.Config
	MySQL      *sql.Config
	DMMySQL    *sql.Config
	Ecode      *ecode.Config
	HTTPClient *HTTPClient
	GRPCClient map[string]*GRPCConf
	AntiSpam   map[string]*antispam.Config
	Tmap       map[string]string
	URLs       map[string]string
	Comment    *Comment
	Infoc      *infoc.Config
	Search     *Search
	Notices    []*v1.NoticeOverview
	Upload     *Upload
}

//Upload ..
type Upload struct {
	HTTPSchema string
}

// Set .
func (c *Config) Set(text string) error {
	if _, err := toml.Decode(text, c); err != nil {
		panic(err)
	}
	if c.Redis != nil {
		for _, anti := range c.AntiSpam {
			anti.Redis = c.Redis
		}
	}
	return nil
}

// Comment 评论配置
type Comment struct {
	Type       int64
	DebugID    int64
	CloseRead  bool
	CloseWrite bool
}

// Search 搜索配置
type Search struct {
	Host string
}

// HTTPClient conf
type HTTPClient struct {
	Normal *bm.ClientConfig
	Slow   *bm.ClientConfig
}

//GRPCConf .
type GRPCConf struct {
	WardenConf *warden.ClientConfig
	Addr       string
}

// Init init conf
func Init() (err error) {
	if err = paladin.Init(); err != nil {
		return
	}
	if err = paladin.Watch("video-c.toml", Conf); err != nil {
		return
	}
	if err = paladin.Watch("app_setting.toml", App); err != nil {
		return
	}
	if err = paladin.Watch("upload_filter.toml", Filter); err != nil {
		return
	}
	return
}
