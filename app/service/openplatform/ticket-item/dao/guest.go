package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/common/openplatform/random"
	item "go-common/app/service/openplatform/ticket-item/api/grpc/v1"
	"go-common/app/service/openplatform/ticket-item/model"
	"go-common/library/database/elastic"
	"go-common/library/ecode"
	"go-common/library/log"
)

// AddGuest 添加新嘉宾
func (d *Dao) AddGuest(c context.Context, info *item.GuestInfoRequest) (ret bool, err error) {
	// 此处guest_id调用订单号生成器
	guestID := random.Uniqid(19)

	// 插入新数据
	if dbErr := d.db.Create(&model.Guest{
		Name:        info.Name,
		GuestImg:    info.GuestImg,
		Description: info.Description,
		GuestID:     guestID,
	}).Error; dbErr != nil {
		log.Error("insertion failed")
		return false, ecode.NothingFound
	}

	return true, nil
}

// UpdateGuest 编辑嘉宾信息
func (d *Dao) UpdateGuest(c context.Context, info *item.GuestInfoRequest) (res bool, err error) {
	// find original guest with id
	var oriGuest model.Guest

	if dbErr := d.db.First(&oriGuest, info.ID).Error; dbErr != nil {
		log.Error("guest(%d) not found", info.ID)
		return false, ecode.NothingFound
	}

	// update guest with new info (using map can update the column with empty string)
	updateErr := d.db.Model(&oriGuest).Updates(
		map[string]interface{}{
			"name":        info.Name,
			"guest_img":   info.GuestImg,
			"description": info.Description,
		}).Error
	if updateErr != nil {
		log.Error("UPDATE FAILED")
		return false, ecode.NotModified
	}

	return true, nil

}

// GuestStatus 更新嘉宾状态
func (d *Dao) GuestStatus(c context.Context, id int64, status int8) (res bool, err error) {
	// find original guest with id
	var oriGuest model.Guest

	if dbErr := d.db.First(&oriGuest, id).Error; dbErr != nil {
		log.Error("guest (%d) not found", id)
		return false, ecode.NothingFound
	}

	// update guest with new info (using map can update the column with empty string)
	updateErr := d.db.Model(&oriGuest).Update("status", status).Error
	if updateErr != nil {
		log.Error("CHANGE STATUS FAILED")
		return false, ecode.NotModified
	}

	return true, nil
}

// GetGuests 获取项目下所有嘉宾
func (d *Dao) GetGuests(c context.Context, pid int64) (res []*model.Guest, err error) {
	// 获取项目下所有guestID
	var projectGuests []*model.ProjectGuest
	if dbErr := d.db.Select("guest_id").Order("position").Where("project_id = ? and delete_status = 0", pid).Find(&projectGuests).Error; dbErr != nil {
		log.Error("获取项目嘉宾id失败: %s", dbErr)
		err = ecode.NothingFound
		return
	}

	var guestIDs []int64
	for idx := range projectGuests {
		guestIDs = append(guestIDs, projectGuests[idx].GuestID)
	}

	// 获取未禁用嘉宾信息 并用id当键值组成map
	var guestInfo []*model.Guest
	if dbErr := d.db.Where("status = 1 and id in (?)", guestIDs).Find(&guestInfo).Error; dbErr != nil {
		log.Error("获取项目嘉宾信息失败: %s", dbErr)
		err = ecode.NothingFound
		return
	}

	guestInfoMap := make(map[int64]*model.Guest)
	for _, value := range guestInfo {
		guestInfoMap[value.ID] = value
	}

	// 根据guestIDs重新排序嘉宾信息
	for _, value := range guestIDs {
		if guestInfoMap[value].Status == 1 {
			res = append(res, guestInfoMap[value])
		}
	}

	return

}

// GuestSearch guest es search.
func (d *Dao) GuestSearch(c context.Context, arg *model.GuestSearchParam) (guests *model.GuestSearchList, err error) {
	r := d.es.NewRequest("ticket_guest").
		Index("ticket_guest").
		Order("ctime", elastic.OrderDesc).
		Order("id", elastic.OrderAsc).
		Ps(arg.Ps).Pn(arg.Pn)

	if arg.Keyword != "" {
		if id, err1 := strconv.Atoi(arg.Keyword); err1 == nil {
			r = r.WhereEq("id", id)
		} else {
			r = r.WhereLike([]string{"name"}, []string{arg.Keyword}, true, elastic.LikeLevelLow)
		}
	}
	log.Info(fmt.Sprintf("%s/x/admin/search/query?%s", d.c.URL.ElasticHost, r.Params()))

	guests = new(model.GuestSearchList)
	err = r.Scan(c, guests)
	if err != nil {
		log.Error("GuestSearch() Scan() error(%v)", err)
	}

	return
}
