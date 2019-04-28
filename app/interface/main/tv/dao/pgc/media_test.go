package pgc

import (
	"context"
	"testing"

	"go-common/app/interface/main/tv/model"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoMedia(t *testing.T) {
	var (
		ctx     = context.Background()
		tvParam = &model.MediaParam{
			SeasonID:  296,
			TrackPath: "1",
			MobiAPP:   "android_tv_yst",
			Platform:  "android",
			Build:     101300,
		}
		normalStr  = `{"code":0,"message":"success","result":{"cover":"http://i0.hdslb.com/bfs/bangumi/b2662854abfffc742a505853274425ac33bced24.jpg","dialog":{"btn_right":{"title":"成为大会员","type":"vip"},"desc":"","title":"开通大会员抢先看"},"episodes":[{"aid":10100146,"badge":"会员","badge_type":0,"cid":10114326,"cover":"http://i0.hdslb.com/bfs/archive/0e40e69bdc340e1c773f8813652935145ffd756c.jpg","duration":0,"ep_id":115281,"episode_status":13,"from":"bangumi","index":"大会员专享的","index_title":"大会员专享的","mid":452156,"page":1,"pub_real_time":"2018-08-24 20:27:00","section_id":23781,"section_type":0,"share_url":"https://m.bilibili.com/bangumi/play/ep115281","vid":""},{"aid":10110399,"badge":"会员","badge_type":0,"cid":10131912,"cover":"http://i0.hdslb.com/bfs/archive/0e40e69bdc340e1c773f8813652935145ffd756c.jpg","duration":0,"ep_id":117023,"episode_status":13,"from":"bangumi","index":"付费2","index_title":"付费2作为覆盖点","mid":27515256,"page":1,"pub_real_time":"2018-09-21 08:00:00","section_id":23781,"section_type":0,"share_url":"https://m.bilibili.com/bangumi/play/ep117023","vid":""},{"aid":10110399,"badge":"会员","badge_type":0,"cid":10131912,"cover":"http://i0.hdslb.com/bfs/archive/0e40e69bdc340e1c773f8813652935145ffd756c.jpg","duration":0,"ep_id":117024,"episode_status":13,"from":"bangumi","index":"付费3","index_title":"付费3作为测试覆盖点","mid":27515256,"page":1,"pub_real_time":"2018-09-20 08:00:00","section_id":23781,"section_type":0,"share_url":"https://m.bilibili.com/bangumi/play/ep117024","vid":""},{"aid":10110399,"badge":"会员","badge_type":0,"cid":10131912,"cover":"http://i0.hdslb.com/bfs/archive/0e40e69bdc340e1c773f8813652935145ffd756c.jpg","duration":0,"ep_id":117025,"episode_status":13,"from":"bangumi","index":"付费ep4","index_title":"付费ep4作为测试覆盖点","mid":27515256,"page":1,"pub_real_time":"2018-09-21 12:00:00","section_id":23781,"section_type":0,"share_url":"https://m.bilibili.com/bangumi/play/ep117025","vid":""}],"evaluate":"自动化:大会员专享纯付费ep,自动化测试使用番,请勿修改","is_new_danmaku":1,"link":"http://www.bilibili.com/bangumi/media/md2130684526/","media_id":2130684526,"mid":-1,"mode":2,"music_menus":[{"cover_url":"http://uat-i0.hdslb.com/bfs/static/a8f59108b204df22093dbc15d416a80bd8ea7249.jpg","intro":"","is_off":0,"is_pay":1,"menu_id":59054,"play_num":0,"title":"ff云の泣 - 青玉案（唱：云の泣 词：小鱼萝莉）"}],"newest_ep":{"desc":"连载中, 每周三更新","id":117025,"index":"付费ep4","is_new":0,"pub_real_time":"1537502400000"},"paster":{"aid":0,"allow_jump":0,"cid":0,"duration":0,"type":0,"url":""},"payment":{"pay_tip":{"primary":{"sub_title":"","title":"开通大会员抢先看","type":1,"url":""}},"pay_type":{"allow_discount":0,"allow_pack":0,"allow_ticket":0,"allow_time_limit":0,"allow_vip_discount":0},"price":"0.0","promotion":"","tip":"大会员专享观看特权哦~","vip_first_promotion":"","vip_first_switch":"1","vip_promotion":"开通大会员抢先看"},"player_icon":{"ctime":1540525201,"hash1":"4f688a632ce9dfe167ac3021ab9fb9f3","hash2":"ec06112f8b96d9c8d592593249d47c34","url1":"http://uat-i0.hdslb.com/bfs/archive/29cb6cf28b86c069ce7ef9aa051f84a1cb3d5148.json","url2":"http://uat-i0.hdslb.com/bfs/archive/737f375f10216213e203f807097bce03c211669d.json"},"publish":{"is_finish":0,"is_started":1,"pub_time":"2018-09-10 10:00:00","pub_time_show":"09月10日10:00","weekday":3},"record":"","rights":{"allow_bp":0,"allow_download":0,"allow_review":1,"area_limit":0,"ban_area_show":1,"copyright":"bilibili","is_preview":0},"season_id":20001,"season_status":13,"season_title":"全ep大会员专享","season_type":4,"seasons":[{"is_new":0,"season_id":33729,"season_title":"player专用的","title":"自动化:专门测试player的付费预览类,不要有任何修改,任何时候!"},{"is_new":0,"season_id":33382,"season_title":"免费仅大陆可看加iphone限制","title":"自动化:免费仅大陆可看加iphone限制,自动化测试使用番,请勿修改"},{"is_new":0,"season_id":33405,"season_title":"大会员&付费混杂ep","title":"自动化:大会员&付费混杂ep,自动化测试使用番,请勿修改"},{"is_new":0,"season_id":20001,"season_title":"全ep大会员专享","title":"自动化:大会员专享纯付费ep,自动化测试使用番,请勿修改"},{"is_new":0,"season_id":20009,"season_title":"付费&会员限时承包免费","title":"自动化:付费&会员，承包免费，请勿修改"},{"is_new":0,"season_id":33413,"season_title":"付费抢先限时","title":"专门测试player的国创-付费抢先限时"},{"is_new":0,"season_id":33427,"season_title":"霹雳付费","title":"自动化:霹雳付费,纯付费ep,自动化测试使用番,请勿修改"},{"is_new":0,"season_id":33406,"season_title":"免费但限制ipad平台","title":"自动化:免费仅亚洲非日本加ipad限制,自动化测试使用番,请勿修改"},{"is_new":0,"season_id":33801,"season_title":"免费仅大陆可看加安装限制","title":"自动化:免费仅大陆可看加android限制,自动化测试使用番,请勿修改"}],"series_id":3822,"share_url":"http://www.bilibili.com/bangumi/play/ss20001","square_cover":"http://i0.hdslb.com/bfs/bangumi/f163df1ab57e3310709941648572625ab6ce8db0.jpg","stat":{"coins":0,"danmakus":0,"favorites":0,"reply":0,"share":0,"views":0},"title":"自动化:大会员专享纯付费ep,自动化测试使用番,请勿修改","total_ep":23}}`
		httpErrStr = `{"code":-400}`
	)
	convey.Convey("Media", t, func(cx convey.C) {
		cx.Convey("Then err should be nil.detail should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIMedia).Reply(200).JSON(normalStr)
			detail, err := d.Media(ctx, tvParam)
			cx.So(err, convey.ShouldBeNil)
			cx.So(detail, convey.ShouldNotBeNil)
		})
		cx.Convey("Http Err, err should not be nil", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIMedia).Reply(200).JSON(httpErrStr)
			_, err := d.Media(ctx, tvParam)
			cx.So(err, convey.ShouldNotBeNil)
			cx.So(err, convey.ShouldEqual, ecode.Int(-400))
			httpMock("GET", d.conf.Host.APIMedia).Reply(-500).JSON(httpErrStr)
			_, err = d.Media(ctx, tvParam)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
	convey.Convey("MediaV2", t, func(cx convey.C) {
		cx.Convey("Then err should be nil.detail should not be nil.", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIMediaV2).Reply(200).JSON(normalStr)
			detail, err := d.MediaV2(ctx, tvParam)
			cx.So(err, convey.ShouldBeNil)
			cx.So(detail, convey.ShouldNotBeNil)
		})
		cx.Convey("Http Err, err should not be nil", func(cx convey.C) {
			httpMock("GET", d.conf.Host.APIMediaV2).Reply(200).JSON(httpErrStr)
			_, err := d.MediaV2(ctx, tvParam)
			cx.So(err, convey.ShouldNotBeNil)
			cx.So(err, convey.ShouldEqual, ecode.Int(-400))
			httpMock("GET", d.conf.Host.APIMediaV2).Reply(-500).JSON(httpErrStr)
			_, err = d.MediaV2(ctx, tvParam)
			cx.So(err, convey.ShouldNotBeNil)
		})
	})
}
