package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	"go-common/app/service/main/feed/model"
	feed "go-common/app/service/main/feed/model"
	xtime "go-common/library/time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FeedValue(t *testing.T) {
	var (
		arc  = api.Arc{Aid: 1}
		arc2 = api.Arc{Aid: 2}
		f    feed.Feed
	)

	Convey("with fold avs", t, func() {
		f = feed.Feed{ID: 100, Fold: []*api.Arc{&arc, &arc2}}
		So(appFeedValue(&f), ShouldEqual, "0,100,1,2")
	})
	Convey("without fold avs", t, func() {
		f = feed.Feed{ID: 1}
		So(appFeedValue(&f), ShouldEqual, "0,1")
	})

	Convey("bangumi", t, func() {
		f = feed.Feed{ID: 100, Type: feed.BangumiType}
		So(appFeedValue(&f), ShouldEqual, "1,100")
	})
}

func Test_RecoverFeed(t *testing.T) {
	var (
		arc  = api.Arc{Aid: 1}
		arc2 = api.Arc{Aid: 2}
		b    = feed.Feed{ID: 100, Type: feed.BangumiType}
		f    feed.Feed
	)

	Convey("bangumi", t, func() {
		r, err := recoverFeed("1,100")
		So(r, ShouldResemble, &b)
		So(err, ShouldBeNil)
	})

	Convey("with fold avs", t, func() {
		f = feed.Feed{ID: 100, Fold: []*api.Arc{&arc, &arc2}}
		r, err := recoverFeed("0,100,1,2")
		So(r, ShouldResemble, &f)
		So(err, ShouldBeNil)
	})

	Convey("without fold avs", t, func() {
		f = feed.Feed{ID: 100}
		r, err := recoverFeed("0,100")
		So(r, ShouldResemble, &f)
		So(err, ShouldBeNil)
	})
}

func Test_pingRedis(t *testing.T) {
	Convey("ping redis", t, func() {
		So(d.pingRedis(context.TODO()), ShouldBeNil)
	})
}

func Test_LastAccessCache(t *testing.T) {
	var (
		mid = int64(1)
		ts  = int64(100)
		err error
	)
	Convey("add cache", t, func() {
		err = d.AddLastAccessCache(context.TODO(), model.TypeApp, mid, ts)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			t1, err := d.LastAccessCache(context.TODO(), model.TypeApp, mid)
			So(t1, ShouldEqual, t1)
			So(err, ShouldBeNil)
		})
	})
}

func Test_FeedCache(t *testing.T) {
	var (
		mid     = int64(1)
		now     = time.Now().Unix()
		err     error
		a1      = api.Arc{Aid: 1, PubDate: xtime.Time(now)}
		a2      = api.Arc{Aid: 2, PubDate: xtime.Time(now - 1000)}
		a3      = api.Arc{Aid: 3}
		bangumi = feed.Bangumi{SeasonID: 100}
		f       = feed.Feed{ID: 1, Archive: &a1, PubDate: a1.PubDate, Fold: []*api.Arc{&a3}}
		f1      = feed.Feed{ID: 2, Archive: &a2, PubDate: a2.PubDate}
		b       = feed.Feed{ID: 100, Type: feed.BangumiType, Bangumi: &bangumi}
		feeds   = []*feed.Feed{&f, &f1, &b}
	)
	Convey("add cache", t, func() {
		for name, client := range map[string]int{"app": model.TypeApp, "web": model.TypeWeb} {
			err = d.AddFeedCache(context.TODO(), client, mid, feeds)
			So(err, ShouldBeNil)
			Convey(name+"get cache", func() {
				res, bids, err := d.FeedCache(context.TODO(), client, mid, 0, 0)
				So(res, ShouldResemble, []*feed.Feed{{ID: f.ID, Fold: []*api.Arc{&a3}, PubDate: f.PubDate}})
				So(bids, ShouldBeEmpty)
				So(err, ShouldBeNil)
			})

			Convey(name+"get cache when end > length", func() {
				res, bids, err := d.FeedCache(context.TODO(), client, mid, 0, 10)
				So(res, ShouldResemble, []*feed.Feed{
					{ID: a1.Aid, Fold: []*api.Arc{&a3}, PubDate: a1.PubDate},
					{ID: a2.Aid, PubDate: a2.PubDate},
					{ID: 100, Type: feed.BangumiType},
				})
				So(bids, ShouldResemble, []int64{100})
				So(err, ShouldBeNil)
			})

			Convey(name+"expire cache", func() {
				ok, err := d.ExpireFeedCache(context.TODO(), client, mid)
				So(ok, ShouldEqual, true)
				So(err, ShouldBeNil)
			})

			Convey(name+"purge cache", func() {
				err := d.PurgeFeedCache(context.TODO(), client, mid)
				So(err, ShouldBeNil)
			})
		}
	})
}

func Test_UppersCache(t *testing.T) {
	var (
		mid  = int64(1)
		mid2 = int64(2)
		now  = time.Now().Unix()
		err  error
		a1   = archive.AidPubTime{Aid: 1, PubDate: xtime.Time(now), Copyright: 1}
		a2   = archive.AidPubTime{Aid: 2, PubDate: xtime.Time(now - 1), Copyright: 0}
		a3   = archive.AidPubTime{Aid: 3, PubDate: xtime.Time(now - 2), Copyright: 0}
	)
	Convey("add cache", t, func() {
		err = d.AddUpperCaches(context.TODO(), map[int64][]*archive.AidPubTime{mid: {&a1, &a2}, mid2: {&a3}})
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			_, err := d.UppersCaches(context.TODO(), []int64{mid, mid2}, 0, 2)
			So(err, ShouldBeNil)
			// So(res, ShouldResemble, map[int64][]*archive.AidPubTime{mid: {&a1, &a2}, mid2: {&a3}})
		})
		Convey("expire cache", func() {
			res, err := d.ExpireUppersCache(context.TODO(), []int64{mid})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, map[int64]bool{mid: true})
		})
		Convey("get expired cache", func() {
			d.redisExpireUpper = 0
			res, err := d.ExpireUppersCache(context.TODO(), []int64{mid})
			So(err, ShouldBeNil)
			So(res, ShouldResemble, map[int64]bool{mid: false})

			_, err = d.UppersCaches(context.TODO(), []int64{mid}, 0, 2)
			So(err, ShouldBeNil)
			// So(nres, ShouldResemble, map[int64][]*archive.AidPubTime{mid: {&a1, &a2}})
		})

		Convey("purge cache", func() {
			err := d.DelUpperCache(context.TODO(), mid, a1.Aid)
			So(err, ShouldBeNil)
		})
	})
}

func Test_ArchiveFeedCache(t *testing.T) {
	var (
		mid = int64(1)
		now = time.Now().Unix()
		err error
		a1  = api.Arc{Aid: 1, PubDate: xtime.Time(now), Author: api.Author{Mid: mid}}
		a2  = api.Arc{Aid: 2, PubDate: xtime.Time(now - 1), Author: api.Author{Mid: mid}}
		a3  = api.Arc{Aid: 3, PubDate: xtime.Time(now - 2), Author: api.Author{Mid: mid}}
		f1  = feed.Feed{ID: a1.Aid, Archive: &a1, PubDate: a1.PubDate, Fold: []*api.Arc{&a3}}
		f2  = feed.Feed{ID: a2.Aid, Archive: &a2, PubDate: a2.PubDate}
		fs  = []*feed.Feed{&f1, &f2}
	)
	Convey("add cache", t, func() {
		err = d.AddArchiveFeedCache(context.TODO(), mid, fs)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			as, err := d.ArchiveFeedCache(context.TODO(), mid, 0, 2)
			So(as, ShouldResemble, []*feed.Feed{
				{ID: a1.Aid, PubDate: a1.PubDate, Fold: []*api.Arc{{Aid: 3}}},
				{ID: a2.Aid, PubDate: a2.PubDate},
			})
			So(err, ShouldBeNil)
		})
	})
}

func Test_BangumiFeedCache(t *testing.T) {
	var (
		mid = int64(1)
		err error
		b1  = feed.Bangumi{SeasonID: 100}
		b2  = feed.Bangumi{SeasonID: 200}
		f1  = feed.Feed{ID: b1.SeasonID, Type: feed.BangumiType, Bangumi: &b1}
		f2  = feed.Feed{ID: b2.SeasonID, Type: feed.BangumiType, Bangumi: &b2}
		fs  = []*feed.Feed{&f1, &f2}
	)
	Convey("add cache", t, func() {
		err = d.AddBangumiFeedCache(context.TODO(), mid, fs)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			res, err := d.BangumiFeedCache(context.TODO(), mid, 0, 2)
			So(res, ShouldResemble, []int64{b2.SeasonID, b1.SeasonID})
			So(err, ShouldBeNil)
		})
	})
}

func Test_UnreadCountCache(t *testing.T) {
	var (
		mid   = int64(1)
		count = 100
		err   error
	)
	Convey("add cache", t, func() {
		err = d.AddUnreadCountCache(context.TODO(), model.TypeApp, mid, count)
		So(err, ShouldBeNil)
		Convey("get cache", func() {
			c, err := d.UnreadCountCache(context.TODO(), model.TypeApp, mid)
			So(c, ShouldEqual, count)
			So(err, ShouldBeNil)
		})
		Convey("get wrong cache", func() {
			c, err := d.UnreadCountCache(context.TODO(), model.TypeWeb, mid)
			So(c, ShouldEqual, 0)
			So(err, ShouldBeNil)
		})
	})
}
