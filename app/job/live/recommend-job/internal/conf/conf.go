package conf

import (
	"errors"
	"flag"

	"go-common/library/cache/redis"
	"go-common/library/conf"
	ecode "go-common/library/ecode/tip"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/middleware/verify"
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
	Log         *log.Config
	BM          *bm.ServerConfig
	Verify      *verify.Config
	Tracer      *trace.Config
	Redis       *redis.Config
	Ecode       *ecode.Config
	ItemCFJob   *JobConfig
	Hadoop      *HadoopConfig
	UserAreaJob *JobConfig
}

// HadoopConfig ...
type HadoopConfig struct {
	HadoopDir string
	TarUrl    string
}

// JobConfig ...
type JobConfig struct {
	Schedule string
	// 多少个goroutine同时去写redis数据
	WorkerNum int
	// 如果指定了文件，使用这个文件
	// 如果没有指定，自动去下载文件
	// 可以是本地文件地址，或者http文件地址
	InputFile string
	// 在hadoop里面文件的路径，带日期
	HadoopFile string
	// Hadoop下载到本地的路径，带日期
	LocalTmpFile string
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
