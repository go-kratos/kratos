package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoRecomData(t *testing.T) {
	var (
		ctx        = context.Background()
		appInfo    = d.conf.TVApp
		sid        = "20001"
		stype      = "4"
		normalStr  = `{"code":0,"message":"success","result":{"from":1,"list":[{"badge":"","badge_type":0,"cover":"http://i0.hdslb.com/bfs/bangumi/3a315d3dd6223adbd69d3361db9474f2d865aafd.jpg","follow_count":1864,"index_show":"全52集","is_finish":1,"is_started":1,"newest_ep_cover":"http://i0.hdslb.com/bfs/archive/aa02da8dc49f1cdfc8562efc50d455a14194287a.jpg","newest_ep_index":"52","play_count":494484,"season_id":22112,"season_status":2,"season_type":5,"season_type_name":"电视剧","title":"猎场DVD版","total_count":52,"url":"http://www.bilibili.com/bangumi/play/ss22112"}],"season_id":20001,"title":""}}`
		httpErrStr = `{"code":-400}`
		emptyStr   = `{"code":0}`
	)
	convey.Convey("RecomData", t, func(c convey.C) {
		c.Convey("Then err should be nil.result should not be nil.", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIRecom).Reply(200).JSON(normalStr)
			result, err := d.RecomData(ctx, appInfo, sid, stype)
			c.So(err, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("Http Err", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIRecom).Reply(200).JSON(httpErrStr)
			_, err := d.RecomData(ctx, appInfo, sid, stype)
			c.So(err, convey.ShouldNotBeNil)
			fmt.Println(err)
		})
		c.Convey("Empty Result Err", func(c convey.C) {
			httpMock("GET", d.conf.Host.APIRecom).Reply(200).JSON(emptyStr)
			_, err := d.RecomData(ctx, appInfo, sid, stype)
			c.So(err, convey.ShouldNotBeNil)
			fmt.Println(err)
		})
	})
}
