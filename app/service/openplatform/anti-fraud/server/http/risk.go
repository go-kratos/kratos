package http

import (
	grpcv1 "go-common/app/service/openplatform/anti-fraud/api/grpc/v1"
	"go-common/app/service/openplatform/anti-fraud/api/http/v1"
	"go-common/app/service/openplatform/anti-fraud/model"
	bm "go-common/library/net/http/blademaster"
)

//riskCheck 风险检查
func riskCheck(c *bm.Context) {
	params := new(grpcv1.RiskCheckRequest)
	if err := c.Bind(params); err != nil {
		return
	}

	res, err := svc.RiskCheck(c, params)

	svc.Ticket2PayShield(c, params, res)

	if res.Rank == model.RankDoubt {
		res.Rank = model.RankNormal
		res.Method = model.MethodPass
		res.Desc = model.CheckPass
	}

	c.JSON(res, err)
}

//riskCheckV2 风险检查
func riskCheckV2(c *bm.Context) {
	params := new(v1.RiskCheckV2Request)
	if err := c.Bind(params); err != nil {
		return
	}
	res, err := svc.RiskCheckV2(c, params)

	if err == nil {
		svc.PayShield(c, params, res)
	}

	c.JSON(res, err)
}

// ipList IP列表接口
func ipList(c *bm.Context) {
	params := new(v1.IPListRequest)

	c.JSON(svc.IPList(c, params))
}

//ipDetail ip详情列表
func ipDetail(c *bm.Context) {
	params := new(v1.IPDetailRequest)
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.IPDetail(c, params))
}

// uidList uid列表接口
func uidList(c *bm.Context) {
	params := new(v1.UIDListRequest)

	c.JSON(svc.UIDList(c, params))
}

//uidDetail uid详情列表
func uidDetail(c *bm.Context) {
	params := new(v1.UIDDetailRequest)
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.UIDDetail(c, params))
}

// ipBlack ip黑名单
func ipBlack(c *bm.Context) {
	params := new(v1.IPBlackRequest)
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.IPBlack(c, params))
}

// uidBlack uid黑名单
func uidBlack(c *bm.Context) {
	params := new(v1.UIDBlackRequest)
	if err := c.Bind(params); err != nil {
		return
	}

	c.JSON(svc.UIDBlack(c, params))
}
