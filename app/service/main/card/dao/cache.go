package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/card/model"
)

func (d *Dao) cacheSFEquip(id int64) string {
	return strconv.FormatInt(id, 10)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch=50 -max_group=10 -batch_err=continue -nullcache=&model.UserEquip{CardID:-1} -check_null_code=$!=nil&&$.CardID==-1
	Equips(c context.Context, keys []int64) (map[int64]*model.UserEquip, error)
	// cache: -nullcache=&model.UserEquip{CardID:-1} -check_null_code=$!=nil&&$.CardID==-1 -singleflight=true
	Equip(c context.Context, key int64) (*model.UserEquip, error)
}
