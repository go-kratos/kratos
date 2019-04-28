package server

import (
	"go-common/app/service/main/usersuit/model"
	"go-common/library/net/rpc/context"
)

// MedalHomeInfo return user mdeal home info.
func (r *RPC) MedalHomeInfo(c context.Context, arg *model.ArgMid, res *[]*model.MedalHomeInfo) (err error) {
	*res, err = r.s.MedalHomeInfo(c, arg.Mid)
	return
}

// MedalUserInfo return medal user info.
func (r *RPC) MedalUserInfo(c context.Context, arg *model.ArgMedalUserInfo, res *model.MedalUserInfo) (err error) {
	var mui *model.MedalUserInfo
	if mui, err = r.s.MedalUserInfo(c, arg.Mid, arg.IP); err == nil && mui != nil {
		*res = *mui
	}
	return
}

// MedalInstall install or uninstall medal.
func (r *RPC) MedalInstall(c context.Context, arg *model.ArgMedalInstall, res *struct{}) (err error) {
	err = r.s.MedalInstall(c, arg.Mid, arg.Nid, arg.IsActivated)
	return
}

// MedalPopup return medal popup.
func (r *RPC) MedalPopup(c context.Context, arg *model.ArgMid, res *model.MedalPopup) (err error) {
	var mp *model.MedalPopup
	if mp, err = r.s.MedalPopup(c, arg.Mid); err == nil && mp != nil {
		*res = *mp
	}
	return
}

// MedalMyInfo return medal my info.
func (r *RPC) MedalMyInfo(c context.Context, arg *model.ArgMid, res *[]*model.MedalMyInfos) (err error) {
	*res, err = r.s.MedalMyInfo(c, arg.Mid)
	return
}

// MedalAllInfo return medal all info.
func (r *RPC) MedalAllInfo(c context.Context, arg *model.ArgMid, res *model.MedalAllInfos) (err error) {
	var mai *model.MedalAllInfos
	if mai, err = r.s.MedalAllInfo(c, arg.Mid); err == nil && mai != nil {
		*res = *mai
	}
	return
}

// MedalGrant send a medal to user.
func (r *RPC) MedalGrant(c context.Context, arg *model.ArgMIDNID, res *struct{}) (err error) {
	err = r.s.MedalGet(c, arg.MID, arg.NID)
	return
}

// MedalActivated get the user activated medal info.
func (r *RPC) MedalActivated(c context.Context, arg *model.ArgMid, res *model.MedalInfo) (err error) {
	var ma *model.MedalInfo
	if ma, err = r.s.MedalActivated(c, arg.Mid); err == nil && ma != nil {
		*res = *ma
	}
	return
}

// MedalActivatedMulti Multi get get the user activated medal info(at most 50).
func (r *RPC) MedalActivatedMulti(c context.Context, arg *model.ArgMids, res *map[int64]*model.MedalInfo) (err error) {
	*res, err = r.s.MedalActivatedMulti(c, arg.Mids)
	return
}
