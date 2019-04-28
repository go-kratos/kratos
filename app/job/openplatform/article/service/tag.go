package service

import (
	"context"
	"strings"

	"go-common/app/interface/main/tag/model"
	"go-common/app/job/openplatform/article/dao"
	"go-common/library/log"
)

// tags gets article tags.
func (s *Service) tags(c context.Context, aid int64) (res string, err error) {
	var (
		tags map[int64][]*model.Tag
		ts   []string
		arg  = &model.ArgResTags{Type: model.PicResType, Oids: []int64{aid}}
	)
	if tags, err = s.tagRPC.ResTags(c, arg); err != nil {
		log.Error("s.Tags(%d) error(%+v)", aid, err)
		dao.PromError("rpc:获取Tag")
		return
	} else if tags == nil || len(tags[aid]) == 0 {
		return
	}
	for _, t := range tags[aid] {
		ts = append(ts, t.Name)
	}
	res = strings.Join(ts, ",")
	return
}
