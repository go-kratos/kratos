package http

import (
	"flag"
	"fmt"

	"go-common/app/service/openplatform/anti-fraud/conf"
	httpx "go-common/library/net/http/blademaster"

	_ "github.com/smartystreets/goconvey/convey"
)

var client *httpx.Client

func init() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(fmt.Errorf("conf.Init() error(%v)", err))
	}
	client = httpx.NewClient(conf.Conf.HTTPClient.Read)
	Init(conf.Conf)
}
