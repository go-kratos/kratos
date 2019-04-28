package service

import (
	"context"
	"encoding/json"

	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/task"
	"go-common/library/xstr"
)

// QueryConfigs .
func (s *Service) QueryConfigs(c context.Context, queryParams *task.QueryParams) (ls []*task.Config, count int64, err error) {
	ls, count, err = s.gorm.QueryConfigs(c, queryParams)
	if err != nil || count == 0 {
		return
	}

	switch queryParams.ConfType {
	case task.TaskConfigAssign:
	case task.TaskConfigRangeWeight:
	case task.TaskConfigEqualWeight:
		if len(queryParams.IDFilter) > 0 || len(queryParams.TypeFilter) > 0 {
			confItem := &task.EqualWeightConfig{}
			lsres := []*task.Config{}
			for _, item := range ls {
				if err = json.Unmarshal([]byte(item.ConfJSON), confItem); err != nil {
					continue
				}
				if len(queryParams.IDFilter) > 0 {
					ids1, _ := xstr.SplitInts(confItem.IDs)
					ids2, _ := xstr.SplitInts(queryParams.IDFilter)
					if len(mergeSlice(ids1, ids2)) == 0 {
						continue
					}
				}
				if len(queryParams.TypeFilter) > 0 {
					ids, _ := xstr.SplitInts(queryParams.TypeFilter)
					if len(mergeSlice(ids, []int64{int64(confItem.Type)})) == 0 {
						continue
					}
				}
				lsres = append(lsres, item)
			}
			return lsres, count, err
		}
	}
	return ls, count, err
}

// UpdateConfig .
func (s *Service) UpdateConfig(c context.Context, config *task.Config) (err error) {
	return s.gorm.UpdateConfig(c, config)
}

// SetStateConfig .
func (s *Service) SetStateConfig(c context.Context, id int64, state int8) (err error) {
	return s.gorm.SetStateConfig(c, id, state)
}

// AddConfig .
func (s *Service) AddConfig(c context.Context, config *task.Config, confJSON interface{}) (err error) {
	return s.gorm.AddConfig(c, config, confJSON)
}

// DeleteConfig .
func (s *Service) DeleteConfig(c context.Context, id int64) (err error) {
	return s.gorm.DeleteConfig(c, id)
}

// WeightLog .
func (s *Service) WeightLog(c context.Context, taskid int64, pn, ps int) (ls []*model.WeightLog, count int, err error) {
	ls, count, err = s.searchWeightLog(c, taskid, pn, ps)
	if err != nil || len(ls) == 0 {
		return
	}

	var Name string
	if mid := ls[0].Mid; mid > 0 {
		if info, _ := s.rpc.Info3(c, mid); info != nil {
			Name = info.Name
		}
	}

	for _, item := range ls {
		item.MemName = Name
	}
	return
}
