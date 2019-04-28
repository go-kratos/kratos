package dao

import (
	"context"
	"fmt"

	"encoding/json"
	"go-common/app/admin/openplatform/sug/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"time"
)

const (
	_getItem     = "SELECT `items_id`,`name`,`brief`,`img` FROM `items` WHERE `items_id` = ? AND `is_lastest_version` = 1"
	_insertMatch = "INSERT INTO `sug_filter` (`season_id`,`items_id`,`type`,`sort`,`season_name`,`items_name`,`head_pic`,`sug_pic`) VALUE (?,?,?,?,?,?,?,?)"
	_updateMatch = "UPDATE `sug_filter` SET `type` = ?,`season_name`=?,`items_name`=?,`sort`=?,`sug_pic`=? WHERE `season_id` = ? AND `items_id` = ?"
	_matchExist  = "SELECT `type` FROM `sug_filter` WHERE `season_id` = ? AND `items_id` = ?"
)

var _selectMatch = "SELECT `season_id`,`items_id`,`sort`,`season_name`,`items_name`,`head_pic`,`sug_pic` FROM `sug_filter` where type = 1"

// GetItem get mall items from db.
func (d *Dao) GetItem(c context.Context, itemsID int64) (item model.Item, err error) {
	row := d.dbMall.QueryRow(c, _getItem, itemsID)
	if err = row.Scan(&item.ItemsID, &item.Name, &item.Brief, &item.Img); err != nil {
		log.Error("Get item %d error,err(%v)", itemsID, err)
		return
	}
	if item.Brief == "" {
		item.Brief = item.Name
		item.Name = ""
	}
	var imgListArr []string
	if err = json.Unmarshal([]byte(item.Img), &imgListArr); err != nil {
		log.Error("get first img err[%s] (%v)", item.Img, err)
		return
	}
	item.Img = imgListArr[0]
	return
}

// UpdateMatch update match.
func (d *Dao) UpdateMatch(c context.Context, season model.Season, item model.Item, typeInt int8, sugPic string) (affect int64, err error) {
	res, err := d.dbTicket.Exec(c, _updateMatch, typeInt, season.Title, item.Name, time.Now().Unix(), sugPic, season.ID, item.ItemsID)
	if err != nil {
		log.Error("row.exec error(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// InsertMatch insert match.
func (d *Dao) InsertMatch(c context.Context, season model.Season, item model.Item, typeInt int8, sort int64, location string) (affect int64, err error) {
	res, err := d.dbTicket.Exec(c, _insertMatch, season.ID, item.ItemsID, typeInt, sort, season.Title, item.Name, item.Img, location)
	if err != nil {
		log.Error("row.exec error(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// GetMatchType get match type.
func (d *Dao) GetMatchType(c context.Context, seasonID, itemsID int64) (matchType int8, err error) {
	row := d.dbTicket.QueryRow(c, _matchExist, seasonID, itemsID)
	if err = row.Scan(&matchType); err != nil {
		if err == sql.ErrNoRows {
			return -1, nil
		}
		return
	}
	return
}

// SearchV2 search match from db
func (d *Dao) SearchV2(c context.Context, params *model.Search) (list []model.SugList, err error) {
	var rows *sql.Rows
	sqlStr := _selectMatch
	if params.ItemsID > 0 {
		sqlStr += fmt.Sprintf(" and items_id = %d", params.ItemsID)
	}
	if params.SeasonID > 0 {
		sqlStr += fmt.Sprintf(" and season_id = %d", params.SeasonID)
	}
	if rows, err = d.dbTicket.Query(c, sqlStr); err != nil {
		log.Error("d._selectMatchSQL.Query error(%v)", err)
		return
	}
	for rows.Next() {
		sug := new(model.SugList)
		if err = rows.Scan(&sug.SeasonId, &sug.ItemsID, &sug.Score, &sug.SeasonName, &sug.ItemsName, &sug.PicURL, &sug.SugURL); err != nil {
			log.Error("row.Scan() error(%v)", err)
			continue
		}
		list = append(list, *sug)
	}
	return
}
