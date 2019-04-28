package mcndao

import (
	"testing"

	"go-common/app/interface/main/mcn/model"
	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/ecode"

	"github.com/smartystreets/goconvey/convey"
)

func TestMcndaoGetMcnSignState(t *testing.T) {
	convey.Convey("GetMcnSignState", t, func(ctx convey.C) {
		var (
			fields = "*"
			mcnMid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			mcn, state, err := d.GetMcnSignState(fields, mcnMid)
			ctx.Convey("Then err should be nil.mcn,state should not be nil.", func(ctx convey.C) {

				ctx.So(err, convey.ShouldEqual, -404)
				ctx.So(state, convey.ShouldNotBeNil)
				ctx.So(mcn, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoGetUpBind(t *testing.T) {
	convey.Convey("GetUpBind", t, func(ctx convey.C) {
		var (
			query = "1=?"
			args  = "0"
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			upList, err := d.GetUpBind(query, args)
			ctx.Convey("Then err should be nil.upList should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(upList, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoBindUp(t *testing.T) {
	convey.Convey("BindUp", t, func(ctx convey.C) {
		var (
			up   = &mcnmodel.McnUp{}
			sign = &mcnmodel.McnSign{}
			arg  = &mcnmodel.McnBindUpApplyReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, affectedRow, err := d.BindUp(up, sign, arg)
			ctx.Convey("Then err should be nil.result,affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, 82010)
				ctx.So(affectedRow, convey.ShouldEqual, 0)
				ctx.So(result, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoUpdateBindUp(t *testing.T) {
	convey.Convey("UpdateBindUp", t, func(ctx convey.C) {
		var (
			values map[string]interface{}
			query  = interface{}(0)
			args   = interface{}(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			affectedRow, err := d.UpdateBindUp(values, query, args)
			ctx.Convey("Then err should be nil.affectedRow should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(affectedRow, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestMcndaoUpConfirm(t *testing.T) {
	convey.Convey("UpConfirm", t, func(ctx convey.C) {
		var (
			arg   = &mcnmodel.McnUpConfirmReq{}
			state model.MCNUPState
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := d.UpConfirm(arg, state)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestMcndaoGetBindInfo(t *testing.T) {
	convey.Convey("GetBindInfo", t, func(ctx convey.C) {
		var (
			arg = &mcnmodel.McnUpGetBindReq{}
		)
		arg.BindID = 10
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetBindInfo(arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				if err == ecode.NothingFound {
					ctx.So(err, convey.ShouldNotBeNil)
				} else {
					ctx.So(err, convey.ShouldBeNil)
				}
				if res != nil {
					ctx.So(res, convey.ShouldNotBeNil)
				} else {
					ctx.So(res, convey.ShouldBeNil)
				}
			})
		})
	})
}

func TestMcndaoGetMcnOldInfo(t *testing.T) {
	convey.Convey("GetMcnOldInfo", t, func(ctx convey.C) {
		var (
			mcnMid = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := d.GetMcnOldInfo(mcnMid)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldEqual, -404)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}
