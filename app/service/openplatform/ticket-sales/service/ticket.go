package service

import (
	"context"
	"time"

	"go-common/app/service/openplatform/ticket-sales/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/log"
)

// TicketView 电子票详情
func (s *Service) TicketView(ctx context.Context, req *v1.TicketViewRequest) (res *v1.TicketViewResponse, err error) {
	res = new(v1.TicketViewResponse)

	var tickets []*model.Ticket
	if req.OrderID > 0 {
		tickets, err = s.dao.TicketsByOrderID(ctx, req.OrderID)
	} else if req.ScreenID > 0 && req.UID > 0 {
		tickets, err = s.dao.TicketsByScreen(ctx, req.ScreenID, req.UID)
	} else if len(req.ID) > 0 {
		var ticketMap map[int64]*model.Ticket
		ticketMap, err = s.dao.TicketsByID(ctx, req.ID)
		for _, ticket := range ticketMap {
			tickets = append(tickets, ticket)
		}
	} else {
		return
	}
	if err != nil {
		log.Error("s.TicketView() error(%v)", err)
		return
	}

	var expireTicketIDs []int64
	var expireTickets []*model.Ticket
	for _, ticket := range tickets {
		tmp := &v1.TicketItem{}
		tmp.ID = ticket.ID
		tmp.UID = ticket.UID
		tmp.OID = ticket.OID
		tmp.SID = ticket.SID
		tmp.Price = ticket.Price
		tmp.Src = ticket.Src
		tmp.Type = ticket.Type
		tmp.Status = ticket.Status
		tmp.Qr = ticket.Qr
		tmp.RefID = ticket.RefID
		tmp.SkuID = ticket.SkuID
		tmp.SeatID = ticket.SeatID
		tmp.Seat = ticket.Seat
		tmp.RefundApplyTime = ticket.RefundApplyTime
		tmp.ETime = ticket.ETime
		tmp.CTime = ticket.CTime
		tmp.MTime = ticket.MTime
		res.Tickets = append(res.Tickets, tmp)

		// 如果票过期了但是状态是正常需要修改状态
		// 这里过期时间取 etime 字段，之前 php 代码取的场次的结束时间，这俩时间有脚本会同步
		if ticket.Status == consts.TkStatusUnchecked && int64(ticket.ETime) <= time.Now().Unix() {
			ticket.Status = consts.TkStatusExpired
			expireTickets = append(expireTickets, ticket)
			expireTicketIDs = append(expireTicketIDs, ticket.ID)
		}
	}

	if len(expireTicketIDs) > 0 {
		if err = s.dao.UpdateTicketStatus(ctx, consts.TkStatusExpired, expireTicketIDs...); err != nil {
			log.Error("s.TicketView() s.dao.UpdateTicketStatus() error(%v)", err)
			return
		}
		s.dao.DelTicketCache(ctx, expireTickets...)
	}
	return
}

// TicketSend 电子票赠送信息
func (s *Service) TicketSend(ctx context.Context, req *v1.TicketSendRequest) (res *v1.TicketSendResponse, err error) {
	ticketSends := make(map[int64]*model.TicketSend)
	if len(req.SendTID) > 0 {
		ticketSends, err = s.dao.TicketSend(ctx, req.SendTID, consts.TIDTypeSend)
	} else if len(req.RecvTID) > 0 {
		ticketSends, err = s.dao.TicketSend(ctx, req.RecvTID, consts.TIDTypeRecv)
	}

	if err != nil {
		log.Error("s.TicketSend(%+v) s.dao.TicketSend() error(%v)", req, err)
		return
	}

	res = &v1.TicketSendResponse{}
	for _, ticketSend := range ticketSends {
		item := &v1.TicketSendItem{}
		item.ID = ticketSend.ID
		item.SID = ticketSend.SID
		item.SendTID = ticketSend.SendTID
		item.RecvTID = ticketSend.RecvTID
		item.SendUID = ticketSend.SendUID
		item.RecvUID = ticketSend.RecvUID
		item.RecvTel = ticketSend.RecvTel
		item.Status = ticketSend.Status
		item.CTime = ticketSend.CTime
		item.MTime = ticketSend.MTime
		item.OID = ticketSend.OID
		res.TicketSends = append(res.TicketSends, item)
	}
	return
}
