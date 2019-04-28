package bws

import (
	"context"
	"fmt"

	bwsmdl "go-common/app/interface/main/activity/model/bws"
)

func midKey(id int64) string {
	return fmt.Sprintf("u_m_%d", id)
}

func keyKey(key string) string {
	return fmt.Sprintf("u_k_%s", key)
}
func pointsKey(id int64) string {
	return fmt.Sprintf("b_p_%d", id)
}

func achievesKey(id int64) string {
	return fmt.Sprintf("b_a_%d", id)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	//cache: -sync=true
	UsersMid(c context.Context, key int64) (*bwsmdl.Users, error)
	//cache: -sync=true
	UsersKey(c context.Context, key string) (*bwsmdl.Users, error)
	//cache: -sync=true
	Points(c context.Context, bid int64) (*bwsmdl.Points, error)
	//cache: -sync=true
	Achievements(c context.Context, bid int64) (*bwsmdl.Achievements, error)
	//cache: -sync=true
	UserAchieves(c context.Context, bid int64, key string) ([]*bwsmdl.UserAchieve, error)
	//cache: -sync=true
	UserPoints(c context.Context, bid int64, key string) ([]*bwsmdl.UserPoint, error)
	//cache: -sync=true
	AchieveCounts(c context.Context, bid int64, day string) ([]*bwsmdl.CountAchieves, error)
}

//go:generate $GOPATH/src/go-common/app/tool/cache/mc
type _mc interface {
	//mc: -key=midKey
	CacheUsersMid(c context.Context, key int64) (*bwsmdl.Users, error)
	//mc: -key=midKey -expire=d.mcExpire -encode=pb
	AddCacheUsersMid(c context.Context, key int64, value *bwsmdl.Users) error
	//mc: -key=midKey
	DelCacheUsersMid(c context.Context, key int64) error
	//mc: -key=keyKey
	CacheUsersKey(c context.Context, key string) (*bwsmdl.Users, error)
	//mc: -key=keyKey -expire=d.mcExpire -encode=pb
	AddCacheUsersKey(c context.Context, key string, value *bwsmdl.Users) error
	//mc: -key=keyKey
	DelCacheUsersKey(c context.Context, key string) error
	//mc: -key=pointsKey
	CachePoints(c context.Context, key int64) (*bwsmdl.Points, error)
	//mc: -key=pointsKey -expire=d.mcExpire -encode=pb
	AddCachePoints(c context.Context, key int64, value *bwsmdl.Points) error
	//mc: -key=pointsKey
	DelCachePoints(c context.Context, key int64) error
	//mc: -key=achievesKey
	CacheAchievements(c context.Context, key int64) (*bwsmdl.Achievements, error)
	//mc: -key=achievesKey -expire=d.mcExpire -encode=pb
	AddCacheAchievements(c context.Context, key int64, value *bwsmdl.Achievements) error
	//mc: -key=achievesKey
	DelCacheAchievements(c context.Context, key int64) error
}
