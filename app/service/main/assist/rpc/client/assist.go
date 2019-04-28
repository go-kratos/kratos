package client

import (
	"context"

	"go-common/library/net/rpc"

	model "go-common/app/service/main/assist/model/assist"
)

const (
	_assists         = "RPC.Assists"
	_assistIDs       = "RPC.AssistIDs"
	_assist          = "RPC.Assist"
	_assistExit      = "RPC.AssistExit"
	_addAssist       = "RPC.AddAssist"
	_delAssist       = "RPC.DelAssist"
	_assistLogs      = "RPC.AssistLogs"
	_assistLogInfo   = "RPC.AssistLogInfo"
	_assistLogAdd    = "RPC.AssistLogAdd"
	_assistLogCancel = "RPC.AssistLogCancel"
	_assistUps       = "RPC.AssistUps"
)

const (
	_appid = "archive.service.assist"
)

var (
	_noArg = &struct{}{}
)

// Service def
type Service struct {
	client *rpc.Client2
}

// New def
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return s
}

// Assists def
func (s *Service) Assists(c context.Context, arg *model.ArgAssists) (res []*model.Assist, err error) {
	err = s.client.Call(c, _assists, arg, &res)
	return
}

// AssistIDs def
func (s *Service) AssistIDs(c context.Context, arg *model.ArgAssists) (res []int64, err error) {
	err = s.client.Call(c, _assistIDs, arg, &res)
	return
}

// Assist def
func (s *Service) Assist(c context.Context, arg *model.ArgAssist) (res *model.AssistRes, err error) {
	res = new(model.AssistRes)
	err = s.client.Call(c, _assist, arg, &res)
	return
}

// AddAssist def
func (s *Service) AddAssist(c context.Context, arg *model.ArgAssist) (err error) {
	err = s.client.Call(c, _addAssist, arg, _noArg)
	return
}

// DelAssist def
func (s *Service) DelAssist(c context.Context, arg *model.ArgAssist) (err error) {
	err = s.client.Call(c, _delAssist, arg, _noArg)
	return
}

// AssistLogs def
func (s *Service) AssistLogs(c context.Context, arg *model.ArgAssistLogs) (res []*model.Log, err error) {
	err = s.client.Call(c, _assistLogs, arg, &res)
	return
}

// AssistLogInfo def
func (s *Service) AssistLogInfo(c context.Context, arg *model.ArgAssistLog) (res *model.Log, err error) {
	res = new(model.Log)
	err = s.client.Call(c, _assistLogInfo, arg, &res)
	return
}

// AssistLogAdd def
func (s *Service) AssistLogAdd(c context.Context, arg *model.ArgAssistLogAdd) (err error) {
	err = s.client.Call(c, _assistLogAdd, arg, _noArg)
	return
}

// AssistLogCancel def
func (s *Service) AssistLogCancel(c context.Context, arg *model.ArgAssistLog) (err error) {
	err = s.client.Call(c, _assistLogCancel, arg, _noArg)
	return
}

// AssistUps def
func (s *Service) AssistUps(c context.Context, arg *model.ArgAssistUps) (res *model.AssistUpsPager, err error) {
	res = new(model.AssistUpsPager)
	err = s.client.Call(c, _assistUps, arg, &res)
	return
}

// AssistExit def
func (s *Service) AssistExit(c context.Context, arg *model.ArgAssist) (err error) {
	err = s.client.Call(c, _assistExit, arg, _noArg)
	return
}
