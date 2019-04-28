package http

import (
	"context"
	"encoding/json"
	"flag"
	"net/url"
	"path/filepath"
	"testing"

	"go-common/app/service/main/archive/conf"
	"go-common/app/service/main/archive/service"
	ghttp "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	s      *service.Service
	client *ghttp.Client
)

func init() {
	dir, _ := filepath.Abs("../cmd/archive-service-test.toml")
	flag.Set("conf", dir)
	conf.Init()
	s = service.New(conf.Conf)
	Init(conf.Conf, s)
	client = ghttp.NewClient(conf.Conf.PlayerClient)
}

func Test_Archive(t *testing.T) {
	Convey("/x/internal/v2/archive", t, func() {
		p := url.Values{}
		p.Set("aid", "10098813")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Get(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}

func Test_ArchiveView(t *testing.T) {
	Convey("/x/internal/v2/archive/view", t, func() {
		p := url.Values{}
		p.Set("aid", "10098813")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Get(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive/view", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}

func Test_ArchiveViews(t *testing.T) {
	Convey("/x/internal/v2/archive/views", t, func() {
		p := url.Values{}
		p.Set("aids", "10098813,10098825,10098813")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Get(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive/views", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}

func Test_RegionArcs(t *testing.T) {
	Convey("/x/internal/v2/archive/region", t, func() {
		p := url.Values{}
		p.Set("rid", "182")
		p.Set("ps", "20")
		p.Set("pn", "1")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Get(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive/region", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}

func Test_ShareAdd(t *testing.T) {
	Convey("/x/internal/v2/archive/share/add", t, func() {
		p := url.Values{}
		p.Set("aid", "5463554")
		p.Set("mid", "1684013")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Post(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive/share/add", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}

func Test_UpCount(t *testing.T) {
	Convey("/x/internal/v2/archive/up/count", t, func() {
		p := url.Values{}
		p.Set("mid", "27515232")
		var res struct {
			Code int             `json:"code"`
			Data json.RawMessage `json:"data"`
		}
		err := client.Get(context.TODO(), "http://0.0.0.0:6081/x/internal/v2/archive/up/count", "", p, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldBeZeroValue)
		Printf("code(%d) data(%s)\n", res.Code, res.Data)
	})
}
