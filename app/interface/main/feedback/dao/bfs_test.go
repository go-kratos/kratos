package dao

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/bouk/monkey"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpload(t *testing.T) {
	convey.Convey("Upload", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			fileName = ""
			fileType = ""
			expire   = int64(0)
			body     io.Reader
		)
		bfsConf := d.c.Bfs
		url := fmt.Sprintf(bfsConf.Addr+_uploadURL, bfsConf.Bucket, fileName)
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			httpMock("PUT", url).Reply(200).SetHeaders(map[string]string{
				"Code":     "200",
				"Location": "SomePlace",
			})
			location, err := d.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should be nil.location should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(location, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When d.bfsClient.Do gets error", func(ctx convey.C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.bfsClient), "Do",
				func(_ *http.Client, _ *http.Request) (*http.Response, error) {
					return nil, fmt.Errorf("d.bfsClient.Do Error")
				})
			defer guard.Unpatch()
			_, err := d.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request status != 200", func(ctx convey.C) {
			httpMock("PUT", url).Reply(404)
			_, err := d.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
		ctx.Convey("When http request Code in header != 200", func(ctx convey.C) {
			httpMock("PUT", url).Reply(404).SetHeaders(map[string]string{
				"Code":     "404",
				"Location": "SomePlace",
			})
			_, err := d.Upload(c, fileName, fileType, expire, body)
			ctx.Convey("Then err should not be nil", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoauthorize(t *testing.T) {
	var (
		key    = "1234"
		secret = "2345"
		method = ""
		bucket = ""
		file   = ""
		expire = int64(0)
	)
	convey.Convey("authorize", t, func(ctx convey.C) {
		authorization := authorize(key, secret, method, bucket, file, expire)
		ctx.Convey("Then authorization should not be nil.", func(ctx convey.C) {
			ctx.So(authorization, convey.ShouldNotBeNil)
		})
	})
}
