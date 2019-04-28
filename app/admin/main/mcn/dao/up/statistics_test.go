package up

import (
	"context"
	"testing"

	"go-common/app/admin/main/mcn/model"
	xtime "go-common/library/time"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestUpMcnDataOverview(t *testing.T) {
	convey.Convey("McnDataOverview", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = xtime.Time(1542124800)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			m, err := d.McnDataOverview(c, date)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(m, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpMcnRankFansOverview(t *testing.T) {
	convey.Convey("McnRankFansOverview", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataType = model.DataType(2)
			date     = xtime.Time(1542124800)
			topLen   = int(5)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mrf, mids, err := d.McnRankFansOverview(c, dataType, date, topLen)
			ctx.Convey("Then err should be nil.mrf,mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldBeNil)
				ctx.So(mrf, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnRankArchiveLikesOverview(t *testing.T) {
	convey.Convey("McnRankArchiveLikesOverview", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			dataType = model.DataType(2)
			date     = xtime.Time(1542124800)
			topLen   = int(5)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ras, mids, avids, tids, err := d.McnRankArchiveLikesOverview(c, dataType, date, topLen)
			ctx.Convey("Then err should be nil.ras,mids,avids,tids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tids, convey.ShouldBeNil)
				ctx.So(avids, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldBeNil)
				ctx.So(ras, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpMcnDataTypeSummary(t *testing.T) {
	convey.Convey("McnDataTypeSummary", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			date = xtime.Time(1542124800)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			mmd, tids, err := d.McnDataTypeSummary(c, date)
			ctx.Convey("Then err should be nil.mmd,tids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tids, convey.ShouldBeNil)
				ctx.So(mmd, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpArcTopDataStatistics(t *testing.T) {
	convey.Convey("ArcTopDataStatistics", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnGetRankReq{}
		)
		arg.SignID = 214
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `
			{
				"message":"0",
				"code":0,
				"data":{
					"type_list":[
						{
							"tid":1,
							"name":"视频"
						}
					],
					"result":[
						{
							"data_type":1,
							"likes_increase":13,
							"likes_accumulate":13,
							"play_increase":7,
							"archive_id":10110514,
							"archive_title":"不同清晰度",
							"pic":"http://i1.hdslb.com/bfs/archive/3348cb2cb34423f936916444a0a77e59f9daf1d",
							"tid_name":"日常",
							"tid":21,
							"ctime":1535362150,
							"author":{
								"face":"http://static.hdslb.com/images/member/noface.gif",
								"mid":27515266,
								"name":"Testeew还觉得是发货"
							},
							"stat":{
								"view":0
							}
						}
					]
				},
				"ttl":1
			}`
			httpMock("GET", d.arcTopURL).Reply(200).JSON(result)
			reply, err := d.ArcTopDataStatistics(c, arg)
			ctx.Convey("Then err should be nil.reply should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reply, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDataFans(t *testing.T) {
	convey.Convey("DataFans", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnCommonReq{SignID: 1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `{
				"message":"",
				"code":0,
				"data":{
					"fans_all":0,
					"fans_inc":0,
					"act_fans":0,
					"fans_dec_all":0,
					"fans_dec":0,
					"view_fans_rate":0,
					"act_fans_rate":0,
					"reply_fans_rate":0,
					"danmu_fans_rate":0,
					"coin_fans_rate":0,
					"like_fans_rate":0,
					"fav_fans_rate":0,
					"share_fans_rate":0,
					"live_gift_fans_rate":0,
					"live_danmu_fans_rate":0
				}
			}`
			httpMock("GET", d.dataFansURL).Reply(200).JSON(result)
			reply, err := d.DataFans(c, arg)
			ctx.Convey("Then err should be nil.reply should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reply, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDataFansBaseAttr(t *testing.T) {
	convey.Convey("DataFansBaseAttr", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnCommonReq{SignID: 1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `{
				"message":"",
				"code":0,
				"data":{
					"fans_sex":{
						"male":0,
						"female":0
					},
					"fans_age":{
						"a":0,
						"b":0,
						"c":0,
						"d":0
					},
					"fans_play_way":{
						"app":0,
						"pc":0,
						"outside":0,
						"other":0
					}
				}
			}`
			httpMock("GET", d.dataFansBaseAttrURL).Reply(200).JSON(result)
			sex, age, playWay, err := d.DataFansBaseAttr(c, arg)
			ctx.Convey("Then err should be nil.sex,age,playWay should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(playWay, convey.ShouldNotBeNil)
				ctx.So(age, convey.ShouldNotBeNil)
				ctx.So(sex, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDataFansArea(t *testing.T) {
	convey.Convey("DataFansArea", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnCommonReq{SignID: 1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `{
				"message":"",
				"code":0,
				"data":{
					"result":[
						{
							"province":"",
							"user":0
						}
					]
				}
			}`
			httpMock("GET", d.dataFansAreaURL).Reply(200).JSON(result)
			reply, err := d.DataFansArea(c, arg)
			ctx.Convey("Then err should be nil.reply should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reply, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDataFansType(t *testing.T) {
	convey.Convey("DataFansType", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnCommonReq{SignID: 1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `{
				"message":"",
				"code":0,
				"data":{
					"result":[
						{
							"type_id":0,
							"user":0,
							"type_name":""
						}
					]
				}
			}`
			httpMock("GET", d.dataFansTypeURL).Reply(200).JSON(result)
			reply, err := d.DataFansType(c, arg)
			ctx.Convey("Then err should be nil.reply should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reply, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpDataFansTag(t *testing.T) {
	convey.Convey("DataFansTag", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.McnCommonReq{SignID: 1}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer gock.OffAll()
			result := `{
				"message":"",
				"code":0,
				"data":{
					"result":[
						{
							"tag_id":0,
							"user":0,
							"tag_name":""
						}
					]
				}
			}`
			httpMock("GET", d.dataFansTagURL).Reply(200).JSON(result)
			reply, err := d.DataFansTag(c, arg)
			ctx.Convey("Then err should be nil.reply should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(reply, convey.ShouldNotBeNil)
			})
		})
	})
}
