package http

import (
	"fmt"
	"net/url"
	"testing"

	"go-common/app/admin/main/apm/model/canal"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	_listuri      = "%s/x/admin/apm/canal"
	_scanurl      = "%s/x/admin/apm/canal/scan"
	_canalediturl = "%s/x/admin/apm/canal/apply/edit"
	_addrallurl   = "%s/x/admin/apm/canal/addrs"
	_adduri       = "%s/x/admin/apm/canal/add"
	_deleteuri    = "%s/x/admin/apm/canal/delete"
	_edituri      = "%s/x/admin/apm/canal/edit"
	_applyurl     = "%s/x/admin/apm/canal/apply"
)

func TestCanalList(t *testing.T) {
	Convey("TestCanalList", t, func() {
		params := url.Values{}
		params.Set("project", "main.web-svr")
		res := new(struct {
			Code    int          `json:"code"`
			Message string       `json:"message"`
			Data    *canal.Paper `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_listuri, _domain), "", params, &res)
		t.Logf("res:%+v", res.Data)
		So(res.Code, ShouldEqual, 0)
	})
}

func TestCanalAdd(t *testing.T) {
	Convey("TestCanalAdd exists", t, func() {
		params := url.Values{}
		params.Set("addr", "172.16.33.866:3308")
		params.Set("bin_name", "fss")
		params.Set("bin_pos", "1")
		params.Set("remark", "admin")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_adduri, _domain), "", params, &res)
		t.Logf("res:%+v", res)
		//So(res.Message, ShouldContainSubstring, "exist")
		//So(res.Code, ShouldEqual, -400)

	})
}

func TestCanalEdit(t *testing.T) {
	Convey("TestCanalEdit", t, func() {
		params := url.Values{}
		params.Set("id", "8384")
		params.Set("bin_name", "rtwew")
		params.Set("bin_pos", "20")
		params.Set("remark", "inesww")
		params.Set("project", "main.web-svrq")
		params.Set("leader", "dsa")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_edituri, _domain), "", params, &res)
		t.Logf("res:%+v", res)
	})
}

func TestCanalDelete(t *testing.T) {
	Convey("TestCanalDelete", t, func() {
		params := url.Values{}
		params.Set("addr", "ewe")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_deleteuri, _domain), "", params, &res)
		t.Logf("res:%+v", res)
	})
}

func TestScanFromConfig(t *testing.T) {
	Convey("TestScanFromConfig", t, func() {
		params := url.Values{}
		params.Set("addr", "10.20.30.34:8902")
		res := new(struct {
			Code    int            `json:"code"`
			Message string         `json:"message"`
			Data    *canal.Results `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_scanurl, _domain), "", params, &res)
		t.Logf("res:%+v", res.Data.Document)

	})
}

func TestApplyLists(t *testing.T) {
	Convey("TestApplyList", t, func() {
		params := url.Values{}
		//params.Set("addr", "172.16.33.243:3308")
		params.Set("project", "main.web-svr")
		params.Set("status", "1")

		res := new(struct {
			Code    int          `json:"code"`
			Message string       `json:"message"`
			Data    *canal.Paper `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_applyurl, _domain), "", params, &res)
		t.Logf("res:%v", res.Data)

	})
}

func TestApplyConfigEdit(t *testing.T) {
	Convey("TestApplyConfigEdit", t, func() {
		params := url.Values{}
		params.Set("addr", "10.20.30.37:8902")
		params.Set("databases", _jsonstring)
		params.Set("mark", "fss")
		params.Set("user", "admin")
		params.Set("password", "admin")
		params.Set("project", "main.web-svr")
		params.Set("leader", "fss")
		res := Response{}
		_ = requests("POST", fmt.Sprintf(_canalediturl, _domain), "", params, &res)
		t.Logf("res:%v", res)

	})
}

func TestCanalAddrAll(t *testing.T) {
	Convey("TestCanalAddrAll", t, func() {
		params := url.Values{}
		res := new(struct {
			Code    int      `json:"code"`
			Message string   `json:"message"`
			Data    []string `json:"data"`
		})
		_ = requests("GET", fmt.Sprintf(_addrallurl, _domain), "", params, &res)
		t.Logf("res:%v", res.Data)
	})
}
