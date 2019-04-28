package pay

import (
	"context"
	"flag"
	"go-common/app/interface/main/creative/conf"
	"go-common/app/interface/main/creative/model/archive"
	"os"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

var (
	d *Dao
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
	d.client.SetTransport(gock.DefaultTransport)
	// d.es.client.SetTransport(gock.DefaultTransport)
	return r
}

func Test_AssReg(t *testing.T) {
	var (
		c        = context.TODO()
		err      error
		AID      = int64(10110788)
		ip       = "127.0.0.1"
		ass      *archive.PayAsset
		registed bool
	)
	convey.Convey("AssReg", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", d.assURI).Reply(200).JSON(`{"code":20061}`)
		ass, registed, err = d.Ass(c, AID, ip)
		ctx.Convey("Then err should be nil.au should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldNotBeNil)
			ctx.So(ass, convey.ShouldBeNil)
			ctx.So(registed, convey.ShouldBeFalse)
		})
	})
}
func Test_UserAcceptProtocol(t *testing.T) {
	var (
		c          = context.TODO()
		err        error
		mid        = int64(2089809)
		protocolID = "iamhashstringforp"
		accept     bool
	)
	convey.Convey("UserAcceptProtocol", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("GET", `http://uat-manager.bilibili.co/x/admin/search/query?business=log_user_action&query={"fields":["protocol_id"],"from":"log_user_action_83_all","order_score_first":true,"order_random_seed":"","highlight":false,"pn":1,"ps":2000,"order":[{"ctime":"desc"}],"where":{"eq":{"mid":2089809,"protocol_id":"iamhashstringforp"}}}`).
			Reply(200).
			JSON(`{"code":0,"message":"0","ttl":1,"data":{"order":"ctime","sort":"desc","result":[],"debug":null,"page":{"num":1,"size":2000,"total":0}}}`)
		accept, err = d.UserAcceptProtocol(c, protocolID, mid)
		ctx.Convey("Then err should be nil.au should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(accept, convey.ShouldBeFalse)
		})
	})
}
