package bfs

import (
	"context"
	"flag"
	"github.com/bouk/monkey"
	"github.com/smartystreets/goconvey/convey"
	"go-common/app/interface/main/creative/conf"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"testing"
)

var (
	d          *Dao
	defaultImg = "https://avatars3.githubusercontent.com/u/12002442"
)

func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.archive.creative")
		flag.Set("conf_token", "96b6a6c10bb311e894c14a552f48fef8")
		flag.Set("tree_id", "2305")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/creative.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	d = New(conf.Conf)
	m.Run()
	os.Exit(0)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	d.client.Transport = gock.DefaultTransport
	return r
}

func TestBfsUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileType = ""
			bs       = []byte("")
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			httpMock(d.c.BFS.Method, d.c.BFS.URL).Reply(200).SetHeaders(mockHeader).JSON("mockByte")
			location, err := d.Upload(c, fileType, bs)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBfsUploadByFile(t *testing.T) {
	convey.Convey("UploadByFile", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			imgpath = defaultImg
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			monkey.Patch(ioutil.ReadFile, func(_ string) ([]byte, error) {
				return []byte{}, nil
			})
			httpMock(_method, _url).Reply(200).SetHeaders(mockHeader)
			location, err := d.UploadByFile(c, imgpath)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBfsCapture(t *testing.T) {
	convey.Convey("Capture", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			url = defaultImg
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			httpMock("GET", url).Reply(200).SetHeaders(mockHeader)
			loc, size, err := d.Capture(c, url)
			ctx.Convey("Then err should be nil.loc,size should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(size, convey.ShouldNotBeNil)
				ctx.So(loc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestBfscheckURL(t *testing.T) {
	convey.Convey("checkURL", t, func(ctx convey.C) {
		var (
			url = defaultImg
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := checkURL(url)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestBfsisPublicIP(t *testing.T) {
	convey.Convey("isPublicIP", t, func(ctx convey.C) {
		var IP = net.ParseIP("127.0.0.1")
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			p1 := isPublicIP(IP)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}
