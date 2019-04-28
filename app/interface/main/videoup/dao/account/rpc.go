package account

import (
	"context"

	accapi "go-common/app/service/main/account/api"
	relaMdl "go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Profile get profile from rpc
func (d *Dao) Profile(c context.Context, mid int64, ip string) (p *accapi.Profile, err error) {
	arg := &accapi.MidReq{
		Mid: mid,
	}
	var rpcRes *accapi.ProfileReply
	if rpcRes, err = d.acc.Profile3(c, arg); err != nil {
		log.Error("d.acc.Profile3 error(%v) | mid(%d) ip(%s) arg(%v)", err, mid, ip, arg)
		err = ecode.CreativeAccServiceErr
	}
	if rpcRes != nil {
		p = rpcRes.Profile
	}
	return
}

// Cards get cards from rpc
func (d *Dao) Cards(c context.Context, mids []int64, ip string) (cards map[int64]*accapi.Card, err error) {
	var res *accapi.CardsReply
	if len(mids) == 0 {
		return
	}
	arg := &accapi.MidsReq{
		Mids: mids,
	}
	if res, err = d.acc.Cards3(c, arg); err != nil {
		log.Error("d.acc.Cards3 error(%v) | mids(%v) ip(%s) arg(%v)", err, mids, ip, arg)
		err = ecode.CreativeAccServiceErr
	}
	if res != nil {
		cards = res.Cards
	}
	return
}

// Infos get infos from rpc
func (d *Dao) Infos(c context.Context, mids []int64, ip string) (infos map[int64]*accapi.Info, err error) {
	var res *accapi.InfosReply
	arg := &accapi.MidsReq{
		Mids: mids,
	}
	infos = make(map[int64]*accapi.Info)
	if res, err = d.acc.Infos3(c, arg); err != nil {
		log.Error("d.acc.Infos3 error(%v) | mids(%v) ip(%s) arg(%v)", err, mids, ip, arg)
		err = ecode.CreativeAccServiceErr
	}
	if res != nil {
		infos = res.Infos
	}
	return
}

// Relations get all relation state.
func (d *Dao) Relations(c context.Context, mid int64, fids []int64, ip string) (res map[int64]int, err error) {
	var rls map[int64]*relaMdl.Following
	if rls, err = d.rela.Relations(c, &relaMdl.ArgRelations{Mid: mid, Fids: fids, RealIP: ip}); err != nil {
		log.Error("d.rela.Relations mid(%d)|ip(%s)|error(%v)", mid, ip, err)
		err = ecode.CreativeAccServiceErr
		return
	}
	if len(rls) == 0 {
		log.Info("d.rela.Relations mid(%d)|ip(%s)", mid, ip)
		return
	}
	res = make(map[int64]int, len(rls))
	for _, v := range rls {
		res[v.Mid] = int(v.Attribute)
	}
	log.Info("d.rela.Relations mid(%d)|res(%+v)|rls(%+v)|ip(%s)", mid, res, rls, ip)
	return
}
