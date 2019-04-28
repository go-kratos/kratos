package dao

import (
	"context"
	"testing"
	"time"

	"go-common/app/interface/main/answer/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestAnswerQidListKey(t *testing.T) {
	convey.Convey("answerQidListKey", t, func() {
		key := answerQidListKey(7593623, model.BaseExtraPassQ)
		convey.So(key, convey.ShouldNotBeNil)
	})
}

func TestDaopingMC(t *testing.T) {
	convey.Convey("pingMC", t, func() {
		err := d.pingMC(context.Background())
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSetExpireCache(t *testing.T) {
	at := &model.AnswerTime{Etimes: 1}
	convey.Convey("SetExpireCache", t, func() {
		err := d.SetExpireCache(context.Background(), mid, at)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("ExpireCache", t, func() {
		at, err := d.ExpireCache(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(at, convey.ShouldNotBeNil)
	})
}

func TestDaoDelExpireCache(t *testing.T) {
	convey.Convey("DelExpireCache", t, func() {
		err := d.DelExpireCache(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSetHistoryCache(t *testing.T) {
	his := &model.AnswerHistory{Hid: time.Now().Unix()}
	convey.Convey("SetHistoryCache", t, func() {
		err := d.SetHistoryCache(context.Background(), mid, his)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("HistoryCache", t, func() {
		ah, err := d.HistoryCache(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ah, convey.ShouldNotBeNil)
	})
}

func TestDaoSetIdsCache(t *testing.T) {
	convey.Convey("SetIdsCache", t, func() {
		err := d.SetIdsCache(context.Background(), mid, []int64{1, 2, 3}, 0)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("IdsCache", t, func() {
		ids, err := d.IdsCache(context.Background(), mid, 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ids, convey.ShouldNotBeNil)
	})
	convey.Convey("DelHistoryCache", t, func() {
		err := d.DelHistoryCache(context.Background(), mid)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoDelIdsCache(t *testing.T) {
	convey.Convey("DelIdsCache", t, func() {
		err := d.DelIdsCache(context.Background(), 0, 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoSetBlockCache(t *testing.T) {
	convey.Convey("SetBlockCache", t, func() {
		err := d.SetBlockCache(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
	})
}

func TestDaoCheckBlockCache(t *testing.T) {
	convey.Convey("CheckBlockCache", t, func() {
		exist, err := d.CheckBlockCache(context.Background(), 0)
		convey.So(err, convey.ShouldBeNil)
		convey.So(exist, convey.ShouldNotBeNil)
	})
}

func TestDaoHidCache(t *testing.T) {
	hid := time.Now().Unix()
	his := &model.AnswerHistory{Hid: hid}
	convey.Convey("SetHidCache", t, func() {
		err := d.SetHidCache(context.Background(), his)
		convey.So(err, convey.ShouldBeNil)
	})
	convey.Convey("HidCache", t, func() {
		ah, err := d.HidCache(context.Background(), hid)
		convey.So(err, convey.ShouldBeNil)
		convey.So(ah, convey.ShouldNotBeNil)
	})
}
