package dao

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_addQuestionBankBind = "INSERT INTO question_bank_bind (target_item, target_item_type, bank_id, use_in_time, source, is_deleted) VALUES %s"

	_getQuestionBankBindWithDeletedNew = " SELECT a.target_item,a.target_item_type,a.bank_id,a.use_in_time,b.id,b.qb_name" +
		" FROM question_bank_bind as a  join question_bank as b" +
		" ON a.bank_id = b.qb_id" +
		" where a.target_item IN(%s) and a.target_item_type = ? and a.source = ?"

	_getQuestionBankBind = _getQuestionBankBindWithDeletedNew + " AND a.is_deleted = 0"

	_updateQuestionBankBind = "UPDATE question_bank_bind set bank_id = ?, use_in_time = ?, is_deleted = 0" +
		" WHERE source = ?" +
		" AND target_item_type = ?" +
		" AND target_item IN(%s)"

	_deleteQuestionBankBind = "UPDATE question_bank_bind SET is_deleted = 1 WHERE target_item IN(%s)" +
		" AND target_item_type = ?" +
		" AND source = ?"

	_getBindItem   = "SELECT id, target_item, target_item_type, bank_id, use_in_time, source, is_deleted, ctime, mtime FROM question_bank_bind WHERE bank_id = ? AND is_deleted = 0 limit ?,?"
	_countBindItem = "SELECT COUNT(1) as count FROM question_bank_bind WHERE bank_id = ? AND is_deleted = 0"
)

// AddBankBind add
func (d *Dao) AddBankBind(c context.Context, update []model.ArgQuestionBankBind, insert []model.ArgQuestionBankBind) (err error) {
	lenUpdate := len(update)
	lenInsert := len(insert)
	if lenUpdate == 0 && lenInsert == 0 {
		return
	}
	updateItems := map[int64]*model.ArgQuestionBankBindToDb{}
	insertItems := map[int64]*model.ArgQuestionBankBindToDb{}
	for _, v := range update {
		if _, ok := updateItems[v.QsBId]; !ok {
			updateItems[v.QsBId] = &model.ArgQuestionBankBindToDb{
				QsBId:          v.QsBId,
				Source:         v.Source,
				TargetItemType: v.TargetItemType,
				UseInTime:      v.UseInTime,
				TargetItems:    []string{v.TargetItems},
			}
		} else {
			updateItems[v.QsBId].TargetItems = append(updateItems[v.QsBId].TargetItems, v.TargetItems)
		}
		d.DelTargetItemBindCache(c, v.TargetItems)
	}

	for _, v := range insert {
		if _, ok := insertItems[v.QsBId]; !ok {
			insertItems[v.QsBId] = &model.ArgQuestionBankBindToDb{
				QsBId:          v.QsBId,
				Source:         v.Source,
				TargetItemType: v.TargetItemType,
				UseInTime:      v.UseInTime,
				TargetItems:    []string{v.TargetItems},
			}
		} else {
			insertItems[v.QsBId].TargetItems = append(insertItems[v.QsBId].TargetItems, v.TargetItems)
		}

	}

	tx, err := d.db.Begin(c)
	if err != nil {
		log.Error("d.AddBankBind(%v) error(%v)", err)
		return
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Error("tx.Rollback() error(%v)", err)
			}
			return
		}

		if err = tx.Commit(); err != nil {
			log.Error("tx.Commit() error(%v)", err)
			return
		}
	}()

	if len(updateItems) > 0 {
		for _, v := range updateItems {
			sql := fmt.Sprintf(_updateQuestionBankBind, `"`+strings.Join(v.TargetItems, `","`)+`"`)
			_, err = tx.Exec(sql, v.QsBId, v.UseInTime, v.Source, v.TargetItemType)
			if err != nil {
				log.Error("d.AddBankBind(%v, %v) error(%v)", update, insert, err)
				return
			}
		}
	}

	if len(insertItems) > 0 {
		for _, v := range insertItems {
			lenInsert = len(v.TargetItems)
			if lenInsert > 0 {
				placeholder := strings.Trim(strings.Repeat("(?, ?, ?, ?, ?, ?),", lenInsert), ",")
				values := make([]interface{}, 0)
				for _, ins := range v.TargetItems {
					values = append(values, ins, v.TargetItemType, v.QsBId, v.UseInTime, v.Source, 0)
				}
				_, err = tx.Exec(fmt.Sprintf(_addQuestionBankBind, placeholder), values...)
				if err != nil {
					log.Error("d.AddBankBind() tx.Exec(%s) error(%v)", fmt.Sprintf(_addQuestionBankBind, values), err)
					return
				}
			}
		}
	}

	return
}

// GetBankBind 查询绑定关系
func (d *Dao) GetBankBind(c context.Context, source, targetItemType int8, targetItem []string, withDeleted bool) (list []*model.QuestionBankBind, err error) {
	list = make([]*model.QuestionBankBind, 0)
	if len(targetItem) == 0 {
		return
	}

	sql := _getQuestionBankBind
	if withDeleted {
		sql = _getQuestionBankBindWithDeletedNew
	}

	sql = fmt.Sprintf(sql, `"`+strings.Join(targetItem, `","`)+`"`)
	rows, err := d.db.Query(c, sql, targetItemType, source)
	if err != nil {
		log.Error("d.GetBankBind(%v, %v, %v) db.Query() error(%v)", source, targetItemType, targetItem, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tmp := &model.QuestionBankBind{QuestionBank: new(model.QuestionBank)}
		err = rows.Scan(&tmp.TargetItem, &tmp.TargetItemType, &tmp.QsBId, &tmp.UseInTime, &tmp.ID, &tmp.QuestionBank.QBName)
		if err != nil {
			return
		}

		list = append(list, tmp)
	}

	return
}

// GetBindBank 获取已绑定的题库信息
func (d *Dao) GetBindBank(c context.Context, source, targetItemType int8, targetItem []string) (binds []*model.QuestionBankBind, err error) {
	// list = make([]*model.QuestionBankBind, 0)

	binds, err = d.GetBankBind(c, source, targetItemType, targetItem, false)
	if len(binds) == 0 {
		return
	}

	var bankIds []int64
	for _, bind := range binds {
		bankIds = append(bankIds, bind.QsBId)
	}

	banks, err := d.GetQusBankListByIds(c, bankIds)
	if err != nil {
		return
	}

	for _, bind := range binds {
		for _, bank := range banks {
			if bind.QsBId == bank.QsBId {
				bind.QuestionBank = bank
				break
			}
		}
	}

	return
}

// CountBindItem cnt
func (d *Dao) CountBindItem(c context.Context, bankID int64) (count int64, err error) {
	if err = d.db.QueryRow(c, _countBindItem, bankID).Scan(&count); err != nil {
		log.Error("d.CountBindItem(%d) error(%v)", bankID, err)
	}
	return
}

// QuestionBankUnbind 解绑
func (d *Dao) QuestionBankUnbind(c context.Context, delIds []int64, targetType int8, source int8) (err error) {
	if len(delIds) > 0 {
		sql := fmt.Sprintf(_deleteQuestionBankBind, xstr.JoinInts(delIds))
		_, err = d.db.Exec(c, sql, targetType, source)
		if err != nil {
			log.Error("d.QuestionBankUnbind() tx.Exec(%s) error(%v)", fmt.Sprintf(_deleteQuestionBankBind, xstr.JoinInts(delIds)), err)
			return
		}
	}
	for _, delID := range delIds {
		d.DelTargetItemBindCache(c, strconv.FormatInt(delID, 10))
	}

	return
}

// GetBindItem itm
func (d *Dao) GetBindItem(c context.Context, bankID int64, page, pageSize int) (list []*model.QuestionBankBind, total int64, err error) {
	list = make([]*model.QuestionBankBind, 0)
	rows, err := d.db.Query(c, _getBindItem, bankID, (page-1)*pageSize, pageSize)
	if err != nil {
		log.Error("d.GetBindItem(%d) error(%v)", bankID, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		tmp := &model.QuestionBankBind{}
		err = rows.Scan(&tmp.ID, &tmp.TargetItem, &tmp.TargetItemType, &tmp.QsBId, &tmp.UseInTime, &tmp.Source, &tmp.IsDeleted, &tmp.Ctime, &tmp.Mtime)
		if err != nil {
			return
		}

		list = append(list, tmp)
	}

	total, err = d.CountBindItem(c, bankID)
	if err != nil {
		log.Error("d.GetBindItem(%d, %d, %d) error(%v)", bankID, page, pageSize, err)
		return
	}

	return
}
