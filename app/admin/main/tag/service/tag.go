package service

import (
	"context"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	_emtpyTagInfo  = new(model.TagInfo)
	_emtpyTagCount = new(model.TagCount)
	_emtpyTag      = new(model.Tag)
)

// TagList TagList.
func (s *Service) TagList(c context.Context, esTag *model.ESTag) (res *model.MngSearchTagList, err error) {
	return s.dao.ESearchTag(c, esTag)
}

// TagEdit tag edit.
func (s *Service) TagEdit(c context.Context, tid int64, tp int32, tname, content string) (err error) {
	var tag *model.Tag
	if err = s.dao.Filter(c, tname); err != nil {
		return
	}
	if tid > 0 {
		if tag, err = s.dao.Tag(c, tid); err != nil {
			return
		}
		if tag == nil {
			return ecode.TagNotExist
		}
		if _, err = s.editTag(c, tid, tp, tname, content); err != nil {
			log.Error("s.editTag(%d,%s) error(%v)", tid, tname, err)
		}
		return
	}
	_, err = s.addTag(c, tp, tname, content)
	return
}

func (s *Service) editTag(c context.Context, tid int64, tp int32, tname, content string) (id int64, err error) {
	var tag *model.Tag
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil {
		err = ecode.TagNotExist
		return
	}
	if tag.ID != tid || tag.Name != tname {
		return
	}
	tag.Content = content
	tag.Type = tp
	if id, err = s.dao.UpdateTag(c, tag); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelTagCache(ctx, tid, tname)
		s.dao.UpdateESearchTag(ctx, tag)
	})
	return
}

func (s *Service) addTag(c context.Context, tp int32, tname, content string) (tid int64, err error) {
	var (
		tag    *model.Tag
		affect int64
	)
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag != nil {
		return tag.ID, ecode.TagAlreadyExist
	}
	tag = &model.Tag{
		Name:    tname,
		Type:    model.TypeBiliContent,
		Content: content,
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("add tag tran error(%v)", err)
		return
	}
	if tag.ID, err = s.dao.TxInsertTag(tx, tag); err != nil || tag.ID <= 0 {
		err = ecode.TagAddFail
		tx.Rollback()
		return
	}
	if affect, err = s.dao.TxInsertTagCount(tx, tag.ID); err != nil || affect <= 0 {
		err = ecode.TagAddFail
		tx.Rollback()
		return
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.UpdateESearchTag(ctx, tag)
	})
	return
}

// TagInfo tag detail info, include tag info and synonym tag infos, used by tid.
func (s *Service) TagInfo(c context.Context, tid int64) (res *model.TagInfo, err error) {
	var (
		tag   *model.Tag
		count *model.TagCount
	)
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return _emtpyTagInfo, ecode.TagNotExist
	}
	count, _ = s.dao.TagCount(c, tid)
	if count == nil {
		count = _emtpyTagCount
	}
	res = &model.TagInfo{
		ID:      tag.ID,
		Type:    tag.Type,
		Name:    tag.Name,
		Cover:   tag.Cover,
		Content: tag.Content,
		Verify:  tag.Verify,
		Attr:    tag.Attr,
		State:   tag.State,
		Bind:    count.Bind,
		Sub:     count.Sub,
		CTime:   tag.CTime,
		MTime:   tag.MTime,
	}
	return
}

// TagInfoByName TagInfoByName.
func (s *Service) TagInfoByName(c context.Context, tname string) (res *model.TagInfo, err error) {
	var (
		tag   *model.Tag
		count *model.TagCount
	)
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil {
		return _emtpyTagInfo, ecode.TagNotExist
	}
	if count, _ = s.dao.TagCount(c, tag.ID); count == nil {
		count = _emtpyTagCount
	}
	res = &model.TagInfo{
		ID:      tag.ID,
		Type:    tag.Type,
		Name:    tag.Name,
		Cover:   tag.Cover,
		Content: tag.Content,
		Verify:  tag.Verify,
		Attr:    tag.Attr,
		State:   tag.State,
		Bind:    count.Bind,
		Sub:     count.Sub,
		CTime:   tag.CTime,
		MTime:   tag.MTime,
	}
	return
}

// TagState tag state.
func (s *Service) TagState(c context.Context, tid int64, state int32) (err error) {
	if _, err := s.dao.UpTagState(c, tid, state); err != nil {
		return ecode.TagChangeStateFail
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		tag, err := s.dao.Tag(ctx, tid)
		if err == nil && tag != nil {
			s.dao.DelTagCache(ctx, tid, tag.Name)
		}
		s.dao.UpdateESearchTag(ctx, &model.Tag{ID: tid, State: state, Verify: model.VerifyUnknown, Type: model.TypeUnknow})
	})
	return
}

// TagVerify TagVerify.
func (s *Service) TagVerify(c context.Context, tid int64) (err error) {
	var (
		affect int64
		tag    *model.Tag
	)
	if tag, err = s.dao.Tag(c, tid); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	if affect, err = s.dao.UpVerifyState(c, tid, model.VerifyDone); err != nil {
		return
	}
	if affect <= 0 {
		return ecode.TagAlreadyExamined
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelTagCache(ctx, tid, tag.Name)
		tag.Verify = model.VerifyDone
		s.dao.UpdateESearchTag(ctx, tag)
	})
	return
}

// TagCheck tag check.
func (s *Service) TagCheck(c context.Context, tid int64, tname string) (res *model.Tag, err error) {
	if tid <= 0 {
		res, err = s.dao.TagByName(c, tname)
	} else {
		res, err = s.dao.Tag(c, tid)
	}
	if err != nil {
		return
	}
	if res == nil || res.State != model.StateNormal {
		err = ecode.TagNotExist
	}
	return
}
