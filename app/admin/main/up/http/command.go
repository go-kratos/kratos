package http

import (
	"context"
	"go-common/app/admin/main/up/model"
	"go-common/library/net/http/blademaster"
)

func commandRefreshUpRank(c *blademaster.Context) {
	httpQueryFunc(new(model.CommandCommonArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.CommandRefreshUpRank(context, arg.(*model.CommandCommonArg))
		},
		"SignCheckTask")(c)
}
