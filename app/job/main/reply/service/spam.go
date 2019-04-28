package service

import (
	"context"
	"encoding/json"

	accmdl "go-common/app/service/main/account/api"
	"go-common/library/ecode"
	"go-common/library/log"
)

type spamMessage struct {
	Mid  int64 `json:"mid"`
	IsUp bool  `json:"is_up"`
	Tp   int8  `json:"tp"`
}

func (s *Service) getLevel(c context.Context, mid int64) (cur int, err error) {
	arg := &accmdl.MidReq{Mid: mid}
	res, err := s.accSrv.Card3(c, arg)
	if err != nil {
		return
	}
	cur = int(res.Card.Level)
	return
}

func (s *Service) addRecReply(c context.Context, msg *consumerMsg) {
	var (
		d    spamMessage
		exp  int
		code = ecode.OK
	)
	if err := json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	count, err := s.spam.IncrReply(c, d.Mid, d.IsUp)
	if err != nil {
		log.Error("spam.IncrReply(%d) error(%v)", d.Mid, err)
		return
	}
	if d.IsUp && count >= 20 {
		exp = 5 * 60 // 5min
		code = ecode.ReplyDeniedAsCaptcha
	} else if count >= 5 {
		exp = 5 * 60 // 5min
		code = ecode.ReplyDeniedAsCaptcha
	}
	if code == ecode.OK {
		return
	}
	if err = s.spam.SetReplyRecSpam(c, d.Mid, code.Code(), exp); err == nil {
		log.Info("spam.SetReplyRecSpam(%d, %d, %d)", d.Mid, code, exp)
	} else {
		log.Error("spam.SetReplyRecSpam(%d, %d, %d), err (%v)", d.Mid, code, exp, err)
	}
}

func (s *Service) addDailyReply(c context.Context, msg *consumerMsg) {
	var d spamMessage
	if err := json.Unmarshal([]byte(msg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	ttl, err := s.spam.TTLDailyReply(c, d.Mid)
	if err != nil {
		log.Error("spam.TTLDailyReply(%d) error(%v)", d.Mid, err)
		return
	}
	count, err := s.spam.IncrDailyReply(c, d.Mid)
	if err != nil {
		log.Error("spam.IncrDailyReply(%d) error(%v)", d.Mid, err)
		return
	}
	if ttl == -2 || ttl == -1 {
		ttl = 24 * 60 * 60 // one day
		if err = s.spam.ExpireDailyReply(c, d.Mid, ttl); err != nil {
			log.Error("spam.ExpireDailyReply(%d) error(%v)", d.Mid, err)
		}
	}
	var code ecode.Codes
	// 23 BBQ  22 火鸟
	if d.Tp == 23 || d.Tp == 22 {
		if count <= 1000 {
			return
		}
		code = ecode.ReplyDeniedAsCD
	} else {
		lv, err := s.getLevel(c, d.Mid)
		if err != nil {
			log.Error("s.getLevel(%d) error(%v)", d.Mid, err)
			return
		}
		switch {
		case lv <= 1 && count < 25:
			return
		case lv == 2 && count < 250:
			return
		case lv == 3 && count < 300:
			return
		case lv == 4 && count < 400:
			return
		case lv >= 5 && count < 800:
			return
		}
		code = ecode.ReplyDeniedAsCaptcha
		if count >= 1000 {
			code = ecode.ReplyDeniedAsCD
		}
	}
	if err = s.spam.SetReplyDailySpam(c, d.Mid, code.Code(), ttl); err == nil {
		log.Info("spam.SetReplyDailySpam(%d, %d, %d)", d.Mid, code, ttl)
	} else {
		log.Error("spam.SetReplyDailySpam(%d, %d, %d) error(%v)", d.Mid, code, ttl, err)
	}
}

func (s *Service) recAct(c context.Context, cmsg *consumerMsg) {
	const _exp = 60
	var d spamMessage
	if err := json.Unmarshal([]byte(cmsg.Data), &d); err != nil {
		log.Error("json.Unmarshal() error(%v)", err)
		return
	}
	count, err := s.spam.IncrAct(c, d.Mid)
	if err != nil {
		log.Error("spam.IncUserRecAct(%d) error(%v)", d.Mid, err)
		return
	}
	if count < 15 {
		return
	}
	if err = s.spam.SetActionRecSpam(c, d.Mid, ecode.ReplyForbidAction.Code(), _exp); err == nil {
		log.Info("spam.SetActRecSpam(%d, %d, %d)", d.Mid, ecode.ReplyForbidAction, _exp)
	} else {
		log.Error("spam.SetActRecSpam(%d, %d, %d) error(%v)", d.Mid, ecode.ReplyForbidAction, _exp, err)
	}
}
