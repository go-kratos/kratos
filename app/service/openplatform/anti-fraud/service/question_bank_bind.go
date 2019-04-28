package service

import (
	"context"
	"strconv"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// QuestionBankBind 绑定题库
func (s *Service) QuestionBankBind(c context.Context, args *model.ArgQuestionBankBinds) (err error) {
	// 检查题库是否存在
	qbIds := make(map[int64]int64)
	bindInfo := make(map[string]model.ArgQuestionBankBind)
	//todo add del item
	var targetItems []string
	var source int8
	var targetItemType int8
	var validateRst = true
	deleteItem := model.ArgQuestionBankUnbind{}
	for _, oneBindInfo := range args.BandBinds {
		if oneBindInfo.QsBId == 0 {
			deleteItem.Source = oneBindInfo.Source
			deleteItem.TargetItemType = oneBindInfo.TargetItemType
			tmpid, _ := strconv.ParseInt(oneBindInfo.TargetItems, 10, 64)
			deleteItem.TargetItems = append(deleteItem.TargetItems, tmpid)
		} else {
			qbIds[oneBindInfo.QsBId] = oneBindInfo.QsBId
			bindInfo[oneBindInfo.TargetItems] = oneBindInfo
			targetItems = append(targetItems, oneBindInfo.TargetItems)
			source = oneBindInfo.Source
			targetItemType = oneBindInfo.TargetItemType
		}
	}
	//验证题库数量或者题库名称
	if len(qbIds) != 0 {
		if validateRst, _ = s.validateBankIds(c, qbIds); !validateRst {
			err = ecode.NotEnoughQuestion
			return
		}
	}
	binds, err := s.d.GetBankBind(c, source, targetItemType, targetItems, true)
	if err != nil {
		log.Error("删除题库错误")
		return
	}

	if len(deleteItem.TargetItems) != 0 {
		_ = s.d.QuestionBankUnbind(c, deleteItem.TargetItems, deleteItem.TargetItemType, deleteItem.Source)
	}
	updateItem := make([]model.ArgQuestionBankBind, 0)
	insertItem := make([]model.ArgQuestionBankBind, 0)
	// 分类
	for _, targetItem := range targetItems {
		for _, bind := range binds {
			if bind.TargetItem == targetItem {
				// 新旧不一样才需要更新
				if bindInfo[targetItem].QsBId != bind.QsBId || bindInfo[targetItem].UseInTime != bind.QsBId || bind.IsDeleted != 0 {
					updateItem = append(updateItem, bindInfo[targetItem])
				}

				goto next
			}
		}
		insertItem = append(insertItem, bindInfo[targetItem])
	next:
	}

	err = s.d.AddBankBind(c, updateItem, insertItem)
	return
}

// GetQuestionBankBind 获取绑定信息
func (s *Service) GetQuestionBankBind(c context.Context, args *model.ArgGetBankBind) (list []*model.QuestionBankBind, err error) {
	list, err = s.d.GetBindBank(c, args.Source, args.TargetItemType, args.TargetItems)
	return
}

// GetBindItem 绑定题目查询
func (s *Service) GetBindItem(c context.Context, args *model.ArgGetBindItems) (resp model.RespList, err error) {
	list, total, err := s.d.GetBindItem(c, args.QsBId, args.PageNo, args.PageSize)

	resp = model.RespList{}
	resp.PageSize = args.PageSize
	resp.PageNo = args.PageNo
	resp.Total = total
	resp.Items = list
	return
}

// QuestionBankUnbind 解绑
func (s *Service) QuestionBankUnbind(c context.Context, args *model.ArgQuestionBankUnbind) (err error) {
	err = s.d.QuestionBankUnbind(c, args.TargetItems, args.TargetItemType, args.Source)
	return
}

// validateBankIds 小于三个
func (s *Service) validateBankIds(c context.Context, qbIds map[int64]int64) (res bool, err error) {
	res = false
	list, err := s.d.GetBankInfoByQBid(c, qbIds)
	if len(list) != len(qbIds) {
		return
	}
	for _, v := range list {
		if v.TotalCnt < 3 {
			return
		}
	}
	res = true
	return
}
