package dao

import (
	"context"
	"testing"

	"go-common/app/service/main/tag/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoreplace(t *testing.T) {
	convey.Convey("replace", t, func(ctx convey.C) {
		var (
			name = "搞笑"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := replace(name)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTag(t *testing.T) {
	convey.Convey("Tag", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(2)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no, err := d.Tag(c, tid)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTags(t *testing.T) {
	convey.Convey("Tags", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Tags(c, tids)
		})
	})
}

func TestDaoTagByName(t *testing.T) {
	convey.Convey("TagByName", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			name = "数码"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			no, err := d.TagByName(c, name)
			ctx.Convey("Then err should be nil.no should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(no, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagsByNames(t *testing.T) {
	convey.Convey("TagsByNames", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			names = []string{"数码", "搞笑", "IDC"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tags, err := d.TagsByNames(c, names)
			ctx.Convey("Then err should be nil.tags should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tags, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxAddTag(t *testing.T) {
	convey.Convey("TxAddTag", t, func(ctx convey.C) {
		no := &model.Tag{
			Name: "tag-service test",
		}
		tx, _ := d.BeginTran(context.Background())
		defer tx.Rollback()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TxAddTag(tx, no)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
	})
}

func TestDaoCreateTag(t *testing.T) {
	convey.Convey("CreateTag", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			no = &model.Tag{
				Name: "unit test 233",
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.CreateTag(c, no)
		})
	})
}

func TestDaoCreateTags(t *testing.T) {
	convey.Convey("CreateTags", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			ts = []*model.Tag{
				{
					Name: "unit test 233",
				},
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.CreateTags(c, ts)
		})
	})
}

func TestDaoCount(t *testing.T) {
	convey.Convey("Count", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tid = int64(1833)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Count(c, tid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoCounts(t *testing.T) {
	convey.Convey("Counts", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			tids = []int64{1833}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Counts(c, tids)
		})
	})
}

func TestDaoTxUpTagBindCount(t *testing.T) {
	convey.Convey("TxUpTagBindCount", t, func(ctx convey.C) {
		var (
			tid   = int64(1833)
			count = int64(233)
		)
		tx, _ := d.BeginTran(context.Background())
		defer tx.Rollback()

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpTagBindCount(tx, tid, count)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoTxUpTagSubCount(t *testing.T) {
	convey.Convey("TxUpTagSubCount", t, func(ctx convey.C) {
		var (
			tid   = int64(1833)
			count = int64(23)
		)
		tx, _ := d.BeginTran(context.Background())
		defer tx.Rollback()

		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.TxUpTagSubCount(tx, tid, count)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldBeGreaterThanOrEqualTo, 0)
			})
		})
		tx.Commit()
	})
}

func TestDaoHots(t *testing.T) {
	convey.Convey("Hots", t, func(ctx convey.C) {
		var (
			c       = context.Background()
			rid     = int64(20)
			hotType = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Hots(c, rid, hotType)
		})
	})
}

func TestDaoRids(t *testing.T) {
	convey.Convey("Rids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			cm, pridMap, rids, err := d.Rids(c)
			ctx.Convey("Then err should be nil.cm,pridMap,rids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rids, convey.ShouldNotBeNil)
				ctx.So(pridMap, convey.ShouldNotBeNil)
				ctx.So(cm, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoHotMap(t *testing.T) {
	convey.Convey("HotMap", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.HotMap(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoPrids(t *testing.T) {
	convey.Convey("Prids", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Prids(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagGroup(t *testing.T) {
	convey.Convey("TagGroup", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			resMap, err := d.TagGroup(c)
			ctx.Convey("Then err should be nil.resMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(resMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecommandTagFilter(t *testing.T) {
	convey.Convey("RecommandTagFilter", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ts, err := d.RecommandTagFilter(c)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoRecommandTagTop(t *testing.T) {
	convey.Convey("RecommandTagTop", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ts, err := d.RecommandTagTop(c)
			ctx.Convey("Then err should be nil.ts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpTagState(t *testing.T) {
	convey.Convey("UpTagState", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tid   = int64(1833)
			state = int32(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affected, err := d.UpTagState(c, tid, state)
			ctx.Convey("Then err should be nil.affected should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affected, convey.ShouldNotBeNil)
			})
		})
	})
}
