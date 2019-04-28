package service

import (
	"context"

	"go-common/app/service/main/tag/model"
	xsql "go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Like .
func (s *Service) Like(c context.Context, mid, oid, tid int64, typ int32, ip string) (err error) {
	var (
		action int32
		rows   int64
	)
	if action, err = s.action(c, mid, oid, tid, typ); err != nil {
		return
	}
	ra := &model.ResourceAction{
		Mid:    mid,
		Oid:    oid,
		Tid:    tid,
		Type:   typ,
		Action: model.UserActionLike,
	}
	switch action {
	case model.UserActionNormal: // 没有顶过
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.likeAction(c, mid, oid, tid, typ, 1); err != nil {
			return
		}
	case model.UserActionLike: // 顶 -> noaction
		ra.Action = model.UserActionNormal
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.likeAction(c, mid, oid, tid, typ, -1); err != nil {
			return
		}
	case model.UserActionHate: // 踩 -> 顶
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.likeAction(c, mid, oid, tid, typ, 1); err != nil {
			return
		}
		if err = s.hateAction(c, mid, oid, tid, typ, -1); err != nil {
			return
		}
	default:
		log.Warn("user Like service is no action!")
		return
	}
	s.cacheCh.Save(func() {
		s.dao.AddActionCache(context.Background(), mid, oid, tid, typ, ra.Action)
		rt, err := s.dao.Resource(context.Background(), oid, tid, typ)
		if err != nil {
			return
		}
		s.dao.AddResourceCache(context.Background(), rt)
		s.dao.DelResTagCache(context.Background(), oid, typ)
	})
	return
}

// Hate user hate under res.
func (s *Service) Hate(c context.Context, mid, oid, tid int64, typ int32, ip string) (err error) {
	var (
		action int32
		rows   int64
	)
	if action, err = s.action(c, mid, oid, tid, typ); err != nil {
		return
	}
	ra := &model.ResourceAction{
		Mid:    mid,
		Oid:    oid,
		Tid:    tid,
		Type:   typ,
		Action: model.UserActionHate,
	}
	switch action {
	case model.UserActionNormal: // 没有操作
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.hateAction(c, mid, oid, tid, typ, 1); err != nil {
			return
		}
	case model.UserActionHate: // 踩 -> no action
		ra.Action = model.UserActionNormal
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.hateAction(c, mid, oid, tid, typ, -1); err != nil {
			return
		}
	case model.UserActionLike: // 顶 -> 踩
		if rows, err = s.dao.AddAction(c, ra); err != nil || rows == 0 {
			return
		}
		if err = s.hateAction(c, mid, oid, tid, typ, 1); err != nil {
			return
		}
		if err = s.likeAction(c, mid, oid, tid, typ, -1); err != nil {
			return
		}
	default:
		log.Warn("user hate service is no action!")
		return
	}
	s.cacheCh.Save(func() {
		s.dao.AddActionCache(context.Background(), mid, oid, tid, typ, ra.Action)
		rt, err := s.dao.Resource(context.Background(), oid, tid, typ)
		if err != nil {
			return
		}
		s.dao.AddResourceCache(context.Background(), rt)
		s.dao.DelResTagCache(context.Background(), oid, typ)
	})
	return
}

func (s *Service) likeAction(c context.Context, mid, oid, tid int64, typ, number int32) (err error) {
	var (
		rt *model.Resource
		tx *xsql.Tx
	)
	if rt, err = s.dao.Resource(c, oid, tid, typ); err != nil {
		return
	}
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.TxUpResLike(tx, oid, tid, typ, number); err != nil {
		tx.Rollback()
		return
	}
	// 顶100 次自动锁定
	switch number {
	case 1:
		if int64(rt.Like+number) >= s.conf.Tag.LikeLimitToLock && !rt.Locked() {
			if _, err = s.dao.TxUpResAttr(tx, oid, tid, model.ResAttrLocked, model.AttrYes, typ); err != nil {
				tx.Rollback()
				return
			}
		}
	case -1:
		if int64(rt.Like+number) < s.conf.Tag.LikeLimitToLock && rt.Locked() {
			if _, err = s.dao.TxUpResAttr(tx, oid, tid, model.ResAttrLocked, model.AttrNo, typ); err != nil {
				tx.Rollback()
				return
			}
		}
	default:
		return ecode.RequestErr
	}
	err = tx.Commit()
	return
}

func (s *Service) hateAction(c context.Context, mid, oid, tid int64, typ, number int32) (err error) {
	var tx *xsql.Tx
	if tx, err = s.dao.BeginTran(c); err != nil {
		return
	}
	if _, err = s.dao.TxUpResHate(tx, oid, tid, typ, number); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// Action .
func (s *Service) Action(c context.Context, mid, oid, tid int64, typ int32) (action int32, err error) {
	return s.action(c, mid, oid, tid, typ)
}

func (s *Service) action(c context.Context, mid, oid, tid int64, typ int32) (action int32, err error) {
	var ok bool
	if ok, err = s.dao.ExpireAction(c, mid, typ); err != nil {
		return
	}
	if ok {
		action, err = s.dao.ActionCache(c, mid, oid, tid, typ)
		return
	}
	var acts []*model.ResourceAction
	if acts, err = s.dao.Actions(c, mid, typ); err != nil {
		return
	}
	if len(acts) == 0 {
		acts = append(acts, &model.ResourceAction{Oid: oid, Tid: -1, Action: model.UserActionNormal})
	} else {
		for _, act := range acts {
			if act.Type == typ && act.Oid == oid && act.Tid == tid {
				action = act.Action
				break
			}
		}
	}
	s.cache.Save(func() {
		s.dao.AddActionsCache(context.Background(), mid, typ, acts)
	})
	return
}

// ActionMap .
func (s *Service) ActionMap(c context.Context, mid, oid int64, typ int32, tids []int64) (res map[int64]int32, err error) {
	var (
		ok   bool
		acts []*model.ResourceAction
	)
	if ok, err = s.dao.ExpireAction(c, mid, typ); err != nil {
		return
	}
	if ok {
		res, err = s.dao.ActionsCache(c, mid, oid, typ, tids)
		return
	}
	if acts, err = s.dao.Actions(c, mid, typ); err != nil {
		return
	}
	if len(acts) == 0 {
		acts = append(acts, &model.ResourceAction{
			Oid:    oid,
			Tid:    -1,
			Action: model.UserActionNormal,
		})
	} else {
		res = make(map[int64]int32)
		for _, act := range acts {
			if act.Type == typ && act.Oid == oid {
				res[act.Tid] = act.Action
			}
		}
	}
	s.cache.Save(func() {
		s.dao.AddActionsCache(context.Background(), mid, typ, acts)
	})
	return
}
