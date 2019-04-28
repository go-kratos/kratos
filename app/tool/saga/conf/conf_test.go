package conf

import (
	"flag"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func init() {
	var err error
	flag.Set("conf", "../cmd/saga-test.toml")
	if err = Init(); err != nil {
		panic(err)
	}
}

func TestConf(t *testing.T) {
	Convey("test conf", t, func() {
		So(Conf, ShouldNotBeNil)
		So(Conf.Property, ShouldNotBeNil)
		So(Conf.Property.Gitlab, ShouldNotBeNil)
		So(Conf.Property.Repos, ShouldNotBeNil)
		So(Conf.Property.Repos[0].MinReviewer, ShouldEqual, 1)
		So(Conf.Property.Repos[1].MinReviewer, ShouldEqual, 0)
		So(Conf.Property.Repos[0].AuthBranches[0], ShouldEqual, "master")
		So(Conf.Property.Repos[1].AuthBranches[0], ShouldEqual, "master")
		So(Conf.Property.Mail, ShouldNotBeNil)
		So(Conf.Property.HealthCheck, ShouldNotBeNil)
		So(Conf.Property.HealthCheck.AlertAddrs, ShouldNotBeEmpty)
		So(Conf.Property.Wechat, ShouldNotBeNil)
	})
}
