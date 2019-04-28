package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/card/model"
	"go-common/library/log"
)

// AllGroup all group.
func (s *Service) AllGroup(c context.Context, mid int64) (res *model.AllGroupResp, err error) {
	res = new(model.AllGroupResp)
	if res.UserCard, err = s.UserCard(c, mid); err != nil {
		return
	}
	ls := []*model.GroupInfo{}
	for _, v := range s.cardgroupmap {
		ls = append(ls, &model.GroupInfo{
			GroupID:   v.ID,
			GroupName: v.Name,
			OrderNum:  v.OrderNum,
			Cards:     s.cardgidmap[v.ID],
		})
	}
	sort.Slice(ls, func(i int, j int) bool {
		return ls[i].OrderNum > ls[j].OrderNum
	})
	res.List = ls
	return
}

func (s *Service) loadGroup() (err error) {
	var (
		c   = context.Background()
		res []*model.CardGroup
	)
	if res, err = s.dao.EffectiveGroups(c); err != nil {
		return
	}
	tmp := make(map[int64]*model.CardGroup, len(res))
	for _, v := range res {
		tmp[v.ID] = v
	}
	s.cardgroupmap = tmp
	return
}

// loadcardgroupproc load cards into memory.
func (s *Service) loadcardgroupproc() {
	defer func() {
		if x := recover(); x != nil {
			log.Error("loadcardgroupproc panic(%v)", x)
			go s.loadcardproc()
		}
	}()
	for {
		if err := s.loadGroup(); err != nil {
			time.Sleep(60 * time.Second)
			continue
		}
		time.Sleep(5 * time.Minute)
	}
}
