package http

import (
	"context"
	"go-common/app/admin/main/up/model/datamodel"
	"go-common/library/net/http/blademaster"
)

func dataGetFanSummary(c *blademaster.Context) {
	httpQueryFunc(new(datamodel.GetFansSummaryArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.DataService().GetFansSummary(context, arg.(*datamodel.GetFansSummaryArg))
		},
		"dataGetFanSummary")(c)
}

func dataRelationFansHistory(c *blademaster.Context) {
	httpQueryFunc(new(datamodel.GetRelationFansHistoryArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.DataService().GetRelationFansDay(context, arg.(*datamodel.GetRelationFansHistoryArg))
		},
		"dataRelationFansHistory")(c)
}

func dataGetUpArchiveInfo(c *blademaster.Context) {
	httpQueryFunc(new(datamodel.GetUpArchiveInfoArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.DataService().GetUpArchiveInfo(context, arg.(*datamodel.GetUpArchiveInfoArg))
		},
		"dataGetUpArchiveInfo")(c)
}

func dataGetUpArchiveTagInfo(c *blademaster.Context) {
	httpQueryFunc(new(datamodel.GetUpArchiveTagInfoArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.DataService().GetUpArchiveTagInfo(context, arg.(*datamodel.GetUpArchiveTagInfoArg))
		},
		"dataGetUpArchiveTagInfo")(c)
}

func dataGetUpViewInfo(c *blademaster.Context) {
	httpQueryFunc(new(datamodel.GetUpViewInfoArg),
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return Svc.Crmservice.DataService().GetUpViewInfo(context, arg.(*datamodel.GetUpViewInfoArg))
		},
		"dataGetUpViewInfo")(c)
}
