package http

import (
	"context"

	"go-common/app/interface/main/mcn/model/mcnmodel"
	"go-common/library/net/http/blademaster"
)

func mcnGetRankArchiveLikesAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetRankAPIReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.McnGetRankArchiveLikesAPI(context, arg.(*mcnmodel.McnGetRankAPIReq))
		},
		"McnGetRankArchiveLikesAPI",
		nil,
		nil,
	)(c)
}

func getMcnSummaryAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetDataSummaryReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnSummaryAPI(context, arg.(*mcnmodel.McnGetDataSummaryReq))
		},
		"GetMcnSummaryAPI",
		nil,
		nil,
	)(c)
}

func getIndexIncAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetIndexIncReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetIndexIncAPI(context, arg.(*mcnmodel.McnGetIndexIncReq))
		},
		"GetIndexIncAPI",
		nil,
		nil,
	)(c)
}

func getIndexSourceAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetIndexSourceReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetIndexSourceAPI(context, arg.(*mcnmodel.McnGetIndexSourceReq))
		},
		"GetIndexSourceAPI",
		nil,
		nil,
	)(c)
}

func getPlaySourceAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetPlaySourceReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetPlaySourceAPI(context, arg.(*mcnmodel.McnGetPlaySourceReq))
		},
		"GetPlaySourceAPI",
		nil,
		nil,
	)(c)
}

func getMcnFansAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansAPI(context, arg.(*mcnmodel.McnGetMcnFansReq))
		},
		"GetMcnFansAPI",
		nil,
		nil,
	)(c)
}

func getMcnFansIncAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansIncReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansIncAPI(context, arg.(*mcnmodel.McnGetMcnFansIncReq))
		},
		"GetMcnFansIncAPI",
		nil,
		nil,
	)(c)
}

func getMcnFansDecAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansDecReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansDecAPI(context, arg.(*mcnmodel.McnGetMcnFansDecReq))
		},
		"GetMcnFansDecAPI",
		nil,
		nil,
	)(c)
}

func getMcnFansAttentionWayAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetMcnFansAttentionWayReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetMcnFansAttentionWayAPI(context, arg.(*mcnmodel.McnGetMcnFansAttentionWayReq))
		},
		"GetMcnFansAttentionWayAPI",
		nil,
		nil,
	)(c)
}

func getFansBaseFansAttrAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetBaseFansAttrReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansBaseFansAttrAPI(context, arg.(*mcnmodel.McnGetBaseFansAttrReq))
		},
		"GetFansBaseFansAttrAPI",
		nil,
		nil,
	)(c)
}

func getFansAreaAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansAreaReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansAreaAPI(context, arg.(*mcnmodel.McnGetFansAreaReq))
		},
		"GetFansAreaAPI",
		nil,
		nil,
	)(c)
}

func getFansTypeAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansTypeReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansTypeAPI(context, arg.(*mcnmodel.McnGetFansTypeReq))
		},
		"GetFansTypeAPI",
		nil,
		nil,
	)(c)
}

func getFansTagAPI(c *blademaster.Context) {
	oarg := new(mcnmodel.McnGetFansTagReq)
	httpGetFunc(
		oarg,
		func(context context.Context, arg interface{}) (res interface{}, err error) {
			return srv.GetFansTagAPI(context, arg.(*mcnmodel.McnGetFansTagReq))
		},
		"GetFansTagAPI",
		nil,
		nil,
	)(c)
}
