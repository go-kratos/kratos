package dao

import (
	"context"
	"time"

	"go-common/app/job/openplatform/open-market/model"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	_fetchProjectList = "SELECT `id`,`name`,`start_time`,`type` FROM `project` WHERE `status` = 1 AND `type` IN (1,3,4,5,6,8)"

	//_fetchUserWishFirst = "SELECT `mtime` FROM `user_wish` WHERE `item_id` = ? ORDER BY `mtime` LIMIT 1"
	_fetchUserWishByDay = "SELECT DISTINCT `id` FROM `user_wish` WHERE `item_id` = ? AND `mtime` BETWEEN ? AND ?"

	//_fetchUserFavoriteFirst = "SELECT `mtime` FROM `user_favorite` WHERE `item_id` = ? AND `status` = 1 AND `type` = 1 ORDER BY `mtime` LIMIT 1"
	_fetchUserFavoriteByDay = "SELECT DISTINCT `id` FROM `user_favorite` WHERE `item_id` = ? AND `status` = 1 AND `type` = 1  AND `mtime` BETWEEN ? AND ? "

	_fetchStockIfSoldOut = "SELECT `sku_id` FROM `sku_stock` WHERE `item_id` = ? AND `stock` <> 0  LIMIT 1"
	//_fetchStockSoldOutDay = "SELECT `mtime` FROM `sku_stock` WHERE `item_id` = ? AND `stock` = 0 ORDER BY `mtime` DESC LIMIT 1"
)

//FetchProject query project list from mysql
func (d *Dao) FetchProject(c context.Context) (projectList []*model.Project, err error) {
	var rows *sql.Rows
	if rows, err = d.ticketDB.Query(c, _fetchProjectList); err != nil {
		log.Error("d._fetchProjectListSQL.Query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		project := new(model.Project)
		if err = rows.Scan(&project.ID, &project.Name, &project.StartTime, &project.Type); err != nil {
			log.Error("row.Scan() error(%v)", err)
			projectList = nil
			return
		}
		//if project.ID == 11644 || project.ID == 11678 || project.ID == 11643 || project.ID == 11650 {
		//	continue
		//}
		if soldOut := d.CheckStock(context.TODO(), project.ID); soldOut {
			continue
		}
		projectList = append(projectList, project)
	}
	err = rows.Err()
	return
}

//CheckStock if a project sold out,return date
func (d *Dao) CheckStock(c context.Context, projectID int32) (soldOut bool) {
	row := d.ticketDB.QueryRow(c, _fetchStockIfSoldOut, projectID)
	var skuID int
	if err := row.Scan(&skuID); err != nil {
		if err == sql.ErrNoRows {
			return true
		}
		log.Error("row.Scan() error(%v)", err)
	}
	return
}

//WishData fetch wishcount by project and date
func (d *Dao) WishData(c context.Context, projectID int32, startTimeUnix int64) (wishData map[int32]int64, err error) {
	var (
		startTime  string
		firstTime  string
		firstDay   time.Time
		startDay   time.Time
		daysBefore = -1
	)
	wishData = make(map[int32]int64)
	startTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Format(_dateFormat)
	firstTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Add(time.Hour * 24 * -30).Format(_dateFormat)
	startDay, _ = time.Parse(_dateFormat, startTime)
	firstDay, _ = time.Parse(_dateFormat, firstTime)
	for {
		daysBefore++
		startDay = startDay.Add(time.Hour * -24)
		if !(startDay.Before(firstDay)) {
			var rows *sql.Rows
			count := 0
			if rows, err = d.ticketDB.Query(c, _fetchUserWishByDay, projectID, startDay.Format(_dateFormat), startDay.Add(time.Hour*24).Format(_dateFormat)); err != nil {
				log.Error("wish count query error(%v)", err)
				return
			}
			for rows.Next() {
				count++
			}
			err = rows.Err()
			wishData[int32(daysBefore)] = int64(count)
			rows.Close()
			continue
		}
		break
	}
	return
}

//FavoriteData fetch favoritecount by project and date
func (d *Dao) FavoriteData(c context.Context, projectID int32, startTimeUnix int64) (favoriteData map[int32]int64, err error) {
	var (
		startTime  string
		firstTime  string
		firstDay   time.Time
		startDay   time.Time
		daysBefore = -1
	)
	favoriteData = make(map[int32]int64)
	startTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Format(_dateFormat)
	firstTime = time.Unix(startTimeUnix, 0).Add(time.Hour * 24).Add(time.Hour * 24 * -30).Format(_dateFormat)
	startDay, _ = time.Parse(_dateFormat, startTime)
	firstDay, _ = time.Parse(_dateFormat, firstTime)
	for {
		daysBefore++
		startDay = startDay.Add(time.Hour * -24)
		if !(startDay.Before(firstDay)) {
			var rows *sql.Rows
			count := 0
			if rows, err = d.ticketDB.Query(c, _fetchUserFavoriteByDay, projectID, startDay.Format(_dateFormat), startDay.Add(time.Hour*24).Format(_dateFormat)); err != nil {
				log.Error("wish count query error(%v)", err)
				return
			}
			for rows.Next() {
				count++
			}
			err = rows.Err()
			favoriteData[int32(daysBefore)] = int64(count)
			rows.Close()
			continue
		}
		break
	}
	return
}
