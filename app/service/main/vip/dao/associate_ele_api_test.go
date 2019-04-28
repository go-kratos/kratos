package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/vip/model"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

// func TestDaoEleOauthGenerateAccessToken(t *testing.T) {
// 	convey.Convey("TestDaoEleOauthGenerateAccessToken", t, func(ctx convey.C) {
// 		defer gock.OffAll()
// 		httpMock("POST", _oauthGenerateAccessTokenURI).Reply(200).JSON(`{"code":0}`)
// 		_, err := d.EleOauthGenerateAccessToken(context.Background(), &model.ArgEleAccessToken{
// 			AuthCode: "abc"})
// 		ctx.So(err, convey.ShouldBeNil)
// 	})
// }

func TestDaoEleUnionReceivePrizes(t *testing.T) {
	convey.Convey("TestDaoEleUnionReceivePrizes", t, func(ctx convey.C) {
		_, err := d.EleUnionReceivePrizes(context.Background(), &model.ArgEleReceivePrizes{
			ElemeOpenID: "o8f999ad5d724b4a2ljbp7cm",
			BliOpenID:   "e11303e8c26268a6cbdc2dc7fce55199",
			SourceID:    "1",
		})
		ctx.So(err, convey.ShouldEqual, ecode.VipEleUnionReqErr)
	})
}

func TestDaoEleUnionUpdateOpenID(t *testing.T) {
	convey.Convey("TestDaoEleUnionUpdateOpenID", t, func(ctx convey.C) {
		_, err := d.EleUnionUpdateOpenID(context.Background(), &model.ArgEleUnionUpdateOpenID{
			ElemeOpenID: "oe1822de904ceed07bpvCUOU",
			BliOpenID:   "bdca8b71e7a6726885d40a395bf9ccd1",
		})
		ctx.So(err, convey.ShouldEqual, ecode.VipEleUnionReqErr)
	})
}

// func TestDaoEleCanPurchase(t *testing.T) {
// 	convey.Convey("TestDaoEleCanPurchase", t, func(ctx convey.C) {
// 		_, err := d.EleCanPurchase(context.Background(), &model.ArgEleCanPurchase{
// 			ElemeOpenID: "oe1822de904ceed07bpvCUOU",
// 			BliOpenID:   "8551e7e9827f759d7fe13bc41b8a7613",
// 			UserIP:      "10.23.167.25",
// 			VipType:     2,
// 		})
// 		ctx.So(err, convey.ShouldBeNil)
// 	})
// }

func TestDaoEleUnionMobile(t *testing.T) {
	convey.Convey("TestDaoEleUnionMobile", t, func(ctx convey.C) {
		_, err := d.EleUnionMobile(context.Background(), &model.ArgEleUnionMobile{
			ElemeOpenID: "o8f999ad5d724b4a2ljbp7cm",
			BliOpenID:   "e11303e8c26268a6cbdc2dc7fce55199",
		})
		ctx.So(err, convey.ShouldEqual, ecode.VipEleUnionReqErr)
	})
}

// func TestDaoEleBindUnion(t *testing.T) {
// 	convey.Convey("TestDaoEleBindUnion", t, func(ctx convey.C) {
// 		res, err := d.EleBindUnion(context.Background(), &model.ArgEleBindUnion{
// 			ElemeOpenID: "o8f999ad5d724b4a2ljbp7cm",
// 			BliOpenID:   "e11303e8c26268a6cbdc2dc7fce55199",
// 			VipType:     2,
// 			SourceID:    "123456789",
// 			UserIP:      "121.46.231.66",
// 		})
// 		ctx.Convey("TestDaoEleBindUnion ", func(ctx convey.C) {
// 			ctx.So(err, convey.ShouldBeNil)
// 			ctx.So(res, convey.ShouldNotBeNil)
// 		})
// 	})
// }

func TestDaoEleRedPackages(t *testing.T) {
	convey.Convey("TestDaoEleRedPackages", t, func(ctx convey.C) {
		_, err := d.EleRedPackages(context.Background())
		ctx.Convey("TestDaoEleRedPackages ", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.VipEleUnionReqErr)
		})
	})
}

func TestDaoEleSpecailFoods(t *testing.T) {
	convey.Convey("TestDaoEleSpecailFoods", t, func(ctx convey.C) {
		_, err := d.EleSpecailFoods(context.Background())
		ctx.Convey("TestDaoEleSpecailFoods ", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.VipEleUnionReqErr)
		})
	})
}
