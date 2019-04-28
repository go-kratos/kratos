package conf

import (
	"flag"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-service-test.toml")
	var err error
	if err = Init(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestConf(t *testing.T) {
	Convey("Test conf", t, func() {
		So(Conf.MySQL, ShouldNotBeNil)
		So(Conf.Memcache, ShouldNotBeNil)
		So(Conf.BM, ShouldNotBeNil)
		So(Conf.RPCServer, ShouldNotBeNil)
		So(Conf.Log, ShouldNotBeNil)

		So(Conf.Property, ShouldNotBeNil)
		t.Log(Conf.Property.MaxHitSplitSize)
		So(Conf.Property.AI, ShouldNotBeNil)
		So(Conf.Property.AIHost, ShouldNotBeNil)
		t.Log(Conf.Property.AI)
		t.Log(Conf.Property.AIHost)
	})
}
