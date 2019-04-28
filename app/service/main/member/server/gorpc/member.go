package gorpc

import (
	"go-common/app/service/main/member/model"
	"go-common/library/net/rpc/context"
)

// Base get user base info.
func (r *RPC) Base(c context.Context, arg *model.ArgMemberMid, res *model.BaseInfo) (err error) {
	var v *model.BaseInfo
	if v, err = r.s.BaseInfo(c, arg.Mid); err == nil && res != nil {
		*res = *v
	}
	return
}

// Bases get batch user base info.
func (r *RPC) Bases(c context.Context, arg *model.ArgMemberMids, res *map[int64]*model.BaseInfo) (err error) {
	*res, err = r.s.BatchBaseInfo(c, arg.Mids)
	return
}

// Member get member info.
func (r *RPC) Member(c context.Context, arg *model.ArgMemberMid, res *model.Member) (err error) {
	var v *model.Member
	if v, err = r.s.Member(c, arg.Mid); err == nil && res != nil {
		*res = *v
	}
	return
}

// Members get batch member info.
func (r *RPC) Members(c context.Context, arg *model.ArgMemberMids, res *map[int64]*model.Member) (err error) {
	*res, err = r.s.Members(c, arg.Mids)
	return
}

// NickUpdated get nickUpdated.
func (r *RPC) NickUpdated(c context.Context, arg *model.ArgMemberMid, res *bool) (err error) {
	*res, err = r.s.NickUpdated(c, arg.Mid)
	return
}

// SetNickUpdated set nickUpdated.
func (r *RPC) SetNickUpdated(c context.Context, arg *model.ArgMemberMid, res *struct{}) (err error) {
	err = r.s.SetNickUpdated(c, arg.Mid)
	return
}

// SetOfficialDoc set official doc.
func (r *RPC) SetOfficialDoc(c context.Context, arg *model.ArgOfficialDoc, res *struct{}) (err error) {
	err = r.s.SetOfficialDoc(c, arg)
	return
}

// SetSex set sex.
func (r *RPC) SetSex(c context.Context, arg *model.ArgUpdateSex, res *struct{}) (err error) {
	err = r.s.SetSex(c, arg.Mid, arg.Sex)
	return
}

// SetName set name.
func (r *RPC) SetName(c context.Context, arg *model.ArgUpdateUname, res *struct{}) (err error) {
	err = r.s.SetName(c, arg.Mid, arg.Name)
	return
}

// SetFace set face.
func (r *RPC) SetFace(c context.Context, arg *model.ArgUpdateFace, res *struct{}) (err error) {
	err = r.s.SetFace(c, arg.Mid, arg.Face)
	return
}

// SetRank set rank.
func (r *RPC) SetRank(c context.Context, arg *model.ArgUpdateRank, res *struct{}) (err error) {
	err = r.s.SetRank(c, arg.Mid, arg.Rank)
	return
}

// SetBirthday set birthday.
func (r *RPC) SetBirthday(c context.Context, arg *model.ArgUpdateBirthday, res *struct{}) (err error) {
	err = r.s.SetBirthday(c, arg.Mid, arg.Birthday)
	return
}

// SetSign set sign.
func (r *RPC) SetSign(c context.Context, arg *model.ArgUpdateSign, res *struct{}) (err error) {
	err = r.s.SetSign(c, arg.Mid, arg.Sign)
	return
}

// OfficialDoc is.
func (r *RPC) OfficialDoc(c context.Context, arg *model.ArgMid, res *model.OfficialDoc) (err error) {
	var od *model.OfficialDoc
	if od, err = r.s.OfficialDoc(c, arg.Mid); err == nil && od != nil {
		*res = *od
	}
	return
}
