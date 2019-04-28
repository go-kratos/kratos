package service

import (
	"context"
	"go-common/app/service/main/up/api/v1"
	"go-common/app/service/main/up/model"
)

// GetUpInfoActive get up info active
func (s *Service) GetUpInfoActive(ctx context.Context, req *model.UpInfoActiveReq) (res *model.UpInfoActiveReply, err error) {
	return s.up.UpInfoActive(ctx, req.Mid)
}

// GetUpsInfoActive get ups info active
func (s *Service) GetUpsInfoActive(ctx context.Context, req *model.UpsInfoActiveReq) (res map[int64]*model.UpInfoActiveReply, err error) {
	return s.up.UpsInfoActive(ctx, req.Mids)
}

// GetHighAllyUps service get high ally ups
func (s *Service) GetHighAllyUps(ctx context.Context, req *v1.HighAllyUpsReq) (res *v1.HighAllyUpsReply, err error) {
	signUps, err := s.up.GetHighAllyUps(ctx, req.Mids)
	if err != nil {
		return
	}

	res = new(v1.HighAllyUpsReply)
	res.Lists = make(map[int64]*v1.SignUp)
	for _, signUp := range signUps {
		res.Lists[signUp.Mid] = &v1.SignUp{
			Mid:       signUp.Mid,
			State:     int32(signUp.State),
			BeginDate: signUp.BeginDate,
			EndDate:   signUp.EndDate,
		}
	}
	return
}
