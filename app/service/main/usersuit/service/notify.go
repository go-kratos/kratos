package service

import (
	"context"
	"strconv"

	"go-common/app/service/main/usersuit/model"

	"github.com/pkg/errors"
)

func (s *Service) accNotify(c context.Context, mid int64, action string) (err error) {
	msg := &model.AccountNotify{UID: mid, Type: "updateUsersuit", Action: action}
	if err = s.accountNotifyPub.Send(c, strconv.FormatInt(msg.UID, 10), msg); err != nil {
		err = errors.Errorf("mid(%d) s.accountNotifyPub.Send(%+v) error(%+v)", msg.UID, msg, err)
	}
	return
}
