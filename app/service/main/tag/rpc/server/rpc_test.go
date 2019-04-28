package server

import (
	"context"
	"testing"
	"time"

	"go-common/app/service/main/tag/model"
	rpc "go-common/app/service/main/tag/rpc/client"

	. "github.com/smartystreets/goconvey/convey"
)

func WithRPC(f func(client *rpc.Service)) func() {
	return func() {
		client := rpc.New(nil)
		f(client)
	}
}

func Test_ResTags(t *testing.T) {
	Convey("ResTags", t, WithRPC(func(client *rpc.Service) {
		var (
			c = context.TODO()
		)
		argResTags := &model.ArgResTags{Oid: 1, Type: model.ResTypeMusic}
		client.ResTags(c, argResTags)
		argID := &model.ArgID{
			Mid: 14771787,
			ID:  5406908,
		}
		argIDs := &model.ArgIDs{
			Mid: 14771787,
			IDs: []int64{5406908, 6156592},
		}
		client.InfoByID(c, argID)
		client.InfoByIDs(c, argIDs)
		argName := &model.ArgName{
			Mid:  14771787,
			Name: "亚顺",
		}
		argNames := &model.ArgNames{
			Mid:   14771787,
			Names: []string{"帅哥", "亚顺"},
		}
		client.InfoByName(c, argName)
		argCheckName := &model.ArgCheckName{
			Name: "亚顺",
			Type: 1,
			Now:  time.Now(),
		}
		client.CheckName(c, argCheckName)
		client.Count(c, argID)
		client.Counts(c, argIDs)
		client.InfoByNames(c, argNames)
		argSub := &model.ArgSub{
			Mid:   14771787,
			Pn:    1,
			Ps:    2,
			Order: 1,
		}
		client.SubTags(c, argSub)
		argResTagLog := &model.ArgResTagLog{
			Oid:  8476387,
			Type: 3,
			Mid:  14771787,
			Pn:   1,
			Ps:   2,
		}
		client.ResTagLog(c, argResTagLog)
		ac := &model.ArgCustomSub{
			Type: 3,
			Mid:  14771787,
			Tids: []int64{12, 324, 432},
		}
		client.AddCustomSubTag(c, ac)
		ass := &model.ArgSub{
			Mid:   14771787,
			Type:  3,
			Pn:    1,
			Ps:    2,
			Order: 1,
		}
		client.CustomSubTag(c, ass)
		as := &model.ArgAddSub{
			Mid:  14771787,
			Tids: []int64{1},
		}
		client.AddSub(c, as)
		acs := &model.ArgCancelSub{
			Mid: 14771787,
			Tid: 1,
		}
		client.CancelSub(c, acs)
		ara := &model.ArgResAction{
			Oid:  8476387,
			Type: 3,
			Mid:  14771787,
			Tid:  1,
		}
		client.Like(c, ara)
		client.Hate(c, ara)
	}))
}
