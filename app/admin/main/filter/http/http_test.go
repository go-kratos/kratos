package http

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"testing"
	"time"

	"go-common/app/admin/main/filter/conf"
	"go-common/app/admin/main/filter/model"
	"go-common/app/admin/main/filter/service"
	"go-common/library/log"
	xhttp "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	domain = "localhost:7382"
	// domain = "api.bilibili.co"

	_filterAdd    = "http://" + domain + "/x/admin/filter/add"
	_filterDel    = "http://" + domain + "/x/admin/filter/del"
	_filterList   = "http://" + domain + "/x/admin/filter/list"
	_filterSearch = "http://" + domain + "/x/admin/filter/search"
	_filterLog    = "http://" + domain + "/x/admin/filter/log"
	_filterEdit   = "http://" + domain + "/x/admin/filter/edit"
	_filterGet    = "http://" + domain + "/x/admin/filter/get"
	// _filterOrigin  = "http://" + domain + "/x/admin/filter/origin"
	// _filterOrigins = "http://" + domain + "/x/admin/filter/origins"

	_filterKeyAdd      = "http://" + domain + "/x/admin/filter/key/add"
	_filterKeyDel      = "http://" + domain + "/x/admin/filter/key/del"
	_filterKeyEditInfo = "http://" + domain + "/x/admin/filter/key/editinfo"
	_filterKeyEdit     = "http://" + domain + "/x/admin/filter/key/edit"
	_filterKeySearch   = "http://" + domain + "/x/admin/filter/key/search"
	_filterKeyLog      = "http://" + domain + "/x/admin/filter/key/log"

	_whiteAdd      = "http://" + domain + "/x/admin/filter/white/add"
	_whiteDel      = "http://" + domain + "/x/admin/filter/white/del"
	_whiteSearch   = "http://" + domain + "/x/admin/filter/white/search"
	_whiteEditInfo = "http://" + domain + "/x/admin/filter/white/editinfo"
	_whiteEdit     = "http://" + domain + "/x/admin/filter/white/edit"
	_whiteLog      = "http://" + domain + "/x/admin/filter/white/log"
)

var (
	client *xhttp.Client
	ctx    = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-admin-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	client = xhttp.NewClient(conf.Conf.HTTPClient.Normal)

	log.Init(conf.Conf.Log)
	srv := service.New(conf.Conf)
	Init(srv)
	time.Sleep(time.Second)

	m.Run()
}

func TestFilter(t *testing.T) {
	var fid int64
	Convey("Filter add", t, func() {
		params := url.Values{}
		params.Set("rule", "unit test")
		params.Set("area", "common,reply")
		params.Set("mode", fmt.Sprintf("%d", model.StrMode))
		params.Set("level", "20")
		params.Set("comment", "unit test")
		params.Set("adid", "2333")
		params.Set("name", "2333")
		params.Set("stime", fmt.Sprintf("%d", time.Now().Unix()))
		params.Set("etime", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
		params.Set("tpid", "0")
		params.Set("source", "0")
		params.Set("key_type", "0")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterAdd, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Filter search", t, func() {
		params := url.Values{}
		params.Set("msg", "unit test")
		params.Set("area", "common")
		params.Set("source", "0")
		params.Set("filter_type", "0")
		params.Set("stage", "0")
		params.Set("pn", "1")
		params.Set("ps", "10")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rules []model.FilterInfo `json:"rules"`
				Total int                `json:"total"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filterSearch, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data.Total, ShouldEqual, 1)
		So(res.Data.Rules, ShouldHaveLength, 1)
		So(res.Data.Rules[0].ID, ShouldBeGreaterThan, 0)
		fid = res.Data.Rules[0].ID
	})

	Convey("Filter get", t, func() {
		params := url.Values{}
		params.Set("id", fmt.Sprintf("%d", fid))

		var res struct {
			Code int              `json:"code"`
			Data model.FilterInfo `json:"data"`
		}

		err := client.Get(ctx, _filterGet, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data.ID, ShouldEqual, fid)
		So(res.Data.Filter, ShouldEqual, "unit test")
		So(res.Data.Level, ShouldEqual, 20)
		So(res.Data.Mode, ShouldEqual, model.StrMode)
	})

	Convey("Filter log", t, func() {
		params := url.Values{}
		params.Set("id", fmt.Sprintf("%d", fid))

		var res struct {
			Code int         `json:"code"`
			Data []model.Log `json:"data"`
		}

		err := client.Get(ctx, _filterLog, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeEmpty)
	})

	Convey("Filter edit", t, func() {
		params := url.Values{}
		params.Set("rule", "unit test")
		params.Set("area", "common,article")
		params.Set("mode", fmt.Sprintf("%d", model.RegMode))
		params.Set("level", "30")
		params.Set("comment", "unit test")
		params.Set("adid", "2333")
		params.Set("reason", "unit test edit")
		params.Set("id", fmt.Sprintf("%d", fid))
		params.Set("name", "2333")
		params.Set("stime", fmt.Sprintf("%d", time.Now().Unix()))
		params.Set("etime", fmt.Sprintf("%d", time.Now().Add(time.Second).Unix()))
		params.Set("tpid", "0")
		params.Set("source", "0")
		params.Set("key_type", "0")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterEdit, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Filter get edited", t, func() {
		params := url.Values{}
		params.Set("id", fmt.Sprintf("%d", fid))

		var res struct {
			Code int              `json:"code"`
			Data model.FilterInfo `json:"data"`
		}

		err := client.Get(ctx, _filterGet, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data.ID, ShouldEqual, fid)
		So(res.Data.Level, ShouldEqual, 30)
		So(res.Data.Mode, ShouldEqual, model.RegMode)
		So(res.Data.Filter, ShouldEqual, "unit test")
	})

	Convey("Filter del", t, func() {
		params := url.Values{}
		params.Set("fid", fmt.Sprintf("%d", fid))
		params.Set("adid", "2333")
		params.Set("name", "2333")
		params.Set("reason", "unit test del")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterDel, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Filter list", t, func() {
		params := url.Values{}
		params.Set("ps", "1")
		params.Add("pn", "20")
		params.Set("area", "common")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rules []model.FilterInfo `json:"rules"`
				Total int                `json:"total"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filterList, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
	})
}

// func testOrigin(t *testing.T) {
// 	Convey("Filter origin", t, func() {
// 		params := url.Values{}
// 		params.Set("id", "1")
// 		params.Set("area", "article")

// 		var res struct {
// 			Code int           `json:"code"`
// 			Data model.Message `json:"data"`
// 		}

// 		err := client.Get(ctx, _filterOrigin, "", params, &res)
// 		So(err, ShouldBeNil)
// 		So(res, ShouldNotBeNil)
// 		So(res.Code, ShouldEqual, 0)
// 	})

// 	Convey("Filter origins", t, func() {
// 		params := url.Values{}
// 		params.Set("ids", "1,2,3")
// 		params.Set("area", "reply")

// 		var res struct {
// 			Code int                     `json:"code"`
// 			Data map[int64]model.Message `json:"data"`
// 		}

// 		err := client.Get(ctx, _filterOrigins, "", params, &res)
// 		So(err, ShouldBeNil)
// 		So(res, ShouldNotBeNil)
// 		So(res.Code, ShouldEqual, 0)
// 	})
// }

func TestKey(t *testing.T) {
	var fid int64
	Convey("Key add", t, func() {
		params := url.Values{}
		params.Set("area", "danmu")
		params.Set("key", "aid:1")
		params.Set("rule", "unit test")
		params.Set("mode", fmt.Sprintf("%d", model.StrMode))
		params.Set("level", "20")
		params.Set("comment", "unit test")
		params.Set("adid", "2333")
		params.Set("name", "2333")
		params.Set("stime", fmt.Sprintf("%d", time.Now().Unix()))
		params.Set("etime", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterKeyAdd, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Key search", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("comment", "unit test")
		params.Set("ps", "10")
		params.Set("pn", "1")
		params.Set("state", fmt.Sprintf("%d", model.FilterStateNormal))

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rules []model.FilterInfo `json:"rules"`
				Total int                `json:"totle"`
				PN    int                `json:"pn"`
				PS    int                `json:"ps"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filterKeySearch, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data.PN, ShouldEqual, 1)
		So(res.Data.Rules, ShouldHaveLength, 1)
		So(res.Data.Rules[0].ID, ShouldBeGreaterThan, 0)
		fid = res.Data.Rules[0].ID
	})

	Convey("Key editinfo", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("fid", fmt.Sprintf("%d", fid))

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rule model.KeyInfo `json:"rule"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filterKeyEditInfo, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data.Rule, ShouldNotBeNil)
		So(res.Data.Rule.ID, ShouldEqual, fid)
		So(res.Data.Rule.Key, ShouldEqual, "aid:1")
		So(res.Data.Rule.Mode, ShouldEqual, model.StrMode)
		So(res.Data.Rule.Level, ShouldEqual, 20)
	})

	Convey("Key log", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")

		var res struct {
			Code int         `json:"code"`
			Data []model.Log `json:"data"`
		}

		err := client.Get(ctx, _filterKeyLog, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeEmpty)
	})

	Convey("Key edit", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("area", "danmu")
		params.Set("rule", "unit test")
		params.Set("name", "2333")
		params.Set("fid", fmt.Sprintf("%d", fid))
		params.Set("mode", fmt.Sprintf("%d", model.RegMode))
		params.Set("level", "30")
		params.Set("stime", fmt.Sprintf("%d", time.Now().Unix()))
		params.Set("etime", fmt.Sprintf("%d", time.Now().Add(time.Minute).Unix()))
		params.Set("adid", "2333")
		params.Set("comment", "unit test")
		params.Set("reason", "unit test edit")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterKeyEdit, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Key editinfo", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("fid", fmt.Sprintf("%d", fid))

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rule model.KeyInfo `json:"rule"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filterKeyEditInfo, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data.Rule, ShouldNotBeNil)
		So(res.Data.Rule.ID, ShouldEqual, fid)
		So(res.Data.Rule.Key, ShouldEqual, "aid:1")
		So(res.Data.Rule.Mode, ShouldEqual, model.RegMode)
		So(res.Data.Rule.Level, ShouldEqual, 30)
	})

	Convey("Key del", t, func() {
		params := url.Values{}
		params.Set("key", "aid:2")
		params.Set("name", "2333")
		params.Set("fid", fmt.Sprintf("%d", fid))
		params.Set("adid", "2333")
		params.Set("comment", "unit test")
		params.Set("reason", "unit test del")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _filterKeyDel, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})
}

func TestWhite(t *testing.T) {
	var wid int64
	Convey("White add", t, func() {
		params := url.Values{}
		params.Set("filter", "unit test")
		params.Set("mode", fmt.Sprintf("%d", model.StrMode))
		params.Set("area", "reply")
		params.Set("tpid", "0")
		params.Set("adid", "2333")
		params.Set("name", "2333")
		params.Set("comment", "unit test")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _whiteAdd, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("White search", t, func() {
		params := url.Values{}
		params.Set("filter", "unit test")
		params.Set("area", "reply")
		params.Set("pn", "1")
		params.Set("ps", "20")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Rules []model.WhiteInfo `json:"rules"`
				Total int               `json:"total"`
				PN    int               `json:"pn"`
				PS    int               `json:"ps"`
			} `json:"data"`
		}

		err := client.Get(ctx, _whiteSearch, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data.Total, ShouldEqual, 1)
		So(res.Data.Rules, ShouldHaveLength, 1)
		So(res.Data.Rules[0].ID, ShouldBeGreaterThan, 0)
		So(res.Data.Rules[0].Content, ShouldEqual, "unit test")
		So(res.Data.Rules[0].Mode, ShouldEqual, model.StrMode)
		wid = res.Data.Rules[0].ID
	})

	Convey("White editinfo", t, func() {
		params := url.Values{}
		params.Set("filter_id", fmt.Sprintf("%d", wid))

		var res struct {
			Code int             `json:"code"`
			Data model.WhiteInfo `json:"data"`
		}

		err := client.Get(ctx, _whiteEditInfo, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data.ID, ShouldEqual, wid)
		So(res.Data.Content, ShouldEqual, "unit test")
		So(res.Data.Mode, ShouldEqual, model.StrMode)
	})

	Convey("White log", t, func() {
		params := url.Values{}
		params.Set("filter_id", fmt.Sprintf("%d", wid))

		var res struct {
			Code int         `json:"code"`
			Data []model.Log `json:"data"`
		}

		err := client.Get(ctx, _whiteLog, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeEmpty)
	})

	Convey("White edit", t, func() {
		params := url.Values{}
		params.Set("filter", "unit test")
		params.Set("area", "article")
		params.Set("tpid", "0")
		params.Set("mode", fmt.Sprintf("%d", model.RegMode))
		params.Set("reason", "unit test edit")
		params.Set("comment", "unit test")
		params.Set("adid", "2333")
		params.Set("name", "2333")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _whiteEdit, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("White editinfo", t, func() {
		params := url.Values{}
		params.Set("filter_id", fmt.Sprintf("%d", wid))

		var res struct {
			Code int             `json:"code"`
			Data model.WhiteInfo `json:"data"`
		}

		err := client.Get(ctx, _whiteEditInfo, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data.ID, ShouldEqual, wid)
		So(res.Data.Content, ShouldEqual, "unit test")
		So(res.Data.Areas, ShouldHaveLength, 1)
		So(res.Data.Areas[0], ShouldEqual, "article")
		So(res.Data.Mode, ShouldEqual, model.RegMode)
	})
	Convey("White del", t, func() {
		params := url.Values{}
		params.Set("filter_id", fmt.Sprintf("%d", wid))
		params.Set("adid", "2333")
		params.Set("name", "2333")
		params.Set("reason", "unit test del")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Post(ctx, _whiteDel, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})
}
