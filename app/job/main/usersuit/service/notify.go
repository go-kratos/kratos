package service

import (
	"context"
	"strconv"

	"go-common/app/job/main/usersuit/model"
	"go-common/library/log"
)

func (s *Service) accNotify(c context.Context, uid int64, Action string) (err error) {
	msg := &model.AccountNotify{UID: uid, Type: "update", Action: Action}
	if err = s.accountNotifyPub.Send(c, strconv.FormatInt(msg.UID, 10), msg); err != nil {
		log.Error("mid(%d) s.accountNotifyPub.Send(%+v,%s) error(%v)", msg.UID, msg, Action, err)
	}
	log.Info("mid(%d) s.accountNotifyPub.Send(%+v,%s)", msg.UID, msg, Action)
	return
}
