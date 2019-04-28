package service

import (
	"context"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"go-common/app/interface/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

// ResTags 资源下tag列表
func (s *Service) ResTags(c context.Context, oids []int64, mid int64, typ int8) (tm map[int64][]*model.Tag, err error) {
	tm = make(map[int64][]*model.Tag, len(oids))
	for _, oid := range oids {
		var ts []*model.Tag
		ts, _, err = s.resTagsService(c, oid, mid, int32(typ))
		if err != nil {
			continue
		}
		tm[oid] = ts
	}
	return
}

// UpResBind .
func (s *Service) UpResBind(c context.Context, oid, mid int64, tNames []string, typ int8, now time.Time) (err error) {
	var (
		tids []int64
		tMap map[string]int64
	)
	_, tMap, err = s.addNewTags(c, 0, tNames, now)
	if err != nil {
		return
	}
	for _, v := range tMap {
		tids = append(tids, v)
	}
	err = s.platformUpBind(c, oid, mid, tids, int32(typ))
	return
}

// ResAdminBind .
func (s *Service) ResAdminBind(c context.Context, oid, mid int64, tNames []string, typ int8, now time.Time) (err error) {
	var (
		tids []int64
		tMap map[string]int64
	)
	_, tMap, err = s.addNewTags(c, 0, tNames, now)
	if err != nil {
		return
	}
	for _, v := range tMap {
		tids = append(tids, v)
	}
	err = s.platformAdminBind(c, oid, mid, tids, int32(typ))
	return
}

var (
	regHTML      = regexp.MustCompile(`(?i)\<.*?(script|href|img|src)+.*?\>`)
	regSymbol    = regexp.MustCompile(`^[g\pP|g\pS]+$`)
	regZeroWidth = regexp.MustCompile(`[\x{200b}]+`)
)

// CheckName check tag name .
func (s *Service) CheckName(name string) (dst string, err error) {
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
	dst = replace(dst)
	dst = strings.TrimSpace(dst)
	if dst == "" || len([]rune(dst)) > model.TnameMaxLen {
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
