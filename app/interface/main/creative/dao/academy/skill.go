package academy

import (
	"context"
	"fmt"

	"go-common/app/interface/main/creative/model/academy"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	// select
	_getOccSQL         = "SELECT id, rank, name, `desc`, main_step, main_software, logo FROM academy_occupation ORDER BY rank ASC"
	_getSkillSQL       = "SELECT id, oid, name, `desc` FROM academy_skill ORDER BY id ASC"
	_getSkillArcSQL    = "SELECT id, aid, type, pid, skid, sid FROM academy_arc_skill WHERE state=0"
	_getSkillArcCntSQL = "SELECT count(*) FROM academy_arc_skill WHERE state=0"
)

//Occupations get all occupation.
func (d *Dao) Occupations(c context.Context) (res []*academy.Occupation, err error) {
	rows, err := d.db.Query(c, _getOccSQL)
	if err != nil {
		log.Error("Occupations d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.Occupation, 0)
	for rows.Next() {
		o := &academy.Occupation{}
		if err = rows.Scan(&o.ID, &o.Rank, &o.Name, &o.Desc, &o.MainStep, &o.MainSoftWare, &o.Logo); err != nil {
			log.Error("Occupations rows.Scan error(%v)", err)
			return
		}
		res = append(res, o)
	}
	return
}

//Skills get all Skill.
func (d *Dao) Skills(c context.Context) (res []*academy.Skill, err error) {
	rows, err := d.db.Query(c, _getSkillSQL)
	if err != nil {
		log.Error("Skills d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.Skill, 0)
	for rows.Next() {
		o := &academy.Skill{}
		if err = rows.Scan(&o.ID, &o.OID, &o.Name, &o.Desc); err != nil {
			log.Error("Skills rows.Scan error(%v)", err)
			return
		}
		res = append(res, o)
	}
	return
}

//SkillArcs get all SkillArc.
func (d *Dao) SkillArcs(c context.Context, pids, skids, sids []int64, offset, limit int) (res []*academy.SkillArc, err error) {
	var (
		whereStr = _getSkillArcSQL
		limiStr  = " ORDER BY id ASC LIMIT ?,?"
		rows     *sql.Rows
	)

	if len(pids) > 0 {
		whereStr += fmt.Sprintf(" AND pid IN (%s)"+limiStr, xstr.JoinInts(pids))
	} else if len(skids) > 0 {
		whereStr += fmt.Sprintf(" AND skid IN (%s)"+limiStr, xstr.JoinInts(skids))
	} else if len(sids) > 0 {
		whereStr += fmt.Sprintf(" AND sid IN (%s)"+limiStr, xstr.JoinInts(sids))
	} else {
		whereStr += limiStr
	}

	rows, err = d.db.Query(c, whereStr, offset, limit)
	if err != nil {
		log.Error("SkillArcs d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	res = make([]*academy.SkillArc, 0)
	for rows.Next() {
		o := &academy.SkillArc{}
		if err = rows.Scan(&o.ID, &o.AID, &o.Type, &o.PID, &o.SkID, &o.SID); err != nil {
			log.Error("SkillArcs rows.Scan error(%v)", err)
			return
		}
		res = append(res, o)
	}

	return
}

//SkillArcCount get all skill achive count.
func (d *Dao) SkillArcCount(c context.Context, pids, skids, sids []int64) (count int, err error) {

	var (
		whereStr = _getSkillArcCntSQL
	)

	if len(pids) > 0 {
		whereStr += fmt.Sprintf(" AND pid IN (%s)", xstr.JoinInts(pids))
	} else if len(skids) > 0 {
		whereStr += fmt.Sprintf(" AND skid IN (%s)", xstr.JoinInts(skids))
	} else if len(sids) > 0 {
		whereStr += fmt.Sprintf(" AND sid IN (%s)", xstr.JoinInts(sids))
	}

	if err = d.db.QueryRow(c, whereStr).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.db.QueryRow error(%v)", err)
	}

	return
}

// PlayAdd  add play archive.
func (d *Dao) PlayAdd(c context.Context, p *academy.Play) (id int64, err error) {
	_inPlaySQL := "INSERT INTO academy_playlist (mid, aid, business, watch, ctime, mtime) VALUES (?, ?, ?, ?, ?,?) ON DUPLICATE KEY UPDATE state=0,watch=?,business=?,mtime=?"

	res, err := d.db.Exec(c, _inPlaySQL, p.MID, p.AID, p.Business, p.Watch, p.CTime, p.MTime, p.Watch, p.Business, p.MTime)
	if err != nil {
		log.Error("PlayAdd d.db.Exec error(%v)", err)
		return
	}
	id, err = res.RowsAffected()
	return
}

// PlayDel add play archive.
func (d *Dao) PlayDel(c context.Context, p *academy.Play) (id int64, err error) {
	_upPlaySQL := "UPDATE academy_playlist SET state=? WHERE mid=? AND aid=? AND business=?"

	res, err := d.db.Exec(c, _upPlaySQL, -1, p.MID, p.AID, p.Business) //-1 删除
	if err != nil {
		log.Error("PlayDel d.db.Exec error(%v)", err)
		return
	}
	id, err = res.RowsAffected()
	return
}

//Plays get all play by mid.
func (d *Dao) Plays(c context.Context, mid int64, offset, limit int) (res []*academy.Play, err error) {
	_getPlaySSQL := "SELECT mid, aid, business, watch, ctime, mtime FROM academy_playlist WHERE state=0 AND mid=? ORDER BY mtime DESC LIMIT ?,?"

	rows, err := d.db.Query(c, _getPlaySSQL, mid, offset, limit)
	if err != nil {
		log.Error("Plays d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()
	res = make([]*academy.Play, 0)

	for rows.Next() {
		o := &academy.Play{}
		if err = rows.Scan(&o.MID, &o.AID, &o.Business, &o.Watch, &o.CTime, &o.MTime); err != nil {
			log.Error("Plays rows.Scan error(%v)", err)
			return
		}
		res = append(res, o)
	}
	return
}

//PlayCount get all play achive count.
func (d *Dao) PlayCount(c context.Context, mid int64) (count int, err error) {
	_getPlayCntSQL := "SELECT count(*) FROM academy_playlist WHERE state=0 AND mid=?"

	if err = d.db.QueryRow(c, _getPlayCntSQL, mid).Scan(&count); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.db.QueryRow error(%v)", err)
	}
	return
}

//Play get play achive info.
func (d *Dao) Play(c context.Context, p *academy.Play) (res *academy.Play, err error) {
	_getPlaySQL := "SELECT mid, aid, business, watch, ctime, mtime FROM academy_playlist WHERE mid=? AND aid=? AND business=? "

	res = &academy.Play{}
	if err = d.db.QueryRow(c, _getPlaySQL, p.MID, p.AID, p.Business).Scan(&res.MID, &res.AID, &res.Business, &res.Watch, &res.CTime, &res.MTime); err != nil {
		if err == sql.ErrNoRows {
			err = nil
			return
		}
		log.Error("d.db.QueryRow error(%v)", err)
	}
	return
}
