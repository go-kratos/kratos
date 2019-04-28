package dao

import (
	"context"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"go-common/app/job/main/growup/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_noCommonTagSQL       = "SELECT tag_info.id, tag_info.category_id, tag_info.ratio, tag_up_info.mid FROM tag_info, tag_up_info WHERE tag_info.is_deleted = 0 AND tag_up_info.is_deleted = 0 AND tag_info.start_at <= ? AND tag_info.end_at >= ? AND tag_info.id = tag_up_info.tag_id AND tag_info.is_common = 0 AND activity_id = 0"
	_commonTagSQL         = "SELECT id, category_id, ratio FROM tag_info WHERE is_deleted = 0 AND tag_info.start_at <= ? AND tag_info.end_at >= ? AND is_common = 1 AND activity_id = 0"
	_activityTagSQL       = "SELECT id, activity_id, category_id, business_id, ratio, is_common FROM tag_info WHERE start_at <= ? AND end_at >= ? AND activity_id > 0 AND is_deleted = 0"
	_getTagTotalIncomeSQL = "SELECT id, tag, total_income FROM tag_info WHERE id IN (%s)"

	_getTagAvSQL         = "SELECT id,av_id,inc_charge,is_deleted FROM av_daily_charge_%s WHERE tag_id = ? AND date = ? AND mid = ? LIMIT ?,?"
	_commonAvSQL         = "SELECT id,av_id,inc_charge,is_deleted FROM av_daily_charge_%s WHERE tag_id = ? AND date = ? LIMIT ?,?"
	_activityMIDExistSQL = "SELECT tag_id, mid FROM tag_up_info WHERE tag_id = ? AND mid = ? AND is_deleted = 0"

	// delete
	_delAvRatioSQL     = "DELETE FROM av_charge_ratio LIMIT ?"
	_delUpIncomeSQL    = "DELETE FROM up_tag_income LIMIT ?"
	_deleteActivitySQL = "DELETE FROM activity_info LIMIT ?"

	// insert
	_addAvTagCaSQL = "INSERT INTO av_charge_ratio(tag_id, av_id, ratio) VALUES %s"

	// update
	_updateTagInfoSQL   = "UPDATE tag_info SET %s WHERE id = ?"
	_updateTagUpInfoSQL = "INSERT INTO tag_up_info(tag_id, mid, total_income) VALUES %s ON DUPLICATE KEY UPDATE total_income = total_income + values(total_income);"
)

// DelAvRatio clear av_charge_ratio table
func (d *Dao) DelAvRatio(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delAvRatioSQL, limit)
	if err != nil {
		log.Error("d.dao.CleanAvRatio error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelIncome clear up_tag_income table
func (d *Dao) DelIncome(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _delUpIncomeSQL, limit)
	if err != nil {
		log.Error("d.dao.CleanIncome error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// DelActivity clear activity_info table.
func (d *Dao) DelActivity(c context.Context, limit int64) (rows int64, err error) {
	res, err := d.db.Exec(c, _deleteActivitySQL, limit)
	if err != nil {
		log.Error("growup-job dao.DelActivity exec error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// CommonTagInfo get common tag infos
func (d *Dao) CommonTagInfo(c context.Context, startAt time.Time) (tagInfos []*model.TagUpInfo, err error) {
	rows, err := d.db.Query(c, _commonTagSQL, startAt, startAt)
	if err != nil {
		log.Error("dao.CommonTagInfo Query (%v), error(%v)", _commonTagSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tagInfo := &model.TagUpInfo{}
		err = rows.Scan(&tagInfo.TagID, &tagInfo.Category, &tagInfo.Ratio)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		tagInfos = append(tagInfos, tagInfo)
	}
	return
}

// NoCommonTagInfo get no common tag infos
func (d *Dao) NoCommonTagInfo(c context.Context, startAt time.Time) (tagInfos []*model.TagUpInfo, err error) {
	rows, err := d.db.Query(c, _noCommonTagSQL, startAt, startAt)
	if err != nil {
		log.Error("dao.NoCommonTagInfo Query (%v), error(%v)", _noCommonTagSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		tagInfo := &model.TagUpInfo{}
		err = rows.Scan(&tagInfo.TagID, &tagInfo.Category, &tagInfo.Ratio, &tagInfo.MID)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		tagInfos = append(tagInfos, tagInfo)
	}
	return
}

// AIDsByMID get aids by mid and tagID
func (d *Dao) AIDsByMID(c context.Context, mid int64, tagID int, offset int64, limit int64, month string, date time.Time) (aids []*model.AID, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getTagAvSQL, month), tagID, date, mid, offset, limit)
	if err != nil {
		log.Error("dao.AIDs Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		aid := &model.AID{}
		err = rows.Scan(&aid.ID, &aid.AvID, &aid.IncCharge, &aid.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	return
}

// AIDs get aids by tagID
func (d *Dao) AIDs(c context.Context, tagID int, offset int64, limit int64, month string, date time.Time) (aids []*model.AID, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_commonAvSQL, month), tagID, date, offset, limit)
	if err != nil {
		log.Error("dao.AIDs Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		aid := &model.AID{}
		err = rows.Scan(&aid.ID, &aid.AvID, &aid.IncCharge, &aid.IsDeleted)
		if err != nil {
			log.Error("rows scan error(%v)", err)
			return
		}
		aids = append(aids, aid)
	}
	log.Info("aids.length: %d", len(aids))
	return
}

// InsertRatio insert av charge ratio
func (d *Dao) InsertRatio(c context.Context, values string) (rows int64, err error) {
	res, err := d.db.Exec(c, fmt.Sprintf(_addAvTagCaSQL, values))
	if err != nil {
		log.Error("d.db.Exec (%s),error(%v)", _addAvTagCaSQL, err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateTagInfo update tag_info total_income.
func (d *Dao) TxUpdateTagInfo(tx *sql.Tx, tagID int64, query string) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateTagInfoSQL, query), tagID)
	if err != nil {
		log.Error("dao.TxUpdateTagInfo error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// TxUpdateTagUpInfo update tag_up_info totalIncome.
func (d *Dao) TxUpdateTagUpInfo(tx *sql.Tx, query string) (rows int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateTagUpInfoSQL, query))
	if err != nil {
		log.Error("dao.TxUpdateTagUpInfo error(%v)", err)
		return
	}
	return res.RowsAffected()
}

// ActivityTagInfo get activity_id tag.
func (d *Dao) ActivityTagInfo(c context.Context, date time.Time) (infos []*model.TagUpInfo, err error) {
	infos = make([]*model.TagUpInfo, 0)
	rows, err := d.db.Query(c, _activityTagSQL, date, date)
	if err != nil {
		log.Error("dao. ActivityTagInfo Query(%s) error(%v)", _activityTagSQL, err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.TagUpInfo{}
		if err = rows.Scan(&a.TagID, &a.ActivityID, &a.Category, &a.Business, &a.Ratio, &a.IsCommon); err != nil {
			log.Error("dao. ActiveTagInfo scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	return
}

// GetActivityAVInfo get avid by http.
func (d *Dao) GetActivityAVInfo(c context.Context, pn, ps int, activities []int64) (info []*model.ActivityAVInfo, total int, err error) {
	info = make([]*model.ActivityAVInfo, 0)
	params := url.Values{}
	params.Set("search_type", "manager_arc")
	params.Set("pagesize", strconv.Itoa(ps))
	params.Set("page", strconv.Itoa(pn))
	var ms string
	for _, a := range activities {
		ms += strconv.FormatInt(a, 10)
		ms += ","
	}
	ms = strings.Trim(ms, ",")
	params.Set("mission_id", ms)
	var res struct {
		Code   int                     `json:"code"`
		Result []*model.ActivityAVInfo `json:"result"`
		Total  int                     `json:"total"`
	}

	if err = d.client.Get(c, d.archiveURL, "", params, &res); err != nil {
		log.Error("growup-job GetActivityAVInfo GET error, url(%s) error(%v)", d.archiveURL+"?"+params.Encode(), err)
		return
	}
	if res.Code != 0 {
		log.Error("growup-job GetActivityAVInfo code != 0. res.Code(%d) | url(%s) error(%v)", res.Code, d.archiveURL+"?"+params.Encode(), err)
		err = ecode.GrowupGetActivityError
		return
	}
	info = res.Result
	total = res.Total
	return
}

// ActivityMIDExist activity mid exist.
func (d *Dao) ActivityMIDExist(c context.Context, tagID, mid int64) (info *model.TagInfo, err error) {
	info = new(model.TagInfo)
	row := d.db.QueryRow(c, _activityMIDExistSQL, tagID, mid)
	err = row.Scan(&info.ID, &info.AVID)
	return
}

// GetAllTypes get all types.
func (d *Dao) GetAllTypes(c context.Context) (rmap map[int16]int16, err error) {
	var res struct {
		Code    int                        `json:"code"`
		Message string                     `json:"message"`
		Data    map[int16]*model.TypesInfo `json:"data"`
	}
	if err = d.client.Get(c, d.typeURL, "", nil, &res); err != nil {
		log.Error("dao.GetAllTypes GET error(%v) | typesURI(%s)", err, d.typeURL)
		return
	}
	if res.Code != 0 {
		log.Error("dao.GetAllTypes code != 0. res.Code(%d) | typesURI(%s) res(%v)", res.Code, d.typeURL, res)
		err = ecode.GrowupGetTypeError
		return
	}
	rmap = make(map[int16]int16, len(res.Data))
	for _, v := range res.Data {
		if v.PID != 0 {
			rmap[v.ID] = v.PID
		}
	}
	return
}

// GetTagTotalIncome get tag and totalIncome from tag_info.
func (d *Dao) GetTagTotalIncome(c context.Context, tagIDs []int64) (infos []*model.TagInfo, err error) {
	rows, err := d.db.Query(c, fmt.Sprintf(_getTagTotalIncomeSQL, xstr.JoinInts(tagIDs)))
	if err != nil {
		log.Error("dao.GetTagTotalIncome query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		a := &model.TagInfo{}
		if err = rows.Scan(&a.ID, &a.Tag, &a.TotalIncome); err != nil {
			log.Error("dao.GetTagTotalIncome scan error(%v)", err)
			return
		}
		infos = append(infos, a)
	}
	err = rows.Err()
	return
}
