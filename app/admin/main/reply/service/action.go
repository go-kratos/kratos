package service

import (
	"context"
	"fmt"

	"go-common/app/admin/main/reply/model"
	thumbup "go-common/app/service/main/thumbup/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ActionCount return action exact count.
func (s *Service) ActionCount(c context.Context, rpID, oid, adminID int64, typ int32) (like, hate int32, err error) {
	rp, err := s.dao.Reply(c, oid, rpID)
	if err != nil {
		return
	}
	if rp == nil {
		err = ecode.ReplyNotExist
		return
	}
	like = rp.Like
	hate = rp.Hate
	return
}

// UpActionLike update action like.
func (s *Service) UpActionLike(c context.Context, rpID, oid, adminID int64, typ, count int32, remark string) (err error) {
	rp, err := s.dao.Reply(c, oid, rpID)
	if err != nil {
		return
	}
	if rp == nil {
		err = ecode.ReplyNotExist
		return
	}
	if _, err = s.thumbupClient.UpdateCount(c, &thumbup.UpdateCountReq{
		Business:   "reply",
		OriginID:   rp.Oid,
		MessageID:  rpID,
		LikeChange: int64(count),
		Operator:   fmt.Sprintf("%d", adminID),
	}); err != nil {
		log.Error("s.thumbupClient.UpdateCount (%d,%d,%d) failed!err:=%v", oid, rpID, int64(count), err)
		return
	}
	rp.Like += count
	if rp.Like < 0 {
		rp.Like = 0
	}
	if err = s.addReplyIndex(c, rp); err != nil {
		log.Error("s.addReplyIndex(%d,%d,%d) error(%v)", rp.ID, rp.Oid, rp.Type, err)
	}
	if err = s.dao.DelReplyCache(c, rp.ID); err != nil {
		log.Error("s.dao.DeleteReplyCache(%d) error(%v)", rp.ID, err)
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, rp.State)
	})
	return
}

// UpActionHate update action hate.
func (s *Service) UpActionHate(c context.Context, rpID, oid, adminID int64, typ, count int32, remark string) (err error) {
	rp, err := s.dao.Reply(c, oid, rpID)
	if err != nil {
		return
	}
	if rp == nil {
		err = ecode.ReplyNotExist
		return
	}
	if _, err = s.thumbupClient.UpdateCount(c, &thumbup.UpdateCountReq{
		Business:      "reply",
		OriginID:      rp.Oid,
		MessageID:     rpID,
		DislikeChange: int64(count),
		Operator:      fmt.Sprintf("%d", adminID),
	}); err != nil {
		log.Error("s.thumbupClient.UpdateCount (%d,%d,%d) failed!err:=%v", oid, rpID, int64(count), err)
		return
	}
	rp.Hate += count
	if rp.Hate < 0 {
		rp.Hate = 0
	}
	if err = s.addReplyIndex(c, rp); err != nil {
		log.Error("s.addReplyIndex(%d,%d,%d) error(%v)", rp.ID, rp.Oid, rp.Type, err)
	}
	if err = s.dao.DelReplyCache(c, rp.ID); err != nil {
		log.Error("s.dao.DeleteReplyCache(%d) error(%v)", rp.ID, err)
	}
	s.cache.Do(c, func(ctx context.Context) {
		s.pubSearchReply(ctx, map[int64]*model.Reply{rp.ID: rp}, rp.State)
	})
	return
}
