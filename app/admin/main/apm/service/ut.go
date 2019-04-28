package service

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"go-common/app/admin/main/apm/conf"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
	xtime "go-common/library/time"

	"github.com/jinzhu/gorm"
)

const (
	_utMerge  = "ut_merge"
	_utCommit = "ut_commit"
	_utpkg    = "ut_pkganls"
)

// UtList return ut_list by merge_id or username
func (s *Service) UtList(c context.Context, arg *ut.MergeReq) (mrInfs []*ut.Merge, count int, err error) {
	var (
		mtime = time.Now().Add(-time.Hour * 24 * 90)
		hql   = fmt.Sprintf("mtime>'%s'", mtime)
	)
	if arg.MergeID != 0 {
		hql += fmt.Sprintf(" and merge_id=%d ", arg.MergeID)
	}
	if arg.UserName != "" {
		hql += fmt.Sprintf(" and username='%s' ", arg.UserName)
	}
	if arg.IsMerged != 0 {
		hql += fmt.Sprintf(" and is_merged='%d' ", arg.IsMerged)
	}
	if err = s.DB.Table(_utMerge).Where(hql).Count(&count).Error; err != nil {
		log.Error("service.UtList error(%v)", err)
		return
	}
	if err = s.DB.Where(hql).Offset((arg.Pn - 1) * arg.Ps).Limit(arg.Ps).Order("mtime desc").Find(&mrInfs).Error; err != nil {
		log.Error("service.UtList error(%v)", err)
		return
	}
	for _, mr := range mrInfs {
		var commit = &ut.Commit{}
		if err = s.DB.Where("merge_id=?", mr.MergeID).Order("id desc").First(commit).Error; err != nil {
			log.Error("service.UtList mergeID(%d) error(%v)", mr.MergeID, err)
			return
		}
		if err = s.DB.Select(`max(commit_id) as commit_id, substring_index(pkg,"/",5) as pkg, ROUND(AVG(coverage/100),2) as coverage, SUM(assertions) as assertions, SUM(panics) as panics, SUM(passed) as passed, ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate, SUM(failures) as failures, MAX(mtime) as mtime`).
			Where(`commit_id=? and (pkg!=substring_index(pkg, "/", 5) or pkg like "go-common/library/%")`, commit.CommitID).
			Group(`substring_index(pkg, "/", 5)`).Find(&commit.PkgAnls).Error; err != nil {
			log.Error("service.UtList commitID(%s) error(%v)", commit.CommitID, err)
			return
		}
		mr.Commit = commit
	}
	return
}

// UtDetailList get ut pkganls by commit_id&pkg
func (s *Service) UtDetailList(c context.Context, arg *ut.DetailReq) (utpkgs []*ut.PkgAnls, err error) {
	var (
		mtime = time.Now().Add(-time.Hour * 24 * 90)
		hql   = fmt.Sprintf("mtime>'%s'", mtime)
	)
	if arg.CommitID != "" {
		hql += fmt.Sprintf(" and commit_id='%s'", arg.CommitID)
	}
	if arg.PKG != "" {
		hql += fmt.Sprintf(" and substring_index(pkg,'/',5)='%s' and (pkg!=substring_index(pkg,'/',5) or pkg like 'go-common/library/%%')", arg.PKG)
	}
	if err = s.dao.DB.Select("id, commit_id, merge_id, pkg, ROUND(coverage/100,2) as coverage, ROUND(passed/assertions*100,2) as pass_rate,panics,failures,skipped,passed,assertions,html_url,report_url,mtime,ctime").Where(hql).Find(&utpkgs).Error; err != nil {
		log.Error("service.UTDetilList commitID(%s) project(%s) error(%v)", arg.CommitID, arg.PKG, err)
		return
	}
	return
}

// UtHistoryCommit get commits history list by merge_id
func (s *Service) UtHistoryCommit(c context.Context, arg *ut.HistoryCommitReq) (utcmts []*ut.Commit, count int, err error) {
	var (
		mtime = time.Now().Add(-time.Hour * 24 * 90)
		hql   = fmt.Sprintf("mtime>'%s'", mtime)
	)
	hql += fmt.Sprintf(" and merge_id=%d", arg.MergeID)
	if err = s.DB.Table(_utCommit).Where(hql).Count(&count).Error; err != nil {
		log.Error("service.UTHistoryCommit mergeID(%d) error(%v)", arg.MergeID, err)
		return
	}
	if err = s.dao.DB.Where(hql).Offset((arg.Pn - 1) * arg.Ps).Limit(arg.Ps).Order("mtime desc").Find(&utcmts).Error; err != nil {
		log.Error("service.UTHistoryCommit mergeID(%d) error(%v)", arg.MergeID, err)
		return
	}
	for _, commit := range utcmts {
		if err = s.dao.DB.Select(`max(commit_id) as commit_id, substring_index(pkg,"/",5) as pkg, ROUND(AVG(coverage/100),2) as coverage, SUM(assertions) as assertions, SUM(panics) as panics, SUM(passed) as passed, ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate, SUM(failures) as failures, MAX(mtime) as mtime`).
			Where("commit_id=? and (pkg!=substring_index(pkg,'/',5) or pkg like 'go-common/library/%')", commit.CommitID).Group(`substring_index(pkg, "/", 5)`).
			Order(`substring_index(pkg, "/", 5)`).Find(&commit.PkgAnls).Error; err != nil {
			log.Error("service.UTHistoryCommit commitID(%s) error(%v)", commit.CommitID, err)
			return
		}
	}
	for i := 0; i < len(utcmts)-1; i++ {
		for _, curPkg := range utcmts[i].PkgAnls {
			if utcmts[i+1] == nil || len(utcmts[i+1].PkgAnls) == 0 {
				continue
			}
			for _, prePkg := range utcmts[i+1].PkgAnls {
				if prePkg.PKG == curPkg.PKG {
					curPkg.CovChange, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", curPkg.Coverage-prePkg.Coverage), 64)
				}
			}
		}
	}
	return
}

// AddUT store ut data with a transaction
func (s *Service) AddUT(c context.Context, pkg *ut.PkgAnls, files []*ut.File, request *ut.UploadRes, dataURL, reportURL, HTMLURL string) (err error) {
	var (
		tx = s.DB.Begin()
		// pkgCoverage float64
	)
	if err = s.AddUTMerge(c, tx, request); err != nil {
		tx.Rollback()
		return
	}
	if err = s.AddUTCommit(c, tx, request); err != nil {
		tx.Rollback()
		return
	}
	if err = s.AddUTPKG(c, tx, pkg, request, reportURL, HTMLURL, dataURL); err != nil {
		tx.Rollback()
		return
	}
	if err = s.AddUTFiles(c, tx, files); err != nil {
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

// AddUTMerge insert merge if not exist
func (s *Service) AddUTMerge(c context.Context, tx *gorm.DB, request *ut.UploadRes) (err error) {
	var (
		uts   = make([]*ut.Merge, 0)
		merge = &ut.Merge{
			UserName: request.UserName,
			MergeID:  request.MergeID,
		}
	)
	if err = tx.Where("merge_id=?", merge.MergeID).Order("mtime desc").Find(&uts).Error; err != nil {
		log.Error("service.CreateUtList err(%v)", err)
		return
	}
	if len(uts) > 0 {
		return
	}
	if err = tx.Create(merge).Error; err != nil {
		log.Error("service.CreateUtInfo err (%v)", err)
	}
	return
}

// AddUTPKG if not exists pkg.commit_id then insert,otherwise,update
func (s *Service) AddUTPKG(c context.Context, tx *gorm.DB, pkg *ut.PkgAnls, request *ut.UploadRes, reportURL, HTMLURL, dataURL string) (err error) {
	var (
		pkgs = make([]*ut.PkgAnls, 0)
		in   = map[string]interface{}{
			"merge_id":   request.MergeID,
			"failures":   pkg.Failures,
			"passed":     pkg.Passed,
			"assertions": pkg.Assertions,
			"panics":     pkg.Panics,
			"skipped":    pkg.Skipped,
			"coverage":   pkg.Coverage,
			"html_url":   HTMLURL,
			"report_url": reportURL,
			"data_url":   dataURL,
			"mtime":      xtime.Time(time.Now().Unix())}
	)
	pkg.MergeID = request.MergeID
	pkg.PKG = request.PKG
	pkg.CommitID = request.CommitID
	pkg.HTMLURL = HTMLURL
	pkg.ReportURL = reportURL
	pkg.DataURL = dataURL
	if err = tx.Where("commit_id=? and pkg=?", pkg.CommitID, pkg.PKG).Find(&pkgs).Error; err != nil {
		log.Error("service.CreateUtPKG query error(%v)", err)
		return
	}
	if len(pkgs) > 0 {
		log.Info("s.CreateUtPKG commit_id(%s) pkg(%+v)", pkg.CommitID, pkg)
		err = tx.Table(_utpkg).Where("commit_id=? and pkg=?", pkg.CommitID, pkg.PKG).Update(in).Error
	} else {
		err = tx.Create(pkg).Error
	}
	if err != nil {
		log.Error("service.CreateUtPKG create error(%v)", err)
	}
	return
}

// AddUTFiles add files in table ut_file
func (s *Service) AddUTFiles(c context.Context, tx *gorm.DB, files []*ut.File) (err error) {
	var (
		recordFiles []*ut.File
	)
	if len(files) == 0 {
		err = fmt.Errorf("length of files is 0")
		log.Error("apmSvc.AddUTFiles Error(%v)", err)
		return
	}
	if err = tx.Where("commit_id=? and pkg=?", files[0].CommitID, files[0].PKG).Find(&recordFiles).Error; err != nil {
		log.Error("apmSvc.AddUTFiles Error(%v)", err)
		return
	}
	if len(recordFiles) > 0 {
		return
	}
	for _, file := range files {
		if err = tx.Create(file).Error; err != nil {
			log.Error("service.AddUTFiles create error(%v)", err)
			return
		}
	}
	return
}

// AddUTCommit insert commit if not exist
func (s *Service) AddUTCommit(c context.Context, tx *gorm.DB, request *ut.UploadRes) (err error) {
	var (
		cts = make([]*ut.Commit, 0)
		ct  = &ut.Commit{
			CommitID: request.CommitID,
			MergeID:  request.MergeID,
		}
	)
	if request.Author == "" || request.Author == "null" {
		ct.UserName = request.UserName
	} else {
		ct.UserName = request.Author
	}
	if err = tx.Where("commit_id=?", ct.CommitID).Find(&cts).Error; err != nil {
		log.Error("s.CreateUtCommit query error(%v)", err)
		return
	}
	if len(cts) > 0 {
		log.Info("s.AddUTCommit(%d) update merge request", ct.MergeID)
		err = tx.Table(_utCommit).Where("commit_id=?", ct.CommitID).Update("merge_id", ct.MergeID).Error
	} else {
		err = tx.Create(ct).Error
	}
	if err != nil {
		log.Error("s.CreateUtCommit(%+v) error(%v)", ct, err)
	}
	return
}

// SetMerged set is_merged = 1 in ut_merge
func (s *Service) SetMerged(c context.Context, mrid int64) (err error) {
	if err = s.dao.DB.Table(_utMerge).Where("merge_id=?", mrid).Update("is_merged", 1).Error; err != nil {
		log.Error("s.SetMerged error(%v)", err)
		return
	}
	return
}

// CheckUT .coverage>=30% && pass_rate=100% && increase >=0
func (s *Service) CheckUT(c context.Context, cid string) (t *ut.Tyrant, err error) {
	var (
		pkgs = make([]*ut.PkgAnls, 0)
	)
	if err = s.DB.Where("commit_id=? and (pkg!=substring_index(pkg,'/',5) or pkg like 'go-common/library/%')", cid).Find(&pkgs).Error; err != nil {
		log.Error("s.Tyrant query by commit_id error(%v)", err)
		return
	}
	if len(pkgs) == 0 {
		err = fmt.Errorf("pkgs is none")
		log.Error("s.Tyrant query by commit_id error(%v)", err)
		return
	}
	for _, pkg := range pkgs {
		if t, err = s.tyrant(pkg); err != nil {
			return
		}
		if t.Tyrant {
			break
		}
	}
	return
}

// tyrant determine whether the standard is up to standard
func (s *Service) tyrant(pkg *ut.PkgAnls) (t *ut.Tyrant, err error) {
	var (
		pkgs     = make([]*ut.PkgAnls, 0)
		curCover = pkg.Coverage / 100
		preCover float64
	)
	if err = s.DB.Select("commit_id, coverage").Where("pkg=?", pkg.PKG).Order("coverage desc").Find(&pkgs).Error; err != nil {
		log.Error("s.Tyrant query by pkg error(%v)", err)
		return
	}
	if len(pkgs) == 0 {
		return
	} else if len(pkgs) == 1 || pkg.CommitID != pkgs[0].CommitID {
		preCover = pkgs[0].Coverage / 100
	} else if pkg.CommitID == pkgs[0].CommitID {
		preCover = pkgs[1].Coverage / 100
	}
	t = &ut.Tyrant{
		Package:  pkg.PKG,
		Coverage: curCover,
		PassRate: float64(pkg.Passed) / float64(pkg.Assertions) * 100,
		Tyrant:   curCover < float64(s.c.UTBaseLine.Coverage) || pkg.Passed != pkg.Assertions,
		Standard: s.c.UTBaseLine.Coverage,
		Increase: curCover - preCover,
		LastCID:  pkgs[0].CommitID,
	}
	return
}

//QATrend give single user's coverage ,passrate and score data
func (s *Service) QATrend(c context.Context, arg *ut.QATrendReq) (trend *ut.QATrendResp, err error) {
	var (
		date, group, order string
		details            []*ut.PkgAnls
		where              = "1=1 and (ut_pkganls.pkg!=substring_index(ut_pkganls.pkg,'/',5) or ut_pkganls.pkg like 'go-common/library/%')"
	)
	if arg.User != "" {
		where += fmt.Sprintf(" and ut_commit.username = '%s'", arg.User)
	}
	if arg.StartTime != 0 && arg.EndTime != 0 {
		where += fmt.Sprintf(" and ut_pkganls.mtime >= '%s' and ut_pkganls.mtime <= '%s'", time.Unix(arg.StartTime, 0).Format("2006-01-02 15:04:05"), time.Unix(arg.EndTime, 0).Format("2006-01-02 15:04:05"))
	} else {
		where += fmt.Sprintf(" and ut_pkganls.mtime >= '%s'", time.Now().AddDate(0, 0, -arg.LastTime))
	}
	if arg.Period == "hour" {
		order += "SUBSTRING(ut_pkganls.mtime,11,2) asc"
	} else {
		order += "mtime asc"
	}
	if arg.Period == "day" {
		group += "LEFT(ut_pkganls.mtime,10)"
	} else {
		group += fmt.Sprintf("%s(ut_pkganls.mtime)", arg.Period)
	}
	if err = s.DB.Table(_utpkg).Joins("join ut_commit on ut_pkganls.commit_id = ut_commit.commit_id").
		Select(`ROUND(AVG(coverage/100),2) AS coverage,ROUND(SUM(passed)/SUM(assertions)*100,2) AS pass_rate, ROUND((AVG(coverage)/100*(SUM(passed)/SUM(assertions))),2) AS score, MAX(ut_pkganls.mtime) AS mtime,GROUP_CONCAT(Distinct ut_commit.commit_id) AS cids`).
		Where(where).Group(group).Order(order).Find(&details).Error; err != nil {
		log.Error("service.QATrend find detail error(%v)", err)
		return
	}
	trend = &ut.QATrendResp{BaseLine: conf.Conf.UTBaseLine.Coverage}
	for _, detail := range details {
		switch arg.Period {
		case "hour":
			date = fmt.Sprintf("%d时", detail.MTime.Time().Hour())
		case "month":
			date = detail.MTime.Time().Month().String()
		default:
			date = detail.MTime.Time().Format("01-02")
		}
		trend.Dates = append(trend.Dates, date)
		trend.Coverages = append(trend.Coverages, detail.Coverage)
		trend.PassRates = append(trend.PassRates, detail.PassRate)
		trend.Scores = append(trend.Scores, detail.Score)
		trend.CommitIDs = append(trend.CommitIDs, detail.Cids)
	}
	return
}

//UTGernalCommit get singe or all users' general commit infos
func (s *Service) UTGernalCommit(c context.Context, commits string) (cmInfos []*ut.CommitInfo, err error) {
	if err = s.DB.Table(_utpkg).Joins("join ut_commit on ut_pkganls.commit_id = ut_commit.commit_id").
		Select("ROUND(AVG(coverage/100),2) as coverage,ROUND(SUM(passed)/SUM(assertions)*100,2) as pass_rate,ut_pkganls.merge_id,ut_pkganls.commit_id,ut_pkganls.mtime").
		Where("ut_pkganls.commit_id in (?) and (pkg!=substring_index(ut_pkganls.pkg,'/',5)", strings.Split(commits, ",")).Group("commit_id").Order("mtime DESC").Find(&cmInfos).Error; err != nil {
		log.Error("service.UTGernalCommit get cmInfos error(%v)", err)
		return
	}
	for i, cmInfo := range cmInfos {
		cmInfos[i].GitlabCommit, _ = s.dao.GitLabCommits(context.Background(), cmInfo.CommitID)
	}
	return
}

// CommentOnMR create or modify comment on relavent merge request
func (s *Service) CommentOnMR(c context.Context, projID, mrID int, msg string) (err error) {
	var (
		mr     = &ut.Merge{}
		noteID int
	)
	if err = s.DB.Where("merge_id=?", mrID).First(mr).Error; err != nil {
		log.Error("apmSvc.GitReport error search ut_merge error(%v)", err)
		return
	}
	if mr.NoteID == 0 {
		if noteID, err = s.gitlab.CreateMRNote(projID, mrID, msg); err != nil {
			log.Error("apmSvc.GitReport error CreateMRNote error(%v)", err)
			return
		}
		if err = s.DB.Table(_utMerge).Where("merge_id=?", mrID).Update("note_id", noteID).Error; err != nil {
			log.Error("apmSvc.GitReport error update ut_merge error(%v)", err)
			return
		}
	} else {
		if err = s.gitlab.UpdateMRNote(projID, mrID, mr.NoteID, msg); err != nil {
			log.Error("apmSvc.GitReport error CreateMRNote error(%v)", err)
			return
		}
	}
	return
}

// CommitHistory .
func (s *Service) CommitHistory(c context.Context, username string, times int64) (pkgs []*ut.PkgAnls, err error) {
	pkgs = make([]*ut.PkgAnls, 0)
	if err = s.DB.Table("ut_pkganls p").
		Select("p.merge_id, p.commit_id, group_concat( p.pkg ) AS pkg,group_concat( ROUND( p.coverage / 100, 2 ) ) AS coverages,group_concat( ROUND( p.passed / p.assertions * 100, 2 ) ) AS pass_rates,p.mtime").
		Joins("left join ut_commit c on c.commit_id=p.commit_id").Where(" c.username=? AND (p.pkg!=substring_index(p.pkg,'/',5) or p.pkg like 'go-common/library/%')", username).
		Group("p.commit_id").Order("p.mtime desc").Limit(times).Find(&pkgs).Error; err != nil {
		log.Error("s.CommitHistory query error(%v)", err)
	}
	return
}
