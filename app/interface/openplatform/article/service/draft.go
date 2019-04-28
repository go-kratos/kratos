package service

import (
	"context"

	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s Service) checkDraftAuthor(c context.Context, aid, mid int64) (d *artmdl.Draft, err error) {
	if d, err = s.dao.ArtDraft(c, mid, aid); err != nil {
		return
	}
	if d == nil {
		err = ecode.NothingFound
		return
	}
	if d.Author.Mid != mid {
		err = ecode.ArtCreationMIDErr
	}
	return
}

// ArtDraft get article draft by id.
func (s *Service) ArtDraft(c context.Context, aid, mid int64) (res *artmdl.Draft, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	if res, err = s.checkDraftAuthor(c, aid, mid); err != nil {
		return
	}
	if res.ListID > 0 {
		res.List, _ = s.dao.List(c, res.ListID)
	}
	return
}

// AddArtDraft add article draft .
func (s *Service) AddArtDraft(c context.Context, a *artmdl.Draft) (id int64, err error) {
	log.Infov(c, log.KV("AddArtDraft", a))
	if err = s.checkPrivilege(c, a.Author.Mid); err != nil {
		return
	}
	a.Content = xssFilter(a.Content)
	if err = s.preDraftCheck(c, a); err != nil {
		return
	}
	if _, err = s.checkList(c, a.Author.Mid, a.ListID); err != nil {
		return
	}
	// if a.ListID > 0 {
	// 	var exist bool
	// 	for _, cid := range s.novelCIDs() {
	// 		if cid == a.Category.ID {
	// 			exist = true
	// 			break
	// 		}
	// 	}
	// 	if !exist {
	// 		a.ListID = 0
	// 	}
	// }
	if a.ID > 0 {
		if _, err = s.checkDraftAuthor(c, a.ID, a.Author.Mid); err != nil {
			return
		}
		_, err = s.dao.AddArtDraft(c, a)
		id = a.ID
		log.Info("update draft success mid(%d) aid(%d)", a.Author.Mid, a.ID)
		return
	}
	var total int
	if total, err = s.dao.CountUpperDraft(c, a.Author.Mid); err != nil {
		return
	}
	if total > s.c.Article.UpperDraftLimit {
		err = ecode.ArtCreationDraftFull
		return
	}
	id, err = s.dao.AddArtDraft(c, a)
	return
}

// DelArtDraft deletes article draft.
func (s *Service) DelArtDraft(c context.Context, aid, mid int64) (err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	if _, err = s.checkDraftAuthor(c, aid, mid); err != nil {
		return
	}
	err = s.dao.DelArtDraft(c, mid, aid)
	return
}

// UpperDrafts batch get draft by mid.
func (s *Service) UpperDrafts(c context.Context, mid int64, pn, ps int) (res *artmdl.Drafts, err error) {
	if err = s.checkPrivilege(c, mid); err != nil {
		return
	}
	var (
		total int
		ds    []*artmdl.Draft
		start = (pn - 1) * ps
		page  = &artmdl.ArtPage{Pn: pn, Ps: ps}
	)
	res = &artmdl.Drafts{Page: page}
	if total, err = s.dao.CountUpperDraft(c, mid); err != nil {
		return
	} else if total == 0 {
		return
	}
	page.Total = total
	if ds, err = s.dao.UpperDrafts(c, mid, start, ps); err != nil {
		return
	}
	for _, v := range ds {
		// 用户没选择分区，返回空分区信息
		if v.Category.ID == 0 {
			continue
		}
		var pid int64
		if pid, err = s.CategoryToRoot(v.Category.ID); err != nil {
			log.Error("s.CategoryToRoot(%d) error(%+v)", v.Category.ID, err)
			err = nil
			continue
		}
		v.Category = s.categoriesMap[pid]
		v.List, _ = s.dao.List(c, v.ListID)
	}
	res.Drafts = ds
	return
}
