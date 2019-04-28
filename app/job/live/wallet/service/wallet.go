package Service

import (
	"context"
	"encoding/json"
	"go-common/app/job/live/wallet/model"
	"go-common/library/log"
)

func (s *Service) mergeData(nwMsg []byte, oldMsg []byte, action string) {
	if action == "update" {
		// do noting
	} else if action == "insert" {
		userNew := &model.User{}
		if err := json.Unmarshal(nwMsg, userNew); err != nil {
			log.Error("json.Unmarshal(%s) error(%v)", string(nwMsg), err)
			return
		}
		log.Info("new user %d", userNew.Uid)
		s.dao.InitWallet(context.TODO(), userNew)
		//s.dao.InitWallet(context.TODO(), userNew.Uid, userNew.Gold, userNew.IapGold, userNew.Silver)
		s.dao.DelWalletCache(context.TODO(), userNew.Uid)
	}
}
