package v1

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
		confPath = "../cmd/config-service-example.toml"
		conf     *conf.Config
	)
	Convey("get apply", t, func() {
		_, err := toml.DecodeFile(confPath, &conf)
		So(err, ShouldBeNil)
	})
	return New(conf)
}

func TestService_CheckVersionest(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "zjx_test"
		hostname = "test_host"
		bver     = "v1.0.0"
		ip       = "123"
		version  = int64(-1)
		env      = "10"
		token    = "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc"
		appoint  = int64(97)
	)
	svr := svr(t)
	rhost := &model.Host{Service: svrName, Name: hostname, BuildVersion: bver, IP: ip, ConfigVersion: version, Appoint: appoint, Customize: "test"}
	Convey("get tag id by name", t, func() {
		event, err := svr.CheckVersion(c, rhost, env, token)
		So(err, ShouldBeNil)
		So(event, ShouldNotBeEmpty)
		Convey("get tag id by name", func() {
			e := <-event
			So(e, ShouldNotBeEmpty)
		})
	})
}

func TestService_Hosts(t *testing.T) {
	svr := svr(t)
	Convey("should get hosts", t, func() {
		_, err := svr.Hosts(context.TODO(), "zjx_test", "10")
		So(err, ShouldBeNil)
	})
}

func TestService_Config(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "zjx_test"
		hostname = "test_host"
		bver     = "v1.0.0"
		version  = int64(78)
		env      = "10"
		token    = "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc"
	)
	svr := svr(t)
	service := &model.Service{Name: svrName, BuildVersion: bver, Env: env, Token: token, Version: version, Host: hostname}
	Convey("should get hosts", t, func() {
		conf, err := svr.Config(c, service)
		So(err, ShouldBeNil)
		So(conf, ShouldNotBeEmpty)
	})
}

func TestService_Config2(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "config_test"
		hostname = "test_host"
		bver     = "shsb-docker-1"
		version  = int64(199)
		env      = "10"
		token    = "qmVUPwNXnNfcSpuyqbiIBb0H4GcbSZFV"
	)
	service := &model.Service{Name: svrName, BuildVersion: bver, Env: env, Token: token, Version: version, Host: hostname}
	svr := svr(t)
	Convey("should get Config2", t, func() {
		conf, err := svr.Config2(c, service)
		So(err, ShouldBeNil)
		So(conf, ShouldNotBeEmpty)
	})
}

func TestService_Push(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "zjx_test"
		bver    = "v1.0.0"
		version = int64(113)
		env     = "10"
	)
	service := &model.Service{Name: svrName, BuildVersion: bver, Version: version, Env: env}
	svr := svr(t)
	Convey("should get Config2", t, func() {
		err := svr.Push(c, service)
		So(err, ShouldBeNil)
	})
}

func TestService_SetToken(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "zjx_test"
		env     = "10"
		token   = "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc"
	)
	svr := svr(t)
	Convey("should get Config2", t, func() {
		svr.SetToken(c, svrName, env, token)
	})
}
func TestService_ClearHost(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "zjx_test"
		env     = "10"
	)
	svr := svr(t)
	Convey("should  clear host", t, func() {
		err := svr.ClearHost(c, svrName, env)
		So(err, ShouldBeNil)
	})
}
func TestService_VersionSuccess(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "zjx_test"
		bver    = "v1.0.0"
		env     = "10"
	)
	svr := svr(t)
	Convey("should  clear host", t, func() {
		vers, err := svr.VersionSuccess(c, svrName, env, bver)
		So(err, ShouldBeNil)
		So(vers, ShouldNotBeEmpty)
	})
}
func TestService_Builds(t *testing.T) {
	var (
		c       = context.TODO()
		svrName = "zjx_test"
		env     = "10"
		err     error
		builds  []string
	)
	svr := svr(t)
	Convey("should  clear host", t, func() {
		builds, err = svr.Builds(c, svrName, env)
		So(err, ShouldBeNil)
		So(builds, ShouldNotBeEmpty)
	})
}

func TestService_File(t *testing.T) {
	var (
		c        = context.TODO()
		svrName  = "zjx_test"
		bver     = "v1.0.0"
		env      = "10"
		fileName = "test.toml"
		token    = "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc"
		ver      = int64(74)
		err      error
	)
	service := &model.Service{Name: svrName, BuildVersion: bver, Env: env, File: fileName, Token: token, Version: ver}
	svr := svr(t)
	Convey("should  clear host", t, func() {
		_, err = svr.File(c, service)
		So(err, ShouldBeNil)
	})
}

func TestService_AddConfigs(t *testing.T) {
	svr := svr(t)
	Convey("should  clear host", t, func() {
		err := svr.AddConfigs(context.TODO(), "zjx_test", "10", "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc", "zjx", map[string]string{"aa": "bb"})
		So(err, ShouldBeNil)
	})
}

func TestService_UpdateConfigs(t *testing.T) {
	svr := svr(t)
	Convey("should  clear host", t, func() {
		err := svr.UpdateConfigs(context.TODO(), "zjx_test", "10", "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc", "zjx", 491, map[string]string{"test": "test123"})
		So(err, ShouldBeNil)
	})
}

func TestService_CopyConfigs(t *testing.T) {
	svr := svr(t)
	Convey("should  clear host", t, func() {
		_, err := svr.CopyConfigs(context.TODO(), "zjx_test", "10", "AXiLBa3Bww3inhfm6qx7g0zLY6WkLSZc", "zjx", "shsb-docker-1")
		So(err, ShouldBeNil)
	})
}

func TestService_VersionIng(t *testing.T) {
	svr := svr(t)
	Convey("should  clear host", t, func() {
		_, err := svr.VersionIng(context.TODO(), "zjx_test1", "10")
		So(err, ShouldBeNil)
	})
}
