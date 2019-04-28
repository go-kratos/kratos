package dao

import (
	"context"

	"go-common/app/interface/main/esports/model"
)

//go:generate $GOPATH/src/go-common/app/tool/cache/gen
type _cache interface {
	// cache
	EpContests(c context.Context, ids []int64) (map[int64]*model.Contest, error)
	// cache
	EpSeasons(c context.Context, ids []int64) (map[int64]*model.Season, error)
	// cache
	EpTeams(c context.Context, ids []int64) (map[int64]*model.Team, error)
	// cache
	EpContestsData(c context.Context, ids []int64) (map[int64][]*model.ContestsData, error)
}
