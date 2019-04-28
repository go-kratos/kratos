package service

import (
	"context"

	accmdl "go-common/app/service/main/account/api"
	feedmdl "go-common/app/service/main/reply-feed/api"

	"go-common/library/log"
)

// nextID return next reply id.
func (s *Service) nextID(c context.Context) (int64, error) {
	rpID, err := s.seqSrv.ID32(c, s.seqArg)
	if err != nil {
		log.Error("s.seqSrv.ID(%v) error(%v)", s.seqArg, err)
		return 0, err
	}
	return int64(rpID), err
}

func (s *Service) userInfo(c context.Context, mid int64) (*accmdl.Profile, error) {
	arg := &accmdl.MidReq{
		Mid: mid,
	}
	res, err := s.acc.Profile3(c, arg)
	if err != nil {
		log.Error("s.acc.UserInfo(%d) error(%v)", mid, err)
		return nil, err
	}
	return res.Profile, nil
}

func (s *Service) replyFeed(c context.Context, mid int64, pn, ps int) (res *feedmdl.ReplyRes, err error) {
	req := &feedmdl.ReplyReq{Mid: mid, Pn: int32(pn), Ps: int32(ps)}
	if res, err = s.feedClient.Reply(c, req); err != nil {
		log.Error("s.feedClient.Reply(%v) error(%v)", req, err)
		return nil, err
	}
	return
}

func (s *Service) replyHotFeed(c context.Context, mid, oid int64, tp, pn, ps int) (res *feedmdl.HotReplyRes, err error) {
	req := &feedmdl.HotReplyReq{Mid: mid, Oid: oid, Tp: int32(tp), Pn: int32(pn), Ps: int32(ps)}
	if res, err = s.feedClient.HotReply(c, req); err != nil {
		log.Error("s.feedClient.HotReply(%v) error(%v)", req, err)
		return nil, err
	}
	return
}
