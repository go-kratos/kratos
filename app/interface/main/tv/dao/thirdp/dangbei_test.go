package thirdp

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model/thirdp"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_GetDBeiCount(t *testing.T) {
	Convey("TestDao_GetDBeiCount", t, WithDao(func(d *Dao) {
		count, err := d.GetThirdpCnt(ctx, DBeiUGC)
		So(err, ShouldBeNil)
		So(count, ShouldBeGreaterThan, 0)
		fmt.Println(count)
	}))
}

func TestDao_SetDBeiCount(t *testing.T) {
	Convey("TestDao_SetDBeiCount", t, WithDao(func(d *Dao) {
		err := d.SetThirdpCnt(ctx, 6667, DBeiUGC)
		So(err, ShouldBeNil)
	}))
}

func TestDao_ThirdpCnt(t *testing.T) {
	Convey("TestDao_ThirdpCnt", t, WithDao(func(d *Dao) {
		cntDBeiPGC, err2 := d.ThirdpCnt(ctx, DBeiPGC)
		So(err2, ShouldBeNil)
		So(cntDBeiPGC, ShouldBeGreaterThan, 0)
		fmt.Println(cntDBeiPGC)
		cntDBeiUGC, err := d.ThirdpCnt(ctx, DBeiUGC)
		So(err, ShouldBeNil)
		So(cntDBeiUGC, ShouldBeGreaterThan, 0)
		fmt.Println(cntDBeiUGC)
		cntMangoPGC, err3 := d.ThirdpCnt(ctx, MangoPGC)
		So(err3, ShouldBeNil)
		So(cntMangoPGC, ShouldBeGreaterThan, 0)
		fmt.Println(cntMangoPGC)
		cntMangoUGC, err4 := d.ThirdpCnt(ctx, MangoUGC)
		So(err4, ShouldBeNil)
		So(cntMangoUGC, ShouldBeGreaterThan, 0)
		fmt.Println(cntMangoUGC)
		So(cntDBeiPGC, ShouldBeLessThanOrEqualTo, cntMangoPGC)
		So(cntDBeiUGC, ShouldBeLessThanOrEqualTo, cntMangoUGC)
	}))
}

func TestDao_ThirdpPages(t *testing.T) {
	Convey("TestDao_ThirdpPages", t, WithDao(func(d *Dao) {
		typesDBei := []string{
			DBeiPGC, DBeiUGC,
		}
		for _, v := range typesDBei {
			fmt.Println("--- ", v, " ---")
			req := &thirdp.ReqDBeiPages{
				LastID: 255,
				Ps:     10,
				TypeC:  v,
			}
			sids, myID, err := d.DBeiPages(ctx, req)
			So(err, ShouldBeNil)
			fmt.Println(sids)
			So(myID, ShouldBeGreaterThan, 0)
			fmt.Println(myID)
		}
		typesMango := []string{
			MangoPGC, MangoUGC,
		}
		for _, v := range typesMango {
			fmt.Println("--- ", v, " ---")
			req := &thirdp.ReqDBeiPages{
				LastID: 255,
				Ps:     10,
				TypeC:  v,
			}
			sids, myID, err := d.MangoPages(ctx, req)
			So(err, ShouldBeNil)
			str, _ := json.Marshal(sids)
			fmt.Println(string(str))
			So(myID, ShouldBeGreaterThan, 0)
			fmt.Println(myID)
		}
	}))
}

func TestDao_thirdpOffset(t *testing.T) {
	Convey("TestDao_thirdpOffset", t, WithDao(func(d *Dao) {
		lastMax, err := d.thirdpOffset(context.Background(), 1, 50, MangoPGC)
		So(lastMax, ShouldEqual, 0)
		So(err, ShouldBeNil)
		types := []string{
			DBeiPGC, DBeiUGC, MangoPGC, MangoUGC,
		}
		for _, v := range types {
			lastMaxType, errP := d.thirdpOffset(context.Background(), 2, 50, v)
			So(lastMaxType, ShouldBeGreaterThan, 0)
			So(errP, ShouldBeNil)
			fmt.Println(fmt.Sprintf("[Type - %s] [Max - %d]", v, lastMaxType))
		}
	}))
}

func TestDao_SetPageID(t *testing.T) {
	Convey("TestDao_SetPageID", t, WithDao(func(d *Dao) {
		err := d.SetPageID(ctx, &thirdp.ReqPageID{
			TypeC: DBeiUGC,
			ID:    8088,
			Page:  3,
		})
		So(err, ShouldBeNil)
	}))
}

func TestDao_GetPageID(t *testing.T) {
	Convey("TestDao_GetPageID", t, WithDao(func(d *Dao) {
		biggestID, err := d.getPageID(ctx, 7, DBeiUGC)
		So(err, ShouldBeNil)
		So(biggestID, ShouldBeGreaterThan, 0)
		fmt.Println(biggestID)
	}))
}

func TestDao_LoadPageID(t *testing.T) {
	Convey("TestDao_LoadPageID", t, WithDao(func(d *Dao) {
		types := []string{
			DBeiPGC, DBeiUGC, MangoPGC, MangoUGC,
		}
		for _, v := range types {
			req := &thirdp.ReqDBeiPages{
				Page:  7,
				Ps:    50,
				TypeC: v,
			}
			biggestID, err := d.LoadPageID(ctx, req)
			So(err, ShouldBeNil)
			So(biggestID, ShouldBeGreaterThan, 0)
			fmt.Println("Type: ", v, " Max: ", biggestID)
		}
	}))
}
