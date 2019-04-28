package medal

import (
	"go-common/app/service/main/usersuit/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestMedalactivatedNidKey(t *testing.T) {
	convey.Convey("activatedNidKey", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := activatedNidKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalownersKey(t *testing.T) {
	convey.Convey("ownersKey", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := ownersKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalRedPointKey(t *testing.T) {
	convey.Convey("RedPointKey", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := RedPointKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalPopupKey(t *testing.T) {
	convey.Convey("PopupKey", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := PopupKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalpingMC(t *testing.T) {
	convey.Convey("pingMC", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.pingMC(c)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalSetMedalOwnersache(t *testing.T) {
	convey.Convey("SetMedalOwnersache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
			nos = []*model.MedalOwner{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetMedalOwnersache(c, mid, nos)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalMedalOwnersCache(t *testing.T) {
	convey.Convey("MedalOwnersCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
			nos = []*model.MedalOwner{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.SetMedalOwnersache(c, mid, nos)
			res, notFound, err := d.MedalOwnersCache(c, mid)
			ctx.Convey("Then err should be nil.res,notFound should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(notFound, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalDelMedalOwnersCache(t *testing.T) {
	convey.Convey("DelMedalOwnersCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMedalOwnersCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalmedalsActivatedCache(t *testing.T) {
	convey.Convey("medalsActivatedCache", t, func(ctx convey.C) {
		var (
			mids = []int64{88889017}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nids, missed, err := d.medalsActivatedCache(c, mids)
			ctx.Convey("Then err should be nil.nids,missed should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(missed, convey.ShouldNotBeNil)
				ctx.So(nids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalSetMedalActivatedCache(t *testing.T) {
	convey.Convey("SetMedalActivatedCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
			nid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetMedalActivatedCache(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalmedalActivatedCache(t *testing.T) {
	convey.Convey("medalActivatedCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nid, notFound, err := d.medalActivatedCache(c, mid)
			ctx.Convey("Then err should be nil.nid,notFound should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(notFound, convey.ShouldNotBeNil)
				ctx.So(nid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalDelMedalActivatedCache(t *testing.T) {
	convey.Convey("DelMedalActivatedCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelMedalActivatedCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalPopupCache(t *testing.T) {
	convey.Convey("PopupCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nid, err := d.PopupCache(c, mid)
			ctx.Convey("Then err should be nil.nid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalSetPopupCache(t *testing.T) {
	convey.Convey("SetPopupCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
			nid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetPopupCache(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalDelPopupCache(t *testing.T) {
	convey.Convey("DelPopupCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelPopupCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalRedPointCache(t *testing.T) {
	convey.Convey("RedPointCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			nid, err := d.RedPointCache(c, mid)
			ctx.Convey("Then err should be nil.nid should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nid, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMedalSetRedPointCache(t *testing.T) {
	convey.Convey("SetRedPointCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
			nid = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetRedPointCache(c, mid, nid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMedalDelRedPointCache(t *testing.T) {
	convey.Convey("DelRedPointCache", t, func(ctx convey.C) {
		var (
			mid = int64(88889017)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelRedPointCache(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
