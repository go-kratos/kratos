package service

import (
	"fmt"

	"context"
	"go-common/app/admin/main/tv/model"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/xstr"
)

const (
	_arcOnline  = 1
	_arcOffline = 2
)

// typeBubbleSort sort type
func typeBubbleSort(pTypes []model.UgcType) (pSortTypes []model.UgcType) {
	flag := true
	for i := 0; i < len(pTypes)-1; i++ {
		flag = true
		for j := 0; j < len(pTypes)-i-1; j++ {
			if pTypes[j].ID > pTypes[j+1].ID {
				pTypes[j], pTypes[j+1] = pTypes[j+1], pTypes[j]
				flag = false
			}
		}
		if flag {
			break
		}
	}
	pSortTypes = pTypes
	return
}

func (s *Service) existArcTps(passed bool) (existTypes map[int32]int, err error) {
	var (
		arcs []*model.Archive
		db   = s.DB.Where("deleted = ?", 0)
	)
	if passed {
		db = db.Where("result = ?", 1)
	}
	existTypes = make(map[int32]int)
	if err = db.Select("DISTINCT(typeid)").Find(&arcs).Error; err != nil {
		log.Error("DistinctType Error %v", err)
		return
	}
	for _, v := range arcs {
		existTypes[v.TypeID] = 1
	}
	return
}

//arcTp return archive type list
func (s *Service) arcTp(passed bool) (pTypes []model.UgcType, err error) {
	var (
		cTypeList  = make(map[int32][]model.UgcCType)
		oriPTypes  []model.UgcType
		existTypes map[int32]int
	)
	typeList := s.ArcTypes
	if existTypes, err = s.existArcTps(passed); err != nil {
		return
	}
	//make parent and child node sperate
	for _, v := range typeList {
		if v.Pid == 0 {
			oriPTypes = append(oriPTypes, model.UgcType{
				ID:   v.ID,
				Name: v.Name,
			})
		} else {
			cType := model.UgcCType{
				Pid:  v.Pid,
				ID:   v.ID,
				Name: v.Name,
			}
			if _, ok := existTypes[v.ID]; ok {
				cTypeList[v.Pid] = append(cTypeList[v.Pid], cType)
			}
		}
	}
	for _, v := range oriPTypes {
		if cValue, ok := cTypeList[v.ID]; ok {
			v.Children = cValue
			pTypes = append(pTypes, v)
		}
	}
	pTypes = typeBubbleSort(pTypes)
	return
}

func (s *Service) loadTps() {
	var (
		data = &model.AvailTps{}
		err  error
	)
	if data.AllTps, err = s.arcTp(false); err != nil {
		log.Error("loadTps Passed Err %v", err)
		return
	}
	if data.PassedTps, err = s.arcTp(true); err != nil {
		log.Error("loadTps All Err %v", err)
		return
	}
	if len(data.AllTps) > 0 || len(data.PassedTps) > 0 {
		s.avaiTps = data
	}
}

// GetTps get cms used types data
func (s *Service) GetTps(c context.Context, passed bool) (data []model.UgcType, err error) {
	if s.avaiTps == nil {
		err = ecode.ServiceUnavailable
		return
	}
	if passed {
		data = s.avaiTps.PassedTps
	} else {
		data = s.avaiTps.AllTps
	}
	return
}

//GetArchivePid get archive pid with child id
func (s *Service) GetArchivePid(id int32) (pid int32) {
	if value, ok := s.ArcTypes[id]; ok {
		pid = value.Pid
		return
	}
	return 0
}

func (s *Service) midTreat(param *model.ArcListParam) (mids []int64) {
	if param.Mid != 0 {
		return []int64{param.Mid}
	}
	if param.UpName != "" {
		var data []*model.Upper
		if err := s.DB.Where("ori_name LIKE ?", "%"+param.UpName+"%").Where("deleted = 0").Find(&data).Error; err != nil {
			log.Error("ArchiveList MidTreat UpName %s, Err %v", param.UpName, err)
			return
		}
		if len(data) > 0 {
			for _, v := range data {
				mids = append(mids, v.MID)
			}
		}
	}
	return
}

// ArchiveList is used for getting archive list
func (s *Service) ArchiveList(c *bm.Context, param *model.ArcListParam) (pager *model.ArcPager, err error) {
	var (
		archives []*model.ArcDB
		reqES    = new(model.ReqArcES)
		data     *model.EsUgcResult
		aids     []int64
		mids     []int64
		upsInfo  map[int64]string
	)
	reqES.FromArcListParam(param, s.typeidsTreat(param.Typeid, param.Pid))
	reqES.Mids = s.midTreat(param)
	pager = new(model.ArcPager)
	if data, err = s.dao.ArcES(c, reqES); err != nil {
		log.Error("ArchiveList Req %v, Err %v", param, err)
		return
	}
	pager.Page = data.Page
	if len(data.Result) == 0 {
		return
	}
	for _, v := range data.Result {
		aids = append(aids, v.AID)
		mids = append(mids, v.MID)
	}
	if err = s.DB.Order("mtime " + reqES.MtimeSort()).Where(fmt.Sprintf("aid IN (%s)", xstr.JoinInts(aids))).Find(&archives).Error; err != nil {
		log.Error("s.ArchiveList Find archives error(%v)", err)
		return
	}
	if upsInfo, err = s.pickUps(mids); err != nil {
		return
	}
	for _, v := range archives {
		item := v.ToList(s.GetArchivePid(v.TypeID))
		if name, ok := upsInfo[v.MID]; ok {
			item.UpName = name
		}
		pager.Items = append(pager.Items, item)
	}
	return
}

func (s *Service) pickUps(mids []int64) (res map[int64]string, err error) {
	if len(mids) == 0 {
		return
	}
	var resSlice []*model.CmsUpper
	res = make(map[int64]string, len(mids))
	if err = s.DB.Where(fmt.Sprintf("mid IN (%s)", xstr.JoinInts(mids))).Where("deleted = 0").Find(&resSlice).Error; err != nil {
		log.Error("pickUps Mids %v, Err %v", mids, err)
		return
	}
	for _, v := range resSlice {
		res[v.MID] = v.OriName
	}
	return
}

// ArcAction is used for online ugc archive
func (s *Service) ArcAction(ids []int64, action int) (err error) {
	var (
		w        = map[string]interface{}{"deleted": 0, "result": 1}
		tx       = s.DB.Model(&model.Archive{}).Begin()
		actValid int
	)
	if action == _arcOnline {
		actValid = 1
	} else if action == _arcOffline {
		actValid = 0
	} else {
		return ecode.TvDangbeiWrongType
	}
	for _, v := range ids {
		arch := model.Archive{}
		if errDB := tx.Where(w).Where("id=?", v).First(&arch).Error; errDB != nil {
			err = fmt.Errorf("找不到id为%v的数据", v)
			log.Error("s.ArcAction First error(%v)", err)
			tx.Rollback()
			return
		}
		if errDB := tx.Where("id=?", v).
			Update("valid", actValid).Error; errDB != nil {
			err = errDB
			log.Error("s.ArcAction Update error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()
	return
}

// ArcUpdate is used for update ugc archive
func (s *Service) ArcUpdate(id int64, cover string, content string, title string) (err error) {
	up := map[string]interface{}{
		"cover":   cover,
		"content": content,
		"title":   title,
	}
	if err = s.DB.Model(&model.Archive{}).Where("id=?", id).Update(up).Error; err != nil {
		log.Error("s.ArcUpdate Update error(%v)", err)
		return
	}
	return
}
