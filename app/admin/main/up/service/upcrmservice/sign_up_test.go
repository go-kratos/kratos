package upcrmservice

import (
	"context"
	"testing"

	"go-common/app/admin/main/up/model/signmodel"

	"github.com/jinzhu/gorm"
	"github.com/smartystreets/goconvey/convey"
)

func TestUpcrmserviceSignUpAuditLogs(t *testing.T) {
	convey.Convey("SignUpAuditLogs", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &signmodel.SignOpSearchArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.SignUpAuditLogs(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignAdd(t *testing.T) {
	convey.Convey("SignAdd", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.SignUpArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.SignAdd(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignUpdate(t *testing.T) {
	convey.Convey("SignUpdate", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.SignUpArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.SignUpdate(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceViolationAdd(t *testing.T) {
	convey.Convey("ViolationAdd", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.ViolationArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ViolationAdd(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceViolationRetract(t *testing.T) {
	convey.Convey("ViolationRetract", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.IDArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ViolationRetract(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceViolationList(t *testing.T) {
	convey.Convey("ViolationList", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.PageArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.ViolationList(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceAbsenceAdd(t *testing.T) {
	convey.Convey("AbsenceAdd", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.AbsenceArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.AbsenceAdd(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmservicegetOrCreateTaskHistory(t *testing.T) {
	convey.Convey("getOrCreateTaskHistory", t, func(ctx convey.C) {
		var (
			tx     = &gorm.DB{}
			signID = int64(0)
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.getOrCreateTaskHistory(tx, signID)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceAbsenceRetract(t *testing.T) {
	convey.Convey("AbsenceRetract", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.IDArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.AbsenceRetract(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceAbsenceList(t *testing.T) {
	convey.Convey("AbsenceList", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.PageArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.AbsenceList(context, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceViewCheck(t *testing.T) {
	convey.Convey("ViewCheck", t, func(ctx convey.C) {
		var (
			context = context.Background()
			arg     = &signmodel.PowerCheckArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.ViewCheck(context, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceaddPayInfo(t *testing.T) {
	convey.Convey("addPayInfo", t, func(ctx convey.C) {
		var (
			tx  = &gorm.DB{}
			arg = &signmodel.SignPayInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := s.addPayInfo(tx, arg)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceaddTaskInfo(t *testing.T) {
	convey.Convey("addTaskInfo", t, func(ctx convey.C) {
		var (
			tx  = &gorm.DB{}
			arg = &signmodel.SignTaskInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := s.addTaskInfo(tx, arg)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceaddContractInfo(t *testing.T) {
	convey.Convey("addContractInfo", t, func(ctx convey.C) {
		var (
			tx  = &gorm.DB{}
			arg = &signmodel.SignContractInfoArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			id, err := s.addContractInfo(tx, arg)
			ctx.Convey("Then err should be nil.id should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(id, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignQuery(t *testing.T) {
	convey.Convey("SignQuery", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &signmodel.SignQueryArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.SignQuery(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignQueryID(t *testing.T) {
	convey.Convey("SignQueryID", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &signmodel.SignIDArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			res, err := s.SignQueryID(c, arg)
			ctx.Convey("Then err should be nil.res should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(res, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignPayComplete(t *testing.T) {
	convey.Convey("SignPayComplete", t, func(ctx convey.C) {
		var (
			con = context.Background()
			arg = &signmodel.SignPayCompleteArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.SignPayComplete(con, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}

func TestUpcrmserviceSignCheckExist(t *testing.T) {
	convey.Convey("SignCheckExist", t, func(ctx convey.C) {
		var (
			c   = context.Background()
			arg = &signmodel.SignCheckExsitArg{}
		)
		ctx.Convey("When everything goes positive", func(ctx convey.C) {
			result, err := s.SignCheckExist(c, arg)
			ctx.Convey("Then err should be nil.result should not be nil.", func(ctx convey.C) {
				ctx.So(err, convey.ShouldBeNil)
				ctx.So(result, convey.ShouldNotBeNil)
			})
		})
	})
}
