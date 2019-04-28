package service

import (
	"context"

	"go-common/app/admin/main/card/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// CardsByGid get cards by gid.
func (s *Service) CardsByGid(c context.Context, gid int64) ([]*model.Card, error) {
	return s.dao.CardsByGid(c, gid)
}

// UpdateCardState update card state.
func (s *Service) UpdateCardState(c context.Context, req *model.ArgState) error {
	return s.dao.UpdateCardState(c, req.ID, req.State)
}

// DeleteCard delete card.
func (s *Service) DeleteCard(c context.Context, id int64) error {
	return s.dao.DeleteCard(c, id)
}

// UpdateGroupState update group state.
func (s *Service) UpdateGroupState(c context.Context, req *model.ArgState) error {
	return s.dao.UpdateGroupState(c, req.ID, req.State)
}

// DeleteGroup delete group.
func (s *Service) DeleteGroup(c context.Context, id int64) error {
	return s.dao.DeleteGroup(c, id)
}

// GroupList group list.
func (s *Service) GroupList(c context.Context, req *model.ArgQueryGroup) (res []*model.CardGroup, err error) {
	if res, err = s.dao.Groups(c, req); err != nil {
		return
	}
	if len(res) <= 0 {
		return
	}
	var cs []*model.Card
	if cs, err = s.dao.Cards(c); err != nil {
		return
	}
	tmp := make(map[int64][]*model.Card, len(res))
	for _, v := range cs {
		if len(tmp[v.GroupID]) <= 0 {
			tmp[v.GroupID] = []*model.Card{}
		}
		tmp[v.GroupID] = append(tmp[v.GroupID], v)
	}
	for _, v := range res {
		v.Cards = tmp[v.ID]
	}
	return
}

// CardOrderChange card order change.
func (s *Service) CardOrderChange(c context.Context, req *model.ArgIds) (err error) {
	var cs []*model.Card
	if cs, err = s.dao.CardsByIds(c, req.Ids); err != nil {
		return
	}
	if len(req.Ids) != len(cs) {
		err = ecode.CardIDNotFoundErr
		return
	}
	orders := make(map[int]*model.Card, len(cs))
	for i, v := range cs {
		orders[i] = v
	}
	us := []*model.Card{}
	for i, v := range req.Ids {
		if orders[i].ID != v {
			us = append(us, &model.Card{ID: v, OrderNum: orders[i].OrderNum})
		}
	}
	if len(us) > 0 {
		err = s.dao.BatchUpdateCardOrder(c, us)
	}
	return
}

// GroupOrderChange group order change.
func (s *Service) GroupOrderChange(c context.Context, req *model.ArgIds) (err error) {
	var cs []*model.CardGroup
	if cs, err = s.dao.GroupsByIds(c, req.Ids); err != nil {
		return
	}
	if len(req.Ids) != len(cs) {
		err = ecode.CardIDNotFoundErr
		return
	}
	orders := make(map[int]*model.CardGroup, len(cs))
	for i, v := range cs {
		orders[i] = v
	}
	us := []*model.CardGroup{}
	for i, v := range req.Ids {
		if orders[i].ID != v {
			us = append(us, &model.CardGroup{ID: v, OrderNum: orders[i].OrderNum})
		}
	}
	if len(us) > 0 {
		err = s.dao.BatchUpdateCardGroupOrder(c, us)
	}
	return
}

// AddCard add card.
func (s *Service) AddCard(c context.Context, req *model.AddCard) (err error) {
	var exist *model.Card
	if exist, err = s.dao.CardByName(req.Name); err != nil || exist != nil {
		return ecode.CardNameExistErr
	}
	var g errgroup.Group
	g.Go(func() (err error) {
		if req.CardURL, err = s.dao.Upload(c, "", req.CardFileType, req.CardBody, s.c.Bfs); err != nil {
			log.Error("d.Upload iconURL(%+v) error(%v)", req, err)
		}
		return
	})
	g.Go(func() (err error) {
		if req.BigCradURL, err = s.dao.Upload(c, "", req.BigCardFileType, req.BigCardBody, s.c.Bfs); err != nil {
			log.Error("d.Upload bigCardURL(%+v) error(%v)", req, err)
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	if req.CardURL == "" || req.BigCradURL == "" {
		err = ecode.CardFileUploadFaildErr
		return
	}
	var order int64
	if order, err = s.dao.MaxCardOrder(); err != nil {
		return
	}
	order++
	req.OrderNum = order
	err = s.dao.AddCard(req)
	return
}

// UpdateCard update card.
func (s *Service) UpdateCard(c context.Context, req *model.UpdateCard) (err error) {
	var g errgroup.Group
	g.Go(func() (err error) {
		if req.CardFileType != "" {
			if req.CardURL, err = s.dao.Upload(c, "", req.CardFileType, req.CardBody, s.c.Bfs); err != nil {
				log.Error("d.Upload iconURL(%+v) error(%v)", req, err)
			}
		}
		return
	})
	g.Go(func() (err error) {
		if req.BigCardFileType != "" {
			if req.BigCradURL, err = s.dao.Upload(c, "", req.BigCardFileType, req.BigCardBody, s.c.Bfs); err != nil {
				log.Error("d.Upload bigCardURL(%+v) error(%v)", req, err)
			}
		}
		return
	})
	if err = g.Wait(); err != nil {
		return
	}
	err = s.dao.UpdateCard(req)
	return
}

// AddGroup add group.
func (s *Service) AddGroup(c context.Context, req *model.AddGroup) (err error) {
	var exist *model.CardGroup
	if exist, err = s.dao.GroupByName(req.Name); err != nil {
		log.Error("s.dao.GroupByName(%+v) error(%v)", req, err)
		return
	}
	if exist != nil {
		return ecode.CardGroupNameExistErr
	}
	var order int64
	if order, err = s.dao.MaxGroupOrder(); err != nil {
		return
	}
	order++
	req.OrderNum = order
	return s.dao.AddGroup(c, req)
}

// UpdateGroup update group.
func (s *Service) UpdateGroup(c context.Context, req *model.UpdateGroup) error {
	return s.dao.UpdateGroup(c, req)
}
