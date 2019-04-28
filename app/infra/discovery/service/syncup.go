package service

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"go-common/app/infra/discovery/conf"
	"go-common/app/infra/discovery/dao"
	"go-common/app/infra/discovery/model"
	libenv "go-common/library/conf/env"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_fetchAllURL = "http://%s/discovery/fetch/all"
)

// syncUp populates the registry information from a peer eureka node.
func (s *Service) syncUp() (err error) {
	for _, node := range s.nodes.AllNodes() {
		if s.nodes.Myself(node.Addr) {
			continue
		}
		uri := fmt.Sprintf(_fetchAllURL, node.Addr)
		var res struct {
			Code int                          `json:"code"`
			Data map[string][]*model.Instance `json:"data"`
		}
		if err = s.client.Get(context.TODO(), uri, "", nil, &res); err != nil {
			log.Error("e.client.Get(%v) error(%v)", uri, err)
			continue
		}
		if res.Code != 0 {
			log.Error("service syncup from(%s) failed ", uri)
			continue
		}
		for _, is := range res.Data {
			for _, i := range is {
				s.tLock.RLock()
				appid, ok := s.tree[i.Treeid]
				s.tLock.RUnlock()
				if !ok || appid != i.Appid {
					s.tLock.Lock()
					s.tree[i.Treeid] = i.Appid
					s.tLock.Unlock()
				}
				s.registry.Register(i, i.LatestTimestamp)
			}
		}
		// NOTE: no return, make sure that all instances from other nodes register into self.
	}
	s.nodes.UP()
	return
}

func (s *Service) regSelf() context.CancelFunc {
	ctx, cancel := context.WithCancel(context.Background())
	now := time.Now().UnixNano()
	ins := &model.Instance{
		Region:   libenv.Region,
		Zone:     libenv.Zone,
		Env:      libenv.DeployEnv,
		Hostname: libenv.Hostname,
		Appid:    model.AppID,
		Addrs: []string{
			"http://" + s.c.BM.Inner.Addr,
		},
		Status:          model.InstanceStatusUP,
		RegTimestamp:    now,
		UpTimestamp:     now,
		LatestTimestamp: now,
		RenewTimestamp:  now,
		DirtyTimestamp:  now,
	}
	s.Register(ctx, ins, now, false)
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				arg := &model.ArgRenew{
					Appid:    model.AppID,
					Region:   libenv.Region,
					Zone:     libenv.Zone,
					Env:      libenv.DeployEnv,
					Hostname: libenv.Hostname,
				}
				if _, err := s.Renew(ctx, arg); err != nil && ecode.NothingFound.Equal(err) {
					s.Register(ctx, ins, now, false)
				}
			case <-ctx.Done():
				arg := &model.ArgCancel{
					Appid:    model.AppID,
					Region:   libenv.Region,
					Zone:     libenv.Zone,
					Env:      libenv.DeployEnv,
					Hostname: libenv.Hostname,
				}
				if err := s.Cancel(context.Background(), arg); err != nil {
					log.Error("s.Cancel(%+v) error(%v)", arg, err)
				}
				return
			}
		}
	}()
	return cancel
}

func (s *Service) nodesproc() {
	var (
		lastTs int64
	)
	for {
		arg := &model.ArgPolls{
			Appid:           []string{model.AppID},
			Region:          libenv.Region,
			Env:             libenv.DeployEnv,
			Hostname:        libenv.Hostname,
			LatestTimestamp: []int64{lastTs},
		}
		ch, _, err := s.registry.Polls(arg)
		if err != nil && err != ecode.NotModified {
			log.Error("s.registry(%v) error(%v)", arg, err)
			time.Sleep(time.Second)
			continue
		}
		apps := <-ch
		ins, ok := apps[model.AppID]
		if !ok || ins == nil {
			return
		}
		var (
			nodes []string
			zones = make(map[string][]string)
		)
		for _, in := range ins.Instances {
			for _, addr := range in.Addrs {
				u, err := url.Parse(addr)
				if err == nil && u.Scheme == "http" {
					if in.Zone == libenv.Zone {
						nodes = append(nodes, u.Host)
					} else if _, ok := s.c.Zones[in.Zone]; ok {
						zones[in.Zone] = append(zones[in.Zone], u.Host)
					}
				}
			}
		}
		lastTs = ins.LatestTimestamp
		log.Info("discovery changed nodes:%v zones:%v", nodes, zones)
		c := new(conf.Config)
		*c = *s.c
		c.Nodes = nodes
		c.Zones = zones
		s.nodes = dao.NewNodes(c)
		s.nodes.UP()
	}
}
