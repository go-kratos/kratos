package gorpc

import (
	"context"

	"go-common/app/service/main/member/model"
	"go-common/library/net/rpc"
)

const (
	_appid = "account.service.member"
)

var (
	_noRes     = &struct{}{}
	_      RPC = &Service{}
)

// Service is a question service.
type Service struct {
	client *rpc.Client2
}

// New new a question service.
func New(c *rpc.ClientConfig) (s *Service) {
	s = &Service{}
	s.client = rpc.NewDiscoveryCli(_appid, c)
	return
}

//go:generate mockgen -source member.go  -destination mock.go -package member

// RPC is
type RPC interface {
	Exp(c context.Context, arg *model.ArgMid2) (res *model.LevelInfo, err error)
	Level(c context.Context, arg *model.ArgMid2) (res *model.LevelInfo, err error)
	Log(c context.Context, arg *model.ArgMid2) (res []*model.UserLog, err error)
	Stat(c context.Context, arg *model.ArgMid2) (res *model.ExpStat, err error)
	UpdateExp(c context.Context, arg *model.ArgAddExp) (err error)
	Base(c context.Context, arg *model.ArgMemberMid) (res *model.BaseInfo, err error)
	Bases(c context.Context, arg *model.ArgMemberMids) (res map[int64]*model.BaseInfo, err error)
	Member(c context.Context, arg *model.ArgMemberMid) (res *model.Member, err error)
	Members(c context.Context, arg *model.ArgMemberMids) (res map[int64]*model.Member, err error)
	NickUpdated(c context.Context, arg *model.ArgMemberMid) (res bool, err error)
	SetNickUpdated(c context.Context, arg *model.ArgMemberMid) (err error)
	SetOfficialDoc(c context.Context, arg *model.ArgOfficialDoc) (err error)
	SetName(c context.Context, arg *model.ArgUpdateUname) (err error)
	SetSign(c context.Context, arg *model.ArgUpdateSign) (err error)
	SetBirthday(c context.Context, arg *model.ArgUpdateBirthday) (err error)
	SetFace(c context.Context, arg *model.ArgUpdateFace) (err error)
	SetSex(c context.Context, arg *model.ArgUpdateSex) (err error)
	SetRank(c context.Context, arg *model.ArgUpdateRank) (err error)
	OfficialDoc(c context.Context, arg *model.ArgMid) (res *model.OfficialDoc, err error)
	Moral(c context.Context, arg *model.ArgMemberMid) (res *model.Moral, err error)
	MoralLog(c context.Context, arg *model.ArgMemberMid) (res []*model.UserLog, err error)
	AddMoral(c context.Context, arg *model.ArgUpdateMoral) (err error)
	BatchAddMoral(c context.Context, arg *model.ArgUpdateMorals) (res map[int64]int64, err error)
	AddUserMonitor(c context.Context, arg *model.ArgAddUserMonitor) error
	IsInMonitor(c context.Context, arg *model.ArgMid) (bool, error)
	AddPropertyReview(c context.Context, arg *model.ArgAddPropertyReview) error
	RealnameStatus(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameStatus, err error)
	RealnameApplyStatus(c context.Context, arg *model.ArgMemberMid) (res *model.RealnameApplyStatusInfo, err error)
	RealnameTelCapture(c context.Context, arg *model.ArgMemberMid) (err error)
	RealnameApply(c context.Context, arg *model.ArgRealnameApply) (err error)
}
