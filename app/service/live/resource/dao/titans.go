package dao

import (
	"context"
	"github.com/jinzhu/gorm"
	"go-common/app/service/live/resource/api/http/v1"
	"go-common/app/service/live/resource/model"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "time"
)

var (
	_statusDel = 4
	_statusOn  = 1
	_statusOff = 2
)

// CheckParams 校验key及team
func (d *Dao) CheckParams(c context.Context, team int64, keyword string) (err error) {
	err = nil
	if "" == keyword {
		err = ecode.ResourceParamErr
	}
	return err
}

// TeamKeyword 用于封装业务参数结构
type TeamKeyword struct {
	Team    int64
	Keyword string
}

// SelectByTeamIndex 单个key的value查询
func (d *Dao) SelectByTeamIndex(c context.Context, team int64, keyword string, id int64) (res *model.SundryConfig, err error) {
	res = &model.SundryConfig{}
	resMid := &model.SundyConfigObject{}
	if 0 != id {
		err = d.rsDB.Model(&model.SundyConfigObject{}).Where("id=?", id).Where("status != ?", _statusDel).Find(&resMid).Error
	} else {
		err = d.rsDB.Model(&model.SundyConfigObject{}).Where("team=?", team).Where("keyword=?", keyword).Where("status != ?", _statusDel).Find(&resMid).Error
	}

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("[live.titans.dao| selectByGroupIndex] d.db.Exec err: %v", err)
		return
	}
	res.Id = resMid.Id
	res.Team = resMid.Team
	res.Keyword = resMid.Keyword
	res.Value = resMid.Value
	res.Status = resMid.Status
	res.Ctime = resMid.Ctime.Time().Format("2006-01-02 15:04:05")
	res.Mtime = resMid.Mtime.Time().Format("2006-01-02 15:04:05")

	err = nil
	return
}

// SelectByParams 管理后台通过条件查询
func (d *Dao) SelectByParams(c context.Context, id int64, team int64, keyword string, name string, status int64, page int64, pageSize int64) (res []*model.SundryConfig, count int64, err error) {
	var (
		Items []*model.SundyConfigObject
	)

	condition := d.rsDB

	//模型配置
	if -1 == team {
		condition = condition.Where("team = ?", 0)
	}
	//除模型配置外的全部配置
	if 0 == team {
		condition = condition.Where("team != ?", 0)
	}
	if 0 != team && -1 != team {
		condition = condition.Where("team = ?", team)
	}

	if 0 != id {
		condition = condition.Where("id = ?", id)
	}
	if "" != keyword {
		condition = condition.Where("keyword = ?", keyword)
	}
	if "" != name {
		condition = condition.Where("name like '%" + name + "%'")
	}
	if 0 != status {
		condition = condition.Where("status = ?", status)
	}
	condition = condition.Where("status != ?", _statusDel)
	sModel := condition.Model(&model.SundryConfig{})
	sModel.Count(&count)

	iOffset := (page - 1) * pageSize
	err = sModel.Offset(iOffset).Limit(pageSize).Order("mtime DESC, id DESC").Find(&Items).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("[live.titans.dao| select by params] d.db.Exec err: %v", err)
		return
	}
	for _, v := range Items {
		item := &model.SundryConfig{
			Id:      v.Id,
			Team:    v.Team,
			Keyword: v.Keyword,
			Name:    v.Name,
			Value:   v.Value,
			Ctime:   v.Ctime.Time().Format("2006-01-02 15:04:05"),
			Mtime:   v.Mtime.Time().Format("2006-01-02 15:04:05"),
			Status:  v.Status,
		}
		res = append(res, item)
	}
	return
}

// InsertRecord 管理后台插入/编辑一条记录
func (d *Dao) InsertRecord(c context.Context, team int64, keyword string, value string, name string, status int64, oid int64) (id int64, count int64, err error) {
	//创建模板
	if team == -1 {
		team = 0
	}
	//查询唯一索引
	err = d.rsDB.Model(&model.SundryConfig{}).Where("team=? and keyword=? and id != ?", team, keyword, oid).Count(&count).Error
	if nil != err {
		return
	}
	if 0 != count {
		return
	}
	setMaps := &model.InsertMaps{
		Team:    team,
		Keyword: keyword,
		Value:   value,
		Name:    name,
		Status:  status,
	}
	newRecord := &model.SundryConfig{}

	//编辑
	if oid != 0 {
		err = d.rsDB.Model(&model.SundryConfig{}).Where("id = ?", oid).Update(setMaps).Error
		id = oid
		return
	}
	setMaps.Status = int64(_statusOff)
	err = d.rsDB.Create(setMaps).Error
	if err != nil {
		log.Error("[live.titans.dao| insertRecord] d.db.Exec err: %v", err)
		return
	}
	d.rsDB.Where("keyword = ?", keyword).Find(&newRecord)
	id = newRecord.Id
	return
}

// SelectByLikes 业务方的请求sql
func (d *Dao) SelectByLikes(c context.Context, teams []int64, teamKeys []*TeamKeyword) (Items []*model.SundyConfigObject, err error) {
	Items = []*model.SundyConfigObject{}
	if len(teams) == 0 && len(teamKeys) == 0 {
		return
	}

	/** sql 式执行
	sql := "select * from ap_sundry_config where status = 1"
	if len(teams) != 0 {
		teamStr := ""
		for i, num := range teams {
			teamStr += strconv.Itoa(int(num))
			if i+1 != len(teams) {
				teamStr += ","
			}
		}
		sql += " and team in (" + teamStr + ") "
	}
	if len(teamKeys) != 0 {
		for _, v := range teamKeys {
			sql = sql + " or (team =" + strconv.Itoa(int(v.Team)) + " and keyword = '" + v.Keyword + "') "
		}
	}
	rows, err := d.db.Query(c, sql)
	if err != nil {
		return
	}
	for rows.Next() {
		resMid := &model.SundyConfigObject{}
		rows.Scan(&resMid.Id, &resMid.Team, &resMid.Keyword, &resMid.Name, &resMid.Value, &resMid.Ctime, &resMid.Mtime, &resMid.Status)
		Items = append(Items, resMid)
	}*/
	/** orm **/
	condition := d.rsDB.Where("status = ?", 1)
	if len(teams) != 0 {
		condition = condition.Where("team in (?)", teams)
	}
	if len(teamKeys) != 0 {
		for _, v := range teamKeys {
			condition = condition.Or("team = ? and keyword = ?", v.Team, v.Keyword)
		}
	}
	err = condition.Find(&Items).Error
	return
}

// FormatTime 时间格式化
func (d *Dao) FormatTime(c context.Context, timeStrUtc string) (timeStr string) {
	ctime, _ := xtime.ParseInLocation("2006-01-02T15:04:05+08:00", timeStrUtc, xtime.Local)
	timeStr = ctime.Format("2006-01-02 15:04:05")
	return timeStr
}

// InsertServiceConfig 插入服务配置
func (d *Dao) InsertServiceConfig(c context.Context, oid int64, treeName string, treePath string, treeId int64, service string, keyword string, template int64, name string, value string, status int64) (id int64, err error) {
	checkRes := &model.ServiceConfigObject{}
	err = d.rsDB.Model(&model.ServiceConfigObject{}).Where("tree_id=?", treeId).Where("keyword=?", keyword).Find(&checkRes).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		log.Error("[live.titans.dao| insertServiceRecord] d.db.Exec err: %v", err)
		return
	}
	if 0 != checkRes.Id && oid != checkRes.Id {
		id = -1
		return
	}

	if oid != 0 {
		//编辑
		updateMaps := &model.UpdateServiceConfig{
			Service:  service,
			Keyword:  keyword,
			Template: template,
			Value:    value,
			Name:     name,
			Status:   status,
		}
		err = d.rsDB.Model(&model.ServiceConfigObject{}).Where("id=?", oid).Update(updateMaps).Error
		id = oid
	} else {
		setMaps := &model.InsertServiceConfig{
			TreeName: treeName,
			TreePath: treePath,
			TreeId:   treeId,
			Service:  service,
			Keyword:  keyword,
			Template: template,
			Value:    value,
			Name:     name,
			Status:   status,
		}
		err = d.rsDB.Model(&model.ServiceConfigObject{}).Create(setMaps).Error
		newRecord := &model.ServiceConfigObject{}
		d.rsDB.Model(&model.ServiceConfigObject{}).Where("tree_id = ?", treeId).Where("keyword=?", keyword).Find(&newRecord)
		id = newRecord.Id
	}
	return
}

// GetServiceConfig 通过tree_id 获取配置
func (d *Dao) GetServiceConfig(c context.Context, treeId int64) (value map[string]string, err error) {
	value = make(map[string]string)
	res := []*model.ServiceConfigObject{}
	err = d.rsDB.Model(&model.ServiceConfigObject{}).Where("tree_id = ?", treeId).Where("status= ?", _statusOn).Find(&res).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	for _, v := range res {
		value[v.Keyword] = v.Value
	}
	return
}

// GetServiceConfigList 管理后台获取服务配置列表
func (d *Dao) GetServiceConfigList(c context.Context, treeName string, treeId int64, keyword string, service string, page int64, pageSize int64, name string, status int64) (list []*v1.MList, totalNum int64, err error) {
	list = []*v1.MList{}
	Items := make([]*model.ServiceConfigObject, 0)
	totalNum = 0
	condition := d.rsDB
	condition = condition.Where("tree_name=?", treeName)
	if 1 == status {
		condition = condition.Where("status=?", status)
	}
	if status != 0 && status != int64(_statusOn) {
		condition = condition.Where("status !=? ", _statusOn)
	}

	if 0 != treeId {
		condition = condition.Where("tree_id =? ", treeId)
	}

	if "" != name {
		condition = condition.Where("name like '%" + name + "%'")
	}

	if "" != keyword {
		condition = condition.Where("keyword like '%" + keyword + "%'")
	}

	if "" != service {
		condition = condition.Where("service =?", service)
	}

	sModel := condition.Model(&model.ServiceConfigObject{})
	sModel.Count(&totalNum)

	iOffset := (page - 1) * pageSize
	err = sModel.Offset(iOffset).Limit(pageSize).Order("mtime DESC, id DESC").Find(&Items).Error
	if nil != err {
		log.Error("sql error select by params")
		return
	}
	for _, v := range Items {
		item := &v1.MList{
			Id:       v.Id,
			TreeName: v.TreeName,
			TreePath: v.TreePath,
			TreeId:   v.TreeId,
			Service:  v.Service,
			Keyword:  v.Keyword,
			Template: v.Template,
			Name:     v.Name,
			Value:    v.Value,
			Ctime:    v.Ctime.Time().Format("2006-01-02 15:04:05"),
			Mtime:    v.Mtime.Time().Format("2006-01-02 15:04:05"),
			Status:   v.Status,
		}
		list = append(list, item)
	}
	return
}

// GetTreeIds 获取treeName对应的tree_ids
func (d *Dao) GetTreeIds(c context.Context, treeName string) (list []int64, err error) {
	list = make([]int64, 0)
	query := "select distinct tree_id  from ap_services_config where status=1 and tree_name = ?"
	rows, err := d.db.Query(c, query, treeName)
	if err != nil && err != sql.ErrNoRows {
		log.Error("[live.titans.dao| getDiscoveryIds] d.db.Exec err: %v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		item := &model.ServiceModel{}
		err = rows.Scan(&item.TreeId)
		if err != nil {
			return
		}
		if item.TreeId != 0 {
			list = append(list, item.TreeId)
		}
	}
	return
}

// GetEasyRecord 获取运营配置
func (d *Dao) GetEasyRecord(c context.Context, treeName string) (list *v1.MList) {
	list = &v1.MList{}
	query := "select id, value from ap_services_config where tree_name = ? and tree_id = 0"
	rows := d.db.QueryRow(c, query, treeName)
	err := rows.Scan(&list.Id, &list.Value)
	if err != nil {
		log.Error("[live.titans.dao| GetEasyRecord] d.db.Exec err: %v", err)
		return
	}
	return
}

// SetEasyRecord 设置运营配置
func (d *Dao) SetEasyRecord(c context.Context, treeName string, value string, id int64) (nId int64, err error) {
	record := &model.InsertServiceConfig{
		TreeName: treeName,
		TreePath: treeName,
		TreeId:   0,
		Name:     treeName + "运营操作列表",
		Value:    value,
		Status:   int64(_statusOn),
	}
	nId = id
	if id != 0 {
		err = d.rsDB.Model(&model.ServiceConfigObject{}).Where("id = ?", id).Update(record).Error
	} else {
		err = d.rsDB.Model(&model.ServiceConfigObject{}).Create(record).Error
		if err != nil {
			return
		}
		newRecord := &model.ServiceConfigObject{}
		err = d.rsDB.Model(&model.ServiceConfigObject{}).Where("tree_name = ?", treeName).Where("tree_id = ?", 0).Find(&newRecord).Error
		if err != nil {
			return
		}
		nId = newRecord.Id
	}
	return
}

// GetListByIds 通过ids获取配置列表
func (d *Dao) GetListByIds(c context.Context, ids []int64) (list []*model.ServiceModel, err error) {
	list = []*model.ServiceModel{}
	Items := make([]*model.ServiceConfigObject, 0)
	condition := d.rsDB

	condition = condition.Where("id in (?)", ids)

	sModel := condition.Model(&model.ServiceConfigObject{})

	err = sModel.Find(&Items).Error
	if nil != err {
		log.Error("sql error select by params")
		return
	}
	for _, v := range Items {
		item := &model.ServiceModel{
			Id:       v.Id,
			TreeName: v.TreeName,
			TreePath: v.TreePath,
			TreeId:   v.TreeId,
			Service:  v.Service,
			Keyword:  v.Keyword,
			Template: v.Template,
			Name:     v.Name,
			Value:    v.Value,
			Ctime:    v.Ctime.Time().Format("2006-01-02 15:04:05"),
			Mtime:    v.Mtime.Time().Format("2006-01-02 15:04:05"),
			Status:   v.Status,
		}
		list = append(list, item)
	}
	return
}
