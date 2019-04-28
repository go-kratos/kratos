package dao

import (
	"context"

	acc "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
)

// AccountCards get cards by mids
func (d *Dao) AccountCards(ctx context.Context, mids []int64) (cards map[int64]*accmdl.Card, err error) {
	var (
		req = &acc.MidsReq{
			Mids: mids,
		}
		reply *acc.CardsReply
	)
	if reply, err = d.accountAPI.Cards3(ctx, req); err != nil {
		return
	}
	cards = reply.Cards
	return
}
