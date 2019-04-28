package dao

import (
	"context"
	"math/rand"
	"testing"

	"go-common/app/service/main/share/model"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoShares(t *testing.T) {
	var (
		c    = context.TODO()
		oids = []int64{1, 2}
		tp   = int(0)
	)
	convey.Convey("Shares", t, func(ctx convey.C) {
		shares, err := d.Shares(c, oids, tp)
		ctx.Convey("Then err should be nil.shares should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(shares, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoShareCount(t *testing.T) {
	var (
		c   = context.TODO()
		oid = int64(rand.Intn(100000000))
		tp  = int(2)
	)
	convey.Convey("ShareCount", t, func(ctx convey.C) {
		count, err := d.ShareCount(c, oid, tp)
		ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(count, convey.ShouldNotBeNil)
		})
	})
}

func TestDaoAdd(t *testing.T) {
	convey.Convey("Add", t, func(ctx convey.C) {
		oid := rand.Intn(1000000)
		mid := rand.Intn(1000000)
		p := &model.ShareParams{
			OID: int64(oid),
			MID: int64(mid),
			TP:  int(3),
			IP:  "",
		}
		shared, err := d.Add(context.Background(), p)
		if err == ecode.ShareAlreadyAdd {
			err = nil
		}
		ctx.Convey("Then err should be nil.shared should not be nil.", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			ctx.So(shared, convey.ShouldNotBeNil)
		})
	})
}
