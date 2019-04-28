package app

import (
	"context"
	"encoding/json"

	"database/sql"
	"go-common/app/interface/main/creative/model/app"
	"go-common/library/log"
)

const (
	// select app portal
	_getMulSQL          = "SELECT subtitle,whiteexp,id,more,buildexp,build,compare,state,pos,title,platform,icon,url,mark,type,ctime,mtime,ptime FROM app_portal WHERE state=1 AND type=? ORDER BY pos ASC"
	_addMaterialDataSQL = "INSERT INTO archive_material (aid, cid, type, data) VALUES (?,?,?,?)"
	_getMaterialDataSQL = " SELECT data from archive_material  where aid=? and cid=? and type=? limit 1"
)

// Portals for portal config.
func (d *Dao) Portals(c context.Context, ty int) (apt []*app.PortalMeta, err error) {
	rows, err := d.db.Query(c, _getMulSQL, ty)
	if err != nil {
		log.Error("mysqlDB.Query error(%v)", err)
		return
	}
	defer rows.Close()
	apt = make([]*app.PortalMeta, 0)
	for rows.Next() {
		a := &app.PortalMeta{
			BuildExps: make([]*app.BuildExp, 0),
			WhiteExps: make([]*app.WhiteExp, 0),
		}
		var (
			buildExpJSON json.RawMessage
			whiteExpJSON json.RawMessage
		)
		if err = rows.Scan(
			&a.SubTitle,
			&whiteExpJSON,
			&a.ID,
			&a.More,
			&buildExpJSON,
			&a.Build,
			&a.Compare,
			&a.State,
			&a.Pos,
			&a.Title,
			&a.Platform,
			&a.Icon,
			&a.URL,
			&a.Mark,
			&a.Type,
			&a.CTime,
			&a.MTime,
			&a.PTime,
		); err != nil {
			log.Error("row.Scan error(%v)", err)
			return
		}
		if len(buildExpJSON) > 0 {
			var buildExps []*app.BuildExp
			if err = json.Unmarshal([]byte(buildExpJSON), &buildExps); err != nil {
				log.Error("json.Unmarshal buildExps failed error(%v)", err)
				continue
			}
			a.BuildExps = buildExps
		}
		if len(whiteExpJSON) > 0 {
			var whiteExps []*app.WhiteExp
			if err = json.Unmarshal([]byte(whiteExpJSON), &whiteExps); err != nil {
				log.Error("json.Unmarshal whiteExps failed error(%v)", err)
				continue
			}
			a.WhiteExps = whiteExps
		}
		apt = append(apt, a)
	}
	return
}

// AddMaterialData fn
func (d *Dao) AddMaterialData(c context.Context, data *app.EditorData) (err error) {
	_, err = d.db.Exec(c, _addMaterialDataSQL, data.AID, data.CID, data.Type, data.Data)
	if err != nil {
		log.Error("d.db.Exec error(%v)", err)
		return
	}
	return
}

// GetMaterialData fn
func (d *Dao) GetMaterialData(c context.Context, data *app.EditorData) (dt []byte, err error) {
	row := d.db.QueryRow(c, _getMaterialDataSQL, data.AID, data.CID, data.Type)
	if err = row.Scan(&dt); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		} else {
			log.Error("row scan error(%v)", err)
		}
	}
	return
}
