package operation

import (
	"context"
	"encoding/json"
	operMdl "go-common/app/interface/main/creative/model/operation"
	"go-common/library/ecode"
	"go-common/library/log"
	"strconv"
	"strings"
	"time"
)

// 0 => web+app; 1=> app; 2=> web; 100=>全平台
const (
	remarkShowBanner = "show_banner"
)

// Tool get a tool down.
func (s *Service) Tool(c context.Context, ty int8) (ops []*operMdl.Operation, err error) {
	var tyStr string
	if ty == 0 || ty == 1 { //创作中心页的tool下载
		tyStr = "icon"
	} else if ty == 2 { //投稿页的tool下载
		tyStr = "submit_icon"
	}
	if op, ok := s.toolCache[tyStr]; ok {
		ops = op
	}
	return
}

// WebOperations get full operations.
func (s *Service) WebOperations(c context.Context) (ops map[string][]*operMdl.Operation, err error) {
	if s.operCache == nil {
		err = ecode.NothingFound
		return
	}
	ops = make(map[string][]*operMdl.Operation)
	for _, v := range s.operCache {
		if v.Platform == 0 || v.Platform == 2 {
			o := &operMdl.Operation{}
			o.ID = v.ID
			o.Ty = v.Ty
			o.Rank = v.Rank
			o.Pic = v.Pic
			o.Link = v.Link
			o.Content = v.Content
			o.Remark = v.Remark
			o.Note = v.Note
			o.Stime = v.Stime
			o.Etime = v.Etime
			if o.Ty == "play" || o.Ty == "collect_arc" {
				o.Ty = "board"
			}
			ops[o.Ty] = append(ops[o.Ty], o)
		}
	}
	for _, ty := range operMdl.FullTypes() {
		trimTy := strings.Trim(ty, "'")
		if _, ok := ops[trimTy]; !ok {
			ops[trimTy] = []*operMdl.Operation{}
		}
	}
	return
}

// AppBanner get app index flexslider; filter by platform + business + remark
func (s *Service) AppBanner(c context.Context) (bns []*operMdl.Banner, cbns []*operMdl.BannerCreator, err error) {
	bns = make([]*operMdl.Banner, 0, len(s.operCache))
	cbns = make([]*operMdl.BannerCreator, 0, len(s.operCache))
	for _, v := range s.operCache {
		// 平台切除web平台;
		// 允许征稿启示也显示在所有APP的Banner上;
		// 必须开启remark=="show_banner"的业务校验;
		if (v.Ty != "play" && v.Ty != "collect_arc") ||
			v.Platform == 2 ||
			v.Remark != remarkShowBanner {
			continue
		}
		bn := &operMdl.Banner{}
		cbn := &operMdl.BannerCreator{}
		bn.Ty = v.Ty
		bn.Rank = v.Rank
		cbn.Ty = v.Ty
		cbn.Rank, _ = strconv.Atoi(v.Rank)
		if v.AppPic == "" {
			pics := []string{}
			if err = json.Unmarshal([]byte(v.Pic), &pics); err != nil {
				log.Error("json.Unmarshal(%v) error(%v)", string(v.Pic), err)
			}
			bn.Pic = pics[len(pics)-1]
			cbn.Pic = pics[len(pics)-1]
		} else {
			bn.Pic = v.AppPic
			cbn.Pic = v.AppPic
		}
		bn.Link = v.Link
		bn.Content = v.Content
		cbn.Link = v.Link
		cbn.Content = v.Content
		bns = append(bns, bn)
		cbns = append(cbns, cbn)
	}
	return
}

// CreatorOperationList get operations list.
func (s *Service) CreatorOperationList(c context.Context, pn, ps int) (list *operMdl.BannerList, err error) {
	if s.operCache == nil {
		err = ecode.NothingFound
		return
	}
	list = &operMdl.BannerList{Pn: pn, Ps: ps}
	// notice: s.CreatorRelOperCache["play"] 已经进行合并
	play, ok := s.CreatorRelOperCache["play"]
	if !ok {
		return
	}
	bcs := make([]*operMdl.BannerCreator, 0)
	for _, v := range play {
		bc := &operMdl.BannerCreator{}
		bc.Ty = v.Ty
		bc.Rank, _ = strconv.Atoi(v.Rank)
		if v.AppPic == "" {
			pics := []string{}
			if err = json.Unmarshal([]byte(v.Pic), &pics); err != nil {
				log.Error("json.Unmarshal(%v) error(%v)", string(v.Pic), err)
			}
			bc.Pic = pics[len(pics)-1]
		} else {
			bc.Pic = v.AppPic
		}
		bc.Link = v.Link
		bc.Content = v.Content
		st, _ := time.Parse("2006-01-02 15:04:05", v.Stime)
		bc.Stime = st.Unix()
		et, _ := time.Parse("2006-01-02 15:04:05", v.Etime)
		bc.Etime = et.Unix()
		bcs = append(bcs, bc)
	}
	total := len(bcs)
	list.Total = total
	start := (pn - 1) * ps
	end := pn * ps
	if total <= start {
		list.BannerCreator = make([]*operMdl.BannerCreator, 0)
	} else if total <= end {
		list.BannerCreator = bcs[start:total]
	} else {
		list.BannerCreator = bcs[start:end]
	}
	return
}

// AppOperationList get operations list.
func (s *Service) AppOperationList(c context.Context, pn, ps int, tp string) (list *operMdl.BannerList, err error) {
	tpOK := false
	for _, fullTp := range operMdl.FullTypes() {
		trimTy := strings.Trim(fullTp, "'")
		if trimTy == tp {
			tpOK = true
			break
		}
	}
	if s.allRelOperCache == nil || !tpOK {
		return
	}
	list = &operMdl.BannerList{Pn: pn, Ps: ps}
	vals := make([]*operMdl.Operation, 0)
	for _, v := range s.allRelOperCache {
		if v.Ty == tp && (v.Platform == 0 || v.Platform == 1) {
			vals = append(vals, v)
		}
	}
	if len(vals) == 0 {
		return
	}
	bcs := make([]*operMdl.BannerCreator, 0)
	for _, v := range vals {
		bc := &operMdl.BannerCreator{}
		bc.Ty = v.Ty
		bc.Rank, _ = strconv.Atoi(v.Rank)
		if v.AppPic == "" {
			pics := []string{}
			if err = json.Unmarshal([]byte(v.Pic), &pics); err != nil {
				log.Error("json.Unmarshal(%v) error(%v)", string(v.Pic), err)
			}
			bc.Pic = pics[len(pics)-1]
		} else {
			bc.Pic = v.AppPic
		}
		bc.Link = v.Link
		bc.Content = v.Content
		st, _ := time.Parse("2006-01-02 15:04:05", v.Stime)
		bc.Stime = st.Unix()
		et, _ := time.Parse("2006-01-02 15:04:05", v.Etime)
		bc.Etime = et.Unix()
		bcs = append(bcs, bc)
	}
	total := len(bcs)
	list.Total = total
	start := (pn - 1) * ps
	end := pn * ps
	if total <= start {
		list.BannerCreator = make([]*operMdl.BannerCreator, 0)
	} else if total <= end {
		list.BannerCreator = bcs[start:total]
	} else {
		list.BannerCreator = bcs[start:end]
	}
	return
}
