package dao

import (
	"context"
	"fmt"
	"go-common/app/admin/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_addAwardSQL = `INSERT INTO special_award
		(award_id,award_name,cycle_start,cycle_end,announce_date,display_status,total_quota,total_bonus,open_status,open_time,created_by)
		VALUES (%d,'%s','%s','%s','%s',?,?,?,1,'%s','%s')`
	_updateAwardSQL          = "UPDATE special_award SET %s WHERE award_id=?"
	_queryAwardSQL           = `SELECT id,award_id,award_name,cycle_start,cycle_end,announce_date,display_status,total_quota,total_bonus,open_status,open_time,created_by,ctime FROM special_award`
	_pageListAwardSQL        = _queryAwardSQL + " ORDER BY id LIMIT ?,?"
	_selectAwardSQL          = _queryAwardSQL + " WHERE award_id=?"
	_selectAwardForUpdateSQL = _queryAwardSQL + " WHERE award_id=? FOR UPDATE"

	_listDivisionSQL     = "SELECT award_id,division_id,division_name,tag_id FROM special_award_division WHERE %s AND is_deleted=0"
	_listPrizeSQL        = "SELECT award_id,prize_id,bonus,quota FROM special_award_prize WHERE award_id = ? AND is_deleted=0"
	_listWinnerSQL       = "SELECT mid,division_id,prize_id,tag_id FROM special_award_winner WHERE award_id = ? AND is_deleted=0 %s"
	_listRecordSQL       = "SELECT award_id,mid,tag_id FROM special_award_record WHERE award_id = ? AND is_deleted=0 %s"
	_listResourceSQL     = "SELECT resource_type,resource_index,content FROM special_award_resource WHERE award_id=? AND is_deleted=0"
	_groupCountWinnerSQL = "SELECT award_id, count(distinct mid) from special_award_winner where %s AND is_deleted=0 GROUP BY award_id"
	_countWinnerSQL      = "SELECT COUNT(mid) from special_award_winner WHERE award_id = ? %s AND is_deleted=0"
	_countAwardSQL       = "SELECT COUNT(distinct award_id) from special_award"
)

// AddAward tx
func AddAward(tx *sql.Tx, awardID int64, awardName, cycleStart, cycleEnd, announceDate, openTime string,
	displayStatus, totalWinner, totalBonus int, creater string) (int64, error) {
	res, err := tx.Exec(fmt.Sprintf(_addAwardSQL, awardID, awardName, cycleStart, cycleEnd, announceDate, openTime, creater),
		displayStatus, totalWinner, totalBonus)
	if err != nil {
		log.Error("dao#AddAward, tx.Exec err(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// Award query
func (d *Dao) Award(c context.Context, awardID int64) (data *model.Award, err error) {
	data = &model.Award{}
	err = d.rddb.QueryRow(c, _selectAwardSQL, awardID).Scan(
		&data.ID, &data.AwardID, &data.AwardName, &data.CycleStart, &data.CycleEnd, &data.AnnounceDate,
		&data.DisplayStatus, &data.TotalQuota, &data.TotalBonus,
		&data.OpenStatus, &data.OpenTime, &data.CreatedBy, &data.CTime)
	if err != nil {
		if err == sql.ErrNoRows {
			data = nil
			err = nil
			return
		}
		log.Error("dao.Award, d.rddb.QueryRow err(%v)", err)
	}
	return
}

// SelectAwardForUpdate tx
func SelectAwardForUpdate(tx *sql.Tx, awardID int64) (award *model.Award, err error) {
	award = &model.Award{}
	err = tx.QueryRow(_selectAwardForUpdateSQL, awardID).Scan(
		&award.ID, &award.AwardID, &award.AwardName, &award.CycleStart, &award.CycleEnd, &award.AnnounceDate,
		&award.DisplayStatus, &award.TotalQuota, &award.TotalBonus,
		&award.OpenStatus, &award.OpenTime, &award.CreatedBy, &award.CTime)
	if err != nil {
		if err == sql.ErrNoRows {
			err = nil
			award = nil
			return
		}
		log.Error("dao.SelectAwardForUpdate awardID{%d} err(%v)", awardID, err)
	}
	return
}

// UpdateAward tx
func UpdateAward(tx *sql.Tx, awardID int64, values string) (int64, error) {
	if values == "" {
		return 0, nil
	}
	res, err := tx.Exec(fmt.Sprintf(_updateAwardSQL, values), awardID)
	if err != nil {
		log.Error("dao#UpdateAward, tx.Exec err(%v)", err)
		return 0, err
	}
	return res.RowsAffected()
}

// Award tx query
func Award(tx *sql.Tx, awardID int64) (data *model.Award, err error) {
	data = &model.Award{}
	err = tx.QueryRow(_selectAwardSQL, awardID).Scan(&data.ID,
		&data.AwardID, &data.AwardName, &data.CycleStart, &data.CycleEnd, &data.AnnounceDate,
		&data.DisplayStatus, &data.TotalQuota, &data.TotalBonus,
		&data.OpenStatus, &data.OpenTime, &data.CreatedBy, &data.CTime)
	if err != nil {
		if err == sql.ErrNoRows {
			data = nil
			err = nil
			return
		}
		log.Error("dao.Award, tx.QueryRow err(%v)", err)
	}
	return
}

// ListAwardsDivision tx
func ListAwardsDivision(tx *sql.Tx, where string) (res []*model.AwardDivision, err error) {
	res = make([]*model.AwardDivision, 0)
	rows, err := tx.Query(fmt.Sprintf(_listDivisionSQL, where))
	if err != nil {
		log.Error("dao.ListAwardsDivision tx.Query where{%s} err(%v)", where, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardDivision{}
		if err = rows.Scan(&data.AwardID, &data.DivisionID, &data.DivisionName, &data.TagID); err != nil {
			log.Error("dao.ListAwardsDivision rows.Scan where{%s} err(%v)", where, err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// ListDivision tx
func ListDivision(tx *sql.Tx, awardID int64) (res []*model.AwardDivision, err error) {
	return ListAwardsDivision(tx, fmt.Sprintf("award_id=%d", awardID))
}

// DivisionInfo tx query divisionID-to-divisionModel k-v pairs
func DivisionInfo(tx *sql.Tx, awardID int64) (res map[int64]*model.AwardDivision, err error) {
	res = make(map[int64]*model.AwardDivision)
	rows, err := tx.Query(fmt.Sprintf(_listDivisionSQL, fmt.Sprintf("award_id=%d", awardID)))
	if err != nil {
		log.Error("dao.DivisionInfo tx.Query awardID{%d} err(%v)", awardID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardDivision{}
		if err = rows.Scan(&data.AwardID, &data.DivisionID, &data.DivisionName, &data.TagID); err != nil {
			log.Error("dao.DivisionInfo rows.Scan awardID{%d} err(%v)", awardID, err)
			return
		}
		res[data.DivisionID] = data
	}
	err = rows.Err()
	return
}

// ListPrize tx
func ListPrize(tx *sql.Tx, awardID int64) (res []*model.AwardPrize, err error) {
	res = make([]*model.AwardPrize, 0)
	rows, err := tx.Query(_listPrizeSQL, awardID)
	if err != nil {
		log.Error("dao.ListPrize, tx.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardPrize{}
		if err = rows.Scan(&data.AwardID, &data.PrizeID, &data.Bonus, &data.Quota); err != nil {
			log.Error("dao.ListPrize, rows.Scan err(%v)", err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// ListResource tx resource_type->resource_index->resource_content
func ListResource(tx *sql.Tx, awardID int64) (res map[int]map[int]string, err error) {
	res = make(map[int]map[int]string)
	rows, err := tx.Query(_listResourceSQL, awardID)
	if err != nil {
		log.Error("dao.ListResource, tx.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var tp, index int
		var content string
		if err = rows.Scan(&tp, &index, &content); err != nil {
			log.Error("dao.ListPrize, rows.Scan err(%v)", err)
			return
		}
		if _, ok := res[tp]; !ok {
			res[tp] = make(map[int]string)
		}
		res[tp][index] = content
	}
	err = rows.Err()
	return
}

// PrizeInfo tx query prizeID-to-prizeModel k-v pairs
func PrizeInfo(tx *sql.Tx, awardID int64) (res map[int64]*model.AwardPrize, err error) {
	res = make(map[int64]*model.AwardPrize)
	rows, err := tx.Query(_listPrizeSQL, awardID)
	if err != nil {
		log.Error("dao.PrizeInfo, tx.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardPrize{}
		if err = rows.Scan(&data.AwardID, &data.PrizeID, &data.Bonus, &data.Quota); err != nil {
			log.Error("dao.PrizeInfo, rows.Scan err(%v)", err)
			return
		}
		res[data.PrizeID] = data
	}
	err = rows.Err()
	return
}

// CountAward tx
func CountAward(tx *sql.Tx) (total int64, err error) {
	err = tx.QueryRow(_countAwardSQL).Scan(&total)
	if err != nil {
		log.Error("dao.CountAward, tx.Query err(%v)", err)
	}
	return
}

// CountAwardWinner tx
func CountAwardWinner(tx *sql.Tx, awardID int64, where string) (total int64, err error) {
	err = tx.QueryRow(fmt.Sprintf(_countWinnerSQL, where), awardID).Scan(&total)
	if err != nil {
		log.Error("dao.CountAwardWinner, awardID(%d), where(%s) err(%v)", awardID, where, err)
	}
	return
}

// GroupCountAwardWinner tx query awardID-to-winnerCount k-v pairs
func GroupCountAwardWinner(tx *sql.Tx, where string) (res map[int64]int, err error) {
	res = make(map[int64]int)
	rows, err := tx.Query(fmt.Sprintf(_groupCountWinnerSQL, where))
	if err != nil {
		log.Error("dao.GroupCountAwardWinner, tx.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var awardID int64
		var winnerC int
		if err = rows.Scan(&awardID, &winnerC); err != nil {
			log.Error("dao.GroupCountAwardWinner, rows.Scan err(%v)", err)
			return
		}
		res[awardID] = winnerC
	}
	return
}

// AwardWinnerAll tx query winnerModel list
func AwardWinnerAll(tx *sql.Tx, awardID int64) (res []*model.AwardWinner, err error) {
	res = make([]*model.AwardWinner, 0)
	rows, err := tx.Query(fmt.Sprintf(_listWinnerSQL, ""), awardID)
	if err != nil {
		log.Error("dao.AwardWinnerAll, tx.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardWinner{}
		if err = rows.Scan(&data.MID, &data.DivisionID, &data.PrizeID, &data.TagID); err != nil {
			log.Error("dao.AwardWinnerAll, rows.Scan err(%v)", err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// AwardDivisionInfo .
func (d *Dao) AwardDivisionInfo(c context.Context, awardID int64) (res map[int64]*model.AwardDivision, err error) {
	res = make(map[int64]*model.AwardDivision)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_listDivisionSQL, fmt.Sprintf("award_id=%d", awardID)))
	if err != nil {
		log.Error("dao.DivisionInfo tx.Query awardID{%d} err(%v)", awardID, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardDivision{}
		if err = rows.Scan(&data.AwardID, &data.DivisionID, &data.DivisionName, &data.TagID); err != nil {
			log.Error("dao.DivisionInfo rows.Scan awardID{%d} err(%v)", awardID, err)
			return
		}
		res[data.DivisionID] = data
	}
	err = rows.Err()
	return
}

// ListAwardRecord .
func (d *Dao) ListAwardRecord(c context.Context, awardID int64, where string) (res []*model.AwardRecord, err error) {
	res = make([]*model.AwardRecord, 0)
	rows, err := d.rddb.Query(c, fmt.Sprintf(_listRecordSQL, where), awardID)
	if err != nil {
		log.Error("dao.ListAwardRecord, db.Query awardID(%d) where(%s) err(%v)", awardID, where, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardRecord{}
		if err = rows.Scan(&data.AwardID, &data.MID, &data.TagID); err != nil {
			log.Error("dao.ListAwardRecord, rows.Scan awardID(%d) where(%s) err(%v)", awardID, where, err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// QueryAwardWinner tx query mid-to-winnerModel k-v pairs
func QueryAwardWinner(tx *sql.Tx, awardID int64, where string) (res map[int64]*model.AwardWinner, err error) {
	res = make(map[int64]*model.AwardWinner)
	str := fmt.Sprintf(_listWinnerSQL, where)
	rows, err := tx.Query(str, awardID)
	if err != nil {
		log.Error("dao.QueryAwardWinner, tx.Query awardID(%d) sql(%s) err(%v)", awardID, str, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.AwardWinner{}
		if err = rows.Scan(&data.MID, &data.DivisionID, &data.PrizeID, &data.TagID); err != nil {
			log.Error("dao.QueryAwardWinner, rows.Scan awardID(%d) sql(%s) err(%v)", awardID, str, err)
			return
		}
		res[data.MID] = data
	}
	err = rows.Err()
	return
}

// ListAward .
func (d *Dao) ListAward(c context.Context) (res []*model.Award, err error) {
	res = make([]*model.Award, 0)
	rows, err := d.rddb.Query(c, _queryAwardSQL)
	if err != nil {
		log.Error("dao.ListAward, db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.Award{}
		if err = rows.Scan(
			&data.ID,
			&data.AwardID, &data.AwardName, &data.CycleStart, &data.CycleEnd, &data.AnnounceDate,
			&data.DisplayStatus, &data.TotalQuota, &data.TotalBonus,
			&data.OpenStatus, &data.OpenTime, &data.CreatedBy, &data.CTime); err != nil {
			log.Error("dao.ListAward, rows.Scan err(%v)", err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// ListAward tx
func ListAward(tx *sql.Tx, from, limit int) (res []*model.Award, err error) {
	res = make([]*model.Award, 0)
	rows, err := tx.Query(_pageListAwardSQL, from, limit)
	if err != nil {
		log.Error("dao.ListAward, db.Query err(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		data := &model.Award{}
		if err = rows.Scan(&data.ID,
			&data.AwardID, &data.AwardName, &data.CycleStart, &data.CycleEnd, &data.AnnounceDate,
			&data.DisplayStatus, &data.TotalQuota, &data.TotalBonus,
			&data.OpenStatus, &data.OpenTime, &data.CreatedBy, &data.CTime); err != nil {
			log.Error("dao.ListAward, rows.Scan err(%v)", err)
			return
		}
		res = append(res, data)
	}
	err = rows.Err()
	return
}

// DelWinner tx
func DelWinner(tx *sql.Tx, awardID int64, where string) (rows int64, err error) {
	w := fmt.Sprintf("award_id=%d %s", awardID, where)
	return delDataTemplate(tx, "special_award_winner", w, "dao#DelWinner")
}

// DelWinnerAll tx
func DelWinnerAll(tx *sql.Tx, awardID int64) (int64, error) {
	where := fmt.Sprintf("award_id=%d", awardID)
	return delDataTemplate(tx, "special_award_winner", where, "dao#DelWinnerAll")
}

// DelDivisionAll tx
func DelDivisionAll(tx *sql.Tx, awardID int64) (int64, error) {
	where := fmt.Sprintf("award_id=%d", awardID)
	return delDataTemplate(tx, "special_award_division", where, "dao#DelDivisionAll")
}

// DelDivisionsExclude tx
func DelDivisionsExclude(tx *sql.Tx, awardID int64, divisionIDs []int64) (int64, error) {
	where := fmt.Sprintf("award_id=%d AND division_id NOT IN (%s)", awardID, xstr.JoinInts(divisionIDs))
	return delDataTemplate(tx, "special_award_division", where, "dao#DelDivisionsExclude")
}

// DelPrizeAll tx
func DelPrizeAll(tx *sql.Tx, awardID int64) (int64, error) {
	where := fmt.Sprintf("award_id=%d", awardID)
	return delDataTemplate(tx, "special_award_prize", where, "dao#DelPrizeAll")
}

// DelPrizesExclude tx
func DelPrizesExclude(tx *sql.Tx, awardID int64, prizeIDs []int64) (int64, error) {
	where := fmt.Sprintf("award_id=%d AND prize_id NOT IN (%s)", awardID, xstr.JoinInts(prizeIDs))
	return delDataTemplate(tx, "special_award_prize", where, "dao#DelPrizesExclude")
}

// DelResources tx
func DelResources(tx *sql.Tx, where string) (int64, error) {
	return delDataTemplate(tx, "special_award_resource", where, "dao#DelResource")
}

// SaveWinners tx
func SaveWinners(tx *sql.Tx, fields, values string) (rows int64, err error) {
	rows, err = saveDataTemplate(tx, "special_award_winner", fields, values, "dao#SaveWinners")
	return
}

// SaveDivisions tx
func SaveDivisions(tx *sql.Tx, fields, values string) (rows int64, err error) {
	rows, err = saveDataTemplate(tx, "special_award_division", fields, values, "dao#SaveDivision")
	return
}

// SaveResource tx
func SaveResource(tx *sql.Tx, fields, values string) (rows int64, err error) {
	rows, err = saveDataTemplate(tx, "special_award_resource", fields, values, "dao#SaveResource")
	return
}

// SavePrizes tx
func SavePrizes(tx *sql.Tx, fields, values string) (rows int64, err error) {
	rows, err = saveDataTemplate(tx, "special_award_prize", fields, values, "dao#SavePrizes")
	return
}

func saveDataTemplate(tx *sql.Tx, tableName, fields, values, funcName string) (rows int64, err error) {
	str := fmt.Sprintf("INSERT INTO %s(%s) VALUES %s", tableName, fields, values)
	res, err := tx.Exec(str)
	if err != nil {
		log.Error("%s exec(%s) err(%v)", funcName, str, err)
		return
	}
	return res.RowsAffected()
}

func delDataTemplate(tx *sql.Tx, tableName, where, funcName string) (rows int64, err error) {
	str := fmt.Sprintf("UPDATE %s SET is_deleted = 1 WHERE %s", tableName, where)
	res, err := tx.Exec(str)
	if err != nil {
		log.Error("%s exec(%s) err(%v)", funcName, str, err)
		return
	}
	return res.RowsAffected()
}
