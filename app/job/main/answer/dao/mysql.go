package dao

import (
	"context"

	"go-common/app/job/main/answer/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_questionByIDSQL      = "SELECT mid,question,ans,status FROM answer_extra_question WHERE id=?"
	_insQsSQL             = "INSERT INTO answer_extra_question(question,ans,av_id,status,source,state,origin_id) VALUES(?,?,?,?,?,?,?)"
	_updateQsSQL          = "UPDATE answer_extra_question SET ans=?,status=?,isdel=? WHERE origin_id=?"
	_questionExtraTypeSQL = "SELECT origin_id,ans,question,av_id,status,source FROM answer_extra_question WHERE isdel=1 and state=0 LIMIT ?"
	_updateStateSQL       = "UPDATE answer_extra_question SET state=? WHERE origin_id=?"
)

// ByID get question by id.
func (d *Dao) ByID(c context.Context, id int64) (que *model.LabourQs, err error) {
	var row = d.db.QueryRow(c, _questionByIDSQL, id)
	que = new(model.LabourQs)
	if err = row.Scan(&que.Mid, &que.Question, &que.Ans, &que.Status); err != nil {
		if err == sql.ErrNoRows {
			que = nil
			err = nil
			return
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

// AddQs add labour question log.
func (d *Dao) AddQs(c context.Context, qs *model.LabourQs) (err error) {
	if _, err = d.db.Exec(c, _insQsSQL, qs.Question, qs.Ans, qs.AvID, qs.Status, qs.Source, qs.State, qs.ID); err != nil {
		log.Error("AddQs: db.Exec(as:%v) error(%v)", qs, err)
	}
	return
}

// UpdateQs update question.
func (d *Dao) UpdateQs(c context.Context, que *model.LabourQs) (err error) {
	if _, err = d.db.Exec(c, _updateQsSQL, que.Ans, que.Status, que.Isdel, que.ID); err != nil {
		log.Error("setQs: db.Exec(%v) error(%v)", que, err)
	}
	return
}

// QidsExtraByState get extra question ids by check
func (d *Dao) QidsExtraByState(c context.Context, size int) (res []*model.LabourQs, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(c, _questionExtraTypeSQL, size); err != nil {
		log.Error("d._questionExtraTypeSQL.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		r := new(model.LabourQs)
		if err = rows.Scan(&r.ID, &r.Ans, &r.Question, &r.AvID, &r.Status, &r.Source); err != nil {
			log.Error("row.Scan() error(%v)", err)
			res = nil
			return
		}
		res = append(res, r)
	}
	err = rows.Err()
	return
}

// UpdateState update state.
func (d *Dao) UpdateState(c context.Context, que *model.LabourQs) (err error) {
	if _, err = d.db.Exec(c, _updateStateSQL, que.State, que.ID); err != nil {
		log.Error("UpdateState: db.Exec(%v) error(%v)", que, err)
	}
	return
}
