package client

import (
	"context"

	"go-common/app/service/main/resource/model"
	"go-common/library/net/rpc"
)

var (
	_noArg = &struct{}{}
)

const (
	_resourceAll   = "RPC.ResourceAll"
	_assignmentAll = "RPC.AssignmentAll"
	_defBanner     = "RPC.DefBanner"
	_resource      = "RPC.Resource"
	_resources     = "RPC.Resources"
	_banners       = "RPC.Banners"
	_pasterAPP     = "RPC.PasterAPP"
	_indexIcon     = "RPC.IndexIcon"
	_playerIcon    = "RPC.PlayerIcon"
	_cmtbox        = "RPC.Cmtbox"
	_sidebars      = "RPC.SideBars"
	_abtest        = "RPC.AbTest"
	_pasterCID     = "RPC.PasterCID"

	// app id
	_appid = "resource.service"
)

// Service is resource rpc client.
type Service struct {
	client *rpc.Client2
}

// New new a resource rpc client.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

// ResourceAll get all resource.
func (s *Service) ResourceAll(c context.Context) (res []*model.Resource, err error) {
	err = s.client.Call(c, _resourceAll, _noArg, &res)
	return
}

// AssignmentAll get all assignment.
func (s *Service) AssignmentAll(c context.Context) (res []*model.Assignment, err error) {
	err = s.client.Call(c, _assignmentAll, _noArg, &res)
	return
}

// DefBanner get default banner.
func (s *Service) DefBanner(c context.Context) (res *model.Assignment, err error) {
	res = new(model.Assignment)
	err = s.client.Call(c, _defBanner, _noArg, res)
	return
}

// Resource get resource.
func (s *Service) Resource(c context.Context, arg *model.ArgRes) (res *model.Resource, err error) {
	res = new(model.Resource)
	err = s.client.Call(c, _resource, arg, res)
	return
}

// Resources get resource.
func (s *Service) Resources(c context.Context, arg *model.ArgRess) (res map[int]*model.Resource, err error) {
	err = s.client.Call(c, _resources, arg, &res)
	return
}

// Banners get banners.
func (s *Service) Banners(c context.Context, arg *model.ArgBanner) (res *model.Banners, err error) {
	err = s.client.Call(c, _banners, arg, &res)
	return
}

// PasterAPP get paster.
func (s *Service) PasterAPP(c context.Context, arg *model.ArgPaster) (res *model.Paster, err error) {
	res = new(model.Paster)
	err = s.client.Call(c, _pasterAPP, arg, res)
	return
}

// IndexIcon get index icons.
func (s *Service) IndexIcon(c context.Context) (res map[string][]*model.IndexIcon, err error) {
	err = s.client.Call(c, _indexIcon, _noArg, &res)
	return
}

// PlayerIcon get palyer config.
func (s *Service) PlayerIcon(c context.Context) (res *model.PlayerIcon, err error) {
	res = new(model.PlayerIcon)
	err = s.client.Call(c, _playerIcon, _noArg, res)
	return
}

// Cmtbox get live box
func (s *Service) Cmtbox(c context.Context, arg *model.ArgCmtbox) (res *model.Cmtbox, err error) {
	err = s.client.Call(c, _cmtbox, arg, &res)
	return
}

// SideBars get side bar.
func (s *Service) SideBars(c context.Context) (res *model.SideBars, err error) {
	res = new(model.SideBars)
	err = s.client.Call(c, _sidebars, _noArg, res)
	return
}

// AbTest get abtest.
func (s *Service) AbTest(c context.Context, arg *model.ArgAbTest) (res map[string]*model.AbTest, err error) {
	err = s.client.Call(c, _abtest, arg, &res)
	return
}

// PasterCID get all Paster's cid.
func (s *Service) PasterCID(c context.Context) (res map[int64]int64, err error) {
	err = s.client.Call(c, _pasterCID, _noArg, &res)
	return
}
