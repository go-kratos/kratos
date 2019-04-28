package ugc

import (
	"context"
	ugcmdl "go-common/app/job/main/tv/model/ugc"
	"testing"

	"fmt"

	"encoding/json"

	"github.com/smartystreets/goconvey/convey"
)

func TestUgcPpVideos(t *testing.T) {
	var (
		c      = context.Background()
		videos = []int64{}
	)
	convey.Convey("PpVideos", t, func(ctx convey.C) {
		err := d.PpVideos(c, videos)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestUgcFinishVideos(t *testing.T) {
	var (
		c      = context.Background()
		videos = []*ugcmdl.SimpleVideo{}
		aid    = int64(0)
	)
	convey.Convey("FinishVideos", t, func(ctx convey.C) {
		err := d.FinishVideos(c, videos, aid)
		ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
		})
	})
}

func TestDao_ParseArc(t *testing.T) {
	convey.Convey("TestDao_ParseArc", t, WithDao(func(d *Dao) {
		var (
			query = "SELECT aid FROM ugc_archive WHERE deleted = 0 LIMIT 1"
			arc   = &ugcmdl.ArcCMS{}
		)
		d.DB.QueryRow(ctx, query).Scan(&arc.AID)
		res, err := d.ParseArc(ctx, int64(arc.AID))
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		fmt.Println(res)
	}))
}

func TestDao_ParseVideos(t *testing.T) {
	convey.Convey("TestDao_ParseVideos", t, WithDao(func(d *Dao) {
		var (
			query = "SELECT aid FROM ugc_video WHERE deleted = 0 AND submit = 1 AND cid < 12780000 LIMIT 1"
			aid   int64
		)
		if err := d.DB.QueryRow(ctx, query).Scan(&aid); err != nil {
			fmt.Println("No to submit data")
			return
		}
		fmt.Println(aid)
		res, err := d.ParseVideos(ctx, int64(aid), 2)
		convey.So(err, convey.ShouldBeNil)
		convey.So(len(res), convey.ShouldBeGreaterThan, 0)
		str, _ := json.Marshal(res)
		fmt.Println(string(str))
		for _, v := range res {
			convey.So(len(v), convey.ShouldBeGreaterThan, 0)
		}
	}))
}

func TestDao_ShouldAudit(t *testing.T) {
	convey.Convey("TestDao_ShouldAudit", t, WithDao(func(d *Dao) {
		var (
			query = "SELECT aid FROM ugc_video v WHERE v.submit = 1 " + _videoCond + " LIMIT 1"
			aid   int64
		)
		fmt.Println(query)
		d.DB.QueryRow(ctx, fmt.Sprintf(query, d.criCID)).Scan(&aid)
		if aid == 0 {
			fmt.Println("db empty")
			return
		}
		res, err := d.ShouldAudit(ctx, aid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldBeTrue)
		fmt.Println(aid, " ", res)
	}))
}

func TestDao_VideoSubmit(t *testing.T) {
	convey.Convey("TestDao_VideoSubmit", t, WithDao(func(d *Dao) {
		var (
			query = "SELECT aid FROM ugc_video v WHERE v.submit = 0 " + _videoCond + " LIMIT 1"
			aid   int64
		)
		fmt.Println(query)
		d.DB.QueryRow(ctx, fmt.Sprintf(query, d.criCID)).Scan(&aid)
		if aid == 0 {
			fmt.Println("db empty")
			return
		}
		res, err := d.VideoSubmit(ctx, aid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(res, convey.ShouldNotBeNil)
		fmt.Println(aid, " ", res)
	}))
}
