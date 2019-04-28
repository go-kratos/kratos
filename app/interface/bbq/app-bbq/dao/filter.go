package dao

import (
	"context"

	filter "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/log"
)

// 用于filter
const (
	FilterAreaAccount = "BBQ_account"
	FilterAreaVideo   = "BBQ_video"
	FilterAreaReply   = "BBQ_reply"
	FilterAreaSearch  = "BBQ_search"
	FilterAreaDanmu   = "BBQ_danmu"
	FilterLevel       = 20
)

// Filter .
func (d *Dao) Filter(ctx context.Context, content string, area string) (level int32, err error) {
	req := new(filter.FilterReq)
	req.Message = content
	req.Area = area
	reply, err := d.filterClient.Filter(ctx, req)
	if err != nil {
		log.Errorv(ctx, log.KV("log", "filter fail : req="+req.String()))
		return
	}
	level = reply.Level
	log.V(1).Infov(ctx, log.KV("log", "get filter reply="+reply.String()))
	return
}
