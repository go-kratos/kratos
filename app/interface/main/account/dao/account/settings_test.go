package account

import (
	"context"
	"go-common/library/ecode"
	"net/url"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestAccountUpdateBirthday(t *testing.T) {
	var (
		c        = context.Background()
		mid      = int64(2345)
		ip       = ""
		birthday = "2018-01-01"
	)
	convey.Convey("UpdateBirthday", t, func(ctx convey.C) {
		err := d.UpdateBirthday(c, mid, ip, birthday)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

// func TestAccountUpdateSex(t *testing.T) {
// 	var (
// 		c   = context.Background()
// 		mid = int64(0)
// 		sex = int64(0)
// 		ip  = ""
// 	)
// 	convey.Convey("UpdateSex", t, func(ctx convey.C) {
// 		err := d.UpdateSex(c, mid, sex, ip)
// 		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 			ctx.So(err, convey.ShouldBeNil)
// 		})
// 	})
// }

// func TestAccountUpdatePerson(t *testing.T) {
// 	var (
// 		c          = context.Background()
// 		mid        = int64(2345)
// 		birthday   = "2018-01-01"
// 		ip         = ""
// 		datingType = int64(1)
// 		marital    = int64(2)
// 		place      = int64(3)
// 	)
// 	convey.Convey("UpdatePerson", t, func(ctx convey.C) {
// 		err := d.UpdatePerson(c, mid, birthday, ip, datingType, marital, place)
// 		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 			ctx.So(err, convey.ShouldBeNil)
// 		})
// 	})
// }

// func TestAccountUpdateUname(t *testing.T) {
// 	var (
// 		c            = context.Background()
// 		mid          = int64(2345)
// 		ip           = ""
// 		uname        = "dangerous" + strconv.FormatInt(rand.Intn(10))
// 		isUpNickFree = true
// 	)
// 	convey.Convey("UpdateUname", t, func(ctx convey.C) {
// 		err := d.UpdateUname(c, mid, ip, uname, isUpNickFree)
// 		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
// 			ctx.So(err, convey.ShouldBeNil)
// 		})
// 	})
// }

// func TestAccountUpdateFace(t *testing.T) {
// 	var (
// 		c         = context.Background()
// 		mid       = int64(0)
// 		face      = ""
// 		ip        = ""
// 		cookie    = ""
// 		accessKey = ""
// 	)
// 	convey.Convey("UpdateFace", t, func(ctx convey.C) {
// 		faceURL, err := d.UpdateFace(c, mid, face, ip, cookie, accessKey)
// 		ctx.Convey("Then err should be nil.faceURL should not be nil.", func(ctx convey.C) {
// 			ctx.So(err, convey.ShouldBeNil)
// 			ctx.So(faceURL, convey.ShouldNotBeNil)
// 		})
// 	})
// }

func TestAccountParseJavaCode(t *testing.T) {
	var (
		code = int(-617)
	)
	convey.Convey("ParseJavaCode", t, func(ctx convey.C) {
		err := ParseJavaCode(code)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldEqual, ecode.UpdateUnameHadLocked)
		})
	})
}

func TestAccountsign(t *testing.T) {
	var (
		params = &url.Values{}
	)
	convey.Convey("sign", t, func(ctx convey.C) {
		query := d.sign(params)
		ctx.Convey("Then query should not be nil.", func(ctx convey.C) {
			ctx.So(query, convey.ShouldNotBeNil)
		})
	})
}

func TestAccountIdentifyInfo(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2345)
		ip  = ""
	)
	convey.Convey("IdentifyInfo", t, func(ctx convey.C) {
		idt, err := d.IdentifyInfo(c, mid, ip)
		ctx.Convey("Then err should be nil.idt should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(idt, convey.ShouldNotBeNil)
		})
	})
}
