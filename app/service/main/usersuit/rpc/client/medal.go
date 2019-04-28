package client

import (
	"context"

	"go-common/app/service/main/usersuit/model"
)

const (
	_medalHomeInfo       = "RPC.MedalHomeInfo"
	_medalUserInfo       = "RPC.MedalUserInfo"
	_medalInstall        = "RPC.MedalInstall"
	_medalPopup          = "RPC.MedalPopup"
	_medalMyInfo         = "RPC.MedalMyInfo"
	_medalAllInfo        = "RPC.MedalAllInfo"
	_medalGrant          = "RPC.MedalGrant"
	_medalActivated      = "RPC.MedalActivated"
	_medalActivatedMulti = "RPC.MedalActivatedMulti"
)

// MedalHomeInfo return user mdeal home info.
func (s *Service2) MedalHomeInfo(c context.Context, arg *model.ArgMid) (res []*model.MedalHomeInfo, err error) {
	err = s.client.Call(c, _medalHomeInfo, arg, &res)
	return
}

// MedalUserInfo return medal user info.
func (s *Service2) MedalUserInfo(c context.Context, arg *model.ArgMedalUserInfo) (res *model.MedalUserInfo, err error) {
	res = new(model.MedalUserInfo)
	err = s.client.Call(c, _medalUserInfo, arg, res)
	return
}

// MedalInstall install or uninstall medal.
func (s *Service2) MedalInstall(c context.Context, arg *model.ArgMedalInstall) (err error) {
	err = s.client.Call(c, _medalInstall, arg, _noRes)
	return
}

// MedalPopup return medal popup.
func (s *Service2) MedalPopup(c context.Context, arg *model.ArgMid) (res *model.MedalPopup, err error) {
	res = new(model.MedalPopup)
	err = s.client.Call(c, _medalPopup, arg, res)
	return
}

// MedalMyInfo return medal my info.
func (s *Service2) MedalMyInfo(c context.Context, arg *model.ArgMid) (res []*model.MedalMyInfos, err error) {
	err = s.client.Call(c, _medalMyInfo, arg, &res)
	return
}

// MedalAllInfo return medal all info.
func (s *Service2) MedalAllInfo(c context.Context, arg *model.ArgMid) (res *model.MedalAllInfos, err error) {
	err = s.client.Call(c, _medalAllInfo, arg, &res)
	return
}

// MedalGrant send a medal to user.
func (s *Service2) MedalGrant(c context.Context, arg *model.ArgMIDNID) (err error) {
	err = s.client.Call(c, _medalGrant, arg, _noRes)
	return
}

// MedalActivated get the user activated medal info.
func (s *Service2) MedalActivated(c context.Context, arg *model.ArgMid) (res *model.MedalInfo, err error) {
	err = s.client.Call(c, _medalActivated, arg, &res)
	return
}

// MedalActivatedMulti Multi get the user activated medal info(at most 50).
func (s *Service2) MedalActivatedMulti(c context.Context, arg *model.ArgMids) (res map[int64]*model.MedalInfo, err error) {
	err = s.client.Call(c, _medalActivatedMulti, arg, &res)
	return
}
