package service

import (
	"context"
	"go-common/library/log"
	"strings"

	"go-common/app/interface/main/tag/model"
	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
)

// bindTags add tags for article by upper.
func (s *Service) bindTags(c context.Context, mid, aid int64, tags []string, ip string) (err error) {
	arg := &model.ArgBind{Type: model.PicResType, Mid: mid, Oid: aid, Names: tags}
	if err = s.tagRPC.UpBind(c, arg); err != nil {
		dao.PromError("rpc:bind tag")
		log.Error("s.tagRPC.UpBind(%v) error(%+v)", arg, err)
	}
	return
}

// Tags gets article tags.
func (s *Service) Tags(c context.Context, aid int64, skipAct bool) (res []*artmdl.Tag, err error) {
	var (
		tags map[int64][]*model.Tag
		arg  = &model.ArgResTags{Type: model.PicResType, Oids: []int64{aid}}
	)
	if tags, err = s.tagRPC.ResTags(c, arg); err != nil {
		log.Error("s.Tags(%d) error(%+v)", aid, err)
		dao.PromError("rpc:获取Tag")
		return
	} else if tags == nil || len(tags[aid]) == 0 {
		return
	}
	for _, t := range tags[aid] {
		if skipAct && t.Type == 4 {
			continue
		}
		tag := &artmdl.Tag{Tid: t.ID, Name: t.Name}
		res = append(res, tag)
	}
	return
}

// BindTags bind tags with activity
func (s *Service) BindTags(c context.Context, mid, aid int64, tags []string, ip string, activityID int64) (err error) {
	activity := s.activities[activityID]
	if (activityID) > 0 && (activity != nil) && (activity.Tags != "") {
		actTags := strings.Split(activity.Tags, ",")
		tags = mergeActivityTags(tags, actTags)
	}
	return s.bindTags(c, mid, aid, tags, ip)
}

func mergeActivityTags(tags, actTags []string) (res []string) {
	m := map[string]bool{}
	for _, t := range tags {
		m[t] = true
		res = append(res, t)
	}
	for _, t := range actTags {
		if !m[t] {
			res = append(res, t)
		}
	}
	return
}
