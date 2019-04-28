package service

import (
	"context"

	"go-common/library/log"
	"go-common/library/net/metadata"
)

const (
	_banned      = 1
	_unBindTel   = "unBindTel"
	_updateTel   = "updateTel"
	_updateEmail = "updateEmail"
	_updateUname = "updateUname"
	_updateFace  = "updateFace"
	_updateIDEN  = "updateIDEN"
)

// PurgeUser purge  user info
func (s *Service) PurgeUser(c context.Context, mid int64, action string) (err error) {
	ui, err := s.UserInfoAsyn(c, mid)
	if err != nil {
		log.Error("purge user(%d) failed", mid)
		return
	}
	if ui == nil || ui.State == _banned {
		log.Error("ui not fund or banned:%v", ui)
		return
	}
	var effect string
	switch action {
	case _unBindTel:
		effect = "解除绑定手机"
	case _updateTel:
		effect = "更新绑定手机"
	case _updateEmail:
		effect = "更新绑定邮箱"
	case _updateUname, _updateFace:
		effect = "修改账号资料"
	case _updateIDEN:
		effect = "完成实名认证"
	default:
		log.Error("unhandle action(%s) for mid(%d)", action, mid)
	}
	ip := metadata.String(c, metadata.RemoteIP)
	s.updatescore(func() {
		s.UpdateUserScore(context.TODO(), mid, ip, effect)
	})
	return
}
