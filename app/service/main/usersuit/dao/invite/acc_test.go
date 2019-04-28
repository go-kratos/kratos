package dao

import (
	"testing"

	"github.com/smartystreets/goconvey/convey"
	"gopkg.in/h2non/gock.v1"
)

func TestDaoBeFormal(t *testing.T) {
	convey.Convey("BeFormal", t, func(ctx convey.C) {
		var (
			mid    = int64(88888970)
			cookie = "fts=1532412723; buvid3=9FE379B1-4742-4414-83BC-070DC6D5C30128620infoc; rpdid=oxlilxsoqdoskipoimqw; UM_distinctid=164cb9399eb124-053863286996f2-163b6952-1fa400-164cb9399ec15e; LIVE_BUVID=dfbc64db5d227b750394fc08dc27ad5e; LIVE_BUVID__ckMd5=612509d11892d4bd; im_notify_type_27956255=0; sid=66taifil; im_local_unread_27956255=0; im_seqno_27956255=4; CNZZDATA2724999=cnzz_eid%3D937547062-1534212022-%26ntime%3D1534212022; pgv_pvi=1353490432; finger=14bc3c4e; DedeUserID=27956255; DedeUserID__ckMd5=2dfe07340adf4bc9; SESSDATA=3f325bec%2C1538121389%2Cb8e04f4a; bili_jct=008e9de4d371cf91f089a66b77e1d8d7; PHPSESSID=3btocfs7f8mo3u8cpkhfv155e0; bp_t_offset_27956255=157578003585063279; _dfcaptcha=aa24be76e8c462474766bc1f4db916a0"
			ip     = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("POST", d.beFormalURI).Reply(0).JSON(`{"code":0}`)
			err := d.BeFormal(c, mid, cookie, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldNotBeNil)
			})
		})
	})
}
