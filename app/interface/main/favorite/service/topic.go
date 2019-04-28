package service

import (
	"context"

	"go-common/app/interface/main/favorite/model"
	favmdl "go-common/app/service/main/favorite/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var _emptyTopics = []*model.Topic{}

// AddFavTopic add fav topic
func (s *Service) AddFavTopic(c context.Context, mid, tpID int64, ck, ak string) (err error) {
	tpIDs := []int64{tpID}
	tps, err := s.topicDao.TopicMap(c, tpIDs, false, nil)
	if err != nil {
		log.Error("s.topic.Get(%v)", err)
		return
	}
	if len(tps) == 0 {
		err = ecode.TopicNotExist
		return
	}
	if err = s.AddFavRPC(c, favmdl.TypeTopic, mid, tpID, 0); err != nil {
		log.Error(" s.AddFavRPC(%d,%d) error(%v)", mid, tpID, err)
	}
	return
}

// DelFavTopic del fav topic
func (s *Service) DelFavTopic(c context.Context, mid, tpID int64) (err error) {
	if err = s.DelFavRPC(c, favmdl.TypeTopic, mid, tpID, 0); err != nil {
		log.Error("s.DelFavRPC(%d,%d) error(%v)", mid, tpID, err)
	}
	return
}

// IsTopicFavoured topic is favoured.
func (s *Service) IsTopicFavoured(c context.Context, mid, tpID int64) (faved bool, err error) {
	typ := favmdl.TypeTopic
	if faved, err = s.IsFavRPC(c, typ, mid, tpID); err != nil {
		log.Error("s.IsFavRPC(%d,%d,%d) error(%v)", typ, mid, tpID, err)
	}
	return
}

// FavTopics get fav topics
func (s *Service) FavTopics(c context.Context, mid int64, pn, ps int, appInfo *model.AppInfo) (res *model.TopicList, err error) {
	res = &model.TopicList{}
	res.PageNum = pn
	res.PageSize = ps
	typ := favmdl.TypeTopic
	favs, err := s.FavoritesRPC(c, typ, mid, mid, 0, 0, "", "", pn, ps)
	if err != nil {
		log.Error("s.Favorites(%d,%d,%d,%d,%d,%d,%s) error(%v)", typ, mid, 0, pn, ps, err)
		return
	}
	res.Total = int64(favs.Page.Count)
	var oids []int64
	for _, fav := range favs.List {
		oids = append(oids, fav.Oid)
	}
	if res.Total == 0 {
		res.List = _emptyTopics
		return
	}
	topics, err := s.topicDao.TopicMap(c, oids, false, appInfo)
	if err != nil {
		log.Error("s.topic.MuliGet error(%v)", err)
		return
	}
	for _, fav := range favs.List {
		if topic, ok := topics[fav.Oid]; ok {
			topic.FavAt = fav.MTime
			topic.MID = mid
			res.List = append(res.List, topic)
		}
	}
	return
}
