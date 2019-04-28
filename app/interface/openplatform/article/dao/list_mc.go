package dao

import (
	"context"
	"fmt"

	"go-common/app/interface/openplatform/article/model"
)

func listArtsKey(id int64) string {
	return fmt.Sprintf("art_rl1_arts_%d", id)
}

func listKey(id int64) string {
	return fmt.Sprintf("art_rll_%d", id)
}

func articleListKey(aid int64) string {
	return fmt.Sprintf("art_rlal_%d", aid)
}

func upListsKey(mid int64) string {
	return fmt.Sprintf("art_uplists_%d", mid)
}

func listReadCountKey(id int64) string {
	return fmt.Sprintf("art_lrc_%d", id)
}

func slideArticlesKey(buvid string) string {
	return fmt.Sprintf("art_slidelists_%s", buvid)
}

// ListArtsCacheMap get read list articles cache
func (d *Dao) ListArtsCacheMap(c context.Context, id int64) (res map[int64]*model.ListArtMeta, err error) {
	var arts []*model.ListArtMeta
	if arts, err = d.CacheListArts(c, id); err != nil {
		return
	}
	for _, art := range arts {
		if res == nil {
			res = make(map[int64]*model.ListArtMeta)
		}
		res[art.ID] = art
	}
	return
}

// SetArticleListCache set article list cache
func (d *Dao) SetArticleListCache(c context.Context, listID int64, arts []*model.ListArtMeta) (err error) {
	if len(arts) == 0 {
		return
	}
	m := make(map[int64]int64)
	for _, art := range arts {
		m[art.ID] = listID
	}
	err = d.SetArticlesListCache(c, m)
	return
}
