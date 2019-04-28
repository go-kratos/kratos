package creative

import (
	"context"
	"go-common/app/interface/main/videoup/model/archive"
	"testing"

	"github.com/smartystreets/goconvey/convey"
	gock "gopkg.in/h2non/gock.v1"
)

func TestCreativeSetWatermark(t *testing.T) {
	convey.Convey("SetWatermark", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(2089809)
			state = int8(1)
			ty    = int8(3)
			pos   = int8(1)
			ip    = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("Post", d.setWatermarkURL).Reply(200).JSON(`{"code":0,"data":{"id":53,"mid":27515310,"uname":"1qs314567","state":2,"type":2,"position":4,"url":"http://i0.hdslb.com/bfs/article/578a61e7caf47b5deaa80940b2806f1d9ce53dde.png","md5":"79d02f08c2b7b0d2b30a6b6a4f61c97e","info":"{\"width\":325,\"height\":50}","ctime":1499764110,"mtime":1499828335},"message":"","ttl":1}`)
			err := d.SetWatermark(c, mid, state, ty, pos, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestCreativeUploadMaterial(t *testing.T) {
	convey.Convey("UploadMaterial", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			editors = []*archive.Editor{}
			aid     = int64(10110826)
			mid     = int64(2089809)
			ip      = "127.0.0.1"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer gock.OffAll()
			httpMock("Post", d.uploadMaterialURL).Reply(200).JSON(`{"code":0,"data":""}`)
			err := d.UploadMaterial(c, editors, aid, mid, ip)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
