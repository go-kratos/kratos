package http

import (
	"context"

	"go-common/app/admin/main/up/model/signmodel"
	"go-common/library/net/http/blademaster"
)

func signUpAuditLogs(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.SignOpSearchArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignUpAuditLogs(context, arg.(*signmodel.SignOpSearchArg))
		},
		"SignUpAuditLogs")(c)
}

func signAdd(c *blademaster.Context) {
	httpPostFunc(new(signmodel.SignUpArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignAdd(context, arg.(*signmodel.SignUpArg))
		},
		"SignAdd")(c)
}

func signUpdate(c *blademaster.Context) {
	httpPostFunc(new(signmodel.SignUpArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignUpdate(context, arg.(*signmodel.SignUpArg))
		},
		"SignUp")(c)
}

func violationAdd(c *blademaster.Context) {
	httpPostFunc(new(signmodel.ViolationArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ViolationAdd(context, arg.(*signmodel.ViolationArg))
		},
		"ViolationAdd")(c)
}

func violationRetract(c *blademaster.Context) {
	httpPostFunc(new(signmodel.IDArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ViolationRetract(context, arg.(*signmodel.IDArg))
		},
		"ViolationRetract")(c)
}

func violationList(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.PageArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ViolationList(context, arg.(*signmodel.PageArg))
		},
		"ViolationList")(c)
}

func absenceAdd(c *blademaster.Context) {
	httpPostFunc(new(signmodel.AbsenceArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.AbsenceAdd(context, arg.(*signmodel.AbsenceArg))
		},
		"AbsenceAdd")(c)
}

func absenceRetract(c *blademaster.Context) {
	httpPostFunc(new(signmodel.IDArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.AbsenceRetract(context, arg.(*signmodel.IDArg))
		},
		"AbsenceRetract")(c)
}

func absenceList(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.PageArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.AbsenceList(context, arg.(*signmodel.PageArg))
		},
		"AbsenceList")(c)
}

func viewCheck(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.PowerCheckArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.ViewCheck(context, arg.(*signmodel.PowerCheckArg))
		},
		"ViewCheck")(c)
}

func signQuery(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.SignQueryArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignQuery(context, arg.(*signmodel.SignQueryArg))
		},
		"SignQuery")(c)
}

func signQueryID(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.SignIDArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignQueryID(context, arg.(*signmodel.SignIDArg))
		},
		"SignQueryID")(c)
}

func signPayComplete(c *blademaster.Context) {
	httpPostFunc(new(signmodel.SignPayCompleteArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignPayComplete(context, arg.(*signmodel.SignPayCompleteArg))
		},
		"SignPayComplete")(c)
}

func signCheckExist(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.SignCheckExsitArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.SignCheckExist(context, arg.(*signmodel.SignCheckExsitArg))
		},
		"SignCheckExist")(c)
}

func countrys(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.CommonArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.Countrys(context, arg.(*signmodel.CommonArg))
		},
		"SignCheckExist")(c)
}

func tids(c *blademaster.Context) {
	httpQueryFunc(new(signmodel.CommonArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.Tids(context, arg.(*signmodel.CommonArg))
		},
		"SignCheckExist")(c)
}
