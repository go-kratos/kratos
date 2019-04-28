package resource

import (
	"context"
	"fmt"
	"strings"
	"time"

	"database/sql"
	"go-common/app/service/main/resource/model"
	xsql "go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_headerResIds = "142,925,926,927,1576,1580,1584,1588,1592,1596,1600,1604,1608,1612,1616,1620,1622,1634,1920,2210,2260"
)

var (
	_allResSQL    = `SELECT id,platform,name,parent,counter,position,rule,size,preview,description,mark,ctime,mtime,level,type,is_ad FROM resource ORDER BY counter desc,position ASC`
	_allAssignSQL = fmt.Sprintf(`SELECT id,name,contract_id,resource_id,pic,litpic,url,rule,weight,agency,price,atype,username FROM resource_assignment 
		WHERE resource_group_id=0 AND stime<? AND etime>? AND state=0 AND resource_id IN (%s) ORDER BY weight,stime desc`, _headerResIds)
	_allAssignNewSQL = `SELECT ra.id,rm.id,rm.name,ra.contract_id,ra.resource_id,rm.pic,rm.litpic,rm.url,ra.rule,ra.position,
		ra.agency,ra.price,ra.stime,ra.etime,ra.apply_group_id,rm.ctime,rm.mtime,rm.atype,ra.username,rm.player_category FROM resource_assignment AS ra,resource_material AS rm 
		WHERE ra.resource_group_id>0 AND ra.category=0 AND ra.stime<? AND ra.etime>? AND ra.state=0 AND ra.audit_state IN (2,3,4) AND 
		ra.id=rm.resource_assignment_id AND rm.audit_state=2 AND rm.category=0 ORDER BY ra.position ASC,ra.weight DESC,rm.mtime DESC`
	_categoryAssignSQL = fmt.Sprintf(`SELECT ra.id,rm.id,rm.name,ra.contract_id,ra.resource_id,rm.pic,rm.litpic,rm.url,ra.rule,ra.position,ra.agency,ra.price,
		ra.stime,ra.etime,ra.apply_group_id,rm.ctime,rm.mtime,rm.atype,ra.username,rm.player_category FROM resource_assignment AS ra,resource_material AS rm 
		WHERE ra.id=rm.resource_assignment_id AND rm.id IN (SELECT max(rm.id) FROM resource_assignment AS ra,resource_material AS rm WHERE ra.resource_group_id>0 
		AND ra.category=1 AND ra.position_id NOT IN (%s) AND ra.stime<? AND ra.etime>? AND ra.state=0 AND ra.audit_state IN (2,3,4) AND ra.id=rm.resource_assignment_id AND 
		rm.audit_state=2 AND rm.category=1 GROUP BY rm.resource_assignment_id) ORDER BY rand()`, _headerResIds)
	_defBannerSQL = `SELECT id,name,contract_id,resource_id,pic,litpic,url,rule,weight,agency,price,atype,username FROM default_one WHERE state=0`
	// index-icon
	_indexIconSQL = `SELECT id,type,title,state,link,icon,weight,user_name,sttime,endtime,deltime,ctime,mtime FROM icon WHERE state=1 AND deltime=0 AND (type=1 OR (type=2 AND sttime>0))`
	_playIconSQL  = `SELECT icon1,hash1,icon2,hash2,stime FROM bar_icon WHERE stime<? AND etime>? AND is_deleted=0`
	// cmtbox
	_cmtboxSQL = `SELECT id,load_cid,server,port,size_factor,speed_factor,max_onscreen,style,style_param,top_margin,state,ctime,mtime FROM cmtbox WHERE state=1`
	// update resource assignment etime
	_updateResourceAssignmentEtime = `UPDATE resource_assignment SET etime=? WHERE id=?`
	// update resource apply status
	_updateResourceApplyStatus = `UPDATE resource_apply SET audit_state=? WHERE apply_group_id IN (%s)`
	// insert resource logs
	_inResourceLogger = `INSERT INTO resource_logger (uname,uid,module,oid,content) VALUES (?,?,?,?,?)`
)

// Resources get resource infos from db
func (d *Dao) Resources(c context.Context) (rscs []*model.Resource, err error) {
	var size sql.NullString
	rows, err := d.db.Query(c, _allResSQL)
	if err != nil {
		log.Error("d.Resources query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		rsc := &model.Resource{}
		if err = rows.Scan(&rsc.ID, &rsc.Platform, &rsc.Name, &rsc.Parent, &rsc.Counter, &rsc.Position, &rsc.Rule, &size, &rsc.Previce,
			&rsc.Desc, &rsc.Mark, &rsc.CTime, &rsc.MTime, &rsc.Level, &rsc.Type, &rsc.IsAd); err != nil {
			log.Error("Resources rows.Scan err (%v)", err)
			return
		}
		rsc.Size = size.String
		rscs = append(rscs, rsc)
	}
	err = rows.Err()
	return
}

// Assignment get assigment from db
func (d *Dao) Assignment(c context.Context) (asgs []*model.Assignment, err error) {
	rows, err := d.db.Query(c, _allAssignSQL, time.Now(), time.Now())
	if err != nil {
		log.Error("d.Assignment query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		asg := &model.Assignment{}
		if err = rows.Scan(&asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic,
			&asg.URL, &asg.Rule, &asg.Weight, &asg.Agency, &asg.Price, &asg.Atype, &asg.Username); err != nil {
			log.Error("Assignment rows.Scan err (%v)", err)
			return
		}
		asg.AsgID = asg.ID
		asgs = append(asgs, asg)
	}
	err = rows.Err()
	return
}

// AssignmentNew get resource_assigment from new db
func (d *Dao) AssignmentNew(c context.Context) (asgs []*model.Assignment, err error) {
	var (
		ok bool
		pm map[string]string
	)
	rows, err := d.db.Query(c, _allAssignNewSQL, time.Now(), time.Now())
	if err != nil {
		log.Error("d.AssignmentNew query error (%v)", err)
		return
	}
	defer rows.Close()
	pm = make(map[string]string)
	for rows.Next() {
		asg := &model.Assignment{}
		if err = rows.Scan(&asg.AsgID, &asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic,
			&asg.URL, &asg.Rule, &asg.Weight, &asg.Agency, &asg.Price, &asg.STime, &asg.ETime, &asg.ApplyGroupID, &asg.CTime, &asg.MTime, &asg.Atype, &asg.Username, &asg.PlayerCategory); err != nil {
			log.Error("AssignmentNew rows.Scan err (%v)", err)
			return
		}
		if (asg.ResID == 2054) || (asg.ResID == 2055) || (asg.ResID == 2056) ||
			(asg.ResID == 2073) || (asg.ResID == 2074) || (asg.ResID == 2075) ||
			(asg.ResID == 1671) || (asg.ResID == 1672) || (asg.ResID == 1673) ||
			(asg.ResID == 2315) || (asg.ResID == 2316) || (asg.ResID == 2317) ||
			(asg.ResID == 2489) || (asg.ResID == 2490) || (asg.ResID == 2491) ||
			(asg.ResID == 2459) || (asg.ResID == 2460) || (asg.ResID == 2461) ||
			(asg.ResID == 2469) || (asg.ResID == 2470) || (asg.ResID == 2471) ||
			(asg.ResID == 2479) || (asg.ResID == 2480) || (asg.ResID == 2481) ||
			(asg.ResID == 2499) || (asg.ResID == 2500) || (asg.ResID == 2501) ||
			(asg.ResID == 2606) || (asg.ResID == 2607) || (asg.ResID == 2608) || (asg.ResID == 2609) || (asg.ResID == 2610) ||
			(asg.ResID == 2618) || (asg.ResID == 2619) || (asg.ResID == 2620) || (asg.ResID == 2621) || (asg.ResID == 2622) || (asg.ResID == 2623) ||
			(asg.ResID == 2556) || (asg.ResID == 2557) || (asg.ResID == 2558) || (asg.ResID == 2559) || (asg.ResID == 2560) ||
			(asg.ResID == 2991) || (asg.ResID == 2992) || (asg.ResID == 2993) {
			asg.ContractID = "rec_video"
		}
		pindex := fmt.Sprintf("%d_%d", asg.ResID, asg.Weight)
		if _, ok = pm[pindex]; ok {
			continue
		}
		asgs = append(asgs, asg)
		pm[pindex] = pindex
	}
	err = rows.Err()
	return
}

// CategoryAssignment get recommend resource_assigment from db
func (d *Dao) CategoryAssignment(c context.Context) (asgs []*model.Assignment, err error) {
	rows, err := d.db.Query(c, _categoryAssignSQL, time.Now(), time.Now())
	if err != nil {
		log.Error("d.CategoryAssignment query error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		asg := &model.Assignment{}
		if err = rows.Scan(&asg.AsgID, &asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic,
			&asg.URL, &asg.Rule, &asg.Weight, &asg.Agency, &asg.Price, &asg.STime, &asg.ETime, &asg.ApplyGroupID, &asg.CTime, &asg.MTime, &asg.Atype, &asg.Username, &asg.PlayerCategory); err != nil {
			log.Error("CategoryAssignment rows.Scan err (%v)", err)
			return
		}
		if (asg.ResID == 2048) || (asg.ResID == 2066) || (asg.ResID == 1670) || (asg.ResID == 2308) || (asg.ResID == 2521) || (asg.ResID == 2979) {
			asg.ContractID = "rec_video"
		}
		asgs = append(asgs, asg)
	}
	err = rows.Err()
	return
}

// DefaultBanner get default banner info
func (d *Dao) DefaultBanner(c context.Context) (asg *model.Assignment, err error) {
	row := d.db.QueryRow(c, _defBannerSQL)
	asg = &model.Assignment{}
	if err = row.Scan(&asg.ID, &asg.Name, &asg.ContractID, &asg.ResID, &asg.Pic, &asg.LitPic,
		&asg.URL, &asg.Rule, &asg.Weight, &asg.Agency, &asg.Price, &asg.Atype, &asg.Username); err != nil {
		if err == sql.ErrNoRows {
			asg = nil
			err = nil
		} else {
			log.Error("d.DefaultBanner.Scan error(%v)", err)
		}
	}
	return
}

// IndexIcon get index icon.
func (d *Dao) IndexIcon(c context.Context) (icons map[int][]*model.IndexIcon, err error) {
	rows, err := d.db.Query(c, _indexIconSQL)
	if err != nil {
		log.Error("d.IndexIcon query error (%v)", err)
		return
	}
	defer rows.Close()
	icons = make(map[int][]*model.IndexIcon)
	for rows.Next() {
		var link string
		icon := &model.IndexIcon{}
		if err = rows.Scan(&icon.ID, &icon.Type, &icon.Title, &icon.State, &link, &icon.Icon,
			&icon.Weight, &icon.UserName, &icon.StTime, &icon.EndTime, &icon.DelTime, &icon.CTime, &icon.MTime); err != nil {
			log.Error("IndexIcon rows.Scan err (%v)", err)
			return
		}
		icon.Links = strings.Split(link, ",")
		icons[icon.Type] = append(icons[icon.Type], icon)
	}
	err = rows.Err()
	return
}

// PlayerIcon get play icon
func (d *Dao) PlayerIcon(c context.Context) (re *model.PlayerIcon, err error) {
	row := d.db.QueryRow(c, _playIconSQL, time.Now(), time.Now())
	re = &model.PlayerIcon{}
	if err = row.Scan(&re.URL1, &re.Hash1, &re.URL2, &re.Hash2, &re.CTime); err != nil {
		if err == sql.ErrNoRows {
			re = nil
			err = nil
		} else {
			log.Error("d.PlayerIcon.Scan error(%v)", err)
		}
	}
	return
}

// Cmtbox sql live danmaku box
func (d *Dao) Cmtbox(c context.Context) (res map[int64]*model.Cmtbox, err error) {
	rows, err := d.db.Query(c, _cmtboxSQL)
	if err != nil {
		log.Error("d.db.Query error (%v)", err)
		return
	}
	defer rows.Close()
	res = make(map[int64]*model.Cmtbox)
	for rows.Next() {
		re := &model.Cmtbox{}
		if err = rows.Scan(&re.ID, &re.LoadCID, &re.Server, &re.Port, &re.SizeFactor, &re.SpeedFactor, &re.MaxOnscreen,
			&re.Style, &re.StyleParam, &re.TopMargin, &re.State, &re.CTime, &re.MTime); err != nil {
			log.Error("Cmtbox rows.Scan err (%v)", err)
			return
		}
		res[re.ID] = re
	}
	err = rows.Err()
	return
}

// TxOffLine off line resource
func (d *Dao) TxOffLine(tx *xsql.Tx, id int) (row int64, err error) {
	res, err := tx.Exec(_updateResourceAssignmentEtime, time.Now(), id)
	if err != nil {
		log.Error("TxOffLine tx.Exec() error(%v)", err)
		return
	}
	row, err = res.RowsAffected()
	return
}

// TxFreeApply free apply
func (d *Dao) TxFreeApply(tx *xsql.Tx, ids []string) (row int64, err error) {
	res, err := tx.Exec(fmt.Sprintf(_updateResourceApplyStatus, strings.Join(ids, ",")), model.ApplyNoAssignment)
	if err != nil {
		log.Error("TxFreeApply tx.Exec() error(%v)", err)
		return
	}
	row, err = res.RowsAffected()
	return
}

// TxInResourceLogger add resource log
func (d *Dao) TxInResourceLogger(tx *xsql.Tx, module, content string, oid int) (row int64, err error) {
	res, err := tx.Exec(_inResourceLogger, "rejob", 1203, module, oid, content)
	if err != nil {
		log.Error("TxInResourceLogger tx.Exec() error(%v)", err)
		return
	}
	row, err = res.RowsAffected()
	return
}
