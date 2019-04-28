package service

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/archive/api"
	"go-common/app/service/main/archive/model/archive"
	feed "go-common/app/service/main/feed/model"
	xtime "go-common/library/time"

	"strconv"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_fold(t *testing.T) {
	Convey("fold archives", t, WithBlankService(func(svf *Service) {
		var t, t2, t3 time.Time
		t, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 00:00:00")
		t2, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 03:00:00")
		t3, _ = time.Parse("2006-01-02 15:04:05", "2017-03-01 05:00:00")
		arc := archive.AidPubTime{Aid: 1, Copyright: 0, PubDate: xtime.Time(t.Unix())}
		arc1 := archive.AidPubTime{Aid: 1, Copyright: 1, PubDate: xtime.Time(t.Unix())}
		arc2 := archive.AidPubTime{Aid: 2, Copyright: 1, PubDate: xtime.Time(t2.Unix())}
		arc3 := archive.AidPubTime{Aid: 3, Copyright: 1, PubDate: xtime.Time(t3.Unix())}

		Convey("fold reprinted archives", func() {
			arcs := []*archive.AidPubTime{&arc2, &arc, &arc3}
			res := []*feed.Feed{
				&feed.Feed{ID: arc3.Aid, PubDate: arc3.PubDate},
				&feed.Feed{ID: arc2.Aid, Fold: []*api.Arc{&api.Arc{Aid: arc.Aid}}, PubDate: arc2.PubDate},
			}
			So(svf.fold(arcs), ShouldResemble, res)
		})

		Convey("not fold original archives", func() {
			arcs := []*archive.AidPubTime{&arc2, &arc1, &arc3}
			res := []*feed.Feed{
				&feed.Feed{ID: arc3.Aid, PubDate: arc3.PubDate},
				&feed.Feed{ID: arc2.Aid, PubDate: arc2.PubDate},
				&feed.Feed{ID: arc1.Aid, PubDate: arc1.PubDate},
			}
			So(svf.fold(arcs), ShouldResemble, res)
		})
	}))
}

func Test_Feed(t *testing.T) {
	for name, client := range map[string]bool{"app": true, "web": false} {
		for _, mid := range []int64{_mid, _bangumiMid} {
			midStr := strconv.FormatInt(mid, 10)
			Convey(name+midStr+" with fold return feed", t, WithService(t, func(svf *Service) {
				res, err := svf.Feed(context.TODO(), client, mid, 1, 2, _ip)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeEmpty)
				Convey(name+"return feed for page 2", func() {
					time.Sleep(time.Millisecond * 300) // wait cache ready
					res, err := svf.Feed(context.TODO(), client, mid, 2, 2, _ip)
					So(err, ShouldBeNil)
					So(res, ShouldNotBeEmpty)
				})
			}))

			Convey(name+midStr+" without fold return feed", t, WithService(t, func(svf *Service) {
				res, err := svf.Feed(context.TODO(), client, mid, 1, 2, _ip)
				So(err, ShouldBeNil)
				So(res, ShouldNotBeEmpty)
				Convey(name+"return feed for page 2", func() {
					time.Sleep(time.Millisecond * 300) // wait cache ready
					res, err := svf.Feed(context.TODO(), client, mid, 2, 2, _ip)
					So(err, ShouldBeNil)
					So(res, ShouldNotBeEmpty)
				})
			}))
		}

		Convey(name+"user don't have attention ups and bangumis", t, WithService(t, func(svf *Service) {
			midStr := strconv.FormatInt(_blankMid, 10)
			Convey("with fold return feed", func() {
				res, err := svf.Feed(context.TODO(), client, _blankMid, 1, 2, _ip)
				So(err, ShouldBeNil)
				So(res, ShouldBeEmpty)
				Convey(name+"return feed for page 2", func() {
					time.Sleep(time.Millisecond * 300) // wait cache ready
					res, err := svf.Feed(context.TODO(), client, _blankMid, 2, 2, _ip)
					So(err, ShouldBeNil)
					So(res, ShouldBeEmpty)
				})
			})

			Convey(name+midStr+" without fold return feed", func() {
				res, err := svf.Feed(context.TODO(), client, _blankMid, 1, 2, _ip)
				So(err, ShouldBeNil)
				So(res, ShouldBeEmpty)
				Convey(name+"return feed for page 2", func() {
					time.Sleep(time.Millisecond * 300) // wait cache ready
					res, err := svf.Feed(context.TODO(), client, _blankMid, 2, 2, _ip)
					So(err, ShouldBeNil)
					So(res, ShouldBeEmpty)
				})
			})
		}))
	}
}

func Test_PurgeFeedCache(t *testing.T) {
	Convey("should return without err", t, WithService(t, func(svf *Service) {
		err := svf.PurgeFeedCache(context.TODO(), _mid, _ip)
		So(err, ShouldBeNil)
	}))
}

func Test_bangumiFeed(t *testing.T) {
	Convey("should return without err", t, WithService(t, func(svf *Service) {
		feeds, err := svf.bangumiFeedFromSeason(context.TODO(), []int64{_seasonID}, _ip)
		So(err, ShouldBeNil)
		So(feeds, ShouldNotBeEmpty)
	}))
}

func Test_fillArchiveFeeds(t *testing.T) {
	Convey("fill feeds", t, WithService(t, func(svf *Service) {
		f1 := &feed.Feed{ID: _arc1.Aid, PubDate: _arc1.PubDate, Fold: []*api.Arc{&api.Arc{Aid: _arc2.Aid}}}
		f2 := &feed.Feed{ID: _arc1.Aid, PubDate: _arc1.PubDate, Archive: _arc1, Fold: []*api.Arc{_arc2}}
		bangumi := &feed.Feed{ID: 1, Type: feed.BangumiType}
		fs := []*feed.Feed{f1, bangumi}
		expt := []*feed.Feed{f2, bangumi}
		feeds, _ := svf.fillArchiveFeeds(context.TODO(), fs, _ip)
		So(feeds, ShouldResemble, expt)
	}))
}

func Test_replaceFeeds(t *testing.T) {
	Convey("replace feeds", t, WithBlankService(func(svf *Service) {
		arc1 := api.Arc{Aid: 1, Copyright: 1}
		arc2 := api.Arc{Aid: 2, Copyright: 1}
		arc4 := api.Arc{Aid: 4, Copyright: 1}
		bangumi := feed.Bangumi{SeasonID: 1}
		f1 := &feed.Feed{Archive: &arc1, ID: arc1.Aid}
		f2 := &feed.Feed{Archive: &arc2, ID: arc2.Aid}
		f3 := &feed.Feed{Bangumi: &bangumi, ID: bangumi.SeasonID, Type: feed.BangumiType}
		blankf3 := &feed.Feed{ID: bangumi.SeasonID, Type: feed.BangumiType}
		f4 := &feed.Feed{Archive: &arc4, ID: arc4.Aid}
		blankf4 := &feed.Feed{ID: arc4.Aid}

		Convey("replace bangumi feed", func() {
			resfs := []*feed.Feed{f1, blankf3, f2}
			fs := []*feed.Feed{f3}
			svf.replaceFeeds(resfs, fs)
			So(resfs, ShouldResemble, []*feed.Feed{f1, f3, f2})
		})

		Convey("replace archive feed", func() {
			resfs := []*feed.Feed{f1, blankf4, f2}
			fs := []*feed.Feed{f4}
			svf.replaceFeeds(resfs, fs)
			So(resfs, ShouldResemble, []*feed.Feed{f1, f4, f2})
		})

		Convey("blank feed", func() {
			resfs := []*feed.Feed{f1, f2}
			fs := []*feed.Feed{}
			svf.replaceFeeds(resfs, fs)
			So(resfs, ShouldResemble, []*feed.Feed{f1, f2})
		})
	}))
}
