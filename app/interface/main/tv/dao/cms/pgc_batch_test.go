package cms

import (
	"context"
	"fmt"
	"go-common/app/interface/main/tv/model"
	"testing"

	"github.com/smartystreets/goconvey/convey"
)

func TestCmssnCMSCacheKey(t *testing.T) {
	var (
		sid = int64(0)
	)
	convey.Convey("snCMSCacheKey", t, func(c convey.C) {
		p1 := snCMSCacheKey(sid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsepCMSCacheKey(t *testing.T) {
	var (
		epid = int64(0)
	)
	convey.Convey("epCMSCacheKey", t, func(c convey.C) {
		p1 := epCMSCacheKey(epid)
		c.Convey("Then p1 should not be nil.", func(c convey.C) {
			c.So(p1, convey.ShouldNotBeNil)
		})
	})
}

func TestCmskeysTreat(t *testing.T) {
	var (
		ids     = []int64{1, 2, 3}
		keyFunc = snCMSCacheKey
	)
	convey.Convey("keysTreat", t, func(c convey.C) {
		idmap, allKeys := keysTreat(ids, keyFunc)
		c.Convey("Then idmap,allKeys should not be nil.", func(c convey.C) {
			c.So(allKeys, convey.ShouldNotBeNil)
			c.So(idmap, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsmissedTreat(t *testing.T) {
	var (
		idmap     map[string]int64
		lenCached = int(0)
	)
	convey.Convey("missedTreat", t, func(c convey.C) {
		missed := missedTreat(idmap, lenCached)
		c.Convey("Then missed should not be nil.", func(c convey.C) {
			c.So(missed, convey.ShouldNotBeNil)
		})
	})
}

func TestCmsEpMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
	)
	convey.Convey("EpMetaCache", t, func(c convey.C) {
		c.Convey("Empty Input", func(c convey.C) {
			cached, missed, err := d.EpMetaCache(ctx, []int64{})
			c.So(err, convey.ShouldBeNil)
			c.So(len(missed), convey.ShouldBeZeroValue)
			c.So(len(cached), convey.ShouldBeZeroValue)
		})
		c.Convey("Normal Situation", func(c convey.C) {
			epids, err := pickIDs(d.db, _pickEpids)
			if err != nil || len(epids) == 0 {
				fmt.Println("empty epids")
				return
			}
			cached, missed, err := d.EpMetaCache(ctx, epids)
			c.So(err, convey.ShouldBeNil)
			c.So(len(missed)+len(cached), convey.ShouldNotEqual, 0)
		})
	})
}

func TestCmsAddSeasonMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
		vs  = &model.SeasonCMS{}
	)
	convey.Convey("AddSeasonMetaCache", t, func(c convey.C) {
		err := d.AddSeasonMetaCache(ctx, vs)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}

func TestCmsAddEpMetaCache(t *testing.T) {
	var (
		ctx = context.Background()
		vs  = &model.EpCMS{}
	)
	convey.Convey("AddEpMetaCache", t, func(c convey.C) {
		err := d.AddEpMetaCache(ctx, vs)
		c.Convey("Then err should be nil.", func(c convey.C) {
			c.So(err, convey.ShouldBeNil)
		})
	})
}
