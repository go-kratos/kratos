package service

import (
	"context"

	"go-common/app/admin/main/reply/model"
	"go-common/library/log"
)

// EmojiPackageList return all emojipackages and emojis
func (s *Service) EmojiPackageList(c context.Context) (resp []*model.EmojiPackage, err error) {
	packs, err := s.dao.EmojiPackageList(c)
	if err != nil {
		log.Error("service.EmojiPackageList err (%v)", err)
		return
	}
	for _, pack := range packs {
		eList, err := s.dao.EmojiListByPid(c, pack.ID)
		if err != nil {
			return nil, err
		}
		pack.Emojis = eList
		resp = append(resp, pack)
	}
	return
}

// CreateEmojiPackage CreateEmojiPackage
func (s *Service) CreateEmojiPackage(c context.Context, name string, url string, sort int32, remark string, state int32) (id int64, err error) {
	id, err = s.dao.CreateEmojiPackage(c, name, url, sort, remark, state)
	if err != nil {
		log.Error("service.CreateEmojiPackage err (%v)", err)
	}
	return
}

// UpEmojiPackageSort UpEmojiPackageSort
func (s *Service) UpEmojiPackageSort(c context.Context, ids string) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	err = s.dao.UpEmojiPackageSort(tx, ids)
	if err != nil {
		tx.Rollback()
		log.Error("service.UpEmojiPackageSort err (%v)", err)
		return
	}
	return tx.Commit()
}

// UpEmojiPackage UpEmojiPackage
func (s *Service) UpEmojiPackage(c context.Context, name string, url string, remark string, state int32, id int64) (idx int64, err error) {
	if state == 2 { // delete package
		_, err = s.dao.DelEmojiPackage(c, id)
		if err != nil {
			log.Error("service.DelEmojiPackage err (%v)", err)
			return 0, err
		}
		idx, err = s.dao.DelEmojiByPid(c, id)
	} else {
		idx, err = s.dao.UpEmojiPackage(c, name, url, remark, state, id)
	}
	if err != nil {
		log.Error("service.UpEmojiPackage err (%v)", err)
	}
	return
}
