package conf

import (
	"flag"
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func showConf() {
	for _, v := range Conf.Runner {
		fmt.Println(v.URL, v.Token)
	}
}

func init() {
	var err error
	flag.Set("conf", "../service/agent/runners.toml")
	if err = Init(); err != nil {
		panic(err)
	}
	showConf()
}

func TestConf(t *testing.T) {
	Convey("test conf", t, func() {
		So(Conf.Runner, ShouldNotBeEmpty)
	})
}
