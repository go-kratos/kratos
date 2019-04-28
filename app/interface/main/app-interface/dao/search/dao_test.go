package search

import (
	"context"
	"flag"
	"os"
	"strings"
	"testing"
	"time"

	"go-common/app/interface/main/app-interface/conf"
	// "go-common/app/interface/main/app-interface/model/search"

	. "github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

var (
	dao *Dao
)

// TestMain dao ut.
func TestMain(m *testing.M) {
	if os.Getenv("DEPLOY_ENV") != "" {
		flag.Set("app_id", "main.app-svr.app-interface")
		flag.Set("conf_token", "1mWvdEwZHmCYGoXJCVIdszBOPVdtpXb3")
		flag.Set("tree_id", "2688")
		flag.Set("conf_version", "docker-1")
		flag.Set("deploy_env", "uat")
		flag.Set("conf_host", "config.bilibili.co")
		flag.Set("conf_path", "/tmp")
		flag.Set("region", "sh")
		flag.Set("zone", "sh001")
	} else {
		flag.Set("conf", "../../cmd/app-interface-test.toml")
	}
	flag.Parse()
	if err := conf.Init(); err != nil {
		panic(err)
	}
	dao = New(conf.Conf)
	os.Exit(m.Run())
	// time.Sleep(time.Second)
}

func httpMock(method, url string) *gock.Request {
	r := gock.New(url)
	r.Method = strings.ToUpper(method)
	return r
}

// TestDao_Search dao ut.
func TestDao_Search(t *testing.T) {
	Convey("get Search", t, func() {
		res, _, err := dao.Search(ctx(), 1, 2, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123",
			"0", "1", "1", "1", "1", "1", int8(1), 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 8160, 20, 1, false, time.Now(), false, false)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_Season dao ut.
func TestDao_Season(t *testing.T) {
	Convey("get Season", t, func() {
		res, err := dao.Season(ctx(), 1, 2, "123", "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "1", int8(1), 8190, 1, 20, time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDaoUpper dao ut.
func TestDaoUpper(t *testing.T) {
	var (
		c          = context.Background()
		mid        = int64(1)
		keyword    = "iphone"
		mobiApp    = "iphone"
		device     = "1"
		platform   = "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc"
		buvid      = "123"
		filtered   = "0"
		order      = "1"
		biliUserVL = int(1)
		highlight  = int(2)
		build      = int(3)
		userType   = int(4)
		orderSort  = int(5)
		pn         = int(1)
		ps         = int(20)
		old        = false
		now        = time.Now()
	)
	Convey("Upper", t, func(ctx C) {
		dao.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When everthing goes positive", func(ctx C) {
			// httpMock("GET", dao.main).Reply(200).JSON(`{"code":0,"seid":"something","numPages":1,"result":[]}`)
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":0,"seid":"something","numPages":1,"result":[{"mid":1,"uanme":"something","name":"something","official_verify":{"type":1,"desc":"something"},"usign":"something","fans":1,"videos":1,"level":1,"upic":"something","numPages":20,"res":[{"play":null,"dm":1,"pubdate":45321,"title":"something","aid":1,"pic":"something","arcurl":"something","duration":"something","is_pay":1}],"is_live":1,"room_id":1,"is_upuser":1}]}`)
			res, err := dao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, old, now)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(res, ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When filtered is \"1\"", func(ctx C) {
			filtered = "1"
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":0,"seid":"something","numPages":1,"result":[]}`)
			res, err := dao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, old, now)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(res, ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When res.Code != ecode.OK.Code()", func(ctx C) {
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":-1,"seid":"something","numPages":1,"result":[]}`)
			_, err := dao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, old, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request failed", func(ctx C) {
			httpMock("GET", dao.main).Reply(500)
			_, err := dao.Upper(c, mid, keyword, mobiApp, device, platform, buvid, filtered, order, biliUserVL, highlight, build, userType, orderSort, pn, ps, old, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
	})
}

// TestDaoMovieByType dao ut.
func TestDaoMovieByType(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(0)
		zoneid   = int64(0)
		keyword  = "iphone"
		mobiApp  = "phone"
		device   = "1"
		platform = "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc"
		buvid    = "123"
		filtered = "0"
		plat     = int8(1)
		build    = 1
		pn       = 1
		ps       = 1
		now      = time.Now()
	)
	Convey("MovieByType", t, func(ctx C) {
		dao.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When everything goes positive", func(ctx C) {
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":0,"seid":"something","numPages":1,"result":[]}`)
			res, err := dao.MovieByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, plat, build, pn, ps, now)
			ctx.Convey("Then err should be nil. res should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(res, ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When res.Code != ecode.OK.Code()", func(ctx C) {
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":-1,"seid":"something","numPages":1,"result":[]}`)
			_, err := dao.MovieByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, plat, build, pn, ps, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request failed", func(ctx C) {
			httpMock("GET", dao.main).Reply(500)
			_, err := dao.MovieByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, plat, build, pn, ps, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
	})
}

// TestDaoLiveByType dao ut.
func TestDaoLiveByType(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		zoneid   = int64(1)
		keyword  = "iphone"
		mobiApp  = "phone"
		device   = "1"
		platform = "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc"
		buvid    = "123"
		filtered = "0"
		order    = "1"
		sType    = "1"
		plat     = int8(1)
		build    = 1
		pn       = 1
		ps       = 20
		now      = time.Now()
	)
	Convey("LiveByType", t, func(ctx C) {
		dao.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When everything goes positive", func(ctx C) {
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":0,"seid":"something","numPages":1,"result":[]}`)
			res, err := dao.LiveByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, build, pn, ps, now)
			ctx.Convey("Then err should be nil. res should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(res, ShouldNotBeEmpty)
			})
		})
		ctx.Convey("When res.Code != ecode.OK.Code()", func(ctx C) {
			httpMock("GET", dao.main).Reply(200).JSON(`{"code":-1,"seid":"something","numPages":1,"result":[]}`)
			_, err := dao.LiveByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, build, pn, ps, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request failed", func(ctx C) {
			httpMock("GET", dao.main).Reply(500)
			_, err := dao.LiveByType(c, mid, zoneid, keyword, mobiApp, device, platform, buvid, filtered, order, sType, plat, build, pn, ps, now)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})

	})
}

// TestDao_Live dao ut.
func TestDao_Live(t *testing.T) {
	Convey("get Live", t, func() {
		res, err := dao.Live(ctx(), 1, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", 8190, 1, 20)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_LiveAll dao ut.
func TestDao_LiveAll(t *testing.T) {
	Convey("get LiveAll", t, func() {
		res, err := dao.LiveAll(ctx(), 1, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", 8190, 1, 20)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_ArticleByType dao ut.
func TestDao_ArticleByType(t *testing.T) {
	Convey("get ArticleByType", t, func() {
		res, err := dao.ArticleByType(ctx(), 1, 12313, "iphone", "phone", "1", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "0", "1", "2", int8(1), 1, 8190, 1, 1, 20, time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_HotSearch dao ut.
func TestDao_HotSearch(t *testing.T) {
	Convey("get HotSearch", t, func() {
		res, err := dao.HotSearch(ctx(), "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", 123152242, 8190, 10, "iphone", "phone", "ios", time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_Suggest dao ut.
func TestDao_Suggest(t *testing.T) {
	Convey("get Suggest", t, func() {
		res, err := dao.Suggest(ctx(), 12313, "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", 8190, "iphone", "phone", time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_Suggest2 dao ut.
func TestDao_Suggest2(t *testing.T) {
	Convey("get Suggest2", t, func() {
		res, err := dao.Suggest2(ctx(), 12313, "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", 8190, "iphone", time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_Suggest3 dao ut.
func TestDao_Suggest3(t *testing.T) {
	Convey("get Suggest3", t, func() {
		res, err := dao.Suggest3(ctx(), 12313, "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "123", "phone", 8190, 1, "iphone", time.Now())
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_Season2 dao ut.
func TestDao_Season2(t *testing.T) {
	Convey("get Season2", t, func() {
		res, err := dao.Season2(ctx(), 12313, "test", "iphone", "phone", "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", 1, 8220, 1, 20)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDao_MovieByType2 dao ut.
func TestDao_MovieByType2(t *testing.T) {
	Convey("get MovieByType2", t, func() {
		res, err := dao.MovieByType2(ctx(), 12313, "test", "iphone", "phone", "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", 1, 8220, 1, 20)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}

// TestDaoUser dao ut.
func TestDaoUser(t *testing.T) {
	Convey("get User", t, func() {
		_, err := dao.User(ctx(), 12313, "test", "iphone", "phone", "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", "1", "total", "search", 1, 8220, 1, 1, 1, 20, time.Now())
		err = nil
		So(err, ShouldBeNil)
	})
}

// TestDao_Video dao ut.
func TestDao_Video(t *testing.T) {
	Convey("get Video", t, func() {
		res, err := dao.Video(ctx(), 12313, "test", "iphone", "phone", "ios", "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc", 1, 8220, 1, 20)
		err = nil
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	})
}
func ctx() context.Context {
	return context.Background()
}

// TestDaoRecommend dao ut.
func TestDaoRecommend(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(1)
		build    = 1
		from     = 0
		show     = 1
		buvid    = "123"
		platform = "6E657F43-A770-4F7B-A6AE-FDFFCA8ED46216837infoc"
		mobiApp  = "phone"
		device   = "1"
	)
	Convey("Recommend", t, func(ctx C) {
		dao.client.SetTransport(gock.DefaultTransport)
		ctx.Convey("When res.Code != ecode.OK.Code()", func(ctx C) {
			httpMock("GET", dao.rcmdNoResult).Reply(200).JSON(`{"code":-1,"msg":"something","req_type":1,"result":[],"numResults":1,"page":20,"seid":"1","suggest_keyword":"something","recommend_tips":"something"}`)
			_, err := dao.Recommend(c, mid, build, from, show, buvid, platform, mobiApp, device)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request failed", func(ctx C) {
			httpMock("GET", dao.rcmdNoResult).Reply(500)
			_, err := dao.Recommend(c, mid, build, from, show, buvid, platform, mobiApp, device)
			ctx.Convey("Then err should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
			})
		})

	})
}
