package dao

import "go-common/app/admin/ep/melloi/model"

//AddPtestJob add PtestJob
func (d *Dao) AddPtestJob(ptestJob *model.PtestJob) (ptestJobID int, err error) {
	err = d.DB.Create(ptestJob).Error
	ptestJobID = ptestJob.ID
	return
}

//CountPtestJob add PtestJob
func (d *Dao) CountPtestJob(ptestJob *model.PtestJob) (total int, err error) {
	err = d.DB.Table(model.PtestJob{}.TableName()).Where(ptestJob).Count(&total).Error
	return
}

//QueryPtestJob query PtestJob
func (d *Dao) QueryPtestJob(ptestJob *model.PtestJob) (pJob []*model.PtestJob, err error) {
	err = d.DB.Table(model.PtestJob{}.TableName()).Where(ptestJob).Order("ctime asc").Find(&pJob).Error
	return
}

//DeletePtestJob delete PtestJob
func (d *Dao) DeletePtestJob(id int) error {
	return d.DB.Model(&model.PtestJob{}).Where("ID = ?", id).Update("active", -1).Error
}

//DeletePtestJobByJobName delete PtestJob
func (d *Dao) DeletePtestJobByJobName(jobName string) error {
	return d.DB.Model(&model.PtestJob{}).Where("job_name = ?", jobName).Update("active", -1).Error
}

//UpdatePtestJob Update PtestJob
func (d *Dao) UpdatePtestJob(ptestJob model.PtestJob) error {
	return d.DB.Model(&model.PtestJob{}).Where("id = ?", ptestJob.ID).Updates(ptestJob).Error
}
