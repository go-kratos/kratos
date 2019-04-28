package dao

import (
	"context"
	"database/sql"
	"fmt"

	"go-common/app/admin/main/answer/model"
	"go-common/library/log"
)

const (
	_questionTable = "ans_v3_question"
	_typeTable     = "ans_v3_question_type"
	_queHistory    = "answer_history_%d"
)

// QueByID by id.
func (d *Dao) QueByID(c context.Context, id int64) (res *model.Question, err error) {
	que := &model.QuestionDB{}
	if err := d.db.Table(_questionTable).Where("id=?", id).First(que).Error; err != nil {
		return nil, err
	}
	res = &model.Question{QuestionDB: que, Ans: []string{que.Ans1, que.Ans2, que.Ans3, que.Ans4}}
	return
}

// ByIDs by id.
func (d *Dao) ByIDs(c context.Context, IDs []int64) (res []*model.QuestionDB, err error) {
	if err := d.db.Table(_questionTable).Where("id in (?)", IDs).Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// IDsByState by id.
func (d *Dao) IDsByState(c context.Context) (res []int64, err error) {
	var (
		rows *sql.Rows
	)
	if rows, err = d.db.Table(_questionTable).Select("id").Where("state = ?", 1).Rows(); err != nil {
		return nil, err
	}
	for rows.Next() {
		var qid int64
		if err = rows.Scan(&qid); err != nil {
			return
		}
		res = append(res, qid)
	}
	return
}

// QuestionAdd add register question.
func (d *Dao) QuestionAdd(c context.Context, q *model.QuestionDB) (aff int64, err error) {
	db := d.db.Save(q)
	if err = db.Error; err != nil {
		return
	}
	aff = db.RowsAffected
	return
}

// QuestionEdit edit register question.
func (d *Dao) QuestionEdit(c context.Context, arg *model.QuestionDB) (aff int64, err error) {
	que := map[string]interface{}{
		"question": arg.Question,
		"ans1":     arg.Ans1,
		"ans2":     arg.Ans2,
		"ans3":     arg.Ans3,
		"ans4":     arg.Ans4,
		"operator": arg.Operator,
	}
	db := d.db.Table(_questionTable).Omit("ctime, operator").Where("id = ?", arg.ID).Updates(que)
	if err = db.Error; err != nil {
		log.Error("%+v", err)
		return
	}
	aff = db.RowsAffected
	return
}

// QuestionList .
func (d *Dao) QuestionList(c context.Context, arg *model.ArgQue) (res []*model.QuestionDB, err error) {
	db := d.db.Table(_questionTable)
	if arg.TypeID != 0 {
		db = db.Where("type_id=?", arg.TypeID)
	}
	if arg.State != -1 {
		db = db.Where("state=?", arg.State)
	}
	if len(arg.Question) != 0 {
		db = db.Where("question LIKE '%%" + arg.Question + "%%'")
	}
	db = db.Offset((arg.Pn - 1) * arg.Ps).Limit(arg.Ps).Order("id desc")
	if err = db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// QuestionCount question page total count.
func (d *Dao) QuestionCount(c context.Context, arg *model.ArgQue) (res int64, err error) {
	db := d.db.Table(_questionTable)
	if arg.TypeID != 0 {
		db = db.Where("type_id=?", arg.TypeID)
	}
	if arg.State != -1 {
		db = db.Where("state=?", arg.State)
	}
	if len(arg.Question) != 0 {
		db = db.Where("question LIKE '%%" + arg.Question + "%%'")
	}
	if err = db.Count(&res).Error; err != nil {
		return 0, err
	}
	return
}

// UpdateStatus update question state.
func (d *Dao) UpdateStatus(c context.Context, state int8, qid int64, operator string) (aff int64, err error) {
	val := map[string]interface{}{
		"state":    state,
		"operator": operator,
	}
	db := d.db.Table(_questionTable).Where("id=?", qid).Updates(val)
	if err = db.Error; err != nil {
		return
	}
	aff = db.RowsAffected
	return

}

// Types get all types.
func (d *Dao) Types(c context.Context) (res []*model.TypeInfo, err error) {
	db := d.db.Table(_typeTable)
	if err = db.Where("parentid != 0").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// TypeSave add register question type.
func (d *Dao) TypeSave(c context.Context, t *model.TypeInfo) (aff int64, err error) {
	db := d.db.Save(t)
	if err = db.Error; err != nil {
		return
	}
	aff = db.RowsAffected
	return
}

// BaseQS .
func (d *Dao) BaseQS(c context.Context) (res []*model.QuestionDB, err error) {
	db := d.db.Table("ans_register_question")
	db = db.Where("type_id=6 AND state=1")
	if err = db.Omit("id").Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// AllQS .
func (d *Dao) AllQS(c context.Context) (res []*model.QuestionDB, err error) {
	db := d.db.Table("ans_v3_question")
	if err = db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// InsBaseQs .
func (d *Dao) InsBaseQs(c context.Context, qs *model.QuestionDB) (lastID int64, err error) {
	db := d.db.Table(_questionTable)
	qs.TypeID = 36
	qs.State = 1
	db = db.Create(qs)
	if err = db.Error; err != nil {
		return
	}
	lastID = qs.ID
	return
}

// QueHistory .
func (d *Dao) QueHistory(c context.Context, arg *model.ArgHistory) (res []*model.AnswerHistoryDB, err error) {
	db := d.db.Table(fmt.Sprintf(_queHistory, arg.Mid%10))
	db = db.Where("mid = ?", arg.Mid).Offset((arg.Pn - 1) * arg.Ps).Limit(arg.Ps).Order("id desc")
	if err = db.Find(&res).Error; err != nil {
		return nil, err
	}
	return
}

// HistoryCount .
func (d *Dao) HistoryCount(c context.Context, arg *model.ArgHistory) (res int64, err error) {
	db := d.db.Table(fmt.Sprintf(_queHistory, arg.Mid%10))
	if err = db.Where("mid = ?", arg.Mid).Count(&res).Error; err != nil {
		return 0, err
	}
	return
}
