package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
)

// DashCurveGraph is a curve graph for leader to show the team members' code coverage
func (s *Service) DashCurveGraph(c context.Context, username string, req *ut.PCurveReq) (res []*ut.PCurveResp, err error) {
	s.appsCache.Lock()
	paths := s.appsCache.PathsByOwner(username)
	s.appsCache.Unlock()

	if err = s.DB.Table("ut_pkganls").Select("SUBSTRING_INDEX(pkg, '/', 5) AS pkg, ROUND(AVG(coverage/100),2) as coverage, SUM(assertions) as assertions, SUM(panics) as panics, SUM(passed) as passed, ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate, SUM(failures) as failures, mtime").
		Where("mtime BETWEEN ? AND ? AND SUBSTRING_INDEX(pkg, '/', 5) IN (?) AND (pkg!=SUBSTRING_INDEX(pkg, '/', 5) OR p.pkg like 'go-common/library/%')", time.Unix(req.StartTime, 0).Format("2006-01-02"), time.Unix(req.EndTime, 0).Format("2006-01-02"), paths).
		Group("date(mtime),SUBSTRING_INDEX(pkg, '/', 5)").Order("mtime DESC").Find(&res).Error; err != nil {
		log.Error("s.ProjectCurveGraph execute sql error(%v)", err)
		return
	}
	return
}

// DashGraphDetail project graph detail for members.
func (s *Service) DashGraphDetail(c context.Context, username string, req *ut.PCurveReq) (res []*ut.PCurveDetailResp, err error) {
	var (
		paths []string
	)
	if req.Path == "" || req.Path == "All" {
		s.appsCache.Lock()
		paths = s.appsCache.PathsByOwner(username)
		s.appsCache.Unlock()
	} else {
		paths = append(paths, req.Path)
	}
	if err = s.DB.Table("ut_pkganls").Select("ROUND(AVG(coverage/100),2) as coverage, SUM(assertions) as assertions, SUM(panics) as panics, SUM(passed) as passed, ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate, SUM(failures) as failures, ut_commit.mtime, ut_commit.username").
		Joins("INNER JOIN  ut_commit ON ut_commit.commit_id=ut_pkganls.commit_id").
		Where("ut_commit.mtime BETWEEN ? AND ? AND SUBSTRING_INDEX(pkg, '/', 5) IN (?) AND pkg!=SUBSTRING_INDEX(pkg, '/', 5)", time.Unix(req.StartTime, 0).Format("2006-01-02"), time.Unix(req.EndTime, 0).Format("2006-01-02"), paths).
		Group("ut_commit.username").Find(&res).Error; err != nil {
		log.Error("s.ProjectGraphDetail execute sql error(%v)", err)
		return
	}
	return
}

// DashGraphDetailSingle project graph detail for Single member.
func (s *Service) DashGraphDetailSingle(c context.Context, username string, req *ut.PCurveReq) (res []*ut.PCurveDetailResp, err error) {
	s.appsCache.Lock()
	paths := s.appsCache.PathsByOwner(username)
	s.appsCache.Unlock()

	if err = s.DB.Table("ut_pkganls").Select("SUBSTRING_INDEX(pkg, '/', 5) AS pkg, ROUND(AVG(coverage/100),2) as coverage, SUM(assertions) as assertions, SUM(panics) as panics, SUM(passed) as passed, ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate, SUM(failures) as failures, ut_commit.mtime, ut_commit.username").
		Joins("INNER JOIN  ut_commit ON ut_commit.commit_id=ut_pkganls.commit_id").
		Where("ut_commit.mtime BETWEEN ? AND ? AND SUBSTRING_INDEX(pkg, '/', 5) IN (?) AND ut_commit.username=? AND (pkg!=SUBSTRING_INDEX(pkg, '/', 5) OR p.pkg like 'go-common/library/%')", time.Unix(req.StartTime, 0).Format("2006-01-02"), time.Unix(req.EndTime, 0).Format("2006-01-02"), paths, req.User).
		Group("SUBSTRING_INDEX(pkg,'/',5)").Find(&res).Error; err != nil {
		log.Error("s.ProjectGraphDetailSingle execute sql error(%v)", err)
		return
	}
	return
}

// DashPkgsTree get all the lastest committed pkgs by username
func (s *Service) DashPkgsTree(c context.Context, path string, username string) (pkgs []*ut.PkgAnls, err error) {
	var (
		mtime = time.Now().Add(-time.Hour * 24 * 30) //two week
		count = strconv.Itoa(strings.Count(path, "/") + 1)
		hql   = fmt.Sprintf("select id,commit_id,merge_id,concat(substring_index(pkg,'/',%s),'/') as pkg,ROUND(AVG(coverage/100),2) as coverage,ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate,sum(panics) as panics,sum(failures) as failures,sum(skipped) as skipped,sum(passed) as passed,sum(assertions) as assertions,html_url,report_url,ctime,mtime from `ut_pkganls` as tmp where tmp.mtime>'%s' and tmp.mtime = (select max(`ut_pkganls`.mtime) from `ut_pkganls` inner join `ut_commit` on ut_commit.`commit_id`=ut_pkganls.`commit_id` where pkg = tmp.pkg and username='%s') group by substring_index(pkg,'/',%s) having pkg like '%s%%'",
			count, mtime, username, count, path)
	)
	if err = s.DB.Raw(hql).Find(&pkgs).Error; err != nil {
		log.Error("apm.Svc DashboardPkgs error(%v)", err)
		return
	}
	for _, pkg := range pkgs {
		if strings.HasSuffix(pkg.PKG, ".go") {
			continue
		}
		if path == pkg.PKG {
			if strings.HasPrefix(path, "go-common/app") && strings.Count(path, "/") < 6 {
				continue
			}
			var files []*ut.PkgAnls
			if files, err = s.dao.ParseUTFiles(c, pkg.HTMLURL); err != nil {
				log.Error("apm.Svc DashboardPkgs error(%v)", err)
				return
			}
			if len(files) == 0 {
				return nil, fmt.Errorf("Get .go files error. Please check your hosts")
			}
			*pkg = *files[0]
			pkgs = append(pkgs, files[1:]...)
		}
	}
	return
}
