package service

import (
	"context"
	"database/sql"
	"go-common/app/service/live/wallet/model"
	"go-common/library/ecode"
)

const TX_STATUS_SUCC = 0   //事务执行成功
const TX_STATUS_FAILED = 1 //事务执行失败
type QueryHandler struct {
}

func (handler *QueryHandler) NeedCheckUid() bool {
	return false
}

func (handler *QueryHandler) NeedTransactionMutex() bool {
	return false
}
func (handler *QueryHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	tid, _ := params[1].(string)
	var record *model.CoinStreamRecord
	if uid == 0 {
		record, err = ws.s.dao.GetCoinStreamByTid(ws.c, tid)
	} else {
		record, err = ws.s.dao.GetCoinStreamByUidTid(ws.c, uid, tid)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			err = ecode.NothingFound
		} else {
			err = ecode.ServerErr
		}
		return
	}

	var (
		result bool
		opType model.ServiceType
	)

	opType = model.ServiceType(record.OpType)
	switch opType {
	case model.PAYTYPE:
		result = record.OpReason == 0
	case model.RECHARGETYPE:
		result = record.OpReason == 0 || record.OpResult == model.STREAM_OP_REASON_POST_QUERY_FAILED
	case model.EXCHANGETYPE:
		result = record.OpReason == 0 || record.OpResult == model.STREAM_OP_REASON_EXECUTE_UNKNOWN
	default:
	}

	res := &model.QueryResp{}
	if result {
		res.Status = TX_STATUS_SUCC
	} else {
		res.Status = TX_STATUS_FAILED
	}

	v = res

	return
}

func (s *Service) Query(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	handler := QueryHandler{}
	return s.execByHandler(&handler, c, basicParam, uid, params...)
}
