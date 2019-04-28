package conf

import (
	"flag"
	"os"

	"go-common/library/conf"
	"go-common/library/log"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
)

// Runner ...
type Runner struct {
	URL   string
	Token string
	Name  string
}

type config struct {
	Offline bool
	Runner  []Runner
}

// Conf ...
var (
	confPath  string
	client    *conf.Client
	filename  string
	reload    chan bool
	Conf      config
	RunnerMap map[string]Runner
	HostName  string
)

func load() (err error) {
	var (
		s  string
		ok bool
	)
	if s, ok = client.Value(filename); !ok {
		err = errors.Errorf("load config center error [%s]", filename)
		return
	}

	log.Info("config data: %s", s)
	Conf = config{}
	if _, err = toml.Decode(s, &Conf); err != nil {
		err = errors.Wrapf(err, "could not decode config err(%+v)", err)
		return
	}
	return
}

func configCenter() (err error) {
	if filename == "" {
		return errors.New("filename is null,please set the environment(RUNNER_REPONAME)")
	}
	if client, err = conf.New(); err != nil {
		panic(err)
	}
	if err = load(); err != nil {
		return
	}
	client.Watch(filename)
	go func() {
		for range client.Event() {
			log.Info("config reload")

			if load() != nil {
				log.Error("config reload error (%v)", err)
			} else {
				reload <- true
			}
		}
	}()
	return
}

func osHostName() string {
	var (
		host string
		err  error
	)
	if host, err = os.Hostname(); err != nil {
		host = "unKown-Hostname"
		log.Error("%+v", errors.WithStack(err))
	}
	return host
}

func init() {
	flag.StringVar(&confPath, "conf", "", "config path")
}

// Init ...
func Init() (err error) {
	filename = os.Getenv("CONF_FILENAME")
	RunnerMap = make(map[string]Runner)
	HostName = osHostName()
	reload = make(chan bool, 1)
	if confPath == "" {
		err = configCenter()
		return
	}
	if _, err = toml.DecodeFile(confPath, &Conf); err != nil {
		log.Error("toml.DecodeFile(%s) err(%+v)", confPath, err)
		return
	}
	return
}

// ReloadEvent ...
func ReloadEvent() <-chan bool {
	return reload
}
