package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoPgcCards(t *testing.T) {
	var (
		ctx         = context.Background()
		ids         = "113,296"
		mockStr     = `{"code":0,"message":"success","result":{"113":{"badge":"","badge_type":0,"cover":"http://i0.hdslb.com/bfs/bangumi/1c3da58352df66fd409831f7d5e13594c804144e.jpg","is_finish":1,"is_play":1,"is_started":1,"media_id":113,"new_ep":{"cover":"http://i0.hdslb.com/bfs/bangumi/88c8d5f6f149880fa31f6d039ed96f58e7a150f1.jpg","id":113506,"index_show":"全13话"},"rating":{"count":1141,"score":9.5},"rights":{"allow_review":1},"season_id":113,"season_title":"TV","season_type":1,"season_type_name":"番剧","stat":{"danmaku":298320,"follow":277453,"view":6978766},"title":"我们大家的河合庄","total_count":13},"296":{"badge":"","badge_type":0,"cover":"http://i0.hdslb.com/bfs/bangumi/37e22d2feafdf9ea5ad2b39860bd0205fb5a2d1d.png","is_finish":1,"is_play":1,"is_started":1,"media_id":296,"new_ep":{"cover":"http://i0.hdslb.com/bfs/bangumi/ae475384c527366a6ef07b414e1e0364695c2aa8.jpg","id":27915,"index_show":"全13话"},"rating":{"count":643,"score":9.2},"rights":{"allow_review":1},"season_id":296,"season_title":"TV","season_type":1,"season_type_name":"番剧","stat":{"danmaku":342835,"follow":253085,"view":5879216},"title":"天体的秩序","total_count":13}}}`
		httpErrStr  = `{"code":-500}`
		EmptyResStr = `{"code":0}`
	)
	convey.Convey("PgcCards", t, func(c convey.C) {
		c.Convey("Normal Situation, Then err should be nil.result should not be nil.", func(c convey.C) {
			httpMock("GET", d.conf.Host.APINewindex).Reply(200).JSON(mockStr)
			result, err := d.PgcCards(ctx, ids)
			c.So(err, convey.ShouldBeNil)
			c.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("Http code Err Situation, Err can't be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APINewindex).Reply(200).JSON(httpErrStr)
			_, err := d.PgcCards(ctx, ids)
			cx.So(err, convey.ShouldNotBeNil)
			fmt.Println(err)
		})
		c.Convey("Empty Res Situation, Then err Should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APINewindex).Reply(200).JSON(EmptyResStr)
			_, err := d.PgcCards(ctx, ids)
			cx.So(err, convey.ShouldNotBeNil)
			fmt.Println(err)
		})
	})
}
