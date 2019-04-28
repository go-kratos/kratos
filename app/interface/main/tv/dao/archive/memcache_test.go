package archive

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model/view"
	arcwar "go-common/app/service/main/archive/api"

	"github.com/smartystreets/goconvey/convey"
)

func TestArchivekeyRl(t *testing.T) {
	var (
		aid = int64(123)
	)
	convey.Convey("keyRl", t, func(ctx convey.C) {
		p1 := keyRl(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivekeyView(t *testing.T) {
	var (
		aid = int64(123)
	)
	convey.Convey("keyView", t, func(ctx convey.C) {
		p1 := keyView(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchivekeyArc(t *testing.T) {
	var (
		aid = int64(123)
	)
	convey.Convey("keyArc", t, func(ctx convey.C) {
		p1 := keyArc(aid)
		ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
			ctx.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveAddArcCache(t *testing.T) {
	var (
		aid = int64(123)
		arc = &arcwar.Arc{
			Aid: aid,
		}
	)
	convey.Convey("AddArcCache", t, func(ctx convey.C) {
		d.AddArcCache(aid, arc)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestArchiveAddRelatesCache(t *testing.T) {
	var (
		aid = int64(123)
		rls = []*view.Relate{
			{
				Aid: aid,
			},
		}
	)
	convey.Convey("AddRelatesCache", t, func(ctx convey.C) {
		d.AddRelatesCache(aid, rls)
		d.addRelatesCache(context.TODO(), aid, rls)
		ctx.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestArchiveAddViewCache(t *testing.T) {
	convey.Convey("AddViewCache", t, func(c convey.C) {
		aid, errGet := getPassAid(d.db)
		if errGet != nil {
			fmt.Println(errGet)
			return
		}
		vp, err := d.view3(context.Background(), aid)
		fmt.Println(vp, " ", aid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(vp, convey.ShouldNotBeNil)
		d.AddViewCache(aid, vp)
		c.Convey("No return values", func(ctx convey.C) {
		})
	})
}

func TestArchiveaddViewCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(123)
		vp  = &arcwar.ViewReply{}
	)
	convey.Convey("addViewCache", t, func(ctx convey.C) {
		err := d.addViewCache(c, aid, vp)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveRelatesCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(123)
	)
	convey.Convey("RelatesCache", t, func(ctx convey.C) {
		rls, err := d.RelatesCache(c, aid)
		ctx.Convey("Then err should be nil.rls should not be nil.", func(ctx convey.C) {
			fmt.Println(rls)
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(rls, convey.ShouldNotBeNil)
		})
	})
}

func TestArchiveaddRelatesCache(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(123)
		rls = []*view.Relate{
			{Aid: aid},
		}
	)
	convey.Convey("addRelatesCache", t, func(ctx convey.C) {
		err := d.addRelatesCache(c, aid, rls)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchiveaddArcCache(t *testing.T) {
	var (
		c      = context.Background()
		aid    = int64(123)
		cached = &arcwar.Arc{Aid: 123}
	)
	convey.Convey("addArcCache", t, func(ctx convey.C) {
		err := d.addArcCache(c, aid, cached)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestArchivearcsCache(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{123}
	)
	convey.Convey("arcsCache", t, func(ctx convey.C) {
		cached, missed, err := d.arcsCache(c, aids)
		ctx.Convey("Then err should be nil.cached,missed should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(len(missed)+len(cached), convey.ShouldEqual, len(aids))
			fmt.Println(cached)
			fmt.Println(missed)
		})
	})
}
