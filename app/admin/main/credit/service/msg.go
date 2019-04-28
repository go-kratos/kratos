package service

import (
	"context"
	"time"

	blkmodel "go-common/app/admin/main/credit/model/blocked"
	"go-common/library/log"
)

func (s *Service) msgproc() {
	// NOTE: chan
	s.wg.Add(1)
	go func() {
		var (
			c              = context.TODO()
			sysMsg         *blkmodel.SysMsg
			ok             bool
			title, content string
		)
		defer s.wg.Done()
		for {
			if sysMsg, ok = <-s.MsgCh; !ok {
				log.Info("msgproc s.msgCh proc stop")
				return
			}
			if sysMsg == nil {
				select {
				case <-time.After(3 * time.Minute):
					continue
				case <-s.stop:
					return
				}
			}
			title, content = blkmodel.MsgInfo(sysMsg)
			if err := s.msgDao.SendSysMsg(c, sysMsg.MID, title, content); err != nil {
				log.Info("mid(%d) title(%s) content(%s) send sysMsg error(%v)", sysMsg.MID, title, content, err)
				continue
			}
			time.Sleep(100 * time.Millisecond)
		}
	}()
}
