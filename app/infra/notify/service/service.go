package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"go-common/app/infra/notify/conf"
	"go-common/app/infra/notify/dao"
	"go-common/app/infra/notify/model"
	"go-common/app/infra/notify/notify"
	"go-common/library/conf/env"
	"go-common/library/log"
)

// Service struct
type Service struct {
	c        *conf.Config
	dao      *dao.Dao
	plock    sync.RWMutex
	pubConfs map[string]*model.Pub
	subs     map[string]*notify.Sub
	pubs     map[string]*notify.Pub
}

// New init
func New(c *conf.Config) (s *Service) {
	s = &Service{
		c:        c,
		dao:      dao.New(c),
		pubConfs: make(map[string]*model.Pub),
		subs:     make(map[string]*notify.Sub),
		pubs:     make(map[string]*notify.Pub),
	}
	err := s.loadNotify()
	if err != nil {
		return
	}
	go s.notifyproc()
	go s.loadPub()
	go s.retryproc()
	return s
}

func (s *Service) loadPub() {
	for {
		pubs, err := s.dao.LoadPub(context.TODO())
		if err != nil {
			log.Error("load pub info err %v", err)
			time.Sleep(time.Minute)
			continue
		}
		ps := make(map[string]*model.Pub, len(pubs))
		for _, p := range pubs {
			ps[key(p.Group, p.Topic)] = p
		}
		s.pubConfs = ps
		time.Sleep(time.Minute)
	}
}

// TODO():auto reload and update.
func (s *Service) loadNotify() (err error) {
	watcher, err := s.dao.LoadNotify(context.TODO(), env.Zone)
	if err != nil {
		log.Error("load notify err %v", err)
		return
	}
	subs := make(map[string]*notify.Sub, len(watcher))
	for _, w := range watcher {
		if sub, ok := s.subs[key(w.Group, w.Topic)]; ok && !sub.Closed() && !sub.IsUpdate(w) {
			subs[key(w.Group, w.Topic)] = sub
		} else {
			n, err := s.newSub(w)
			if err != nil {
				log.Error("create notify topic(%s) group(%s) err(%v)", w.Topic, w.Group, err)
				continue
			}
			subs[key(w.Group, w.Topic)] = n
			log.Info("new sub %s %s", w.Group, w.Topic)
		}
	}
	// close subs not subscribe any more.
	for k, sub := range s.subs {
		if _, ok := subs[k]; !ok {
			sub.Close()
			log.Info("close sub not subscribe any %s", k)
		}
	}
	s.subs = subs
	return
}

func (s *Service) newSub(w *model.Watcher) (*notify.Sub, error) {
	var err error
	if w.Filter {
		w.Filters, err = s.dao.Filters(context.TODO(), w.ID)
		if err != nil {
			log.Error("s.dao.Filters err(%v)", err)
		}
	}
	return notify.NewSub(w, s.dao, s.c)
}

func (s *Service) notifyproc() {
	for {
		time.Sleep(time.Minute)
		s.loadNotify()
	}
}

func key(group, topic string) string {
	return fmt.Sprintf("%s_%s", group, topic)
}

// Ping Service
func (s *Service) Ping(c context.Context) (err error) {
	return s.dao.Ping(c)
}

// Close Service
func (s *Service) Close() {
	s.dao.Close()
}

func (s *Service) retryproc() {
	for {
		fs, err := s.dao.LoadFailBk(context.TODO())
		if err != nil {
			log.Error("s.loadFailBk err (%v)", err)
			time.Sleep(time.Minute)
			continue
		}
		for _, f := range fs {
			if n, ok := s.subs[key(f.Group, f.Topic)]; ok && !n.Closed() {
				n.AddRty(f.Msg, f.ID, f.Index)
			}
		}
		time.Sleep(time.Minute * 10)
	}
}
