package dao

import (
	"strconv"
	"strings"

	"go-common/app/admin/ep/melloi/conf"
	"go-common/app/admin/ep/melloi/model"
)

// AddScene Add Scene
func (d *Dao) AddScene(scene *model.Scene) (sceneId int, err error) {
	scene.IsActive = true
	err = d.DB.Table("scene").Create(scene).Error
	sceneId = scene.ID
	return
}

// QueryDraft Query Draft
func (d *Dao) QueryDraft(scene *model.Scene) (res *model.QueryDraft, err error) {
	res = &model.QueryDraft{}
	err = d.DB.Table("scene").Select("id as scene_id, scene_name").Where(scene).Where("is_draft = 1 and is_active = 1").Count(&res.Total).Order("mtime desc").Find(&res.Drafts).Error
	return
}

// UpdateScene Update Scene
func (d *Dao) UpdateScene(scene *model.Scene, scriptIDList []int) (fusing int, err error) {
	if err = d.DB.Exec("update scene set fusing = ? where id = ?", scene.Fusing, scene.ID).Error; err != nil {
		return
	}
	if err = d.DB.Model(&model.Scene{}).Where("ID = ?", scene.ID).Update(scene).Error; err != nil {
		return
	}
	if scene.IsBatch { //选择已有接口点击确定时执行
		if scene.Fusing == 0 {
			err = d.DB.Model(&model.Script{}).Where("scene_id = ? and id in (?)", scene.ID, scriptIDList).Update("fusing", conf.Conf.Melloi.DefaultFusing).Error
		} else {
			err = d.DB.Model(&model.Script{}).Where("scene_id = ? and id in (?)", scene.ID, scriptIDList).Update("fusing", scene.Fusing).Error
		}
	} else if scene.IsUpdate { //
		err = d.DB.Model(&model.Script{}).Where("scene_id = ?", scene.ID).Update("fusing", scene.Fusing).Error
	}
	fusing = scene.Fusing
	return
}

// SaveScene Save Scene
func (d *Dao) SaveScene(scene *model.Scene) error {
	return d.DB.Model(&model.Scene{}).Where("ID = ?", scene.ID).Update("is_draft", 0).Update("scene_type", scene.SceneType).Error
}

// SaveOrder Save Order
func (d *Dao) SaveOrder(reqList []*model.GroupOrder, scene *model.Scene) (err error) {
	for _, req := range reqList {
		if scene.SceneType == 1 {
			err = d.DB.Model(&model.Script{}).Where("ID = ?", req.ID).Update("group_id", req.GroupID).Error
		} else {
			err = d.DB.Model(&model.Script{}).Where("ID = ?", req.ID).Update("run_order", req.RunOrder).Error
		}
	}
	return
}

// QueryGroupId Query GroupId
func (d *Dao) QueryGroupId(script *model.Script) (res *model.QueryRelation, groupId int, err error) {
	res = &model.QueryRelation{}
	//select group_id, count(group_id) as count from script where scene_id = 1 and group_id in (select group_id from script where id = 777) group by group_id
	err = d.DB.Table("script").Select("group_id").Where("id = ?", script.ID).Find(&res.RelationList).Error
	groupId = res.RelationList[0].GroupID
	return
}

// QueryRelation Query Relation
func (d *Dao) QueryRelation(groupId int, script *model.Script) (res *model.QueryRelation, err error) {
	res = &model.QueryRelation{}
	err = d.DB.Table("script").Select("group_id, count(group_id) as count").Where("scene_id = ? and group_id = ?", script.SceneID, groupId).Find(&res.RelationList).Error
	return
}

// QueryAPI Query API
func (d *Dao) QueryAPI(scene *model.Scene) (res *model.QueryAPIs, err error) {
	res = &model.QueryAPIs{}
	if err = d.DB.Table(model.Script{}.TableName()).Select("id").Where("scene_id = ? and active = 1", scene.ID).Count(&res.Total).Error; err != nil {
		return
	}
	if res.Total > 0 {
		//select scene.scene_name, scene.scene_type, scene.department, scene.project, scene.app, script.group_id, script.run_order, script.id as test_id, script.test_name
		//from scene INNER JOIN script on scene.id = script.scene_id and script.active = 1 where script.scene_id = 1
		gDB := d.DB.Table(model.Scene{}.TableName()).Select("scene.id as scene_id, scene.scene_name, scene.scene_type, scene.department, " +
			"scene.project, scene.app, script.group_id, script.run_order, script.id, script.test_name, script.url, script.output_params, script.threads_sum, script.load_time").
			Joins("inner join script on scene.id = script.scene_id and script.active = 1")
		if scene.SceneType == 2 {
			err = gDB.Where(scene).Count(&res.Total).Order("group_id, run_order").Find(&res).Find(&res.APIs).Error
			res.SceneType = 2
		} else {
			err = gDB.Where(scene).Count(&res.Total).Order("script.id").Find(&res).Find(&res.APIs).Error
			res.SceneType = 1
		}
	} else {
		gDB := d.DB.Table(model.Scene{}.TableName()).Select("id as scene_id, scene_name, department, project, app")
		err = gDB.Where("id = ?", scene.ID).Find(&res).Error
		//查询无结果时，APIs返回空数组
		res.APIs = []*model.TestAPI{}
	}
	return
}

// DeleteAPI Delete API
func (d *Dao) DeleteAPI(script *model.Script) error {
	return d.DB.Model(&model.Script{}).Where("ID = ?", script.ID).Update("active", 0).Error
}

// AddConfig Add Config
func (d *Dao) AddConfig(script *model.Script) error {
	return d.DB.Exec("update script set threads_sum = ?, load_time = ?, ready_time = ? where scene_id = ? and group_id = ?", script.ThreadsSum, script.LoadTime, script.ReadyTime, script.SceneID, script.GroupID).Error
}

// QueryTree Query Tree
func (d *Dao) QueryTree(script *model.Script) (res *model.ShowTree, err error) {
	res = &model.ShowTree{}
	err = d.DB.Table(model.Script{}.TableName()).Select("department, project, app").Where(script).Count(&res.IsShow).Order("id").Limit(1).Find(&res.Tree).Error
	if res.IsShow > 1 {
		res.IsShow = 1
	}
	return
}

//QueryScenesByPage query scripts by page
func (d *Dao) QueryScenesByPage(scene *model.Scene, pn, ps int32, treeNodes []string) (qsr *model.QuerySceneResponse, err error) {
	qsr = &model.QuerySceneResponse{}
	gDB := d.DB.Table(model.Scene{}.TableName()).Where("app in (?)", treeNodes).Where("is_draft = 0")
	if scene.Department != "" && scene.Project != "" && scene.APP != "" {
		gDB = gDB.Where("department = ? and project = ? and app = ?", scene.Department, scene.Project, scene.APP)
	}
	if scene.Department != "" && scene.Project != "" {
		gDB = gDB.Where("department = ? and project = ?", scene.Department, scene.Project)
	}
	if scene.Department != "" {
		gDB = gDB.Where("department = ?", scene.Department)
	}
	if scene.SceneName != "" {
		gDB = gDB.Where("scene_name LIKE ? ", "%"+scene.SceneName+"%")
	}
	if scene.UserName != "" {
		gDB = gDB.Where("user_name LIKE ? ", "%"+scene.UserName+"%")
	}

	if scene.SceneName == "" && scene.UserName == "" && scene.Department == "" && scene.Project == "" && scene.APP == "" {
		gDB = gDB.Where(scene)
	}

	err = gDB.Where("is_active = 1").Count(&qsr.TotalSize).Order("ctime desc").Offset((pn - 1) * ps).Limit(ps).
		Find(&qsr.Scenes).Error
	qsr.PageSize = ps /**/
	qsr.PageNum = pn
	return
}

//QueryScenesByPageWhiteName query scripts by page white name
func (d *Dao) QueryScenesByPageWhiteName(scene *model.Scene, pn, ps int32) (qsr *model.QuerySceneResponse, err error) {
	qsr = &model.QuerySceneResponse{}
	gDB := d.DB.Table(model.Scene{}.TableName()).Where("is_draft = 0")
	if scene.SceneName != "" {
		gDB = gDB.Where("scene_name LIKE ? ", "%"+scene.SceneName+"%")
	}
	if scene.UserName != "" {
		gDB = gDB.Where("user_name LIKE ? ", "%"+scene.UserName+"%")
	}

	if scene.SceneName == "" && scene.UserName == "" {
		gDB = gDB.Where(scene)
	}

	err = gDB.Where("is_active = 1").Count(&qsr.TotalSize).Order("mtime desc").Offset((pn - 1) * ps).Limit(ps).
		Find(&qsr.Scenes).Error
	qsr.PageSize = ps /**/
	qsr.PageNum = pn
	return
}

//QueryScenes query scene
func (d *Dao) QueryScenes(scene *model.Scene, pn int, ps int) (scenes []*model.Scene, err error) {
	err = d.DB.Table(model.Scene{}.TableName()).Where(scene).Order("mtime desc").Offset((pn - 1) * ps).Limit(ps).Find(&scenes).Error
	return
}

// QueryExistAPI Query Exist API
func (d *Dao) QueryExistAPI(script *model.Script, pageNum int32, pageSize int32, sceneId int, treeNodes []string) (res *model.APIInfoList, err error) {
	res = &model.APIInfoList{}
	gDB := d.DB.Table(model.Script{}.TableName()).Where("app in (?)", treeNodes)
	if pageSize == 0 && pageNum == 0 {
		err = gDB.Where(script).Where("active = 1 and scene_id = 0").Count(&res.TotalSize).Order("id desc").Find(&res.ScriptList).Error
	} else {
		err = gDB.Where(script).Where("active = 1 and scene_id = 0").Count(&res.TotalSize).Order("id desc").Offset((pageNum - 1) * pageSize).Limit(pageSize).Find(&res.ScriptList).Error
		res.PageNum = pageNum
		res.PageSize = pageSize
	}
	res.SceneID = sceneId
	return
}

// QueryGroup Query Group
func (d *Dao) QueryGroup(sceneId int) (res *model.GroupList, err error) {
	res = &model.GroupList{}
	gDB := d.DB.Table(model.Script{}.TableName())
	//select group_id, count(group_id) from script where scene_id = 1 GROUP BY group_id order by group_id
	err = gDB.Select("group_id, threads_sum, load_time, ready_time").Where("scene_id = ? and active = 1", sceneId).Group("group_id").Order("group_id").Find(&res.GroupList).Error
	return
}

// QueryPreview Query Preview
func (d *Dao) QueryPreview(sceneId, groupId int) (res *model.PreviewList, err error) {
	res = &model.PreviewList{}
	gDB := d.DB.Table(model.Script{}.TableName())
	//select id as test_id, test_name, run_order from script where scene_id = 1 and group_id = 1 order by run_order
	err = gDB.Select("id, test_name, run_order, group_id, const_timer, random_timer").Where("scene_id = ? and group_id = ? and active = 1", sceneId, groupId).Order("run_order").Find(&res.PreList).Error
	return
}

// QueryUsefulParams Query Useful Params
func (d *Dao) QueryUsefulParams(sceneId int) (res *model.UsefulParamsList, err error) {
	res = &model.UsefulParamsList{}
	gDB := d.DB.Table(model.Script{}.TableName())
	err = gDB.Select("output_params").Where("scene_id = ?", sceneId).Where("active = 1 and output_params != '[{\"\":\"\"}]'").Order("id desc").Find(&res.ParamsList).Error
	return
}

// UpdateBindScene Update Bind Scene
func (d *Dao) UpdateBindScene(bindScene *model.BindScene) (err error) {
	var id int
	tempList := strings.Split(bindScene.ID, ",")
	for _, tempId := range tempList {
		if id, err = strconv.Atoi(tempId); err != nil {
			return err
		}
		err = d.DB.Model(&model.Script{}).Where("id =?", id).Updates(model.Script{SceneID: bindScene.SceneID}).Error
	}
	return
}

// QueryDrawRelation Query Draw Relation
func (d *Dao) QueryDrawRelation(scene *model.Scene) (res *model.SaveOrderReq, err error) {
	res = &model.SaveOrderReq{}
	gDB := d.DB.Table(model.Script{}.TableName())
	err = gDB.Select("id, test_name, group_id, run_order").Where("scene_id = ? and active = 1", scene.ID).Order("group_id, run_order").Find(&res.GroupOrderList).Error
	return
}

// DeleteDraft Delete Draft
func (d *Dao) DeleteDraft(scene *model.Scene) error {
	if scene.ID == 0 {
		return d.DB.Model(&model.Scene{}).Where("user_name = ? and is_draft = 1", scene.UserName).Update("is_active", 0).Error
	}
	return d.DB.Model(&model.Scene{}).Where("user_name = ? and id = ? and is_draft = 1", scene.UserName, scene.ID).Update("is_active", 0).Error
}

// QueryConfig Query Config
func (d *Dao) QueryConfig(script *model.Script) (res *model.GroupInfo, err error) {
	res = &model.GroupInfo{}
	gDB := d.DB.Table(model.Script{}.TableName())
	// select distinct group_id, threads_sum, ready_time, load_time from script where scene_id = 282 and active = 1 and group_id = 2
	err = gDB.Select("distinct group_id, threads_sum, ready_time, load_time").Where("scene_id = ? and group_id = ? and active = 1", script.SceneID, script.GroupID).Find(&res).Error
	return
}

// DeleteScene Delete Scene
func (d *Dao) DeleteScene(scene *model.Scene) error {
	return d.DB.Model(&model.Scene{}).Where("id = ?", scene.ID).Update("is_active", 0).Error
}

// QueryFusing Query Fusing
func (d *Dao) QueryFusing(script *model.Script) (res *model.FusingInfoList, err error) {
	res = &model.FusingInfoList{}
	gDB := d.DB.Table(model.Script{}.TableName())
	//select id as test_id, test_name, run_order from script where scene_id = 1 and group_id = 1 order by run_order
	err = gDB.Select("fusing").Where("scene_id = ? and active = 1", script.SceneID).Find(&res.FusingList).Error
	return
}
