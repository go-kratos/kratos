package service

import (
	"context"
	"time"

	"go-common/app/interface/main/tag/model"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
)

// InfoByID .
func (s *Service) InfoByID(c context.Context, mid, tid int64) (t *model.Tag, err error) {
	return s.info(c, mid, tid)
}

// infoByID .
func (s *Service) info(c context.Context, mid, tid int64) (t *model.Tag, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	tag, err := s.tag(c, tid, mid)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.tag() tid:%d,mid:%d,ip:%s,err:%v", tid, mid, ip, err)
		return
	}
	if tag == nil {
		err = ecode.TagNotExist
		return
	}
	t = &model.Tag{
		ID:           tag.ID,
		Type:         int8(tag.Type),
		Name:         tag.Name,
		Cover:        tag.Cover,
		Content:      tag.Content,
		Attribute:    int8(tag.Attr),
		State:        int8(tag.State),
		IsAtten:      int8(tag.Attention),
		HeadCover:    tag.HeadCover,
		ShortContent: tag.ShortContent,
		CTime:        tag.CTime,
		MTime:        tag.MTime,
	}
	r, _ := s.count(c, tid, mid)
	if r != nil {
		t.Count.Atten = int(r.Sub)
		t.Count.Use = int(r.Bind)
	}
	return t, nil
}

// MinfoByIDs .
func (s *Service) MinfoByIDs(c context.Context, mid int64, tids []int64) (ts []*model.Tag, err error) {
	return s.infos(c, mid, tids)
}

// MinfoByIDs .
func (s *Service) infos(c context.Context, mid int64, tids []int64) (ts []*model.Tag, err error) {
	if len(tids) == 0 {
		err = ecode.TagNotExist
		return
	}
	// TODO
	if tids[0] == -1 {
		err = ecode.TagNotExist
		return
	}
	tagMap, err := s.dao.TagMap(c, tids, mid)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.infos() tids(%+v) mid(%d) error(%v)", tids, mid, err)
		return
	}
	for _, tag := range tagMap {
		t := &model.Tag{
			ID:           tag.Id,
			Name:         tag.Name,
			Cover:        tag.Cover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			HeadCover:    tag.HeadCover,
			Type:         int8(tag.Type),
			State:        int8(tag.State),
			CTime:        tag.Ctime,
			MTime:        tag.Mtime,
			IsAtten:      int8(tag.Attention),
			Attribute:    int8(tag.Attr),
		}
		t.Count.Atten = int(tag.Sub)
		t.Count.Use = int(tag.Bind)
		ts = append(ts, t)
	}
	return
}

// InfoByName .
func (s *Service) InfoByName(c context.Context, mid int64, name string) (t *model.Tag, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	tag, err := s.tagName(c, mid, name)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.tag() name:%s,mid:%d,ip:%s,err:%v", name, mid, ip, err)
		return
	}
	if tag == nil {
		err = ecode.TagNotExist
		return
	}
	t = &model.Tag{
		ID:           tag.ID,
		Type:         int8(tag.Type),
		Name:         tag.Name,
		Cover:        tag.Cover,
		Content:      tag.Content,
		Attribute:    int8(tag.Attr),
		State:        int8(tag.State),
		IsAtten:      int8(tag.Attention),
		HeadCover:    tag.HeadCover,
		ShortContent: tag.ShortContent,
		CTime:        tag.CTime,
		MTime:        tag.MTime,
	}
	r, _ := s.count(c, tag.ID, mid)
	if r != nil {
		t.Count.Atten = int(r.Sub)
		t.Count.Use = int(r.Bind)
	}
	return t, nil
}

// MinfoByNames .
func (s *Service) MinfoByNames(c context.Context, mid int64, names []string) (ts []*model.Tag, err error) {
	var (
		tids []int64
		t    *model.Tag
		ip   = metadata.String(c, metadata.RemoteIP)
	)
	tags, err := s.tagNames(c, mid, names)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.tag() names:%+v,mid:%d,ip:%s,err:%v", names, mid, ip, err)
		return
	}
	for _, v := range tags {
		tids = append(tids, v.ID)
	}
	cm, _ := s.countMap(c, tids, mid)
	for _, v := range tags {
		t = &model.Tag{
			ID:           v.ID,
			Type:         int8(v.Type),
			Name:         v.Name,
			Cover:        v.Cover,
			Content:      v.Content,
			Attribute:    int8(v.Attr),
			State:        int8(v.State),
			IsAtten:      int8(v.Attention),
			HeadCover:    v.HeadCover,
			ShortContent: v.ShortContent,
			CTime:        v.CTime,
			MTime:        v.MTime,
		}
		if len(cm) > 0 {
			if r, ok := cm[v.ID]; ok {
				t.Count.Atten = int(r.Sub)
				t.Count.Use = int(r.Bind)
			}
		}
		ts = append(ts, t)
	}
	return
}

// CheckTag .
func (s *Service) CheckTag(c context.Context, mid int64, name string, now time.Time) (t *model.Tag, err error) {
	if err = s.filter(c, name, now); err != nil {
		return
	}
	t, err = s.resCheckTag(c, mid, name, model.UpTag, now)
	return
}

// HotTags .
func (s *Service) HotTags(c context.Context, mid, rid, hotType int64) (mhs []*model.HotTags, err error) {
	var (
		hs   *model.HotTags
		hots []*model.HotTag
	)
	rids, ok := s.ridMap[rid]
	if !ok {
		rids = append(rids, rid)
	}
	for _, r := range rids {
		key := s.hotTagKey(r, hotType)
		s.hLock.RLock()
		hs, ok = s.hotTag[key]
		s.hLock.RUnlock()
		if !ok {
			if hots, err = s.hots(c, r, hotType); err != nil {
				return
			}
			if hots == nil {
				err = ecode.TagNotExist
				return
			}
			hs = &model.HotTags{Rid: r, Tags: hots}
			s.hLock.Lock()
			s.hotTag[key] = hs
			s.hLock.Unlock()
		}
		mhs = append(mhs, hs)
	}
	if mid > 0 {
		var attenMap map[int64]*model.Tag
		if attenMap, err = s.attened(c, mid); err != nil {
			return
		}
		nmhs := make([]*model.HotTags, 0, len(mhs))
		for _, mh := range mhs {
			nhs := &model.HotTags{Rid: mh.Rid}
			for _, t := range mh.Tags {
				nht := &model.HotTag{}
				*nht = *t // NOTE: make sure copy, is_atten can be change by every mid
				if _, ok := attenMap[t.Tid]; ok {
					nht.IsAtten = 1
				}
				nhs.Tags = append(nhs.Tags, nht)
			}
			nmhs = append(nmhs, nhs)
		}
		mhs = nmhs
	}
	return
}

// SimilarTags .
func (s *Service) SimilarTags(c context.Context, rid, tid int64) (sts []*model.SimilarTag, err error) {
	var (
		ok      bool
		key     string
		tname   string
		stsTids []int64
		tag     *model.Tag
		tags    []*model.Tag
		tagMap  map[int64]string
		st      *model.SimilarTag
	)
	if tag, err = s.info(c, 0, tid); err != nil {
		return
	}
	if tag == nil {
		return nil, ecode.TagNotExist
	}
	if tag.State == model.TagStateDel || tag.State == model.TagStateHide {
		return nil, ecode.TagIsSealing
	}
	key = s.similarKey(rid, tid)
	s.sLock.RLock()
	sts, ok = s.simlars[key]
	s.sLock.RUnlock()
	if ok {
		return
	}
	if sts, err = s.similars(c, rid, tid); err != nil {
		log.Error("s.tag.Similars(%d,%d,%s) error(%v)", rid, tid, metadata.String(c, metadata.RemoteIP), err)
		err = nil
		return
	}
	for _, st = range sts {
		stsTids = append(stsTids, st.Tid)
	}
	if tags, err = s.infos(c, 0, stsTids); err != nil {
		return
	}
	if len(tags) == 0 {
		err = ecode.TagRankingSimNotExist
		return
	}
	tagMap = make(map[int64]string)
	for _, tag = range tags {
		tagMap[tag.ID] = tag.Name
	}
	for _, st = range sts {
		if tname, ok = tagMap[st.Tid]; ok {
			st.Tname = tname
		}
	}
	if len(sts) > 0 {
		s.sLock.Lock()
		s.simlars[key] = sts
		s.sLock.Unlock()
	}
	return
}

// ChangeSim .
func (s *Service) ChangeSim(c context.Context, tid int64) (sis []*model.SimilarTag, err error) {
	var (
		tids []int64
		tag  *model.Tag
		tags []*model.Tag
	)
	if tag, err = s.info(c, 0, tid); err != nil {
		return
	}
	if tag == nil {
		return nil, ecode.TagNotExist
	}
	if tag.State == model.TagStateDel || tag.State == model.TagStateHide {
		return nil, ecode.TagIsSealing
	}
	tids, _ = s.similarsTids(c, tid)
	if len(tids) > 0 {
		if tags, err = s.infos(c, 0, tids); err != nil {
			return
		}
		for _, t := range tags {
			sis = append(sis, &model.SimilarTag{
				Tid:    t.ID,
				Tname:  t.Name,
				TCover: t.Cover,
				Tatten: t.Count.Atten,
			})
		}
	}
	if len(sis) == 0 {
		sis = []*model.SimilarTag{}
	}
	return
}

// AddActivityTag check activity tag if new add it
func (s *Service) AddActivityTag(c context.Context, name string, now time.Time) (err error) {
	t, _ := s.InfoByName(c, 0, name)
	if t != nil {
		if t.Type != int8(rpcModel.TypeOfficailActivity) {
			err = ecode.TagAlreadyExist
		}
		return
	}
	rt := &rpcModel.Tag{
		Name: name,
		Type: rpcModel.TypeOfficailActivity,
	}
	if err = s.createTag(c, rt); err != nil {
		return
	}
	_, err = s.InfoByName(c, 0, name)
	return
}

// RecommandTag .
func (s *Service) RecommandTag(c context.Context) (res map[int64]map[string][]*rpcModel.UploadTag, err error) {
	return s.recommandTagService(c)
}

// TagGroup .
func (s *Service) TagGroup(c context.Context) (rs []*rpcModel.Synonym, err error) {
	return s.tagGroup(c)
}

// HotMap .
func (s *Service) HotMap(c context.Context) (res map[int16][]int64, err error) {
	return s.hotMap(c)
}

// Prids .
func (s *Service) Prids(c context.Context) (res []int64, err error) {
	return s.prids(c)
}

// getChangeTag .
func (s *Service) getChangeTag(sourceTags, convertTags []string, regionTags []string) (allTags, addTags, delTags, defaultNames []string) {
	sourceTagMap := make(map[string]struct{}, len(sourceTags))
	convertTagMap := make(map[string]struct{}, len(convertTags))
	regionTagMap := make(map[string]struct{}, len(regionTags))
	for _, tag := range sourceTags {
		sourceTagMap[tag] = struct{}{}
	}
	for _, tag := range convertTags {
		convertTagMap[tag] = struct{}{}
	}
	for _, tag := range regionTags {
		regionTagMap[tag] = struct{}{}
	}
	// add
	for tag := range convertTagMap {
		if _, ok := sourceTagMap[tag]; !ok {
			addTags = append(addTags, tag)
		}
	}
	// del
	for tag := range sourceTagMap {
		if _, ok := convertTagMap[tag]; !ok {
			delTags = append(delTags, tag)
		}
	}
	// default
	for _, tag := range regionTags {
		if _, ok := convertTagMap[tag]; !ok {
			defaultNames = append(defaultNames, tag)
		}
	}
	allTags = append(allTags, sourceTags...)
	allTags = append(allTags, addTags...)
	allTags = append(allTags, defaultNames...)
	return
}

// TopicTags TopicTags.
func (s *Service) TopicTags(c context.Context, mid int64, names []string, now time.Time) (res []*model.Tag, err error) {
	res, _, err = s.addNewTags(c, mid, names, now)
	return
}

// addNewTags TopicTags.
func (s *Service) addNewTags(c context.Context, mid int64, names []string, now time.Time) (res []*model.Tag, tMap map[string]int64, err error) {
	names, err = s.mfilter(c, names, now)
	if err != nil {
		return
	}
	if len(names) == 0 {
		return _emptyTs, nil, nil
	}
	if res, err = s.MinfoByNames(c, mid, names); err != nil {
		if !ecode.TagNotExist.Equal(err) {
			log.Error("s.MinfoByNames(%v) error(%v)", names, err)
			return nil, nil, err
		}
		err = nil
	}
	tMap = make(map[string]int64, len(names))
	tnameMap := make(map[string]string)
	for _, name := range names {
		tnameMap[name] = name
	}
	for _, et := range res {
		delete(tnameMap, et.Name)
		tMap[et.Name] = et.ID
	}
	if len(tnameMap) != 0 {
		var (
			ts []*rpcModel.Tag
			tn []string
		)
		for _, v := range tnameMap {
			tn = append(tn, v)
			ts = append(ts, &rpcModel.Tag{
				Name: v,
			})
		}
		if err = s.createTags(c, ts); err != nil {
			return
		}
		rtags, _ := s.MinfoByNames(c, mid, tn)
		for _, rt := range rtags {
			res = append(res, rt)
			tMap[rt.Name] = rt.ID
		}
	}
	if len(res) == 0 {
		res = _emptyTs
	}
	return
}

// resCheckTag .
func (s *Service) resCheckTag(c context.Context, mid int64, name string, tp int8, now time.Time) (t *model.Tag, err error) {
	if t, err = s.InfoByName(c, 0, name); err != nil {
		if !ecode.TagNotExist.Equal(err) {
			return
		}
		err = nil
	}
	if t == nil {
		if err = s.createTag(c, &rpcModel.Tag{Name: name}); err != nil {
			return
		}
		if t, err = s.InfoByName(c, 0, name); err != nil {
			return
		}
	}
	// 不正常tag不可加
	if t.State != model.TagStateNormal {
		err = ecode.TagIsSealing
		return
	}
	// 活动tag不可加
	if t.Type == model.OfficailActiveTag {
		err = ecode.TagIsOfficailTag
	}
	return
}

// TagTop web-interface tag top, include tag info, and similar tags.
func (s *Service) TagTop(c context.Context, arg *model.ReqTagTop) (res *model.TagTop, err error) {
	var tag *model.Tag
	if arg.Tid > 0 {
		tag, err = s.info(c, arg.Mid, arg.Tid)
	} else {
		tag, err = s.InfoByName(c, arg.Mid, arg.TName)
	}
	if err != nil {
		return
	}
	if tag == nil {
		return nil, ecode.TagNotExist
	}
	if tag.State == model.TagStateDel || tag.State == model.TagStateHide {
		return nil, ecode.TagIsSealing
	}
	res = &model.TagTop{
		Tag:      tag,
		Similars: make([]*model.SimilarTag, 0),
	}
	var (
		tids []int64
		tags []*model.Tag
	)
	tids, _ = s.similarsTids(c, tag.ID)
	if len(tids) > 0 {
		tags, _ = s.infos(c, model.NoneUserID, tids)
		for _, tag := range tags {
			similarTag := &model.SimilarTag{
				Tid:   tag.ID,
				Tname: tag.Name,
			}
			res.Similars = append(res.Similars, similarTag)
		}
	}
	return
}
