package v2

import (
	"context"
	"strconv"

	v2pb "go-common/app/interface/live/app-interface/api/http/v2"
)

// 首页banner
func (s *IndexService) getIndexBanner(ctx context.Context, platform string, device string, build int64) (resp []*v2pb.MBanner) {
	bizList := map[int64]int64{
		0: _bannerType,
	}
	moduleList := s.GetAllModuleInfoMapFromCache(ctx)
	for biz, moduleType := range bizList {
		for _, moduleInfo := range moduleList[moduleType] {
			bannerList, err := s.roomexDao.GetBanner(ctx, biz, 0, platform, device, build)
			if err != nil {
				continue
			}
			res := &v2pb.MBanner{}
			list := make([]*v2pb.PicItem, 0)

			for _, banner := range bannerList {
				id, _ := strconv.Atoi(banner.Id)
				list = append(list, &v2pb.PicItem{
					Id:      int64(id),
					Link:    banner.Link,
					Pic:     banner.Pic,
					Title:   banner.Title,
					Content: "",
				})
			}

			res.ModuleInfo = moduleInfo
			res.List = list

			resp = append(resp, res)
		}
	}
	return
}
