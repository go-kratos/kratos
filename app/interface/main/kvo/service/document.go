package service

import (
	"context"
	"encoding/json"
	"time"

	"go-common/app/interface/main/kvo/model"
	"go-common/app/interface/main/kvo/model/module"

	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

// Document get document
func (s *Service) Document(c context.Context, mid int64, moduleKey string, timestamp int64, checkSum int64) (setting *module.Setting, err error) {
	var (
		uc          *model.UserConf
		rm          json.RawMessage
		moduleKeyID int
	)
	if moduleKeyID = module.VerifyModuleKey(moduleKey); moduleKeyID == 0 {
		err = ecode.RequestErr
		return
	}
	uc, err = s.userConf(c, mid, moduleKeyID)
	if err != nil {
		return
	}
	if uc.CheckSum == 0 || uc.Timestamp == 0 {
		err = ecode.NotModified
		return
	}
	// 数据没有变动
	if uc.CheckSum == checkSum && uc.Timestamp == timestamp {
		err = ecode.NotModified
		return
	}
	rm, err = s.document(c, uc.CheckSum)
	if err != nil {
		return
	}
	setting = &module.Setting{
		Timestamp: uc.Timestamp,
		CheckSum:  uc.CheckSum,
		Data:      rm,
	}
	return
}

func (s *Service) userConf(c context.Context, mid int64, moduleKeyID int) (uc *model.UserConf, err error) {
	uc, err = s.da.UserConfCache(c, mid, moduleKeyID)
	if err != nil {
		log.Error("service.userConf.UserConfCache(%v,%v) err:%v", mid, moduleKeyID, err)
	}
	if uc != nil {
		s.sp.Incr("user_conf_cached")
		return
	}
	uc, err = s.da.UserConf(c, mid, moduleKeyID)
	if err != nil {
		log.Error("service.userConf(%v,%v) err:%v", mid, moduleKeyID, err)
		return
	}
	if uc == nil {
		uc = &model.UserConf{
			Mid:       mid,
			ModuleKey: moduleKeyID,
		}
		s.sp.Incr("default_user_conf")
	}
	s.sp.Incr("user_conf_missed")
	s.updateUcCache(mid, moduleKeyID)
	return
}

func (s *Service) document(c context.Context, checkSum int64) (rm json.RawMessage, err error) {
	var (
		doc *model.Document
	)
	rm, err = s.da.DocumentCache(c, checkSum)
	if err != nil {
		log.Error("service.document.DocumentCache(%v) err:%v", checkSum, err)
	}
	if rm != nil {
		s.sp.Incr("document_cached")
		return
	}
	doc, err = s.da.Document(c, checkSum)
	if err != nil {
		log.Error("service.document(%v) err:%v", checkSum, err)
		return
	}
	if doc == nil {
		err = ecode.NothingFound
		s.sp.Incr("user_conf_document_error")
		return
	}
	s.sp.Incr("document_missed")
	rm = json.RawMessage(doc.Doc)
	s.da.SetDocumentCache(c, checkSum, rm)
	return
}

// AddDocument add a user document
func (s *Service) AddDocument(c context.Context, mid int64, moduleKey string, data string, timestamp int64, oldSum int64, now time.Time) (resp *model.UserConf, err error) {
	var (
		uc          *model.UserConf
		doc         *model.Document
		rm          json.RawMessage
		checkSum    int64
		tx          *sql.Tx
		moduleKeyID int
	)
	if moduleKeyID = module.VerifyModuleKey(moduleKey); moduleKeyID == 0 {
		return nil, ecode.RequestErr
	}
	if rm, checkSum, err = module.Result(moduleKeyID, data); err != nil {
		log.Error("service.GetModule(%v,%s) err:%v", moduleKey, data, err)
		return nil, ecode.RequestErr
	}
	if len(rm) > s.docLimit {
		err = ecode.KvoDataOverLimit
		return
	}
	if uc, err = s.da.UserConf(c, mid, moduleKeyID); err != nil {
		log.Error("service.AddDocument.UserConf(%v,%v) err:%v", mid, moduleKeyID, err)
		return
	}
	s.updateUcCache(mid, moduleKeyID)
	if uc != nil {
		if uc.Timestamp != timestamp {
			err = ecode.KvoTimestampErr
			log.Error("service.AddDocument.CompareTimeStamp(%v,%v,%v) err:%v", mid, uc.Timestamp, timestamp, err)
			return
		}
		if uc.CheckSum != oldSum {
			err = ecode.KvoCheckSumErr
			return
		}
		if uc.CheckSum == checkSum {
			err = ecode.NotModified
			return
		}
	}
	// trans
	tx, err = s.da.BeginTx(c)
	if err != nil {
		log.Error("s.da.BeginTx err:%v", err)
		return
	}
	if err = s.da.TxUpUserConf(c, tx, mid, moduleKeyID, checkSum, now); err != nil {
		log.Error("s.da.TxUpUserConf(%v,%v,%v) error(%v)", mid, moduleKeyID, checkSum, err)
		tx.Rollback()
		return
	}
	doc, err = s.da.Document(c, checkSum)
	if err != nil {
		tx.Rollback()
		return
	}
	if doc == nil {
		if err = s.da.TxUpDocuement(c, tx, checkSum, string(rm), now); err != nil {
			log.Error("s.da.TxUpDocuement(%v,%v,%v) error(%v)", mid, moduleKeyID, checkSum, err)
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit(), error(%v)", err)
	}
	resp = &model.UserConf{
		CheckSum:  checkSum,
		Timestamp: now.Unix(),
	}
	s.updateUcCache(mid, moduleKeyID)
	return
}
