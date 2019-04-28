package service

import (
	"context"
	"encoding/json"
	"math"
	"strconv"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

// add to blacklist reason
const (
	// _all          = 0
	// _stopIncome   = 1
	// _breachRecord = 2
	_porder       = 3
	_executeOrder = 4
)

// InitBlacklistMID init av_black_list mid
func (s *Service) InitBlacklistMID(c context.Context) (err error) {
	blacklist, err := s.listBlacklist(c, "mid = 0")
	if err != nil {
		log.Error("s.listBlacklist error(%v)", err)
		return
	}

	avIDs := make([]int64, 0)
	for _, b := range blacklist {
		if b.MID == 0 {
			avIDs = append(avIDs, b.AvID)
		}
	}
	if len(avIDs) == 0 {
		return
	}
	m, err := s.GetAvsMID(c, avIDs)
	if err != nil {
		log.Error("GetAvsMID error(%v)", err)
		return
	}

	for i := 0; i < len(blacklist); i++ {
		blacklist[i].MID = m[blacklist[i].AvID]
	}
	_, err = s.updateBlacklistBatch(c, blacklist)
	return
}

func (s *Service) listBlacklist(c context.Context, query string) (list []*model.Blacklist, err error) {
	from, limit := 0, 2000
	var b []*model.Blacklist
	for {
		b, err = s.dao.ListBlacklist(c, query, from, limit)
		if err != nil {
			return
		}
		list = append(list, b...)
		if len(b) < limit {
			break
		}
		from += len(b)
	}
	return
}

// UpdateBlacklist update blacklist
func (s *Service) UpdateBlacklist(c context.Context) (err error) {
	defer func() {
		GetTaskService().SetTaskStatus(c, TaskBlacklist, time.Now().AddDate(0, 0, -1).Format(_layout), err)
	}()

	blacklist := make([]*model.Blacklist, 0)
	porders, err := s.getNewPorder(c)
	if err != nil {
		log.Error("s.getNewPorder error(%v)", err)
		return
	}
	log.Info("Get new porder %d", len(porders))
	blacklist = append(blacklist, porders...)

	executeOrders, err := s.getNewExecuteOrder(c)
	if err != nil {
		log.Error("s.getNewExecuteOrder error(%v)", err)
		return
	}
	log.Info("Get new execute order %d", len(executeOrders))

	blacklist = append(blacklist, executeOrders...)
	count, err := s.updateBlacklistBatch(c, blacklist)
	if err != nil {
		log.Error("s.updateBlacklistBatch error(%v)", err)
		return
	}
	log.Info("Add %d list into blacklist", count)
	return
}

func (s *Service) updateBlacklistBatch(c context.Context, blacklist []*model.Blacklist) (count int64, err error) {
	ups, err := s.getHasSignUpInfo(c)
	if err != nil {
		log.Error("s.dao.GetHasSignUpInfo error(%v)", err)
		return
	}

	for i := 0; i < len(blacklist); i++ {
		if nickname, ok := ups[blacklist[i].MID]; ok {
			blacklist[i].HasSigned = 1
			blacklist[i].Nickname = nickname
		}
	}
	return s.dao.AddBlacklistBatch(c, blacklist)
}

func (s *Service) getHasSignUpInfo(c context.Context) (m map[int64]string, err error) {
	m = make(map[int64]string)
	offset, limit := 0, 2000
	for {
		err = s.dao.GetHasSignUpInfo(c, offset, limit, m)
		if err != nil {
			log.Error("s.dao.GetHasSignUpInfo error(%v)", err)
			return
		}
		offset += limit
		if len(m) < offset {
			break
		}
	}
	return
}

func (s *Service) getNewPorder(c context.Context) (blacklist []*model.Blacklist, err error) {
	beginTime, err := s.dao.GetLastCtime(c, _porder)
	if err != nil {
		log.Error("s.dao.GetLastCtime error(%v)", err)
		return
	}
	if beginTime != 0 {
		beginTime -= 10 * 60 // pre 10min
	}
	endTime := time.Now().Unix()

	porders, err := s.getPorder(beginTime, endTime)
	if err != nil {
		log.Error("get Porder error(%v)", err)
		return
	}

	// get porder mid
	avIds := []int64{}
	for _, b := range porders {
		avIds = append(avIds, b.AID)
	}
	m, err := s.GetAvsMID(c, avIds)
	if err != nil {
		log.Error("s.dao.GetAvsMID error(%v)", err)
		return
	}

	blacklist = make([]*model.Blacklist, len(porders))
	for i := 0; i < len(porders); i++ {
		blacklist[i] = &model.Blacklist{
			AvID:   porders[i].AID,
			MID:    m[porders[i].AID],
			Reason: _porder,
		}
	}

	return
}

// GetAvsMID get avs mid from api
func (s *Service) GetAvsMID(c context.Context, avs []int64) (avsMap map[int64]int64, err error) {
	avsMap = make(map[int64]int64)
	if len(avs) == 0 {
		return
	}
	start, limit := 0, 10
	if limit > len(avs) {
		limit = len(avs)
	}
	for start+limit <= len(avs) {
		if err = s.getAvsMID(c, avs[start:start+limit], avsMap); err != nil {
			return
		}
		start += limit
		if start < len(avs) && start+limit > len(avs) {
			limit = len(avs) - start
		}
	}
	log.Info("Get avs(%d) from archiveURL", len(avsMap))
	return
}

func (s *Service) getAvsMID(c context.Context, avs []int64, avsMap map[int64]int64) (err error) {
	params := map[string]string{"aids": xstr.JoinInts(avs)}
	body, err := s.HTTPClient("GET", s.conf.Host.Archives, params, time.Now().Unix())
	if err != nil {
		log.Error("s.HTTPClient error(%v)", err)
		return
	}

	res := model.ArchiveRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("json.Unmarshal body %s error(%v)", string(body), err)
		return
	}

	for _, archive := range res.Data {
		avsMap[archive.AID] = archive.Owner.MID
	}
	return
}

func (s *Service) getPorder(begin, end int64) (porders []*model.Porder, err error) {
	params := map[string]string{
		"begin": strconv.FormatInt(begin, 10),
		"end":   strconv.FormatInt(end, 10),
	}
	body, err := s.HTTPClient("GET", s.conf.Host.Porder, params, time.Now().UnixNano()/int64(math.Pow(10, 6)))
	if err != nil {
		log.Error("s.HTTPClient error(%v)", err)
		return
	}

	res := model.PorderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		log.Error("json.Unmarshal body %s error(%v)", string(body), err)
		return
	}

	porders = res.Data
	return
}

func (s *Service) getNewExecuteOrder(c context.Context) (blacklist []*model.Blacklist, err error) {
	beginTime, err := s.dao.GetLastCtime(c, _executeOrder)
	if err != nil {
		log.Error("s.dao.GetLastCtime error(%v)", err)
		return
	}

	if beginTime != 0 {
		beginTime -= 10 * 60 // pre 10min
	}

	executeOrders, err := s.dao.GetExecuteOrder(c, time.Unix(beginTime, 0), time.Now())
	if err != nil {
		log.Error("s.dao.GetExecuteOrder error(%v)", err)
		return
	}

	blacklist = make([]*model.Blacklist, len(executeOrders))
	for i := 0; i < len(executeOrders); i++ {
		blacklist[i] = &model.Blacklist{
			AvID:   executeOrders[i].AvID,
			MID:    executeOrders[i].MID,
			Reason: _executeOrder,
		}
	}
	return
}
