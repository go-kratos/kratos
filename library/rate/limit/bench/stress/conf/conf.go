package conf

import (
	"flag"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/trace"

	"github.com/BurntSushi/toml"
)

// global var
var (
	confPath string
	// Conf config
	Conf = &Config{}
)

// Config config set
type Config struct {
	// elk
	Log *log.Config
	// http
	BM *HTTPServers
	// tracer
	Tracer *trace.Config
	// redis
	Redis *redis.Config
	// memcache
	Memcache *memcache.Config
	// MySQL
	MySQL *sql.Config
	// ecode
	Ecode *ecode.Config
}

// HTTPServers Http Servers
type HTTPServers struct {
	Outer *bm.ServerConfig
	Inner *bm.ServerConfig
	Local *bm.ServerConfig
}

func init() {
	flag.StringVar(&confPath, "conf", "./stress-test.toml", "default config path")
}

// Init init conf
func Init() error {
	if confPath != "" {
		return local()
	}
	s := `# This is a TOML document. Boom

version = "1.0.0"
user = "nobody"
pid = "/tmp/stress.pid"
dir = "./"
perf = "0.0.0.0:6420"
trace = false
debug = false


[log]
#dir = "/data/log/stress"
 #[log.agent]
 # taskID = "000161"
 # proto = "unixgram"
 # addr = "/var/run/lancer/collector.sock"
 # chan = 10240

[bm]
	[bm.inner]
	addr = "0.0.0.0:9001"
	timeout = "1s"
	[bm.local]
	addr = "0.0.0.0:9002"
	timeout = "1s"`
	_, err := toml.Decode(s, &Conf)
	return err
}

func local() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
