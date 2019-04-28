package ugcpay

import (
	"context"
	"fmt"

	"go-common/app/interface/main/app-view/conf"
	"go-common/app/interface/main/app-view/model/view"
	ugcpayrpc "go-common/app/service/main/ugcpay/api/grpc/v1"

	"github.com/pkg/errors"
)

// Dao is ugcpay dao
type Dao struct {
	// grpc
	ugcpayRPC ugcpayrpc.UGCPayClient
}

// New ugcpay dao
func New(c *conf.Config) (d *Dao) {
	d = &Dao{}
	var err error
	d.ugcpayRPC, err = ugcpayrpc.NewClient(nil)
	if err != nil {
		panic(fmt.Sprintf("ugcpay NewClient error(%v)", err))
	}
	return
}

// AssetRelationDetail ugcpay
func (d *Dao) AssetRelationDetail(c context.Context, mid, aid int64, platform string) (res *view.Asset, err error) {
	var (
		arg   = &ugcpayrpc.AssetRelationDetailReq{Mid: mid, Oid: aid, Otype: "archive", Currency: "bp"}
		asset *ugcpayrpc.AssetRelationDetailResp
	)
	if asset, err = d.ugcpayRPC.AssetRelationDetail(c, arg); err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	res = &view.Asset{}
	switch asset.RelationState {
	case "paid":
		res.Paid = 1
	}
	if price, ok := asset.AssetPlatformPrice[platform]; ok {
		res.Price = price
	} else {
		res.Price = asset.AssetPrice
	}
	res.Msg.Desc1 = "本视频为付费内容，支付" + fmt.Sprintf("%0.2f", float64(res.Price)/100) + "币即可观看"
	res.Msg.Desc2 = "用一点点奖励来支持UP主们创作吧"
	return
}
