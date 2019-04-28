package http

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"path/filepath"
	"testing"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/service"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	_domain      = "http://127.0.0.1:8000"
	_contentType = "application/x-www-form-urlencoded"
	_cookie      = "username=fengshanshan; _AJSESSIONID=cf400491a236da90f27fb1b9bb9c4e2d; sven-apm=994ae2b6d290e584488443f9cc4733fbee7a88a4cd376135b1295f4cf81231de"
	_realIP      = "172.16.33.134"
	_configuri   = "%s/x/admin/apm/canal/apply/config"
	_jsonstring  = `[{ "schema":"123","table":[  {"name":"abc","primarykey":["order_id","new_id"],"omitfield":["new","old"]} , {"name":"def","primarykey":["order_id","new_id"],"omitfield":["new","old"] } ,{"name":"sfg","primarykey":["order_id","new_id"],"omitfield":["new","old"]} ],"databus": { "group": "LiveTime-LiveLive-P","addr": "172.16.33.158:6205"}},{ "schema":"456","table":[  {"name":"abc" ,"primarykey":["order_id","new_id"],"omitfield":["new","old"]} , {"name":"def" } ,{"name":"sfg"} ], "databus": {"group": "AccAnswer-MainManager-S","addr": "172.16.33.158:6205"}}]`
)

func init() {
	dir, _ := filepath.Abs("../cmd/apm-admin-test.toml")
	flag.Parse()
	flag.Set("conf", dir)
	conf.Init()
	log.Init(conf.Conf.Log)
	apmSvc = service.New(conf.Conf)
	Init(conf.Conf, apmSvc)
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data"`
}

func requests(method, uri, realIP string, params url.Values, res interface{}) (err error) {

	client := xhttp.NewClient(conf.Conf.HTTPClient)
	req, err := client.NewRequest(method, uri, realIP, params)
	if err != nil {
		return
	}
	req.Header.Set("X-BACKEND-BILI-REAL-IP", _realIP)
	req.Header.Set("Content-Type", _contentType)
	req.Header.Set("Cookie", _cookie)
	if err = client.Do(context.TODO(), req, &res); err != nil {
		return
	}
	return
}

func TestApplyDetailToConfig(t *testing.T) {
	Convey("TestApply register nonconfig and no canal_apply", t, func() {
		params := url.Values{}
		params.Set("addr", "10.20.30.34:8902")
		params.Set("databases", _jsonstring)
		params.Set("mark", "demo")
		params.Set("user", "admin")
		params.Set("password", "admin")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")

		res := Response{}
		_ = requests("POST", fmt.Sprintf(_configuri, _domain), "", params, &res)
		So(res.Code, ShouldEqual, 70015)
	})

	Convey("TestApplyAddrIllegal", t, func() {
		params := url.Values{}
		params.Set("addr", "172.16.33.2553308")
		params.Set("databases", _jsonstring)
		params.Set("mark", "demo")
		params.Set("user", "admin")
		params.Set("passwd", "admin")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")

		res := Response{}
		_ = requests("POST", fmt.Sprintf(_configuri, _domain), "", params, &res)
		So(res.Code, ShouldEqual, 70002)
		//So(res.Message, ShouldContainSubstring, "addr参数不合法")
	})

	Convey("TestApplyExist", t, func() {
		params := url.Values{}
		params.Set("addr", "10.20.30.34:8902")
		params.Set("databases", _jsonstring)
		params.Set("mark", "demo")
		params.Set("user", "admin")
		params.Set("passwd", "admin")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_configuri, _domain), "", params, &res)
		So(res.Code, ShouldEqual, 0)
		//So(res.Message, ShouldContainSubstring, "己提交申请")
	})

	Convey("TestApplyAddRequestParamError", t, func() {
		params := url.Values{}
		params.Set("leader", "fss")
		params.Set("remark", "test")
		params.Set("project", "main.web-svr")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_configuri, _domain), "", params, &res)
		So(res.Code, ShouldEqual, -400)
	})
}

func TestApplyExist(t *testing.T) {
	Convey("TestApply register nonconfig and no canal_apply", t, func() {
		params := url.Values{}
		params.Set("addr", "10.20.30.34:8902")
		params.Set("databases", _jsonstring)
		params.Set("mark", "demo")
		params.Set("user", "a")
		params.Set("password", "a")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")

		res := Response{}
		_ = requests("POST", fmt.Sprintf(_configuri, _domain), "", params, &res)
		t.Logf("%+v", res)
		So(res.Code, ShouldEqual, 70015)
	})
}
