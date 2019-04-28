package v1

import (
	"context"
	"go-common/library/ecode"
	"go-common/library/net/metadata"

	v1pb "go-common/app/service/live/xuser/api/grpc/v1"
	expm "go-common/app/service/live/xuser/model/exp"

	XanchorV1 "go-common/app/service/live/dao-anchor/api/grpc/v1"
	"go-common/library/log"
)

func (s *UserExpService) asyncUpdateDBUexp(ctx context.Context, req *v1pb.UserExpChunk) error {
	ctxNodeadline := metadata.WithContext(ctx)
	f := func(c context.Context) {
		reqStartTime := s.RecordTimeCost()
		retryTime := 1
		_, err := s.dao.AddUexp(c, req.Uid, req.Num)
		if err != nil {
			for ; retryTime <= _retryAddExpTimes; retryTime++ {
				_, err = s.dao.AddUexp(c, req.Uid, req.Num)
				if err != nil {
					break
				}
			}
			if retryTime >= _retryAddExpTimes {
				reqUpdateDBFail := s.RecordTimeCost()
				log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_USEREXP_UPDATEDB_RETRYALL_FAIL+"|asyncUpdateDBUexp|DB error(%v),params(%v)|耗时:%dms", err, req, reqUpdateDBFail-reqStartTime)
				err = ecode.XUserAddUserExpUpdateDBError
				return
			}
			log.Info(_errorServiceLogPrefix+"|asyncUpdateDBUexp"+_INFO_ADD_USEREXP_RETRY_UPDATEDB+"|DB retry %d times success,params(%v)", retryTime, req)
		}

		s.asyncCLearExpCache(c, req)
		QueryConfig := s.getQueryStatus()
		if QueryConfig == 1 {
			reqBeforeQueryDBTime := s.RecordTimeCost()
			expResult, err := s.dao.Exp(c, req.Uid)
			reqAfterQueryDBTime := s.RecordTimeCost()
			if err != nil {
				log.Error(_errorServiceLogPrefix+"|asyncUpdateDBUexp|"+_ERROR_QUERY_AFTER_ADD_USEREXP_UPDATEDB_FAIL+"|更新db后再查询error|Query error(%v),params(%v),resp(%v)|耗时:%dms", err, req, expResult, reqAfterQueryDBTime-reqBeforeQueryDBTime)
			} else {
				expResultList := make([]*expm.Exp, 0)
				expResultList = append(expResultList, expResult)
				item := s.FormatLevel(expResultList)
				s.asyncSetExpCache(c, item)
				logDesc, ok := _addExpReqBizMap[req.ReqBiz]
				if !ok {
					logDesc = "默认渠道"
				}
				s.addUserExpLog(req.Uid, req.Num, item, logDesc)
			}
		}
	}
	if runErr := s.dbUpdater.Do(ctxNodeadline, func(c context.Context) {
		f(ctxNodeadline)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ASYNADD_UEXP_NUM_FAIL+"|asyncSetExpCache|error(%v),run cache is full(%v)", runErr)
		f(ctxNodeadline)
	}
	return nil
}

func (s *UserExpService) asyncUpdateDBRexp(ctx context.Context, req *v1pb.UserExpChunk) error {
	ctxNodeadline := metadata.WithContext(ctx)
	f := func(c context.Context) {
		reqStartTime := s.RecordTimeCost()
		retryTime := 1
		_, err := s.dao.AddRexp(c, req.Uid, req.Num)
		if err != nil {
			for ; retryTime <= _retryAddExpTimes; retryTime++ {
				_, err = s.dao.AddRexp(c, req.Uid, req.Num)
				if err != nil {
					break
				}
			}
			if retryTime >= _retryAddExpTimes {
				reqUpdateDBFail := s.RecordTimeCost()
				log.Error(_errorServiceLogPrefix+"|"+_ERROR_ADD_USEREXP_UPDATEDB_RETRYALL_FAIL+"|asyncUpdateDBRexp|DB error(%v),params(%v)|耗时:%dms", err, req, reqUpdateDBFail-reqStartTime)
				err = ecode.XUserAddUserExpUpdateDBError
				return
			}
			log.Info(_errorServiceLogPrefix+"|asyncUpdateDBRexp"+_INFO_ADD_USEREXP_RETRY_UPDATEDB+"|DB retry %d times success,params(%v)", retryTime, req)
		}
		params := &XanchorV1.AnchorIncreReq{ReqId: "AnchorIncreReq", Fields: []string{"exp"}, Uid: req.Uid, Exp: req.Num}
		err = s.xuserDao.UpdateAnchorInfo(c, params)
		if err != nil {
			reqUpdateDaoDBFail := s.RecordTimeCost()
			log.Error(_errorServiceLogPrefix+"|asyncUpdateDBRexp|"+_ERROR_ASYNADD_REXP_NUM_FAIL+"|主播经验双写error|Query error(%v),params(%v)|耗时:%dms", err, req, reqUpdateDaoDBFail-reqStartTime)
		}

		s.asyncCLearExpCache(c, req)
		QueryConfig := s.getQueryStatus()
		if QueryConfig == 1 {
			reqBeforeQueryDBTime := s.RecordTimeCost()
			expResult, err := s.dao.Exp(c, req.Uid)
			reqAfterQueryDBTime := s.RecordTimeCost()
			if err != nil {
				log.Error(_errorServiceLogPrefix+"|asyncUpdateDBRexp|"+_ERROR_QUERY_AFTER_ADD_USEREXP_UPDATEDB_FAIL+"|更新db后再查询error|Query error(%v),params(%v),resp(%v)|耗时:%dms", err, req, expResult, reqAfterQueryDBTime-reqBeforeQueryDBTime)
			} else {
				expResultList := make([]*expm.Exp, 0)
				expResultList = append(expResultList, expResult)
				item := s.FormatLevel(expResultList)
				s.asyncSetExpCache(c, item)
				logDesc, ok := _addExpReqBizMap[req.ReqBiz]
				if !ok {
					logDesc = "默认渠道"
				}
				s.addAnchorExpLog(req.Uid, req.Num, item, logDesc)
			}
		}
	}
	if runErr := s.dbUpdater.Do(ctxNodeadline, func(c context.Context) {
		f(ctxNodeadline)
	}); runErr != nil {
		log.Error(_errorServiceLogPrefix+"|"+_ERROR_ASYNADD_UEXP_NUM_FAIL+"|asyncSetExpCache|error(%v),run cache is full(%v)", runErr)
		f(ctxNodeadline)
	}
	return nil
}
