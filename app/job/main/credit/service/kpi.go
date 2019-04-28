package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/job/main/credit/model"
	"go-common/library/log"
)

// KPIReward send reward to user by kpi info.
func (s *Service) KPIReward(c context.Context, nwMsg []byte, oldMsg []byte) (err error) {
	var (
		mr          = &model.Kpi{}
		res         model.Kpi
		nameplentID int64
		coins       float64
	)
	if err = json.Unmarshal(nwMsg, mr); err != nil {
		log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
		return
	}
	if res, err = s.dao.KPIInfo(c, mr.ID); err != nil {
		log.Error("json.KPIInf(%d) error(%v)", mr.ID, err)
		return
	}
	if res.HandlerStatus == 1 {
		return
	}
	s.dao.SendMsg(c, mr.Mid, _msgTitle, _msgContext)
	coins = model.KpiCoinsRate(mr.Rate)
	if coins > 0 {
		s.dao.AddMoney(c, mr.Mid, coins, model.KPICoinsReason)
	}
	s.dao.UpdateKPIHandlerStatus(c, mr.ID)
	if pend, ok := model.LevelPendantByKPI(int8(mr.Rate)); ok {
		expired := time.Now().AddDate(0, 0, 30)
		if err = s.dao.UpdateJuryExpired(c, mr.Mid, expired); err != nil {
			return
		}
		if len(pend) > 0 {
			if err = s.dao.SendPendant(c, mr.Mid, pend, model.KPIDefealtPendSendDays); err != nil {
				log.Error("s.dao.SendPendant err(%v)", err)
			}
			if err == nil {
				s.dao.UpdateKPIPendentStatus(c, mr.ID)
			}
			if mr.Mid%50 == 1 {
				time.Sleep(time.Second)
			}
		}
	}
	num, err := s.dao.CountKPIRate(c, mr.Mid)
	if err != nil {
		log.Error("s.dao.CountKPIRate(mid:%d) err(%v)", mr.Mid, err)
		return
	}
	nameplentID = model.KpiPlateIDRateTimes(num)
	if nameplentID != 0 {
		for i := 0; i <= 5; i++ {
			if err = s.dao.SendMedal(c, mr.Mid, nameplentID); err != nil {
				log.Error("s.dao.SendMedal(mid:%d nameNameplatid:%d) err(%v)", mr.Mid, nameplentID, err)
				continue
			}
			break
		}
	}
	return
}
