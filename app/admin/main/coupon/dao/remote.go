package dao

import (
	"context"

	"go-common/app/admin/main/coupon/model"
	seasongrpc "go-common/app/service/openplatform/pgc-season/api/grpc/season/v1"

	"github.com/pkg/errors"
)

//GetPGCInfo get pgc info.
func (d *Dao) GetPGCInfo(c context.Context, oid int32) (r *model.PGCInfoResq, err error) {
	var (
		params *seasongrpc.SeasonInfoReq
		oids   = make([]int32, 0)
		reply  *seasongrpc.CardsInfoReply
	)
	oids = append(oids, oid)
	params = &seasongrpc.SeasonInfoReq{
		SeasonIds: oids,
	}
	if reply, err = d.rpcClient.Cards(c, params); err != nil {
		err = errors.WithStack(err)
		return
	}
	if proto, success := reply.Cards[oid]; success {
		r = new(model.PGCInfoResq)
		r.Title = proto.Title
	}
	return
}
