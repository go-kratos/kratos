package block

import (
	"context"

	account "go-common/app/service/main/account/api"
	mdlfigure "go-common/app/service/main/figure/model"
	mdlspy "go-common/app/service/main/spy/model"
	"go-common/library/net/metadata"

	"github.com/pkg/errors"
)

// SpyScore .
func (s *Service) SpyScore(c context.Context, mid int64) (score int8, err error) {
	var (
		arg = &mdlspy.ArgUserScore{
			Mid: mid,
		}
		res *mdlspy.UserScore
	)
	if res, err = s.spyRPC.UserScore(c, arg); err != nil {
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
func (s *Service) FigureRank(c context.Context, mid int64) (rank int8, err error) {
	var (
		arg = &mdlfigure.ArgUserFigure{
			Mid: mid,
		}
		res *mdlfigure.FigureWithRank
	)
	if res, err = s.figureRPC.UserFigure(c, arg); err != nil {
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
func (s *Service) AccountInfo(c context.Context, mid int64) (nickname string, tel int32, level int32, regTime int64, err error) {
	var (
		arg = &account.MidReq{
			Mid:    mid,
			RealIp: metadata.String(c, metadata.RemoteIP),
		}
		profileReply *account.ProfileReply
	)
	if profileReply, err = s.accountClient.Profile3(c, arg); err != nil {
		err = errors.WithStack(err)
		return
	}
	res := profileReply.Profile
	nickname = res.Name
	tel = res.TelStatus
	level = res.Level
	regTime = int64(res.JoinTime)
	return
}
