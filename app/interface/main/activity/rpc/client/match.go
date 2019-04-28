package client

import (
	"context"

	matmdl "go-common/app/interface/main/activity/model/like"
	"go-common/library/net/rpc"
)

const (
	_matchs            = "RPC.Matchs"
	_subjectUp         = "RPC.SubjectUp"
	_likeUp            = "RPC.LikeUp"
	_addLikeCtimeCache = "RPC.AddLikeCtimeCache"
	_delLikeCtimeCache = "RPC.DelLikeCtimeCache"
	_actSubject        = "RPC.ActSubject"
	_actProtocol       = "RPC.ActProtocol"
)
const (
	_appid = "activity.service"
)

var (
	_noArg = &struct{}{}
)

// Service struct info.
type Service struct {
	client *rpc.Client2
}

// New new service instance and return.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// Matchs receive  matchs
func (s *Service) Matchs(c context.Context, arg *matmdl.ArgMatch) (res []*matmdl.Match, err error) {
	err = s.client.Call(c, _matchs, arg, &res)
	return
}

// SubjectUp update act_subject cache.
func (s *Service) SubjectUp(c context.Context, arg *matmdl.ArgSubjectUp) (err error) {
	return s.client.Call(c, _subjectUp, arg, _noArg)
}

// LikeUp update likes cache
func (s *Service) LikeUp(c context.Context, arg *matmdl.ArgLikeUp) (err error) {
	return s.client.Call(c, _likeUp, arg, _noArg)
}

//AddLikeCtimeCache add like ctime cache
func (s *Service) AddLikeCtimeCache(c context.Context, arg *matmdl.ArgLikeUp) (err error) {
	return s.client.Call(c, _addLikeCtimeCache, arg, _noArg)
}

// DelLikeCtimeCache del like ctime cache
func (s *Service) DelLikeCtimeCache(c context.Context, arg *matmdl.ArgLikeItem) (err error) {
	return s.client.Call(c, _delLikeCtimeCache, arg, _noArg)
}

// ActSubject get act subject item.
func (s *Service) ActSubject(c context.Context, arg *matmdl.ArgActSubject) (res *matmdl.SubjectItem, err error) {
	err = s.client.Call(c, _actSubject, arg, &res)
	return
}

// ActProtocol get protocol message
func (s *Service) ActProtocol(c context.Context, arg *matmdl.ArgActProtocol) (res *matmdl.SubProtocol, err error) {
	err = s.client.Call(c, _actProtocol, arg, &res)
	return
}
