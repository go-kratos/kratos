package pendant

import (
	"strconv"
	"testing"
	"time"

	"go-common/app/service/main/usersuit/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestPendantPendantGroupInfo(t *testing.T) {
	convey.Convey("PendantGroupInfo", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PendantGroupInfo(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantGroupByID(t *testing.T) {
	convey.Convey("GroupByID", t, func(ctx convey.C) {
		var (
			gid = int64(4)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GroupByID(c, gid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantGIDRefPID(t *testing.T) {
	convey.Convey("GIDRefPID", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			gidMap, pidMap, err := d.GIDRefPID(c)
			ctx.Convey("Then err should be nil.gidMap,pidMap should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(pidMap, convey.ShouldNotBeNil)
				ctx.So(gidMap, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPendantList(t *testing.T) {
	convey.Convey("PendantList", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PendantList(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPendants(t *testing.T) {
	convey.Convey("Pendants", t, func(ctx convey.C) {
		var (
			pids = []int64{4}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.Pendants(c, pids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPendantInfo(t *testing.T) {
	convey.Convey("PendantInfo", t, func(ctx convey.C) {
		var (
			pid = int64(4)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PendantInfo(c, pid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPendantPrice(t *testing.T) {
	convey.Convey("PendantPrice", t, func(ctx convey.C) {
		var (
			pid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PendantPrice(c, pid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantgetOrderInfoSQL(t *testing.T) {
	convey.Convey("getOrderInfoSQL", t, func(ctx convey.C) {
		var (
			arg = &model.ArgOrderHistory{}
			tp  = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sql, values := d.getOrderInfoSQL(c, arg, tp)
			ctx.Convey("Then sql,values should not be nil.", func(ctx convey.C) {
				ctx.So(values, convey.ShouldNotBeNil)
				ctx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantOrderInfo(t *testing.T) {
	convey.Convey("OrderInfo", t, func(ctx convey.C) {
		var (
			arg = &model.ArgOrderHistory{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, count, err := d.OrderInfo(c, arg)
			ctx.Convey("Then err should be nil.res,count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantOrderInfoByID(t *testing.T) {
	convey.Convey("OrderInfoByID", t, func(ctx convey.C) {
		var (
			orderID = ""
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.OrderInfoByID(c, orderID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantAddOrderInfo(t *testing.T) {
	convey.Convey("AddOrderInfo", t, func(ctx convey.C) {
		var (
			arg = &model.PendantOrderInfo{Mid: 650454, OrderID: strconv.FormatInt(time.Now().Unix(), 10)}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.AddOrderInfo(c, arg)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantTxAddOrderInfo(t *testing.T) {
	convey.Convey("TxAddOrderInfo", t, func(ctx convey.C) {
		var (
			arg   = &model.PendantOrderInfo{Mid: 650454, OrderID: strconv.FormatInt(time.Now().UnixNano(), 10)}
			tx, _ = d.BeginTran(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TxAddOrderInfo(c, arg, tx)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantUpdateOrderInfo(t *testing.T) {
	convey.Convey("UpdateOrderInfo", t, func(ctx convey.C) {
		var (
			arg = &model.PendantOrderInfo{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.UpdateOrderInfo(c, arg)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantTxUpdateOrderInfo(t *testing.T) {
	convey.Convey("TxUpdateOrderInfo", t, func(ctx convey.C) {
		var (
			arg   = &model.PendantOrderInfo{}
			tx, _ = d.BeginTran(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TxUpdateOrderInfo(c, arg, tx)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPackageByMid(t *testing.T) {
	convey.Convey("PackageByMid", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PackageByMid(c, mid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantPackageByID(t *testing.T) {
	convey.Convey("PackageByID", t, func(ctx convey.C) {
		var (
			mid = int64(650454)
			pid = int64(21)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.PackageByID(c, mid, pid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantEquipByMid(t *testing.T) {
	convey.Convey("EquipByMid", t, func(ctx convey.C) {
		var (
			mid = int64(88888929)
			no  = int64(44)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, noRow, err := d.EquipByMid(c, mid, no)
			ctx.Convey("Then err should be nil.res,noRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(noRow, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantEquipByMids(t *testing.T) {
	convey.Convey("EquipByMids", t, func(ctx convey.C) {
		var (
			mids = []int64{650454}
			no   = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.EquipByMids(c, mids, no)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantAddEquip(t *testing.T) {
	convey.Convey("AddEquip", t, func(ctx convey.C) {
		var (
			arg = &model.PendantEquip{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			n, err := d.AddEquip(c, arg)
			ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(n, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantTxUpdatePackageInfo(t *testing.T) {
	convey.Convey("TxUpdatePackageInfo", t, func(ctx convey.C) {
		var (
			arg   = &model.PendantPackage{Mid: 88888929, Pid: 2, Status: 1}
			tx, _ = d.BeginTran(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			n, err := d.TxUpdatePackageInfo(c, arg, tx)
			ctx.Convey("Then err should be nil.n should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(n, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantCheckPackageExpire(t *testing.T) {
	convey.Convey("CheckPackageExpire", t, func(ctx convey.C) {
		var (
			mid     = int64(650454)
			expires = int64(2147483647)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.CheckPackageExpire(c, mid, expires)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantBeginTran(t *testing.T) {
	convey.Convey("BeginTran", t, func(ctx convey.C) {
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.BeginTran(c)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantTxAddPackage(t *testing.T) {
	convey.Convey("TxAddPackage", t, func(ctx convey.C) {
		var (
			arg   = &model.PendantPackage{Mid: time.Now().Unix(), Pid: 4}
			tx, _ = d.BeginTran(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TxAddPackage(c, arg, tx)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestPendantTxAddHistory(t *testing.T) {
	convey.Convey("TxAddHistory", t, func(ctx convey.C) {
		var (
			arg   = &model.PendantHistory{}
			tx, _ = d.BeginTran(c)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := d.TxAddHistory(c, arg, tx)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}
