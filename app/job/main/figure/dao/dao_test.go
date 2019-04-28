package dao

import (
	"flag"

	"path/filepath"
	"sync"
	"time"

	"go-common/app/job/main/figure/conf"
	"go-common/library/log"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	once sync.Once
	d    *Dao
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
	log.Init(conf.Conf.Log)
	defer log.Close()
}

func startService() {
	initConf()
	if d == nil {
		d = New(conf.Conf)
	}
	time.Sleep(time.Second * 2)
}

func CleanDB() {
}

func init() {
	dir, _ := filepath.Abs("../cmd/figure-job-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	if d == nil {
		d = New(conf.Conf)
	}
}

func WithDao(f func(d *Dao)) func() {
	return func() {
		Reset(func() { CleanDB() })
		f(d)
	}
}
