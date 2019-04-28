package stat

import (
	"context"
	"testing"

	favmdl "go-common/app/service/main/favorite/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestStatfolderStatMcKey(t *testing.T) {
	convey.Convey("folderStatMcKey", t, func(convCtx convey.C) {
		var (
			table = ""
			fid   = int64(111)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			p1 := folderStatMcKey(table, fid)
			convCtx.Convey("Then p1 should not be nil.", func(convCtx convey.C) {
				convCtx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestStatSetFolderStatMc(t *testing.T) {
	convey.Convey("SetFolderStatMc", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(111)
			s  = &favmdl.Folder{}
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			err := d.SetFolderStatMc(c, id, s)
			convCtx.Convey("Then err should be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestStatFolderStatMc(t *testing.T) {
	convey.Convey("FolderStatMc", t, func(convCtx convey.C) {
		var (
			c  = context.Background()
			id = int64(111)
		)
		convCtx.Convey("When everything goes positive", func(convCtx convey.C) {
			f, err := d.FolderStatMc(c, id)
			convCtx.Convey("Then err should be nil.f should not be nil.", func(convCtx convey.C) {
				convCtx.So(err, convey.ShouldBeNil)
				convCtx.So(f, convey.ShouldNotBeNil)
			})
		})
	})
}
