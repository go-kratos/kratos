package service

import (
	"context"
	"time"

	"go-common/app/interface/main/shorturl/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// ShortCache get short by cache
func (s *Service) ShortCache(c context.Context, short string) (su *model.ShortUrl, err error) {
	if su, err = s.shortd.Cache(c, short); err != nil {
		log.Error("s.shrtd.Cache(%s) error(%v)", short, err)
		return
	}
	if su != nil {
		return
	}
	if su, err = s.shortd.Short(c, short); err != nil {
		log.Error("s.shortd.Short(%s) error(%v)", short, err)
		return
	}
	//  add cache
	if su == nil {
		s.shortd.AddCache(func() {
			s.shortd.SetEmptyCache(context.TODO(), short)
		})
		return
	}
	s.shortd.AddCache(func() {
		s.shortd.SetCache(context.TODO(), su)
	})
	return
}

// ShortByID model.ShortUrl by short id
func (s *Service) ShortByID(c context.Context, id int64) (su *model.ShortUrl, err error) {
	su, err = s.shortd.ShortbyID(c, id)
	if err != nil || su == nil {
		log.Error("ShortByID.Get(%d) error(%v)", id, err)
		err = ecode.ShortURLNotFound
		return
	}
	return
}

// Add new url to db
func (s *Service) Add(c context.Context, mid int64, long string) (short string, err error) {
	su := &model.ShortUrl{
		Long:  long,
		Mid:   mid,
		State: model.StateNormal,
		CTime: xtime.Time(time.Now().Unix()),
	}
	ss := model.Generate(long) // detemine no repeat
	for _, str := range ss {
		su.Short = str
		if su.ID, err = s.shortd.InShort(c, su); err != nil {
			log.Error("s.shortd.InShort error(%v)", err)
			err = ecode.ServerErr
			return
		}
		if su.ID == 0 {
			var shortRes *model.ShortUrl
			// already exist
			if shortRes, err = s.shortd.Short(c, str); err == nil && shortRes != nil {
				if shortRes.Long == long {
					if _, err = s.shortd.UpdateState(c, shortRes.ID, mid, model.StateNormal); err != nil {
						break
					}
					short = shortRes.Short
					break
				} else {
					continue
				}
			}
		} else {
			short = str
			break
		}
	}
	if short == "" {
		err = ecode.ShortURLNotFound
		return
	}
	s.shortd.AddCache(func() {
		var su, err = s.shortd.Short(context.TODO(), short)
		if err == nil && su != nil {
			s.shortd.SetCache(context.TODO(), su)
		}
	})
	return
}

// ShortUpdate model.ShortUpdate
func (s *Service) ShortUpdate(c context.Context, id, mid int64, long string) (err error) {
	if _, err = s.shortd.ShortUp(c, id, mid, long); err != nil {
		log.Error("s.shortd.ShortUp error(%v)", err)
		return
	}
	s.shortd.AddCache(func() {
		var short *model.ShortUrl
		if short, err = s.shortd.ShortbyID(context.TODO(), id); err != nil {
			log.Error("s.shortd.ShortByID(%d) error(%v)", id, err)
			return
		}
		s.shortd.SetCache(context.TODO(), short)
	})
	return
}

// ShortDel model.ShortDel
func (s *Service) ShortDel(c context.Context, id, mid int64, now time.Time) (err error) {
	if _, err = s.shortd.UpdateState(c, id, mid, model.StateDelted); err != nil {
		log.Error("s.shortd.ShortDel error(%v)", err)
	}
	s.shortd.AddCache(func() {
		var short *model.ShortUrl
		if short, err = s.shortd.ShortbyID(context.TODO(), id); err != nil {
			log.Error("s.shortd.ShortByID(%d) error(%v)", id, err)
			return
		}
		s.shortd.SetCache(context.TODO(), short)
	})
	return
}

// ShortCount model.ShortCount count
func (s *Service) ShortCount(c context.Context, mid int64, long string) (count int, err error) {
	if count, err = s.shortd.ShortCount(c, mid, long); err != nil {
		log.Error("s.shortd.ShortState error(%v)", err)
	}
	return
}

// ShortLimit get short_url list
func (s *Service) ShortLimit(c context.Context, pn, ps int, mid int64, long string) (res []*model.ShortUrl, err error) {
	if res, err = s.shortd.ShortLimit(c, (pn-1)*ps, ps, mid, long); err != nil {
		log.Error("s.shortd.ShortLimit error(%v)", err)
	}
	return
}
