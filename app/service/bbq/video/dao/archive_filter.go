package dao

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/service/bbq/video/conf"
	"go-common/app/service/bbq/video/model"
	mainAccount "go-common/app/service/main/account/api"
	mainArchive "go-common/app/service/main/archive/api"
	"go-common/library/log"
)

// Filter .
type Filter interface {
	doFilter(*model.Archive) bool
}

// ArchiveFilters .
type ArchiveFilters struct {
	filter []Filter
}

// NewArchiveFilters .
func NewArchiveFilters(f ...Filter) *ArchiveFilters {
	return &ArchiveFilters{
		filter: f,
	}
}

// DoFilter 过滤稿件
func (af *ArchiveFilters) DoFilter(a *model.Archive) (res bool) {
	res = false
	for _, v := range af.filter {
		res = res || v.doFilter(a)
		if res {
			return
		}
	}
	return
}

// ArchiveDetailFilter 稿件详情过滤
type ArchiveDetailFilter struct {
	rule        *conf.ArchiveRule
	archiveType map[int32]*mainArchive.Tp
}

func (f *ArchiveDetailFilter) doFilter(a *model.Archive) bool {
	c := context.Background()
	// 过滤状态
	if a.State != f.rule.State {
		log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter state want[%d] actually[%d]", f.rule.State, a.State)))
		return true
	}

	// 过滤转载
	if a.Copyright != 1 {
		log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter copyright want[1] actually[%d]", a.Copyright)))
		return true
	}

	// 过滤access
	if a.Access == f.rule.NotAccess {
		log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter access not[%d] actually[%d]", f.rule.NotAccess, a.Access)))
		return true
	}

	// 过滤TID
	if subT, ok := f.archiveType[a.TypeID]; ok {
		a.TID = subT.Pid
		for _, t := range f.rule.TID {
			if a.TID == t {
				log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter TID not[%d] actually[%d]", t, a.TID)))
				return true
			}
		}
	}

	// 过滤SubTID
	a.SubTID = a.TypeID
	for _, t := range f.rule.SubTID {
		if a.TypeID == t {
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter TID not[%d] actually[%d]", t, a.TypeID)))
			return true
		}
	}

	// 过滤Title
	for _, v := range f.rule.Titles {
		if strings.Contains(a.Title, v) {
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter title[%s] contains[%s]", a.Title, v)))
			return true
		}
	}

	// 过滤Content
	for _, v := range f.rule.Contents {
		if strings.Contains(a.Content, v) {
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter content[%s] contains[%s]", a.Content, v)))
			return true
		}
	}

	return false
}

// ArchiveUpFilter 稿件详情过滤
type ArchiveUpFilter struct {
	rule          *conf.UPRule
	accountClient mainAccount.AccountClient
}

func (f *ArchiveUpFilter) doFilter(a *model.Archive) bool {
	c := context.Background()
	// 过滤UP主mid
	for _, v := range f.rule.MID {
		if a.MID == v {
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter mid[%d] contains[%d]", a.MID, v)))
			return true
		}
	}

	response, err := f.accountClient.Card3(c, &mainAccount.MidReq{Mid: a.MID})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter err[%v]", err)))
		return true
	}
	// 过滤UP主uname
	for _, v := range f.rule.UName {
		if strings.Contains(response.Card.Name, v) {
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchiveUpFilter doFilter uname[%s] contains[%s]", response.Card.Name, v)))
			return true
		}
	}

	return false
}

// ArchivePageFilter 稿件尺寸过滤
type ArchivePageFilter struct {
	rule          *conf.PageRule
	archiveClient mainArchive.ArchiveClient
}

func (f *ArchivePageFilter) doFilter(a *model.Archive) bool {
	c := context.Background()
	response, err := f.archiveClient.View(c, &mainArchive.ViewRequest{
		Aid: a.AID,
	})
	if err != nil {
		log.Errorv(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter err[%v]", err)))
		return true
	}

	for _, page := range response.Pages {
		// Rotate 不准，暂时不用
		// if page.Dimension.Rotate != 1 {
		// 	// 必须是竖屏
		// 	continue
		// }
		if page.Dimension.Width >= page.Dimension.Height {
			// 高度必须大于宽度
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter arc[%v] page[%v] filter width >= height", *(response.Arc), *page)))
			continue
		}
		yx := float32(page.Dimension.Height) / float32(page.Dimension.Width)
		if yx < f.rule.MinYX || yx > f.rule.MaxYX {
			// 高宽比大于 1.75 小于 1.80
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter arc[%v] page[%v] filter want[%v, %v] actually[%v]", *(response.Arc), *page, f.rule.MinYX, f.rule.MaxYX, yx)))
			continue
		}
		if page.Dimension.Width < f.rule.MinX {
			// 不低于720p
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter arc[%v] page[%v] filter want[%v] actually[%v]", *(response.Arc), *page, f.rule.MinX, page.Dimension.Width)))
			continue
		}
		if page.Duration <= f.rule.MinDuration || page.Duration > f.rule.MaxDuration {
			// 过滤时长,大于15s 小于60s
			log.Infov(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter arc[%v] page[%v] filter want(%v, %v) actually[%v]", *(response.Arc), *page, f.rule.MinDuration, f.rule.MaxDuration, page.Duration)))
			continue
		}
		a.CID = page.Cid
	}
	log.Infov(c, log.KV("log", fmt.Sprintf("ArchivePageFilter doFilter arc[%v]", *(a))))

	return a.CID == 0
}
