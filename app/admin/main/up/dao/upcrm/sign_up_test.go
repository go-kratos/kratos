package upcrm

import (
	"testing"
	"time"

	"go-common/app/admin/main/up/model/signmodel"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmInsertSignUp(t *testing.T) {
	convey.Convey("InsertSignUp", t, func(ctx convey.C) {
		var (
			db = d.crmdb
			up = &signmodel.SignUp{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.InsertSignUp(db, up)
			d.crmdb.Delete(up)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmInsertPayInfo(t *testing.T) {
	convey.Convey("InsertPayInfo", t, func(ctx convey.C) {
		var (
			db   = d.crmdb
			info = &signmodel.SignPay{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.InsertPayInfo(db, info)
			d.crmdb.Delete(info)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmInsertTaskInfo(t *testing.T) {
	convey.Convey("InsertTaskInfo", t, func(ctx convey.C) {
		var (
			db   = d.crmdb
			info = &signmodel.SignTask{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.InsertTaskInfo(db, info)
			d.crmdb.Delete(info)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmInsertContractInfo(t *testing.T) {
	convey.Convey("InsertContractInfo", t, func(ctx convey.C) {
		var (
			db   = d.crmdb
			info = &signmodel.SignContract{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.InsertContractInfo(db, info)
			d.crmdb.Delete(info)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmDelPayInfo(t *testing.T) {
	convey.Convey("DelPayInfo", t, func(ctx convey.C) {
		var (
			db  = d.crmdb
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.DelPayInfo(db, ids)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmDelTaskInfo(t *testing.T) {
	convey.Convey("DelTaskInfo", t, func(ctx convey.C) {
		var (
			db  = d.crmdb
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.DelTaskInfo(db, ids)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmDelSignContract(t *testing.T) {
	convey.Convey("DelSignContract", t, func(ctx convey.C) {
		var (
			db  = d.crmdb
			ids = []int64{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.DelSignContract(db, ids)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmSignUpID(t *testing.T) {
	convey.Convey("SignUpID", t, func(ctx convey.C) {
		var (
			sigID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			su, msp, mst, msc, err := d.SignUpID(sigID)
			ctx.Convey("Then err should be nil.su,msp,mst,msc should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(msc, convey.ShouldBeNil)
				ctx.So(mst, convey.ShouldBeNil)
				ctx.So(msp, convey.ShouldBeNil)
				ctx.So(su, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpcrmGetSignIDByCondition(t *testing.T) {
	convey.Convey("GetSignIDByCondition", t, func(ctx convey.C) {
		var (
			arg = &signmodel.SignQueryArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			signIDs, err := d.GetSignIDByCondition(arg)
			ctx.Convey("Then err should be nil.signIDs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(signIDs, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpcrmGetSignUpByID(t *testing.T) {
	convey.Convey("GetSignUpByID", t, func(ctx convey.C) {
		var (
			signID = []uint32{}
			order  = ""
			offset = int(0)
			limit  = int(0)
			query  = interface{}(0)
			args   = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetSignUpByID(signID, order, offset, limit, query, args)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetSignUpCount(t *testing.T) {
	convey.Convey("GetSignUpCount", t, func(ctx convey.C) {
		var (
			query = ""
			args  = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count := d.GetSignUpCount(query, args)
			ctx.Convey("Then count should not be nil.", func(ctx convey.C) {
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetTask(t *testing.T) {
	convey.Convey("GetTask", t, func(ctx convey.C) {
		var (
			signID = []uint32{}
			state  = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetTask(signID, state)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetPay(t *testing.T) {
	convey.Convey("GetPay", t, func(ctx convey.C) {
		var (
			signID = []uint32{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetPay(signID)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetContract(t *testing.T) {
	convey.Convey("GetContract", t, func(ctx convey.C) {
		var (
			signID = []uint32{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetContract(signID)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmPayComplete(t *testing.T) {
	convey.Convey("PayComplete", t, func(ctx convey.C) {
		var (
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.PayComplete(ids)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetDueSignUp(t *testing.T) {
	convey.Convey("GetDueSignUp", t, func(ctx convey.C) {
		var (
			now             = time.Now()
			expireAfterDays = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetDueSignUp(now, expireAfterDays)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmGetDuePay(t *testing.T) {
	convey.Convey("GetDuePay", t, func(ctx convey.C) {
		var (
			now             = time.Now()
			expireAfterDays = int(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := d.GetDuePay(now, expireAfterDays)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmUpdateEmailState(t *testing.T) {
	convey.Convey("UpdateEmailState", t, func(ctx convey.C) {
		var (
			table = "sign_up"
			ids   = []int64{}
			state = int8(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.UpdateEmailState(table, ids, state)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldEqual, 0)
			})
		})
	})
}

func TestUpcrmCheckUpHasValidContract(t *testing.T) {
	convey.Convey("CheckUpHasValidContract", t, func(ctx convey.C) {
		var (
			mid  = int64(0)
			date = time.Now()
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			exist, err := d.CheckUpHasValidContract(mid, date)
			ctx.Convey("Then err should be nil.exist should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(exist, convey.ShouldNotBeNil)
			})
		})
	})
}
