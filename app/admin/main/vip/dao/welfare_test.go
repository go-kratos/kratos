package dao

import (
	"context"
	"testing"

	"go-common/app/admin/main/vip/model"
	"go-common/library/log"

	"github.com/smartystreets/goconvey/convey"
)

func TestDaoWelfareTypeAdd(t *testing.T) {
	convey.Convey("WelfareTypeAdd", t, func(ctx convey.C) {
		var (
			wt = &model.WelfareType{ID: 1000000, Name: "utTestAdd"}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareTypeAdd(wt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareTypeUpd(t *testing.T) {
	convey.Convey("WelfareTypeUpd", t, func(ctx convey.C) {
		var (
			wt = &model.WelfareType{ID: 1000000, Name: "utTestUpd"}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareTypeUpd(wt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareTypeState(t *testing.T) {
	convey.Convey("WelfareTypeState", t, func(ctx convey.C) {
		var (
			tx       = d.BeginGormTran(context.Background())
			id       = int(1000000)
			state    = int(1)
			operId   = int(0)
			operName = "utTest"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareTypeState(tx, id, state, operId, operName)
			tx.Commit()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareTypeList(t *testing.T) {
	convey.Convey("WelfareTypeList", t, func(ctx convey.C) {
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			wts, err := d.WelfareTypeList()
			ctx.Convey("Then err should be nil.wts should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(wts, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWelfareAdd(t *testing.T) {
	convey.Convey("WelfareAdd", t, func(ctx convey.C) {
		var (
			wt = &model.Welfare{ID: 100000, Tid: 100000, WelfareName: "utTestAdd"}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareAdd(wt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareUpd(t *testing.T) {
	convey.Convey("WelfareUpd", t, func(ctx convey.C) {
		var (
			wt = &model.WelfareReq{ID: 100000, Tid: 100000, WelfareName: "utTestUpd"}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareUpd(wt)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareState(t *testing.T) {
	convey.Convey("WelfareState", t, func(ctx convey.C) {
		var (
			id       = int(100000)
			state    = int(1)
			operId   = int(0)
			operName = "utTest"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareState(id, state, operId, operName)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoResetWelfareTid(t *testing.T) {
	convey.Convey("ResetWelfareTid", t, func(ctx convey.C) {
		var (
			tx  = d.BeginGormTran(context.Background())
			tid = int(100000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.ResetWelfareTid(tx, tid)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareList(t *testing.T) {
	convey.Convey("WelfareList", t, func(ctx convey.C) {
		var (
			tid = int(100000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ws, err := d.WelfareList(tid)
			ctx.Convey("Then err should be nil.ws should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ws, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWelfareBatchSave(t *testing.T) {
	convey.Convey("WelfareBatchSave", t, func(ctx convey.C) {
		var (
			wcb = &model.WelfareCodeBatch{ID: 100000, BatchName: "utTest"}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareBatchSave(wcb)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareBatchList(t *testing.T) {
	convey.Convey("WelfareBatchList", t, func(ctx convey.C) {
		var (
			wid = int(100000)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			wbs, err := d.WelfareBatchList(wid)
			ctx.Convey("Then err should be nil.wbs should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(wbs, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestDaoWelfareBatchState(t *testing.T) {
	convey.Convey("WelfareBatchState", t, func(ctx convey.C) {
		var (
			tx       = d.BeginGormTran(context.Background())
			id       = int(100000)
			state    = int(1)
			operId   = int(0)
			operName = "utTest"
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareBatchState(tx, id, state, operId, operName)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaoWelfareCodeBatchInsert(t *testing.T) {
	convey.Convey("WelfareCodeBatchInsert", t, func(ctx convey.C) {
		var (
			wcs = []*model.WelfareCode{}
			wc  = new(model.WelfareCode)
		)
		wc.Code = "utTest"
		wc.Bid = 1000000
		wc.Wid = 1000000
		wcs = append(wcs, wc)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareCodeBatchInsert(wcs)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
		d.delCode()
	})
}

func (d *Dao) delCode() {
	codeInfo := &model.WelfareCode{
		Bid:  1000000,
		Wid:  1000000,
		Code: "utTest",
	}
	if err := d.vip.Where("bid = ? and wid =? and code = ?", codeInfo.Bid, codeInfo.Wid, codeInfo.Code).Delete(codeInfo).Error; err != nil {
		log.Error("DelUser error(%v)", err)
	}
}

func TestDaoWelfareCodeStatus(t *testing.T) {
	convey.Convey("WelfareCodeStatus", t, func(ctx convey.C) {
		var (
			tx    = d.BeginGormTran(context.Background())
			bid   = int(100000)
			state = int(1)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := d.WelfareCodeStatus(tx, bid, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestDaogetBatchInsertSQL(t *testing.T) {
	convey.Convey("getBatchInsertSQL", t, func(ctx convey.C) {
		var (
			buff = []*model.WelfareCode{}
		)
		buff = append(buff, new(model.WelfareCode))
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			stmt, valueArgs := getBatchInsertSQL(buff)
			ctx.Convey("Then stmt,valueArgs should not be nil.", func(ctx convey.C) {
				ctx.So(valueArgs, convey.ShouldNotBeNil)
				ctx.So(stmt, convey.ShouldNotBeNil)
			})
		})
	})
}
