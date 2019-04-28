package dao

import (
	"context"
	"fmt"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoHeaderData(t *testing.T) {
	var (
		ctx            = context.Background()
		appInfo        = d.conf.TVApp
		normalStr      = `{"code":0,"message":"success","result":{"cn":[{"cover":"http://i0.hdslb.com/bfs/bangumi/1fb653e1e9a825ef1487216c834d3f72c647a8aa.jpg","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/4bbd70afe81fd1c483ca61a038a8a983358ef3ac.jpg","id":254310,"index":"34","index_show":"连载中"},"season_id":24438,"title":"小绿和小蓝"}],"documentary":[{"cover":"http://i0.hdslb.com/bfs/bangumi/30daefc74e09e9cb1a3e8771226002efe186aaf6.jpg","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/efe3a69efe4779662d4eff11f12e4b94cc3e52b6.jpg","id":253992,"index":"花絮12","index_show":"连载中"},"season_id":25810,"title":"历史那些事"}],"jp":[{"cover":"http://i0.hdslb.com/bfs/bangumi/a4c0e0ccc44fe3949a734f546cf5bb07da925bad.png","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/523e3b8faee592c683ea55cfdcdb26106651fdea.jpg","id":250465,"index":"6","index_show":"连载中"},"season_id":25739,"title":"关于我转生变成史莱姆这档事"}],"movie":[{"cover":"http://i0.hdslb.com/bfs/bangumi/75c7528cbf3254dd20a4512376ced74733ab98ef.jpg","new_ep":{"cover":"http://i0.hdslb.com/bfs/archive/8beff40b27076475353d25ca590d0932791f71d6.jpg","id":253907,"index":"中文","index_show":"全1话"},"season_id":25944,"title":"黑子的篮球 LAST GAME"}],"tv":[{"cover":"http://i0.hdslb.com/bfs/bangumi/f79e3d4bb24d73db6407d771a4293d737376a11a.jpg","new_ep":{"cover":"http://i1.hdslb.com/bfs/archive/e24e3c851a81d46d3706d7d5401d63e0f6d9453e.jpg","id":143246,"index ":"30 ","index_show ":"全30话 "},"season_id ":21448,"title ":"亮剑"}]}}`
		httpErrStr     = `{"code": -400}`
		missingZoneStr = `{"code":0,"message":"success","result":{}}`
	)
	convey.Convey("HeaderData Normal Situation", t, func(c convey.C) {
		c.Convey("Then err should be nil.result should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIIndex).Reply(200).JSON(normalStr)
			result, err := d.HeaderData(ctx, appInfo)
			cx.So(err, convey.ShouldBeNil)
			cx.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("Http err, Then err Should not be nil", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIIndex).Reply(200).JSON(httpErrStr)
			_, err := d.HeaderData(ctx, appInfo)
			cx.So(err, convey.ShouldNotBeNil)
			fmt.Println(err)
		})
		c.Convey("Then err should be about the zone", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIIndex).Reply(200).JSON(missingZoneStr)
			_, err := d.HeaderData(ctx, appInfo)
			cx.So(err, convey.ShouldNotBeNil)
			cx.So(err.Error(), convey.ShouldEqual, "Result Miss Data: jp")
			fmt.Println(err)
		})
	})
}

func TestDaoFollowData(t *testing.T) {
	var (
		ctx        = context.Background()
		appInfo    = d.conf.TVApp
		accessKey  = "9dfa21e5d98f0be3410974f6894b72af"
		normalStr  = `{"code":0,"count":"59","message":"success","pages":"6","result":[{"actor":[],"alias":"ラーメン大好き小泉さん","allow_bp":"1","allow_download":"0","area":"日本","bangumi_id":"3905","bangumi_title":"爱吃拉面的小泉同学","brief":"今天，她也在某处享用着拉面——\n冷淡而沉默，不和他人亲近的神秘转学生，小泉同学。她其实是每天追求着美...","copyright":"bilibili","cover":"http://i0.hdslb.com/bfs/bangumi/e4dcc80598a4133af6b1a880bb8006fa5346f31b.jpg","danmaku_count":"436189","ed_jump":5,"episodes":[],"evaluate":"","favorites":"975275","is_finish":"1","is_started":1,"last_time":"2018-03-22 19:00:00.0","limitGroupId":1193,"new_cover":"http://i0.hdslb.com/bfs/archive/a06bbc48bc0c5c95afe0d2875d12376fc764807a.jpg","new_ep":{"av_id":"21083439","cover":"http://i0.hdslb.com/bfs/archive/a06bbc48bc0c5c95afe0d2875d12376fc764807a.jpg","danmaku":"34591825","episode_id":"164993","episode_status":2,"from":"bangumi","index":"12","index_title":"名古屋 / 再会","page":"1","up":{},"update_time":"2018-03-22 19:00:00.0"},"newest_ep_id":"164993","newest_ep_index":"12","play_count":"14835307","progress":"全12话","pub_string":"","pub_time":"2018-01-04 19:00:00","related_seasons":[],"season_id":"21728","season_status":2,"season_title":"TV","seasons":[],"share_url":"http://bangumi.bilibili.com/anime/21728/","spid":"0","squareCover":"http://i0.hdslb.com/bfs/bangumi/a160be35d676316efeb53157001b8df257a48a77.jpg","staff":"","tag2s":[],"tags":[{"cover":"http://i0.hdslb.com/bfs/bangumi/29baa04d18505c775c131e0d0db0bf8704cc61bb.jpg","tag_id":"21","tag_name":"治愈"},{"cover":"http://i0.hdslb.com/bfs/bangumi/0a89e6fc2da1a8f714ef3872d3b58df7137eb195.jpg","tag_id":"106","tag_name":"美食"},{"cover":"http://i0.hdslb.com/bfs/bangumi/3121473d5dd03a9bcccb8490034207e724e731b3.jpg","tag_id":"135","tag_name":"漫画改"}],"title":"爱吃拉面的小泉同学","total_count":"12","trailerAid":"-1","update_pattern":"","user_season":{"attention":"0","bp":0,"last_ep_index":"","last_time":"0","report_ts":0},"watchingCount":"0","weekday":"4"}]}`
		httpErrStr = `{code": -400}`
	)
	convey.Convey("FollowData", t, func(c convey.C) {
		c.Convey("Then err should be nil.result should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIFollow).Reply(200).JSON(normalStr)
			result, err := d.FollowData(ctx, appInfo, accessKey)
			cx.So(err, convey.ShouldBeNil)
			cx.So(result, convey.ShouldNotBeNil)
		})
		c.Convey("Then err should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIFollow).Reply(200).JSON(httpErrStr)
			_, err := d.FollowData(ctx, appInfo, accessKey)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
}
