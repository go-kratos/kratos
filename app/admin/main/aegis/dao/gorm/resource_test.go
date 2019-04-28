package gorm

import (
	"context"
	"testing"

	"go-common/app/admin/main/aegis/model/resource"

	"github.com/smartystreets/goconvey/convey"
)

func TestGormListHelperForTask(t *testing.T) {
	convey.Convey("ListHelperForTask", t, func(ctx convey.C) {
		var (
			c    = context.Background()
			rids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.ListHelperForTask(c, rids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestGormResourceByOID(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("ResourceByOID", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ResourceByOID(c, "", 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormOidByRID(t *testing.T) {
	var (
		c = context.Background()
	)
	convey.Convey("OidByRID", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.OidByRID(c, 0)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormTxUpdateResult(t *testing.T) {
	convey.Convey("TxUpdateResult", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTx(context.TODO())
			rid    = int64(0)
			res    = map[string]interface{}{"extra1": 1}
			rscRes = &resource.Result{}
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpdateResult(tx, rid, res, rscRes)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormResourceRes(t *testing.T) {
	convey.Convey("ResourceRes", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ResourceRes(context.TODO(), 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormResByOID(t *testing.T) {
	convey.Convey("ResByOID", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.ResByOID(context.TODO(), 0, "")
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormUpdateResource(t *testing.T) {
	convey.Convey("UpdateResource", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.UpdateResource(context.TODO(), 1, "xyz", map[string]interface{}{"mid": 123})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestGormResOIDByID(t *testing.T) {
	convey.Convey("ResOIDByID", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ResOIDByID(context.TODO(), []int64{1, 2})
			t.Logf("res(%+v)", res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormResIDByOID(t *testing.T) {
	convey.Convey("ResIDByOID", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.ResIDByOID(context.TODO(), 1, []string{"1", "2"})
			t.Logf("res(%+v)", res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormResourceHit(t *testing.T) {
	convey.Convey("ResourceHit", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			_, err := d.ResourceHit(context.TODO(), []int64{})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				convey.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormResultByRID(t *testing.T) {
	convey.Convey("ResultByRID", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.ResultByRID(context.TODO(), 50)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestGormRidsByOids(t *testing.T) {
	convey.Convey("RidsByOids", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			d.RidsByOids(cntx, 0, []string{})
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
			})
		})
	})
}

func TestGormTxUpdateResource(t *testing.T) {
	convey.Convey("TxUpdateResource", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTx(context.TODO())
			rid   = int64(0)
			res   = map[string]interface{}{"extra1": 1}
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpdateResource(tx, rid, res)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormTxUpdateState(t *testing.T) {
	convey.Convey("TxUpdateState", t, func(ctx convey.C) {
		var (
			tx, _ = d.BeginTx(context.TODO())
			rids  = []int64{0}
		)
		defer tx.Commit()
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxUpdateState(tx, rids, 0)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormMetaByRID(t *testing.T) {
	convey.Convey("MetaByRID", t, func(ctx convey.C) {
		var (
			rids = []int64{0}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			_, err := d.MetaByRID(cntx, rids)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestGormUpsertByOIDs(t *testing.T) {
	convey.Convey("UpsertByOIDs", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.UpsertByOIDs(context.TODO(), 1, []string{"oid000001", "oid000002"})
			for _, item := range res {
				t.Logf("res(%+v)", item)
			}
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
