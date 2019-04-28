package resource

import (
	"context"

	"go-common/app/admin/main/growup/model"
	"go-common/app/admin/main/growup/util"
	vip "go-common/app/service/main/vip/model"
	"go-common/library/log"
)

const _panelType = "incentive"

// VipProducts returns <vipProductID, goodsInfo> pairs
func VipProducts(c context.Context) (r map[string]*model.GoodsInfo, err error) {
	res, err := vipRPC.VipPanelInfo5(c, &vip.ArgPanel{PanelType: _panelType})
	if err != nil {
		log.Error("VipPanelInfo5 err(%v)", err)
		return
	}

	r = make(map[string]*model.GoodsInfo, len(res.Vps))
	for _, v := range res.Vps {
		r[v.PdID] = &model.GoodsInfo{
			// vip商品唯一标识
			ProductID: v.PdID,
			// vip商品名称
			ProductName: v.PdName,
			// 大会员实时价格 = 激励兑换商品的实时成本价; 单位元转换为单位分
			OriginPrice: int64(util.MulWithRound(v.DPrice, float64(100), 0)),
			// vip会员时长
			Month: v.Month,
		}
	}
	return
}
