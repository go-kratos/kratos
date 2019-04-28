package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/job/main/thumbup/conf"
	"go-common/app/job/main/thumbup/dao"
	"go-common/app/job/main/thumbup/model"
	"go-common/library/queue/databus"
	"go-common/library/queue/databus/databusutil"
	"go-common/library/sync/pipeline"
)

const (
	_retryTimes = 3
)

// Service .Service
type Service struct {
	c              *conf.Config
	dao            *dao.Dao
	waiter         *sync.WaitGroup
	merge          *pipeline.Pipeline
	likeGroup      *databusutil.Group
	itemLikesGroup *databusutil.Group
	userLikesGroup *databusutil.Group
	// businessMap
	businessMap   map[string]*model.Business
	businessIDMap map[int64]*model.Business
	// for 拜年祭
	statMerge *statMerge
}

type statMerge struct {
	Business string
	Target   int64
	Sources  map[int64]bool
}

// New create service instance and return.
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:   c,
		dao: dao.New(c),
		// waitGroup
		waiter: new(sync.WaitGroup),
	}
	if c.StatMerge != nil {
		s.statMerge = &statMerge{
			Business: c.StatMerge.Business,
			Target:   c.StatMerge.Target,
			Sources:  make(map[int64]bool),
		}
		for _, id := range c.StatMerge.Sources {
			s.statMerge.Sources[id] = true
		}
	}
	s.loadBusiness()
	s.initLikeGroup(c)
	s.initItemLikesGroup(c)
	s.initUserLikesGroup(c)
	s.initCountsMerge()
	go s.loadBusinessproc()
	return
}

func (s *Service) loadBusinessproc() {
	for {
		s.loadBusiness()
		time.Sleep(time.Minute * 5)
	}
}

func (s *Service) initItemLikesGroup(c *conf.Config) {
	ig := databusutil.NewGroup(c.ItemLikesDatabusutil, databus.New(c.Databus.ItemLikes).Messages())
	ig.New = newItemLikeMsg
	ig.Split = itemLikesSplit
	ig.Do = s.itemLikesDo
	ig.Start()
	s.itemLikesGroup = ig
}

func (s *Service) initUserLikesGroup(c *conf.Config) {
	ug := databusutil.NewGroup(c.UserLikesDatabusutil, databus.New(c.Databus.UserLikes).Messages())
	ug.New = newUserLikeMsg
	ug.Split = userLikesSplit
	ug.Do = s.userLikesDo
	ug.Start()
	s.userLikesGroup = ug
}

func (s *Service) initLikeGroup(c *conf.Config) {
	lg := databusutil.NewGroup(c.LikeDatabusutil, databus.New(c.Databus.Like).Messages())
	lg.New = newLikeMsg
	lg.Split = likeSplit
	lg.Do = s.likeDo
	lg.Start()
	s.likeGroup = lg
}

func (s *Service) initCountsMerge() {
	s.merge = pipeline.NewPipeline(s.c.Merge)
	s.merge.Split = s.countsSplit
	s.merge.Do = s.updateCountsDo
	s.merge.Start()
}

func (s *Service) loadBusiness() {
	var (
		err      error
		business []*model.Business
	)
	businessMap := make(map[string]*model.Business)
	businessIDMap := make(map[int64]*model.Business)
	for {
		if business, err = s.dao.Business(context.TODO()); err != nil {
			time.Sleep(time.Second)
			continue
		}
		for _, b := range business {
			businessMap[b.Name] = b
			businessIDMap[b.ID] = b
		}
		s.businessMap = businessMap
		s.businessIDMap = businessIDMap
		return
	}
}

// Ping .
func (s *Service) Ping(ctx context.Context) error {
	return s.dao.Ping(ctx)
}

// Close .
func (s *Service) Close() {
	s.waiter.Wait()
	s.itemLikesGroup.Close()
	s.userLikesGroup.Close()
	s.likeGroup.Close()
	s.dao.Close()
}
