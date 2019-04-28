package account

import (
	"context"

	accwar "go-common/app/service/main/account/api"
	accmdl "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Card3 get card info by mid
func (d *Dao) Card3(c context.Context, mid int64) (res *accmdl.Card, err error) {
	var (
		arg   = &accwar.MidReq{Mid: mid}
		reply *accwar.CardReply
	)
	if reply, err = d.accClient.Card3(c, arg); err != nil || reply == nil || reply.Card == nil {
		if err != nil {
			log.Error("s.accDao.Info(%d) error(%v)", mid, err)
		}
		err = ecode.AccessDenied
		return
	}
	res = reply.Card
	return
}

// IsVip checks whether the member is vip
func IsVip(card *accmdl.Card) bool {
	if card.Vip.Type == 0 || card.Vip.Status == 0 || card.Vip.Status == 2 || card.Vip.Status == 3 {
		return false
	}
	return true
}
