package service

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_inRuleRegexp    = "^\\d+((,)?\\d+)?$"
	_notInRuleRegexp = "^(\\d+)?$"
)

var (
	_emptyChannelRule = make([]*model.ChannelRule, 0)
)

// ChannelRuleCheckIDs channel rule check.
func (s *Service) ChannelRuleCheckIDs(c context.Context, in []int64, notin int64) (res *model.ChannelRule, err error) {
	var (
		inIDs         []int64
		notinID       int64
		notinRuleName string
		inRuleName    []string
		tagMap        = make(map[int64]*model.Tag, model.ChannelRuleMaxLen)
		allTids       = make([]int64, 0, model.ChannelRuleMaxLen)
	)
	allTids = append(allTids, in...)
	if notin > 0 {
		allTids = append(allTids, notin)
	}

	if _, tagMap, err = s.dao.Tags(c, allTids); err != nil {
		return
	}
	for _, tid := range in {
		k, ok := tagMap[tid]
		if !ok || k.State != model.StateNormal {
			err = ecode.TagNotExist
			return
		}
		inIDs = append(inIDs, k.ID)
		inRuleName = append(inRuleName, k.Name)
	}
	if notin > 0 {
		k, ok := tagMap[notin]
		if !ok || k.State != model.StateNormal {
			err = ecode.TagNotExist
			return
		}
		notinRuleName = k.Name
		notinID = k.ID
	}
	res = &model.ChannelRule{
		InRule:        xstr.JoinInts(inIDs),
		NotInRule:     strconv.FormatInt(notinID, 10),
		InRuleName:    strings.Join(inRuleName, ","),
		NotInRuleName: notinRuleName,
	}
	if notinRuleName != "" {
		res.Name = fmt.Sprintf("%s-%s", strings.Join(inRuleName, "+"), notinRuleName)
	} else {
		res.Name = strings.Join(inRuleName, "+")
	}
	return
}

// ChannelRuleCheckNames channel rule check.
func (s *Service) ChannelRuleCheckNames(c context.Context, in []string, notin string) (res *model.ChannelRule, err error) {
	var (
		inNames   []string
		notinName string
		inIDs     []int64
		notinID   int64
		tagMap    = make(map[string]*model.Tag, model.ChannelRuleMaxLen)
		allNames  = make([]string, 0, model.ChannelRuleMaxLen)
	)
	allNames = append(allNames, in...)
	if notin != "" {
		allNames = append(allNames, notin)
	}
	if _, _, tagMap, err = s.dao.TagByNames(c, allNames); err != nil {
		return
	}
	for _, name := range in {
		k, ok := tagMap[name]
		if !ok || k.State != model.StateNormal {
			err = ecode.TagNotExist
			return
		}
		inIDs = append(inIDs, k.ID)
		inNames = append(inNames, k.Name)
	}
	if notin != "" {
		k, ok := tagMap[notin]
		if !ok || k.State != model.StateNormal {
			err = ecode.TagNotExist
			return
		}
		notinID = k.ID
		notinName = k.Name
	}
	res = &model.ChannelRule{
		InRule:        xstr.JoinInts(inIDs),
		NotInRule:     strconv.FormatInt(notinID, 10),
		InRuleName:    strings.Join(inNames, ","),
		NotInRuleName: notinName,
	}
	if notinName != "" {
		res.Name = fmt.Sprintf("%s-%s", strings.Join(inNames, "+"), notinName)
	} else {
		res.Name = strings.Join(inNames, "+")
	}
	return
}

func ruleCheck(rule *model.ChannelRule) bool {
	if ok, _ := regexp.MatchString(_inRuleRegexp, rule.InRule); !ok {
		return false
	}
	if rule.NotInRule != "" {
		if ok, _ := regexp.MatchString(_notInRuleRegexp, rule.NotInRule); !ok {
			return false
		}
	}
	return true
}

// RuleChecks RuleChecks.
func (s *Service) RuleChecks(rules []*model.ChannelRule) (res []*model.ChannelRule, err error) {
	for _, rule := range rules {
		if !ruleCheck(rule) {
			err = ecode.ChanRuleCanotUse
			return
		}
		res = append(res, rule)
	}
	return
}

// ChannelRule ChannelRule.
func (s *Service) ChannelRule(c context.Context, tid int64) ([]*model.ChannelRule, error) {
	rules, tids, err := s.dao.ChannelRule(c, tid)
	if err != nil {
		return nil, err
	}
	if len(tids) == 0 {
		return _emptyChannelRule, nil
	}
	_, tagMap, err := s.dao.Tags(c, tids)
	if err != nil {
		return nil, err
	}
	var res []*model.ChannelRule
	for _, rule := range rules {
		if rule.InRule == "" {
			continue
		}
		var (
			inRule    []string
			notinRule string
		)
		if rule.NotInRule != "" {
			notinTid, err := strconv.ParseInt(rule.NotInRule, 10, 64)
			if err != nil {
				log.Error("get channel rule info strconv.ParseInt(%s) error(%v)", rule.NotInRule, err)
				err = nil
				continue
			}
			v, b := tagMap[notinTid]
			if notinTid != 0 && b {
				notinRule = v.Name
			}
		}
		ts, err := xstr.SplitInts(rule.InRule)
		if err != nil {
			log.Error("get channel rule info xstr.SplitInts(%s) error(%v)", rule.InRule, err)
			err = nil
			continue
		}
		for _, t := range ts {
			k, ok := tagMap[t]
			if t == 0 || !ok {
				continue
			}
			inRule = append(inRule, k.Name)
		}
		cr := &model.ChannelRule{
			ID:            rule.ID,
			Tid:           rule.Tid,
			InRule:        rule.InRule,
			NotInRule:     rule.NotInRule,
			InRuleName:    strings.Join(inRule, ","),
			NotInRuleName: notinRule,
			State:         rule.State,
			CTime:         rule.CTime,
			MTime:         rule.MTime,
			Editor:        rule.Editor,
		}
		if notinRule != "" {
			cr.Name = fmt.Sprintf("%s-%s", strings.Join(inRule, "+"), notinRule)
		} else {
			cr.Name = strings.Join(inRule, "+")
		}
		res = append(res, cr)
	}
	sort.Sort(model.ChannelRuleSort(res))
	return res, nil
}
