package dao

import (
	"bytes"
	"context"
	"fmt"
	"strconv"

	"go-common/app/service/openplatform/anti-fraud/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_QusBankInfo          = "SELECT qb_id ,qb_name,cd_time,max_retry_time,is_deleted FROM question_bank where qb_id= ?  and is_deleted= 0 "
	_addQusBank           = "insert into  question_bank(qb_id ,qb_name,cd_time,max_retry_time,is_deleted) values(?,?,?,?,?)"
	_qusBanklist          = "SELECT qb_id ,qb_name,cd_time,max_retry_time,is_deleted FROM question_bank where  is_deleted= 0  and total_cnt < ? and qb_id IN(%s) order by id desc  "
	_qusBankcnt           = "SELECT  COUNT(*)  FROM question_bank where  is_deleted= 0   "
	_delQusBank           = "update question_bank set is_deleted = ? where qb_id =? "
	_searchQusBank        = "SELECT qb_id ,qb_name,cd_time,max_retry_time,is_deleted FROM question_bank where  is_deleted= 0 and qb_name  LIKE ? "
	_updateQusBank        = "update question_bank set qb_name = ? ,max_retry_time=? ,cd_time = ? where qb_id =? "
	_qusBanklistByids     = "SELECT qb_id, qb_name, cd_time, max_retry_time, is_deleted FROM question_bank WHERE is_deleted = 0 and qb_id IN(%s)"
	_getStaticQusBankList = "SELECT id, qb_id,qb_name,cd_time,max_retry_time,is_deleted,total_cnt ,easy_cnt,normal_cnt,hard_cnt from question_bank where is_deleted =0  "
	_getQbID              = "SELECT qb_id FROM question_bank where id= ? and is_deleted= 0"
	_getQusBankByID       = "select  a.qb_id,a.qb_name,a.cd_time,a.max_retry_time , a.is_deleted, count(if(b.is_deleted = 0,true,null)) total " +
		"from question_bank a  left join question b on a.qb_id = b.question_bank_id WHERE a.qb_id IN(%s) GROUP BY a.qb_id HAVING a.is_deleted = 0   "
	_getQusBankCnt = "SELECT sum(case when difficulty = 1 then 1 else 0 end )  'easy', sum(case when difficulty = 2 then 1 else 0 end )  'normal'," +
		"sum(case when difficulty = 3  then 1 else 0 end )  'hard',count(id)  'total' from question WHERE is_deleted =0 and  `question_bank_id`= ?"
	_updateQusBankCnt = "update question_bank set easy_cnt = ? ,normal_cnt=? ,hard_cnt = ?,total_cnt = ? where qb_id =? "
)

// GetQusBankInfo info
func (d *Dao) GetQusBankInfo(c context.Context, qbid int64) (oi *model.QuestionBank, err error) {
	oi = &model.QuestionBank{}
	row := d.db.QueryRow(c, _QusBankInfo, qbid)
	if err = row.Scan(&oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted); err != nil {
		if err == sql.ErrNoRows {
			oi = nil
			err = nil
		} else {
			log.Error("row.Scan error(%v)", err)
		}
	}
	return
}

// InsertQusBank add
func (d *Dao) InsertQusBank(c context.Context, oi *model.QuestionBank) (lastID int64, err error) {
	res, err := d.db.Exec(c, _addQusBank, oi.QsBId, oi.QBName, oi.CdTime, oi.MaxRetryTime, oi.IsDeleted)
	if err != nil {
		log.Error("[dao.question|GetQusBankList] d.db.Query err: %v", err)
		return
	}
	lastID, err = res.LastInsertId()
	return
}

// StatisticsQusBank 统计bank
func (d *Dao) StatisticsQusBank(c context.Context, offset int, limitnum int, name string) (res []*model.QusBankSt, err error) {
	res = make([]*model.QusBankSt, 0)

	var rows *sql.Rows
	var _sql string
	if name == "" {
		_sql = _getStaticQusBankList + " order by id desc   limit ? ,?   "
		rows, err = d.db.Query(c, _sql, offset, limitnum)

	} else {
		_sql = _getStaticQusBankList + " and   qb_name   LIKE ?  order by id desc   limit ? ,?   "
		name = "%" + name + "%"
		rows, err = d.db.Query(c, _sql, name, offset, limitnum)
	}

	if err != nil {
		log.Error("[dao.question|GetQusBankList] d.db.Query err: %v %d,%d", err, offset, limitnum)
		return
	}

	defer rows.Close()
	for rows.Next() {
		oi := &model.QusBankSt{}
		if err = rows.Scan(&oi.ID, &oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted, &oi.TotalCnt, &oi.EasyCnt, &oi.NormalCnt, &oi.HardCnt); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// GetQusBankList list
func (d *Dao) GetQusBankList(c context.Context, cnt int, ids []int64) (res []*model.QuestionBank, err error) {
	res = make([]*model.QuestionBank, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_qusBanklist, xstr.JoinInts([]int64(ids))), cnt)
	if err != nil {
		log.Error("[dao.question|GetQusBankList] d.db.Query err: %v %d,%d", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		oi := &model.QuestionBank{}
		if err = rows.Scan(&oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// DelQusBank del
func (d *Dao) DelQusBank(c context.Context, qbID int64, status int8) (affect int64, err error) {
	res, err := d.db.Exec(c, _delQusBank, status, qbID)
	if err != nil {
		log.Error("d.UpdateQusBank(qbid:%d, dmid:%d) error(%v)", qbID, status, err)
		return
	}
	s, e := res.RowsAffected()
	log.Error("d.UpdateQusBank(qbxxxxxid:%d, dxxxxxmid:%d) error(%v)", qbID, s, e)

	return res.RowsAffected()
}

// UpdateQusBank update
func (d *Dao) UpdateQusBank(c context.Context, qbID int64, name string, trytime int64, cdtime int64) (affect int64, err error) {
	res, err := d.db.Exec(c, _updateQusBank, name, trytime, cdtime, qbID)
	if err != nil {
		log.Error("d.UpdateQusBank(qbid:%d, dmid:%d) error(%v)", qbID, trytime, err)
		return
	}
	return res.RowsAffected()
}

// GetQusBankCount cnt
func (d *Dao) GetQusBankCount(c context.Context, name string) (total int64, err error) {

	var cntSQL string
	if name == "" {
		cntSQL = _qusBankcnt
		err = d.db.QueryRow(c, cntSQL).Scan(&total)
	} else {
		cntSQL = _qusBankcnt + " and  qb_name   LIKE ? "
		name = "%" + name + "%"
		err = d.db.QueryRow(c, cntSQL, name).Scan(&total)
	}
	if err != nil {
		log.Error("d.GetQusBankCount error(%v)", err)
		return
	}
	return

}

// BankSearch search
func (d *Dao) BankSearch(c context.Context, name string) (res []*model.QuestionBank, err error) {
	res = make([]*model.QuestionBank, 0)
	rows, err := d.db.Query(c, _searchQusBank, "%"+name+"%")
	if err != nil {
		log.Error("[dao.question|BankSearch] d.db.Query err: %v", err)
		return
	}

	defer rows.Close()
	for rows.Next() {
		oi := &model.QuestionBank{}
		if err = rows.Scan(&oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// GetQusBankListByIds ids
func (d *Dao) GetQusBankListByIds(c context.Context, ids []int64) (res []*model.QuestionBank, err error) {
	res = make([]*model.QuestionBank, 0)
	rows, err := d.db.Query(c, fmt.Sprintf(_qusBanklistByids, xstr.JoinInts([]int64(ids))))
	if err != nil {
		log.Error("[dao.question|GetQusBankListByIds] d.db.Query err: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		oi := &model.QuestionBank{}
		if err = rows.Scan(&oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted); err != nil {
			log.Error("[dao.question|GetOrder] rows.Scan err: %v", res)
			return
		}
		res = append(res, oi)
	}
	return
}

// GetQBId id
func (d *Dao) GetQBId(c context.Context, id int64) (qbID int64, err error) {
	err = d.db.QueryRow(c, _getQbID, id).Scan(&qbID)
	if err != nil {
		log.Error("d.GetQBId error(%v)", err.Error())
		return
	}
	return
}

// GetBankInfoByQBid byids
func (d *Dao) GetBankInfoByQBid(c context.Context, qbID map[int64]int64) (res []*model.QusBankSt, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getQusBankByID, d.CoverStr(qbID)))
	fmt.Println(fmt.Sprintf(_getQusBankByID, d.CoverStr(qbID)))
	if err != nil {
		log.Error(fmt.Sprintf("d.mysql.Query(%s) error(%+v)", _getQusBankByID, err.Error()))
		return
	}
	defer rows.Close()
	for rows.Next() {
		oi := &model.QusBankSt{}
		if err = rows.Scan(&oi.QsBId, &oi.QBName, &oi.CdTime, &oi.MaxRetryTime, &oi.IsDeleted, &oi.TotalCnt); err != nil {
			log.Error(fmt.Sprintf("d.mysql.Query(%s) error(%+v)", _getQusBankByID, err.Error()))
			return
		}
		res = append(res, oi)
	}
	return
}

// CoverStr str
func (d *Dao) CoverStr(strs map[int64]int64) string {
	var buf = bytes.NewBuffer(nil)
	for _, str := range strs {
		buf.WriteString("'")
		buf.WriteString(strconv.FormatInt(str, 10))
		buf.WriteString("'")
		buf.WriteString(",")
	}
	buf.Truncate(buf.Len() - 1)
	return buf.String()
}

// UpdateQsBankCnt 更新数量
func (d *Dao) UpdateQsBankCnt(c context.Context, qid int64) (effid int64, err error) {

	oi := &model.QusBankCnt{}
	err = d.db.QueryRow(c, _getQusBankCnt, qid).Scan(&oi.EasyCnt, &oi.NormalCnt, &oi.HardCnt, &oi.TotalCnt)

	if err != nil {
		log.Error("d.UpdateQsBankCnt error(%v)", err.Error())
		return
	}
	res, err := d.db.Exec(c, _updateQusBankCnt, &oi.EasyCnt, &oi.NormalCnt, &oi.HardCnt, &oi.TotalCnt, qid)
	if err != nil {
		log.Error("d._updateQusBankCnt(qbid:%d, dmid:%d) error(%v)", err)
		return
	}
	return res.RowsAffected()
}
