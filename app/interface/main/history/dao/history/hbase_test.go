package history

import (
	"context"
	"testing"

	"go-common/app/interface/main/history/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestHistoryhashRowKey(t *testing.T) {
	convey.Convey("hashRowKey", t, func(ctx convey.C) {
		var (
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := hashRowKey(mid)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistorycolumn(t *testing.T) {
	convey.Convey("column", t, func(ctx convey.C) {
		var (
			aid = int64(14771787)
			typ = int8(3)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := d.column(aid, typ)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryAdd(t *testing.T) {
	convey.Convey("Add", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			h   = &model.History{Mid: 14771787, Aid: 32767458, TP: 3, Pro: 20, Unix: 1540976376, DT: 3}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Add(c, mid, h)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryAddMap(t *testing.T) {
	convey.Convey("AddMap", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			h   = &model.History{Mid: 14771787, Aid: 32767458, TP: 3, Pro: 20, Unix: 1540976376, DT: 3}
			hs  = map[int64]*model.History{14771787: h}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.AddMap(c, mid, hs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryAidsMap(t *testing.T) {
	convey.Convey("AidsMap", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{14771787}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			his, err := d.AidsMap(c, mid, aids)
			ctx.Convey("Then err should be nil.his should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(his, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryMap(t *testing.T) {
	convey.Convey("Map", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			his, err := d.Map(c, mid)
			ctx.Convey("Then err should be nil.his should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(his, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestHistoryDelAids(t *testing.T) {
	convey.Convey("DelAids", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			mid  = int64(14771787)
			aids = []int64{32767458}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.DelAids(c, mid, aids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestHistoryget(t *testing.T) {
	convey.Convey("delete", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, _, err := d.get(c, tableInfo, hashRowKey(mid))
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistorydelete(t *testing.T) {
	convey.Convey("delete", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			mid       = int64(14771787)
			delColumn = map[string][]byte{"info:14771787": []byte{}}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.delete(c, tableInfo, hashRowKey(mid), delColumn)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
func TestHistoryDelete(t *testing.T) {
	convey.Convey("Delete", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
			h   = &model.History{Mid: 14771787, Aid: 32767458, TP: 3, Pro: 20, Unix: 1540976376, DT: 3}
			his = []*model.History{h}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Delete(c, mid, his)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryClear(t *testing.T) {
	convey.Convey("Clear", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.Clear(c, mid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistorySetInfoShadow(t *testing.T) {
	convey.Convey("SetInfoShadow", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mid   = int64(14771787)
			value = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetInfoShadow(c, mid, value)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestHistoryInfoShadow(t *testing.T) {
	convey.Convey("InfoShadow", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(14771787)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sw, err := d.InfoShadow(c, mid)
			ctx.Convey("Then err should be nil.sw should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(sw, convey.ShouldNotBeNil)
			})
		})
	})
}
