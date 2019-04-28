package http

import (
	"flag"
	"path/filepath"
	"time"

	"go-common/app/service/main/up/conf"
)

func init() {
	dir, _ := filepath.Abs("../cmd/upcredit-service.toml")
	flag.Set("conf", dir)
	conf.Init()
	// Init(conf.Conf)
	time.Sleep(time.Second)
}
