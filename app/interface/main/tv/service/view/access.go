package view

import (
	"context"

	"go-common/app/interface/main/tv/dao/account"
	"go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
)

// checkAceess check user Aceess
func (s *Service) checkAceess(c context.Context, mid, aid int64, state, access int, ak, ip string) (err error) {
	if state >= 0 && access == 0 {
		return
	}
	if state < 0 {
		if state == archive.StateForbidFixed {
			log.Warn("archive(%d) is fixed", aid)
		} else if state == archive.StateForbidUpDelete {
			log.Warn("archive(%d) is deleted", aid)
		} else {
			log.Warn("mid(%d) have not access view not pass archive(%d) ", mid, aid)
		}
		err = ecode.NothingFound
		return
	}
	if mid == 0 {
		log.Warn("not login can not view(%d) state(%d) access(%d) mid(%d)", aid, state, access, mid)
		err = ecode.AccessDenied
		s.prom.Incr("no_login_access")
		return
	}
	card, err := s.accDao.Card3(c, mid)
	if err != nil {
		log.Warn("s.accDao.Info failed can not view(%d) state(%d) access(%d)", aid, state, access)
		s.prom.Incr("err_login_access")
		return
	}
	if access > 0 && int(card.Rank) < access && !account.IsVip(card) {
		err = ecode.AccessDenied
		log.Warn("mid(%d) rank(%d) vip(tp:%d,status:%d) have not access(%d) view archive(%d) ", mid, card.Rank, card.Vip.Type, card.Vip.Status, access, aid)
		s.prom.Incr("login_access")
	}
	return
}
