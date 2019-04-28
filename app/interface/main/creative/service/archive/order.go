package archive

import (
	"context"
	"time"

	"go-common/app/interface/main/creative/model/archive"
	"go-common/app/interface/main/creative/model/game"
	"go-common/app/interface/main/creative/model/order"
	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// ExecuteOrders fn
func (s *Service) ExecuteOrders(c context.Context, mid int64, ip string) (orders []*order.Order, err error) {
	if orders, err = s.order.ExecuteOrders(c, mid, ip); err != nil {
		log.Error("s.order.ExecuteOrders mid(%d) err(%v)", mid, err)
		return
	}
	return
}

// Oasis for orders.
func (s *Service) Oasis(c context.Context, mid int64, ip string) (oa *order.Oasis, err error) {
	if oa, err = s.order.Oasis(c, mid, ip); err != nil {
		log.Error("s.order.Oasis mid(%d) err(%v)", mid, err)
		return
	}
	return
}

// ArcOrderGameInfo fn
// platform 1:android;2:ios
func (s *Service) ArcOrderGameInfo(c context.Context, aid int64, platform int, ip string) (gameInfo *game.Info, err error) {
	var (
		orderID, gameBaseID int64
		beginDate           xtime.Time
	)
	if orderID, _, gameBaseID, err = s.order.OrderByAid(c, aid); err != nil {
		log.Error("s.order.OrderByAid aid(%d)|ip(%s)|err(%+v)", aid, ip, err)
		err = ecode.NothingFound
		return
	}
	if gameBaseID == 0 || orderID == 0 {
		log.Error("s.order.OrderByAid aid(%d)|ip(%s) not found", aid, ip)
		err = ecode.NothingFound
		return
	}
	if gameInfo, err = s.game.Info(c, gameBaseID, platform, ip); err != nil {
		log.Error("s.game.Info aid(%d)|orderID(%d)|gameBaseID(%d)|ip(%s)|err(%+v)", aid, orderID, gameBaseID, ip, err)
		err = ecode.NothingFound
		return
	}
	if !gameInfo.IsOnline {
		log.Error("s.game.Info IsOnline is false aid(%d)|orderID(%d)|gameBaseID(%d)|ip(%s)|err(%+v)", aid, orderID, gameBaseID, ip, err)
		err = ecode.NothingFound
		return
	}
	if beginDate, err = s.order.LaunchTime(c, orderID, ip); err != nil {
		log.Error("s.order.LaunchTime aid(%d)|orderID(%d)|gameBaseID(%d)|ip(%s)|err(%+v)", aid, orderID, gameBaseID, ip, err)
		return
	}
	gameInfo.BeginDate = beginDate
	gameInfo.BaseID = gameBaseID
	return
}

// ArcCommercial fn
func (s *Service) ArcCommercial(c context.Context, aid int64, ip string) (cm *archive.Commercial, err error) {
	var (
		orderID, gameBaseID int64
		beginDate           xtime.Time
		a                   *api.Arc
		cache               = true
	)
	// try cache
	if cm, err = s.arc.ArcCMCache(c, aid); err != nil {
		err = nil
		cache = false
	} else if cm != nil {
		s.pCacheHit.Incr("cmarc_cache")
		cache = false
		return
	}
	s.pCacheMiss.Incr("cmarc_cache")
	// get archive
	a, err = s.arc.Archive(c, aid, ip)
	if err != nil {
		log.Error("arcCommercial aid(%d)|ip(%s)|err(%+v)", aid, ip, err)
		err = ecode.NothingFound
		return
	}
	if a == nil {
		log.Error("arcCommercial nil aid(%d)|ip(%s)|err(%+v)", aid, ip, err)
		err = ecode.NothingFound
		return
	}
	// add cache
	defer func() {
		if cache {
			s.addCache(func() {
				if cm == nil {
					cm = &archive.Commercial{}
				}
				s.arc.AddArcCMCache(context.Background(), aid, cm)
			})
		}
	}()
	// check order or porder
	if a.OrderID > 0 {
		// order
		if orderID, _, gameBaseID, err = s.order.OrderByAid(c, aid); err != nil {
			log.Error("arcCommercial aid(%d)|ip(%s)|error(%+v)", aid, ip, err)
			err = ecode.NothingFound
			return
		}
		if gameBaseID == 0 || orderID == 0 {
			log.Error("arcCommercial aid(%d)|ip(%s) not found", aid, ip)
			err = ecode.NothingFound
			return
		}
		if beginDate, err = s.order.LaunchTime(c, orderID, ip); err != nil {
			log.Error("arcCommercial get launch time failed. aid(%d)|orderID(%d)|gameBaseID(%d)|ip(%s)|error(%+v)", aid, orderID, gameBaseID, ip, err)
			return
		}
		// check time
		if time.Now().Unix() < beginDate.Time().Unix() {
			log.Error("arcCommercial launch time invalid. aid(%d)|orderID(%d)|gameBaseID(%d) beginDate(%+v)", aid, orderID, gameBaseID, beginDate)
			err = ecode.NothingFound
			return
		}
		cm = &archive.Commercial{}
		cm.AID = a.Aid
		cm.OrderID = orderID
		cm.GameID = gameBaseID
	} else if a.AttrVal(arcMdl.AttrBitIsPorder) == arcMdl.AttrYes {
		// porder
		var pd *archive.Porder
		if pd, err = s.arc.Porder(c, aid); err != nil {
			log.Error("arcCommercial aid(%d) error(%v)", aid, err)
			err = ecode.NothingFound
			return
		}
		if pd == nil {
			log.Error("arcCommercial porder not show aid(%d)", aid)
			err = ecode.CreativePorderForbidShowFront
			return
		}
		cm = &archive.Commercial{}
		cm.AID = a.Aid
		cm.POrderID = pd.ID
		cm.GameID = pd.BrandID
	} else {
		log.Error("arcCommercial is not commercial. aid(%d)|ip(%s)|", aid, ip)
		err = ecode.NothingFound
		return
	}
	return
}

// UpValidate func
func (s *Service) UpValidate(c context.Context, mid int64, ip string) (uv *order.UpValidate, err error) {
	if uv, err = s.order.UpValidate(c, mid, ip); err != nil {
		log.Error("s.order.UpValidate mid(%d)|ip(%s)|err(%v)", mid, ip, err)
	}
	return
}
