package service

import (
	"context"
	"strconv"
	"time"

	"go-common/app/job/main/workflow/model"
	"go-common/library/log"
)

// taskExpireProc task expire.
func (s *Service) taskExpireproc(c context.Context, dealType int) {
	var businessExpireMap = make(map[int64]int)
	for _, v := range s.businessAttr {
		if v.AssignType == 1 {
			continue
		}
		businessExpireMap[v.ID] = v.ExpireTime
	}
	sParams := searchParams(c, dealType, model.ListAfter, s.businessAttr)
	for {
		var expireCids []int64
		cLists, err := s.challByIDs(c, sParams)
		if err != nil {
			log.Error("s.challByIDs error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		if len(cLists) <= 0 {
			time.Sleep(time.Second * 30)
			continue
		}
		now := time.Now()
		for k, cl := range cLists {
			if cl.DispatchState != model.QueueState {
				continue
			}
			dispatchTime := cl.DispatchTime.Format("2006-01-02 15:04:05")
			expireTime := businessExpireMap[cl.Business]
			m, _ := time.ParseDuration("-" + strconv.Itoa(expireTime) + "m")
			if now.Add(m).Format("2006-01-02 15:04:05") > dispatchTime {
				expireCids = append(expireCids, k)
			}
		}
		if len(expireCids) > 0 {
			log.Info("expire cids is %v", expireCids)
			assignAdminid := int64(0)
			newDispatchState := s.dispatchState(c, dealType, model.ListAfter, cLists[expireCids[0]].DispatchState)
			err := s.dao.UpDispatchStateAdminIDByIds(c, expireCids, newDispatchState, assignAdminid)
			if err != nil {
				log.Error("s.dao.UpDispatchStateAdminIDByIds error(%v)", err)
				time.Sleep(time.Second * 3)
				continue
			}
		}
		time.Sleep(time.Second * 30)
	}
}
