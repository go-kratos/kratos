package service

import (
	"context"
	"regexp"
	"strings"
	"unicode/utf8"

	"go-common/app/job/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

var (
	regHTML      = regexp.MustCompile(`(?i)\<.*?(script|href|img|src)+.*?\>`)
	regSymbol    = regexp.MustCompile(`^[g\pP|g\pS]+$`)
	regZeroWidth = regexp.MustCompile(`[\x{200b}]+`)
)

// CheckName check tag name .
func (s *Service) checkName(name string) (dst string, err error) {
	if !utf8.ValidString(name) {
		log.Error("utf8.ValidString(%s)", name)
		err = ecode.RequestErr
		return
	}
	index := regHTML.FindAllString(name, -1)
	if len(index) > 0 {
		err = ecode.RequestErr
		return
	}
	dst = regZeroWidth.ReplaceAllString(name, "")
	dst = replace(strings.TrimSpace(dst))
	if dst == "" || len([]rune(dst)) > model.TNameMaxLen {
		log.Error("name == nil or length > max length: %s", name)
		err = ecode.RequestErr
		return
	}
	if regSymbol.MatchString(dst) {
		log.Error("cant not contain continuous symbol(%s)", name)
		err = ecode.RequestErr
	}
	return
}

func replace(name string) string {
	var (
		there bool
		rb    []byte
	)
	sb := []byte(name)
	for _, b := range sb {
		if b < 0x20 || b == 0x7f {
			there = true
			continue
		}
		rb = append(rb, b)
	}
	if there {
		log.Warn("There are invisible characters,tag Name(%v)", sb)
	}
	return string(rb)
}

// TODO 插入再查？
func (s *Service) createTags(c context.Context, names []string) (res []*model.Tag, err error) {
	if names, err = s.dao.MFilter(c, names); err != nil {
		return
	}
	if len(names) == 0 {
		return
	}
	tags := make([]*model.Tag, 0, len(names))
	for _, v := range names {
		tags = append(tags, &model.Tag{
			Name:  v,
			State: model.TagStateNormal,
		})
	}
	if err = s.dao.InsertTags(c, tags); err != nil {
		return
	}
	res, _, err = s.dao.TagByNames(c, names)
	return
}
