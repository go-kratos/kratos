package service

import (
	"context"
	"sort"
	"time"

	"go-common/app/service/main/vip/model"
)

func (s *Service) loadjointly() (err error) {
	s.jointlyList, err = s.dao.Jointlys(context.TODO(), time.Now().Unix())
	return
}

// Jointly jointly info.
func (s *Service) Jointly(c context.Context) (res []*model.JointlyResp, err error) {
	var (
		now = time.Now().Unix()
		tmp = []*model.Jointly{}
	)
	if len(s.jointlyList) == 0 {
		return
	}
	for _, v := range s.jointlyList {
		if v.StartTime < now && v.EndTime > now {
			tmp = append(tmp, v)
		}
	}
	if len(tmp) == 0 {
		return
	}
	sort.Slice(tmp, func(i int, j int) bool {
		return tmp[i].MTime.Time().After(tmp[j].MTime.Time())
	})
	for _, v := range tmp {
		res = append(res, &model.JointlyResp{
			Title:   v.Title,
			Content: v.Content,
			IsHot:   v.IsHot,
			Link:    v.Link,
			EndTime: v.EndTime,
		})
	}
	return
}
