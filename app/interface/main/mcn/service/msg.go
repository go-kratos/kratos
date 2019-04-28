package service

import (
	"context"

	"go-common/app/admin/main/mcn/model"
	"go-common/library/log"
)

// sendMsg send msg
func (s *Service) sendMsg(arg *model.ArgMsg) {
	s.worker.Add(func() {
		var err error
		mids, title, content, code := arg.MsgInfo(s.msgMap[arg.MSGType])
		if len(mids) == 0 || title == "" || content == "" || code == "" {
			log.Warn("mid(%+v) title(%s) content(%s) code(%s) sth is empty!", mids, title, content, code)
			return
		}
		if err = s.msg.MutliSendSysMsg(context.Background(), mids, code, title, content, ""); err != nil {
			log.Error("s.msg.MutliSendSysMsg(%+v,%s,%s,%s,%s) error(%+v)", mids, code, title, content, "", err)
		}
	})
}

func (s *Service) setMsgTypeMap() {
	s.msgMap = make(map[model.MSGType]*model.MSG, len(s.c.Property.MSG))
	for _, msg := range s.c.Property.MSG {
		s.msgMap[msg.MSGType] = msg
	}
}
