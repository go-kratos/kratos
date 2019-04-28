package dao

import (
	"context"
	"fmt"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoChannelData(t *testing.T) {
	var (
		ctx        = context.Background()
		seasonType = int(0)
		appInfo    = d.conf.TVApp
	)
	convey.Convey("ChannelData", t, func(c convey.C) {
		c.Convey("Then err should be nil.result should not be nil.", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIZone).Reply(200).
				JSON(`{"code":0,"message":"success","result":[{"cover":"http://i0.hdslb.com/bfs/bangumi/75c7528cbf3254dd20a4512376ced74733ab98ef.jpg","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/8beff40b27076475353d25ca590d0932791f71d6.jpg","id":253907,"index":"中文","index_show":"全1话"},"season_id":25944,"title":"黑子的篮球 LAST GAME"},{"cover":"http://i0.hdslb.com/bfs/bangumi/3622ff139d6875c7b196edd3f9db7d0a3883a158.jpg","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/86ec94eb1e00785d15612b700e95eb309710597a.jpg","id":151657,"index":"序列之争","index_show":"全1话"},"season_id":12364,"title":"刀剑神域：序列之争"}]}`)
			result, err := d.ChannelData(ctx, seasonType, appInfo)
			c.So(err, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("We mock http error, The error should not be nil", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIZone).Reply(200).
				JSON(`{"code":-100,"message":"fail","result":[]}`)
			result, err := d.ChannelData(ctx, seasonType, appInfo)
			c.So(err, convey.ShouldNotBeNil)
			c.So(len(result), convey.ShouldBeZeroValue)
		})
		c.Convey("result empty error", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIZone).Reply(200).
				JSON(`{"code":0,"message":"succ","result":[]}`)
			result, err := d.ChannelData(ctx, seasonType, appInfo)
			c.So(err, convey.ShouldEqual, ecode.TvPGCRankEmpty)
			c.So(len(result), convey.ShouldBeZeroValue)
		})
		c.Convey("result not empty but NewEP empty error", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIZone).Reply(200).
				JSON(`{"code":0,"message":"succ","result":[{"cover":"ttt"}]}`)
			result, err := d.ChannelData(ctx, seasonType, appInfo)
			c.So(err, convey.ShouldEqual, ecode.TvPGCRankNewEPNil)
			c.So(len(result), convey.ShouldBeZeroValue)
		})
	})
}

func TestDaoUgcAIData(t *testing.T) {
	var (
		ctx         = context.Background()
		tid         = int16(3)
		normalStr   = `{"note":"统计所有投稿在 2018年11月02日 - 2018年11月09日 的数据综合得分，每日更新一次","source_date":"2018-11-09","code":0,"num":250,"list":[{"aid":35129595,"mid":423442,"pts":571043,"play":277133,"coins":22257,"video_review":2949},{"aid":34994701,"mid":284120,"pts":272817,"play":103881,"coins":7543,"video_review":446},{"aid":35395331,"mid":314791153,"pts":270235,"play":160364,"coins":13475,"video_review":2586}]}`
		httpCodeErr = `{"code":-100}`
		url         = fmt.Sprintf(d.conf.Host.AIUgcType, tid)
	)
	convey.Convey("UgcAIData", t, func(c convey.C) {
		c.Convey("Then err should be nil.result should not be nil.", func(c convey.C) {
			httpMock("GET", url).Reply(200).JSON(normalStr)
			result, err := d.UgcAIData(ctx, tid)
			c.So(err, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("Http code error", func(c convey.C) {
			httpMock("GET", url).Reply(200).JSON(httpCodeErr)
			_, err := d.UgcAIData(ctx, tid)
			c.So(err, convey.ShouldNotBeNil)
		})
	})
}
