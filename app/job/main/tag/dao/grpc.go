package dao

import (
	"context"

	"go-common/app/job/main/tag/model"
	filgrpc "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/log"
)

// MFilter multi filter.
func (d *Dao) MFilter(c context.Context, msgs []string) (checked []string, err error) {
	var (
		res *filgrpc.MFilterReply
		n   = model.TagBatchNumMax
	)
	for len(msgs) > 0 {
		if n > len(msgs) {
			n = len(msgs)
		}
		msgMap := make(map[string]string, n)
		for _, tname := range msgs[:n] {
			msgMap[tname] = tname
		}
		msgs = msgs[n:]
		if res, err = d.filClient.MFilter(c, &filgrpc.MFilterReq{Area: "tag", MsgMap: msgMap}); err != nil {
			log.Error("d.MFilter(%v) error(%v)", msgs, err)
			return
		}
		for name, v := range res.RMap {
			if v.Level <= 10 {
				checked = append(checked, name)
			}
		}
	}
	return
}
