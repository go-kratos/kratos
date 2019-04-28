package dao

import (
	"context"
	"fmt"

	"go-common/app/service/main/thumbup/model"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// 用户点赞列表
	// cache: -singleflight=true -ignores=||start,end
	userLikeList(c context.Context, mid int64, businessID int64, state int8, start, end int) (res []*model.ItemLikeRecord, err error)
}

func (d *Dao) cacheSFuserLikeList(mid, businessID int64, state int8, start, end int) string {
	return fmt.Sprintf("sf_u%v_%v_%v", mid, businessID, state)
}

// UserLikeList 用户点赞列表
func (d *Dao) UserLikeList(c context.Context, mid int64, businessID int64, state int8, start, end int) (res []*model.ItemLikeRecord, err error) {
	var ls []*model.ItemLikeRecord
	ls, err = d.userLikeList(c, mid, businessID, state, start, end)
	for _, x := range ls {
		if x.MessageID != -1 {
			res = append(res, x)
		}
	}
	return
}
