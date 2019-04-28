package tag

import (
	"context"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestTagAllTagInfo(t *testing.T) {
	convey.Convey("AllTagInfo", t, func(ctx convey.C) {
		var (
			c = context.Background()
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(tag) VALUES('test')")
			tagInfos, err := d.AllTagInfo(c)
			ctx.Convey("Then err should be nil.tagInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagInfos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagGetTagInfoByDate(t *testing.T) {
	convey.Convey("GetTagInfoByDate", t, func(ctx convey.C) {
		var (
			c         = context.Background()
			dimension = int(1)
			ctype     = int(1)
			startAt   = "2018-06-01"
			endAt     = "2018-06-01"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(tag,dimension,category_id,business_id,adjust_type,ratio,is_common,start_at,end_at) VALUES('utt', 1, 1, 1, 1, 1, 1, '2018-05-01', '2018-10-01')")
			tagInfos, err := d.GetTagInfoByDate(c, dimension, ctype, startAt, endAt)
			ctx.Convey("Then err should be nil.tagInfos should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(tagInfos, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagTxUpdateTagInfoIncome(t *testing.T) {
	convey.Convey("TxUpdateTagInfoIncome", t, func(ctx convey.C) {
		var (
			c      = context.Background()
			tx, _  = d.db.Begin(c)
			id     = int64(100001)
			income = int64(100)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			defer tx.Commit()
			Exec(c, "INSERT INTO tag_info(id, tag) VALUES(1000001, 'uttt2')")
			rows, err := d.TxUpdateTagInfoIncome(tx, id, income)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestTagUpdateTagUps(t *testing.T) {
	convey.Convey("UpdateTagUps", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			tagID = int64(100001)
			ups   = int(10)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			Exec(c, "INSERT INTO tag_info(id, tag) VALUES(1000001, 'uttt2')")
			rows, err := d.UpdateTagUps(c, tagID, ups)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
