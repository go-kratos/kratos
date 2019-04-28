package service

import (
	"context"
	"regexp"
	"strings"

	"go-common/app/admin/main/filter/model"
	"go-common/library/ecode"
	"go-common/library/log"
)

func (s *Service) checkArea(ctx context.Context, areas []string) (err error) {
	for _, area := range areas {
		var areaInfo *model.Area
		if areaInfo, err = s.dao.AreaByName(ctx, area); err != nil {
			return
		}
		if areaInfo == nil {
			err = ecode.FilterIllegalArea
			return
		}
	}
	return
}

func (s *Service) checkReg(mode int8, rule string) (err error) {
	if mode == model.RegMode {
		if strings.Contains(rule, ".*") {
			err = ecode.FilterRegexpError1
			return
		}
		if strings.Contains(rule, "||") {
			err = ecode.FilterRegexpError2
			return
		}
		if _, err = regexp.Compile(rule); err != nil {
			log.Error("regexp.Compile(%s) err(%v)", rule, err)
			err = ecode.FilterIllegalRegexp
			return
		}
	}
	return
}

// checkWhiteSample 白样本检查，检查是否有大面积误伤正常内容
func (s *Service) checkWhiteSample(mode int8, rule string) (err error) {
	switch mode {
	case model.RegMode:
		var (
			reg      *regexp.Regexp
			hitCount = 0
		)
		if reg, err = regexp.Compile(rule); err != nil {
			log.Error("regexp.Compile(%s) err(%v)", rule, err)
			err = ecode.FilterIllegalRegexp
			return
		}
		for _, content := range s.conf.Property.NormalContents {
			if reg.MatchString(content) {
				hitCount++
			}
			if hitCount*100/len(s.conf.Property.NormalContents) >= s.conf.Property.NormalHitRate {
				err = ecode.FilterWhiteSampleHit
				return
			}
		}
	case model.StrMode:
		hitCount := 0
		for _, content := range s.conf.Property.NormalContents {
			if strings.Contains(content, rule) {
				hitCount++
			}
			if hitCount*100/len(s.conf.Property.NormalContents) >= s.conf.Property.NormalHitRate {
				err = ecode.FilterWhiteSampleHit
				return
			}
		}
	}
	return
}

// checkBlackSample 黑样本检查，检查是否有高危内容失效
func (s *Service) checkBlackSample(mode int8, rule string) (err error) {
	switch mode {
	case model.RegMode:
		var reg *regexp.Regexp
		if reg, err = regexp.Compile(rule); err != nil {
			log.Error("regexp.Compile(%s) err(%v)", rule, err)
			err = ecode.FilterIllegalRegexp
			return
		}
		for _, content := range s.conf.Property.RiskContents {
			if reg.FindString(content) == content {
				err = ecode.FilterBlackSampleHit
				log.Info("checkBlackSample find risk content [%s] hit by [%d:%s]", content, mode, rule)
				return
			}
		}
	case model.StrMode:
		for _, content := range s.conf.Property.RiskContents {
			if content == rule {
				err = ecode.FilterBlackSampleHit
				log.Info("checkBlackSample find risk content [%s] hit by [%d:%s]", content, mode, rule)
				return
			}
		}
		return
	}
	return
}
