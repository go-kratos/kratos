package dao

import (
	"context"
	"go-common/app/admin/main/growup/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoUpdateTagState(t *testing.T) {
	convey.Convey("UpdateTagState", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tagID     = int(100)
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO tag_info(id, is_deleted) VALUES(100, 1) ON DUPLICATE KEY UPDATE is_deleted = 1")
			rows, err := d.UpdateTagState(c, tagID, isDeleted)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertTagUpInfo(t *testing.T) {
	convey.Convey("TxInsertTagUpInfo", t, func(ctx convey.C) {
		var (
			tx, _     = d.BeginTran(context.Background())
			tagID     = int64(100)
			mid       = int64(1000)
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			d.Exec(context.Background(), "DELETE FROM tag_up_info WHERE tag_id = 100")
			rows, err := d.TxInsertTagUpInfo(tx, tagID, mid, isDeleted)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertTagUpInfo(t *testing.T) {
	convey.Convey("InsertTagUpInfo", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tagID     = int64(100)
			mid       = int64(1000)
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "DELETE FROM tag_up_info WHERE tag_id = 100")
			rows, err := d.InsertTagUpInfo(c, tagID, mid, isDeleted)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTagCom(t *testing.T) {
	convey.Convey("UpdateTagCom", t, func(ctx convey.C) {
		var (
			c        = context.Background()
			tagID    = int(100)
			isCommon = int(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO tag_info(id, is_common) VALUES(100, 0) ON DUPLICATE KEY UPDATE is_common = 1")
			rows, err := d.UpdateTagCom(c, tagID, isCommon)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoInsertTag(t *testing.T) {
	convey.Convey("InsertTag", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			tag = &model.TagInfo{
				Tag:      "tt",
				Category: 1,
				Business: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "DELETE FROM tag_info WHERE tag = 'tt'")
			rows, err := d.InsertTag(c, tag)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxInsertTag(t *testing.T) {
	convey.Convey("TxInsertTag", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTran(context.Background())
			tag   = &model.TagInfo{
				Tag:      "tt",
				Category: 1,
				Business: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			d.Exec(context.Background(), "DELETE FROM tag_info WHERE tag = 'tt'")
			rows, err := d.TxInsertTag(tx, tag)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTagInfo(t *testing.T) {
	convey.Convey("UpdateTagInfo", t, func(ctx convey.C) {
		var (
			c  = context.Background()
			no = &model.TagInfo{
				Tag:      "tt",
				Category: 1,
				Business: 3,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpdateTagInfo(c, no)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagInfo(t *testing.T) {
	convey.Convey("GetTagInfo", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = int(101)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO tag_info(id, tag) VALUES(101, 'kkkkk')")
			info, err := d.GetTagInfo(c, tagID)
			ctx.Convey("Then err should be nil.info should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(info, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagInfoByName(t *testing.T) {
	convey.Convey("GetTagInfoByName", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tag       = "ppp"
			dimension = int(0)
			category  = int(0)
			business  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO tag_info(id, tag) VALUES(102, 'ppp')")
			id, err := d.GetTagInfoByName(c, tag, dimension, category, business)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTxGetTagInfoByName(t *testing.T) {
	convey.Convey("TxGetTagInfoByName", t, func(ctx convey.C) {
		var (
			tx, _     = d.BeginTran(context.Background())
			tag       = "ppp"
			dimension = int(0)
			category  = int(0)
			business  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			defer tx.Commit()
			id, err := d.TxGetTagInfoByName(tx, tag, dimension, category, business)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoTagsCount(t *testing.T) {
	convey.Convey("TagsCount", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.TagsCount(c, query)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagInfos(t *testing.T) {
	convey.Convey("GetTagInfos", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			query = ""
			from  = int(0)
			limit = int(10)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			tagInfos, err := d.GetTagInfos(c, query, from, limit)
			ctx.Convey("Then err should be nil.tagInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagInfos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetNickname(t *testing.T) {
	convey.Convey("GetNickname", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			mid = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO up_category_info(mid, nick_name) VALUES(100, 'tt')")
			nickname, err := d.GetNickname(c, mid)
			ctx.Convey("Then err should be nil.nickname should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(nickname, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoGetTagUpInfoMID(t *testing.T) {
	convey.Convey("GetTagUpInfoMID", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			tagID     = int64(100)
			isDeleted = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mids, err := d.GetTagUpInfoMID(c, tagID, isDeleted)
			ctx.Convey("Then err should be nil.mids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoUpdateTagActivity(t *testing.T) {
	convey.Convey("UpdateTagActivity", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			tagID      = int64(102)
			activityID = int64(100)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			d.Exec(c, "INSERT INTO tag_info(id, activity_id) VALUES(102, 100) ON DUPLICATE KEY UPDATE activity_id = 100")
			rows, err := d.UpdateTagActivity(c, tagID, activityID)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
