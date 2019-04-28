package dao

import (
	"context"
	"net/http"
	"testing"
	"time"

	"go-common/app/admin/main/upload/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoNewBfs(t *testing.T) {
	convey.Convey("NewBfs", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctx.Convey("d.Bfs should not be nil", func(ctx convey.C) {
				ctx.So(d.Bfs, convey.ShouldNotBeNil)
			})
		})
	})
}

//func TestDaowaterMark(t *testing.T) {
//	convey.Convey("waterMark", t, func(ctx convey.C) {
//		bts, err := ioutil.ReadFile("/Users/sunsuke/Desktop/hahaha.png")
//		fmt.Println(err)
//		//for _, byt := range bts {
//		//fmt.Print(byt, ",")
//		//}
//
//		var (
//			c = context.Background()
//			// png file
//			data        = bts
//			contentType = http.DetectContentType(data)
//			wmKey       = ""
//			wmText      = "test"
//			paddingX    = int(0)
//			paddingY    = int(0)
//			wmScale     = float64(0)
//		)
//		fmt.Println(contentType)
//		ctx.Convey("When everything goes positive", func(ctx convey.C) {
//			res, err := d.Bfs.waterMark(c, data, contentType, wmKey, wmText, paddingX, paddingY, wmScale)
//			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}

func TestDaoUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			data = []byte("ut test")
			up   = &model.UploadParam{
				Bucket:      "test",
				Auth:        d.Bfs.Authorize("221bce6492eba70f", "6eb80603e85842542f9736eb13b7e3", http.MethodPut, "test", "", time.Now().Unix()),
				ContentType: http.DetectContentType(data),
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			location, etag, err := d.Bfs.Upload(c, up, data)
			ctx.Convey("Then err should be nil.location,etag should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(etag, convey.ShouldNotBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSpecificUpload(t *testing.T) {
	convey.Convey("SpecificUpload", t, func(ctx convey.C) {
		var (
			c           = context.Background()
			data        = []byte("ut specific upload")
			contentType = http.DetectContentType(data)
			auth        = d.Bfs.Authorize("221bce6492eba70f", "6eb80603e85842542f9736eb13b7e3", http.MethodPut, "test", "", time.Now().Unix())
			bucket      = "test"
			fileName    = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			location, etag, err := d.Bfs.SpecificUpload(c, contentType, auth, bucket, fileName, data)
			ctx.Convey("Then err should be nil.location,etag should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(etag, convey.ShouldNotBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoAuthorize(t *testing.T) {
	convey.Convey("Authorize", t, func(ctx convey.C) {
		var (
			key      = "221bce6492eba70f"
			secret   = "6eb80603e85842542f9736eb13b7e3"
			method   = http.MethodPut
			bucket   = "test"
			fileName = ""
			expire   = time.Now().Unix()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			authorization := d.Bfs.Authorize(key, secret, method, bucket, fileName, expire)
			ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
				ctx.So(authorization, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoClose(t *testing.T) {
	convey.Convey("close", t, func(ctx convey.C) {
		d.Close()
	})
}

func TestDaoPing(t *testing.T) {
	convey.Convey("ping", t, func(ctx convey.C) {
		d.Ping(context.Background())
	})
}
