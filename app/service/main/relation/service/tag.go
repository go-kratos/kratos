package service

import (
	"context"
	"regexp"
	"sort"
	"time"

	"go-common/app/service/main/relation/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

var (
	_tagRegexp     = regexp.MustCompile(`^[A-Z0-9a-z\x{4e00}-\x{9fa5}]+$`) // only letter digital and chinese legal.
	_keepedTagName = []string{"默认分组", "特别关注", "悄悄关注", "互相关注", "被任命为骑士团"}
	_tagNumLimit   = 15 + 2
	_tagLenLimit   = 16
)

// Tag get user list  by tagid.
func (s *Service) Tag(c context.Context, mid int64, tagid int64, ip string) (tagInfo []int64, err error) {
	if mid <= 0 {
		return
	}
	var (
		fs []*model.Following
	)
	tagInfo = make([]int64, 0)
	alltags, err := s.tags(c, mid)
	if err != nil {
		return
	}
	if _, ok := alltags[tagid]; !ok && tagid != -10 {
		err = ecode.RelTagNotExist
		return
	}
	if fs, err = s.followings(c, mid); err != nil || fs == nil {
		return
	}
	var list []*model.Following
	for _, f := range fs {
		if !f.Following() {
			continue
		}
		if tagid == 0 {
			if len(f.Tag) == 0 {
				list = append(list, f)
				continue
			}
			var exist = false
			for _, t := range f.Tag {
				if _, ok := alltags[t]; ok {
					exist = true
					break
				}
			}
			if !exist {
				list = append(list, f)
			}
			continue
		}
		for _, ft := range f.Tag {
			if tagid == ft {
				list = append(list, f)
			}
		}
	}
	sort.Slice(list, func(i, j int) bool { return list[i].MTime > list[j].MTime })
	for _, f := range list {
		tagInfo = append(tagInfo, f.Mid)
	}
	return
}

// Tags get tag list.
func (s *Service) Tags(c context.Context, mid int64, ip string) (tags []*model.TagCount, err error) {
	if mid <= 0 {
		return
	}
	if tags, err = s.dao.TagCountCache(c, mid); err != nil {
		return
	} else if tags != nil {
		return
	}
	var (
		fs []*model.Following
	)
	if fs, err = s.followings(c, mid); err != nil || fs == nil {
		return
	}
	alltags, err := s.tags(c, mid)
	if err != nil {
		return
	}
	tc := make(map[int64]int64, len(alltags))
	// init tag count
	for tid := range alltags {
		tc[tid] = 0
	}
	for _, f := range fs {
		var deleted = true
		if !f.Following() {
			continue
		}
		if len(f.Tag) == 0 {
			tc[0]++
		} else {
			for _, v := range f.Tag {
				if _, ok := tc[v]; ok {
					tc[v]++
					deleted = false
				}
			}
			if deleted {
				tc[0]++
			}
		}
	}
	for k, v := range tc {
		tmp := &model.TagCount{Tagid: k, Name: alltags[k].Name, Count: v}
		tags = append(tags, tmp)
	}
	sort.Slice(tags, func(i, j int) bool { return tags[i].Tagid < tags[j].Tagid })
	s.addCache(func() {
		s.dao.SetTagCountCache(context.Background(), mid, tags)
	})
	return
}

// UserTag get user tags.
func (s *Service) UserTag(c context.Context, mid int64, fid int64, ip string) (tags map[int64]string, err error) {
	if mid <= 0 || fid <= 0 {
		return
	}
	if mid == fid {
		return
	}
	var (
		mpf map[int64]*model.Following
		tag *model.TagUser
	)
	tags = make(map[int64]string)
	if mpf, err = s.dao.RelationsCache(c, mid, []int64{fid}); err != nil {
		return
	} else if mpf != nil {
		if tag, ok := mpf[fid]; ok {
			return s.tagidToName(c, mid, tag.Tag)
		}
	}
	if tag, err = s.dao.TagUserByMidFid(c, mid, fid); err != nil {
		return
	} else if tag == nil || len(tag.Tag) == 0 {
		return
	}
	return s.tagidToName(c, mid, tag.Tag)
}

func (s *Service) tagidToName(c context.Context, mid int64, tagids []int64) (ttn map[int64]string, err error) {
	if len(tagids) == 0 {
		return
	}
	alltags, err := s.tags(c, mid)
	if err != nil {
		return
	}
	ttn = make(map[int64]string)
	for _, id := range tagids {
		if tag, ok := alltags[id]; ok {
			ttn[tag.Id] = tag.Name
		}
	}
	return
}

// CreateTag add tag.
func (s *Service) CreateTag(c context.Context, mid int64, tagStr string, ip string) (res int64, err error) {
	if mid <= 0 {
		return
	}
	for _, v := range _keepedTagName {
		if tagStr == v {
			err = ecode.RelTagExisted
			return
		}
	}
	if len([]rune(tagStr)) > _tagLenLimit {
		err = ecode.RelTagLenLimit
		return
	}
	if !s.tagCheck(tagStr) {
		err = ecode.RelTagExistNotAllowedWords
		return
	}
	tags, err := s.dao.Tags(c, mid)
	if err != nil {
		return
	}
	if len(tags) >= _tagNumLimit {
		err = ecode.RelTagNumLimit
		return
	}
	for _, tag := range tags {
		if tag.Name == tagStr {
			err = ecode.RelTagExisted
			return
		}
	}
	res, err = s.dao.AddTag(c, mid, mid, tagStr, time.Now())
	s.addCache(func() {
		s.dao.DelTagCountCache(context.Background(), mid)
		s.dao.DelTagsCache(context.Background(), mid)
	})
	return
}

// UpdateTag update tag name.
func (s *Service) UpdateTag(c context.Context, mid int64, tagID int64, newTag string, ip string) (err error) {
	if mid <= 0 {
		return
	}
	for _, v := range _keepedTagName {
		if newTag == v {
			err = ecode.RelTagExisted
			return
		}
	}
	if len([]rune(newTag)) > _tagLenLimit {
		err = ecode.RelTagLenLimit
		return
	}
	if !s.tagCheck(newTag) {
		err = ecode.RelTagExistNotAllowedWords
		return
	}
	alltags, err := s.tags(c, mid)
	if err != nil {
		return
	}
	if _, ok := alltags[tagID]; !ok {
		err = ecode.RelTagNotExist
		return
	}
	for _, tag := range alltags {
		if tag.Name == newTag {
			err = ecode.RelTagExisted
			return
		}
	}
	if _, err = s.dao.SetTagName(c, tagID, mid, newTag, time.Now()); err != nil {
		return
	}
	s.addCache(func() {
		s.dao.DelTagCountCache(context.Background(), mid)
		s.dao.DelTagsCache(context.Background(), mid)
	})
	return
}

// DelTag del user tg by tagid.
func (s *Service) DelTag(c context.Context, mid int64, tagID int64, ip string) (err error) {
	if mid <= 0 {
		return
	}
	if _, err = s.dao.DelTag(c, mid, tagID); err != nil {
		return
	}
	s.addCache(func() {
		s.dao.DelTagCountCache(context.Background(), mid)
		s.dao.DelTagsCache(context.Background(), mid)
	})
	return
}

// TagsMoveUsers move user to new tags from beforeid.
// if beforeid equal zero,just copy
func (s *Service) TagsMoveUsers(c context.Context, mid, beforeid int64, afterIdsStr, fidStr string, ip string) (err error) {
	if mid <= 0 {
		return
	}
	var (
		fids  []int64
		tids  []int64
		rms   []*model.Following
		mtags map[int64]*model.Tag
	)
	if tids, err = xstr.SplitInts(afterIdsStr); err != nil || len(tids) > _tagNumLimit {
		err = ecode.RequestErr
		return
	}
	if fids, err = xstr.SplitInts(fidStr); err != nil {
		err = ecode.RequestErr
		return
	}
	// 判断是否已经关注
	if rms, err = s.dao.FollowingsIn(c, mid, fids); err != nil {
		return
	}
	if len(rms) != len(fids) {
		err = ecode.RelTagAddFollowingFirst
		return
	}
	for _, v := range rms {
		if !v.Following() {
			err = ecode.RelTagAddFollowingFirst
			return
		}
	}
	if mtags, err = s.dao.Tags(c, mid); err != nil {
		return
	}
	tmpTids := make([]int64, 0)
	for _, tid := range tids {
		if tid == 0 {
			continue
		}
		if tag, ok := mtags[tid]; !ok || tag.Status != 0 {
			err = ecode.RelTagNotExist
			return
		}
		tmpTids = append(tmpTids, tid)
	}
	utags, err := s.dao.UsersTags(c, mid, fids)
	if err != nil {
		return
	}
	for _, fid := range fids {
		var atags []int64
		if tag, ok := utags[fid]; ok {
			btags := make(map[int64]struct{})
			for _, tid := range tag.Tag {
				if mtag, ok := mtags[tid]; ok && mtag.Status == 0 && tid != beforeid {
					btags[mtag.Id] = struct{}{}
				}
			}
			for _, t := range tmpTids {
				btags[t] = struct{}{}
			}
			for tid := range btags {
				atags = append(atags, tid)
			}
		} else {
			atags = tmpTids
		}
		// TODO:add all user in once.
		_, err = s.dao.AddTagUser(c, mid, fid, atags, time.Now())
	}
	return
}

// TagsAddUsers add user to tidStr.
func (s *Service) TagsAddUsers(c context.Context, mid int64, tidStr, fidStr string, ip string) (err error) {
	if mid <= 0 {
		return
	}
	var (
		fids  []int64
		tids  []int64
		rms   []*model.Following
		mtags map[int64]*model.Tag
	)
	if tids, err = xstr.SplitInts(tidStr); err != nil || len(tids) > _tagNumLimit {
		err = ecode.RequestErr
		return
	}
	if fids, err = xstr.SplitInts(fidStr); err != nil {
		err = ecode.RequestErr
		return
	}
	// 判断是否已经关注
	if rms, err = s.dao.FollowingsIn(c, mid, fids); err != nil {
		return
	}
	if len(rms) != len(fids) {
		err = ecode.RelTagAddFollowingFirst
		return
	}
	for _, v := range rms {
		if !v.Following() {
			err = ecode.RelTagAddFollowingFirst
			return
		}
	}
	if mtags, err = s.dao.Tags(c, mid); err != nil {
		return
	}
	tmpTids := make([]int64, 0)
	for _, tid := range tids {
		if tid == 0 {
			continue
		}
		if tag, ok := mtags[tid]; !ok || tag.Status != 0 {
			log.Warn("Invalid tag id: %d: tag not exist", tid)
			continue
			// err = ecode.RelTagNotExist
			// return
		}
		tmpTids = append(tmpTids, tid)
	}
	for _, fid := range fids {
		_, err = s.dao.AddTagUser(c, mid, fid, tmpTids, time.Now())
	}
	return
}

func (s *Service) tagCheck(tag string) bool {
	return _tagRegexp.MatchString(tag)
}

func (s *Service) tags(c context.Context, mid int64) (alltags map[int64]*model.Tag, err error) {
	alltags, err = s.dao.TagsCache(c, mid)
	if err != nil {
		return
	}
	// cache miss.
	if len(alltags) == 0 {
		alltags, err = s.dao.Tags(c, mid)
		if err != nil {
			return
		}
		s.addCache(func() {
			s.dao.SetTagsCache(context.Background(), mid, &model.Tags{Tags: alltags})
		})
	}
	return
}

// DelTagCache all tag related cache.
func (s *Service) DelTagCache(c context.Context, mid int64) (err error) {
	if err = s.DelFollowingCache(c, mid); err != nil {
		return
	}
	if err = s.dao.DelTagCountCache(c, mid); err != nil {
		return
	}
	if err = s.dao.DelTagsCache(c, mid); err != nil {
		return
	}
	return s.dao.DelStatCache(c, mid)
}

// AddSpecial add fid to special list.
func (s *Service) AddSpecial(c context.Context, mid, fid int64) (err error) {
	rl, err := s.Relation(c, mid, fid)
	if err != nil {
		return
	}
	if rl == nil || !rl.Following() {
		err = ecode.RelTagAddFollowingFirst
		return
	}
	var sp bool
	for _, id := range rl.Tag {
		if id == -10 {
			sp = true
		}
	}
	ids := rl.Tag
	if !sp {
		ids = append(ids, -10)
		s.dao.AddTagUser(c, mid, fid, ids, time.Now())
	} else {
		err = ecode.RelFollowAttrAlreadySet
	}
	return
}

// DelSpecial del fid from special list.
func (s *Service) DelSpecial(c context.Context, mid, fid int64) (err error) {
	rl, err := s.Relation(c, mid, fid)
	if err != nil {
		return
	}
	if rl == nil || !rl.Following() {
		err = ecode.RelTagAddFollowingFirst
		return
	}
	var ids []int64
	for _, id := range rl.Tag {
		if id != -10 {
			ids = append(ids, id)
		}
	}
	// no special before.
	if len(ids) == len(rl.Tag) {
		err = ecode.RelFollowAttrNotSet
		return
	}
	s.dao.AddTagUser(c, mid, fid, ids, time.Now())
	return
}

// Special get special list.
func (s *Service) Special(c context.Context, mid int64) (list []int64, err error) {
	return s.Tag(c, mid, -10, "")
}
