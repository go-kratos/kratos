package service

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"go-common/app/admin/main/appstatic/conf"
	"go-common/app/admin/main/appstatic/model"

	. "github.com/smartystreets/goconvey/convey"
)

var srv *Service

func init() {
	dir, _ := filepath.Abs("../cmd/appstatic-admin-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	srv = New(conf.Conf)
	time.Sleep(time.Second)
}

func WithService(f func(s *Service)) func() {
	return func() {
		f(srv)
	}
}

func TestService_Ping(t *testing.T) {
	Convey("Ping", t, WithService(func(svf *Service) {
		svf.Ping(context.Background())
		fmt.Println("service ping successfully")
	}))
}

func TestService_Wait(t *testing.T) {
	Convey("Ping", t, WithService(func(svf *Service) {
		svf.Wait()
		fmt.Println("service wait successfully")
	}))
}

func TestService_Close(t *testing.T) {
	Convey("Close", t, WithService(func(svf *Service) {
		svf.Close()
		fmt.Println("service closed successfully")
	}))
}

func TestService_GenerateVer(t *testing.T) {
	Convey("Generate Version", t, WithService(func(svf *Service) {
		resID, version, err := svf.GenerateVer("myTestRes", &model.Limit{}, &model.FileInfo{
			Name: "testResFile",
			Size: 333,
			Type: "application/zip",
			Md5:  "333",
			URL:  "www.bilibili.com",
		}, &model.ResourcePool{
			ID:   1,
			Name: "resourcefile",
		}, 1)
		So(err, ShouldBeNil)
		So(resID, ShouldBeGreaterThan, 0)
		So(version, ShouldBeGreaterThan, 0)
	}))
}
