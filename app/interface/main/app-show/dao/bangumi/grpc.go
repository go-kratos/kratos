package bangumi

import (
	"context"

	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"

	"github.com/pkg/errors"
)

func (d *Dao) CardsByAids(c context.Context, aids []int64) (res map[int32]*seasongrpc.CardInfoProto, err error) {
	var (
		tmpAids []int32
	)
	for _, aid := range aids {
		tmpAids = append(tmpAids, int32(aid))
	}
	arg := &seasongrpc.SeasonAidReq{Aids: tmpAids}
	info, err := d.rpcClient.CardsByAids(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	res = info.Cards
	return
}
