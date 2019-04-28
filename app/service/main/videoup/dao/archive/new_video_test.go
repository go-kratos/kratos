package archive

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go-common/app/service/main/videoup/model/archive"
)

func TestDao_TxAddVideoCid(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
	)
	Convey("TxAddVideoCid", t, func(ctx C) {
		_, err := d.TxAddVideoCid(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_AddNewVideo(t *testing.T) {
	var (
		c = context.Background()
		v = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
	)
	Convey("AddNewVideo", t, func(ctx C) {
		_, err := d.AddNewVideo(c, v)
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAddNewVideo(t *testing.T) {
	var (
		c = context.Background()
		v = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
		tx, _ = d.BeginTran(c)
	)
	Convey("TxAddNewVideo", t, func(ctx C) {
		_, err := d.TxAddNewVideo(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAddVideoRelation(t *testing.T) {
	var (
		c = context.Background()
		v = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
		tx, _ = d.BeginTran(c)
	)
	Convey("TxAddVideoRelation", t, func(ctx C) {
		_, err := d.TxAddVideoRelation(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVideoRelation(t *testing.T) {
	var (
		c = context.Background()
		v = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVideoRelation", t, func(ctx C) {
		_, err := d.TxUpVideoRelation(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpRelationState(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpRelationState", t, func(ctx C) {
		_, err := d.TxUpRelationState(tx, 23333, 1212, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpVdoStatus(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
	)
	Convey("TxUpVdoStatus", t, func(ctx C) {
		_, err := d.TxUpVdoStatus(tx, 1212, 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}
func TestDao_TxUpNewVideo(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		v     = &archive.Video{
			Aid:   23333,
			Cid:   12121,
			Title: "sssss",
		}
	)
	Convey("TxUpNewVideo", t, func(ctx C) {
		_, err := d.TxUpNewVideo(tx, v)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_NewCidsByFns(t *testing.T) {
	var (
		c  = context.Background()
		vs = []*archive.Video{{
			Filename: "1212121243gf",
		}}
	)
	Convey("NewCidsByFns", t, func(ctx C) {
		_, err := d.NewCidsByFns(c, vs)
		So(err, ShouldBeNil)
	})
}

func TestDao_CheckNewVideosTimeout(t *testing.T) {
	var (
		c  = context.Background()
		fs = []string{"1212121243gf"}
	)
	Convey("CheckNewVideosTimeout", t, func(ctx C) {
		_, _, err := d.CheckNewVideosTimeout(c, fs)
		So(err, ShouldBeNil)
	})
}

func TestDao_ParseDimensions(t *testing.T) {
	Convey("CheckNewVideosTimeout", t, func(ctx C) {
		_, err := d.parseDimensions("1,2,3")
		So(err, ShouldBeNil)
	})
}

func TestArchiveNewVideoFn(t *testing.T) {
	var (
		c        = context.Background()
		filename = "23333333333"
	)
	Convey("NewVideoFn", t, func(ctx C) {
		_, err := d.NewVideoFn(c, filename)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideoByFn(t *testing.T) {
	var (
		c        = context.Background()
		filename = "23333333333"
	)
	Convey("NewVideoByFn", t, func(ctx C) {
		_, err := d.NewVideoByFn(c, filename)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveSimpleArcVideos(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(2333333)
	)
	Convey("SimpleArcVideos", t, func(ctx C) {
		_, err := d.SimpleArcVideos(c, aid)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideos(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(2333333)
	)
	Convey("NewVideos", t, func(ctx C) {
		_, err := d.NewVideos(c, aid)
		ctx.Convey("Then err should be nil.vs should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideoMap(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(23333)
	)
	Convey("NewVideoMap", t, func(ctx C) {
		_, _, err := d.NewVideoMap(c, aid)
		ctx.Convey("Then err should be nil.vm,cvm should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideoByCID(t *testing.T) {
	var (
		c   = context.Background()
		cid = int64(23333)
	)
	Convey("NewVideoByCID", t, func(ctx C) {
		_, err := d.NewVideoByCID(c, cid)
		ctx.Convey("Then err should be nil.v should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideosByCID(t *testing.T) {
	var (
		c    = context.Background()
		cids = []int64{23333}
	)
	Convey("NewVideosByCID", t, func(ctx C) {
		_, err := d.NewVideosByCID(c, cids)
		ctx.Convey("Then err should be nil.vm should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideosByFn(t *testing.T) {
	var (
		c   = context.Background()
		fns = []string{"23333"}
	)
	Convey("NewVideosByFn", t, func(ctx C) {
		_, err := d.NewVideosByFn(c, fns)
		ctx.Convey("Then err should be nil.vm should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveNewVideosReason(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(23333)
	)
	Convey("NewVideosReason", t, func(ctx C) {
		_, err := d.NewVideosReason(c, aid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}
