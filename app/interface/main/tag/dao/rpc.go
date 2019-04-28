package dao

import (
	"context"

	account "go-common/app/service/main/account/api"
	"go-common/library/log"
)

// UserCard account user card info.
func (d *Dao) UserCard(c context.Context, mid int64) (*account.Card, error) {
	arg := &account.MidReq{
		Mid: mid,
	}
	reply, err := d.accRPC.Card3(c, arg)
	if err != nil {
		log.Error("d.dao.UserCard(%v) error(%v)", arg, err)
		return nil, err
	}
	return reply.Card, nil
}

// UserCards account users card info.
func (d *Dao) UserCards(c context.Context, mids []int64) (map[int64]*account.Card, error) {
	arg := &account.MidsReq{
		Mids: mids,
	}
	reply, err := d.accRPC.Cards3(c, arg)
	if err != nil {
		log.Error("d.dao.UserCards(%v) error(%v)", arg, err)
		return nil, err
	}
	return reply.Cards, nil
}

// UserProfile get user profile.
func (d *Dao) UserProfile(c context.Context, mid int64) (*account.Profile, error) {
	arg := &account.MidReq{
		Mid: mid,
	}
	reply, err := d.accRPC.Profile3(c, arg)
	if err != nil {
		log.Error("d.dao.UserProfile(%v) error(%v)", arg, err)
		return nil, err
	}
	return reply.Profile, nil
}
