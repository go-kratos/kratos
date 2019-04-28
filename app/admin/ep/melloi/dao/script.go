package dao

import (
	"time"

	"go-common/app/admin/ep/melloi/model"
)

//QueryScripts query script
func (d *Dao) QueryScripts(script *model.Script, pn int, ps int) (scripts []*model.Script, err error) {
	script.Active = 1
	err = d.DB.Table(model.Script{}.TableName()).Where(script).Order("id asc").Offset((pn - 1) * ps).Limit(ps).Find(&scripts).Error
	return
}

//QueryScriptsByPage query scripts by page
func (d *Dao) QueryScriptsByPage(script *model.Script, pn, ps int32, treeNodes []string) (qsr *model.QueryScriptResponse, err error) {
	script.TestType = 1
	qsr = &model.QueryScriptResponse{}
	gDB := d.DB.Table(model.Script{}.TableName()).Where("app in (?)", treeNodes)
	if script.Department != "" && script.Project != "" && script.App != "" {
		gDB = gDB.Where("department = ? and project = ? and app = ?", script.Department, script.Project, script.App)
	}
	if script.Department != "" && script.Project != "" {
		gDB = gDB.Where("department = ? and project = ?", script.Department, script.Project)
	}
	if script.Department != "" {
		gDB = gDB.Where("department = ?", script.Department)
	}
	if script.TestName != "" {
		gDB = gDB.Where("test_name LIKE ? and test_type = ?", "%"+script.TestName+"%", 1)
	}
	if script.UpdateBy != "" {
		gDB = gDB.Where("update_by LIKE ? and test_type = ?", "%"+script.UpdateBy+"%", 1)
	}

	if script.TestName == "" && script.UpdateBy == "" {
		gDB = gDB.Where(script)
	}

	err = gDB.Where("active = 1").Count(&qsr.TotalSize).Order("ctime desc").Offset((pn - 1) * ps).Limit(ps).
		Find(&qsr.Scripts).Error

	qsr.PageSize = ps /**/
	qsr.PageNum = pn
	return
}

//QueryScriptsByPageWhiteName query by whiteName
func (d *Dao) QueryScriptsByPageWhiteName(script *model.Script, pn, ps int32) (qsr *model.QueryScriptResponse, err error) {
	qsr = &model.QueryScriptResponse{}
	script.TestType = 1
	gDB := d.DB.Table(model.Script{}.TableName())
	if script.TestName != "" {
		gDB = gDB.Where("test_name LIKE ? and test_type = ?", "%"+script.TestName+"%", 1)
	}
	if script.UpdateBy != "" {
		gDB = gDB.Where("update_by LIKE ? and test_type = ?", "%"+script.UpdateBy+"%", 1)
	}

	if script.TestName == "" && script.UpdateBy == "" {
		gDB = gDB.Where(script)
	}

	err = gDB.Where("active = 1").Count(&qsr.TotalSize).Order("ctime desc").Offset((pn - 1) * ps).Limit(ps).
		Find(&qsr.Scripts).Error
	qsr.PageSize = ps /**/
	qsr.PageNum = pn
	return
}

//QueryScriptByID get script by id
func (d *Dao) QueryScriptByID(id int) (script *model.Script, err error) {
	script = &model.Script{}
	err = d.DB.Table(model.Script{}.TableName()).Where("id = ? and active = 1", id).First(script).Error
	return
}

//QueryScriptsInID get script by id
func (d *Dao) QueryScriptsInID(id []int) (scripts []*model.Script, err error) {
	err = d.DB.Table(model.Script{}.TableName()).Where("id in (?) and active = 1", id).Find(&scripts).Error
	return
}

//CountQueryScripts count query scripts
func (d *Dao) CountQueryScripts(script *model.Script) (total int) {
	d.DB.Table(model.Script{}.TableName()).Where(script).Count(&total)
	return
}

//AddScript add script
func (d *Dao) AddScript(script *model.Script) (id int, groupId int, runOrder int, err error) {
	if script.OutputParams == "[]" || script.OutputParams == "" {
		script.OutputParams = "[{\"\":\"\"}]"
	}
	if script.APIHeader == "[]" || script.APIHeader == "" {
		script.APIHeader = "[{\"\":\"\"}]"
	}
	if script.ArgumentString == "[]" || script.AssertionString == "" {
		script.ArgumentString = "[{\"\":\"\"}]"
	}
	script.Active = 1
	err = d.DB.Create(script).Error
	id = script.ID
	groupId = script.GroupID
	runOrder = script.RunOrder
	return
}

//QueryParams query params
func (d *Dao) QueryParams(script *model.Script, scene *model.Scene) (paramList *model.ParamList, err error) {
	paramList = new(model.ParamList)
	if scene.SceneType == 1 || scene.SceneType == 0 {
		err = d.DB.Table(model.Script{}.TableName()).Select("id, group_id, run_order, output_params").Where("scene_id = ? and active = 1", script.SceneID).Order("group_id, run_order").Find(&paramList.ParamList).Error
	} else if scene.SceneType == 2 {
		err = d.DB.Table(model.Script{}.TableName()).Select("id, group_id, run_order, output_params").Where("scene_id = ? and group_id = ? and active = 1", script.SceneID, script.GroupID).Order("group_id, run_order").Find(&paramList.ParamList).Error
	}
	return
}

//AddScriptSnap add script snap
func (d *Dao) AddScriptSnap(scriptSnap *model.ScriptSnap) (id int, err error) {
	scriptSnap.Ctime = time.Now()
	scriptSnap.Mtime = time.Now()
	err = d.DB.Create(scriptSnap).Error
	id = scriptSnap.ID
	return
}

//UpdateScript update script
func (d *Dao) UpdateScript(script *model.Script) (err error) {
	script.Mtime = time.Now()
	if err = d.DB.Exec("update script set is_async = ?, login = ?, keep_alive = ?, use_sign = ?, assertion = ?, conn_time_out = ?, resp_time_out = ?, "+
		"use_data_file = ?, file_name = ?, params_name = ?, delimiter = ?, fusing = ?, use_business_stop = ?, business_stop_percent = ? where id = ?",
		script.IsAsync, script.Login, script.KeepAlive, script.UseSign, script.Assertion, script.ConnTimeOut, script.RespTimeOut,
		script.UseDataFile, script.FileName, script.ParamsName, script.Delimiter, script.Fusing, script.UseBusinessStop, script.BusinessStopPercent, script.ID).Error; err != nil {
		return
	}

	if script.OutputParams == "[]" {
		script.OutputParams = "[{\"\":\"\"}]"
	}
	if script.APIHeader == "[]" {
		script.APIHeader = "[{\"\":\"\"}]"
	}
	if script.ArgumentString == "[]" {
		script.ArgumentString = "[{\"\":\"\"}]"
	}
	//script.Active = 1
	//return d.DB.Model(&model.Script{}).Save(script).Error
	return d.DB.Model(&model.Script{}).Updates(script).Error
}

//UpdateScriptPart update scriptPart
func (d *Dao) UpdateScriptPart(script *model.Script) (err error) {
	script.Mtime = time.Now()
	if script.OutputParams == "[]" {
		script.OutputParams = "[{\"\":\"\"}]"
	}
	if script.APIHeader == "[]" {
		script.APIHeader = "[{\"\":\"\"}]"
	}
	if script.ArgumentString == "[]" {
		script.ArgumentString = "[{\"\":\"\"}]"
	}
	//script.Active = 1
	//return d.DB.Model(&model.Script{}).Save(script).Error
	return d.DB.Model(&model.Script{}).Updates(script).Error
}

//AddScriptTimer add script timer
func (d *Dao) AddScriptTimer(script *model.Script) error {
	return d.DB.Exec("update script set const_timer = ?, random_timer = ? where id = ?", script.ConstTimer, script.RandomTimer, script.ID).Error
}

//DelScript delete script
func (d *Dao) DelScript(id int) error {
	return d.DB.Model(&model.Script{}).Where("id = ?", id).Update("active", 0).Error
}

//QueryScriptSnap query script snap
func (d *Dao) QueryScriptSnap(snap *model.ScriptSnap) (scriptSnap []*model.ScriptSnap, err error) {
	err = d.DB.Table(model.ScriptSnap{}.TableName()).Where(snap).Find(&scriptSnap).Error
	return
}

//ScriptHost query api domain name in script
func (d *Dao) ScriptHost(script *model.Script) (host string, err error) {
	err = d.DB.Create(script).Error
	host = script.Domain
	return
}

////UpdateRunOrder Update Run Order
//func (d *Dao) UpdateRunOrder(runOrder int) error {
//	return d.DB.Model(&model.Script{}).Where("ID=?", script.ID).Updates(script).Error
//}
