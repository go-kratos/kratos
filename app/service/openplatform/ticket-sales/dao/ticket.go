package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go-common/app/service/openplatform/ticket-sales/model"
	"go-common/app/service/openplatform/ticket-sales/model/consts"
	"go-common/library/cache/redis"
	"go-common/library/log"
	"go-common/library/xstr"
)

//票号相关常量
const (
	sqlCountRefundTicket = "SELECT COUNT(*) FROM ticket WHERE oid IN (%s) AND refund_status=?"

	_selectTicketByOrderIDSQL  = "SELECT id, uid, oid, sid, price, src, type, status, qr, ref_id, sku_id, seat_id, seat, etime, refund_apply_time, ctime, mtime FROM ticket WHERE oid=?"
	_selectTicketByScreenIDSQL = "SELECT id, uid, oid, sid, price, src, type, status, qr, ref_id, sku_id, seat_id, seat, etime, refund_apply_time, ctime, mtime FROM ticket WHERE sid=? AND uid=? AND type != ?"
	_selectTicketByIDSQL       = "SELECT id, uid, oid, sid, price, src, type, status, qr, ref_id, sku_id, seat_id, seat, etime, refund_apply_time, ctime, mtime FROM ticket WHERE id IN (%s)"
	_updateTicketStatusSQL     = "UPDATE ticket set status=? WHERE id IN (%s)"
	_selectTicketSendBySendTID = "SELECT id, sid, send_tid, recv_tid, send_uid, recv_uid, recv_tel, status, ctime, mtime, oid FROM ticket_send WHERE send_tid IN (%s)"
	_selectTicketSendByRecvTID = "SELECT id, sid, send_tid, recv_tid, send_uid, recv_uid, recv_tel, status, ctime, mtime, oid FROM ticket_send WHERE recv_tid IN (%s)"
)

//rawRefundTicketCnt 统计用户已退票数
func (d *Dao) rawRefundTicketCnt(ctx context.Context, oids []int64) (cnt int64, err error) {
	lo := len(oids)
	if lo == 0 {
		return
	}
	q := fmt.Sprintf(sqlCountRefundTicket, strings.Repeat(",?", lo)[1:])
	a := make([]interface{}, lo+1)
	a[lo] = consts.TkStatusRefunded
	for k, v := range oids {
		a[k] = v
	}
	err = d.db.QueryRow(ctx, q, a...).Scan(&cnt)
	return
}

// CacheTicketsByOrderID 通过 order_id 获取 tickets 取缓存
func (d *Dao) CacheTicketsByOrderID(c context.Context, orderID int64) (res []*model.Ticket, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", fmt.Sprintf(model.CacheKeyOrderTickets, orderID)))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheTicketsByOrderID(%d) error(%v)", orderID, err)
		return
	}

	if err = json.Unmarshal(reply, &res); err != nil {
		log.Error("d.CacheTicketsByOrderID(%d) json.Unmarshal() error(%v)", orderID, err)
		return
	}
	return
}

// RawTicketsByOrderID 通过 order_id 获取 tickets
func (d *Dao) RawTicketsByOrderID(c context.Context, orderID int64) (res []*model.Ticket, err error) {
	rows, err := d.db.Query(c, _selectTicketByOrderIDSQL, orderID)
	if err != nil {
		log.Error("d.TicketsByOrderID(%v) d.db.Query() error(%v)", orderID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ticket := &model.Ticket{}
		if err = rows.Scan(&ticket.ID, &ticket.UID, &ticket.OID, &ticket.SID, &ticket.Price, &ticket.Src, &ticket.Type,
			&ticket.Status, &ticket.Qr, &ticket.RefID, &ticket.SkuID, &ticket.SeatID,
			&ticket.Seat, &ticket.ETime, &ticket.RefundApplyTime, &ticket.CTime, &ticket.MTime); err != nil {
			log.Error("d.TicketsByOrderID(%v) rows.Scan() error(%v)", orderID, err)
			return
		}
		res = append(res, ticket)
	}
	return
}

// AddCacheTicketsByOrderID 通过 order_id 获取 tickets 写缓存
func (d *Dao) AddCacheTicketsByOrderID(c context.Context, orderID int64, tickets []*model.Ticket) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	val, err := json.Marshal(tickets)
	if err != nil {
		log.Error("d.AddCacheTicketsByScreen() error(%v)", err)
		return
	}
	if _, err = conn.Do("SETEX", fmt.Sprintf(model.CacheKeyOrderTickets, orderID), model.RedisExpireOneDayTmp, val); err != nil {
		log.Error("d.AddCacheTicketsByOrderID(%d, %+v) error(%v)", orderID, tickets, err)
		return
	}
	return
}

// CacheTicketsByScreen 通过 screen_id user_id 获取 tickets
func (d *Dao) CacheTicketsByScreen(c context.Context, screenID int64, UID int64) (res []*model.Ticket, err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	reply, err := redis.Bytes(conn.Do("GET", fmt.Sprintf(model.CacheKeyScreenTickets, screenID, UID)))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheTicketsByScreen(%d, %d) error(%v)", screenID, UID, err)
		return
	}

	if err = json.Unmarshal(reply, &res); err != nil {
		log.Error("d.CacheTicketsByScreen(%d, %d) json.Unmarshal() error(%v)", screenID, UID, err)
		return
	}
	return
}

// RawTicketsByScreen 通过 screen_id user_id 获取 tickets
func (d *Dao) RawTicketsByScreen(c context.Context, screenID int64, UID int64) (res []*model.Ticket, err error) {
	rows, err := d.db.Query(c, _selectTicketByScreenIDSQL, screenID, UID, consts.TkTypeDistrib)
	if err != nil {
		log.Error("d.RawTicketsByScreen(%d, %d) error(%v)", screenID, UID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		ticket := &model.Ticket{}
		if err = rows.Scan(&ticket.ID, &ticket.UID, &ticket.OID, &ticket.SID, &ticket.Price, &ticket.Src, &ticket.Type,
			&ticket.Status, &ticket.Qr, &ticket.RefID, &ticket.SkuID, &ticket.SeatID,
			&ticket.Seat, &ticket.ETime, &ticket.RefundApplyTime, &ticket.CTime, &ticket.MTime); err != nil {
			log.Error("d.RawTicketsByScreen(%d, %d) rows.Scan() error(%v)", screenID, UID, err)
			return
		}
		res = append(res, ticket)
	}
	return
}

// AddCacheTicketsByScreen 通过 screen_id user_id 获取 tickets
func (d *Dao) AddCacheTicketsByScreen(c context.Context, screenID int64, tickets []*model.Ticket, UID int64) (err error) {
	conn := d.redis.Get(c)
	defer conn.Close()

	val, err := json.Marshal(tickets)
	if err != nil {
		log.Error("d.AddCacheTicketsByScreen() error(%v)", err)
		return
	}
	if _, err = conn.Do("SETEX", fmt.Sprintf(model.CacheKeyScreenTickets, screenID, UID), model.RedisExpireOneDayTmp, val); err != nil {
		log.Error("d.AddCacheTicketsByScreen(%d, %d, %+v) error(%v)", screenID, UID, tickets, err)
		return
	}
	return
}

// CacheTicketsByID .
func (d *Dao) CacheTicketsByID(c context.Context, ticketID []int64) (res map[int64]*model.Ticket, err error) {
	if len(ticketID) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	keys := make([]interface{}, 0)
	for _, ID := range ticketID {
		keys = append(keys, fmt.Sprintf(model.CacheKeyTicket, ID))
	}

	reply, err := redis.ByteSlices(conn.Do("MGET", keys...))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheTicketsByID(%v) conn.Do() error(%v)", ticketID, err)
		return
	}

	res = make(map[int64]*model.Ticket)
	for _, item := range reply {
		if len(item) == 0 {
			continue
		}
		ticket := &model.Ticket{}
		if err = json.Unmarshal(item, ticket); err != nil {
			log.Error("d.CacheTicketsByID(%v) json.Unmarshal(%s) error(%v)", ticketID, item, err)
			continue
		}
		res[ticket.ID] = ticket
	}
	return
}

// RawTicketsByID .
func (d *Dao) RawTicketsByID(c context.Context, ticketID []int64) (res map[int64]*model.Ticket, err error) {
	if len(ticketID) == 0 {
		return
	}

	rows, err := d.db.Query(c, fmt.Sprintf(_selectTicketByIDSQL, xstr.JoinInts(ticketID)))
	if err != nil {
		log.Error("d.RawTicketsByID(%v) d.db.Query() error(%v)", ticketID, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Ticket)
	for rows.Next() {
		ticket := &model.Ticket{}
		if err = rows.Scan(&ticket.ID, &ticket.UID, &ticket.OID, &ticket.SID, &ticket.Price, &ticket.Src, &ticket.Type,
			&ticket.Status, &ticket.Qr, &ticket.RefID, &ticket.SkuID, &ticket.SeatID,
			&ticket.Seat, &ticket.ETime, &ticket.RefundApplyTime, &ticket.CTime, &ticket.MTime); err != nil {
			log.Error("d.RawTicketsByID(%v) rows.Scan() error(%v)", ticketID, err)
			return
		}
		res[ticket.ID] = ticket
	}
	return
}

// AddCacheTicketsByID .
func (d *Dao) AddCacheTicketsByID(c context.Context, tickets map[int64]*model.Ticket) (err error) {
	if len(tickets) == 0 {
		return
	}
	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()

	args := make([]interface{}, 0)
	for ID, ticket := range tickets {
		var b []byte
		if b, err = json.Marshal(ticket); err != nil {
			log.Error("d.AddCacheTicketsByID(%+v) json.Marshal(%+v) error(%v)", tickets, ticket, err)
			continue
		}
		args = append(args, fmt.Sprintf(model.CacheKeyTicket, ID), b)
	}

	if err = conn.Send("MSET", args...); err != nil {
		log.Error("d.AddCacheTicketsByID(%+v) conn.Send() error(%v)", tickets, err)
		return
	}
	for ID := range tickets {
		conn.Send("EXPIRE", fmt.Sprintf(model.CacheKeyTicket, ID), model.RedisExpireTenMinTmp)
	}
	return
}

// UpdateTicketStatus 更新票状态
func (d *Dao) UpdateTicketStatus(c context.Context, status int16, ticketID ...int64) (err error) {
	if len(ticketID) == 0 {
		return
	}
	if _, err = d.db.Exec(c, fmt.Sprintf(_updateTicketStatusSQL, xstr.JoinInts(ticketID)), status); err != nil {
		log.Error("d.UpdateTicketStatus(%d, %v) error(%v)", status, ticketID, err)
	}
	return
}

// DelTicketCache 删除单张电子票全部 cache
func (d *Dao) DelTicketCache(c context.Context, tickets ...*model.Ticket) (err error) {
	if len(tickets) == 0 {
		return
	}
	var keys []interface{}
	for _, ticket := range tickets {
		keys = append(
			keys,
			fmt.Sprintf(model.CacheKeyOrderTickets, ticket.OID),
			fmt.Sprintf(model.CacheKeyScreenTickets, ticket.SID, ticket.UID),
			fmt.Sprintf(model.CacheKeyTicket, ticket.ID),
			fmt.Sprintf(model.CacheKeyTicketQr, ticket.Qr),
		)
	}
	if err = d.RedisDel(c, keys...); err != nil {
		log.Error("d.DelTicketCache() d.RedisDel(%v) error(%v)", keys, err)
	}
	return
}

// CacheTicketSend .
func (d *Dao) CacheTicketSend(c context.Context, IDs []int64, TIDType string) (res map[int64]*model.TicketSend, err error) {
	var cacheKey string
	switch TIDType {
	case consts.TIDTypeSend:
		cacheKey = model.CacheKeyTicketSend
	case consts.TIDTypeRecv:
		cacheKey = model.CacheKeyTicketRecv
	default:
		return
	}
	conn := d.redis.Get(c)
	defer conn.Close()

	keys := make([]interface{}, 0)
	for _, ID := range IDs {
		keys = append(keys, fmt.Sprintf(cacheKey, ID))
	}
	reply, err := redis.ByteSlices(conn.Do("MGET", keys...))
	if err != nil {
		if err == redis.ErrNil {
			err = nil
			return
		}
		log.Error("d.CacheTicketSend(%v, %s) conn.Do() error(%v)", IDs, TIDType, err)
		return
	}

	res = make(map[int64]*model.TicketSend)
	for _, item := range reply {
		if len(item) == 0 {
			continue
		}
		tmp := &model.TicketSend{}
		if err = json.Unmarshal(item, tmp); err != nil {
			log.Error("d.CacheTicketSend() json.Unmarshal(%s) error(%v)", item, err)
			continue
		}
		switch TIDType {
		case consts.TIDTypeSend:
			res[tmp.SendTID] = tmp
		case consts.TIDTypeRecv:
			res[tmp.RecvTID] = tmp
		default:
			return
		}
	}
	return
}

// RawTicketSend .
func (d *Dao) RawTicketSend(c context.Context, IDs []int64, TIDType string) (res map[int64]*model.TicketSend, err error) {
	if len(IDs) == 0 {
		return
	}
	var sql string
	switch TIDType {
	case consts.TIDTypeSend:
		sql = _selectTicketSendBySendTID
	case consts.TIDTypeRecv:
		sql = _selectTicketSendByRecvTID
	default:
		return
	}
	rows, err := d.db.Query(c, fmt.Sprintf(sql, xstr.JoinInts(IDs)))
	if err != nil {
		log.Error("d.RawTicketSend(%v, %s) d.db.Query() error(%v)", IDs, TIDType, err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.TicketSend)
	for rows.Next() {
		tmp := &model.TicketSend{}
		if err = rows.Scan(&tmp.ID, &tmp.SID, &tmp.SendTID, &tmp.RecvTID, &tmp.SendUID, &tmp.RecvUID, &tmp.RecvTel, &tmp.Status, &tmp.CTime, &tmp.MTime, &tmp.OID); err != nil {
			log.Error("d.RawTicketSend(%v, %s) rows.Scan() error(%v)", IDs, TIDType, err)
			return
		}
		switch TIDType {
		case consts.TIDTypeSend:
			res[tmp.SendTID] = tmp
		case consts.TIDTypeRecv:
			res[tmp.RecvTID] = tmp
		default:
			return
		}
	}
	return
}

// AddCacheTicketSend .
func (d *Dao) AddCacheTicketSend(c context.Context, tsMap map[int64]*model.TicketSend, TIDType string) (err error) {
	var cacheKey string
	switch TIDType {
	case consts.TIDTypeSend:
		cacheKey = model.CacheKeyTicketSend
	case consts.TIDTypeRecv:
		cacheKey = model.CacheKeyTicketRecv
	default:
		return
	}

	conn := d.redis.Get(c)
	defer func() {
		conn.Flush()
		conn.Close()
	}()

	var args []interface{}
	for ID, item := range tsMap {
		var b []byte
		if b, err = json.Marshal(item); err != nil {
			log.Error("d.AddCacheTicketSend(%v, %s), json.Marshal(%s) error(%v)", tsMap, TIDType, b, err)
			continue
		}
		args = append(args, fmt.Sprintf(cacheKey, ID), b)
	}

	if err = conn.Send("MSET", args...); err != nil {
		log.Error("d.AddCacheTicketsByID(%+v) conn.Send() error(%v)", tsMap, err)
		return
	}
	for ID := range tsMap {
		conn.Send("EXPIRE", fmt.Sprintf(cacheKey, ID), model.RedisExpireOneDayTmp)
	}
	return
}
