package dao

import (
	"context"
	"go-common/app/service/bbq/user/internal/model"

	acc "go-common/app/service/main/account/api"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// RawUserCard 从主站获取用户基础信息
func (d *Dao) RawUserCard(c context.Context, mid int64) (userCard *model.UserCard, err error) {
	req := &acc.MidReq{
		Mid:    mid,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}
	res, err := d.accountClient.Card3(c, req)
	if err != nil {
		log.Error("user card rpc error(%v)", err)
		return
	}
	vipInfo := model.VIPInfo{
		Type:    res.Card.Vip.Type,
		Status:  res.Card.Vip.Status,
		DueDate: res.Card.Vip.DueDate,
	}
	userCard = &model.UserCard{
		MID:     res.Card.Mid,
		Name:    res.Card.Name,
		Sex:     res.Card.Sex,
		Rank:    res.Card.Rank,
		Face:    res.Card.Face,
		Sign:    res.Card.Sign,
		Level:   res.Card.Level,
		VIPInfo: vipInfo,
	}
	return
}

// RawUserCards 从主站获取用户基础信息
func (d *Dao) RawUserCards(c context.Context, mids []int64) (userCards map[int64]*model.UserCard, err error) {
	req := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}
	res, err := d.accountClient.Cards3(c, req)
	if err != nil {
		log.Error("user card rpc error(%v)", err)
		return
	}
	userCards = make(map[int64]*model.UserCard, len(mids))
	for _, card := range res.Cards {
		vipInfo := model.VIPInfo{
			Type:    card.Vip.Type,
			Status:  card.Vip.Status,
			DueDate: card.Vip.DueDate,
		}
		userCard := &model.UserCard{
			MID:     card.Mid,
			Name:    card.Name,
			Sex:     card.Sex,
			Rank:    card.Rank,
			Face:    card.Face,
			Sign:    card.Sign,
			Level:   card.Level,
			VIPInfo: vipInfo,
		}
		userCards[card.Mid] = userCard
	}
	return
}

// RawUserAccCards 批量获取账号信息
func (d *Dao) RawUserAccCards(c context.Context, mids []int64) (res *acc.CardsReply, err error) {
	req := &acc.MidsReq{
		Mids:   mids,
		RealIp: metadata.String(c, metadata.RemoteIP),
	}

	res, err = d.accountClient.Cards3(c, req)
	if err != nil {
		log.Error("d.accountClient.Cards3 err [%v]", err)
	}
	return
}
