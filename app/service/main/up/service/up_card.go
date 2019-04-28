package service

import (
	"context"

	accgrpc "go-common/app/service/main/account/api"
	"go-common/app/service/main/up/dao/global"
	"go-common/app/service/main/up/model"
	"go-common/library/log"
	"go-common/library/net/http/blademaster"
	"go-common/library/net/metadata"
)

// ListCardBase list card mids
func (s *Service) ListCardBase(ctx *blademaster.Context) (mids []int64, err error) {
	mids, err = s.card.ListUpMID(ctx)
	return
}

// GetCardInfo get card content by mid
func (s *Service) GetCardInfo(ctx *blademaster.Context, mid int64) (card *model.UpCard, err error) {
	cards, err := s.GetCardInfoByMids(ctx, []int64{mid})
	if err != nil {
		return
	}
	card = cards[mid]
	return
}

// ListCardDetail page list card content
func (s *Service) ListCardDetail(ctx context.Context, offset uint, size uint) (cards map[int64]*model.UpCard, total int, err error) {
	cards = make(map[int64]*model.UpCard)

	total, err = s.card.CountUpCard(ctx)
	if err != nil {
		log.Error("ListCardDetail CountUpCard err, err=%v", err)
		return
	}

	if total <= 0 {
		return
	}

	infos, err := s.card.ListUpInfo(ctx, offset, size)
	if err != nil {
		log.Error("ListCardDetail ListUpInfo err, err=%v", err)
		return
	}

	if len(infos) <= 0 {
		return
	}

	var mids []int64
	for _, info := range infos {
		mids = append(mids, info.MID)
	}

	midAccountsMap, err := s.getUpAccountsMap(ctx, mids)
	if err != nil {
		log.Error("GetCardInfoByMids getUpAccountsMap err, mids=%d, err=%v", mids, err.Error())
		return
	}
	midImagesMap, err := s.card.MidImagesMap(ctx, mids)
	if err != nil {
		log.Error("card.MidImagesMap err, mids=%d, err=%v", mids, err.Error())
		return
	}
	midVideosMap, err := s.getUpVideosMap(ctx, mids)
	if err != nil {
		log.Error("GetCardInfoByMids getUpVideosMap err, mids=%d, err=%v", mids, err.Error())
		return
	}

	for _, info := range infos {
		card := &model.UpCard{
			UpCardInfo: info,
			Accounts:   []*model.UpCardAccount{},
			Images:     []*model.UpCardImage{},
			Videos:     []*model.UpCardVideo{},
		}
		mid := info.MID

		if data, ok := midAccountsMap[mid]; ok {
			card.Accounts = data
		}
		if data, ok := midImagesMap[mid]; ok {
			card.Images = data
		}
		if data, ok := midVideosMap[mid]; ok {
			card.Videos = data
		}

		cards[mid] = card
	}

	return
}

// GetCardInfoByMids get <mid, card> map by mids
func (s *Service) GetCardInfoByMids(ctx *blademaster.Context, mids []int64) (cards map[int64]*model.UpCard, err error) {
	if len(mids) <= 0 {
		return
	}

	cards = make(map[int64]*model.UpCard)

	midInfoMap, err := s.card.MidUpInfoMap(ctx, mids)
	if err != nil {
		log.Error("card.MidUpInfoMap err, mids=%d, err=%v", mids, err.Error())
		return
	}
	midAccountsMap, err := s.getUpAccountsMap(ctx, mids)
	if err != nil {
		log.Error("GetCardInfoByMids getUpAccountsMap err, mids=%d, err=%v", mids, err.Error())
		return
	}
	midImagesMap, err := s.card.MidImagesMap(ctx, mids)
	if err != nil {
		log.Error("card.MidImagesMap err, mids=%d, err=%v", mids, err.Error())
		return
	}
	midVideosMap, err := s.getUpVideosMap(ctx, mids)
	if err != nil {
		log.Error("GetCardInfoByMids getUpVideosMap err, mids=%d, err=%v", mids, err.Error())
		return
	}

	for _, mid := range mids {
		card := &model.UpCard{
			Accounts: []*model.UpCardAccount{},
			Images:   []*model.UpCardImage{},
			Videos:   []*model.UpCardVideo{},
		}
		if data, ok := midInfoMap[mid]; ok {
			card.UpCardInfo = data
		}
		if data, ok := midAccountsMap[mid]; ok {
			card.Accounts = data
		}
		if data, ok := midImagesMap[mid]; ok {
			card.Images = data
		}
		if data, ok := midVideosMap[mid]; ok {
			card.Videos = data
		}

		cards[mid] = card
	}
	return
}

func (s *Service) getUpAccountsMap(ctx context.Context, mids []int64) (upAccountsMap map[int64][]*model.UpCardAccount, err error) {
	// 全网账号
	upAccountsMap, err = s.card.MidAccountsMap(ctx, mids)
	if err != nil {
		log.Error("card.MidAccountsMap err, mids=%d, err=%v", mids, err)
		return
	}
	var infosReply *accgrpc.InfosReply
	if infosReply, err = global.GetAccClient().Infos3(ctx, &accgrpc.MidsReq{Mids: mids, RealIp: metadata.String(ctx, metadata.RemoteIP)}); err != nil {
		return
	}
	if infosReply == nil || infosReply.Infos == nil {
		return
	}
	for mid, info := range infosReply.Infos {
		for _, account := range upAccountsMap[mid] {
			account.Picture = info.Face
		}
	}

	return
}

func (s *Service) getUpVideosMap(ctx context.Context, mids []int64) (upVideosMap map[int64][]*model.UpCardVideo, err error) {
	upVideosMap = make(map[int64][]*model.UpCardVideo)
	midAvidsMap, err := s.card.MidAvidsMap(ctx, mids)
	if err != nil {
		log.Error("card.MidAvidsMap err, mids=%d, err=%v", mids, err.Error())
		return
	}

	var allAvids []int64
	for _, avids := range midAvidsMap {
		allAvids = append(allAvids, avids...)
	}

	avidVideoMap, err := s.card.AvidVideoMap(ctx, allAvids)
	if err != nil {
		log.Error("card.AvidVideoMap err, mids=%d, err=%v", mids, err.Error())
		return
	}

	for mid, avids := range midAvidsMap {
		for _, avid := range avids {
			if video, ok := avidVideoMap[avid]; ok {
				upVideosMap[mid] = append(upVideosMap[mid], video)
			}
		}
	}

	return
}
