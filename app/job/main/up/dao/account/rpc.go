package account

import (
	"context"

	account "go-common/app/service/main/account/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Profile get profile from rpc
func (d *Dao) Profile(c context.Context, mid int64, ip string) (p *account.Profile, err error) {
	arg := &account.ArgMid{
		Mid: mid,
	}
	if p, err = d.acc.Profile3(c, arg); err != nil {
		log.Error("d.acc.Profile3 error(%v) | mid(%d) ip(%s) arg(%v)", err, mid, ip, arg)
		err = ecode.CreativeAccServiceErr
	}
	return
}

//Infos get up infos
func (d *Dao) Infos(c context.Context, mids []int64, ip string) (infos map[int64]*account.Info, err error) {
	arg := &account.ArgMids{
		Mids: mids,
	}
	if infos, err = d.acc.Infos3(c, arg); err != nil {
		log.Error("d.acc.info3 error(%v) arg(%v)", err, arg)
		err = ecode.CreativeAccServiceErr
	}
	return
}
