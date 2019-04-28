package v2

import (
	"context"
	"testing"

	"go-common/app/infra/config/conf"
	"go-common/app/infra/config/model"

	"github.com/BurntSushi/toml"
	. "github.com/smartystreets/goconvey/convey"
)

func svr(t *testing.T) *Service {
	var (
		confPath = "../../cmd/config-service-example.toml"
		conf     *conf.Config
	)
	Convey("get apply", t, func() {
		_, err := toml.DecodeFile(confPath, &conf)
		So(err, ShouldBeNil)
	})
	return New(conf)
}

func TestService_CheckVersion(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "2888_dev_sh001"
		hostname = "test_host"
		bver     = "server-1"
		ip       = "123"
		version  = int64(-1)
		token    = "47fbbf237b7f11e8992154e1ad15006a"
		appoint  = int64(-1)
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	rhost := &model.Host{Service: svrName, Name: hostname, BuildVersion: bver, IP: ip, ConfigVersion: version, Appoint: appoint, Customize: "test"}
	Convey("get tag id by name", t, func() {
		event, err := svr.CheckVersion(c, rhost, token)
		So(err, ShouldBeNil)
		So(event, ShouldNotBeEmpty)
		Convey("get tag id by name", func() {
			e := <-event
			So(e, ShouldNotBeEmpty)
		})
	})
}

func TestService_Hosts(t *testing.T) {
	var (
		svrName = "2888_dev_sh001"
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should get hosts", t, func() {
		_, err := svr.Hosts(context.TODO(), svrName)
		So(err, ShouldBeNil)
	})
}

func TestService_Config(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "2888_dev_sh001"
		version = int64(49)
		token   = "47fbbf237b7f11e8992154e1ad15006a"
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should get hosts", t, func() {
		conf, err := svr.Config(c, svrName, token, version, nil)
		So(err, ShouldBeNil)
		So(conf, ShouldNotBeEmpty)
	})
}

func TestService_Push(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "2888_dev_sh001"
		bver    = "server-1"
		version = int64(49)
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	service := &model.Service{Name: svrName, BuildVersion: bver, Version: version}
	Convey("should get Config2", t, func() {
		err := svr.Push(c, service)
		So(err, ShouldBeNil)
	})
}

func TestService_SetToken(t *testing.T) {
	var (
		svrName = "2888_dev_sh001"
		token   = "47fbbf237b7f11e8992154e1ad15006a"
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should get Config2", t, func() {
		svr.SetToken(svrName, token)
	})
}

func TestService_ClearHost(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "2888_dev_sh001"
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should  clear host", t, func() {
		err := svr.ClearHost(c, svrName)
		So(err, ShouldBeNil)
	})
}
func TestService_VersionSuccess(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "2888_dev_sh001"
		bver    = "server-1"
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err := svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should  clear host", t, func() {
		vers, err := svr.VersionSuccess(c, svrName, bver)
		So(err, ShouldBeNil)
		So(vers, ShouldNotBeEmpty)
	})
}
func TestService_Builds(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "2888_dev_sh001"
		err     error
		builds  []string
		tmp     string
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err = svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	Convey("should  clear host", t, func() {
		builds, err = svr.Builds(c, svrName)
		So(err, ShouldBeNil)
		So(builds, ShouldNotBeEmpty)
	})
}

func TestService_File(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "2888_dev_sh001"
		bver     = "server-1"
		fileName = "test.toml"
		token    = "47fbbf237b7f11e8992154e1ad15006a"
		ver      = int64(49)
		err      error
		tmp      string
	)
	svr := svr(t)
	Convey("AppService", t, func() {
		tmp, err = svr.AppService("sh001", "dev", "47fbbf237b7f11e8992154e1ad15006a")
		So(err, ShouldBeNil)
		So(tmp, ShouldEqual, svrName)
	})
	service := &model.Service{Name: svrName, BuildVersion: bver, File: fileName, Token: token, Version: ver}
	Convey("should  clear host", t, func() {
		_, err = svr.File(c, service)
		So(err, ShouldBeNil)
	})
}
