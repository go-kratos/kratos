package common

import (
	"context"
	"fmt"

	"go-common/app/admin/main/feed/model/common"
	showModel "go-common/app/admin/main/feed/model/show"
	account "go-common/app/service/main/account/model"
	"go-common/app/service/main/archive/api"
	seasondao "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"
	"go-common/library/ecode"
)

//CardPreview card preview
func (s *Service) CardPreview(c context.Context, cType string, id int64) (title string, err error) {
	var (
		accCard    *account.Card
		appActive  *showModel.AppActive
		eventTopic *showModel.EventTopic
		webCard    *showModel.SearchWebCard
		seaCards   map[int32]*seasondao.CardInfoProto
		arcCard    *api.Arc
	)
	switch cType {
	case common.CardPgc:
		v := []int32{int32(id)}
		if seaCards, err = s.pgcDao.CardsInfoReply(c, v); err != nil {
			return
		}
		if v, ok := seaCards[int32(id)]; ok {
			return v.Title, nil
		}
		return "", fmt.Errorf("无效pgc卡片ID(%d)", id)
	case common.CardAv:
		if arcCard, err = s.arcDao.Archive3(c, id); err != nil {
			if err.Error() == ecode.NothingFound.Error() {
				return "", fmt.Errorf("无效稿件ID(%d)", id)
			}
			return
		}
		return arcCard.Title, nil
	case common.CardUp:
		if accCard, err = s.accDao.Card3(c, id); err != nil {
			if err.Error() == ecode.MemberNotExist.Error() {
				return "", fmt.Errorf("无效up主ID(%d)", id)
			}
			return
		}
		return accCard.Name, nil
	case common.CardChannelTab:
		if appActive, err = s.showDao.AAFindByID(c, int64(id)); err != nil {
			return "", err
		}
		if appActive == nil {
			return "", fmt.Errorf("无效tab卡片ID(%d)", id)
		}
		return appActive.Name, nil
	case common.CardEventTopic:
		if eventTopic, err = s.showDao.ETFindByID(id); err != nil {
			return "", err
		}
		if eventTopic == nil {
			return "", fmt.Errorf("无效事件专题卡片ID(%d)", id)
		}
		return eventTopic.Title, nil
	case common.CardSearchWeb:
		if webCard, err = s.showDao.SWBFindByID(id); err != nil {
			return "", err
		}
		if webCard == nil {
			return "", fmt.Errorf("无效web卡片ID(%d)", id)
		}
		return webCard.Title, nil
	default:
		err = fmt.Errorf("参数错误")
		return "", err
	}
}
