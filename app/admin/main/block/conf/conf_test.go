package conf

import (
	"flag"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	defer os.Exit(0)
	flag.Set("conf", "../cmd/block-admin-test.toml")
	var err error
	if err = Init(); err != nil {
		panic(err)
	}

	m.Run()
}

func TestConf(t *testing.T) {
	Convey("conf", t, func() {
		So(Conf.Property.MSG, ShouldNotBeNil)
		So(Conf.Property.MSG.BlackHouseForever, ShouldNotBeNil)
		So(Conf.Property.MSG.BlackHouseLimit, ShouldNotBeNil)
		So(Conf.Property.MSG.BlockRemove, ShouldNotBeNil)
		So(Conf.Property.MSG.SysForever, ShouldNotBeNil)
		So(Conf.Property.MSG.SysLimit, ShouldNotBeNil)
		t.Log(Conf.Property.MSG)
	})
}
