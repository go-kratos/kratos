package cache

import (
	"sync"
	"time"

	"go-common/app/interface/main/creative/model/tag"
)

var (
	//TagCache tag's cache key=>tag_id, value=>*tag.Meta
	TagCache = make(map[int64]*tag.Meta)
	tagMutex sync.Mutex
)

//ClearTagCache clear all tag cache
func ClearTagCache(tm time.Time) {
	TagCache = make(map[int64]*tag.Meta)
}

//AddTagCache add tag cache
func AddTagCache(meta *tag.Meta) {
	if meta == nil {
		return
	}
	tagMutex.Lock()
	TagCache[meta.TagID] = meta
	tagMutex.Unlock()
}

//GetTagCache get tag cache
func GetTagCache(ids ...int64) (result map[int64]*tag.Meta, leftIDs []int64) {
	result = make(map[int64]*tag.Meta)
	tagMutex.Lock()
	for _, v := range ids {
		var d, ok = TagCache[v]
		if !ok {
			leftIDs = append(leftIDs, v)
			continue
		}
		result[v] = d
	}
	tagMutex.Unlock()

	return
}
