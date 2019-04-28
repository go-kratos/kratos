package service

import (
	"context"
	"strconv"

	"go-common/app/admin/main/usersuit/model"
	"go-common/library/log"
)

func (s *Service) accNotify(c context.Context, uid int64, action string) (err error) {
	msg := &model.AccountNotify{UID: uid, Type: "update", Action: action}
	if err = s.accountNotifyPub.Send(c, strconv.FormatInt(msg.UID, 10), msg); err != nil {
		log.Error("mid(%d) s.accountNotifyPub.Send(%+v) error(%v)", msg.UID, msg, err)
	}
	return
}
