package service

import (
	"go-common/app/interface/main/push-archive/dao"
	"go-common/app/interface/main/push-archive/model"
)

var (
	// 存放实验组mid尾号
	fansTestGroup = make(map[int]struct{})
	// 存放对照组mid尾号
	fansComparisonGroup = make(map[int]struct{})
	// 指定mid放进测试组
	fansTestMids = make(map[int64]struct{})
)

func (s *Service) mappingAbtest() {
	for _, n := range s.c.Abtest.TestGroup {
		fansTestGroup[n] = struct{}{}
	}
	for _, n := range s.c.Abtest.ComparisonGroup {
		fansComparisonGroup[n] = struct{}{}
	}
	for _, n := range s.c.Abtest.TestMids {
		fansTestMids[n] = struct{}{}
	}
}

// 将所有粉丝通过abtest规则，拆分成 (实验流量||对照组流量) && 其余流量
func (s *Service) fansByAbtest(group *dao.FanGroup, fans []int64) (result, others []int64) {
	for _, fan := range fans {
		n := int(fan % 10)
		if group.Hitby == model.GroupDataTypeAbtest {
			if _, ok := fansTestMids[fan]; ok {
				result = append(result, fan)
				continue
			}
			if _, ok := fansTestGroup[n]; ok {
				result = append(result, fan)
				continue
			}
		} else if group.Hitby == model.GroupDataTypeAbComparison {
			if _, ok := fansComparisonGroup[n]; ok {
				result = append(result, fan)
				continue
			}
		}
		others = append(others, fan)
	}
	return
}
