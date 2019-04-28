package dao

import (
	"context"
	"testing"

	"go-common/app/interface/main/playlist/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoplArcHit(t *testing.T) {
	var (
		pid = int64(1)
	)
	convey.Convey("plArcHit", t, func(ctx convey.C) {
		p1 := plArcHit(pid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVideo(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
		aid = int64(13825646)
	)
	convey.Convey("Video", t, func(ctx convey.C) {
		res, err := d.Video(c, pid, aid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoVideos(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
	)
	convey.Convey("Videos", t, func(ctx convey.C) {
		_, err := d.Videos(c, pid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoAddArc(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aid  = int64(3)
		sort = int64(0)
		desc = "abc"
	)
	convey.Convey("AddArc", t, func(ctx convey.C) {
		lastID, err := d.AddArc(c, pid, aid, sort, desc)
		ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lastID, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchAddArc(t *testing.T) {
	var (
		c        = context.Background()
		pid      = int64(1)
		arcSorts = []*model.ArcSort{}
	)
	convey.Convey("BatchAddArc", t, func(ctx convey.C) {
		arcSorts = append(arcSorts, &model.ArcSort{Aid: 13825646, Desc: "abc", Sort: 100})
		lastID, err := d.BatchAddArc(c, pid, arcSorts)
		ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lastID, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDelArc(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
		aid = int64(13825646)
	)
	convey.Convey("DelArc", t, func(ctx convey.C) {
		affected, err := d.DelArc(c, pid, aid)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchDelArc(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aids = []int64{13825646, 11, 22, 33}
	)
	convey.Convey("BatchDelArc", t, func(ctx convey.C) {
		affected, err := d.BatchDelArc(c, pid, aids)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateArcDesc(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aid  = int64(13825646)
		desc = "abc"
	)
	convey.Convey("UpdateArcDesc", t, func(ctx convey.C) {
		affected, err := d.UpdateArcDesc(c, pid, aid, desc)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdateArcSort(t *testing.T) {
	var (
		c    = context.Background()
		pid  = int64(1)
		aid  = int64(13825646)
		sort = int64(0)
	)
	convey.Convey("UpdateArcSort", t, func(ctx convey.C) {
		affected, err := d.UpdateArcSort(c, pid, aid, sort)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoBatchUpdateArcSort(t *testing.T) {
	var (
		c        = context.Background()
		pid      = int64(1)
		arcSorts = []*model.ArcSort{}
	)
	convey.Convey("BatchUpdateArcSort", t, func(ctx convey.C) {
		arcSorts = append(arcSorts, &model.ArcSort{Aid: 1, Desc: "abc"})
		affected, err := d.BatchUpdateArcSort(c, pid, arcSorts)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAdd(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
		fid = int64(1)
	)
	convey.Convey("Add", t, func(ctx convey.C) {
		lastID, err := d.Add(c, mid, fid)
		ctx.Convey("Then err should be nil.lastID should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(lastID, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoDel(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
	)
	convey.Convey("Del", t, func(ctx convey.C) {
		affected, err := d.Del(c, pid)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoUpdate(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
	)
	convey.Convey("Update", t, func(ctx convey.C) {
		affected, err := d.Update(c, pid)
		ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(affected, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPlsByMid(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(2)
	)
	convey.Convey("PlsByMid", t, func(ctx convey.C) {
		res, err := d.PlsByMid(c, mid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoPlsByPid(t *testing.T) {
	var (
		c    = context.Background()
		pids = []int64{1, 2, 3}
	)
	convey.Convey("PlsByPid", t, func(ctx convey.C) {
		_, err := d.PlsByPid(c, pids)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDaoPlByPid(t *testing.T) {
	var (
		c   = context.Background()
		pid = int64(1)
	)
	convey.Convey("PlByPid", t, func(ctx convey.C) {
		res, err := d.PlByPid(c, pid)
		ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(res, convey.ShouldNotBeNil)
		})
	})
}
