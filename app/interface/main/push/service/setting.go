package service

import (
	"context"

	pb "go-common/app/service/main/push/api/grpc/v1"
	"go-common/library/log"
)

// Setting gets user notify setting.
func (s *Service) Setting(ctx context.Context, mid int64) (st map[int32]int32, err error) {
	arg := &pb.SettingRequest{Mid: mid}
	reply, err := s.pushRPC.Setting(ctx, arg)
	if err != nil {
		log.Error("s.pushRPC.Setting(%+v) error(%v)", arg, err)
		return
	}
	return reply.Settings, nil
}

// SetSetting saves setting.
func (s *Service) SetSetting(ctx context.Context, mid int64, typ, val int) (err error) {
	arg := &pb.SetSettingRequest{Mid: mid, Type: int32(typ), Value: int32(val)}
	if _, err = s.pushRPC.SetSetting(ctx, arg); err != nil {
		log.Error("s.pushRPC.SetSetting(%+v) error(%v)", arg, err)
	}
	return
}
