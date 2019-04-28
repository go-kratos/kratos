package ugc

import (
	"fmt"
	"testing"

	"go-common/app/job/main/tv/model/ugc"
	arccli "go-common/app/service/main/archive/api"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUgcTxDelArc(t *testing.T) {
	var (
		tx, _ = d.DB.Begin(ctx)
		aid   = int64(0)
	)
	Convey("TxDelArc", t, func(ctx C) {
		err := d.TxDelArc(tx, aid)
		ctx.Convey("Then err should be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestUgcTxDelVideos(t *testing.T) {
	var (
		tx, _ = d.DB.Begin(ctx)
		aid   = int64(0)
	)
	Convey("TxDelVideos", t, func(ctx C) {
		err := d.TxDelVideos(tx, aid)
		ctx.Convey("Then err should be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestUgcTxDelVideo(t *testing.T) {
	var (
		tx, _ = d.DB.Begin(ctx)
		cid   = int64(0)
	)
	Convey("TxDelVideo", t, func(ctx C) {
		err := d.TxDelVideo(tx, cid)
		ctx.Convey("Then err should be nil.", func(ctx C) {
			ctx.So(err, ShouldBeNil)
		})
	})
}

func TestDao_DelVideoArc(t *testing.T) {
	Convey("TestDao_DelVideoArc", t, WithDao(func(d *Dao) {
		var (
			aid            = int64(99998888)
			cid1           = int64(999988881)
			cid2           = int64(999988882)
			tx, _          = d.DB.Begin(ctx)
			arc            = &arccli.Arc{Aid: aid}
			countVQ        = "SELECT COUNT(1) FROM ugc_video WHERE aid = ? AND deleted = 0"
			countAQ        = "SELECT COUNT(1) FROM ugc_archive WHERE aid = ? AND deleted = 0"
			countV, countA int
			arcValid       bool
		)
		// add archive and two videos
		d.TxImportArc(tx, &arccli.Arc{Aid: aid})
		d.TxMnlVideos(tx, &arccli.ViewReply{
			Arc: arc,
			Pages: []*arccli.Page{
				{
					Cid: cid1,
				},
				{
					Cid: cid2,
				},
			},
		})
		tx.Commit()
		d.DB.QueryRow(ctx, countVQ, aid).Scan(&countV)
		So(countV, ShouldEqual, 2)
		// delete one video, still one active video under the archive, we keep the archive
		_, err := d.DelVideoArc(ctx, &ugc.DelVideos{
			AID:  aid,
			CIDs: []int64{cid1},
		})
		So(err, ShouldBeNil)
		d.DB.QueryRow(ctx, countVQ, aid).Scan(&countV)
		d.DB.QueryRow(ctx, countAQ, aid).Scan(&countA)
		So(countV, ShouldEqual, 1)
		So(countA, ShouldEqual, 1)
		So(err, ShouldBeNil)
		// delete the last video, the archive should also be deleted
		arcValid, err = d.DelVideoArc(ctx, &ugc.DelVideos{
			AID:  aid,
			CIDs: []int64{cid2},
		})
		So(err, ShouldBeNil)
		So(arcValid, ShouldBeFalse)
		d.DB.QueryRow(ctx, countVQ, aid).Scan(&countV)
		d.DB.QueryRow(ctx, countAQ, aid).Scan(&countA)
		So(countV, ShouldEqual, 0)
		So(countA, ShouldEqual, 0)
	}))
}

func TestDao_DelVideos(t *testing.T) {
	Convey("TestDao_DelVideos", t, WithDao(func(d *Dao) {
		var (
			aid       = int64(99998888)
			cid1      = 99998887
			cid2      = 99998886
			insertSQL = "REPLACE INTO ugc_video (aid, cid) VALUES (%d, %d)"
		)
		ress, err2 := d.DB.Exec(ctx, fmt.Sprintf(insertSQL, aid, cid1))
		fmt.Println(fmt.Sprintf(insertSQL, aid, cid1))
		fmt.Println(err2)
		fmt.Println(ress.RowsAffected())
		d.DB.Exec(ctx, fmt.Sprintf(insertSQL, aid, cid2))
		var count int
		d.DB.QueryRow(ctx, "SELECT COUNT(1) FROM ugc_video WHERE aid = ? AND deleted = 0", aid).Scan(&count)
		So(count, ShouldEqual, 2)
		err := d.DelVideos(ctx, aid)
		So(err, ShouldBeNil)
		d.DB.QueryRow(ctx, "SELECT COUNT(1) FROM ugc_video WHERE aid = ? AND deleted = 0", aid).Scan(&count)
		So(count, ShouldEqual, 0)
	}))
}
