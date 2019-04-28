package http

import (
	"context"
	"flag"
	"net/url"
	"os"
	"testing"

	"go-common/app/service/main/filter/conf"
	"go-common/app/service/main/filter/service"
	xhttp "go-common/library/net/http/blademaster"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	domain       = "localhost:6272"
	_filter      = "http://" + domain + "/x/internal/filter"
	_mfilter     = "http://" + domain + "/x/internal/filter/multi"
	_postFilter  = "http://" + domain + "/x/internal/filter/post"
	_postMFilter = "http://" + domain + "/x/internal/filter/mpost"
	_areaMFilter = "http://" + domain + "/x/internal/filter/area/mpost"
	_article     = "http://" + domain + "/x/internal/filter/article"
	_hit         = "http://" + domain + "/x/internal/filter/v2/hit"
	_test        = "http://" + domain + "/x/internal/filter/test"

	_keyFilter = "http://" + domain + "/x/internal/filter/key"
	_dmFilter  = "http://" + domain + "/x/internal/filter/key/dm"
	_dmTest    = "http://" + domain + "/x/internal/filter/key/test"

	// _rubbish     = "http://" + domain + "/x/internal/filter/rubbish/"
	// _postRubbish = "http://" + domain + "/x/internal/filter/rubbish/post"
)

var (
	client *xhttp.Client
	ctx    = context.TODO()
)

func TestMain(m *testing.M) {
	flag.Set("conf", "../cmd/filter-service-test.toml")
	var err error
	if err = conf.Init(); err != nil {
		panic(err)
	}
	client = xhttp.NewClient(&xhttp.ClientConfig{})

	srv := service.New()
	Init(srv)

	os.Exit(m.Run())
}

func TestFilter(t *testing.T) {
	Convey("Filter", t, func() {
		params := url.Values{}
		params.Set("area", "danmu")
		params.Set("msg", "我爱习大")

		var res struct {
			Code int `json:"code"`
			Data struct {
				MSG    string `json:"msg"`
				Level  int    `json:"level"`
				TypeID []int  `json:"typeid"`
			} `json:"data"`
		}

		err := client.Get(ctx, _filter, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
	})

	Convey("FilterPost", t, func() {
		params := url.Values{}
		params.Set("area", "danmu")
		params.Set("msg", "我爱习大")

		var res struct {
			Code int `json:"code"`
			Data struct {
				MSG    string `json:"msg"`
				Level  int    `json:"level"`
				TypeID []int  `json:"typeid"`
			} `json:"data"`
		}

		err := client.Post(ctx, _postFilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
	})

	Convey("MFilter", t, func() {
		params := url.Values{}
		params.Set("msg", "123")
		params.Add("msg", "123")
		params.Set("area", "article")

		var res struct {
			Code int `json:"code"`
			Data []struct {
				MSG    string `json:"msg"`
				Level  int    `json:"level"`
				TypeID []int  `json:"typeid"`
				Limit  int    `json:"limit"`
			} `json:"data"`
		}

		err := client.Get(ctx, _mfilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data, ShouldHaveLength, 2)
	})

	Convey("PostMFilter", t, func() {
		params := url.Values{}
		params.Set("msg", "123")
		params.Add("msg", "456")
		params.Set("area", "article")

		var res struct {
			Code int `json:"code"`
			Data []struct {
				MSG    string `json:"msg"`
				Level  int    `json:"level"`
				TypeID []int  `json:"typeid"`
				Limit  int    `json:"limit"`
			} `json:"data"`
		}

		err := client.Post(ctx, _postMFilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldNotBeNil)
		So(res.Data, ShouldHaveLength, 2)
	})

	Convey("AreaMFilter", t, func() {
		params := url.Values{}
		params.Set("msg", "123")
		params.Add("msg", "456")
		params.Set("area", "article")

		var res struct {
			Code int `json:"code"`
			Data []struct {
				MSG    string `json:"msg"`
				Level  int    `json:"level"`
				TypeID []int  `json:"typeid"`
				Limit  int    `json:"limit"`
			} `json:"data"`
		}

		err := client.Post(ctx, _areaMFilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
		So(res.Data, ShouldHaveLength, 2)
	})

	Convey("Article", t, func() {
		params := url.Values{}
		params.Set("msg", "123")
		params.Set("area", "article")

		var res struct {
			Code int      `json:"code"`
			Data []string `json:"data"`
		}

		err := client.Post(ctx, _article, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Hit", t, func() {
		params := url.Values{}
		params.Set("msg", "test")
		params.Set("area", "article")

		var res struct {
			Code int      `json:"code"`
			Data []string `json:"data"`
		}

		err := client.Post(ctx, _hit, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Filter test", t, func() {
		params := url.Values{}
		params.Set("msg", "test")
		params.Set("area", "article")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Whites interface{} `json:"whits"`
				Hits   interface{} `json:"hits"`
			} `json:"data"`
		}

		err := client.Get(ctx, _test, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})
}

func TestDM(t *testing.T) {
	Convey("Key filter", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("msg", "test")
		params.Set("area", "danmu")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Out   string `json:"out"`
				Level int    `json:"level"`
			} `json:"data"`
		}

		err := client.Get(ctx, _keyFilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Dm filter", t, func() {
		params := url.Values{}
		params.Set("cid", "1")
		params.Set("aid", "1")
		params.Set("sid", "1")
		params.Set("rid", "1")
		params.Set("pid", "1")
		params.Set("msg", "test")
		params.Set("area", "danmu")

		var res struct {
			Code int `json:"code"`
			Data struct {
				Out   string `json:"out"`
				Level int    `json:"level"`
			} `json:"data"`
		}

		err := client.Get(ctx, _dmFilter, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})

	Convey("Dm test", t, func() {
		params := url.Values{}
		params.Set("key", "aid:1")
		params.Set("msg", "test")
		params.Set("area", "danmu")

		var res struct {
			Code int `json:"code"`
		}

		err := client.Get(ctx, _dmTest, "", params, &res)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		So(res.Code, ShouldEqual, 0)
	})
}
