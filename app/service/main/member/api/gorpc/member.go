package gorpc

import (
	"context"

	"go-common/app/service/main/member/model"
)

const (
	_Base           = "RPC.Base"
	_Bases          = "RPC.Bases"
	_Member         = "RPC.Member"
	_Members        = "RPC.Members"
	_UpdateExp      = "RPC.UpdateExp"
	_Exp            = "RPC.Exp"
	_Level          = "RPC.Level"
	_Log            = "RPC.Log"
	_Stat           = "RPC.Stat"
	_NickUpdated    = "RPC.NickUpdated"
	_SetNickUpdated = "RPC.SetNickUpdated"
	_Moral          = "RPC.Moral"
	_MoralLog       = "RPC.MoralLog"

	_SetOfficialDoc = "RPC.SetOfficialDoc"
	_OfficialDoc    = "RPC.OfficialDoc"

	_setName       = "RPC.SetName"
	_setSign       = "RPC.SetSign"
	_setRank       = "RPC.SetRank"
	_setFace       = "RPC.SetFace"
	_setSex        = "RPC.SetSex"
	_setBirthday   = "RPC.SetBirthday"
	_addMoral      = "RPC.AddMoral"
	_batchAddMoral = "RPC.BatchAddMoral"
)

// Exp rpc user exp.
func (s *Service) Exp(c context.Context, arg *model.ArgMid2) (res *model.LevelInfo, err error) {
	res = new(model.LevelInfo)
	err = s.client.Call(c, _Exp, arg, res)
	return
}

// Level user level.
func (s *Service) Level(c context.Context, arg *model.ArgMid2) (res *model.LevelInfo, err error) {
	res = new(model.LevelInfo)
	err = s.client.Call(c, _Level, arg, res)
	return
}

// Log user exp log.
func (s *Service) Log(c context.Context, arg *model.ArgMid2) (res []*model.UserLog, err error) {
	err = s.client.Call(c, _Log, arg, &res)
	return
}

// Stat user exp log.
func (s *Service) Stat(c context.Context, arg *model.ArgMid2) (res *model.ExpStat, err error) {
	err = s.client.Call(c, _Stat, arg, &res)
	return
}

// UpdateExp update user exp.
func (s *Service) UpdateExp(c context.Context, arg *model.ArgAddExp) (err error) {
	err = s.client.Call(c, _UpdateExp, arg, &_noRes)
	return
}

// Base get user base info.
func (s *Service) Base(c context.Context, arg *model.ArgMemberMid) (res *model.BaseInfo, err error) {
	err = s.client.Call(c, _Base, arg, &res)
	return
}

// Bases get batch base info.
func (s *Service) Bases(c context.Context, arg *model.ArgMemberMids) (res map[int64]*model.BaseInfo, err error) {
	err = s.client.Call(c, _Bases, arg, &res)
	return
}

// Member get the full information within member-service.
func (s *Service) Member(c context.Context, arg *model.ArgMemberMid) (res *model.Member, err error) {
	err = s.client.Call(c, _Member, arg, &res)
	return
}

// Members get batch the full information within member-service.
func (s *Service) Members(c context.Context, arg *model.ArgMemberMids) (res map[int64]*model.Member, err error) {
	err = s.client.Call(c, _Members, arg, &res)
	return
}

// NickUpdated get nickUpdated.
func (s *Service) NickUpdated(c context.Context, arg *model.ArgMemberMid) (res bool, err error) {
	err = s.client.Call(c, _NickUpdated, arg, &res)
	return
}

// SetNickUpdated set nickUpdated.
func (s *Service) SetNickUpdated(c context.Context, arg *model.ArgMemberMid) (err error) {
	err = s.client.Call(c, _SetNickUpdated, arg, &_noRes)
	return
}

// SetOfficialDoc set official doc.
func (s *Service) SetOfficialDoc(c context.Context, arg *model.ArgOfficialDoc) (err error) {
	err = s.client.Call(c, _SetOfficialDoc, arg, &_noRes)
	return
}

// SetName set name.
func (s *Service) SetName(c context.Context, arg *model.ArgUpdateUname) (err error) {
	err = s.client.Call(c, _setName, arg, &_noRes)
	return
}

// SetSign set sign.
func (s *Service) SetSign(c context.Context, arg *model.ArgUpdateSign) (err error) {
	err = s.client.Call(c, _setSign, arg, &_noRes)
	return
}

// SetBirthday set birthday.
func (s *Service) SetBirthday(c context.Context, arg *model.ArgUpdateBirthday) (err error) {
	err = s.client.Call(c, _setBirthday, arg, &_noRes)
	return
}

// SetFace set face.
func (s *Service) SetFace(c context.Context, arg *model.ArgUpdateFace) (err error) {
	err = s.client.Call(c, _setFace, arg, &_noRes)
	return
}

// SetSex set sex.
func (s *Service) SetSex(c context.Context, arg *model.ArgUpdateSex) (err error) {
	err = s.client.Call(c, _setSex, arg, &_noRes)
	return
}

// SetRank set rank.
func (s *Service) SetRank(c context.Context, arg *model.ArgUpdateRank) (err error) {
	err = s.client.Call(c, _setRank, arg, &_noRes)
	return
}

// OfficialDoc is.
func (s *Service) OfficialDoc(c context.Context, arg *model.ArgMid) (res *model.OfficialDoc, err error) {
	res = new(model.OfficialDoc)
	err = s.client.Call(c, _OfficialDoc, arg, res)
	return
}
