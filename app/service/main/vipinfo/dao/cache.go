package dao

import (
	"context"

	"go-common/app/service/main/vipinfo/model"
)

// 注意 -check_null_code=$!=nil 生成的代码内容要去掉！
//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache: -batch=50 -max_group=10 -batch_err=continue -nullcache=&model.VipUserInfo{VipType:0} -check_null_code=$!=nil
	Infos(c context.Context, keys []int64) (map[int64]*model.VipUserInfo, error)
	// cache: -nullcache=&model.VipUserInfo{VipType:0} -check_null_code=$!=nil
	Info(c context.Context, key int64) (*model.VipUserInfo, error)
}
