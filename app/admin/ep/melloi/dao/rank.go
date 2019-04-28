package dao

import (
	"go-common/app/admin/ep/melloi/model"
)

//TreesQuery query department and performance test count of department
func (d *Dao) TreesQuery() (res *model.TreeList, err error) {
	res = &model.TreeList{}
	//select department, count(department) as count from report_summary where department != '' GROUP BY department ORDER BY count desc
	err = d.DB.Table("report_summary").Select("department, project, app").Group("app").
		Having("department != ''").Order("department").Scan(&res.TreeList).Error
	return
}

//TreeNumQuery query department
func (d *Dao) TreeNumQuery() (res *model.NumList, err error) { /**/
	res = &model.NumList{}
	//select count(DISTINCT department) as count from report_summary where department != ''
	if err = d.DB.Table("report_summary").Select("department, count(DISTINCT department) as dept_num").
		Where("department != ''").Scan(&res.NumList).Error; err != nil {
		return
	}
	if err = d.DB.Table("report_summary").Select("project, count(DISTINCT project) as pro_num").
		Where("project != ''").Scan(&res.NumList).Error; err != nil {
		return
	}
	if err = d.DB.Table("report_summary").Select("app, count(DISTINCT app) as app_num").
		Where("app != ''").Scan(&res.NumList).Error; err != nil {
		return
	}
	return
}

//TopHttpQuery query performance test top api
func (d *Dao) TopHttpQuery() (res *model.TopAPIRes, err error) {
	res = &model.TopAPIRes{}
	//select s.url, count(r.script_id) as count from report_summary r INNER JOIN script s on r.script_id = s.id GROUP BY s.url having url != '' ORDER BY count desc
	err = d.DB.Limit(10).Table("report_summary").Select("script.url, count(report_summary.script_id) as count").
		Joins("inner join script on report_summary.script_id = script.id").
		Where("report_summary.active = 1 and report_summary.debug != 1 and script.test_type = 1 and report_summary.type = 0").
		Group("script.url").Having("script.url != ''").Order("count desc").Scan(&res.APIList).Error
	return
}

//TopGrpcQuery Top Grpc Query
func (d *Dao) TopGrpcQuery() (res *model.GrpcRes, err error) {
	res = &model.GrpcRes{}
	//select g.service_name, g.request_method, count(r.script_id) as count from report_summary r INNER JOIN grpc g on r.script_id = g.id
	//where g.service_name != '' and r.active = 1 and r.debug != 1 and r.type = 1
	//GROUP BY g.service_name ORDER BY count desc
	err = d.DB.Limit(10).Table("report_summary").Select("grpc.service_name, grpc.request_method, count(report_summary.script_id) as count").
		Joins("inner join grpc on report_summary.script_id = grpc.id").
		Where("report_summary.active = 1 and report_summary.debug != 1 and grpc.service_name != '' and report_summary.type = 1").
		Group("grpc.service_name").Order("count desc").Scan(&res.GrpcList).Error
	return
}

//TopSceneQuery Top Scene Query
func (d *Dao) TopSceneQuery() (res *model.SceneRes, err error) {
	res = &model.SceneRes{}
	//select s.scene_name, count(r.script_id) as count from report_summary r INNER JOIN scene s on r.scene_id = s.id
	//where s.scene_name != '' and s.is_active = 1 and s.is_draft = 0 and r.active = 1 and r.debug != 1 and r.scene_id != 0 and r.type = 2
	//GROUP BY s.scene_name ORDER BY count desc
	err = d.DB.Limit(10).Table("report_summary").Select("scene.department, scene.scene_name, count(report_summary.script_id) as count").
		Joins("inner join scene on report_summary.scene_id = scene.id").
		Where("report_summary.active = 1 and report_summary.debug != 1 and scene.scene_name != '' and scene.is_active = 1 and scene.is_draft = 0 and report_summary.scene_id != 0 and report_summary.type = 2").
		Group("scene.scene_name").Order("count desc").Scan(&res.SceneList).Error
	return
}

//TopDeptQuery query performance test top department
func (d *Dao) TopDeptQuery() (res *model.TopDeptRes, err error) {
	res = &model.TopDeptRes{}
	//select department, count(department) as count from report_summary where department != '' GROUP BY department ORDER BY count desc
	err = d.DB.Limit(10).Table("report_summary").Select("department, count(department) as count").
		Where("report_summary.active = 1 and report_summary.debug != 1").
		Group("department").Having("department != ''").Order("count desc").Scan(&res.DeptList).Error
	return
}

//BuildLineQuery query performance test count by time
func (d *Dao) BuildLineQuery(rank *model.Rank, summary *model.ReportSummary) (res *model.BuildLineRes, err error) {
	res = &model.BuildLineRes{}
	//select DATE_FORMAT(ctime, '%H') as count from report_summary where ctime >= date_sub(now(), interval 24 hour) AND ctime <= NOW()
	//switch rank.TimeDegree {
	//case "H":
	//	err = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m-%d %H') as date").
	//		Where(summary).Where("active = 1 and debug != 1").
	//		Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime).Order("date").
	//		Scan(&res.BuildList).Error
	//case "d":
	//	err = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m-%d') as date").Where(summary).
	//		Where("active = 1 and debug != 1").
	//		Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime).Order("date").Scan(&res.BuildList).Error
	//case "m":
	//	err = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m') as date").Where(summary).
	//		Where("active = 1 and debug != 1").
	//		Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime).Order("date").Scan(&res.BuildList).Error
	//case "Y":
	//	err = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y') as date").Where(summary).
	//		Where("active = 1 and debug != 1").
	//		Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime).Order("date").Scan(&res.BuildList).Error
	//default:
	//	err = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m-%d %H') as date").Where(summary).
	//		Where("active = 1 and debug != 1").
	//		Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime).Order("date").Scan(&res.BuildList).Error
	//}
	gDB := d.DB.Table(model.ReportSummary{}.TableName())
	switch rank.TimeDegree {
	case "H":
		gDB = gDB.Select("DATE_FORMAT(ctime, '%Y-%m-%d %H') as date").
			Where("active = 1 and debug != 1").
			Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime)
	case "d":
		gDB = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m-%d') as date").
			Where("active = 1 and debug != 1").
			Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime)
	case "m":
		gDB = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m') as date").
			Where("active = 1 and debug != 1").
			Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime)
	case "Y":
		gDB = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y') as date").
			Where("active = 1 and debug != 1").
			Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime)
	default:
		gDB = d.DB.Table("report_summary").Select("DATE_FORMAT(ctime, '%Y-%m-%d %H') as date").
			Where("active = 1 and debug != 1").
			Where("ctime >= ? AND ctime <= ?", rank.StartTime, rank.EndTime)
	}
	if rank.SearchAll {
		err = gDB.Where("type in (0, 1, 2)").Order("date").Scan(&res.BuildList).Error
	} else {
		err = gDB.Where("type = ?", summary.Type).Order("date").Scan(&res.BuildList).Error
	}
	return
}

//StateLineQuery query statistic of state
func (d *Dao) StateLineQuery() (res *model.StateLineRes, err error) {
	res = &model.StateLineRes{}
	//select test_status, count(test_status) as count from report_summary GROUP BY test_status
	err = d.DB.Table("report_summary").Select("test_status, count(test_status) as count").Where("test_status != 0").
		Where("active = 1 and debug != 1").Group("test_status").Scan(&res.StateList).Error
	return
}
