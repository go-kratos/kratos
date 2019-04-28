package gorpc

import (
	"context"
	"testing"

	"go-common/app/interface/main/history/model"
	rpcClient "go-common/app/interface/main/history/rpc/client"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx    = context.TODO()
	client *rpcClient.Service
)

func WithRPC(f func(client *rpcClient.Service)) func() {
	return func() {
		client = rpcClient.New(nil)
		f(client)
	}
}
func Test_Histroy_rpc(t *testing.T) {
	Convey("rpc client Add", t, WithRPC(func(client *rpcClient.Service) {
		arg := &model.ArgHistory{
			Mid: 14771787,
			History: &model.History{
				Aid: 17406762,
				TP:  1,
				Pro: 122,
			},
		}
		client.Add(ctx, arg)
	}))
	Convey("rpc client preogress", t, WithRPC(func(client *rpcClient.Service) {
		arg := &model.ArgPro{
			Mid:  14771787,
			Aids: []int64{17406762},
		}
		client.Progress(ctx, arg)
	}))
	Convey("rpc client delete", t, WithRPC(func(client *rpcClient.Service) {
		r := &model.Resource{Oid: 100, Business: "archive"}
		arg := &model.ArgDelete{
			Mid:       14771787,
			Resources: []*model.Resource{r},
		}
		client.Delete(ctx, arg)
	}))
	Convey("rpc client history", t, WithRPC(func(client *rpcClient.Service) {
		arg := &model.ArgHistories{Mid: 14771787}
		client.History(ctx, arg)
	}))
}
