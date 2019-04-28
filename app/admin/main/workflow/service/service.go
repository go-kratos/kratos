package service

import (
	"context"
	"sync"
	"time"

	"go-common/app/admin/main/workflow/dao"
	"go-common/app/admin/main/workflow/model"
	"go-common/library/conf/paladin"
	"go-common/library/log"
)

// Service is service.
type Service struct {
	closed bool
	// dao
	dao *dao.Dao

	wg    sync.WaitGroup
	jobCh chan func()

	// cache
	callbackCache  map[int8]*model.Callback
	busAttrCache   map[int8]*model.BusinessAttr
	reviewTypeName map[int64]string
	businessName   map[string]int8
	tagListCache   map[int8]map[int64]*model.TagMeta //map[bid]map[tid]*model.TagMeta
	roleCache      map[int8]map[int8]string          //map[bid]map[rid]name
	c              *paladin.Map                      // application.toml conf
}

// New is workflow-admin service implementation.
func New() (s *Service) {
	var ac = new(paladin.TOML)
	if err := paladin.Watch("application.toml", ac); err != nil {
		panic(err)
	}
	s = &Service{
		dao:           dao.New(),
		wg:            sync.WaitGroup{},
		callbackCache: make(map[int8]*model.Callback),
		c:             ac,
	}
	s.jobCh = make(chan func(), paladin.Int(s.c.Get("chanSize"), 1024))

	go s.cacheproc()

	s.wg.Add(1)
	go s.jobproc()

	s.loadReviewTypeName()
	return s
}

// Ping check server ok.
func (s *Service) Ping(c context.Context) (err error) {
	err = s.dao.Ping(c)
	return
}

// Close consumer close.
func (s *Service) Close() {
	s.dao.Close()

	close(s.jobCh)
	s.closed = true

	s.wg.Wait()
}

func (s *Service) task(f func()) {
	select {
	case s.jobCh <- f:
	default:
		log.Warn("Failed to enqueue a task due to job channel is full")
	}
}

// jobproc is a job queue for executing closure.
func (s *Service) jobproc() {
	defer s.wg.Done()

	for {
		f, ok := <-s.jobCh
		if !ok {
			log.Info("Stop job proc due to job channel is closed")
			return
		}
		f()
	}
}

// cacheproc goroutine
func (s *Service) cacheproc() {
	for {
		s.loadCallbacks()
		s.loadBusAttrs()
		s.loadTagList()
		s.loadBusinessRole()
		time.Sleep(5 * time.Minute)
	}
}

func (s *Service) loadCallbacks() {
	cbs, err := s.dao.AllCallbacks(context.Background())
	if err != nil {
		log.Error("s.dao.AllCallbacks() error(%v)", err)
		return
	}
	cbMap := make(map[int8]*model.Callback, len(cbs))
	for _, cb := range cbs {
		cbMap[cb.Business] = cb
	}
	s.callbackCache = cbMap
}

// loadBusAttrs returns attributes of business
func (s *Service) loadBusAttrs() (err error) {
	var busAttrs []*model.BusinessAttr
	if err = s.dao.ORM.Table("workflow_business_attr").Find(&busAttrs).Error; err != nil {
		log.Error("init business attr failed(%v)!", err)
		return
	}
	busAttrsMap := make(map[int8]*model.BusinessAttr, len(busAttrs))
	busName := make(map[string]int8)
	for _, attr := range busAttrs {
		busAttrsMap[int8(attr.BID)] = attr
		busName[attr.BusinessName] = int8(attr.BID)
	}
	s.busAttrCache = busAttrsMap
	s.businessName = busName
	return
}

func (s *Service) loadReviewTypeName() {
	s.reviewTypeName = make(map[int64]string)
	s.reviewTypeName[1] = "动画"
	s.reviewTypeName[2] = "电影"
	s.reviewTypeName[3] = "纪录片"
	s.reviewTypeName[4] = "国产动画"
	s.reviewTypeName[5] = "连续剧"
}

func (s *Service) loadTagList() (err error) {
	var tlc map[int8]map[int64]*model.TagMeta
	if tlc, err = s.dao.TagList(context.Background()); err != nil {
		log.Error("init manager tag failed(%v)!", err)
		return
	}
	s.tagListCache = tlc
	return
}

func (s *Service) loadBusinessRole() (err error) {
	var rc map[int8]map[int8]string
	if rc, err = s.dao.LoadRole(context.Background()); err != nil {
		log.Error("init business role failed(%v)!", err)
		return
	}
	s.roleCache = rc
	return
}
