package conf

import (
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/database/orm"
	"go-common/library/database/sql"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/antispam"
	v "go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/time"

	"github.com/BurntSushi/toml"
)

var (
	confPath string
	// Conf init config
	Conf *Config
)

// Config config.
type Config struct {
	// log
	Log *log.Config
	//rpc server2
	RPCServer *rpc.ServerConfig
	// db
	DB *sql.Config
	// redis
	Redis *redis.Config
	// timeout
	PollTimeout time.Duration
	// local cache
	PathCache string
	// orm
	ORM *orm.Config
	//BM
	BM *bm.ServerConfig
	// Antispam
	Antispam *antispam.Config
	Verify   *v.Config
}

func init() {
	flag.StringVar(&confPath, "conf", "./config-service-example.toml", "config path")
}

// Init init.
func Init() (err error) {
	_, err = toml.DecodeFile(confPath, &Conf)
	return
}
