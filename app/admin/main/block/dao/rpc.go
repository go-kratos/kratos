package dao

import (
	"context"

	mdlaccount "go-common/app/service/main/account/model"
	mdlfigure "go-common/app/service/main/figure/model"
	mdlspy "go-common/app/service/main/spy/model"

	"github.com/pkg/errors"
)

// SpyScore .
func (d *Dao) SpyScore(c context.Context, mid int64) (score int8, err error) {
	var (
		arg = &mdlspy.ArgUserScore{
			Mid: mid,
		}
		res *mdlspy.UserScore
	)
	if res, err = d.spyRPC.UserScore(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		return
	}
	score = res.Score
	return
}

// FigureRank .
func (d *Dao) FigureRank(c context.Context, mid int64) (rank int8, err error) {
	var (
		arg = &mdlfigure.ArgUserFigure{
			Mid: mid,
		}
		res *mdlfigure.FigureWithRank
	)
	if res, err = d.figureRPC.UserFigure(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		rank = 100
		return
	}
	rank = res.Percentage
	return
}

// AccountInfo .
func (d *Dao) AccountInfo(c context.Context, mid int64) (nickname string, tel int32, level int32, regTime int64, err error) {
	var (
		arg = &mdlaccount.ArgMid{
			Mid: mid,
		}
		res *mdlaccount.Profile
	)
	if res, err = d.accountRPC.Profile3(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	if res == nil {
		return
	}
	nickname = res.Name
	tel = res.TelStatus
	level = res.Level
	regTime = int64(res.JoinTime)
	return
}
