package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"
	"go-common/library/time"

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
	Rank      *Rank
	Log       *log.Config
	BM        *bm.ServerConfig
	RPCServer *rpc.ServerConfig
	Verify    *verify.Config
	Tracer    *trace.Config
	MySQL     *DB
	Databus   *Databus
	Ecode     *ecode.Config
}

// DB .
type DB struct {
	BilibiliArchive *sql.Config
	ArchiveStat     *sql.Config
	BilibiliTV      *sql.Config
}

// Databus .
type Databus struct {
	StatView    *databus.Config
	Archive     *databus.Config
	UgcTvBinlog *databus.Config
}

// Rank .
type Rank struct {
	SwitchAll  bool
	SwitchIncr bool
	RowsLimit  int
	Ticker     time.Duration
	BatchSleep time.Duration
	BatchStep  time.Duration
	FilePath   string
	FileName   string
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
