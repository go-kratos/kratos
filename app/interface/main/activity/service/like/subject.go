package like

import (
	"context"
	"time"

	ldao "go-common/app/interface/main/activity/dao/like"
	"go-common/app/interface/main/activity/model/like"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/sync/errgroup"
)

// SubjectInitialize act_subject data initialize .
func (s *Service) SubjectInitialize(c context.Context, minSid int64) (err error) {
	if minSid < 0 {
		minSid = 0
	}
	var actSub []*like.SubjectItem
	for {
		if actSub, err = s.dao.SubjectListMoreSid(c, minSid); err != nil {
			log.Error("dao.subjectListMoreSid(%d) error(%+v)", minSid, err)
			break
		}
		// empty slice or nil
		if len(actSub) == 0 {
			log.Info("SubjectInitialize end success")
			break
		}
		for _, sub := range actSub {
			item := sub
			if minSid < item.ID {
				minSid = item.ID
			}
			id := item.ID
			//the activity offline is stored with empty data
			if item.State != ldao.SubjectValidState {
				item = &like.SubjectItem{}
			}
			s.cache.Do(c, func(c context.Context) {
				s.dao.AddCacheActSubject(c, id, item)
			})
		}
	}
	s.cache.Do(c, func(c context.Context) {
		s.SubjectMaxIDInitialize(c)
	})
	return
}

// SubjectMaxIDInitialize Initialize act_subject max id data .
func (s *Service) SubjectMaxIDInitialize(c context.Context) (err error) {
	var actSub *like.SubjectItem
	if actSub, err = s.dao.SubjectMaxID(c); err != nil {
		log.Error(" s.dao.SubjectMaxID() error(%+v)", err)
		return
	}
	if actSub.ID >= 0 {
		if err = s.dao.AddCacheActSubjectMaxID(c, actSub.ID); err != nil {
			log.Error("s.dao.AddCacheActSubjectMaxID(%d) error(%v)", actSub.ID, err)
		}
	}
	return
}

// SubjectUp up act_subject cahce info .
func (s *Service) SubjectUp(c context.Context, sid int64) (err error) {
	var (
		actSub   *like.SubjectItem
		maxSubID int64
	)
	group, ctx := errgroup.WithContext(c)
	group.Go(func() (e error) {
		if actSub, e = s.dao.RawActSubject(ctx, sid); e != nil {
			log.Error("dao.RawActSubject(%d) error(%+v)", sid, e)
		}
		return
	})
	group.Go(func() (e error) {
		if maxSubID, e = s.dao.CacheActSubjectMaxID(ctx); e != nil {
			log.Error("dao.RawActSubject(%d) error(%v)", sid, e)
		}
		return
	})
	if err = group.Wait(); err != nil {
		log.Error("SubjectUp error(%v)", err)
		return
	}
	if actSub.ID == 0 || actSub.State != ldao.SubjectValidState {
		actSub = &like.SubjectItem{}
	}
	if maxSubID < sid {
		s.cache.Do(c, func(c context.Context) {
			s.dao.AddCacheActSubjectMaxID(context.Background(), sid)
		})
	}
	s.cache.Do(c, func(c context.Context) {
		s.dao.AddCacheActSubject(context.Background(), sid, actSub)
	})
	return
}

// SubjectLikeListInitialize Initialize likes list .
func (s *Service) SubjectLikeListInitialize(c context.Context, sid int64) (err error) {
	var (
		actSub *like.SubjectItem
		items  []*like.Item
		lid    = int64(0)
	)
	if actSub, err = s.dao.RawActSubject(c, sid); err != nil {
		log.Error("dao.RawActSubject(%d) error(%+v)", sid, err)
		return
	}
	if actSub.ID == 0 {
		log.Info("SubjectSLikeListInitialize end success")
		return
	}
	for {
		if items, err = s.dao.LikesBySid(c, lid, sid); err != nil {
			log.Error("dao.LikesBySid(%d,%d) error(%+v)", lid, sid, err)
			break
		}
		// empty slice or nil
		if len(items) == 0 {
			log.Info("SubjectSLikeListInitialize end success")
			break
		}
		//Initialize likes ctime cache
		cItems := items
		s.cache.Do(c, func(c context.Context) {
			s.dao.LikeListCtime(c, sid, cItems)
		})
		for _, val := range items {
			if lid < val.ID {
				lid = val.ID
			}
		}
	}
	return
}

// LikeActCountInitialize Initialize like_action cache data .
func (s *Service) LikeActCountInitialize(c context.Context, sid int64) (err error) {
	var (
		actSub      *like.SubjectItem
		items       []*like.Item
		lid         = int64(0)
		types       = make(map[int64]int)
		likeSumItem []*like.LidLikeSum
	)
	if actSub, err = s.dao.RawActSubject(c, sid); err != nil {
		log.Error("dao.RawActSubject(%d) error(%+v)", sid, err)
		return
	}
	if actSub.ID == 0 {
		log.Info("SubjectSLikeListInitialize end success")
		return
	}
	for {
		if items, err = s.dao.LikesBySid(c, lid, sid); err != nil {
			log.Error("dao.LikesBySid(%d,%d) error(%+v)", lid, sid, err)
			break
		}
		if len(items) == 0 {
			log.Info("SubjectSLikeListInitialize end success")
			break
		}
		lidList := make([]int64, 0, len(items))
		for _, val := range items {
			if lid < val.ID {
				lid = val.ID
			}
			lidList = append(lidList, val.ID)
			types[val.ID] = val.Type
		}
		if likeSumItem, err = s.dao.LikeActSums(c, sid, lidList); err != nil {
			log.Error(" s.dao.LikeActSums(%d,%v) error(%+v)", sid, lidList, err)
			return
		}
		if len(likeSumItem) == 0 {
			continue
		}
		lidLike := make(map[int64]int64, len(likeSumItem))
		for _, v := range likeSumItem {
			lidLike[v.Lid] = v.Likes
		}
		eg, errCtx := errgroup.WithContext(c)
		eg.Go(func() (e error) {
			e = s.dao.SetInitializeLikeCache(errCtx, sid, lidLike, types)
			return
		})
		eg.Go(func() (e error) {
			e = s.SetLikeActSum(errCtx, lidLike)
			return
		})
		if err = eg.Wait(); err != nil {
			log.Error("LikeActCountInitialize:eg.Wait() error(%+v)", err)
			return
		}
	}
	return
}

// SetLikeActSum set like_extend sum data
func (s *Service) SetLikeActSum(c context.Context, lidLikes map[int64]int64) (err error) {
	var (
		AddLids []*like.Extend
	)
	if len(lidLikes) == 0 {
		return
	}
	for k, v := range lidLikes {
		AddLids = append(AddLids, &like.Extend{Like: v, Lid: k})
	}
	_, err = s.BatchInsertLikeExtend(c, AddLids)
	return
}

// ActSubject .
func (s *Service) ActSubject(c context.Context, sid int64) (res *like.SubjectItem, err error) {
	if res, err = s.dao.ActSubject(c, sid); err != nil {
		return
	}
	if res == nil {
		err = ecode.NothingFound
	}
	return
}

// ActProtocol .
func (s *Service) ActProtocol(c context.Context, a *like.ArgActProtocol) (res *like.SubProtocol, err error) {
	res = new(like.SubProtocol)
	if res.SubjectItem, err = s.dao.ActSubject(c, a.Sid); err != nil {
		log.Error("s.dao.ActSubject() error(%+v)", err)
		return
	}
	if res.SubjectItem.ID == 0 {
		err = ecode.NothingFound
		return
	}
	now := time.Now().Unix()
	if int64(res.SubjectItem.Stime) <= now && int64(res.SubjectItem.Etime) >= now {
		if res.ActSubjectProtocol, err = s.dao.ActSubjectProtocol(c, a.Sid); err != nil {
			log.Error("s.dao.ActSubjectProtocol(%d) error(%+v)", a.Sid, err)
		}
	}
	return
}
