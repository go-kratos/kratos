package service

import (
	"context"
	"fmt"
	"time"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// GetQusBankInfo 获取题库信息
func (s *Service) GetQusBankInfo(c context.Context, qbid int64) (res *model.QuestionBank, err error) {
	if res, err = s.d.GetQusBankInfo(c, qbid); err != nil {
		log.Info(fmt.Sprintf("获取信息失败: res: %v qbid: %d err:%s", res, qbid, err.Error()))
		return
	}
	return
}

// AddQusBank 添加题库
func (s *Service) AddQusBank(c context.Context, in *model.ArgAddQusBank) (res model.AddReturn, err error) {
	qb := &model.QuestionBank{
		QsBId:        time.Now().UnixNano() / 1e6,
		QBName:       in.QBName,
		CdTime:       in.CdTime,
		MaxRetryTime: in.MaxRetryTime,
	}
	id, err := s.d.InsertQusBank(c, qb)
	if err != nil {
		log.Warn("s.dao.AddQusBank(%v) error(%v)", qb)
		return
	}
	res.ID, err = s.d.GetQBId(c, id)
	return
}

// DelQusBank 删除题库
func (s *Service) DelQusBank(c context.Context, qbid int64, status int8) (res bool, err error) {
	var row int64

	cnt, err := s.d.CountBindItem(c, qbid)
	if cnt > 0 || err != nil {
		err = ecode.BankUsing
		return
	}

	if row, err = s.d.DelQusBank(c, qbid, status); err != nil {
		log.Error("s.dao.DelQusBank(%v) error(%v)", err)
		err = ecode.ParamInvalid
		return
	}
	if row != 0 {
		res = true
	}
	return
}

// UpdataQusBank 更新题库
func (s *Service) UpdataQusBank(c context.Context, id int64, name string, trytime int64, cdtime int64) (res bool, err error) {

	var row int64
	if row, err = s.d.UpdateQusBank(c, id, name, trytime, cdtime); err != nil {
		log.Error("s.dao.DelQusBank(%v) error(%v)", err)
		err = ecode.UpdateError
		return
	}
	if row != 0 {
		err = s.d.DelQusBankCache(c, id)
		if err != nil {
			log.Error("s.dao.DelQusBankCache error(%v)", err)
			return
		}
		res = true
	}
	return
}

// GetQusBanklist 获取题库列表
func (s *Service) GetQusBanklist(c context.Context, pageNo int, pageSize int, name string) (res []*model.QusBankSt, err error) {
	offset, limit := s.PageInfo(pageNo, pageSize)
	if res, err = s.d.StatisticsQusBank(c, offset, limit, name); err != nil {
		return
	}
	return
}

// GetQusBankTotal 获取题库数量
func (s *Service) GetQusBankTotal(c context.Context, name string) (res int64, err error) {
	if res, err = s.d.GetQusBankCount(c, name); err != nil {
		return
	}
	return
}

// PageInfo 分页
func (s *Service) PageInfo(pageNo int, pageSize int) (offset int, limit int) {
	if pageNo >= 0 && pageSize > 0 {
		offset = (pageNo - 1) * pageSize
		limit = pageSize
	} else {
		offset = model.STARTINDEX
		limit = model.PAGESIZE
	}
	return
}

// QusBankCheck 获取total数量小于cnt题库列表
func (s *Service) QusBankCheck(c context.Context, in *model.ArgCheckQus) (res []*model.QuestionBank, err error) {
	if res, err = s.d.GetQusBankList(c, in.Cnt, in.QusIDs); err != nil {
		return
	}
	return
}
