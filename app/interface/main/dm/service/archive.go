package service

import (
	"context"
	"math"

	"go-common/app/service/main/archive/api"
	arcMdl "go-common/app/service/main/archive/model/archive"
	"go-common/library/log"
)

func (s *Service) archiveInfos(c context.Context, aids []int64) (archiveInfos map[int64]*api.Arc) {
	var (
		start, end int
	)
	archiveInfos = map[int64]*api.Arc{}
	if len(aids) <= 0 {
		return
	}
	page := int(math.Ceil(float64(len(aids)) / float64(100)))
	for i := 0; i < page; i++ {
		start = i * 100
		end = (i + 1) * 100
		if end > len(aids) {
			end = len(aids)
		}
		arg := &arcMdl.ArgAids2{Aids: aids[start:end]}
		infos, err := s.acvSvc.Archives3(c, arg)
		if err != nil {
			log.Error("s.arcRPC.Archives3(%v) error(%v)", arg, err)
			return
		}
		for _, info := range infos {
			archiveInfos[info.Aid] = info
		}
	}
	return
}
