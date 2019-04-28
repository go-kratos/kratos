package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/card/model"
	"go-common/library/log"
)

// Card get card info.
func (s *Service) Card(c context.Context, id int64) *model.Card {
	return s.cardmap[id]
}

// CardHots get card hots.
func (s *Service) CardHots(c context.Context) []*model.Card {
	return s.cardhots
}

// CardsByGid get card by gid.
func (s *Service) CardsByGid(c context.Context, gid int64) []*model.Card {
	return s.cardgidmap[gid]
}

func (s *Service) loadCard() (err error) {
	var (
		c   = context.Background()
		res []*model.Card
		ok  bool
	)
	if res, err = s.dao.EffectiveCards(c); err != nil {
		return
	}
	tmp := make(map[int64]*model.Card, len(res))
	htmp := []*model.Card{}
	cgtmp := map[int64][]*model.Card{}
	for _, v := range res {
		if _, ok = s.cardgroupmap[v.GroupID]; !ok {
			continue
		}
		v.CardTypeName = model.CardTypeNameMap[v.CardType]
		tmp[v.ID] = v
		if v.IsHot == model.CardIsHot {
			htmp = append(htmp, v)
		}
		if len(cgtmp[v.GroupID]) == 0 {
			cgtmp[v.GroupID] = []*model.Card{}
		}
		cgtmp[v.GroupID] = append(cgtmp[v.GroupID], v)
	}
	sort.Slice(htmp, func(i int, j int) bool {
		return htmp[i].Mtime.Time().After(htmp[j].Mtime.Time())
	})
	s.cardhots = htmp
	s.cardmap = tmp
	s.cardgidmap = cgtmp
	return
}

// loadcardproc load cards into memory.
func (s *Service) loadcardproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("loadbatchinfoproc panic(%v)", x)
			go s.loadcardproc()
		}
	}()
	for {
		if err := s.loadCard(); err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		time.Sleep(5 * time.Minute)
	}
}
