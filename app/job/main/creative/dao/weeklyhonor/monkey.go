package weeklyhonor

import (
	"context"
	model "go-common/app/interface/main/creative/model/weeklyhonor"
	upgrpc "go-common/app/service/main/up/api/v1"
	"reflect"

	"github.com/bouk/monkey"
)

//MockHonorStat is
func (d *Dao) MockHonorStat(stat *model.HonorStat, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "HonorStat", func(_ *Dao, _ context.Context, _ int64, _ string) (*model.HonorStat, error) {
		return stat, err
	})
}

//MockSendNotify is
func (d *Dao) MockSendNotify(err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "SendNotify", func(_ *Dao, _ context.Context, _ []int64) error {
		return err
	})
}

//MockLatestHonorLogs is
func (d *Dao) MockLatestHonorLogs(hls []*model.HonorLog, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "LatestHonorLogs", func(_ *Dao, _ context.Context, _ []int64) ([]*model.HonorLog, error) {
		return hls, err
	})
}

//MockClickCounts is
func (d *Dao) MockClickCounts(res map[int64]int32, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "ClickCounts", func(_ *Dao, _ context.Context, _ []int64) (map[int64]int32, error) {
		return res, err
	})
}

//MockUpCount is
func (d *Dao) MockUpCount(count int, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpCount", func(_ *Dao, _ context.Context, _ int64) (int, error) {
		return count, err
	})
}

//MockUpActivesList is
func (d *Dao) MockUpActivesList(upActives []*upgrpc.UpActivity, newid int64, err error) (guard *monkey.PatchGuard) {
	return monkey.PatchInstanceMethod(reflect.TypeOf(d), "UpActivesList", func(_ *Dao, _ context.Context, _ int64, _ int) ([]*upgrpc.UpActivity, int64, error) {
		return upActives, newid, err
	})
}
