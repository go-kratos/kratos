package conf

import (
	"errors"
	"flag"

	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/rpc"
	"go-common/library/net/trace"
	"go-common/library/queue/databus"

	"go-common/library/database/hbase.v2"

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
	Log    *log.Config
	BM     *bm.ServerConfig
	Tracer *trace.Config
	Ecode  *ecode.Config

	AccountSummaryHBase    *hbase.Config
	MemberBinLog           *databus.Config
	BlockBinLog            *databus.Config
	PassportBinLog         *databus.Config
	RelationBinLog         *databus.Config
	AccountSummaryProducer *databus.Config
	MemberService          *rpc.ClientConfig
	RelationService        *rpc.ClientConfig
	HTTPClient             *bm.ClientConfig
	Host                   *Host
	FeatureGate            *FeatureGate
	AccountSummary         *AccountSummary

	MemberDB   *sql.Config
	RelationDB *sql.Config
	PassportDB *sql.Config
}

// AccountSummary is
type AccountSummary struct {
	SubProcessWorker   uint64
	SyncRangeStart     int64
	SyncRangeEnd       int64
	SyncRangeWorker    uint64
	InitialWriteWorker uint64
}

// FeatureGate is
type FeatureGate struct {
	DisableSubProcess     bool
	Initial               bool
	InitialMemberBase     bool
	InitialMemberExp      bool
	InitialMemberOfficial bool
	InitialRelationStat   bool
	InitialBlock          bool
	InitialPassport       bool

	SyncRange bool
}

// Host is
type Host struct {
	Passport string
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
