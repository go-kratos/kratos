package service

import (
	"context"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/live/wallet/conf"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
	"math/rand"
	"os"
	"strconv"
)

var (
	once              sync.Once
	s                 *Service
	ctx               = context.Background()
	r                 *rand.Rand
	testValidPlatform = []string{"pc", "android", "ios", "h5"}
)

func initConf() {
	if err := conf.Init(); err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	conf.ConfPath = "../cmd/live-wallet-test.toml"
	once.Do(startService)
	os.Exit(m.Run())
}

func startService() {
	initConf()
	s = New(conf.Conf)
	time.Sleep(time.Millisecond * 100)
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func getTestRandUid() int64 {
	return r.Int63n(10000000)
}

func testWith(f func()) func() {
	once.Do(startService)
	return f
}

func getTestBasicParam(tid string, area int64, bizCode string, source string, bizSource string, metaData string) *model.BasicParam {
	return &model.BasicParam{
		TransactionId: tid,
		Area:          area,
		BizCode:       bizCode,
		Source:        source,
		BizSource:     bizSource,
		MetaData:      metaData,
	}
}

func getTestDefaultBasicParam(tid string) *model.BasicParam {
	return getTestBasicParam(tid, 0, "", "", "", "")
}

func atoiForTest(s string) int {
	v, _ := strconv.Atoi(s)
	return v
}

func getTestWallet(t *testing.T, uid int64, platform string) *model.MelonseedResp {
	v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
	if err != nil {
		t.FailNow()
	}
	return v.(*model.MelonseedResp)

}

func getTestAll(t *testing.T, uid int64, platform string) *model.DetailResp {
	v, err := s.GetAll(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
	if err != nil {
		t.FailNow()
	}
	return v.(*model.DetailResp)

}

func TestService_Get(t *testing.T) {
	Convey("normal get without metal", t, testWith(func() {
		Convey("basic", func() {
			for _, platform := range testValidPlatform {
				uid := getTestRandUid()
				v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
				resp := v.(*model.MelonseedResp)
				So(resp, ShouldNotBeNil)
				So(err, ShouldBeNil)
				So(atoiForTest(resp.Gold), ShouldBeGreaterThanOrEqualTo, 0)
				So(atoiForTest(resp.Silver), ShouldBeGreaterThanOrEqualTo, 0)
			}
		})
	}))

	Convey("params", t, testWith(func() {
		Convey("invalid platform", func() {
			platform := "unknown"
			uid := getTestRandUid()
			v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
			So(v, ShouldBeNil)
			So(err, ShouldEqual, ecode.RequestErr)
		})
		Convey("invalid uid", func() {
			platform := testValidPlatform[0]
			var uid int64 = 0
			v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 0)
			So(v, ShouldBeNil)
			So(err, ShouldEqual, ecode.RequestErr)
		})
	}))

	Convey("metal", t, testWith(func() {
		uid := getTestRandUid()
		platform := testValidPlatform[0]
		v, err := s.Get(ctx, getTestDefaultBasicParam(""), uid, platform, 1)
		resp := v.(*model.MelonseedWithMetalResp)
		So(resp, ShouldNotBeNil)
		So(err, ShouldBeNil)
		So(atoiForTest(resp.Gold), ShouldBeGreaterThanOrEqualTo, 0)
		So(atoiForTest(resp.Silver), ShouldBeGreaterThanOrEqualTo, 0)
		So(atoiForTest(resp.Metal), ShouldBeGreaterThanOrEqualTo, 0)

	}))
}

func TestService_GetAll(t *testing.T) {
	Convey("normal getAll without metal", t, testWith(func() {
		Convey("basic", func() {
			for _, platform := range testValidPlatform {
				uid := getTestRandUid()
				resp := getTestAll(t, uid, platform)
				So(resp, ShouldNotBeNil)
				So(atoiForTest(resp.Gold), ShouldBeGreaterThanOrEqualTo, 0)
				So(atoiForTest(resp.Silver), ShouldBeGreaterThanOrEqualTo, 0)
				So(atoiForTest(resp.SilverPayCnt), ShouldBeGreaterThanOrEqualTo, 0)
				So(atoiForTest(resp.GoldRechargeCnt), ShouldBeGreaterThanOrEqualTo, 0)
				So(atoiForTest(resp.GoldPayCnt), ShouldBeGreaterThanOrEqualTo, 0)
			}
		})
	}))

	Convey("advanced", t, testWithTestUser(func(u *TestUser) {

		resp := getTestAll(t, u.uid, "pc")
		So(resp.CostBase, ShouldEqual, 1000)

		resp = getTestAll(t, u.uid, "pc")
		So(resp.CostBase, ShouldEqual, 1000)

		v, err := s.GetAll(ctx, getTestDefaultBasicParam(""), u.uid, "pc", 1)
		if err != nil {
			t.FailNow()
		}
		respWithMetal := v.(*model.DetailWithMetalResp)
		So(respWithMetal, ShouldNotBeNil)

		So(respWithMetal.CostBase, ShouldEqual, 1000)
	}))
}
