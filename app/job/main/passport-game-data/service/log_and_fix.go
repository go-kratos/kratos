package service

import (
	"context"
	"encoding/json"

	"go-common/app/job/main/passport-game-data/model"
	"go-common/library/log"
)

func (s *Service) fixCloudRecord(c context.Context, newRecord *model.AsoAccount, old *model.AsoAccount) {
	if old == nil {
		affected, err := s.d.AddIgnoreAsoAccount(c, newRecord)
		if err != nil {
			oldStr, _ := json.Marshal(old)
			newStr, _ := json.Marshal(newRecord)
			log.Error("failed to fix cloud record by adding,  old(%s) new(%s) error(%v)", oldStr, newStr, err)
			return
		}
		if affected == 0 {
			oldStr, _ := json.Marshal(old)
			newStr, _ := json.Marshal(newRecord)
			log.Error("failed to fix cloud record by adding because of concurrent update, old(%s) new(%s)", oldStr, newStr)
			return
		}

		oldStr, _ := json.Marshal(old)
		newStr, _ := json.Marshal(newRecord)
		log.Info("fix cloud record by adding ok, old(%s) new(%s)", oldStr, newStr)
		return
	}
	affected, err := s.d.UpdateAsoAccountCloud(c, newRecord, old.Mtime)
	if err != nil {
		oldStr, _ := json.Marshal(old)
		newStr, _ := json.Marshal(newRecord)
		log.Error("failed to fix cloud record by updating, old(%s) new(%s) error(%v)", oldStr, newStr, err)
		return
	}
	if affected == 0 {
		oldStr, _ := json.Marshal(old)
		newStr, _ := json.Marshal(newRecord)
		log.Error("failed to fix cloud record by updating because of concurrent update, old(%s) new(%s)", oldStr, newStr)
		return
	}

	oldStr, _ := json.Marshal(old)
	newStr, _ := json.Marshal(newRecord)
	log.Info("fix cloud record by updating ok, old(%s) new(%s)", oldStr, newStr)
}

func (s *Service) doLog(cloud *model.AsoAccount, local *model.OriginAsoAccount, afterPending bool) {
	localStr := []byte("nil")
	localEncStr := []byte("nil")
	if local != nil {
		localStr, _ = json.Marshal(local)
		localEncStr, _ = json.Marshal(model.Default(local))
	}
	cloudStr := []byte("nil")
	if cloud != nil {
		cloudStr, _ = json.Marshal(cloud)
	}
	if afterPending {
		log.Info("failed to compare, because cloud record is not updated in time, local(%s) local_encrypted(%s) cloud(%s)", localStr, localEncStr, cloudStr)
		return
	}
	log.Info("compare diff, local(%s) local_encrypted(%s) cloud(%s)", localStr, localEncStr, cloudStr)
}
