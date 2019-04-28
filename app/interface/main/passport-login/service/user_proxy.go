package service

import (
	"context"
	"go-common/app/interface/main/passport-login/model"
)

// ProxyCheckUser check user .
func (s *Service) ProxyCheckUser(c context.Context, param *model.ParamLogin) (decodeUser *model.DecodeUser, err error) {
	return s.CheckUser(context.Background(), param.UserName, param.Pwd)
}
