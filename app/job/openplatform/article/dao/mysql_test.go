package dao

import (
	"context"
	"testing"

	artmdl "go-common/app/interface/openplatform/article/model"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_Update(t *testing.T) {
	var (
		c    = context.TODO()
		cnt1 = int64(1)
		cnt2 = int64(2)
		st1  = &artmdl.StatMsg{
			Aid:      888,
			View:     &cnt1,
			Favorite: &cnt1,
			Like:     &cnt1,
			Dislike:  &cnt1,
			Reply:    &cnt1,
			Share:    &cnt1,
		}
		st2 = &artmdl.StatMsg{
			Aid:      888,
			View:     &cnt2,
			Favorite: &cnt2,
			Like:     &cnt2,
			Dislike:  &cnt2,
			Reply:    &cnt2,
			Share:    &cnt2,
		}
	)
	Convey("update stats", t, WithDao(func(d *Dao) {
		rows, err := d.Update(c, st1)
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)

		Convey("get st1", func() {
			stat, err1 := d.Stat(c, 888)
			So(err1, ShouldBeNil)
			So(stat, ShouldResemble, st1)
		})

		rows, err = d.Update(c, st2)
		So(err, ShouldBeNil)
		So(rows, ShouldBeGreaterThan, 0)

		Convey("get st2", func() {
			stat, err1 := d.Stat(c, 888)
			So(err1, ShouldBeNil)
			So(stat, ShouldResemble, st2)
		})
	}))
}

func Test_GameList(t *testing.T) {
	Convey("work", t, WithDao(func(d *Dao) {
		mids, err := d.GameList(context.Background())
		So(err, ShouldBeNil)
		So(mids, ShouldNotBeEmpty)
	}))
}

func Test_NewestArtIDByCategory(t *testing.T) {
	var _dataCategory = int64(6)
	Convey("get data", t, WithDao(func(d *Dao) {
		res, err := d.NewestArtIDByCategory(context.TODO(), []int64{_dataCategory}, 100)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
	Convey("no data", t, WithDao(func(d *Dao) {
		res, err := d.NewestArtIDByCategory(context.TODO(), []int64{1000}, 100)
		So(err, ShouldBeNil)
		So(res, ShouldBeEmpty)
	}))
}

func Test_NewestArtIDs(t *testing.T) {
	Convey("get data", t, WithDao(func(d *Dao) {
		res, err := d.NewestArtIDs(context.TODO(), 100)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}

func Test_SearchArts(t *testing.T) {
	Convey("should get data", t, WithDao(func(d *Dao) {
		_searchInterval = 24 * 3600 * 365
		res, err := d.SearchArts(context.TODO(), 0)
		So(err, ShouldBeNil)
		So(res, ShouldNotBeEmpty)
	}))
}
