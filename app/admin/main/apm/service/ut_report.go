package service

import (
	"context"
	"fmt"
	"go-common/app/admin/main/apm/model/ut"
	"go-common/library/log"
	"sort"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// GitReport post simple report using SAGA account
func (s *Service) GitReport(c context.Context, projID int, mrID int, commitID string) (err error) {
	var (
		pkgs = make([]*ut.PkgAnls, 0)
		row  = `<tr><td colspan="2"><a href="http://sven.bilibili.co/#/ut/detail?commit_id=%s&pkg=%s">%s</a></td><td>%.2f</td><td>%.2f</td><td>%.2f</td><td style="text-align: center">%s</td></tr>`
		msg  = fmt.Sprintf(`<pre><h4>单元测试速报</h4>(Commit:<a href="http://sven.bilibili.co/#/ut?merge_id=%d&pn=1&ps=20">%s</a>)</pre>`, mrID, commitID) + `<table><thead><tr><th colspan="2">包名</th><th>覆盖率(%%)</th><th>通过率(%%)</th><th>覆盖率较历史最高(%%)</th><th>是否合格</th></tr></thead><tbody>%s</tbody>%s</table>`
		rows string
		root = ""
		file = &ut.File{}
	)
	if err = s.DB.Where("commit_id=? AND (pkg!=substring_index(pkg, '/', 5) OR pkg like 'go-common/library/%')", commitID).Find(&pkgs).Error; err != nil {
		log.Error("apmSvc.GitReport query error(%v)", err)
		return
	}
	for _, pkg := range pkgs {
		t, _ := s.tyrant(pkg)
		app := pkg.PKG
		if !strings.Contains(pkg.PKG, "/library") && len(pkg.PKG) >= 5 {
			app = strings.Join(strings.Split(pkg.PKG, "/")[:5], "/")
		}
		rows += fmt.Sprintf(row, pkg.CommitID, app, pkg.PKG, pkg.Coverage/100, t.PassRate, t.Increase, "%s")
		if t.Tyrant {
			rows = fmt.Sprintf(rows, "❌")
		} else {
			rows = fmt.Sprintf(rows, "✔️")
		}
	}
	if err = s.DB.Select(`count(id) as id, commit_id, sum(statements) as statements, sum(covered_statements) as covered_statements`).
		Where(`pkg!=substring_index(pkg, "/", 5) OR ut_pkganls.pkg like 'go-common/library/%'`).Group(`commit_id`).Having("commit_id=?", commitID).First(file).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("apmSvc.GitReport query error(%v)", err)
			return
		}
		err = nil
	} else {
		root = fmt.Sprintf(`<tfoot><tr><td><b>总覆盖率: </b>%.2f%%<br><b>总检测文件数: </b>%d</br></td><td><b>Tracked语句数: </b>%d<br><b>覆盖的语句数: </b>%d</br></td><td colspan="4" align="center">总覆盖率 = 覆盖语句数 / Tracked语句数</td></tr></tfoot>`,
			float64(file.CoveredStatements)/float64(file.Statements)*100, file.ID, file.Statements, file.CoveredStatements)
	}
	if err = s.CommentOnMR(c, projID, mrID, fmt.Sprintf(msg, rows, root)); err != nil {
		log.Error("apmSvc.GitReport call CommentOnMR error(%v)", err)
		return
	}
	return
}

// WechatReport send wechat msg to a group when mr is merged
func (s *Service) WechatReport(c context.Context, mrid int64, cid, src, des string) (err error) {
	var (
		pkgs []*ut.PkgAnls
		mr   = &ut.Merge{}
		foot = fmt.Sprintf("\n*覆盖率增长:该包本次合并最后一次commit与过往已合并记录中最大覆盖率进行比较\nMR:http://git.bilibili.co/platform/go-common/merge_requests/%d\n单测报告:http://sven.bilibili.co/#/ut?merge_id=%d&pn=1&ps=20\n", mrid, mrid)
	)
	if err = s.dao.DB.Where("merge_id=? AND is_merged=?", mrid, 1).First(mr).Error; err != nil {
		log.Error("apmSvc.WechatReport Error(%v)", err)
		return
	}
	msg := "【测试姬】新鲜的单测速报出炉啦ᕕ( ᐛ )ᕗ\n\n" + fmt.Sprintf("由 %s 发起的 MR (%s->%s) 合并成功！\n\n", mr.UserName, src, des)
	if err = s.dao.DB.Where("commit_id=? AND (pkg!=substring_index(pkg, '/', 5) OR ut_pkganls.pkg like 'go-common/library/%')", cid).Find(&pkgs).Error; err != nil || len(pkgs) == 0 {
		log.Error("apmSvc.WechatReport Error(%v)", err)
		return
	}
	for _, pkg := range pkgs {
		var (
			arrow    = ""
			maxPkgs  = make([]ut.PkgAnls, 0)
			increase = pkg.Coverage / 100
			file     = &ut.File{}
			lastFile = &ut.File{}
		)
		if err = s.DB.Table("ut_pkganls").Joins("left join ut_merge ON ut_merge.merge_id=ut_pkganls.merge_id").Select("ut_pkganls.commit_id, ut_pkganls.coverage").Where("ut_pkganls.pkg=? AND ut_pkganls.merge_id!=? AND ut_merge.is_merged=1", pkg.PKG, mrid).Order("coverage desc").Find(&maxPkgs).Error; err != nil {
			log.Error("apmSvc.WechatReport error(%v)", err)
			return
		}
		if len(maxPkgs) != 0 {
			increase = pkg.Coverage/100 - maxPkgs[0].Coverage/100
		}
		if increase < float64(0) {
			arrow = "⬇️"
			if !strings.Contains(foot, "本次合并后有包覆盖率下降了喔~还需要再加油鸭💪~") {
				foot += "\n本次合并后有包覆盖率下降了喔~还需要再加油鸭💪~"
			}
		} else if increase > float64(0) {
			arrow = "⬆️"
		}
		msg += fmt.Sprintf("*%s\n\t覆盖率：%.2f%%\t覆盖率增长：%.2f%%\t%s\n", pkg.PKG, pkg.Coverage/100, increase, arrow)
		if err = s.DB.Select("count(id) as id,commit_id,pkg,sum(statements) as statements,sum(covered_statements) as covered_statements").Group("commit_id,pkg").Having("commit_id=? AND pkg=?", cid, pkg.PKG).First(file).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
				continue
			}
			log.Error("apmSvc.WechatReport error(%v)", err)
			return
		}
		msg += fmt.Sprintf("\tTracked语句数：%d\t覆盖语句数：%d\n", file.Statements, file.CoveredStatements)
		if err = s.DB.Table("ut_file").Joins("left join ut_commit on ut_commit.commit_id = ut_file.commit_id left join ut_merge on ut_merge.merge_id =ut_commit.merge_id").Select("ut_file.commit_id,ut_file.pkg,count(ut_file.id) as id,sum(statements) as statements, sum(covered_statements) as covered_statements,ut_file.mtime").Where("is_merged=1 and ut_merge.merge_id!=?", mrid).Group("commit_id,pkg").Having("pkg=?", pkg.PKG).Order("mtime desc").First(lastFile).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = nil
				continue
			}
			log.Error("apmSvc.WechatReport error(%v)", err)
			return
		}
		msg += fmt.Sprintf("\tTracked语句增长：%d\t覆盖语句增长：%d\n", file.Statements-lastFile.Statements, file.CoveredStatements-lastFile.CoveredStatements)
	}
	if err = s.dao.SendWechatToGroup(c, s.c.WeChat.ChatID, msg+foot); err != nil {
		log.Error("apmSvc.WechatReport Error(%v)", err)
		return
	}
	return
}

// RankWechatReport send rank to wechat group 19:00 everyday
func (s *Service) RankWechatReport(c context.Context) (err error) {
	var (
		topRanks []*ut.RankResp
		msg      = "【测试姬】今日份的单测榜单🎉\n\n"
		root     = "更多详情请戳：http://sven.bilibili.co/#/ut/leaderboard\n感谢辛勤工作的一天的你☕️一起快乐下班吧☕️\n"
	)
	if topRanks, err = s.RankTen(c, "desc"); err != nil {
		log.Error("apmSvc.RankWechatReport Error(%v)", err)
		return
	}
	msg += "🔘TOP 10\n"
	for i, r := range topRanks {
		msg += fmt.Sprintf("%d: %s\t分数: %.2f\n", i+1, r.UserName, r.Score)
	}
	msg += "恭喜以上各位~还请继续保持哟(๑•̀ㅂ•́)و✧\n\n"
	if err = s.dao.SendWechatToGroup(c, s.c.WeChat.ChatID, msg+root); err != nil {
		log.Error("apmSvc.RankWechatReport Error(%v)", err)
		return
	}
	return
}

// SummaryWechatReport send depts' summary to ut wechat group every Friday(19:00)
func (s *Service) SummaryWechatReport(c context.Context) (err error) {
	var (
		depts               = `"main","ep","openplatform","live","video","bbq"`
		msg                 = fmt.Sprintf("【测试姬】Kratos大仓库单元测试周刊 ( %s )\n\n>> 接入情况: \n", time.Now().Format("2006-01-02"))
		covMsg              = ">> 覆盖情况: \n"
		rankMsg             = ">> 本周Top 3应用:\n"
		root                = "\n\n更多汇总信息请姥爷们查看: http://sven.bilibili.co/#/ut/dashboard?tab=project\n更多文档信息: http://info.bilibili.co/pages/viewpage.action?pageId=6947230\n吐槽建议: http://sven.bilibili.co/#/suggestion/list\n"
		tpAcc               = "\t* %s\t接入数: %d\t总数: %d\t接入率: %.2f%%\n"
		tpCov               = "\t* %s\t应用平均覆盖率: %.2f%%\t增长率：%.2f%%\n"
		tpProjRank          = "\t%d. %s \t包平均覆盖率: %.2f%%\n"
		sumTotal, sumAccess int64
		sumCov              float64
	)
	s.appsCache.Lock()
	defer s.appsCache.Unlock()
	for _, dept := range s.appsCache.Dept {
		if !strings.Contains(depts, "\""+dept.Name+"\"") {
			continue
		}
		sumTotal += dept.Total
		sumAccess += dept.Access
		msg += fmt.Sprintf(tpAcc, dept.Name, dept.Access, dept.Total, float64(dept.Access)/float64(dept.Total)*100)
		if dept.Coverage == float64(0) {
			continue
		}
		var preCoverage float64
		if preCoverage, err = s.dao.GetAppCovCache(c, dept.Name); err != nil {
			log.Error("service.SummaryWechatReport GetAppCov Error(%v)", err)
			return
		}
		covMsg += fmt.Sprintf(tpCov, dept.Name, dept.Coverage/100, (dept.Coverage-preCoverage)/100)
	}
	sort.Slice(s.appsCache.Slice, func(i, j int) bool { return s.appsCache.Slice[i].Coverage > s.appsCache.Slice[j].Coverage })
	for k, slice := range s.appsCache.Slice {
		if k <= 2 {
			rankMsg += fmt.Sprintf(tpProjRank, k+1, slice.Path, slice.Coverage/100)
		}
		if slice.HasUt == 1 {
			sumCov += slice.Coverage
		}
	}
	msg += fmt.Sprintf("\t总接入情况 - 接入应用数: %d，应用总数: %d\n\n", sumAccess, sumTotal) + covMsg + fmt.Sprintf("\t总覆盖情况 - 应用平均覆盖率: %.2f%%\n\n", sumCov/float64(sumAccess)/100) + rankMsg
	if err = s.dao.SendWechatToGroup(c, s.c.WeChat.ChatID, msg+root); err != nil {
		log.Error("apmSvc.WechatReport Error(%v)", err)
		return
	}
	return
}
