package danmu

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestDanmuResetUpBanned(t *testing.T) {
	var (
		c     = context.TODO()
		mid   = int64(2089809)
		state = int8(0)
		hash  = "iamahashstring"
		ip    = "127.0.0.1"
	)
	convey.Convey("ResetUpBanned", t, func(ctx convey.C) {
		defer gock.OffAll()
		httpMock("POST", d.assistDmBannedURL).Reply(200).JSON(`{"code":0,"message":"0","ttl":1,"data":""}`)
		err := d.ResetUpBanned(c, mid, state, hash, ip)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}
