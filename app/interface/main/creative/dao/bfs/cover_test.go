package bfs

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

var mockHeader = map[string]string{"location": "mockLocation", "code": "200"}

func TestBfsUpVideoCovers(t *testing.T) {
	var (
		c      = context.TODO()
		covers = []string{
			"http://static.hdslb.com/images/transparent.gif",
		}
	)
	convey.Convey("UpVideoCovers", t, func(ctx convey.C) {
		httpMock(_method, _url).Reply(200).SetHeaders(mockHeader)
		httpMock("GET", covers[0]).Reply(200).JSON("mock byte")
		cvs, err := d.UpVideoCovers(c, covers)
		ctx.Convey("Then err should be nil.cvs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(cvs, convey.ShouldResemble, []string{"mockLocation"})
		})
	})
}

func TestBfsbvcCover(t *testing.T) {
	var (
		url = "http://static.hdslb.com/images/transparent.gif"
	)
	convey.Convey("bvcCover", t, func(ctx convey.C) {
		httpMock("GET", url).Reply(200).JSON("mock byte")
		bs, err := d.bvcCover(url)
		ctx.Convey("Then err should be nil.bs should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(bs, convey.ShouldNotBeNil)
		})
	})
}
