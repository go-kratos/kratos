package tag

import (
	"context"

	"go-common/app/interface/main/tag/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

//CheckChannelReview check whether archive in channel review list
func (d *Dao) CheckChannelReview(c context.Context, aid int64) (in bool, channelIDs string, err error) {
	var res map[int64]*model.ResChannelCheckBack
	arg := &model.ArgResChannel{
		Oids: []int64{aid},
		Type: 3,
	}
	if res, err = d.tagDisRPC.ResChannelCheckBack(c, arg); err != nil {
		log.Error("CheckChannelReview d.tagDisRPC.ResChannelCheckBack error(%v) aid(%d)", err, aid)
		return
	}
	if res != nil && res[aid] != nil {
		in = res[aid].CheckBack == 1
		ids := []int64{}
		for cid := range res[aid].Channels {
			ids = append(ids, cid)
		}
		channelIDs = xstr.JoinInts(ids)
	} else {
		log.Warn("CheckChannelReview response(%+v) for aid(%d) is nil", res, aid)
	}
	return
}
