package cms

import (
	"fmt"
	"testing"

	"go-common/app/interface/main/tv/model"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDao_SetSnCMSCache(t *testing.T) {
	Convey("TestDao_SetSnCMSCache Test", t, WithDao(func(d *Dao) {
		err := d.SetSnCMSCache(ctx, &model.SeasonCMS{
			SeasonID: 6462,
			Cover:    "Test_cover1",
			Title:    "Test_title1",
			Desc:     "Test_desc1",
		})
		So(err, ShouldBeNil)
	}))
}

func TestDao_GetSnCMSCache(t *testing.T) {
	Convey("TestDao_GetSnCMSCache Test", t, WithDao(func(d *Dao) {
		sids, errPick := pickIDs(d.db, _pickSids)
		if errPick != nil || len(sids) == 0 {
			fmt.Println("Empty sids ", errPick)
			return
		}
		sid := sids[0]
		d.LoadSnCMS(ctx, sid)
		season, err := d.GetSnCMSCache(ctx, sid)
		So(err, ShouldBeNil)
		So(season, ShouldNotBeNil)
		fmt.Println(*season)
	}))
}

func TestDao_GetEpCMSCache(t *testing.T) {
	Convey("TestDao_GetEpCMSCache Test", t, WithDao(func(d *Dao) {
		sids, errPick := pickIDs(d.db, _pickEpids)
		if errPick != nil || len(sids) == 0 {
			fmt.Println("Empty sids ", errPick)
			return
		}
		epid := sids[1]
		ep, err := d.GetEpCMSCache(ctx, epid)
		So(err, ShouldBeNil)
		So(ep, ShouldNotBeNil)
		fmt.Println(*ep)
	}))
}

func TestDao_SnCMSCacheKey(t *testing.T) {
	Convey("TestDao_SnCMSCacheKey Test", t, WithDao(func(d *Dao) {
		key := snCMSCacheKey(177)
		So(key, ShouldNotBeBlank)
		fmt.Println(key)
	}))
}

func TestDao_EPCMSCacheKey(t *testing.T) {
	Convey("TestDao_EPCMSCacheKey Test", t, WithDao(func(d *Dao) {
		key := epCMSCacheKey(1)
		So(key, ShouldNotBeBlank)
		fmt.Println(key)
	}))
}

func TestDao_ArcCMSCacheKey(t *testing.T) {
	Convey("TestDao_ArcCMSCacheKey Test", t, WithDao(func(d *Dao) {
		key := d.ArcCMSCacheKey(177)
		So(key, ShouldNotBeBlank)
		fmt.Println(key)
	}))
}
