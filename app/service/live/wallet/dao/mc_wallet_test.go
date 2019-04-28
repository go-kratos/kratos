package dao

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"

	"go-common/app/service/live/wallet/model"
	mc "go-common/library/cache/memcache"
	"math/rand"
	"time"
)

func TestDao_WalletCache(t *testing.T) {
	once.Do(startService)
	Convey("Get Nil Cache", t, func() {

		r = rand.New(rand.NewSource(time.Now().UnixNano()))

		uid := r.Int63n(1000000)

		_, err := d.WalletCache(ctx, uid)

		So(err, ShouldEqual, mc.ErrNotFound)

	})

	Convey("Set And Get And Del", t, func() {

		r = rand.New(rand.NewSource(time.Now().UnixNano()))

		uid := r.Int63n(1000000)

		detail := &model.Detail{Uid: uid, Gold: 1, IapGold: 1, Silver: 1, GoldPayCnt: 1, GoldRechargeCnt: 1, SilverPayCnt: 1}

		mcDetail := &model.McDetail{Detail: detail, Exist: true, Version: d.CacheVersion(ctx)}
		err := d.SetWalletCache(ctx, mcDetail, d.cacheExpire)
		So(err, ShouldBeNil)

		nd, err := d.WalletCache(ctx, uid)
		So(err, ShouldBeNil)
		So(nd.Exist, ShouldEqual, true)
		So(nd.Detail.Gold, ShouldEqual, 1)
		So(nd.Detail.IapGold, ShouldEqual, 1)
		So(nd.Detail.Silver, ShouldEqual, 1)
		So(nd.Detail.GoldPayCnt, ShouldEqual, 1)
		So(nd.Detail.GoldRechargeCnt, ShouldEqual, 1)
		So(nd.Detail.SilverPayCnt, ShouldEqual, 1)
		So(nd.Version, ShouldEqual, d.CacheVersion(ctx))

		err = d.DelWalletCache(ctx, uid)
		So(err, ShouldBeNil)

		nnd, err := d.WalletCache(ctx, uid)
		So(err, ShouldEqual, mc.ErrNotFound)
		So(nnd, ShouldBeNil)

	})

	Convey("version", t, func() {
		Convey("just set", func() {
			uid := int64(1)
			d.DelWalletCache(ctx, uid)
			d.GetDetailByCache(ctx, uid)
			detail, err := d.WalletCache(ctx, uid)
			So(err, ShouldBeNil)
			So(detail.Version, ShouldEqual, 1)
		})
		Convey("old to new", func() {

			uid := int64(1)
			od, err := d.Detail(ctx, uid)
			So(err, ShouldBeNil)
			d.DelWalletCache(ctx, uid)

			detail := &model.Detail{Uid: uid, Gold: od.Gold + 1, IapGold: 1, Silver: 1, GoldPayCnt: 1, GoldRechargeCnt: 1, SilverPayCnt: 1} // fake data

			mcDetail := &model.McDetail{Detail: detail, Exist: true} // old version
			err = d.SetWalletCache(ctx, mcDetail, d.cacheExpire)
			So(err, ShouldBeNil)
			So(mcDetail.Version, ShouldEqual, 0)
			So(mcDetail.Detail.Gold, ShouldEqual, od.Gold+1)
			So(mcDetail.Detail.IapGold, ShouldEqual, 1)
			So(mcDetail.Detail.Silver, ShouldEqual, 1)
			So(mcDetail.Detail.GoldPayCnt, ShouldEqual, 1)

			mcDetail1, err := d.WalletCache(ctx, uid) //
			So(err, ShouldBeNil)
			So(mcDetail1.Detail.Gold, ShouldEqual, od.Gold+1)

			fd, err := d.GetDetailByCache(ctx, uid) // check version if it is old update cache
			So(err, ShouldBeNil)
			So(fd.Gold, ShouldEqual, od.Gold)
			So(fd.Silver, ShouldEqual, od.Silver)
			So(fd.GoldPayCnt, ShouldEqual, od.GoldPayCnt)

			mcDetail2, err := d.WalletCache(ctx, uid)
			So(err, ShouldBeNil)
			So(mcDetail2.Version, ShouldEqual, d.CacheVersion(ctx))
			So(mcDetail2.Detail.Gold, ShouldEqual, od.Gold)

		})
	})
}
