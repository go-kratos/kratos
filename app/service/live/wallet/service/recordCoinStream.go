package service

import (
	"context"
	"encoding/json"
	"go-common/app/service/live/wallet/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

type RecordCoinStreamHandler struct {
}

func (handler *RecordCoinStreamHandler) NeedCheckUid() bool {
	return true
}

const ItemNumCount = 50

func (handler *RecordCoinStreamHandler) NeedTransactionMutex() bool {
	return false
}
func (handler *RecordCoinStreamHandler) BizExecute(ws *WalletService, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	platform, _ := params[0].(string)
	if !model.IsValidPlatform(platform) {
		err = ecode.RequestErr
		return
	}
	arg, _ := params[1].(*model.RecordCoinStreamForm)

	var recordCoinStreamItems []*model.RecordCoinStreamItem

	jsonErr := json.Unmarshal([]byte(arg.Data), &recordCoinStreamItems)
	if jsonErr != nil || len(recordCoinStreamItems) == 0 || len(recordCoinStreamItems) > ItemNumCount {
		err = ecode.RequestErr
		return
	}

	for _, m := range recordCoinStreamItems {
		if !m.IsValid() {
			err = ecode.RequestErr
			return
		}
	}

	goldPay := 0
	goldRecharge := 0
	silverPay := 0
	var detail *model.DetailWithSnapShot
	_, err = ws.s.dao.DoTx(ws.c, func(conn *sql.Tx) (v interface{}, err error) {
		detail, err = ws.s.dao.WalletForUpdate(conn, uid)
		if err != nil {
			return
		}
		for _, m := range recordCoinStreamItems {

			sysCoinType := model.GetSysCoinType(m.CoinType, platform)
			sysCoinTypeNo := model.GetCoinTypeNumber(sysCoinType)
			record := model.NewCoinStream(uid, m.TransactionId, m.ExtendTid, sysCoinTypeNo, m.CoinNum, m.GetOpType(),
				m.Timestamp, basicParam.BizCode, basicParam.Area, basicParam.Source, basicParam.BizSource, basicParam.MetaData)
			model.AddMoreParam2CoinStream(record, basicParam, platform)
			record.Reserved1 = m.Reserved1
			record.OrgCoinNum = m.GetOrgCoinNum()
			record.OpResult = m.GetOpResult()
			_, err = ws.s.dao.NewCoinStreamRecordInTx(conn, record)
			if err != nil {
				break
			}
			if m.IsPayType() {
				if sysCoinTypeNo == model.SysCoinTypeIosGold || sysCoinTypeNo == model.SysCoinTypeGold {
					goldPay = goldPay + int(0-m.CoinNum)
				} else if sysCoinTypeNo == model.SysCoinTypeSilver {
					silverPay = silverPay + int(0-m.CoinNum)
				}
			} else if m.IsRechargeType() && (sysCoinTypeNo == model.SysCoinTypeIosGold || sysCoinTypeNo == model.SysCoinTypeGold) {
				goldRecharge = goldRecharge + int(m.CoinNum)
			}
		}
		_, err = ws.s.dao.ModifyCntInTx(conn, uid, goldPay, goldRecharge, silverPay)
		return
	})

	if err != nil {
		return
	}
	// 事务执行以后　需要发databus 原因在于成就系统依赖于该消息来算用户消费总数如果不发则成就系统无法及时知道用户的消费总数增加
	// 后续成就系统可以通过消费cannal 的方式来重构　后则无需发消息队列
	// ps : 流水可能有多条　但是消息队列只发一条
	var action string
	var number int64
	var coinType string
	needPub := true
	if goldPay > 0 || silverPay > 0 {
		action = "pay"
		if goldPay > 0 {
			coinType = "gold"
			number = int64(goldPay)
		} else {
			coinType = "silver"
			number = int64(silverPay)
		}
	} else if goldRecharge > 0 {
		coinType = "gold"
		number = int64(goldRecharge)
		action = "recharge"
	} else {
		needPub = false
	}
	detail.GoldRechargeCnt += int64(goldRecharge)
	detail.GoldPayCnt += int64(goldPay)
	detail.SilverPayCnt += int64(silverPay)
	if needPub {
		log.Info("RecordCoinStream#ExecSuccess#Pub#action:%s#coinType:%s#number:%d#uid:%d", action, coinType, number, uid)
		ws.s.pubWalletChangeWithDetailSnapShot(ws.c,
			uid, action, number, coinType, platform, "", 0, detail)
	}

	return
}

func (s *Service) RecordCoinStream(c context.Context, basicParam *model.BasicParam, uid int64, params ...interface{}) (v interface{}, err error) {
	handler := RecordCoinStreamHandler{}
	return s.execByHandler(&handler, c, basicParam, uid, params...)
}
