package service

import (
	"context"

	"go-common/app/infra/discovery/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Register a new instance.
func (s *Service) Register(c context.Context, ins *model.Instance, latestTimestamp int64, replication bool) {
	s.registry.Register(ins, latestTimestamp)
	if ins.Treeid != 0 {
		s.tLock.RLock()
		appid, ok := s.tree[ins.Treeid]
		s.tLock.RUnlock()
		if !ok || appid != ins.Appid {
			s.tLock.Lock()
			s.tree[ins.Treeid] = ins.Appid
			s.tLock.Unlock()
		}
	}
	if !replication {
		s.nodes.Replicate(c, model.Register, ins, ins.Zone != s.env.Zone)
	}
}

// Renew marks the given instance of the given app name as renewed, and also marks whether it originated from replication.
func (s *Service) Renew(c context.Context, arg *model.ArgRenew) (i *model.Instance, err error) {
	i, ok := s.registry.Renew(arg)
	if !ok {
		err = ecode.NothingFound
		log.Error("renew appid(%s) hostname(%s) zone(%s) env(%s) error", arg.Appid, arg.Hostname, arg.Zone, arg.Env)
		return
	}
	if !arg.Replication {
		s.nodes.Replicate(c, model.Renew, i, arg.Zone != s.env.Zone)
		return
	}
	if arg.DirtyTimestamp > i.DirtyTimestamp {
		err = ecode.NothingFound
		return
	} else if arg.DirtyTimestamp < i.DirtyTimestamp {
		err = ecode.Conflict
	}
	return
}

// Cancel cancels the registration of an instance.
func (s *Service) Cancel(c context.Context, arg *model.ArgCancel) (err error) {
	i, ok := s.registry.Cancel(arg)
	if !ok {
		err = ecode.NothingFound
		log.Error("cancel appid(%s) hostname(%s) error", arg.Appid, arg.Hostname)
		return
	}
	if !arg.Replication {
		s.nodes.Replicate(c, model.Cancel, i, arg.Zone != s.env.Zone)
	}
	return
}

// FetchAll fetch all instances of all the department.
func (s *Service) FetchAll(c context.Context) (im map[string][]*model.Instance) {
	return s.registry.FetchAll()
}

// Fetchs fetch multi app by appids.
func (s *Service) Fetchs(c context.Context, arg *model.ArgFetchs) (is map[string]*model.InstanceInfo, err error) {
	is = make(map[string]*model.InstanceInfo, len(arg.Appid))
	for _, appid := range arg.Appid {
		i, err := s.registry.Fetch(arg.Zone, arg.Env, appid, 0, arg.Status)
		if err != nil {
			log.Error("Fetchs fetch appid(%s) err", err)
			continue
		}
		is[appid] = i
	}
	return
}

// Fetch fetch all instances by appid.
func (s *Service) Fetch(c context.Context, arg *model.ArgFetch) (info *model.InstanceInfo, err error) {
	var appid string
	if arg.Treeid != 0 {
		s.tLock.RLock()
		appid = s.tree[arg.Treeid]
		s.tLock.RUnlock()
	}
	if appid == "" {
		appid = arg.Appid
	}
	return s.registry.Fetch(arg.Zone, arg.Env, appid, 0, arg.Status)
}

// Polls hangs request and then write instances when that has changes, or return NotModified.
func (s *Service) Polls(c context.Context, arg *model.ArgPolls) (ch chan map[string]*model.InstanceInfo, new bool, err error) {
	var appids []string
	s.tLock.RLock()
	if len(arg.Treeid) > 0 {
		appids = make([]string, 0, len(arg.Treeid))
	}
	for _, tid := range arg.Treeid {
		appid := s.tree[tid]
		appids = append(appids, appid)
	}
	s.tLock.RUnlock()
	if len(appids) != 0 {
		arg.Appid = appids
	}
	return s.registry.Polls(arg)
}

// Polling get polling clients.
func (s *Service) Polling(c context.Context, arg *model.ArgPolling) (res []string, err error) {
	return s.registry.Polling(arg)
}

// DelConns delete conn of host in appid
func (s *Service) DelConns(arg *model.ArgPolls) {
	s.registry.DelConns(arg)
}

// Set set the status of instance by hostnames.
func (s *Service) Set(c context.Context, arg *model.ArgSet) (err error) {
	if ok := s.registry.Set(c, arg); !ok {
		err = ecode.NothingFound
		return
	}
	if !arg.Replication {
		s.nodes.ReplicateSet(c, arg, arg.Zone != s.env.Zone)
	}
	return
}

// Nodes get all nodes of discovery.
func (s *Service) Nodes(c context.Context) (nsi []*model.Node) {
	return s.nodes.Nodes()
}
