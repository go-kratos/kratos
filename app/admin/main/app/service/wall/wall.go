package wall

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/app/conf"
	walldao "go-common/app/admin/main/app/dao/wall"
	"go-common/app/admin/main/app/model/wall"
	"go-common/library/log"
)

// Service wall service.
type Service struct {
	dao *walldao.Dao
}

// New new a wall service.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: walldao.New(c),
	}
	return
}

// Walls select All
func (s *Service) Walls(c context.Context) (res []*wall.Wall, err error) {
	if res, err = s.dao.Walls(c); err != nil {
		log.Error("s.dao.Walls error(%v)", err)
		return
	}
	return
}

// WallByID by id
func (s *Service) WallByID(c context.Context, id int64) (res *wall.Wall, err error) {
	if res, err = s.dao.WallByID(c, id); err != nil {
		log.Error("s.dao.WallByID error(%v)", err)
		return
	}
	return
}

// UpdateWall update wall
func (s *Service) UpdateWall(c context.Context, a *wall.Param, now time.Time) (err error) {
	if err = s.dao.Update(c, a, now); err != nil {
		log.Error("s.dao.Update(%v)", err)
		return
	}
	return
}

// Insert insert
func (s *Service) Insert(c context.Context, a *wall.Param, now time.Time) (err error) {
	if err = s.dao.Insert(c, a, now); err != nil {
		log.Error("s.dao.Insert error(%v)", err)
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
	res, err := s.dao.Walls(c)
	if err != nil {
		log.Error("publish s.dao.Walls error(%v)", err)
		return
	}
	for _, v := range res {
		a := &wall.Param{ID: v.ID}
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
