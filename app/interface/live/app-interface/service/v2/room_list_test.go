package v2

import (
	"flag"
	"go-common/app/interface/live/app-interface/conf"
)

var (
	s *IndexService
)

func init() {
	flag.Set("conf", "../../cmd/test.toml")
	if err := conf.Init(); err != nil {
		panic(err)
	}
	s = NewIndexService(conf.Conf)
}
