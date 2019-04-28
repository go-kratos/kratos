package service

import (
	"context"
	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
)

var (
	_emptySynonymTag = make([]*model.SynonymTag, 0)
	_emptyBasicTag   = make([]*model.BasicTag, 0)
)

// SynonymList SynonymList.
func (s *Service) SynonymList(c context.Context, keyWord string, pn, ps int32) (res []*model.SynonymTag, count int64, err error) {
	var (
		ids, cids []int64
		sTagMap   map[int64]*model.SynonymTag
	)
	start := (pn - 1) * ps
	end := ps
	tagMap := make(map[int64]*model.Tag)
	stagChild := make([]*model.SynonymTag, 0)
	if count, _ = s.dao.SynonymCount(c, keyWord); count <= 0 {
		return
	}
	if sTagMap, res, ids, err = s.dao.Synonyms(c, keyWord, start, end); err != nil {
		return
	}
	if len(ids) > 0 {
		_, stagChild, cids, _ = s.dao.SynonymIDs(c, ids)
	}
	if len(cids) > 0 {
		_, tagMap, _ = s.dao.Tags(c, cids)
	}
	for _, v := range stagChild {
		if p, ok := sTagMap[v.Ptid]; ok {
			if k, b := tagMap[v.Tid]; b {
				p.Adverb = append(p.Adverb, &model.BasicTag{ID: v.Tid, Name: k.Name})
			}
		}
	}
	for _, v := range res {
		if t, ok := sTagMap[v.Tid]; ok {
			v.Adverb = t.Adverb
		} else {
			v.Adverb = _emptyBasicTag
		}
	}
	if len(res) == 0 {
		res = _emptySynonymTag
	}
	return
}

// SynonymAdd SynonymAdd.
func (s *Service) SynonymAdd(c context.Context, uname, tname string, adverb []int64) (err error) {
	var (
		ptid, affect int64
		tag          = new(model.Tag)
		stag         = new(model.SynonymTag)
	)
	if tag, err = s.dao.TagByName(c, tname); err != nil {
		return
	}
	if tag == nil {
		return ecode.TagNotExist
	}
	if stag, err = s.dao.SynonymByName(c, tname); err != nil {
		return
	}
	if stag == nil {
		if affect, err = s.dao.InsertSynonym(c, uname, 0, tag.ID); err != nil {
			return
		}
		if affect <= 0 {
			return ecode.TagOperateFail
		}
		ptid = tag.ID
	} else {
		ptid = stag.Tid
	}
	adverb = s.removeDuplicate(adverb)
	if len(adverb) == 0 {
		return
	}
	if _, err = s.dao.DelSynonymSon(c, ptid, adverb); err != nil {
		return
	}
	if affect, err = s.dao.InsertSynonyms(c, uname, ptid, adverb); err != nil {
		return
	}
	if affect <= 0 {
		return ecode.TagOperateFail
	}
	return
}

// SynonymInfo Synonym info by ID.
func (s *Service) SynonymInfo(c context.Context, tid int64) (res *model.SynonymInfo, err error) {
	var (
		ids  []int64
		stag []*model.SynonymTag
		tags []*model.Tag
	)
	res = new(model.SynonymInfo)
	if stag, err = s.dao.Synonym(c, tid); err != nil {
		return
	}
	for _, st := range stag {
		if st.Ptid != 0 {
			ids = append(ids, st.Ptid)
			ids = append(ids, st.Tid)
		}
	}
	if len(ids) == 0 {
		return
	}
	if tags, _, err = s.dao.Tags(c, ids); err != nil {
		return
	}
	for _, t := range tags {
		if t.ID == tid {
			res.Tag = t
		} else {
			res.Adverb = append(res.Adverb, t)
		}
	}
	return
}

// SynonymDelete del synonym.
func (s *Service) SynonymDelete(c context.Context, tid int64) (err error) {
	var (
		ids, tids []int64
	)
	if _, err = s.dao.DelSynonym(c, tid); err != nil {
		return
	}
	tids = append(tids, tid)
	if len(tids) > 0 {
		if _, _, ids, err = s.dao.SynonymIDs(c, tids); err != nil {
			return
		}
	}
	if len(ids) != 0 {
		if err = s.RemoveSynonymSon(c, tid, ids); err != nil {
			return
		}
	}
	return
}

// RemoveSynonymSon 删除二级子类.
func (s *Service) RemoveSynonymSon(c context.Context, ptid int64, tids []int64) (err error) {
	_, err = s.dao.DelSynonymSon(c, ptid, tids)
	return
}

// SynonymIsExist SynonymIsExist.
func (s *Service) SynonymIsExist(c context.Context, adverb string) (tag *model.Tag, err error) {
	if tag, err = s.dao.TagByName(c, adverb); err != nil {
		return
	}
	if tag == nil {
		return nil, ecode.TagNotExist
	}
	return
}

func (s *Service) removeDuplicate(strs []int64) (res []int64) {
	tempMap := map[int64]byte{} // 存放不重复主键
	for _, str := range strs {
		l := len(tempMap)
		tempMap[str] = 0
		if len(tempMap) != l {
			res = append(res, str)
		}
	}
	return
}
