package service

import (
	"context"
	"sync/atomic"

	"go-common/app/service/main/reply-feed/api"
	"go-common/app/service/main/reply-feed/model"
)

func (s *Service) name(mid int64) (string, bool) {
	if mid == 0 {
		return model.DefaultSlotName, false
	}
	slot, ok := s.midMapping[mid]
	s.statisticsLock.RLock()
	defer s.statisticsLock.RUnlock()
	if ok && slot >= 0 && slot < model.SlotsNum {
		return s.statisticsStats[slot].Name, true
	}
	stat := s.statisticsStats[mid%model.SlotsNum]
	if stat.Name == model.DefaultSlotName || stat.State == model.StateInactive {
		return model.DefaultSlotName, false
	}
	return stat.Name, true
}

func (s *Service) incrHot(req *v1.HotReplyReq) {
	if req.Mid == 0 {
		return
	}
	if req.Pn == 1 {
		if req.Ps == 20 {
			// 用户点击更多热评
			atomic.AddUint32(&s.statisticsStats[req.Mid%model.SlotsNum].HotClick, 1)
		} else {
			// 点开评论区自带的5条热门评论
			return
		}
	}
	atomic.AddUint32(&s.statisticsStats[req.Mid%model.SlotsNum].HotView, 1)
}

func (s *Service) incrTotalView(req *v1.ReplyReq) {
	if req.Mid == 0 {
		return
	}
	atomic.AddUint32(&s.statisticsStats[req.Mid%model.SlotsNum].TotalView, 1)
}

func (s *Service) incrView(req *v1.ReplyReq) {
	if req.Mid == 0 {
		return
	}
	atomic.AddUint32(&s.statisticsStats[req.Mid%model.SlotsNum].View, 1)
}

// HotReply return hot reply
func (s *Service) HotReply(ctx context.Context, req *v1.HotReplyReq) (res *v1.HotReplyRes, err error) {
	var (
		start = (req.Pn - 1) * req.Ps
		end   = start + req.Ps - 1
		ok    bool
		count int
	)
	res = new(v1.HotReplyRes)
	// increment hot view and hot click count
	s.incrHot(req)
	if tp, exists := s.oidWhiteList[req.Oid]; exists && int32(tp) == req.Tp {
		res.Name = model.DefaultSlotName
		return
	}
	res.Name, ok = s.name(req.Mid)
	if !ok {
		return
	}
	if ok, err = s.dao.ExpireReplyZSetRds(ctx, res.Name, req.Oid, int(req.Tp)); err != nil {
		return
	}
	if ok {
		if res.RpIDs, err = s.dao.ReplyZSetRds(ctx, res.Name, req.Oid, int(req.Tp), int(start), int(end)); err != nil {
			return
		}
		if count, err = s.dao.CountReplyZSetRds(ctx, res.Name, req.Oid, int(req.Tp)); err != nil {
			return
		}
		res.Count = int32(count)
	} else {
		// s.eventProducer.Send(ctx, strconv.FormatInt(req.Oid, 10), &model.EventMsg{Action: "re_idx", Oid: req.Oid, Tp: int(req.Tp)})
		res.Name = model.DefaultSlotName
	}
	return
}

// Reply do increment reply view count.
func (s *Service) Reply(ctx context.Context, req *v1.ReplyReq) (res *v1.ReplyRes, err error) {
	res = new(v1.ReplyRes)
	// 用户点开评论区的次数
	if req.Pn == 1 {
		s.incrView(req)
	}
	// 用户在评论区总浏览次数
	s.incrTotalView(req)
	res.Name, _ = s.name(req.Mid)
	return
}
