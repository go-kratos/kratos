package service

import (
	"context"
	"regexp"
	"strings"
	"time"

	search "go-common/app/interface/openplatform/article/model/search"
	"go-common/library/log"
)

const _sourceType = "article"

// Segment .
func (s *Service) Segment(c context.Context, id int32, content string, withTag int, remarks string) (keywords []string, err error) {
	var (
		source  = _sourceType
		trackid = int32(time.Now().Unix())
		res     *search.TagboxResponse
	)
	if withTag == 1 {
		content = strings.Replace(content, "&nbsp;", " ", -1)
		rule := "\\<[\\S\\s]+?\\>"
		reg, _ := regexp.Compile(rule)
		content = reg.ReplaceAllString(content, "")
	}
	req := &search.TagboxRequest{
		Id:         &id,
		SourceType: &source,
		Content:    &content,
		Trackid:    &trackid,
		Remarks:    &remarks,
	}
	if res, err = s.searchRPC.Segment(c, req); err != nil {
		log.Error("s.Segment error(%+v), params(%+v)", err, req)
		return
	}
	if *res.ExecCode != int32(0) {
		log.Error("creation: s.segment id(%d), code(%d)", id, res.ExecCode)
		return
	}
	keywords = res.GetKeywords()
	return
}
