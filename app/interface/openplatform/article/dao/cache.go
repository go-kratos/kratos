package dao

import (
	"context"
	"strconv"

	"go-common/app/interface/openplatform/article/model"
)

func (d *Dao) cacheSFList(id int64) string {
	return strconv.FormatInt(id, 10)
}

func (d *Dao) cacheSFListArts(id int64) string {
	return strconv.FormatInt(id, 10)
}

func (d *Dao) cacheSFUpLists(id int64) string {
	return strconv.FormatInt(id, 10)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -nullcache=&model.List{ID:-1} -check_null_code=$!=nil&&$.ID==-1 -singleflight=true
	List(c context.Context, id int64) (*model.List, error)
	// cache: -batch=100 -max_group=10 -nullcache=&model.List{ID:-1} -check_null_code=$!=nil&&$.ID==-1
	Lists(c context.Context, keys []int64) (map[int64]*model.List, error)
	// cache: -singleflight=true -nullcache=[]*model.ListArtMeta{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	ListArts(c context.Context, id int64) ([]*model.ListArtMeta, error)
	// cache: -nullcache=[]*model.ListArtMeta{{ID:-1}} -check_null_code=len($)==1&&$[0].ID==-1
	ListsArts(c context.Context, ids []int64) (map[int64][]*model.ListArtMeta, error)
	// cache: -nullcache=-1 -batch=100 -max_group=10
	ArtsListID(c context.Context, keys []int64) (map[int64]int64, error)
	// cache: -nullcache=[]int64{-1} -check_null_code=len($)==1&&$[0]==-1 -singleflight=true
	UpLists(c context.Context, mid int64) ([]int64, error)
	// cache: -nullcache=&model.AuthorLimit{Limit:-1} -check_null_code=$!=nil&&$.Limit==-1
	Author(c context.Context, mid int64) (*model.AuthorLimit, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	// 获取文集文章列表缓存
	//mc: -key=listArtsKey
	CacheListArts(c context.Context, id int64) (res []*model.ListArtMeta, err error)
	// 增加文集含有的文章列表缓存
	//mc: -key=listArtsKey -expire=d.mcListArtsExpire
	AddCacheListArts(c context.Context, id int64, arts []*model.ListArtMeta) (err error)
	// 获取文章所属文集
	//mc: -key=articleListKey -type=get
	ArticleListCache(c context.Context, id int64) (res int64, err error)
	// 增加文章所属文集缓存
	//mc: -key=articleListKey -expire=d.mcArtListExpire
	SetArticlesListCache(c context.Context, arts map[int64]int64) (err error)
	//mc: -key=listKey
	CacheList(c context.Context, id int64) (res *model.List, err error)
	//mc: -key=listKey -expire=d.mcListExpire
	AddCacheList(c context.Context, id int64, list *model.List) (err error)
	//mc: -key=listKey
	CacheLists(c context.Context, ids []int64) (res map[int64]*model.List, err error)
	//mc: -key=listKey -expire=d.mcListExpire
	AddCacheLists(c context.Context, lists map[int64]*model.List) (err error)
	//mc: -key=listArtsKey
	CacheListsArts(c context.Context, ids []int64) (res map[int64][]*model.ListArtMeta, err error)
	//mc: -key=listArtsKey -expire=d.mcListArtsExpire
	AddCacheListsArts(c context.Context, arts map[int64][]*model.ListArtMeta) (err error)
	//mc: -key=articleListKey
	CacheArtsListID(c context.Context, ids []int64) (res map[int64]int64, err error)
	//mc: -key=articleListKey -expire=d.mcArtListExpire
	AddCacheArtsListID(c context.Context, arts map[int64]int64) (err error)
	//mc: -key=upListsKey -expire=d.mcUpListsExpire
	AddCacheUpLists(c context.Context, mid int64, lists []int64) (err error)
	//mc: -key=upListsKey
	CacheUpLists(c context.Context, id int64) (res []int64, err error)
	//mc: -key=listReadCountKey -expire=d.mcListReadExpire
	AddCacheListReadCount(c context.Context, id int64, read int64) (err error)
	//mc: -key=listReadCountKey
	CacheListReadCount(c context.Context, id int64) (res int64, err error)
	//mc: -key=listReadCountKey
	CacheListsReadCount(c context.Context, ids []int64) (res map[int64]int64, err error)
	//mc: -key=hotspotsKey -expire=d.mcHotspotExpire
	AddCacheHotspots(c context.Context, hots []*model.Hotspot) (err error)
	//mc: -key=hotspotsKey
	DelCacheHotspots(c context.Context) (err error)
	//mc: -key=hotspotsKey
	cacheHotspots(c context.Context) (res []*model.Hotspot, err error)
	//mc: -key=mcHotspotKey
	CacheHotspot(c context.Context, id int64) (res *model.Hotspot, err error)
	//mc: -key=mcHotspotKey -expire=d.mcHotspotExpire
	AddCacheHotspot(c context.Context, id int64, val *model.Hotspot) (err error)
	// 增加作者状态缓存
	//mc: -key=mcAuthorKey -expire=d.mcAuthorExpire
	AddCacheAuthor(c context.Context, mid int64, author *model.AuthorLimit) (err error)
	//mc: -key=mcAuthorKey
	CacheAuthor(c context.Context, mid int64) (res *model.AuthorLimit, err error)
	//mc: -key=mcAuthorKey
	DelCacheAuthor(c context.Context, mid int64) (err error)
	//mc: -key=slideArticlesKey
	CacheListArtsId(c context.Context, buvid string) (*model.ArticleViewList, error)
	//mc: -key=slideArticlesKey -expire=d.mcArticlesIDExpire
	AddCacheListArtsId(c context.Context, buvid string, val *model.ArticleViewList) error
	//mc: -key=slideArticlesKey
	DelCacheListArtsId(c context.Context, buvid string) error
	//mc: -key=AnniversaryKey -expire=60*60*24*30
	CacheAnniversary(c context.Context, mid int64) (*model.AnniversaryInfo, error)
	//mc: -key=mcTagKey
	CacheAidsByTag(c context.Context, tag int64) (*model.TagArts, error)
	//mc: -key=mcTagKey -expire=d.mcArticleTagExpire
	AddCacheAidsByTag(c context.Context, tag int64, val *model.TagArts) error
	//mc: -key=mcUpStatKey -expire=d.mcUpStatDailyExpire
	CacheUpStatDaily(c context.Context, mid int64) (*model.UpStat, error)
	//mc: -key=mcUpStatKey -expire=d.mcUpStatDailyExpire
	AddCacheUpStatDaily(c context.Context, mid int64, val *model.UpStat) error
}

// RebuildUpListsCache .
func (d *Dao) RebuildUpListsCache(c context.Context, mid int64) (err error) {
	lists, err := d.RawUpLists(c, mid)
	if err != nil {
		return
	}
	return d.AddCacheUpLists(c, mid, lists)
}

// RebuildListReadCountCache .
func (d *Dao) RebuildListReadCountCache(c context.Context, id int64) (err error) {
	res, err := d.RawListReadCount(c, id)
	if err != nil {
		return
	}
	return d.AddCacheListReadCount(c, id, res)
}
