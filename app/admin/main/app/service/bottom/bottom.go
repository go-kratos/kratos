package bottom

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/app/conf"
	bottomdao "go-common/app/admin/main/app/dao/bottom"
	"go-common/app/admin/main/app/model/bottom"
	"go-common/library/log"
)

// Service bottom dao
type Service struct {
	dao *bottomdao.Dao
}

// New new a bottom dao
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: bottomdao.New(c),
	}
	return
}

// Bottoms select bottom all
func (s *Service) Bottoms(c context.Context) (res []*bottom.Bottom, err error) {
	if res, err = s.dao.Bottoms(c); err != nil {
		log.Error("s.dao.Bottoms error(%v)", err)
		return
	}
	return
}

// BottomByID select bottom by id
func (s *Service) BottomByID(c context.Context, id int64) (res *bottom.Bottom, err error) {
	if res, err = s.dao.BottomByID(c, id); err != nil {
		log.Error("s.dao.BottomByID error(%v)", err)
		return
	}
	return
}

// Insert insert
func (s *Service) Insert(c context.Context, a *bottom.Param, now time.Time) (err error) {
	if err = s.dao.Insert(c, a, now); err != nil {
		log.Error("s.dao.Insert error(%v)", err)
		return
	}
	return
}

// Update uodate
func (s *Service) Update(c context.Context, a *bottom.Param, now time.Time) (err error) {
	if err = s.dao.Update(c, a, now); err != nil {
		log.Error("s.dao.Update error(%v)", err)
		return
	}
	return
}

// Publish update state
func (s *Service) Publish(c context.Context, idsStr string, now time.Time) (err error) {
	var ids = map[int64]struct{}{}
	idsArr := strings.Split(idsStr, ",")
	for _, v := range idsArr {
		i, _ := strconv.ParseInt(v, 10, 64)
		ids[i] = struct{}{}
	}
	res, err := s.dao.Bottoms(c)
	if err != nil {
		log.Error("publish s.dao.Bottoms error(%v)", err)
		return
	}
	for _, v := range res {
		a := &bottom.Param{ID: v.ID}
		if _, ok := ids[v.ID]; ok {
			a.State = 1
		} else {
			a.State = 0
		}
		err = s.dao.UpdateByID(c, a, now)
		if err != nil {
			log.Error("publish s.dao.UpdateByID error(%v)", err)
			return
		}
	}
	return
}

//Delete delete by id
func (s *Service) Delete(c context.Context, id int64) (err error) {
	if s.dao.Delete(c, id); err != nil {
		log.Error("s.dao.Delete error(%v)", err)
		return
	}
	return
}
