package dao

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSettlePeriod(t *testing.T) {
	Convey("SettlePeriod", t, func() {
		_, err := d.SettlePeriod(ctx, 2)
		So(err, ShouldBeNil)
	})
}

func TestHitSettlePeriod(t *testing.T) {
	Convey("HitSettlePeriod", t, func() {
		_, err := d.HitSettlePeriod(ctx, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateSettle(t *testing.T) {
	Convey("UpdateSettle", t, func() {
		err := d.UpdateSettle(ctx, 1, 1, 10, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestUpdateCoinCount(t *testing.T) {
	Convey("UpdateCoinCount", t, func() {
		err := d.UpdateCoinCount(ctx, 1, 1, 10, 1, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestEvery10000(t *testing.T) {
	Convey("Every10000", t, func() {
		_, _, err := d.Every10000(ctx, 1, 1)
		So(err, ShouldBeNil)
	})
}

func TestUpsertSettle(t *testing.T) {
	Convey("UpsertSettle", t, func() {
		err := d.UpsertSettle(ctx, 1, 1, 1, 1, 1, time.Now())
		So(err, ShouldBeNil)
	})
}

func TestTotalCoins(t *testing.T) {
	Convey("TotalCoins", t, func() {
		_, err := d.TotalCoins(ctx, 1, time.Now(), time.Now())
		So(err, ShouldBeNil)
	})
}
