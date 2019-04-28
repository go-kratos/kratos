package service

import (
	"context"

	"go-common/app/job/main/account-summary/model"
	"go-common/library/log"
)

// Syncable is
type Syncable interface {
	Key() (string, error)
	Marshal() (map[string][]byte, error)
}

// SyncToHBase is
func (s *Service) SyncToHBase(ctx context.Context, body Syncable) error {
	key, err := body.Key()
	if err != nil {
		log.Error("Failed to struct key with body: %+v: %+v", body, err)
		return err
	}
	data, err := body.Marshal()
	if err != nil {
		log.Error("Failed to sync to hbase with key: %s, body: %+v: %+v", key, body, err)
		return err
	}
	defer s.logging(ctx, key)
	return s.dao.Save(ctx, key, data)
}

func (s *Service) logging(ctx context.Context, key string) {
	sum, _ := s.dao.GetByKey(ctx, key)
	log.Info("Sync to hbase result: key: %s, summary: %+v", key, sum)
}

// SyncOne is
func (s *Service) SyncOne(ctx context.Context, mid int64) error {
	// member
	if err := s.syncMember(ctx, mid); err != nil {
		log.Error("Failed to sync member with mid: %d: %+v", mid, err)
	}

	// relation
	if err := s.syncRelationStat(ctx, mid); err != nil {
		log.Error("Failed to sync relation stat with mid: %d: %+v", mid, err)
	}

	// block
	if err := s.syncBlock(ctx, mid); err != nil {
		log.Error("Failed to sync block with mid: %d: %+v", mid, err)
	}

	// passport
	if err := s.syncPassportSummary(ctx, mid); err != nil {
		log.Error("Failed to sync passport summary with mid: %d: %+v", mid, err)
	}

	return nil
}

// GetOne is
func (s *Service) GetOne(ctx context.Context, mid int64) (*model.AccountSummary, error) {
	return s.dao.GetByKey(ctx, model.MidKey(mid))
}
