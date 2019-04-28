package service

import (
	"context"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// sendMsg send msg
func (s *Service) sendMsg(c context.Context, arg *model.ArgMsg) error {
	var err error
	mids, title, content, code := arg.MsgInfo(s.msgMap[arg.MSGType])
	if len(mids) == 0 || title == "" || content == "" || code == "" {
		log.Warn("mid(%+v) title(%s) content(%s) code(%s) sth is empty!", mids, title, content, code)
		return nil
	}
	if err = s.msg.MutliSendSysMsg(c, mids, code, title, content, metadata.String(c, metadata.RemoteIP)); err != nil {
		log.Error("s.msg.MutliSendSysMsg(%+v,%s,%s,%s,%s) error(%+v)", mids, code, title, content, metadata.String(c, metadata.RemoteIP), err)
	}
	return err
}

func (s *Service) setMsgTypeMap() {
	s.msgMap = make(map[model.MSGType]*model.MSG, len(s.c.Property.MSG))
	for _, msg := range s.c.Property.MSG {
		s.msgMap[msg.MSGType] = msg
	}
}
