package archive

import (
	"context"
	"testing"
	"time"

	"database/sql"
	"fmt"
	"go-common/app/service/main/videoup/model/archive"
	xsql "go-common/library/database/sql"
	"reflect"

	"github.com/bouk/monkey"
	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_TxAddArchive(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		a     = &archive.Archive{
			Mid:    123,
			TypeID: 22,
			Title:  "UT测试",
			Author: "ut",
			Desc:   "UT测试UT测试",
		}
	)
	Convey("TxUpArchiveState", t, func(ctx C) {
		_, err := d.TxAddArchive(tx, a)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAUpArchive(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		a     = &archive.Archive{
			Aid:    10111210,
			Mid:    123,
			TypeID: 22,
			Title:  "UT测试",
			Author: "ut",
			Desc:   "UT测试UT测试",
		}
	)
	Convey("TxUpArchiveState", t, func(ctx C) {
		_, err := d.TxUpArchive(tx, a)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxAUpArchiveMid(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		a     = &archive.Archive{
			Aid:    10111210,
			Mid:    123,
			TypeID: 22,
			Title:  "UT测试",
			Author: "ut",
			Desc:   "UT测试UT测试",
		}
	)
	Convey("TxUpArchiveState", t, func(ctx C) {
		_, err := d.TxUpArchiveMid(tx, a)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestArchiveTxUpArchiveState(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(233333)
		state = int8(0)
	)
	Convey("TxUpArchiveState", t, func(ctx C) {
		_, err := d.TxUpArchiveState(tx, aid, state)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpAddit(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(233333)
	)
	Convey("TxUpAddit", t, func(ctx C) {
		_, err := d.TxUpAddit(tx, aid, 0, 0, 0, 0, []byte{}, "", "", "", "", "", 0)
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpArchiveBiz(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(233333)
	)
	Convey("TxUpArchiveBiz", t, func(ctx C) {
		_, err := d.TxUpArchiveBiz(tx, aid, 0, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpAdditReason(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(233333)
	)
	Convey("TxUpAdditReason", t, func(ctx C) {
		_, err := d.TxUpAdditReason(tx, aid, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

func TestDao_TxUpAdditRedirect(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(233333)
	)
	Convey("TxUpAdditRedirect", t, func(ctx C) {
		_, err := d.TxUpAdditRedirect(tx, aid, "")
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
		So(err, ShouldBeNil)
	})
}

// func TestArchiveTxUpAdditReason(t *testing.T) {
// 	var (
// 		c      = context.Background()
// 		tx, _  = d.BeginTran(c)
// 		aid    = int64(233333)
// 		reason = "2333"
// 	)
// 	Convey("TxUpAdditReason", t, func(ctx C) {
// 		rows, err := d.TxUpAdditReason(tx, aid, reason)
// 		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx C) {
// 			ctx.So(err, ShouldBeNil)
// 			ctx.So(rows, ShouldNotBeNil)
// 		})
// 	})
// }

func TestArchiveTxUpAdditRedirect(t *testing.T) {
	var (
		c           = context.Background()
		tx, _       = d.BeginTran(c)
		aid         = int64(0)
		redirectURL = "233333"
	)
	Convey("TxUpAdditRedirect", t, func(ctx C) {
		rows, err := d.TxUpAdditRedirect(tx, aid, redirectURL)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(rows, ShouldNotBeNil)
		})
	})
}

func TestArchiveTxUpArcAttr(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(0)
		bit   = uint(0)
		val   = int32(0)
	)
	Convey("TxUpArcAttr", t, func(ctx C) {
		rows, err := d.TxUpArcAttr(tx, aid, bit, val)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(rows, ShouldNotBeNil)
		})
	})
}

func TestArchiveTxUpTag(t *testing.T) {
	var (
		c     = context.Background()
		tx, _ = d.BeginTran(c)
		aid   = int64(22333)
		tag   = "2333"
	)
	Convey("TxUpTag", t, func(ctx C) {
		rows, err := d.TxUpTag(tx, aid, tag)
		ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(rows, ShouldNotBeNil)
		})
	})
}

func TestArchiveArchive(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(22333)
	)
	Convey("Archive", t, func(ctx C) {
		_, err := d.Archive(c, aid)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveAddit(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(22333)
		aa  = &archive.ArcMissionParam{}
	)
	Convey("Addit", t, func(ctx C) {
		_, err := d.Addit(c, aid)
		ctx.Convey("Then err should be nil.ad should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
	Convey("Vote", t, func(ctx C) {
		ad, err := d.Vote(c, aid)
		ctx.Convey("Then err should be nil.ad should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(ad, ShouldBeNil)
		})
	})
	Convey("Recos", t, func(ctx C) {
		ad, err := d.Recos(c, aid)
		ctx.Convey("Then err should be nil.ad should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(ad, ShouldBeNil)
		})
	})
	Convey("UpMissionID", t, func(ctx C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.db), "Exec", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (sql.Result, error) {
			return nil, sql.ErrNoRows
		})
		defer guard.Unpatch()
		ad, err := d.UpMissionID(c, aa)
		ctx.Convey("Then err should be nil.ad should not be nil.", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
			ctx.So(ad, ShouldBeZeroValue)
		})
	})
}

func TestArchiveMids(t *testing.T) {
	var (
		c    = context.Background()
		aids = []int64{222}
	)
	Convey("Mids", t, func(ctx C) {
		mm, err := d.Mids(c, aids)
		ctx.Convey("Then err should be nil.mm should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(mm, ShouldNotBeNil)
		})
	})
}

func TestArchiveArchivesUpAll(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(223345)
		offset = int(20)
		ps     = int(1)
	)
	Convey("ArchivesUpAll", t, func(ctx C) {
		_, err := d.ArchivesUpAll(c, mid, offset, ps)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveArchivesUpOpen(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(223345)
		offset = int(20)
		ps     = int(1)
	)
	Convey("ArchivesUpOpen", t, func(ctx C) {
		_, err := d.ArchivesUpOpen(c, mid, offset, ps)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveArchivesUpUnOpen(t *testing.T) {
	var (
		c      = context.Background()
		mid    = int64(223345)
		offset = int(20)
		ps     = int(1)
	)
	Convey("ArchivesUpUnOpen", t, func(ctx C) {
		_, err := d.ArchivesUpUnOpen(c, mid, offset, ps)
		ctx.Convey("Then err should be nil.aids should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchiveArchiveAllUpCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(223345)
	)
	Convey("ArchiveAllUpCount", t, func(ctx C) {
		count, err := d.ArchiveAllUpCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(count, ShouldNotBeNil)
		})
	})
}

func TestArchiveArchiveOpenUpCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(223345)
	)
	Convey("ArchiveOpenUpCount", t, func(ctx C) {
		count, err := d.ArchiveOpenUpCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(count, ShouldNotBeNil)
		})
	})
}

func TestArchiveArchiveUnOpenUpCount(t *testing.T) {
	var (
		c   = context.Background()
		mid = int64(223345)
	)
	Convey("ArchiveUnOpenUpCount", t, func(ctx C) {
		count, err := d.ArchiveUnOpenUpCount(c, mid)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
			ctx.So(count, ShouldNotBeNil)
		})
	})
}

func TestArchiveSimpleArchive(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(222)
	)
	Convey("SimpleArchive", t, func(ctx C) {
		_, err := d.SimpleArchive(c, aid)
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchivePOI(t *testing.T) {
	var (
		c   = context.Background()
		aid = int64(222)
	)
	Convey("poi", t, func(ctx C) {
		data, err := d.POI(c, aid)
		fmt.Println(string(data))
		ctx.Convey("Then err should be nil.a should not be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestArchivePOIAdd(t *testing.T) {
	var (
		c     = context.Background()
		aid   = int64(222)
		tx, _ = d.BeginTran(c)
		err   error
	)
	Convey("add poi err", t, func(ctx C) {
		guard := monkey.PatchInstanceMethod(reflect.TypeOf(tx),
			"Exec",
			func(_ *xsql.Tx, _ string, _ ...interface{}) (sql.Result, error) {
				return nil, fmt.Errorf("tx.Exec error")
			})
		defer guard.Unpatch()
		_, err = d.TxUpArchiveBiz(tx, aid, 1, "2222")
		ctx.Convey("TestArchivePOIAdd.", func(ctx C) {
			ctx.So(err, ShouldNotBeNil)
		})
	})
}

func TestArchiverejectedArchives(t *testing.T) {
	Convey("rejectedArchives", t, func(ctx C) {
		var (
			c              = context.Background()
			mid      int64 = 2089809
			state    int32 = -4
			offset   int32
			limit    int32 = 20
			start, _       = time.Parse("20060102", "20100101")
		)
		ctx.Convey("When everything gose positive", func(ctx C) {
			arcs, count, err := d.RejectedArchives(c, mid, state, offset, limit, &start)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx C) {
				ctx.So(err, ShouldBeNil)
				ctx.So(arcs, ShouldNotBeNil)
				ShouldNotEqual(count, 0)
			})
		})
		ctx.Convey("When no rows found", func(ctx C) {
			guard := monkey.PatchInstanceMethod(reflect.TypeOf(d.slaveDB), "Query", func(_ *xsql.DB, _ context.Context, _ string, _ ...interface{}) (*xsql.Rows, error) {
				return nil, xsql.ErrNoRows
			})
			defer guard.Unpatch()
			arcs, _, err := d.RejectedArchives(c, mid, state, offset, limit, &start)
			ctx.Convey("Then err should be nil.arcs should not be nil.", func(ctx C) {
				ctx.So(err, ShouldNotBeNil)
				ctx.So(arcs, ShouldBeNil)
			})
		})
	})
}
