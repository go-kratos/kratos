package service

import (
	"context"
	"encoding/json"

	"go-common/app/service/main/workflow/model/sobot"
	"go-common/library/ecode"
	"go-common/library/log"
)

// SobotTicketInfo get ticket info
func (s *Service) SobotTicketInfo(c context.Context, ticketID int32) (res json.RawMessage, err error) {
	if res, err = s.sobot.SobotTicketInfo(c, ticketID); err != nil {
		log.Error("s.sobot.SobotAddTicket(%d) error(%v)", ticketID, err)
		return
	}
	return
}

// SobotTicketAdd add ticket
func (s *Service) SobotTicketAdd(c context.Context, tp *sobot.TicketParam) (err error) {
	if !tp.Check() {
		log.Error("s.SobotTicketAdd() params(%+v) error", tp)
		err = ecode.RequestErr
		return
	}
	if err = s.sobot.SobotAddTicket(c, tp); err != nil {
		log.Error("s.sobot.SobotAddTicket(%+v) error(%v)", tp, err)
		return
	}
	return
}

// SobotTicketModify modify ticket status
func (s *Service) SobotTicketModify(c context.Context, tp *sobot.TicketParam) (err error) {
	if !tp.CheckModify() {
		log.Error("s.SobotTicketModify() params(%+v) error", tp)
		err = ecode.RequestErr
		return
	}
	if err = s.sobot.SobotTicketModify(c, tp); err != nil {
		log.Error("s.sobot.SobotTicketModify(%+v) error(%v)", tp, err)
		return
	}
	return
}

// SobotReplyAdd add reply to sobot
func (s *Service) SobotReplyAdd(c context.Context, rp *sobot.ReplyParam) (err error) {
	if !rp.Check() {
		log.Error("s.SobotReplyAdd() params(%+v) error", rp)
		err = ecode.RequestErr
		return
	}
	if err = s.sobot.SobotAddReply(c, rp); err != nil {
		log.Error("s.sobot.SobotAddReply(%+v) error(%v)", rp, err)
		return
	}
	return
}
