package dao

import (
	"context"
	"encoding/json"
	"fmt"

	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_fetchItemList      = "SELECT `items_id`,`name`,`ip_right_id`,`brief`,`img` FROM `items` WHERE  `status` = 1 AND `sub_status` <> 13 AND `ip_right_id` <> 0 AND `is_lastest_version` = 1 limit %d,%d"
	_fetchItem          = "SELECT `items_id`,`name`,`ip_right_id`,`brief`,`img` FROM `items` WHERE  `status` = 1 AND `sub_status` <> 13  AND `is_lastest_version` = 1 and `items_id` = ?"
	_fetchIPRight       = "SELECT `name`,`chs_name`,`alias`,`parent_id` FROM `ip_right` WHERE `ip_right_id` = ?"
	_fetchParentIPRight = "SELECT `name`,`chs_name`,`alias` FROM `ip_right` WHERE `ip_right_id` = ?"

	_getWishCount    = "SELECT COUNT(*) FROM ugc_subject_vote_%d WHERE `subject_id` = ? AND `subject_type` = 1 AND DATE(`ctime`) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)"
	_getCommentCount = "SELECT COUNT(*) FROM ugc_info_%d WHERE `subject_id` = ? AND `subject_type` = 1  AND DATE(`ctime`) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY)"
	_getSaleCount    = "select sum(sku_num) from (select order_id,payment_id,payment_time from order_basic where status in (3, 4, 5, 7) and parent_order_id = 0) as t1 join (select order_id,items_id,sku_id, sku_num from order_sku) as t2 on t1.order_id = t2.order_id  where items_id = ? AND DATE(`ctime`) >= DATE_SUB(CURDATE(), INTERVAL 7 DAY) group by items_id "

	_fetchMatch  = "SELECT `season_id`,`items_id`,`sort` FROM `sug_filter` where type = 1 and sort > 1000000000 limit %d,%d"
	_insertMatch = "INSERT INTO sug_filter (`items_id`,`season_id`,`season_name`,`sort`,`type`,`items_name`,`head_pic`,`sug_pic`) VALUES(?, ?, ?,?,'1',?,?,?)ON DUPLICATE KEY UPDATE items_name = IF(type=0, items_name,?) ,season_name =  IF(type=0, season_name,?) ,sort = IF(sort>1000000000, sort,?),mtime = IF(type=0, mtime,NOW()),head_pic =  IF(type=0, head_pic,?),sug_pic =  IF(type=0, sug_pic,?)"
	_updatePic   = "update sug_filter set sug_pic = ?,items_name = ?,head_pic =? where items_id = ? and season_id = ?"
)

//FetchItem query project list from mysql
func (d *Dao) FetchItem(c context.Context) (itemList []*model.Item, err error) {
	var (
		rows   *sql.Rows
		start  = 0
		offset = 10
	)
	for {
		sqlStr := fmt.Sprintf(_fetchItemList, start, offset)
		start += offset
		if rows, err = d.mallDB.Query(c, sqlStr); err != nil {
			log.Error("d._fetchProjectListSQL.Query error(%v)", err)
			return
		}
		i := 0
		for rows.Next() {
			i++
			item := new(model.Item)
			var imgList string
			if err = rows.Scan(&item.ID, &item.Name, &item.IPRightID, &item.Brief, &imgList); err != nil {
				log.Error("row.Scan() error(%v)", err)
				continue
			}
			if item.Brief == "" {
				item.Brief = item.Name
				item.Name = ""
			}
			var imgListArr []string
			if err = json.Unmarshal([]byte(imgList), &imgListArr); err != nil {
				log.Error("get first img err[%s] (%v)", imgList, err)
				err = nil
				continue
			}
			item.HeadImg = imgListArr[0]
			row := d.mallDB.QueryRow(c, _fetchIPRight, item.IPRightID)
			var name, chs_name, alias, parentID, parentName, parentChsName, parentAlias string
			if err = row.Scan(&name, &chs_name, &alias, &parentID); err != nil {
				err = nil
				continue
			}
			if parentID != "0" {
				row := d.mallDB.QueryRow(c, _fetchParentIPRight, parentID)
				if err = row.Scan(&parentName, &parentChsName, &parentAlias); err != nil {
					log.Error("get parent ip_right_info err (%v)", err)
					err = nil
					continue
				}
			}
			item.Keywords = name + " " + chs_name + " " + alias + " " + parentName + " " + parentChsName + " " + parentAlias
			itemList = append(itemList, item)
		}
		if i < 10 {
			break
		}
	}
	defer rows.Close()
	err = rows.Err()
	return
}

// GetBind select sug from db
func (d *Dao) GetBind(c context.Context) (sugList []*model.Sug, err error) {
	var (
		rows   *sql.Rows
		start  = 0
		offset = 50
	)
	for {
		sqlStr := fmt.Sprintf(_fetchMatch, start, offset)
		start += offset
		if rows, err = d.ticketDB.Query(c, sqlStr); err != nil {
			log.Error("d._fetchSugListSQL.Query error(%v)", err)
			return
		}
		i := 0
		for rows.Next() {
			i++
			sug := &model.Sug{}
			sug.Item = &model.Item{}
			if err = rows.Scan(&sug.SeasonID, &sug.Item.ID, &sug.Score); err != nil {
				log.Error("row.Scan() error(%v)", err)
				continue
			}
			sugList = append(sugList, sug)
		}
		if i < 10 {
			break
		}
	}
	defer rows.Close()
	err = rows.Err()
	return
}

// GetItem get item info
func (d *Dao) GetItem(c context.Context, sug *model.Sug) (err error) {
	var row *sql.Row
	if row = d.mallDB.QueryRow(c, _fetchItem, sug.Item.ID); err != nil {
		log.Error("d._fetchProjectListSQL.Query error(%v)", err)
		return
	}
	var imgList string
	if err = row.Scan(&sug.Item.ID, &sug.Item.Name, &sug.Item.IPRightID, &sug.Item.Brief, &imgList); err != nil {
		log.Error("row.Scan() error(%v)", err)
		return
	}
	if sug.Item.Brief == "" {
		sug.Item.Brief = sug.Item.Name
		sug.Item.Name = ""
	}
	var imgListArr []string
	if err = json.Unmarshal([]byte(imgList), &imgListArr); err != nil {
		log.Error("get first img err[%s] (%v)", imgList, err)
		err = nil
		return
	}
	sug.Item.HeadImg = imgListArr[0]
	row = d.mallDB.QueryRow(c, _fetchIPRight, sug.Item.IPRightID)
	var name, chs_name, alias, parentID, parentName, parentChsName, parentAlias string
	if err = row.Scan(&name, &chs_name, &alias, &parentID); err != nil {
		err = nil
		return
	}
	if parentID != "0" {
		row := d.mallDB.QueryRow(c, _fetchParentIPRight, parentID)
		if err = row.Scan(&parentName, &parentChsName, &parentAlias); err != nil {
			log.Error("get parent ip_right_info err (%v)", err)
			err = nil
			return
		}
	}
	sug.Item.Keywords = name + " " + chs_name + " " + alias + " " + parentName + " " + parentChsName + " " + parentAlias

	return
}
