package service

import (
	"context"
	"strings"
	"time"

	"go-common/app/job/main/workflow/model"
	"go-common/library/log"
)

const (
	_wfKeyPrefix      = "wf_"
	_feedbackDealType = 1
)

// queueProc .
func (s *Service) queueproc(c context.Context, dealType int) {
	for {
		var (
			key     string
			listMap = make(map[string][]int64)
		)
		sParams := searchParams(c, dealType, model.ListBefore, s.businessAttr)
		cLists, err := s.challByIDs(c, sParams)
		if err != nil {
			log.Error("s.challByIDs error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		if len(cLists) > 0 {
			for _, cList := range cLists {
				if cList.DispatchState != model.QueueState {
					continue
				}
				now := time.Now().Format("2006-01-02 15:04:05")
				log.Info("current cid(%d) dispatch_state is (%d) time is (%s)", cList.ID, cList.DispatchState, now)
				key = genKey(c, cList.Business, dealType)
				listMap[key] = append(listMap[key], cList.ID)
			}
			for key, list := range listMap {
				newDispatchState := s.dispatchState(c, dealType, model.ListBefore, cLists[list[0]].DispatchState)
				err := s.dao.UpDispatchStateByIDs(c, list, newDispatchState)
				if err != nil {
					log.Error("s.dao.UpDispatchStateByIDs error(%v)", err)
					time.Sleep(time.Second * 3)
					continue
				}
				now := time.Now().Format("2006-01-02 15:04:05")
				log.Info("this cids(%v) change dispatch_state to (%d) time is (%s)", list, newDispatchState, now)
				err = s.dao.SetList(c, key, list)
				if err != nil {
					log.Error("s.dao.SetList error(%v)", err)
					time.Sleep(time.Second * 3)
					continue
				}
			}
		}
		time.Sleep(time.Second * 15)
	}
}

// repairListProc .
func (s *Service) repairQueueproc(c context.Context, dealType int) {
	s.setCrash(c)

	for {
		time.Sleep(time.Second * 30)
		exist, err := s.dao.IsCrash(c)
		if err != nil {
			log.Error("s.dao.ExistKey error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		if exist {
			continue
		}
		var keySlice []string
		for _, attr := range s.businessAttr {
			var key string
			if attr.AssignType == model.SysAssignType {
				continue
			}
			if dealType == model.FDealType {
				if dealType == attr.DealType {
					key = genKey(c, attr.ID, dealType)
				}
			} else if dealType == model.ADealType {
				key = genKey(c, attr.ID, dealType)
			}
			keySlice = append(keySlice, key)
		}
		sParams := searchParams(c, dealType, model.ListIng, s.businessAttr)
		for _, key := range keySlice {
			var cids []int64
			sParams.Business = strings.Split(key, "_")[1]
			searchRes, err := s.dao.SearchChall(c, sParams)
			if err != nil {
				log.Error("s.dao.SearchChall error(%v)", err)
				time.Sleep(time.Second * 3)
				continue
			}
			searchDataRes := searchRes.Result
			if len(searchDataRes) > 0 {
				for _, r := range searchDataRes {
					cids = append(cids, r.ID)
				}
				err := s.dao.SetList(c, key, cids)
				if err != nil {
					log.Error("s.dao.SetList error(%v)", err)
					time.Sleep(time.Second * 3)
					continue
				}
			}
		}
		s.setCrash(c)
	}
}

// SetConstKey .
func (s *Service) setCrash(c context.Context) {
	for {
		if err := s.dao.SetCrash(c); err != nil {
			log.Error("s.dao.SetString error(%v)", err)
			time.Sleep(time.Second * 3)
			continue
		}
		break
	}
}
