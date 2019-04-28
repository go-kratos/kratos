package capsule

import (
	"context"

	"go-common/app/interface/live/web-ucenter/conf"
	lotteryApi "go-common/app/service/live/xlottery/api/grpc/v1"
)

// Dao dao
type Dao struct {
	client *lotteryApi.Client
}

// New init
func New(c *conf.Config) (dao *Dao) {
	dao = &Dao{}
	client, err := lotteryApi.NewClient(conf.Conf.Warden)
	if err != nil {
		panic(err)
	}
	dao.client = client
	return
}

// GetDetail grpc
func (d *Dao) GetDetail(ctx context.Context, uid int64, from string) (*lotteryApi.CapsuleGetDetailResp, error) {
	return d.client.CapsuleClient.GetDetail(ctx, &lotteryApi.CapsuleGetDetailReq{Uid: uid, From: from})
}

// OpenCapsule grpc
func (d *Dao) OpenCapsule(ctx context.Context, uid int64, otype string, count int64, platform string) (*lotteryApi.CapsuleOpenCapsuleResp, error) {
	return d.client.CapsuleClient.OpenCapsule(ctx, &lotteryApi.CapsuleOpenCapsuleReq{Uid: uid, Type: otype, Count: count, Platform: platform})
}

// GetCapsuleInfo grpc
func (d *Dao) GetCapsuleInfo(ctx context.Context, uid int64, otype, from string) (*lotteryApi.CapsuleGetCapsuleInfoResp, error) {
	return d.client.CapsuleClient.GetCapsuleInfo(ctx, &lotteryApi.CapsuleGetCapsuleInfoReq{Uid: uid, Type: otype, From: from})
}

// OpenCapsuleByType grpc
func (d *Dao) OpenCapsuleByType(ctx context.Context, uid int64, otype string, count int64, platform string) (*lotteryApi.CapsuleOpenCapsuleByTypeResp, error) {
	return d.client.CapsuleClient.OpenCapsuleByType(ctx, &lotteryApi.CapsuleOpenCapsuleByTypeReq{Uid: uid, Type: otype, Count: count, Platform: platform})
}
