package service

import (
	"context"
	"time"

	"go-common/app/interface/openplatform/article/dao"
	artmdl "go-common/app/interface/openplatform/article/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// UpdateSortCache update sort cache
func (s *Service) UpdateSortCache(c context.Context, aid int64, changed [][2]int64, ip string) (err error) {
	var meta *artmdl.Meta
	if meta, err = s.ArticleMeta(c, aid); (err != nil) || (meta == nil) {
		return
	}
	cids := []int64{meta.Category.ID, meta.Category.ParentID}
	var exist bool
	for _, cid := range cids {
		for _, ch := range changed {
			field := int(ch[0])
			value := ch[1]
			if exist, err = s.dao.ExpireSortCache(c, cid, field); err != nil {
				dao.PromError("sort:更新文章排序缓存")
				return
			}
			if exist && s.shouldAddSort(meta.PublishTime, field) {
				if err = s.dao.AddSortCache(c, cid, field, aid, value); err != nil {
					return
				}
			}
		}
	}
	return
}

func (s *Service) shouldAddSort(t xtime.Time, field int) bool {
	if (field != artmdl.FieldLike) && (field != artmdl.FieldReply) && (field != artmdl.FieldFav) && (field != artmdl.FieldView) {
		return true
	}
	limitTime := time.Now().Unix() - s.sortLimitTime
	return int64(t) > limitTime
}

func (s *Service) delArtSortCache(c context.Context, aid int64) (err error) {
	var (
		root, cid int64
	)
	if root, cid, err = s.RootCategory(c, aid); err != nil {
		dao.PromError("sort:删除文章缓存查找分类")
		log.Error("s.RootCategory(%d,%d) error(%+v)", aid, cid, err)
		return
	}
	cids := []int64{root, cid, _recommendCategory}
	err = s.delArtSortCacheFromCid(c, aid, cids...)
	return
}

func (s *Service) delArtSortCacheFromCid(c context.Context, aid int64, cids ...int64) (err error) {
	for _, cid := range cids {
		for _, field := range artmdl.SortFields {
			if err = s.dao.DelSortCache(c, cid, field, aid); err != nil {
				return
			}
		}
	}
	return
}

func (s *Service) addArtSortCache(c context.Context, art *artmdl.Meta) (err error) {
	if art == nil {
		return
	}
	cids := []int64{art.Category.ID, _recommendCategory}
	var oldRoot int64
	if oldRoot, err = s.CategoryToRoot(art.Category.ID); err != nil {
		dao.PromError("sort:增加文章缓存查找分类")
		log.Error("s.CategoryToRoot(%d,%d) error(%+v)", art.ID, art.Category.ID, err)
		err = nil
	} else {
		cids = append(cids, oldRoot)
	}
	if art.Stats == nil {
		if stat, e := s.stat(c, art.ID); (e == nil) && (stat != nil) {
			art.Stats = stat
		} else {
			art.Stats = new(artmdl.Stats)
		}
	}

	var exist bool
	for _, cid := range cids {
		for _, field := range artmdl.SortFields {
			if exist, err = s.dao.ExpireSortCache(c, cid, field); err != nil {
				dao.PromError("sort:增加最新文章")
				return
			}
			if !exist || !s.shouldAddSort(art.PublishTime, field) {
				continue
			}
			var value int64
			switch field {
			case artmdl.FieldNew:
				value = int64(art.PublishTime)
			case artmdl.FieldLike:
				value = art.Stats.Like
			case artmdl.FieldReply:
				value = art.Stats.Reply
			case artmdl.FieldFav:
				value = art.Stats.Favorite
			case artmdl.FieldView:
				value = art.Stats.View
			default:
				dao.PromError("sort:新增最新文章-排序分类错误")
				log.Error("addArtSortCache sort field error: %v", field)
			}
			if err = s.dao.AddSortCache(c, cid, field, art.ID, value); err != nil {
				dao.PromError("sort:新增最新文章")
			}
		}
	}
	return
}
