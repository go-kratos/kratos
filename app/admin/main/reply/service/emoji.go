package service

import (
	"context"
	"strings"

	"go-common/app/admin/main/reply/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// EmojiList EmojiList
func (s *Service) EmojiList(c context.Context, pid int64) (emojis []*model.Emoji, err error) {
	if pid != 0 {
		emojis, err = s.dao.EmojiListByPid(c, pid)
	} else {
		emojis, err = s.dao.EmojiList(c)
	}
	if err != nil {
		log.Error("service.EmojiList error (%v)", err)
	}
	return
}

// CreateEmoji CreateEmoji
func (s *Service) CreateEmoji(c context.Context, pid int64, name string, url string, sort int32, state int32, remark string) (id int64, err error) {
	if !strings.HasPrefix(name, "[") {
		name = "[" + name
	}
	if !strings.HasSuffix(name, "]") {
		name = name + "]"
	}
	id, err = s.dao.CreateEmoji(c, pid, name, url, sort, state, remark)
	if err != nil {
		log.Error("service.CreateEmoji error (%v)", err)
	}
	return
}

// UpEmojiSort UpEmojiSort
func (s *Service) UpEmojiSort(c context.Context, ids string) (err error) {
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		return
	}
	err = s.dao.UpEmojiSort(tx, ids)
	if err != nil {
		tx.Rollback()
		log.Error("service.UpEmojiSort error (%v)", err)
		return
	}
	return tx.Commit()
}

// UpEmojiState UpEmojiState
func (s *Service) UpEmojiState(c context.Context, state int32, id int64) (idx int64, err error) {
	if state == 2 { // delete emoji
		idx, err = s.dao.DelEmojiByID(c, id)
	} else {
		idx, err = s.dao.UpEmojiStateByID(c, state, id)
	}
	if err != nil {
		log.Error("service.UpEmojiState error (%v)", err)
	}
	return
}

// UpEmoji UpEmoji
func (s *Service) UpEmoji(c context.Context, name string, remark string, url string, id int64) (idx int64, err error) {
	idx, err = s.dao.UpEmoji(c, name, remark, url, id)
	if err != nil {
		log.Error("service.UpEmojiState error (%v)", err)
	}
	return
}

// EmojiByName EmojiByName
func (s *Service) EmojiByName(c context.Context, name string) (err error) {
	emojis, e := s.dao.EmojiByName(c, name)
	if e != nil {
		log.Error("service.CreateEmoji EmojiByName err (%v)", e)
		err = e
		return
	}
	if len(emojis) > 0 {
		err = ecode.ReplyEmojiExits
		return
	}
	return
}
