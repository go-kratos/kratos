package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/service/main/rank/model"
)

// Do .
func (s *Service) Do(c context.Context, arg *model.DoReq) error {
	switch arg.Action {
	case "all":
		go s.all(context.Background(), 0, 0)
	case "patchid":
		go s.all(context.Background(), arg.MinID, arg.MaxID)
	case "patchtime":
		timeLayout := "2006-01-02 15:04:05"
		loc, _ := time.LoadLocation("Local")
		beginTime, err := time.ParseInLocation(timeLayout, arg.BeginTime, loc)
		if err != nil {
			return err
		}
		endTime, err := time.ParseInLocation(timeLayout, arg.EndTime, loc)
		if err != nil {
			return err
		}
		go s.patch(context.Background(), beginTime, endTime)
	default:
		return nil
	}
	return nil
}

// Mget .
func (s *Service) Mget(c context.Context, arg *model.MgetReq) (*model.MgetResp, error) {
	res := new(model.MgetResp)
	tmap := make(map[int64]*model.Field)
	for _, id := range arg.Oids {
		field := s.field(id)
		tmap[id] = field
	}
	res.List = tmap
	return res, nil
}

// Sort .
func (s *Service) Sort(c context.Context, arg *model.SortReq) (*model.SortResp, error) {
	var isResult, isDeleted, isValid, isPid bool
	filter := new(model.Field)
	res := new(model.SortResp)
	res.Page = new(model.Page)
	res.Page.Pn = arg.Pn
	res.Page.Ps = arg.Ps
	for k, v := range arg.Filters {
		if k == "result" && v != "" {
			isResult = true
			result, _ := strconv.ParseInt(v, 10, 8)
			filter.Result = int8(result)
		}
		if k == "deleted" && v != "" {
			isDeleted = true
			deleted, _ := strconv.ParseInt(v, 10, 8)
			filter.Deleted = int8(deleted)
		}
		if k == "valid" && v != "" {
			isValid = true
			valid, _ := strconv.ParseInt(v, 10, 8)
			filter.Valid = int8(valid)
		}
		if k == "pid" && v != "" {
			isPid = true
			pid, _ := strconv.ParseInt(v, 10, 16)
			filter.Pid = int16(pid)
		}
	}
	fs := make([]*model.Field, 0)
	for _, oid := range arg.Oids {
		f := s.field(oid)
		if isResult && filter.Result != f.Result {
			continue
		}
		if isDeleted && filter.Deleted != f.Deleted {
			continue
		}
		if isValid && filter.Valid != f.Valid {
			continue
		}
		if isPid && filter.Pid != f.Pid {
			continue
		}
		if !f.Flag {
			continue
		}
		fs = append(fs, f)
	}
	if len(fs) == 0 {
		return res, nil
	}
	// deep copy
	fss := make([]*model.Field, 0)
	for _, v := range fs {
		cv := *v
		fss = append(fss, &cv)
	}
	// sort
	sort.Slice(fss, func(i, j int) bool {
		if arg.Field == "click" {
			if arg.Order == model.RankOrderByAsc {
				return fss[i].Click < fss[j].Click
			}
			return fss[i].Click > fss[j].Click
		}
		if arg.Field == "pubtime" {
			if arg.Order == model.RankOrderByAsc {
				return fss[i].Pubtime < fss[j].Pubtime
			}
			return fss[i].Pubtime > fss[j].Pubtime
		}
		return true
	})

	for _, f := range fss {
		res.Result = append(res.Result, f.Oid)
	}
	res.Page.Total = len(res.Result)
	start := (arg.Pn - 1) * arg.Ps
	end := arg.Pn * arg.Ps
	if start > len(res.Result) {
		res.Result = []int64{}
		return res, nil
	}
	if end > len(res.Result) || end == 0 {
		end = len(res.Result)
	}
	res.Result = res.Result[start:end]
	return res, nil
}

// Group .
func (s *Service) Group(c context.Context, arg *model.GroupReq) (*model.GroupResp, error) {
	res := new(model.GroupResp)
	tmap := make(map[int16]int)
	for _, oid := range arg.Oids {
		f := s.field(oid)
		if !f.Flag {
			continue
		}
		if _, ok := tmap[f.Pid]; ok {
			tmap[f.Pid]++
			continue
		}
		tmap[f.Pid] = 1
	}
	for k, v := range tmap {
		g := new(model.Group)
		g.Key = fmt.Sprintf("%d", k)
		g.Count = v
		res.List = append(res.List, g)
	}
	return res, nil
}
