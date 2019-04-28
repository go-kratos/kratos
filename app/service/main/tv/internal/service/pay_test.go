package service

//
//func TestServiceflushPayParamAsync(t *testing.T) {
//	convey.Convey("flushPayParamAsync", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			token    = ""
//			payParam = &model.PayParam{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			s.flushPayParamAsync(c, token, payParam)
//			ctx.Convey("No return values", func(ctx convey.C) {
//			})
//		})
//	})
//}
//
//func TestServiceflushUserInfoAsync(t *testing.T) {
//	convey.Convey("flushUserInfoAsync", t, func(ctx convey.C) {
//		var (
//			c   = context.Background()
//			mid = int64(0)
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			s.flushUserInfoAsync(c, mid)
//			ctx.Convey("No return values", func(ctx convey.C) {
//			})
//		})
//	})
//}
//
//func TestServicegiveMVipGiftAsync(t *testing.T) {
//	convey.Convey("giveMVipGiftAsync", t, func(ctx convey.C) {
//		var (
//			c       = context.Background()
//			mid     = int64(0)
//			pid     = int32(0)
//			orderNo = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			s.giveMVipGiftAsync(c, mid, pid, orderNo)
//			ctx.Convey("No return values", func(ctx convey.C) {
//			})
//		})
//	})
//}
//
//func TestServiceinitUserInfo(t *testing.T) {
//	convey.Convey("initUserInfo", t, func(ctx convey.C) {
//		var (
//			c     = context.Background()
//			tx, _ = s.dao.BeginTran(c)
//			mid   = int64(0)
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			ui, err := s.initUserInfo(c, tx, mid)
//			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(ui, convey.ShouldNotBeNil)
//			})
//		})
//		ctx.Reset(func() {
//			tx.Commit()
//		})
//	})
//}
//
//func TestServicePayOrder(t *testing.T) {
//	convey.Convey("PayOrder", t, func(ctx convey.C) {
//		var (
//			c  = context.Background()
//			id = int(0)
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			po, err := s.PayOrder(c, id)
//			ctx.Convey("Then err should be nil.po should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(po, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceYstOrderState(t *testing.T) {
//	convey.Convey("YstOrderState", t, func(ctx convey.C) {
//		var (
//			c       = context.Background()
//			seqNo   = ""
//			traceNo = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			res, err := s.YstOrderState(c, seqNo, traceNo)
//			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServicePayPending(t *testing.T) {
//	convey.Convey("PayPending", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			req      = &model.YstPayCallbackReq{}
//			payOrder = &model.PayOrder{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.PayPending(c, req, payOrder)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServicePayFail(t *testing.T) {
//	convey.Convey("PayFail", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			req      = &model.YstPayCallbackReq{}
//			payOrder = &model.PayOrder{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.PayFail(c, req, payOrder)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServicePaySuccess(t *testing.T) {
//	convey.Convey("PaySuccess", t, func(ctx convey.C) {
//		var (
//			c        = context.Background()
//			req      = &model.YstPayCallbackReq{}
//			payOrder = &model.PayOrder{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.PaySuccess(c, req, payOrder)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServicePayCallback(t *testing.T) {
//	convey.Convey("PayCallback", t, func(ctx convey.C) {
//		var (
//			c   = context.Background()
//			req = &model.YstPayCallbackReq{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			res := s.PayCallback(c, req)
//			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
//				ctx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceGiveMVipGift(t *testing.T) {
//	convey.Convey("GiveMVipGift", t, func(ctx convey.C) {
//		var (
//			c       = context.Background()
//			mid     = int64(0)
//			pid     = int32(0)
//			orderNo = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.GiveMVipGift(c, mid, pid, orderNo)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceUserInfoByContractCode(t *testing.T) {
//	convey.Convey("UserInfoByContractCode", t, func(ctx convey.C) {
//		var (
//			c            = context.Background()
//			contractCode = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			ui, err := s.UserInfoByContractCode(c, contractCode)
//			ctx.Convey("Then err should be nil.ui should not be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//				ctx.So(ui, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceSignContract(t *testing.T) {
//	convey.Convey("SignContract", t, func(ctx convey.C) {
//		var (
//			c            = context.Background()
//			contractCode = ""
//			contractId   = ""
//			remark       = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.SignContract(c, contractCode, contractId, remark)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceCancelContract(t *testing.T) {
//	convey.Convey("CancelContract", t, func(ctx convey.C) {
//		var (
//			c            = context.Background()
//			contractCode = ""
//			contractId   = ""
//			remark       = ""
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			err := s.CancelContract(c, contractCode, contractId, remark)
//			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
//				ctx.So(err, convey.ShouldBeNil)
//			})
//		})
//	})
//}
//
//func TestServiceWxContractCallback(t *testing.T) {
//	convey.Convey("WxContractCallback", t, func(ctx convey.C) {
//		var (
//			c   = context.Background()
//			req = &model.WxContractCallbackReq{}
//		)
//		ctx.Convey("When everything gose positive", func(ctx convey.C) {
//			res := s.WxContractCallback(c, req)
//			ctx.Convey("Then res should not be nil.", func(ctx convey.C) {
//				ctx.So(res, convey.ShouldNotBeNil)
//			})
//		})
//	})
//}
