package bangumi

import (
	"context"

	episodegrpc "go-common/app/service/openplatform/pgc-season/api/grpc/episode/v1"

	"github.com/pkg/errors"
)

func (d *Dao) CardsInfoReply(c context.Context, episodeIds []int32) (res map[int32]*episodegrpc.EpisodeCardsProto, err error) {
	arg := &episodegrpc.EpReq{Epids: episodeIds}
	info, err := d.rpcClient.Cards(c, arg)
	if err != nil {
		err = errors.Wrapf(err, "%v", arg)
		return
	}
	res = info.Cards
	return
}
