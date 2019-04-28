package unicom

import (
	"context"
	"time"

	"go-common/library/ecode"
	"go-common/library/log"
)

// PrivilegePack
func (s *Service) Pack(c context.Context, usermob string, mid int64, now time.Time) (msg string, err error) {
	row := s.unicomInfo(c, usermob, now)
	u, ok := row[usermob]
	if !ok || u == nil {
		err = ecode.NothingFound
		msg = "该卡号尚未开通哔哩哔哩专属免流服务"
		return
	}
	var (
		result int64
	)
	userpacke, err := s.userPack(c, usermob)
	if err != nil {
		log.Error("s.userPack error(%v)", err)
		msg = "特权礼包领取失败"
		return
	}
	var (
		unicomPackOk bool
		userPackOk   bool
		user         map[int64]struct{}
	)
	if user, unicomPackOk = userpacke[usermob]; unicomPackOk {
		_, userPackOk = user[mid]
	}
	if !unicomPackOk && !userPackOk {
		result, err = s.dao.InPack(c, usermob, mid)
		if err != nil {
			log.Error("s.pack.InPack(%s, %s) error(%v)", usermob, mid, err)
		} else if result == 0 {
			err = ecode.RequestErr
			msg = "每张卡只能领取一次特权礼包哦"
			log.Error("s.pack.InPackc(%s,%s) error(%v) result==0", usermob, mid, err)
		} else {
			if err = s.live.Pack(c, mid, u.CardType); err != nil {
				msg = "礼包领取失败"
				log.Error("s.live.Pack mid error(%v)", mid, err)
			}
		}
	} else if unicomPackOk {
		err = ecode.NotModified
		msg = "每张卡只能领取一次特权礼包哦"
	} else if userPackOk {
		err = ecode.NotModified
		msg = "每个账号只能领取一次哦"
	} else {
		msg = "领取成功，特权礼包将在1~3个⼯工作⽇日内发放到您的账号"
	}
	return
}

func (s *Service) userPack(c context.Context, usermob string) (pack map[string]map[int64]struct{}, err error) {
	if pack, err = s.dao.Pack(c, usermob); err != nil {
		log.Error("s.pack.Pack usermob:%v error(%v)", usermob, err)
		return
	}
	return
}
