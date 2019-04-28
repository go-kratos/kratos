package service

import (
	"context"
	"fmt"
	"go-common/app/interface/main/videoup/model/archive"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/net/metadata"
	"go-common/library/sync/errgroup"

	"github.com/microcosm-cc/bluemonday"
)

var (
	fieldMap = map[string]string{
		"title":   "稿件标题",
		"source":  "转载来源",
		"desc":    "稿件简介",
		"dynamic": "推荐语",
	}
)

// WebSingleFilter fn.
func (s *Service) WebSingleFilter(c context.Context, content string) (dt *archive.FilterData, hit []string, err error) {
	ip := metadata.String(c, metadata.RemoteIP)
	if dt, hit, err = s.filter.VideoFilter(c, content, ip); err != nil {
		log.Error("s.WebSingleFilter(%s) error(%v)", content, err)
		return
	}
	if len(hit) > 0 {
		log.Warn("s.WebSingleFilter(%s) hit(%v) ip(%s)", content, hit, ip)
	}
	return
}

// WebFilterArcParam fn.
func (s *Service) WebFilterArcParam(c context.Context, ap *archive.ArcParam, ip string) (code int, msg string, data map[string]interface{}) {
	res := make(map[string]string)
	data = make(map[string]interface{})
	g := &errgroup.Group{}
	ctx := context.TODO()
	// title
	g.Go(func() error {
		_, hit, _ := s.filter.VideoFilter(ctx, ap.Title, ip)
		if len(hit) > 0 {
			field := "title"
			fmtMsg := ecode.VideoupFieldFilterForbid.Message()
			res[field] = fmt.Sprintf(fmtMsg, fieldMap[field])
		}
		return nil
	})
	// source
	g.Go(func() error {
		_, hit, _ := s.filter.VideoFilter(ctx, ap.Source, ip)
		if len(hit) > 0 {
			field := "source"
			fmtMsg := ecode.VideoupFieldFilterForbid.Message()
			res[field] = fmt.Sprintf(fmtMsg, fieldMap[field])
		}
		return nil
	})
	// desc
	g.Go(func() error {
		_, hit, _ := s.filter.VideoFilter(ctx, ap.Desc, ip)
		if len(hit) > 0 {
			field := "desc"
			fmtMsg := ecode.VideoupFieldFilterForbid.Message()
			res[field] = fmt.Sprintf(fmtMsg, fieldMap[field])
		}
		return nil
	})
	// dynamic
	g.Go(func() error {
		_, hit, _ := s.filter.VideoFilter(ctx, ap.Dynamic, ip)
		if len(hit) > 0 {
			field := "dynamic"
			fmtMsg := ecode.VideoupFieldFilterForbid.Message()
			res[field] = fmt.Sprintf(fmtMsg, fieldMap[field])
		}
		return nil
	})
	g.Wait()
	if len(res) != 0 {
		code = ecode.VideoupFieldFilterForbid.Code()
		msg = "当前输入包含敏感信息,请修正"
		data["info"] = res
		log.Warn("WebFilterArcParam res(%+v),ap(%+v),ip(%s)", res, ap, ip)
	}
	return
}

// xss filter
func xssFilter(content string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(content)
}
