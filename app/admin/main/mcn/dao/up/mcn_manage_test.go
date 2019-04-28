package up

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestUpTxAddMCNRenewal(t *testing.T) {
	convey.Convey("TxAddMCNRenewal", t, func(ctx convey.C) {
		var (
			tx, err = d.BeginTran(context.Background())
			arg     = &model.MCNSign{MCNMID: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ctx.So(err, convey.ShouldBeNil)
			lastID, err := d.TxAddMCNRenewal(tx, arg)
			defer tx.Rollback()

			ctx.Convey("Then err should be nil. lastID should greater than zero.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(lastID, convey.ShouldBeGreaterThan, 0)
			})
		})
	})
}

func TestUpTxAddMCNPays(t *testing.T) {
	convey.Convey("TxAddMCNPays", t, func(ctx convey.C) {
		var (
			tx, err            = d.BeginTran(context.Background())
			lastID             = int64(1)
			mcnMID             = int64(1)
			payInfo1, payInfo2 = &model.SignPayReq{}, &model.SignPayReq{}
			arg                = []*model.SignPayReq{}
		)
		payInfo1.DueDate = "2018-10-25"
		payInfo1.PayValue = 10000
		arg = append(arg, payInfo1)
		payInfo2.DueDate = "2018-11-25"
		payInfo2.PayValue = 20000
		arg = append(arg, payInfo2)
		ctx.So(err, convey.ShouldBeNil)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxAddMCNPays(tx, lastID, mcnMID, arg)
			defer tx.Rollback()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpTxAddMCNUPs(t *testing.T) {
	convey.Convey("TxAddMCNUPs", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.Background())
			lastID = int64(1)
			mcnMID = int64(1)
			arg    []*model.MCNUP
			up     = &model.MCNUP{
				SignID:          1,
				MCNMID:          1,
				UPMID:           1,
				BeginDate:       0,
				EndDate:         0,
				ContractLink:    "http://www.baidu.com",
				UPAuthLink:      "http://www.baidu.com",
				State:           1,
				StateChangeTime: 0,
				Permission:      1,
			}
		)
		arg = append(arg, up)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.TxAddMCNUPs(tx, lastID, mcnMID, arg)
			defer func() {
				tx.Rollback()
			}()
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestUpUpMCNState(t *testing.T) {
	convey.Convey("UpMCNState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNStateEditReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMCNState(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMCNPay(t *testing.T) {
	convey.Convey("UpMCNPay", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNPayEditReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMCNPay(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMCNPayState(t *testing.T) {
	convey.Convey("UpMCNPayState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNPayStateEditReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMCNPayState(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNList(t *testing.T) {
	convey.Convey("MCNList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNListReq{}
		)
		arg.State = -1
		arg.Order = "s.mtime"
		arg.Sort = "DESC"
		arg.Page = 1
		arg.Size = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, ids, mids, err := d.MCNList(c, arg)
			ctx.Convey("Then err should be nil.res,ids should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(ids, convey.ShouldNotBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpbuildMCNListSQL(t *testing.T) {
	convey.Convey("buildMCNListSQL", t, func(ctx convey.C) {
		var (
			SQLType = ""
			arg     = &model.MCNListReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sql, values := d.buildMCNListSQL(SQLType, arg)
			ctx.Convey("Then sql,values should not be nil.", func(ctx convey.C) {
				ctx.So(values, convey.ShouldNotBeNil)
				ctx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcheckSort(t *testing.T) {
	convey.Convey("checkSort", t, func(ctx convey.C) {
		var (
			arg       = ""
			orderSign bool
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			p1 := checkSort(arg, orderSign)
			ctx.Convey("Then p1 should not be nil.", func(ctx convey.C) {
				ctx.So(p1, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMcnListTotal(t *testing.T) {
	convey.Convey("McnListTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNListReq{
				State: -1,
			}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.MCNListTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNPayInfos(t *testing.T) {
	convey.Convey("MCNPayInfos", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			ids = []int64{1, 2}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MCNPayInfos(c, ids)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpTxMCNRenewalUPs(t *testing.T) {
	convey.Convey("TxMCNRenewalUPs", t, func(ctx convey.C) {
		var (
			tx, _  = d.BeginTran(context.Background())
			signID = int64(1)
			mcnID  = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			ups, err := d.TxMCNRenewalUPs(tx, signID, mcnID)
			defer func() {
				if err != nil {
					tx.Rollback()
				} else {
					tx.Commit()
				}
			}()
			ctx.Convey("Then err should be nil.ups should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if len(ups) == 0 {
					ctx.So(ups, convey.ShouldBeEmpty)
				} else {
					ctx.So(ups, convey.ShouldNotBeNil)
				}
			})
		})
	})
}

func TestUpMCNInfo(t *testing.T) {
	convey.Convey("MCNInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNInfoReq{ID: 1}
		)
		arg.MCNMID = 1
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			m, err := d.MCNInfo(c, arg)
			ctx.Convey("Then err should be nil.m should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				if m == nil {
					ctx.So(m, convey.ShouldBeNil)
				} else {
					ctx.So(m, convey.ShouldNotBeNil)
				}
			})
		})
	})
}

func TestUpMCNUPList(t *testing.T) {
	convey.Convey("MCNUPList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPListReq{}
		)
		arg.DataType = 1
		arg.State = -1
		arg.UpType = -1
		arg.Order = "u.mtime"
		arg.Sort = "DESC"
		arg.Page = 1
		arg.Size = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.MCNUPList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldBeEmpty)
			})
		})
	})
}

func TestUpMCNUPListTotal(t *testing.T) {
	convey.Convey("MCNUPListTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPListReq{}
		)
		arg.DataType = 1
		arg.State = -1
		arg.UpType = -1
		arg.Order = "u.mtime"
		arg.Sort = "DESC"
		arg.Page = 1
		arg.Size = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.MCNUPListTotal(c, arg)
			ctx.Convey("Then err should be nil.count should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpbuildMCNUPListSQL(t *testing.T) {
	convey.Convey("buildMCNUPListSQL", t, func(ctx convey.C) {
		var (
			SQLType = ""
			arg     = &model.MCNUPListReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			sql, values := d.buildMCNUPListSQL(SQLType, arg)
			ctx.Convey("Then sql,values should not be nil.", func(ctx convey.C) {
				ctx.So(values, convey.ShouldNotBeNil)
				ctx.So(sql, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMCNUPState(t *testing.T) {
	convey.Convey("UpMCNUPState", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPStateEditReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.UpMCNUPState(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestUpMcnSignByMCNMID(t *testing.T) {
	convey.Convey("McnSignByMCNMID", t, func(ctx convey.C) {
		var (
			c     = context.Background()
			mcnID = int64(1)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			row, err := d.McnSignByMCNMID(c, mcnID)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(row, convey.ShouldNotBeEmpty)
			})
		})
	})
}
func TestUpMCNCheatList(t *testing.T) {
	convey.Convey("TestUpMCNCheatList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNCheatListReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, mids, err := d.MCNCheatList(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				for k, r := range rows {
					fmt.Printf("%d:%+v \n", k, r)
				}
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
				ctx.So(mids, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNCheatListTotal(t *testing.T) {
	convey.Convey("TestUpMCNCheatListTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNCheatListReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.MCNCheatListTotal(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				fmt.Printf("count:%d \n", count)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNCheatUPList(t *testing.T) {
	convey.Convey("TestUpMCNCheatUPList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNCheatUPListReq{UPMID: 1}
		)
		arg.Page = 1
		arg.Size = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			rows, err := d.MCNCheatUPList(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				for k, r := range rows {
					fmt.Printf("%d:%+v \n", k, r)
				}
				ctx.So(err, convey.ShouldBeNil)
				if len(rows) == 0 {
					ctx.So(rows, convey.ShouldBeEmpty)
				} else {
					ctx.So(rows, convey.ShouldNotBeNil)
				}

			})
		})
	})
}
func TestUpMCNCheatUPListTotal(t *testing.T) {
	convey.Convey("TestUpMCNCheatUPListTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNCheatUPListReq{UPMID: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			count, err := d.MCNCheatUPListTotal(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				fmt.Printf("count:%d \n", count)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNImportUPInfo(t *testing.T) {
	convey.Convey("TestUpMCNImportUPInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNImportUPInfoReq{UPMID: 1}
		)
		arg.SignID = 1
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.MCNImportUPInfo(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				fmt.Printf("res:%+v \n", res)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMCNImportUPRewardSign(t *testing.T) {
	convey.Convey("TestUpUpMCNImportUPRewardSign", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNImportUPRewardSignReq{UPMID: 1}
		)
		arg.SignID = 1
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := d.UpMCNImportUPRewardSign(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				fmt.Printf("res:%+v \n", res)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpUpMCNPermission(t *testing.T) {
	convey.Convey("UpMCNPermission", t, func(ctx convey.C) {
		var (
			c          = context.Background()
			signID     = int64(0)
			permission = uint32(0)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.UpMCNPermission(c, signID, permission)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpMCNIncreaseList(t *testing.T) {
	convey.Convey("TestUpMCNIncreaseList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNIncreaseListReq{}
		)
		arg.SignID = 99
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			rows, err := d.MCNIncreaseList(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				for k, r := range rows {
					fmt.Printf("%d:%+v \n", k, r)
				}
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(rows, convey.ShouldNotBeNil)
			})
		})
	})
}
func TestUpMCNIncreaseListTotal(t *testing.T) {
	convey.Convey("TestUpMCNIncreaseListTotal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNIncreaseListReq{}
		)
		arg.SignID = 99
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			count, err := d.MCNIncreaseListTotal(c, arg)
			ctx.Convey("Then err should be nil.rows should not be nil.", func(ctx convey.C) {
				fmt.Printf("count:%d \n", count)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(count, convey.ShouldNotBeNil)
			})
		})
	})
}
