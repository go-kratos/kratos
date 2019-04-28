package dao

import (
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model"

	"context"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_ModIntervs(t *testing.T) {
	Convey("TestDao_ModIntervs", t, WithDao(func(d *Dao) {
		var modID int
		if err := d.db.QueryRow(ctx,
			"SELECT module_id FROM tv_rank "+
				"WHERE is_deleted = 0 AND module_id != 0 AND category = 5 AND type = 1"+
				"LIMIT 1").Scan(&modID); err != nil || modID == 0 {
			return
		}
		sids, err := d.ModIntervs(ctx, modID, 10)
		fmt.Println(sids)
		So(err, ShouldBeNil)
	}))
}

func TestDao_ZoneIntervs(t *testing.T) {
	Convey("TestDao_ZoneIntervs", t, WithDao(func(d *Dao) {
		// home page
		resp, err := d.ZoneIntervs(ctx, &model.ReqZoneInterv{
			RankType: 0,
			Category: 1,
			Limit:    50,
		})
		for _, v := range resp.Ranks {
			fmt.Println(v)
		}
		So(err, ShouldBeNil)
		So(len(resp.Ranks), ShouldBeGreaterThan, 0)
	}))
}

func TestDao_AllIntervs(t *testing.T) {
	Convey("TestDao_AllIntervs", t, WithDao(func(d *Dao) {
		var countUGC, countPGC int64
		d.db.QueryRow(ctx, "SELECT COUNT(1) FROM tv_rank WHERE is_deleted =0 AND cont_type = 2 ").Scan(&countUGC)
		d.db.QueryRow(ctx, "SELECT COUNT(1) FROM tv_rank WHERE is_deleted =0 AND cont_type != 2 ").Scan(&countPGC)
		sids, aids, err := d.AllIntervs(ctx)
		So(err, ShouldBeNil)
		if countUGC > 0 {
			So(len(aids), ShouldBeGreaterThan, 0)
			fmt.Println(aids)
		} else {
			fmt.Println("empty ugc rank")
		}
		if countPGC > 0 {
			So(len(sids), ShouldBeGreaterThan, 0)
			fmt.Println(sids)
		} else {
			fmt.Println("empty pgc ranks")
		}
	}))
}

func TestDao_RmInterv(t *testing.T) {
	Convey("TestDao_RmInterv", t, WithDao(func(d *Dao) {
		var aid, sid int64
		d.db.QueryRow(ctx, "SELECT cont_id FROM tv_rank WHERE is_deleted =0 AND cont_type = 2 LIMIT 1").Scan(&aid)
		d.db.QueryRow(ctx, "SELECT cont_id FROM tv_rank WHERE is_deleted =0 AND cont_type != 2 LIMIT 1 ").Scan(&sid)
		if aid > 0 {
			err := d.RmInterv(context.Background(), []int64{aid}, []int64{})
			So(err, ShouldBeNil)
			d.db.Exec(ctx, "UPDATE tv_rank SET is_deleted = 0 WHERE cont_id = ?", aid)
		} else {
			fmt.Println("empty ugc rank")
		}
		if sid > 0 {
			err := d.RmInterv(context.Background(), []int64{}, []int64{sid})
			So(err, ShouldBeNil)
			d.db.Exec(ctx, "UPDATE tv_rank SET is_deleted = 0 WHERE cont_id = ?", sid)
		} else {
			fmt.Println("empty pgc ranks")
		}
	}))
}

func TestDao_IdxIntervs(t *testing.T) {
	Convey("TestDao_IdxIntervs", t, WithDao(func(d *Dao) {
		res, err := d.IdxIntervs(ctx)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeNil)
		qq, _ := json.Marshal(res)
		fmt.Println(string(qq))
	}))
}
