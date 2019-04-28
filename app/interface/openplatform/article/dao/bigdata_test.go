package dao

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_SkyHorse(t *testing.T) {
	Convey("normal should get data", t, func() {
		data := `{
			"code": 0,
			"data": [
			{
			"tid": 1652,
			"id": 1,
			"goto": "av",
			"source": "user_group",
			"image_cnt" : 3,
			"av_feature": "a"
			},
			{
			"tid": 8227,
			"id": 2,
			"goto": "av",
			"source": "user_group",
			"av_feature": "b"
			}
			],
			"user_feature": "c"
			}`
		httpMock("GET", d.c.Article.SkyHorseURL).Reply(200).JSON(data)
		res, err := d.SkyHorse(ctx(), 1, 0, "", 1, 20)
		So(err, ShouldBeNil)
		So(res.Data, ShouldNotBeEmpty)
	})
	Convey("-3 should get data", t, func() {
		data := `{
			"code": -3,
			"data": [
			{
			"tid": 1652,
			"id": 1,
			"goto": "av",
			"source": "user_group",
			"image_cnt" : 3,
			"av_feature": "a"
			},
			{
			"tid": 8227,
			"id": 2,
			"goto": "av",
			"source": "user_group",
			"av_feature": "b"
			}
			],
			"user_feature": "c"
			}`
		httpMock("GET", d.c.Article.SkyHorseURL).Reply(200).JSON(data)
		res, err := d.SkyHorse(ctx(), 1, 0, "", 1, 20)
		So(err, ShouldBeNil)
		So(res.Data, ShouldNotBeEmpty)
	})
	Convey("code !=0 or -3 should get error", t, func() {
		data := `{"code":-10}`
		httpMock("GET", d.c.Article.SkyHorseURL).Reply(200).JSON(data)
		res, err := d.SkyHorse(ctx(), 1, 0, "", 1, 20)
		So(err, ShouldNotBeNil)
		So(res.Data, ShouldBeEmpty)
	})
}
