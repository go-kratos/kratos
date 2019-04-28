package dao

import (
	"context"
	xsql "database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/videoup-task/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_reviewCfg    = 3
	_countSQL     = "SELECT COUNT(*) FROM task_json_config"
	_listConfsSQL = "SELECT id,conf_json,conf_type,btime,etime,state,uid,uname,description,mtime FROM task_json_config WHERE conf_type=3"
	_reConfsSQL   = "SELECT id,conf_json,conf_type,btime,etime,state,uid,uname,description,mtime FROM task_json_config WHERE conf_type=3"
	_inConfSQL    = "INSERT INTO task_json_config(conf_json,conf_type,btime,etime,state,uid,uname,description) VALUE (?,?,?,?,?,?,?,?)"
	_upConfSQL    = "UPDATE task_json_config SET conf_json=?,conf_type=?,btime=?,etime=?,state=?,uid=?,uname=?,description=? WHERE id=?"
	_delConfSQL   = "DELETE FROM task_json_config WHERE id=?"

	_reviewSQL   = "SELECT review_form FROM task_review WHERE task_id=?"
	_inReviewSQL = "INSERT INTO task_review(task_id,review_form,uid,uname) VALUE (?,?,?,?)"
)

// ListConfs 配置列表
func (d *Dao) ListConfs(c context.Context, uids []int64, bt, et, sort string, pn, ps int64) (rcs []*model.ReviewConf, count int64, err error) {
	var (
		rows                           *sql.Rows
		countstring, sqlstring, params string
		wherecases                     []string
	)

	if len(uids) > 0 {
		wherecases = append(wherecases, fmt.Sprintf("uid IN (%s)", xstr.JoinInts(uids)))
	}
	if len(bt) > 0 && len(et) > 0 {
		wherecases = append(wherecases, fmt.Sprintf("mtime>='%s' AND mtime<='%s'", bt, et))
	}

	if len(wherecases) > 0 {
		params = " AND " + strings.Join(wherecases, " AND ")
	}
	countstring = _countSQL + " WHERE conf_type=3" + params
	sqlstring = _listConfsSQL + params + fmt.Sprintf(" ORDER BY mtime %s LIMIT %d,%d", sort, (pn-1)*ps, pn*ps)

	if err = d.arcDB.QueryRow(c, countstring).Scan(&count); err != nil {
		log.Error("d.arcDB.QueryRow(%s) error(%v)", countstring, err)
		return
	}
	if count == 0 {
		return
	}

	if rows, err = d.arcDB.Query(c, sqlstring); err != nil {
		log.Error("d.arcDB.Query(%s) error(%v)", sqlstring, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			jsonCfg []byte
			cfgType int8
		)
		trc := &model.ReviewConf{}
		if err = rows.Scan(&trc.ID, &jsonCfg, &cfgType, &trc.Bt, &trc.Et, &trc.State, &trc.UID, &trc.Uname, &trc.Desc, &trc.Mt); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}

		if err = json.Unmarshal(jsonCfg, trc); err != nil {
			log.Error("json.Unmarshal error(%v)", err)
			continue
		}
		trc.Refresh()
		rcs = append(rcs, trc)
	}
	return
}

// ReviewConfs 复审配置
func (d *Dao) ReviewConfs(c context.Context) (rcs []*model.ReviewConf, err error) {
	var rows *sql.Rows

	if rows, err = d.arcDB.Query(c, _reConfsSQL); err != nil {
		log.Error("d.arcDB.Query(%s, %d) error(%v)", _reConfsSQL, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var (
			jsonCfg []byte
			cfgType int8
		)
		trc := &model.ReviewConf{}
		if err = rows.Scan(&trc.ID, &jsonCfg, &cfgType, &trc.Bt, &trc.Et, &trc.State, &trc.UID, &trc.Uname, &trc.Desc, &trc.Mt); err != nil {
			log.Error("rows.Scan error(%v)", err)
			continue
		}

		if err = json.Unmarshal(jsonCfg, trc); err != nil {
			log.Error("json.Unmarshal error(%v)", err)
			continue
		}
		trc.Refresh()
		rcs = append(rcs, trc)
	}
	return
}

// InReviewConf 插入配置
func (d *Dao) InReviewConf(c context.Context, rc *model.ReviewConf) (lastid int64, err error) {
	var (
		res     xsql.Result
		jsonCfg []byte
	)

	v := new(struct {
		Types    []int64 `json:"types" params:"types"`       // 分区
		UpFroms  []int64 `json:"upfroms" params:"upfroms"`   // 投稿来源
		UpGroups []int64 `json:"upgroups" params:"upgroups"` // 用户组
		Uids     []int64 `json:"uids" params:"uids"`         // 指定uid
		FansLow  int64   `json:"fanslow" params:"fanslow"`   // 粉丝数最低值
		FansHigh int64   `json:"fanshigh" params:"fanshigh"` // 粉丝数最高
	})
	v.Types = rc.Types
	v.UpFroms = rc.UpFroms
	v.UpGroups = rc.UpGroups
	v.Uids = rc.Uids
	v.FansLow = rc.FansLow
	v.FansHigh = rc.FansHigh
	if rc.Bt.TimeValue().IsZero() {
		rc.Bt = model.NewFormatTime(time.Now())
	}

	if jsonCfg, err = json.Marshal(v); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", rc, err)
		return
	}

	if res, err = d.arcDB.Exec(c, _inConfSQL, jsonCfg, _reviewCfg, rc.Bt, rc.Et, 0, rc.UID, rc.Uname, rc.Desc); err != nil {
		log.Error("d.arcDB.Exec(%+v) error(%s, %v)", _inConfSQL, rc, err)
		return
	}
	return res.LastInsertId()
}

// UpReviewConf 更新指定配置
func (d *Dao) UpReviewConf(c context.Context, rc *model.ReviewConf) (lastid int64, err error) {
	var (
		res     xsql.Result
		jsonCfg []byte
	)

	v := new(struct {
		Types    []int64 `json:"types" params:"types"`       // 分区
		UpFroms  []int64 `json:"upfroms" params:"upfroms"`   // 投稿来源
		UpGroups []int64 `json:"upgroups" params:"upgroups"` // 用户组
		Uids     []int64 `json:"uids" params:"uids"`         // 指定uid
		FansLow  int64   `json:"fanslow" params:"fanslow"`   // 粉丝数最低值
		FansHigh int64   `json:"fanshigh" params:"fanshigh"` // 粉丝数最高
	})
	v.Types = rc.Types
	v.UpFroms = rc.UpFroms
	v.UpGroups = rc.UpGroups
	v.Uids = rc.Uids
	v.FansLow = rc.FansLow
	v.FansHigh = rc.FansHigh
	if rc.Bt.TimeValue().IsZero() {
		rc.Bt = model.NewFormatTime(time.Now())
	}

	if jsonCfg, err = json.Marshal(v); err != nil {
		log.Error("json.Marshal(%+v) error(%v)", rc, err)
		return
	}

	if res, err = d.arcDB.Exec(c, _upConfSQL, jsonCfg, _reviewCfg, rc.Bt, rc.Et, rc.State, rc.UID, rc.Uname, rc.Desc, rc.ID); err != nil {
		log.Error("d.arcDB.Exec(%s %+v) error(%v)", _upConfSQL, rc, err)
		return
	}
	return res.RowsAffected()
}

// DelReviewConf 删除指定配置
func (d *Dao) DelReviewConf(c context.Context, id int) (lastid int64, err error) {
	var res xsql.Result
	if res, err = d.arcDB.Exec(c, _delConfSQL, id); err != nil {
		log.Error("d.arcDB.Exec(%s %d) error(%v)", _delConfSQL, id, err)
		return
	}
	return res.RowsAffected()
}

// ReviewForm 复审表单
func (d *Dao) ReviewForm(c context.Context, tid int64) (tsf *model.SubmitForm, err error) {
	var form []byte
	if err = d.arcDB.QueryRow(c, _reviewSQL, tid).Scan(&form); err != nil {
		if err == sql.ErrNoRows {
			log.Info("ReviewForm QueryRow empty(%d)", tid)
			err = nil
			return
		}
		log.Error("d.arcDB.QueryRow(%s, %d) error(%v)", _reviewSQL, tid, err)
		return
	}

	tsf = &model.SubmitForm{}
	if err = json.Unmarshal(form, tsf); err != nil {
		log.Error("json.Unmarshal error(%v)", err)
		tsf = nil
	}
	return
}

// InReviewForm insert submit form
func (d *Dao) InReviewForm(c context.Context, sf *model.SubmitForm, uid int64, uname string) (lastid int64, err error) {
	var (
		res xsql.Result
		bsf []byte
	)

	if bsf, err = json.Marshal(sf); err != nil {
		log.Error("json.Marshal error(%v)", err)
		return
	}

	if res, err = d.arcDB.Exec(c, _inReviewSQL, sf.TaskID, bsf, uid, uname); err != nil {
		log.Error("d.arcDB.Exec(%s,%d,%v,%d,%s) error(%v)", _inReviewSQL, sf.TaskID, bsf, uid, uname, err)
		return
	}

	return res.LastInsertId()
}
