package service

import (
	"context"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/workflow/model"
	srvmodel "go-common/app/service/main/workflow/model"
	"go-common/library/log"
)

// searchParams .
func searchParams(c context.Context, dealType, listState int, busAttr []*model.BusinessAttr) (params *model.SearchParams) {
	var businessArr []string
	params = &model.SearchParams{}
	if listState == model.ListBefore {
		params.AssigneeAdminIDs = "0"
		params.AssigneeAdminIDsNot = ""
		switch dealType {
		case model.FDealType:
			params.States = model.FListBeforeStates
			params.BusinessStates = model.FListBeforeBusinessStates
			params.MtimeTo = time.Now().Add(-time.Minute * 1).Format("2006-01-02 15:04:05")
		case model.ADealType:
			params.States = model.AListBeforeStates
		}
	} else if listState == model.ListAfter {
		params.AssigneeAdminIDs = ""
		params.AssigneeAdminIDsNot = "0"
		switch dealType {
		case model.FDealType:
			params.States = model.FListAfterStates
			params.BusinessStates = model.FListAfterBusinessStates
		case model.ADealType:
			params.States = model.AListAfterStates
		}
	} else if listState == model.ListIng {
		params.AssigneeAdminIDs = ""
		params.AssigneeAdminIDsNot = ""
		switch dealType {
		case model.FDealType:
			params.States = model.FListAfterStates
			params.BusinessStates = model.FListAfterBusinessStates
		case model.ADealType:
			params.States = model.AListAfterStates
		}
	}
	for _, attr := range busAttr {
		if attr.AssignType == model.SysAssignType {
			continue
		}
		if dealType == model.ADealType {
			businessArr = append(businessArr, strconv.FormatInt(attr.ID, 10))
		} else {
			if attr.DealType == dealType {
				businessArr = append(businessArr, strconv.FormatInt(attr.ID, 10))
			}
		}
	}
	params.Business = strings.Join(businessArr, ",")
	return
}

// challByIDs .
func (s *Service) challByIDs(c context.Context, params *model.SearchParams) (res map[int64]*model.Chall, err error) {
	var cids []int64
	searchRes, err := s.dao.SearchChall(c, params)
	if err != nil {
		log.Error("s.dao.SearchChall error(%v)", err)
		return
	}
	searchDataRes := searchRes.Result
	if len(searchDataRes) > 0 {
		for _, r := range searchDataRes {
			cids = append(cids, r.ID)
		}
		res, err = s.dao.ChallByIDs(c, cids)
	}
	return
}

// disPatchState .
func (s *Service) dispatchState(c context.Context, dealType, listState, oldDispatchState int) (newDispatchState int64) {
	state := oldDispatchState & srvmodel.QueueState
	if dealType == model.FDealType {
		if listState == model.ListBefore {
			newDispatchState, _ = strconv.ParseInt("f"+strconv.Itoa(state), 16, 64)
		} else if listState == model.ListAfter {
			newDispatchState, _ = strconv.ParseInt("1"+strconv.Itoa(state), 16, 64)
		}
	} else if dealType == model.ADealType {
		if listState == model.ListBefore {
			newDispatchState = int64(srvmodel.QueueState)
		} else if listState == model.ListAfter {
			newDispatchState = int64(srvmodel.QueueStateBefore)
		}
	}
	return
}

// key .
func genKey(c context.Context, business int64, dealType int) (key string) {
	key = _wfKeyPrefix + strconv.FormatInt(business, 10) + "_" + strconv.Itoa(dealType)
	return
}
