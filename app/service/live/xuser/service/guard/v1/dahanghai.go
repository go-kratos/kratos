package v1

import (
	"context"
	"github.com/pkg/errors"
	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	"go-common/app/service/live/xuser/dao/exp"
	"go-common/library/ecode"
	"strconv"
	"time"

	confm "go-common/app/service/live/xuser/conf"
	dhhm "go-common/app/service/live/xuser/model/dhh"
	"go-common/library/log"

	dahanghaiModel "go-common/app/service/live/xuser/model/dhh"
)

const (
	_errorServiceLogPrefix = "xuser.dhh.service"
	_promCacheMissed       = "xuser_dhh_redis:用户守护cache miss"
	_promAnchorCacheMissed = "xuser_dhh_redis:主播侧守护cache miss"
	_promCacheHitAll       = "xuser_dhh_redis:用户守护cache全部命中"
	_promAnchorCacheHitAll = "xuser_dhh_redis:主播侧守护cache全部命中"
	// _retryAddExpTimes      = 3
	_formatTime = "2006-01-02 15:04:05"
	_default    = 0
	_zongdu     = 1
	_tidu       = 2
	_jianzhang  = 3
)

// RParseInt 转换
func RParseInt(inputStr string, defaultValue int64) (output int64) {
	if mid, err := strconv.ParseInt(inputStr, 10, 64); err == nil {
		output = mid
	} else {
		output = defaultValue
	}
	return
}

func (s *GuardService) adaptResultFromRedis(uid int64, targetID int64, sortType int64, dhhList []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfo) {
	resultP1 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP2 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP3 := make([]*v1pb.DaHangHaiInfo, 0)
	for _, v := range dhhList {
		TargetID := RParseInt(v.TargetId, 0)
		PrivilegeType := RParseInt(v.PrivilegeType, 0)
		if TargetID == targetID {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = RParseInt(v.Id, 0)
			respItem.Uid = RParseInt(v.Uid, 0)
			respItem.TargetId = RParseInt(v.TargetId, 0)
			respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			switch PrivilegeType {
			case _zongdu:
				{
					resultP1 = append(resultP1, respItem)
				}
			case _tidu:
				{
					resultP2 = append(resultP2, respItem)
				}
			case _jianzhang:
				{
					resultP3 = append(resultP3, respItem)
				}
			}
		}
	}
	result = s.mergeResult(resultP1, resultP2, resultP3, sortType)
	return
}

func (s *GuardService) filterUIDTopFromRedis(uid int64, dhhList []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	for _, v := range dhhList {
		PrivilegeType := RParseInt(v.PrivilegeType, 0)
		if len(result) <= 0 {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = RParseInt(v.Id, 0)
			respItem.Uid = RParseInt(v.Uid, 0)
			respItem.TargetId = RParseInt(v.TargetId, 0)
			respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			result[uid] = respItem
		} else {
			if result[uid].PrivilegeType > PrivilegeType {
				respItem := &v1pb.DaHangHaiInfo{}
				respItem.Id = RParseInt(v.Id, 0)
				respItem.Uid = RParseInt(v.Uid, 0)
				respItem.TargetId = RParseInt(v.TargetId, 0)
				respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
				respItem.StartTime = v.StartTime
				respItem.ExpiredTime = v.ExpiredTime
				respItem.Ctime = v.Ctime
				respItem.Utime = v.Utime
				result[uid] = respItem
			}
		}
	}
	return
}

func (s *GuardService) getAnchorTopGuardCount(uid int64, dhhList []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	for _, v := range dhhList {
		PrivilegeType := RParseInt(v.PrivilegeType, 0)
		UID := RParseInt(v.Uid, 0)
		if _, exist := result[UID]; !exist {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = RParseInt(v.Id, 0)
			respItem.Uid = RParseInt(v.Uid, 0)
			respItem.TargetId = RParseInt(v.TargetId, 0)
			respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			result[UID] = respItem
		} else {
			if result[UID].PrivilegeType > PrivilegeType {
				respItem := &v1pb.DaHangHaiInfo{}
				respItem.Id = RParseInt(v.Id, 0)
				respItem.Uid = RParseInt(v.Uid, 0)
				respItem.TargetId = RParseInt(v.TargetId, 0)
				respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
				respItem.StartTime = v.StartTime
				respItem.ExpiredTime = v.ExpiredTime
				respItem.Ctime = v.Ctime
				respItem.Utime = v.Utime
				result[UID] = respItem
			}
		}
	}
	return
}

func (s *GuardService) getAnchorTopGuardCountFromDB(uid int64, dhhList []*dhhm.DHHDBTime) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	for _, v := range dhhList {
		PrivilegeType := v.PrivilegeType
		UID := v.UID
		if _, exist := result[UID]; !exist {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = v.ID
			respItem.Uid = v.UID
			respItem.TargetId = v.TargetId
			respItem.PrivilegeType = v.PrivilegeType
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			result[UID] = respItem
		} else {
			if result[UID].PrivilegeType > PrivilegeType {
				respItem := &v1pb.DaHangHaiInfo{}
				respItem.Id = v.ID
				respItem.Uid = v.UID
				respItem.TargetId = v.TargetId
				respItem.PrivilegeType = v.PrivilegeType
				respItem.StartTime = v.StartTime
				respItem.ExpiredTime = v.ExpiredTime
				respItem.Ctime = v.Ctime
				respItem.Utime = v.Utime
				result[UID] = respItem
			}
		}
	}
	return
}

func (s *GuardService) filterResultFromRedis(uid int64, targetMap map[int64]int64, dhhList []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfo) {
	resultP1 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP2 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP3 := make([]*v1pb.DaHangHaiInfo, 0)
	for _, v := range dhhList {
		PrivilegeType := RParseInt(v.PrivilegeType, 0)
		respItem := &v1pb.DaHangHaiInfo{}
		respItem.Id = RParseInt(v.Id, 0)
		respItem.Uid = RParseInt(v.Uid, 0)
		respItem.TargetId = RParseInt(v.TargetId, 0)
		respItem.PrivilegeType = RParseInt(v.PrivilegeType, 0)
		respItem.StartTime = v.StartTime
		respItem.ExpiredTime = v.ExpiredTime
		respItem.Ctime = v.Ctime
		respItem.Utime = v.Utime
		switch PrivilegeType {
		case _zongdu:
			{
				resultP1 = append(resultP1, respItem)
			}
		case _tidu:
			{
				resultP2 = append(resultP2, respItem)
			}
		case _jianzhang:
			{
				resultP3 = append(resultP3, respItem)
			}
		}
	}

	result = s.mergeResultWithTargetMap(resultP1, resultP2, resultP3, targetMap)
	return
}

func (s *GuardService) filterResultFromDB(uid int64, targetMap map[int64]int64, dhhList []*dahanghaiModel.DHHDBTime) (result map[int64]*v1pb.DaHangHaiInfo) {
	resultP1 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP2 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP3 := make([]*v1pb.DaHangHaiInfo, 0)
	for _, v := range dhhList {
		PrivilegeType := v.PrivilegeType
		respItem := &v1pb.DaHangHaiInfo{}
		respItem.Id = v.ID
		respItem.Uid = v.UID
		respItem.TargetId = v.TargetId
		respItem.PrivilegeType = v.PrivilegeType
		respItem.StartTime = v.StartTime
		respItem.ExpiredTime = v.ExpiredTime
		respItem.Ctime = v.Ctime
		respItem.Utime = v.Utime
		switch PrivilegeType {
		case _zongdu:
			{
				resultP1 = append(resultP1, respItem)
			}
		case _tidu:
			{
				resultP2 = append(resultP2, respItem)
			}
		case _jianzhang:
			{
				resultP3 = append(resultP3, respItem)
			}
		}
	}

	result = s.mergeResultWithTargetMap(resultP1, resultP2, resultP3, targetMap)
	return
}

func (s *GuardService) mergeResultForGetByUIDBatch(dbResult map[int64][]*dahanghaiModel.DaHangHaiRedis2, cacheResult []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfoList) {
	result = make(map[int64]*v1pb.DaHangHaiInfoList)
	for _, v := range cacheResult {
		uid := RParseInt(v.Uid, 0)
		if _, exist := result[uid]; !exist {
			result[uid] = &v1pb.DaHangHaiInfoList{}
			result[uid].List = make([]*v1pb.DaHangHaiInfo, 0)
		}
		item := &v1pb.DaHangHaiInfo{}
		item.Id = RParseInt(v.Id, 0)
		item.Uid = RParseInt(v.Uid, 0)
		item.TargetId = RParseInt(v.TargetId, 0)
		item.PrivilegeType = RParseInt(v.PrivilegeType, 0)
		item.Ctime = v.Ctime
		item.Utime = v.Utime
		item.ExpiredTime = v.ExpiredTime
		item.StartTime = v.StartTime
		result[uid].List = append(result[uid].List, item)
	}

	for _, v := range dbResult {
		for _, vv := range v {
			uid := RParseInt(vv.Uid, 0)
			if _, exist := result[uid]; !exist {
				result[uid] = &v1pb.DaHangHaiInfoList{}
				result[uid].List = make([]*v1pb.DaHangHaiInfo, 0)
			}
			item := &v1pb.DaHangHaiInfo{}
			item.Id = RParseInt(vv.Id, 0)
			item.Uid = RParseInt(vv.Uid, 0)
			item.TargetId = RParseInt(vv.TargetId, 0)
			item.PrivilegeType = RParseInt(vv.PrivilegeType, 0)
			item.Ctime = vv.Ctime
			item.Utime = vv.Utime
			item.ExpiredTime = vv.ExpiredTime
			item.StartTime = vv.StartTime
			result[uid].List = append(result[uid].List, item)
		}
	}
	return
}

func (s *GuardService) formatDaHangHaiCache(uid int64, dbData []*dahanghaiModel.DHHDBTime) (dhhList []*dahanghaiModel.DaHangHaiRedis2) {
	dhhList = make([]*dahanghaiModel.DaHangHaiRedis2, 0)
	if len(dbData) <= 0 {
		return
	}
	for _, v := range dbData {
		item := &dahanghaiModel.DaHangHaiRedis2{}
		item.Id = strconv.Itoa(int(v.ID))
		item.Uid = strconv.Itoa(int(v.UID))
		item.TargetId = strconv.Itoa(int(v.TargetId))
		item.PrivilegeType = strconv.Itoa(int(v.PrivilegeType))
		item.StartTime = v.StartTime
		item.ExpiredTime = v.ExpiredTime
		item.Ctime = v.Ctime
		item.Utime = v.Utime
		dhhList = append(dhhList, item)
	}
	return
}

func (s *GuardService) formatDaHangHaiCacheBatch(dbData map[int64][]*dahanghaiModel.DHHDBTime) (dhhMapList map[int64][]*dahanghaiModel.DaHangHaiRedis2) {
	dhhMapList = make(map[int64][]*dahanghaiModel.DaHangHaiRedis2)
	if len(dbData) <= 0 {
		return
	}
	for k, v := range dbData {
		dhhList := make([]*dahanghaiModel.DaHangHaiRedis2, 0)
		for _, vv := range v {
			item := &dahanghaiModel.DaHangHaiRedis2{}
			item.Id = strconv.Itoa(int(vv.ID))
			item.Uid = strconv.Itoa(int(vv.UID))
			item.TargetId = strconv.Itoa(int(vv.TargetId))
			item.PrivilegeType = strconv.Itoa(int(vv.PrivilegeType))
			item.StartTime = vv.StartTime
			item.ExpiredTime = vv.ExpiredTime
			item.Ctime = vv.Ctime
			item.Utime = vv.Utime
			dhhList = append(dhhList, item)
		}
		dhhMapList[k] = dhhList
	}
	return
}

func (s *GuardService) adaptResultFromDB(uid int64, targetID int64, sortType int64, dhhList []*dahanghaiModel.DHHDBTime) (result map[int64]*v1pb.DaHangHaiInfo) {
	resultP1 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP2 := make([]*v1pb.DaHangHaiInfo, 0)
	resultP3 := make([]*v1pb.DaHangHaiInfo, 0)
	for _, v := range dhhList {
		if v.TargetId == targetID {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = v.ID
			respItem.Uid = v.UID
			respItem.TargetId = v.TargetId
			respItem.PrivilegeType = v.PrivilegeType
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			switch v.PrivilegeType {
			case _zongdu:
				{
					resultP1 = append(resultP1, respItem)
				}
			case _tidu:
				{
					resultP2 = append(resultP2, respItem)
				}
			case _jianzhang:
				{
					resultP3 = append(resultP3, respItem)
				}
			}
		}
	}

	result = s.mergeResult(resultP1, resultP2, resultP3, sortType)

	return
}

func (s *GuardService) filterTopFromDB(uid int64, dhhList []*dahanghaiModel.DHHDBTime) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	for _, v := range dhhList {
		PrivilegeType := v.PrivilegeType
		if len(result) <= 0 {
			respItem := &v1pb.DaHangHaiInfo{}
			respItem.Id = v.ID
			respItem.Uid = v.UID
			respItem.TargetId = v.TargetId
			respItem.PrivilegeType = v.PrivilegeType
			respItem.StartTime = v.StartTime
			respItem.ExpiredTime = v.ExpiredTime
			respItem.Ctime = v.Ctime
			respItem.Utime = v.Utime
			result[uid] = respItem
		} else {
			if result[uid].PrivilegeType > PrivilegeType {
				respItem := &v1pb.DaHangHaiInfo{}
				respItem.Id = v.ID
				respItem.Uid = v.UID
				respItem.TargetId = v.TargetId
				respItem.PrivilegeType = v.PrivilegeType
				respItem.StartTime = v.StartTime
				respItem.ExpiredTime = v.ExpiredTime
				respItem.Ctime = v.Ctime
				respItem.Utime = v.Utime
				result[uid] = respItem
			}
		}
	}
	return
}

func (s *GuardService) mergeResult(resultZongdu []*v1pb.DaHangHaiInfo, resultTidu []*v1pb.DaHangHaiInfo, resultJianzhang []*v1pb.DaHangHaiInfo,
	sortType int64) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	switch sortType {
	case _default:
		{
			if len(resultZongdu) > 0 {
				for _, v := range resultZongdu {
					result[v.TargetId] = v
					return
				}
			}

			if len(resultTidu) > 0 {
				for _, v := range resultTidu {
					result[v.TargetId] = v
					return
				}

			}

			if len(resultJianzhang) > 0 {
				for _, v := range resultJianzhang {
					result[v.TargetId] = v
					return
				}
			}
		}
	case _zongdu:
		{
			for _, v := range resultZongdu {
				result[v.TargetId] = v
				return
			}
		}
	case _tidu:
		{
			for _, v := range resultTidu {
				result[v.TargetId] = v
				return
			}
		}
	case _jianzhang:
		{
			for _, v := range resultJianzhang {
				result[v.TargetId] = v
				return
			}
		}
	}
	return
}

func (s *GuardService) mergeResultWithTargetMap(resultZongdu []*v1pb.DaHangHaiInfo, resultTidu []*v1pb.DaHangHaiInfo, resultJianzhang []*v1pb.DaHangHaiInfo,
	targetMap map[int64]int64) (result map[int64]*v1pb.DaHangHaiInfo) {
	result = make(map[int64]*v1pb.DaHangHaiInfo)
	resultArr := make([]*v1pb.DaHangHaiInfo, 0)
	if len(resultZongdu) > 0 {
		resultArr = append(resultArr, resultZongdu...)
	}

	if len(resultTidu) > 0 {
		resultArr = append(resultArr, resultTidu...)
	}

	if len(resultJianzhang) > 0 {
		resultArr = append(resultArr, resultJianzhang...)
	}
	for kmap, vmap := range targetMap {
		switch vmap {
		case _default:
			{
				for _, varr := range resultArr {
					if varr.TargetId == kmap {
						result[varr.TargetId] = varr
						break
					}
				}
			}
		case _zongdu, _tidu, _jianzhang:
			{
				for _, varr := range resultArr {
					if (varr.TargetId == kmap) && (varr.PrivilegeType == vmap) {
						result[varr.TargetId] = varr
						break
					}
				}
			}
		}
	}
	return
}

func (s *GuardService) findMissUIDs(hits map[int64]bool, req []int64) (resp []int64) {
	resp = make([]int64, 0)
	if len(req) <= 0 {
		return
	}
	for _, v := range req {
		if _, exist := hits[v]; !exist {
			resp = append(resp, v)
		}
	}
	return
}

// GetByUIDTargetID ...
// 单uid单targetid
func (s *GuardService) GetByUIDTargetID(ctx context.Context, req *v1pb.GetByUidTargetIdReq) (resp *v1pb.GetByUidTargetIdResp, err error) {
	reqStartTime := confm.RecordTimeCost()

	resp = &v1pb.GetByUidTargetIdResp{}
	resp.Data = make(map[int64]*v1pb.DaHangHaiInfo)
	cacheHealth := true
	uids := make([]int64, 0)
	uids = append(uids, req.Uid)
	dhhResultFromRedis, err := s.dao.GetUIDAllGuardFromRedis(ctx, uids)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDTargetID|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		cacheHealth = false
		return
	} else if dhhResultFromRedis != nil {
		resp.Data = s.adaptResultFromRedis(req.Uid, req.TargetId, req.SortType, dhhResultFromRedis)
		exp.PromCacheHit(_promCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promCacheMissed)

	resultDB, err := s.dao.GetByUIDs(ctx, uids)
	resultDBTime, _ := s.changeDBTime(ctx, resultDB)
	if err != nil {
		reqAfterQueryDBTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUidTargetId|s.dao.Exp|从DB获取大航海信息error|(%v),missedUIDs(%v)|耗时:%dms", err, uids, reqAfterQueryDBTime-reqStartTime)
		return
	}
	dhhResultFromDB := s.formatDaHangHaiCache(req.Uid, resultDBTime)
	// 写入缓存
	if cacheHealth {
		s.asyncSetDHHCache(ctx, dhhResultFromDB, req.Uid)
	}
	resp.Data = s.adaptResultFromDB(req.Uid, req.TargetId, req.SortType, resultDBTime)
	return
}

// ClearUIDCache 清除cache
func (s *GuardService) ClearUIDCache(ctx context.Context, req *v1pb.ClearUIDCacheReq) (resp *v1pb.ClearUIDCacheResp, err error) {
	resp = &v1pb.ClearUIDCacheResp{}
	if req.MagicKey != "jibamao" {
		err = errors.WithMessage(ecode.ResourceParamErr, "禁止通行！")
	}
	s.asyncCLearExpCache(ctx, req.Uid)
	return
}

// GetAnchorRecentTopGuard 获取最近开通总督的用户信息
func (s *GuardService) GetAnchorRecentTopGuard(ctx context.Context, req *v1pb.GetAnchorRecentTopGuardReq) (resp *v1pb.GetAnchorRecentTopGuardResp, err error) {
	resp = &v1pb.GetAnchorRecentTopGuardResp{}
	if req == nil || req.Uid <= 0 {
		return
	}
	guardResultFromRedis, err := s.dao.GetAnchorRecentTopGuardCache(ctx, req.Uid)
	if err != nil {
		err = errors.WithMessage(ecode.XUserGuardFetchRecentTopListFail, "get GetAnchorRecentTopGuard cache failed")
		return
	}
	uids := make(map[int64]int64)
	nowTime := time.Now().Unix()
	for k, v := range guardResultFromRedis {
		if v > nowTime {
			uids[k] = v
		}
	}
	if err != nil {
		return
	}
	retList := make([]*v1pb.GetAnchorRecentTopGuardList, 0)
	for k, v := range uids {
		item := &v1pb.GetAnchorRecentTopGuardList{}
		item.Uid = k
		item.IsOpen = 1
		item.EndTime = v
		retList = append(retList, item)
	}

	resp.List = retList
	resp.Cnt = int64(len(retList))
	return
}

// func (s *GuardService) getTopListGuardPrepareParams(req *v1pb.GetTopListGuardReq) (uid int64, page int64, pageSize int64) {
// 	uid = req.Uid
// 	if req.Page == 0 {
// 		page = 1
// 	} else {
// 		page = req.Page
// 	}
// 	if req.PageSize == 0 {
// 		pageSize = 10
// 	} else {
// 		pageSize = req.PageSize
// 	}
// 	return
// }

// GetTopListGuard ...
// 房间页守护排行榜
func (s *GuardService) GetTopListGuard(ctx context.Context, req *v1pb.GetTopListGuardReq) (resp *v1pb.GetTopListGuardResp, err error) {
	resp = &v1pb.GetTopListGuardResp{}
	// if req == nil {
	// 	return
	// }
	// uid, page, pagesize := s.getTopListGuardPrepareParams(req)
	return
}

// GetTopListGuardNum ...
// 房间页守护数量
func (s *GuardService) GetTopListGuardNum(ctx context.Context, req *v1pb.GetTopListGuardNumReq) (resp *v1pb.GetTopListGuardNumResp, err error) {
	reqStartTime := confm.RecordTimeCost()
	resp = &v1pb.GetTopListGuardNumResp{}
	cacheHealth := true
	uids := make([]int64, 0)
	uids = append(uids, req.Uid)
	dhhResultFromRedis, err := s.dao.GetAnchorAllGuardFromRedis(ctx, uids)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|"+_XUserGetAnchorGuardNumCacheError+"|GetTopListGuardNum|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		cacheHealth = false
		return
	} else if dhhResultFromRedis != nil {
		resp.TotalCount = int64(len(s.getAnchorTopGuardCount(req.Uid, dhhResultFromRedis)))
		exp.PromCacheHit(_promAnchorCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promAnchorCacheMissed)

	resultDB, err := s.dao.GetByAnchorUIDs(ctx, uids)
	resultDBTime, _ := s.changeDBTime(ctx, resultDB)
	if err != nil {
		reqAfterQueryDBTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDForGift|s.dao.dahanghai|从DB获取大航海信息error|(%v),missedUIDs(%v)|耗时:%dms", err, uids, reqAfterQueryDBTime-reqStartTime)
		return
	}
	dhhResultFromDB := s.formatDaHangHaiCache(req.Uid, resultDBTime)
	// 写入缓存
	if cacheHealth {
		s.asyncSetAnchorGuardCache(ctx, dhhResultFromDB, req.Uid)
	}
	resp.TotalCount = int64(len(s.getAnchorTopGuardCountFromDB(req.Uid, resultDBTime)))
	return
}

// GetByTargetIdsBatch ...
// 单uid全量守护
func (s *GuardService) GetByTargetIdsBatch(ctx context.Context, req *v1pb.GetByTargetIdsReq) (resp *v1pb.GetByTargetIdsResp, err error) {
	return
}

// GetByUIDForGift ...
// 单uid全量守护
func (s *GuardService) GetByUIDForGift(ctx context.Context, req *v1pb.GetByUidReq) (resp *v1pb.GetByUidResp, err error) {
	reqStartTime := confm.RecordTimeCost()
	resp = &v1pb.GetByUidResp{}
	resp.Data = make(map[int64]*v1pb.DaHangHaiInfo)
	cacheHealth := true
	uids := make([]int64, 0)
	uids = append(uids, req.Uid)
	dhhResultFromRedis, err := s.dao.GetUIDAllGuardFromRedis(ctx, uids)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|"+_XUserGetDaHangHaiForGiftCacheError+"|GetByUIDForGift|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		cacheHealth = false
		return
	} else if dhhResultFromRedis != nil {
		resp.Data = s.filterUIDTopFromRedis(req.Uid, dhhResultFromRedis)
		exp.PromCacheHit(_promCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promCacheMissed)

	resultDB, err := s.dao.GetByUIDs(ctx, uids)
	resultDBTime, _ := s.changeDBTime(ctx, resultDB)
	if err != nil {
		reqAfterQueryDBTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDForGift|s.dao.dahanghai|从DB获取大航海信息error|(%v),missedUIDs(%v)|耗时:%dms", err, uids, reqAfterQueryDBTime-reqStartTime)
		return
	}
	dhhResultFromDB := s.formatDaHangHaiCache(req.Uid, resultDBTime)
	// 写入缓存
	if cacheHealth {
		s.asyncSetDHHCache(ctx, dhhResultFromDB, req.Uid)
	}
	resp.Data = s.filterTopFromDB(req.Uid, resultDBTime)
	return
}

func (s *GuardService) transToMap(input []*dahanghaiModel.DaHangHaiRedis2) (resp map[int64]bool) {
	resp = make(map[int64]bool)
	if len(input) <= 0 {
		return
	}
	for _, v := range input {
		uid := RParseInt(v.Uid, 0)
		resp[uid] = true
	}
	return
}

func (s *GuardService) adaptRetForGetByUIDBatch(dhhList []*dahanghaiModel.DaHangHaiRedis2) (result map[int64]*v1pb.DaHangHaiInfoList) {
	result = make(map[int64]*v1pb.DaHangHaiInfoList)
	if len(dhhList) <= 0 {
		return
	}
	for _, v := range dhhList {
		uid := RParseInt(v.Uid, 0)
		if _, exist := result[uid]; !exist {
			result[uid] = &v1pb.DaHangHaiInfoList{}
			result[uid].List = make([]*v1pb.DaHangHaiInfo, 0)
		}
		item := &v1pb.DaHangHaiInfo{}
		item.Id = RParseInt(v.Id, 0)
		item.Uid = RParseInt(v.Uid, 0)
		item.TargetId = RParseInt(v.TargetId, 0)
		item.PrivilegeType = RParseInt(v.PrivilegeType, 0)
		item.Ctime = v.Ctime
		item.Utime = v.Utime
		item.ExpiredTime = v.ExpiredTime
		item.StartTime = v.StartTime
		result[uid].List = append(result[uid].List, item)
	}
	return
}

// GetByUIDBatch ...
// 多uid全量守护
func (s *GuardService) GetByUIDBatch(ctx context.Context, req *v1pb.GetByUidBatchReq) (resp *v1pb.GetByUidBatchResp, err error) {
	resp = &v1pb.GetByUidBatchResp{}
	if req == nil {
		return
	}
	reqStartTime := confm.RecordTimeCost()

	resp.Data = make(map[int64]*v1pb.DaHangHaiInfoList)
	cacheHealth := true
	dhhResultFromRedis, err := s.dao.GetUIDAllGuardFromRedisBatch(ctx, req.Uids)
	hitUIDs := s.transToMap(dhhResultFromRedis)
	missUIDS := s.findMissUIDs(hitUIDs, req.Uids)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDBatch|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		cacheHealth = false
		return
	} else if (dhhResultFromRedis != nil) && (len(missUIDS) <= 0) {
		resp.Data = s.adaptRetForGetByUIDBatch(dhhResultFromRedis)
		exp.PromCacheHit(_promCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promCacheMissed)
	resultDB, err := s.dao.GetByUIDsWithMap(ctx, missUIDS)
	resultDBTime, _ := s.changeDBTimeBatch(ctx, resultDB)
	if err != nil {
		reqAfterQueryDBTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDTargetIds|从DB获取大航海信息error|(%v),missedUIDs(%v)|耗时:%dms", err, req.Uids, reqAfterQueryDBTime-reqStartTime)
		return
	}
	dhhResultFromDB := s.formatDaHangHaiCacheBatch(resultDBTime)
	// 写入缓存
	if cacheHealth {
		s.asyncSetDHHCacheBatch(ctx, dhhResultFromDB)
	}
	resp.Data = s.mergeResultForGetByUIDBatch(dhhResultFromDB, dhhResultFromRedis)
	return
}

// GetByUIDTargetIds ...
// 单uid多targetids
func (s *GuardService) GetByUIDTargetIds(ctx context.Context, req *v1pb.GetByUidTargetIdsReq) (resp *v1pb.GetByUidTargetIdsResp, err error) {

	reqStartTime := confm.RecordTimeCost()

	resp = &v1pb.GetByUidTargetIdsResp{}
	resp.Data = make(map[int64]*v1pb.DaHangHaiInfo)
	cacheHealth := true
	uids := make([]int64, 0)
	uids = append(uids, req.Uid)
	dhhResultFromRedis, err := s.dao.GetUIDAllGuardFromRedis(ctx, uids)
	filterMap, _ := s.makeTargetFilterMap(ctx, req)

	// 回源db原则,仅在get成功且miss时回源db!!!
	if err != nil {
		reqAfterQueryMCTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDTargetIds|查询缓存失败,暂不回源db,接口返回err|%dms", err, reqAfterQueryMCTime-reqStartTime)
		cacheHealth = false
		return
	} else if dhhResultFromRedis != nil {
		resp.Data = s.filterResultFromRedis(req.Uid, filterMap, dhhResultFromRedis)
		exp.PromCacheHit(_promCacheHitAll)
		return
	}
	exp.PromCacheMiss(_promCacheMissed)

	resultDB, err := s.dao.GetByUIDs(ctx, uids)
	resultDBTime, _ := s.changeDBTime(ctx, resultDB)
	if err != nil {
		reqAfterQueryDBTime := confm.RecordTimeCost()
		log.Error(_errorServiceLogPrefix+"|GetByUIDTargetIds|从DB获取大航海信息error|(%v),missedUIDs(%v)|耗时:%dms", err, uids, reqAfterQueryDBTime-reqStartTime)
		return
	}
	dhhResultFromDB := s.formatDaHangHaiCache(req.Uid, resultDBTime)
	// 写入缓存
	if cacheHealth {
		s.asyncSetDHHCache(ctx, dhhResultFromDB, req.Uid)
	}
	resp.Data = s.filterResultFromDB(req.Uid, filterMap, resultDBTime)
	return

}

func (s *GuardService) changeDBTime(ctx context.Context, dhhs []*dhhm.DHHDB) (result []*dhhm.DHHDBTime, err error) {
	result = make([]*dhhm.DHHDBTime, 0)
	if len(dhhs) <= 0 {
		return
	}
	for _, v := range dhhs {
		item := &dhhm.DHHDBTime{}
		item.ID = v.ID
		item.UID = v.UID
		item.TargetId = v.TargetId
		item.PrivilegeType = v.PrivilegeType
		item.StartTime = v.StartTime.Format(_formatTime)
		item.ExpiredTime = v.ExpiredTime.Format(_formatTime)
		item.Ctime = v.Ctime.Format(_formatTime)
		item.Utime = v.Utime.Format(_formatTime)
		result = append(result, item)
	}
	return
}

func (s *GuardService) changeDBTimeBatch(ctx context.Context, dhhs map[int64][]*dhhm.DHHDB) (result map[int64][]*dhhm.DHHDBTime, err error) {
	result = make(map[int64][]*dhhm.DHHDBTime)
	if len(dhhs) <= 0 {
		return
	}
	for k, v := range dhhs {
		resultList := make([]*dhhm.DHHDBTime, 0)
		for _, vv := range v {
			item := &dhhm.DHHDBTime{}
			item.ID = vv.ID
			item.UID = vv.UID
			item.TargetId = vv.TargetId
			item.PrivilegeType = vv.PrivilegeType
			item.StartTime = vv.StartTime.Format(_formatTime)
			item.ExpiredTime = vv.ExpiredTime.Format(_formatTime)
			item.Ctime = vv.Ctime.Format(_formatTime)
			item.Utime = vv.Utime.Format(_formatTime)
			resultList = append(resultList, item)
		}
		result[k] = make([]*dhhm.DHHDBTime, 0)
		result[k] = resultList
	}
	return
}

func (s *GuardService) makeTargetFilterMap(ctx context.Context, req *v1pb.GetByUidTargetIdsReq) (resp map[int64]int64, err error) {
	resp = make(map[int64]int64)
	if req == nil {
		return
	}
	for _, v := range req.TargetIDs {
		resp[v.TargetId] = v.SortType
	}
	return
}
