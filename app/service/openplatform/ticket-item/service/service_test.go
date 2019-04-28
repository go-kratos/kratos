package service

import (
	"flag"
	"fmt"
	"path/filepath"

	"go-common/app/service/openplatform/ticket-item/conf"
)

func init() {
	dir, _ := filepath.Abs("../cmd/item.toml")
	flag.Set("conf", dir)
	if err := conf.Init(); err != nil {
		panic(fmt.Sprintf("conf.Init() error(%v)", err))
	}
	s = New(conf.Conf)
}
