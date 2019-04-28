package dao

import (
	"context"

	acc "go-common/app/service/main/account/api"
	"go-common/library/net/metadata"
)

//GetUserCards 获取用户卡片信息
func (d *Dao) GetUserCards(ctx context.Context, mids []int64) (*acc.CardsReply, error) {
	req := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(ctx, metadata.RemoteIP),
	}
	return d.accClient.Cards3(ctx, req)
}
