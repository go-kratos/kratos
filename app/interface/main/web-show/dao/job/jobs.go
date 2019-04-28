package job

import (
	"context"

	"go-common/app/interface/main/web-show/model/job"
	"go-common/library/log"
)

const (
	_selCateSQL = "SELECT id,name,type FROM job_category"
	_selJobsSQL = "SELECT id,name,cate_id,addr_id,duty,demand FROM jobs WHERE status=1"
)

// Categories Dao
func (dao *Dao) Categories(c context.Context) (cs []*job.Category, err error) {
	rows, err := dao.db.Query(c, _selCateSQL)
	if err != nil {
		log.Error("Job.GetAllCate error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &job.Category{}
		if err = rows.Scan(&c.ID, &c.Name, &c.Type); err != nil {
			PromError("Categories", "rows.scan err(%v)", err)
			return
		}
		cs = append(cs, c)
	}
	return
}

// Jobs Dao
func (dao *Dao) Jobs(c context.Context) (js []*job.Job, err error) {
	rows, err := dao.db.Query(c, _selJobsSQL)
	if err != nil {
		log.Error("Job.Jobs error (%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		j := &job.Job{}
		if err = rows.Scan(&j.ID, &j.Name, &j.CateID, &j.AddrID, &j.Duty, &j.Demand); err != nil {
			PromError("Jobs", "rows.scan err(%v)", err)
			return
		}
		js = append(js, j)
	}
	return

}
