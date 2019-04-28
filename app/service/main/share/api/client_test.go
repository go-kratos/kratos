package v1

import (
	"context"
	"testing"

	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

var client ShareClient

func init() {
	var err error
	client, err = NewClient(nil)
	if err != nil {
		panic(err)
	}
}

func TestAddShare(t *testing.T) {
	convey.Convey("TestAddShare", t, func(ctx convey.C) {
		var c = context.Background()
		ctx.Convey("When everything is correct", func(ctx convey.C) {
			reply, err := client.AddShare(c, &AddShareRequest{Oid: 22, Mid: 33, Type: 3})
			ctx.So(err, convey.ShouldBeNil)
			ctx.Printf("%+v\n", reply.Shares)
		})
		ctx.Convey("When error", func(ctx convey.C) {
			reply, err := client.AddShare(c, &AddShareRequest{Oid: 22, Mid: 33, Type: 3})
			ctx.So(err, convey.ShouldEqual, ecode.ShareAlreadyAdd)
			ctx.So(reply, convey.ShouldBeNil)
		})
	})
}
