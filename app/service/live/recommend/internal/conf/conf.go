package conf

import (
	"errors"
	"flag"

	"github.com/BurntSushi/toml"

	"go-common/library/cache/memcache"
	"go-common/library/cache/redis"
	"go-common/library/conf"
	"go-common/library/database/sql"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
	"go-common/library/net/rpc/liverpc"
	"go-common/library/net/trace"
)

var (
	confPath string
	client   *conf.Client
	// Conf config
	Conf = &Config{}
)

// Config .
type Config struct {
	Log           *log.Config
	BM            *bm.ServerConfig
	Verify        *verify.Config
	Tracer        *trace.Config
	Redis         *redis.Config
	Memcache      *memcache.Config
	MySQL         *sql.Config
	Ecode         *ecode.Config
	LiveRpc       map[string]*liverpc.ClientConfig
	Feature       *FeatureConf
	CommonFeature *CommonFeatureConf
}

// CommonFeatureConf 细分维度的特征权重配置
type CommonFeatureConf struct {
	UserAreaInterest NumMatch
	FansNum          RangeSplit
	CornerSign       ReMatch
	Online           RangeSplit
}

// NumMatch 通过int匹配获得权重
type NumMatch struct {
	Type    string
	Values  []int
	Weights []float32
}

// ReMatch 通过文字匹配获得权重
type ReMatch struct {
	Type    string
	Values  []string
	Weights []float32
}

// RangeSplit 通过分割区间获得权重
type RangeSplit struct {
	Type    string
	Values  []int64
	Weights []float32
}

// FeatureConf 特征配置
type FeatureConf struct {
	WeightVector []float32
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
