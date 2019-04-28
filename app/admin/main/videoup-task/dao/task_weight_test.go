package dao

import (
	"context"
	"go-common/app/admin/main/videoup-task/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoGetMaxWeight(t *testing.T) {
	convey.Convey("GetMaxWeight", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			max, err := d.GetMaxWeight(c)
			ctx.Convey("Then err should be nil.max should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(max, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpCwAfterAdd(t *testing.T) {
	convey.Convey("UpCwAfterAdd", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			id   = int64(0)
			desc = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpCwAfterAdd(c, id, desc)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInWeightConf(t *testing.T) {
	convey.Convey("InWeightConf", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcases map[int64]*model.WCItem
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.InWeightConf(c, mcases)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoDelWeightConf(t *testing.T) {
	convey.Convey("DelWeightConf", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.DelWeightConf(c, id)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoListWeightConf(t *testing.T) {
	convey.Convey("ListWeightConf", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			cf = &model.Confs{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ListWeightConf(c, cf)
			ctx.Convey("Then err should be nil.citems should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWeightConf(t *testing.T) {
	convey.Convey("WeightConf", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			items, err := d.WeightConf(c)
			ctx.Convey("Then err should be nil.items should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(items, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaokey(t *testing.T) {
	convey.Convey("key", t, func(ctx convey.C) {
		var (
			id = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := key(id)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetWeightRedis(t *testing.T) {
	convey.Convey("GetWeightRedis", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mcases, err := d.GetWeightRedis(c, ids)
			ctx.Convey("Then err should be nil.mcases should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mcases, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWeightVC(t *testing.T) {
	convey.Convey("WeightVC", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			wvc, err := d.WeightVC(c)
			ctx.Convey("Then err should be nil.wvc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(wvc, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoSetWeightVC(t *testing.T) {
	convey.Convey("SetWeightVC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			wvc  = &model.WeightVC{}
			desc = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.SetWeightVC(c, wvc, desc)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInWeightVC(t *testing.T) {
	convey.Convey("InWeightVC", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			wvc  = &model.WeightVC{}
			desc = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.InWeightVC(c, wvc, desc)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoLWConfigHelp(t *testing.T) {
	convey.Convey("LWConfigHelp", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.LWConfigHelp(c, []int64{1})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoSetWeightRedis(t *testing.T) {
	convey.Convey("SetWeightRedis", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			mcases = map[int64]*model.TaskPriority{
				0: &model.TaskPriority{},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.SetWeightRedis(c, mcases)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
