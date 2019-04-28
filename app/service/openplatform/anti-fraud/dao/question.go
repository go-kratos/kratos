package dao

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_getQusSQL            = "SELECT question_id ,question_type,answer_type,question_name,question_bank_id,difficulty,is_deleted FROM question where question_id=? and is_deleted =0  "
	_addQusSQL            = "INSERT INTO question(question_id ,question_type,answer_type,question_name,question_bank_id,difficulty) VALUES(?,?,?,?,?,?)"
	_getQuslistBywhereSQL = "SELECT question_id ,question_type,answer_type,question_name,question_bank_id,difficulty,is_deleted FROM question  where  is_deleted = 0 and question_bank_id = ? order by id desc   limit ? ,?  "
	_getAllQusByBankID    = "SELECT question_id FROM question WHERE question_bank_id = ? AND is_deleted = 0"
	_getQuslistSQL        = "SELECT question_id ,question_type,answer_type,question_name,question_bank_id,difficulty,is_deleted FROM question  where  is_deleted = ?   limit ? ,?  "
	_getQusCntSQL         = "SELECT  COUNT(*)  FROM question where is_deleted = 0  "
	_delQusSQL            = "UPDATE question set is_deleted = ? where question_id =? "
	_updateQusSQL         = "UPDATE question set question_type = ? ,answer_type=? ,question_name = ? , question_bank_id = ? ,difficulty = ?  where question_id =? and is_deleted =0  "
	_addAnswerSQL         = "INSERT INTO question_answer(answer_id ,question_id,answer_content,is_correct ) VALUES(?,?,?,?)"
	_multyAddAnswerSQL    = "INSERT INTO question_answer(answer_id ,question_id,answer_content,is_correct ) VALUES %s "
	_updateAnswerSQL      = "UPDATE question_answer set answer_content = ? , is_correct =? , is_deleted =0  where answer_id = ? "
	_delAnswerSQL         = "UPDATE question_answer set is_deleted = 1 where question_id =? "
	_getAnswerListSQL     = "select  answer_id ,question_id, answer_content,is_correct  from question_answer where question_id =? and is_deleted =0  "
	_addUserAnswerSQL     = "INSERT INTO question_user_answer(uid ,question_id,platform,source,answers,is_correct ) VALUES(?,?,?,?,?,?)"
	_checkAnswerSQL       = "SELECT count(1) from  question_answer where is_correct =1 and  question_id = ? and  answer_id IN(%s) "
	_getRandomPicSQL      = "SELECT x,y,src from question_verify_pic where id =? "
	_getListPicSQL        = "SELECT id from question_verify_pic limit ? ,? "
	_getPicCntSQL         = "SELECT  COUNT(*)  FROM question_verify_pic "
)

// GetQusInfo info
func (d *Dao) GetQusInfo(c context.Context, qid int64) (oi *model.Question, err error) {
	oi = &model.Question{}
	row := d.db.QueryRow(c, _getQusSQL, qid)
	if err = row.Scan(&oi.QsID, &oi.QsType, &oi.AnswerType, &oi.QsName, &oi.QsBId, &oi.QsDif, &oi.IsDeleted); err != nil {
		if err == sql.ErrNoRows {
			oi = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// InsertQus add
func (d *Dao) InsertQus(c context.Context, oi *model.Question) (lastID int64, err error) {
	res, err := d.db.Exec(c, _addQusSQL, oi.QsID, oi.QsType, oi.AnswerType, oi.QsName, oi.QsBId, oi.QsDif)
	if err != nil {
		log.Error("[dao.question|GetQusList] d.db.Query err: %v", err)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// GetQusList list
func (d *Dao) GetQusList(c context.Context, offset int, limitnum int, qBid int64) (res []*model.Question, err error) {
	res = make([]*model.Question, 0)
	_sql := _getQuslistSQL
	if qBid > 0 {
		_sql = _getQuslistBywhereSQL
	}
	rows, err := d.db.Query(c, _sql, qBid, offset, limitnum)

	if err != nil {
		log.Error("[dao.question|GetQusList] d.db.Query err: %v %d,%d", err, offset, limitnum)
		return
	}

	defer rows.Close()
	for rows.Next() {
		oi := &model.Question{}
		if err = rows.Scan(&oi.QsID, &oi.QsType, &oi.AnswerType, &oi.QsName, &oi.QsBId, &oi.QsDif, &oi.IsDeleted); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// GetQusIds ids
func (d *Dao) GetQusIds(c context.Context, bankID int64) (ids []int64, err error) {
	if ids = d.GetBankQuestionsCache(c, bankID); len(ids) > 1 {
		return
	}
	rows, err := d.db.Query(c, _getAllQusByBankID, bankID)
	if err != nil {
		log.Error("d.GetQusIds(%d) error(%v)", bankID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var temp int64
		if err = rows.Scan(&temp); err != nil {
			log.Error("d.GetQusIds(%d) rows.Scan() error(%v)", bankID, err)
			return
		}
		ids = append(ids, temp)
	}
	d.SetBankQuestionsCache(c, bankID, ids)
	return
}

// DelQus del
func (d *Dao) DelQus(c context.Context, qid int64) (affect int64, err error) {
	res, err := d.db.Exec(c, _delQusSQL, 1, qid)
	if err != nil {
		log.Error("d.DelQus(qbid:%d, dmid:%d) error(%v)", qid, 1, err)
		return
	}
	return res.RowsAffected()
}

// UpdateQus update
func (d *Dao) UpdateQus(c context.Context, update *model.ArgUpdateQus, answers []model.Answer) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateQusSQL, update.Type, update.AnType, update.Name, update.BId, update.Dif, update.QsID)
	if err != nil {
		log.Error("d.UpdateQus(qbid:%d, dmid:%d) error(%v)", update.QsID, update.QsID, err)
		return
	}

	return res.RowsAffected()
}

// GetQusCount cnt
func (d *Dao) GetQusCount(c context.Context, bid int64) (total int64, err error) {
	var cntSQL string
	if bid == 0 {
		cntSQL = _getQusCntSQL
		err = d.db.QueryRow(c, cntSQL).Scan(&total)
	} else {
		cntSQL = _getQusCntSQL + "and question_bank_id = ?"
		err = d.db.QueryRow(c, cntSQL, bid).Scan(&total)
	}
	if err != nil {
		log.Error("d.GetQusCount error(%v)", err)
		return
	}
	return

}

// InserAnwser add
func (d *Dao) InserAnwser(c context.Context, answer *model.AnswerAdd) (affect int64, err error) {
	res, err := d.db.Exec(c, _addAnswerSQL, answer.AnswerID, answer.QsID, answer.AnswerContent, answer.IsCorrect)
	if err != nil {
		log.Error("d.InserAnwser() error(%v)", err)
		return
	}
	affect, err = res.LastInsertId()
	return
}

// MultiAddAnwser add
func (d *Dao) MultiAddAnwser(c context.Context, answers []*model.AnswerAdd) (err error) {
	length := len(answers)
	if length == 0 {
		return
	}
	values := strings.Trim(strings.Repeat("(?, ?, ?, ?),", length), ",")
	args := make([]interface{}, 0)
	for _, ins := range answers {
		AnswerID := time.Now().UnixNano() / 1e6
		time.Sleep(time.Millisecond)
		args = append(args, AnswerID, ins.QsID, ins.AnswerContent, ins.IsCorrect)
	}
	_, err = d.db.Exec(c, fmt.Sprintf(_multyAddAnswerSQL, values), args...)
	if err != nil {
		log.Error("d.InserAnwser() error(%v)", err)
		return
	}

	return
}

// UpdateAnwser upd
func (d *Dao) UpdateAnwser(c context.Context, answer *model.AnswerAdd) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateAnswerSQL, answer.AnswerContent, answer.IsCorrect, answer.AnswerID)
	if err != nil {
		log.Error("d.UpdateAnwser() error(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// DelAnwser del
func (d *Dao) DelAnwser(c context.Context, qusID int64) (affect int64, err error) {
	res, err := d.db.Exec(c, _delAnswerSQL, qusID)
	if err != nil {
		log.Error("d.UpdateAnwser() error(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// GetAnswerList list
func (d *Dao) GetAnswerList(c context.Context, qusID int64) (res []*model.Answer, err error) {
	res = make([]*model.Answer, 0)
	rows, err := d.db.Query(c, _getAnswerListSQL, qusID)
	if err != nil {
		log.Error("[dao.GetAnswerList] d.db.Query err: %v", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		oi := &model.Answer{}
		if err = rows.Scan(&oi.AnswerID, &oi.QsID, &oi.AnswerContent, &oi.IsCorrect); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// AddUserAnwser add
func (d *Dao) AddUserAnwser(c context.Context, answer *model.ArgCheckAnswer, isCorrect int8) (affect int64, err error) {

	ids := xstr.JoinInts(answer.Answers)
	res, err := d.db.Exec(c, _addUserAnswerSQL, answer.UID, answer.QsID, answer.Platform, answer.Source, ids, isCorrect)
	if err != nil {
		log.Error("d.InserAnwser() error(%v)", err)
		return
	}
	affect, err = res.LastInsertId()
	return
}

// CheckAnswer check
func (d *Dao) CheckAnswer(c context.Context, qsid int64, ids []int64) (total int, err error) {
	err = d.db.QueryRow(c, fmt.Sprintf(_checkAnswerSQL, xstr.JoinInts(ids)), qsid).Scan(&total)
	if err != nil {
		log.Error("d.GetQusBankCount error(%v)", err)
		return
	}
	return
}

// GetRandomPic get
func (d *Dao) GetRandomPic(c context.Context, id int) (oi *model.QuestBkPic, err error) {

	oi = &model.QuestBkPic{}

	row := d.db.QueryRow(c, _getRandomPicSQL, id)
	if err = row.Scan(&oi.X, &oi.Y, &oi.Src); err != nil {
		if err == sql.ErrNoRows {
			oi = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}

	}
	return

}

// GetAllPicIds ids
func (d *Dao) GetAllPicIds(c context.Context, offset int, limitnum int) (ids []int, err error) {

	rows, err := d.db.Query(c, _getListPicSQL, offset, limitnum)
	if err != nil {
		log.Error("[dao.GetAllPicIds] d.db.Query err: %v %d,%d", err, offset, limitnum)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var oi int
		if err = rows.Scan(&oi); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", ids)
			return
		}
		ids = append(ids, oi)
	}
	return

}

// GetPicCount cnt
func (d *Dao) GetPicCount(c context.Context) (total int, err error) {
	err = d.db.QueryRow(c, _getPicCntSQL).Scan(&total)
	if err != nil {
		log.Error("d.GetQusCount error(%v)", err)
		return
	}
	return

}
