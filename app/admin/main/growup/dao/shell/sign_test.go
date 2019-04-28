package shell

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoSign(t *testing.T) {
	convey.Convey("Sign", t, func(ctx convey.C) {
		var (
			v struct {
				Name string
			}
			token = "abc"
		)
		v.Name = "aaa"
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := Sign(v, token)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoEncode(t *testing.T) {
	convey.Convey("Encode", t, func(ctx convey.C) {
		var (
			v struct {
				Name string
			}
		)
		v.Name = "aaa"
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := Encode(v)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
