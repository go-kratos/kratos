package bangumi

import (
	"context"

	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"

	"github.com/pkg/errors"
)

// CardsInfoReply pgc cards info
func (d *Dao) CardsInfoReply(c context.Context, seasonIds []int32) (res map[int32]*seasongrpc.CardInfoProto, err error) {
	arg := &seasongrpc.SeasonInfoReq{SeasonIds: seasonIds}
	info, err := d.rpcClient.Cards(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	res = info.Cards
	return
}
