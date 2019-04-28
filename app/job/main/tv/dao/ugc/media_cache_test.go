package ugc

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	ugcmdl "go-common/app/job/main/tv/model/ugc"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_TotalVideos(t *testing.T) {
	Convey("TestDao_TotalVideos", t, WithDao(func(d *Dao) {
		res, err := d.TotalVideos(ctx)
		So(err, ShouldBeNil)
		So(res, ShouldBeGreaterThan, 0)
		fmt.Println(res)
	}))
}

func TestDao_ArcVideoCnt(t *testing.T) {
	Convey("TestDao_ArcVideoCnt", t, WithDao(func(d *Dao) {
		var aid int64
		d.DB.QueryRow(ctx, "select aid from ugc_video where deleted = 0 limit 1").Scan(&aid)
		if aid == 0 {
			fmt.Println("empty arc")
			return
		}
		cnt, errcnt := d.ArcVideoCnt(ctx, aid)
		So(errcnt, ShouldBeNil)
		fmt.Println(aid)
		So(cnt, ShouldBeGreaterThan, 0)
		fmt.Println(cnt)
		resVideos, lastIDVideo, errVd := d.PickArcVideo(ctx, aid, 0, 10)
		So(errVd, ShouldBeNil)
		So(len(resVideos), ShouldBeGreaterThan, 0)
		So(lastIDVideo, ShouldBeGreaterThan, 0)
		str, _ := json.Marshal(resVideos)
		fmt.Println(string(str))
	}))
}

func TestDao_SetArc(t *testing.T) {
	Convey("TestDao_SetArc", t, WithDao(func(d *Dao) {
		err := d.SetArcCMS(ctx, &ugcmdl.ArcCMS{
			Title: "testtest",
			AID:   777,
		})
		So(err, ShouldBeNil)
	}))
}

func TestDao_UpArcsCnt(t *testing.T) {
	Convey("TestDao_UpArcsCnt", t, WithDao(func(d *Dao) {
		var mid int64
		d.DB.QueryRow(ctx, "select mid from ugc_archive where deleted = 0 limit 1").Scan(&mid)
		if mid == 0 {
			fmt.Println("empty arc")
			return
		}
		count, err := d.UpArcsCnt(ctx, mid)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		fmt.Println("mid ", mid, " cnt", count)
		if count > 1 {
			d.DB.Exec(context.Background(), "update ugc_archive set deleted = 1 where mid = ? and deleted = 0 limit 1", mid)
			cntNonDeleted, err2 := d.CountUpArcs(context.Background(), mid)
			So(err2, ShouldBeNil)
			So(count, ShouldBeGreaterThan, cntNonDeleted)
			fmt.Println("all: ", count, " non-deleted: ", cntNonDeleted)
		}
	}))
}

func TestDao_TransFailVideos(t *testing.T) {
	Convey("TestDao_TransFailVideos", t, WithDao(func(d *Dao) {
		query := "SELECT aid FROM ugc_video WHERE  cid > 12780000 AND transcoded = 2 and deleted = 0 limit 1"
		var aid int64
		d.DB.QueryRow(context.Background(), query).Scan(&aid)
		if aid == 0 {
			fmt.Println("Empty archives")
			return
		}
		cids, err := d.TransFailVideos(ctx, aid)
		So(err, ShouldBeNil)
		So(len(cids), ShouldBeGreaterThan, 0)
		fmt.Println("aid ", aid, " cids ", cids)
	}))
}

func TestDao_ActVideos(t *testing.T) {
	Convey("TestDao_ActVideos", t, WithDao(func(d *Dao) {
		var (
			aid = int64(88888888)
			cid = 99999999
		)
		insertSQL := "REPLACE INTO ugc_video (aid, cid, deleted) VALUES (%d, %d, 1)"
		d.DB.Exec(ctx, fmt.Sprintf(insertSQL, aid, cid))
		has, err := d.ActVideos(ctx, aid)
		So(err, ShouldBeNil)
		So(has, ShouldBeFalse)
		d.DB.Exec(ctx, "UPDATE ugc_video SET deleted = 0 WHERE cid = ?", cid)
		has, err = d.ActVideos(ctx, aid)
		So(err, ShouldBeNil)
		So(has, ShouldBeTrue)
	}))
}

func TestDao_PickArcMC(t *testing.T) {
	Convey("TestDao_PickArcMC", t, WithDao(func(d *Dao) {
		pickMid := "select mid from ugc_archive where deleted = 0 group by mid order by count(aid) desc limit 1"
		var mid = 0
		d.DB.QueryRow(ctx, pickMid).Scan(&mid)
		if mid == 0 {
			fmt.Println("empty archive")
			return
		}
		fmt.Println("mid ", mid)
		res1, err1 := d.PickUpArcs(ctx, mid, 0, 5)
		So(err1, ShouldBeNil)
		So(len(res1), ShouldBeGreaterThan, 0)
		res2, err2 := d.PickUpArcs(ctx, mid, 7, 5)
		So(err2, ShouldBeNil)
		So(len(res2), ShouldBeGreaterThan, 0)
		str1, _ := json.Marshal(res1)
		str2, _ := json.Marshal(res2)
		fmt.Println(string(str1))
		fmt.Println(string(str2))
	}))
}
