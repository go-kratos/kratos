package service

import (
	"context"
	"testing"

	"go-common/app/admin/main/config/conf"

	"github.com/BurntSushi/toml"
	. "github.com/smartystreets/goconvey/convey"
)

func svr(t *testing.T) (svr *Service) {
	var (
		confPath = "../cmd/config-admin-example.toml"
		conf     *conf.Config
	)
	Convey("should decodeFile file", t, func() {
		_, err := toml.DecodeFile(confPath, &conf)
		So(err, ShouldBeNil)
	})
	return New(conf)
}

func TestService_UpdateToken(t *testing.T) {
	svr := svr(t)
	Convey("should update token", t, func() {
		err := svr.UpdateToken(context.Background(), "dev", "sh001", 2888)
		So(err, ShouldBeNil)
	})
}

func TestService_AppByName(t *testing.T) {
	svr := svr(t)
	Convey("should get app by name", t, func() {
		app, err := svr.AppByTree(2888, "dev", "sh001")
		So(err, ShouldBeNil)
		So(app, ShouldNotBeEmpty)
	})
}
