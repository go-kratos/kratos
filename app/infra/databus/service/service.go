package service

import (
	"context"
	"time"

	"go-common/app/infra/databus/conf"
	"go-common/app/infra/databus/dao"
	"go-common/app/infra/databus/model"
	"go-common/library/log"
	"go-common/library/stat/prom"
)

const (
	_authUpdateInterval = 1 * time.Minute
)

// Service service instance
type Service struct {
	dao *dao.Dao
	// auth
	auths map[string]*model.Auth
	// the auth of cluster changed
	clusterChan chan model.Auth
	// stats prom
	StatProm  *prom.Prom
	CountProm *prom.Prom
	TimeProm  *prom.Prom
}

// New new and return service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: dao.New(c),
		// cluster
		clusterChan: make(chan model.Auth, 5),
		// stats prom
		StatProm: prom.New().WithState("go_databus_state", []string{"role", "group", "topic", "partition"}),
		// count prom: count consumer and producer partition speed
		CountProm: prom.New().WithState("go_databus_counter", []string{"operation", "group", "topic"}),
		TimeProm:  prom.New().WithTimer("go_databus_timer", []string{"group"}),
	}
	s.fillAuth()
	go s.proc()
	return
}

// Ping check mysql connection
func (s *Service) Ping(c context.Context) error {
	return s.dao.Ping(c)
}

// Close close mysql connection
func (s *Service) Close() {
	if s.dao != nil {
		s.dao.Close()
	}
}

func (s *Service) proc() {
	for {
		s.fillAuth()
		time.Sleep(_authUpdateInterval)
	}
}

func (s *Service) fillAuth() (err error) {
	auths, err := s.dao.Auth(context.Background())
	if err != nil {
		log.Error("service.fillAuth error(%v)", err)
		return
	}
	var changed []*model.Auth
	// check cluster change event
	for group, nw := range auths {
		old, ok := s.auths[group]
		if !ok {
			continue
		}
		if old.Cluster != nw.Cluster {
			changed = append(changed, old)
			log.Info("cluster changed group(%s) topic(%s) oldCluster(%s) newCluster(%s)", old.Group, old.Topic, old.Cluster, nw.Cluster)
		}
	}
	s.auths = auths
	for _, ch := range changed {
		s.clusterChan <- *ch
	}
	return
}

// AuthApp check auth from cache
func (s *Service) AuthApp(group string) (a *model.Auth, ok bool) {
	a, ok = s.auths[group]
	return
}

// ClusterEvent return cluster change event
func (s *Service) ClusterEvent() (group <-chan model.Auth) {
	return s.clusterChan
}
