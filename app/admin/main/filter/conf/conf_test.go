package conf

import (
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-admin-test.toml")
	var err error
	if err = Init(); err != nil {
		panic(err)
	}
}

func TestConf(t *testing.T) {
	Convey("conf", t, func() {
		So(Conf.MySQL, ShouldNotBeNil)
		So(Conf.Memcache, ShouldNotBeNil)
		So(Conf.BM, ShouldNotBeNil)
		So(Conf.Log, ShouldNotBeNil)
		So(Conf.Property, ShouldNotBeNil)
		So(Conf.HTTPClient, ShouldNotBeNil)
		So(Conf.Ecode, ShouldNotBeNil)
		So(Conf.Auth, ShouldNotBeNil)
	})
}
