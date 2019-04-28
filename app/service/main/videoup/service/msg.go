package service

import (
	"context"

	accApi "go-common/app/service/main/account/api"
	"go-common/app/service/main/videoup/model/archive"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// sendMsg send msg
func (s *Service) sendMsg(arg *archive.ArgMsg) {
	s.worker.Add(func() {
		var err error
		mids, title, content, code := arg.MsgInfo(s.msgMap[arg.MSGID])
		log.Info("sendMsg aid (%d) mid(%+v) title(%s) content(%s) code(%s)", arg.Apply.ApplyAID, mids, title, content, code)
		if len(mids) == 0 || title == "" || content == "" || code == "" {
			log.Warn("sendMsg mid(%+v) title(%s) content(%s) code(%s) sth is empty!", mids, title, content, code)
			return
		}
		if err = s.msg.MutliSendSysMsg(context.Background(), mids, code, title, content, ""); err != nil {
			log.Error("sendMsg s.msg.MutliSendSysMsg(%+v,%s,%s,%s,%s) error(%+v)", mids, code, title, content, "", err)
		}
	})
}

func (s *Service) setMsgTypeMap() {
	s.msgMap = make(map[int]*archive.MSG, len(s.c.Property.MSG))
	for _, msg := range s.c.Property.MSG {
		s.msgMap[msg.MSGID] = msg
	}
}

func (s *Service) profile(c context.Context, mid int64) (p *accApi.ProfileStatReply, err error) {
	if p, err = s.accRPC.ProfileWithStat3(c, &accApi.MidReq{Mid: mid}); err != nil {
		p = nil
		log.Error("s.accRPC.ProfileWithStat3(%d) error(%v)", mid, err)
	}
	return
}

// Infos get user info by mids.
func (s *Service) Infos(c context.Context, mids []int64) (users map[int64]*accApi.Info, err error) {
	var (
		arg = &accApi.MidsReq{
			Mids: mids,
		}
		res *accApi.InfosReply
	)
	ip := metadata.String(c, metadata.RemoteIP)
	if res, err = s.accRPC.Infos3(c, arg); err != nil {
		log.Error("s.accRPC.Infos3() error(%v)|ip(%s)", err, ip)
		return
	}
	if res != nil {
		users = res.Infos
	}
	return
}
