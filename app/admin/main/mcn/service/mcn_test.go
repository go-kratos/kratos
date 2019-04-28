package service

import (
	"context"
	"fmt"
	"testing"

	"go-common/app/admin/main/mcn/model"

	"github.com/smartystreets/goconvey/convey"
)

func TestServiceMcnSignEntry(t *testing.T) {
	convey.Convey("McnSignEntry", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignEntryReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.McnSignEntry(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMcnSignList(t *testing.T) {
	convey.Convey("McnSignList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnSignList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnSignOP(t *testing.T) {
	convey.Convey("McnSignOP", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignStateOpReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.McnSignOP(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMcnUPReviewList(t *testing.T) {
	convey.Convey("McnUPReviewList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPStateReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.McnUPReviewList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnUPOP(t *testing.T) {
	convey.Convey("McnUPOP", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPStateOpReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.McnUPOP(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMcnPermitOP(t *testing.T) {
	convey.Convey("McnPermitOP", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNSignPermissionReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.McnPermitOP(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServicegetPermitOpenOrClosed(t *testing.T) {
	convey.Convey("getPermitOpenOrClosed", t, func(ctx convey.C) {
		var (
			a = uint32(5)
			b = uint32(5)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			open, closed := s.getPermitOpenOrClosed(a, b)
			ctx.Convey("Then open,closed should not be nil.", func(ctx convey.C) {
				ctx.So(closed, convey.ShouldNotBeNil)
				ctx.So(open, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServicegetUpPermitString(t *testing.T) {
	convey.Convey("getPermitOpenOrClosed", t, func(ctx convey.C) {
		var (
			a = uint32(5)
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			ps := s.getUpPermitString(a)
			ctx.Convey("Then open,closed should not be nil.", func(ctx convey.C) {
				ctx.So(ps, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnUPPermitList(t *testing.T) {
	convey.Convey("McnUPPermitList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPPermitStateReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.McnUPPermitList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMcnUPPermitOP(t *testing.T) {
	convey.Convey("McnUPPermitOP", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPPermitOPReq{}
		)
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			err := s.McnUPPermitOP(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMCNList(t *testing.T) {
	convey.Convey("MCNList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNListReq{State: -1}
		)
		arg.MCNMID = 12345
		arg.Export = "csv"
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.MCNList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				for k, v := range res.List {
					fmt.Printf("re[%d]:%+v", k, v)
				}
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMCNPayEdit(t *testing.T) {
	convey.Convey("MCNPayEdit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNPayEditReq{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.MCNPayEdit(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMCNPayStateEdit(t *testing.T) {
	convey.Convey("MCNPayStateEdit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNPayStateEditReq{ID: 1, MCNMID: 212895899, SignID: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.MCNPayStateEdit(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMCNStateEdit(t *testing.T) {
	convey.Convey("MCNStateEdit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNStateEditReq{ID: 1, MCNMID: 1212, Action: model.McnAccountRestore}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.MCNStateEdit(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMCNRenewal(t *testing.T) {
	convey.Convey("MCNRenewal", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNRenewalReq{ID: 5, MCNMID: 27515432, BeginDate: "2018-09-22", EndDate: "2019-09-23", ContractLink: "ContractLink"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.MCNRenewal(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}

func TestServiceMCNInfo(t *testing.T) {
	convey.Convey("MCNInfo", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNInfoReq{ID: 9}
		)
		arg.MCNMID = 27515432
		ctx.Convey("When everything gose positive", func(ctx convey.C) {
			res, err := s.MCNInfo(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				fmt.Println(res)
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMCNUPList(t *testing.T) {
	convey.Convey("McnUPList", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPListReq{SignID: 3, State: -1, SortFansCountActive: "asc"}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.MCNUPList(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				for k, v := range res.List {
					fmt.Printf("re[%d]:%+v \n", k, v)
				}
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestServiceMCNUPStateEdit(t *testing.T) {
	convey.Convey("MCNUPStateEdit", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &model.MCNUPStateEditReq{ID: 1, SignID: 1}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			err := s.MCNUPStateEdit(c, arg)
			ctx.Convey("Then err should be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
			})
		})
	})
}
