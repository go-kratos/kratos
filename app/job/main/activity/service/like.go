package service

import (
	"context"
	"encoding/json"

	actmdl "go-common/app/interface/main/activity/model/like"
	l "go-common/app/job/main/activity/model/like"
	"go-common/library/log"
)

// AddLike like add data update cache .
func (s *Service) AddLike(c context.Context, addMsg json.RawMessage) (err error) {
	var (
		likeObj = new(l.Item)
	)
	if err = json.Unmarshal(addMsg, likeObj); err != nil {
		log.Error("AddLike json.Unmarshal(%s) error(%+v)", addMsg, err)
		return
	}
	if err = s.actRPC.LikeUp(c, &actmdl.ArgLikeUp{Lid: likeObj.ID}); err != nil {
		log.Error("s.actRPC.LikeUp(%d) error(%+v)", likeObj.ID, err)
		return
	}
	if err = s.actRPC.AddLikeCtimeCache(c, &actmdl.ArgLikeUp{Lid: likeObj.ID}); err != nil {
		log.Error("s.actRPC.AddLikeCtimeCache(%d) error(%+v)", likeObj.ID, err)
		return
	}
	log.Info("AddLike success s.actRPC.LikeUp(%d)", likeObj.ID)
	return
}

// UpLike update likes data update cahce
func (s *Service) UpLike(c context.Context, newMsg, oldMsg json.RawMessage) (err error) {
	var (
		likeObj = new(l.Item)
		oldObj  = new(l.Item)
	)
	if err = json.Unmarshal(newMsg, likeObj); err != nil {
		log.Error("UpLike json.Unmarshal(%s) error(%+v)", newMsg, err)
		return
	}
	if err = json.Unmarshal(oldMsg, oldObj); err != nil {
		log.Error("UpLike json.Unmarshal(%s) error(%+v)", oldMsg, err)
		return
	}
	if err = s.actRPC.LikeUp(c, &actmdl.ArgLikeUp{Lid: likeObj.ID}); err != nil {
		log.Error("s.actRPC.LikeUp(%d) error(%+v)", likeObj.ID, err)
		return
	}
	if oldObj.State != likeObj.State {
		if likeObj.State == 1 {
			//add ctime cache
			if err = s.actRPC.AddLikeCtimeCache(c, &actmdl.ArgLikeUp{Lid: likeObj.ID}); err != nil {
				log.Error("s.actRPC.AddLikeCtimeCache(%d) error(%+v)", likeObj.ID, err)
				return
			}
		} else {
			//del ctime cahce
			delItem := &actmdl.ArgLikeItem{
				ID:   likeObj.ID,
				Sid:  likeObj.Sid,
				Type: likeObj.Type,
			}
			if err = s.actRPC.DelLikeCtimeCache(c, delItem); err != nil {
				log.Error("s.actRPC.DelLikeCtimeCache(%v) error(%+v)", likeObj, err)
				return
			}
		}
	}
	log.Info("UpLike success s.actRPC.LikeUp(%d)", likeObj.ID)
	return
}

// DelLike delete like update cache
func (s *Service) DelLike(c context.Context, oldMsg json.RawMessage) (err error) {
	var (
		likeObj = new(l.Item)
	)
	if err = json.Unmarshal(oldMsg, likeObj); err != nil {
		log.Error("DelLike json.Unmarshal(%s) error(%+v)", oldMsg, err)
		return
	}
	if err = s.actRPC.LikeUp(c, &actmdl.ArgLikeUp{Lid: likeObj.ID}); err != nil {
		log.Error("s.actRPC.LikeUp(%d) error(%+v)", likeObj.ID, err)
		return
	}
	//del ctime cahce
	delItem := &actmdl.ArgLikeItem{
		ID:   likeObj.ID,
		Sid:  likeObj.Sid,
		Type: likeObj.Type,
	}
	if err = s.actRPC.DelLikeCtimeCache(c, delItem); err != nil {
		log.Error("s.actRPC.DelLikeCtimeCache(%v) error(%+v)", likeObj, err)
		return
	}
	log.Info("DelLike success s.actRPC.LikeUp(%d)", likeObj.ID)
	return
}

// upLikeContent .
func (s *Service) upLikeContent(c context.Context, upMsg json.RawMessage) (err error) {
	var (
		cont = new(l.Content)
	)
	if err = json.Unmarshal(upMsg, cont); err != nil {
		log.Error("upLikeContent json.Unmarshal(%s) error(%+v)", upMsg, err)
		return
	}
	if err = s.dao.SetLikeContent(c, cont.ID); err != nil {
		log.Error("s.dao.SetLikeContent(%d) error(%+v)", cont.ID, err)
	}
	log.Info("upLikeContent success s.dao.SetLikeContent(%d)", cont.ID)
	return
}
