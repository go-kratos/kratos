package member

import (
	"context"

	"go-common/app/service/main/usersuit/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// MedalHomeInfo return user mdeal home info.
func (s *Service) MedalHomeInfo(c context.Context, mid int64) (res []*model.MedalHomeInfo, err error) {
	var arg = &model.ArgMid{Mid: mid}
	res, err = s.usRPC.MedalHomeInfo(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalHomeInfo(%d) error (%v)", mid, err)
		return
	}
	return
}

// MedalUserInfo return medal user info.
func (s *Service) MedalUserInfo(c context.Context, mid int64) (res *model.MedalUserInfo, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	var arg = &model.ArgMedalUserInfo{Mid: mid, IP: ip}
	res, err = s.usRPC.MedalUserInfo(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalUserInfo(%d) error (%v)", mid, err)
		return
	}
	return
}

// MedalInstall install or uninstall medal.
func (s *Service) MedalInstall(c context.Context, mid, nid int64, isActivated int8) (err error) {
	var arg = &model.ArgMedalInstall{Mid: mid, Nid: nid, IsActivated: isActivated}
	err = s.usRPC.MedalInstall(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalInstall(mid:%d nid:%d isActivated:%d) error (%v)", mid, nid, isActivated, err)
		return
	}
	return
}

// MedalPopup return medal popup.
func (s *Service) MedalPopup(c context.Context, mid int64) (res *model.MedalPopup, err error) {
	var arg = &model.ArgMid{Mid: mid}
	res, err = s.usRPC.MedalPopup(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalPopup(mid:%d) error (%v)", mid, err)
		return
	}
	return
}

// MedalMyInfo return medal my info.
func (s *Service) MedalMyInfo(c context.Context, mid int64) (res []*model.MedalMyInfos, err error) {
	var arg = &model.ArgMid{Mid: mid}
	res, err = s.usRPC.MedalMyInfo(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalMyInfo(mid:%d) error (%v)", mid, err)
		return
	}
	return
}

// MedalAllInfo return medal all info.
func (s *Service) MedalAllInfo(c context.Context, mid int64) (res *model.MedalAllInfos, err error) {
	var arg = &model.ArgMid{Mid: mid}
	res, err = s.usRPC.MedalAllInfo(c, arg)
	if err != nil {
		log.Error("s.medalRPC.MedalAllInfo(mid:%d) error (%v)", mid, err)
		return
	}
	return
}
