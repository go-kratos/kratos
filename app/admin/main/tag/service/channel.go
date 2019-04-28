package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"time"

	"go-common/app/admin/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	xtime "go-common/library/time"
	"go-common/library/xstr"
)

// ChannelsByCategory channels by channel_type.
func (s *Service) ChannelsByCategory(c context.Context, id int64, name string) (offline, online []*model.Channel, err error) {
	var (
		tids     []int64
		category *model.ChannelCategory
		channels []*model.Channel
		tagMap   map[int64]*model.Tag
	)
	if name != "" {
		if category, err = s.dao.ChannelCategoryByName(c, name); err != nil {
			return
		}
		if category == nil {
			err = ecode.ChanTypeNotExist
			return
		}
		id = category.ID
	}
	if channels, tids, err = s.dao.ChannelsByType(c, id); err != nil {
		return
	}
	if len(tids) > 0 {
		if _, tagMap, err = s.dao.Tags(c, tids); err != nil {
			return
		}
	}
	offline = make([]*model.Channel, 0)
	online = make([]*model.Channel, 0)
	sort.Sort(model.ChannelSort(channels))
	for _, channel := range channels {
		k, ok := tagMap[channel.ID]
		if !ok {
			continue
		}
		channel.Name = k.Name
		channel.Cover = k.Cover
		channel.Content = k.Content
		channel.INTShield = channel.AttrVal(model.ChannelAttrINT)
		switch channel.State {
		case model.ChanStateOffline:
			offline = append(offline, channel)
		case model.ChanStateCommon:
			online = append(online, channel)
		case model.ChanStateRecomend:
			channel.Top = channel.AttrVal(model.ChannelAttrTop)
			online = append(online, channel)
		default:
		}
	}
	return
}

// AllChannels get all channels.
func (s *Service) AllChannels(c context.Context) (res []*model.Channel, err error) {
	var tids []int64
	if res, tids, err = s.dao.ChannelAll(c); err != nil {
		return
	}
	_, tagMap, err := s.dao.Tags(c, tids)
	if err != nil {
		return
	}
	for _, channel := range res {
		if k, ok := tagMap[channel.ID]; ok {
			channel.Name = k.Name
			channel.CheckBack = channel.AttrVal(model.ChannelAttrCheckBack)
			channel.Activity = channel.AttrVal(model.ChannelAttrActivity)
			channel.INTShield = channel.AttrVal(model.ChannelAttrINT)
		}
	}
	sort.Sort(model.ChannelSort(res))
	return
}

// ChanneList ChanneList.
func (s *Service) ChanneList(c context.Context, param *model.ParamChanneList) (res []*model.Channel, total int32, err error) {
	var (
		sql      []string
		channels []*model.Channel
		tids     []int64
	)
	start := (param.Pn - 1) * param.Ps
	end := param.Ps
	if len(param.IDs) > 0 {
		sql = append(sql, fmt.Sprintf("tid in (%s)", xstr.JoinInts(param.IDs)))
	}
	if param.Operator != "" {
		sql = append(sql, fmt.Sprintf("operator=%q", param.Operator))
	}
	if param.STime != "" {
		sql = append(sql, fmt.Sprintf("ctime>=%q", param.STime))
	}
	if param.ETime != "" {
		sql = append(sql, fmt.Sprintf("ctime<=%q", param.ETime))
	}
	if param.Type >= 0 {
		sql = append(sql, fmt.Sprintf("type=%d", param.Type))
	}
	if param.State >= 0 {
		sql = append(sql, fmt.Sprintf("state=%d", param.State))
	}
	if channels, tids, err = s.dao.ChanneList(c, sql, param.Order, param.Sort, start, end); err != nil {
		return
	}
	total, _ = s.dao.CountChanneList(c, sql)
	tagCountMap, _ := s.dao.TagCounts(c, tids)
	_, tagMap, _ := s.dao.Tags(c, tids)
	countMap, _ := s.channelRuleTagNum(c, tids)
	for _, channel := range channels {
		if param.INTShield >= 0 && channel.AttrVal(model.ChannelAttrINT) != param.INTShield {
			continue
		}
		if k, ok := tagCountMap[channel.ID]; ok {
			channel.Count = k
		}
		if t, ok := tagMap[channel.ID]; ok {
			channel.Name = t.Name
			channel.Cover = t.Cover
			channel.Content = t.Content
			channel.CheckBack = channel.AttrVal(model.ChannelAttrCheckBack)
			channel.Activity = channel.AttrVal(model.ChannelAttrActivity)
			channel.INTShield = channel.AttrVal(model.ChannelAttrINT)
		}
		if v, ok := countMap[channel.ID]; ok {
			channel.TagNums = v
		}
		res = append(res, channel)
	}
	return
}

func (s *Service) channelRuleTagNum(c context.Context, tids []int64) (numMap map[int64]int32, err error) {
	chanRuleMap, err := s.dao.ChannelRules(c, tids)
	if err != nil {
		return nil, err
	}
	numMap = make(map[int64]int32)
	for tid, chanRules := range chanRuleMap {
		nameMap := make(map[int64]int)
		var count int32
		for _, chanRule := range chanRules {
			if chanRule.State != model.ChanRuleNormal {
				continue
			}
			if !ruleCheck(chanRule) {
				continue
			}
			tidB, _ := strconv.ParseInt(chanRule.NotInRule, 10, 64)
			tids, _ := xstr.SplitInts(chanRule.InRule)
			tids = append(tids, tidB)
			for _, id := range tids {
				if id == 0 {
					continue
				}
				if k, ok := nameMap[id]; !ok || k == 0 {
					nameMap[id] = 1
					count++
				}
			}
		}
		numMap[tid] = count
	}
	return
}

// ChannelInfo ChannelInfo.
func (s *Service) ChannelInfo(c context.Context, tid int64, tname string) (res *model.ChannelInfo, err error) {
	var (
		channel  *model.Channel
		tag      *model.Tag
		tagMap   map[int64]*model.Tag
		synonyms = make([]*model.ChannelSynonym, 0, model.ChannelSynonymMax)
	)
	if tid <= 0 {
		tag, err = s.dao.TagByName(c, tname)
	} else {
		tag, err = s.dao.Tag(c, tid)
	}
	if err != nil {
		return
	}
	if tag == nil || tag.State != model.StateNormal {
		err = ecode.ChannelNotExist
		return
	}
	if channel, err = s.dao.Channel(c, tag.ID); err != nil {
		return
	}
	if channel == nil {
		err = ecode.ChannelNotExist
		return
	}
	synonymMap, tids, err := s.dao.ChannelSynonymMap(c, tag.ID)
	if err != nil {
		return
	}
	if len(tids) > 0 {
		if _, tagMap, err = s.dao.Tags(c, tids); err != nil {
			return
		}
	}
	for _, synonym := range synonymMap {
		if synonym.State != model.StateNormal {
			continue
		}
		k, ok := tagMap[synonym.Tid]
		if !ok || k == nil || k.State != model.StateNormal {
			continue
		}
		synonym.TName = k.Name
		synonyms = append(synonyms, synonym)
	}
	sort.Sort(model.ChannelSynonymSort(synonyms))
	res = &model.ChannelInfo{
		ID:           channel.ID,
		Name:         tag.Name,
		Type:         channel.Type,
		Rank:         channel.Rank,
		Operator:     channel.Operator,
		Cover:        tag.Cover,
		HeadCover:    tag.HeadCover,
		ShortContent: tag.ShortContent,
		Content:      tag.Content,
		CheckBack:    channel.AttrVal(model.ChannelAttrCheckBack),
		State:        channel.State,
		Rules:        _emptyChannelRule,
		Activity:     channel.AttrVal(model.ChannelAttrActivity),
		INTShield:    channel.AttrVal(model.ChannelAttrINT),
		Synonyms:     synonyms,
	}
	res.Rules, err = s.ChannelRule(c, tag.ID)
	return
}

// DeleteChannel delete channel.
func (s *Service) DeleteChannel(c context.Context, tid int64, uname string) (err error) {
	var channel *model.Channel
	if channel, err = s.dao.Channel(c, tid); err != nil {
		return
	}
	if channel == nil {
		return ecode.ChannelNotExist
	}
	if channel.State == model.ChanStateStop {
		return ecode.ChannelNoChange
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("channel deleet tran error(%v)", err)
		return
	}
	channel.State = model.ChanStateStop
	channel.Editor = uname
	if _, err = s.dao.TxUpChannel(tx, channel); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpChannelRuleState(tx, tid, model.ChanRuleDelete, uname); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpStateChannelSynonym(tx, tid, model.ChanGroupDelete, uname); err != nil {
		tx.Rollback()
		return
	}
	err = tx.Commit()
	return
}

// StateChannel .
func (s *Service) StateChannel(c context.Context, tid int64, state int32, uname string) (err error) {
	var (
		channel  *model.Channel
		category *model.ChannelCategory
	)
	if channel, err = s.dao.Channel(c, tid); err != nil {
		return
	}
	if channel == nil {
		return ecode.ChannelNotExist
	}
	if channel.State == state {
		return ecode.ChannelNoChange
	}
	if category, err = s.dao.ChannelCategory(c, channel.Type); err != nil {
		return
	}
	channel.State = state
	channel.Editor = uname
	if category.AttrVal(model.ChannelCategoryAttrINT) == model.CategoryStateShieldINT && channel.AttrVal(model.ChannelAttrINT) != model.ChannelStateShieldINT {
		channel.AttrSet(model.ChannelAttrINT, model.ChannelStateShieldINT)
	}
	switch state {
	case model.ChanStateCommon, model.ChanStateOffline:
		channel.AttrSet(model.ChannelAttrTop, model.ChannelTopNo)
	default:
	}
	_, err = s.dao.UpdateChannel(c, channel)
	return
}

// MigrateChannel migrate channel from a type to b type.
func (s *Service) MigrateChannel(c context.Context, tid, tp int64, tname, uname string) (err error) {
	var (
		tag      *model.Tag
		channel  *model.Channel
		category *model.ChannelCategory
	)
	if tid > 0 {
		tag, err = s.dao.Tag(c, tid)
	} else {
		tag, err = s.dao.TagByName(c, tname)
	}
	if err != nil {
		return
	}
	if tag == nil || tag.State == model.StateDel {
		return ecode.TagNotExist
	}
	if tag.State == model.StateShield {
		return ecode.TagAlreadyShield
	}
	if channel, err = s.dao.Channel(c, tag.ID); err != nil {
		return
	}
	if channel == nil || channel.State < model.ChanStateCommon {
		return ecode.ChannelNotExist
	}
	if channel.Type == tp {
		return ecode.ChannelAleadyMigrated
	}
	if category, err = s.dao.ChannelCategory(c, tp); err != nil {
		return
	}
	if category == nil || category.State == model.StateDel {
		return ecode.ChanTypeNotExist
	}
	channel.Rank, _ = s.dao.CountChannelByType(c, tp)
	channel.Editor = uname
	channel.Type = tp
	if category.AttrVal(model.ChannelCategoryAttrINT) == model.CategoryStateShieldINT {
		channel.AttrSet(model.ChannelAttrINT, model.ChannelStateShieldINT)
	}
	_, err = s.dao.UpdateChannel(c, channel)
	return
}

// SortChannels sort channels.
func (s *Service) SortChannels(c context.Context, tp int64, tids []int64) (err error) {
	var (
		category *model.ChannelCategory
		channels []*model.Channel
		tidMap   = make(map[int64]int32, len(tids))
	)
	if category, err = s.dao.ChannelCategory(c, tp); err != nil {
		return
	}
	if category == nil || category.State == model.StateDel {
		return ecode.ChanTypeNotExist
	}
	if channels, _, err = s.dao.ChannelsByType(c, tp); err != nil {
		return
	}
	for index, tid := range tids {
		tidMap[tid] = int32(index)
	}
	for _, channel := range channels {
		k, ok := tidMap[channel.ID]
		if !ok {
			continue
		}
		channel.Rank = k
	}
	_, err = s.dao.UpdateChannels(c, channels)
	return
}

// RecommandChannels recommand channels.
func (s *Service) RecommandChannels(c context.Context) (res []*model.Channel, err error) {
	res = make([]*model.Channel, 0)
	channels, tids, err := s.dao.RecommandChannel(c)
	if err != nil {
		return
	}
	if len(tids) == 0 {
		return
	}
	_, tagMap, err := s.dao.Tags(c, tids)
	if err != nil {
		return
	}
	for _, channel := range channels {
		tag, ok := tagMap[channel.ID]
		if !ok || tag == nil || tag.State != model.StateNormal {
			continue
		}
		channel.Top = channel.AttrVal(model.ChannelAttrTop)
		channel.Name = tag.Name
		res = append(res, channel)
	}
	sort.Sort(model.RecommendChannelSort(res))
	return
}

// SortRecommandChannel sort recommend channels.
func (s *Service) SortRecommandChannel(c context.Context, tops []int64, recommends []int64) (err error) {
	var (
		topNum    int
		channels  []*model.Channel
		topMap    = make(map[int64]int32, len(tops))
		normalMap = make(map[int64]int32, len(recommends))
	)
	if channels, err = s.RecommandChannels(c); err != nil {
		return
	}
	if len(channels) == 0 {
		return ecode.TagOperateFail
	}
	topNum = len(tops)
	for index, tid := range tops {
		topMap[tid] = int32(index)
	}
	for index, tid := range recommends {
		normalMap[tid] = int32(topNum + index)
	}
	for _, channel := range channels {
		k, ok := topMap[channel.ID]
		if ok {
			channel.TopRank = k
			channel.AttrSet(model.ChannelAttrTop, model.ChannelTopYes)
			continue
		}
		v, b := normalMap[channel.ID]
		if b {
			channel.TopRank = v
			channel.AttrSet(model.ChannelAttrTop, model.ChannelTopNo)
			continue
		}
		channel.State = model.ChanStateCommon
	}
	_, err = s.dao.UpdateChannels(c, channels)
	return
}

// MigrateRecommendChannel migrate recommend channel.
func (s *Service) MigrateRecommendChannel(c context.Context, tid int64, tname, uname string) (err error) {
	var (
		tag     *model.Tag
		channel *model.Channel
	)
	if tid > 0 {
		tag, err = s.dao.Tag(c, tid)
	} else {
		tag, err = s.dao.TagByName(c, tname)
	}
	if err != nil {
		return
	}
	if tag == nil || tag.State == model.StateDel {
		return ecode.TagNotExist
	}
	if tag.State == model.StateShield {
		return ecode.TagAlreadyShield
	}
	if channel, err = s.dao.Channel(c, tag.ID); err != nil {
		return
	}
	if channel == nil || channel.State < model.ChanStateCommon {
		return ecode.ChannelNotExist
	}
	if channel.State == model.ChanStateRecomend {
		return ecode.ChannelRecommended
	}
	if channel.TopRank, err = s.dao.CountRecommendChannel(c); err != nil {
		return
	}
	channel.State = model.ChanStateRecomend
	channel.Attr = channel.Attr & 1
	channel.Editor = uname
	_, err = s.dao.UpdateChannel(c, channel)
	return
}

// EditChannel edit channel.
func (s *Service) EditChannel(c context.Context, channelInfo *model.ChannelInfo, uname string) (err error) {
	var (
		tag      *model.Tag
		category *model.ChannelCategory
		channel  *model.Channel
	)
	if channelInfo.ID > 0 {
		tag, err = s.dao.Tag(c, channelInfo.ID)
		if tag == nil {
			return ecode.TagNotExist
		}
	} else {
		err = s.dao.Filter(c, channelInfo.Name)
		if err != nil {
			return
		}
		tag, err = s.dao.TagByName(c, channelInfo.Name)
	}
	if err != nil {
		return ecode.TagOperateFail
	}
	if channelInfo.State == model.ChanStateStop {
		channelInfo.Rules = make([]*model.ChannelRule, 0)
		channelInfo.Synonyms = make([]*model.ChannelSynonym, 0)
	}
	if category, err = s.dao.ChannelCategory(c, channelInfo.Type); err != nil {
		return
	}
	if category == nil || category.State == model.StateDel {
		return ecode.ChanTypeNotExist
	}
	if tag == nil {
		tag = &model.Tag{
			Type:         model.TypeBiliContent,
			Name:         channelInfo.Name,
			Cover:        channelInfo.Cover,
			Content:      channelInfo.Content,
			Verify:       model.VerifyNone,
			State:        model.StateNormal,
			HeadCover:    channelInfo.HeadCover,
			ShortContent: channelInfo.ShortContent,
		}
		if err = s.createTag(c, tag); err != nil {
			return err
		}
	}
	if tag.State == model.StateShield {
		return ecode.TagAlreadyShield
	}
	if tag.State == model.StateDel {
		return ecode.TagAlreadyDelete
	}
	if channel, err = s.dao.Channel(c, tag.ID); err != nil {
		return err
	}
	if channel == nil {
		channel = &model.Channel{
			ID:       tag.ID,
			Type:     channelInfo.Type,
			Operator: uname,
			Cover:    channelInfo.Cover,
			State:    channelInfo.State,
			Editor:   uname,
		}
		if category.AttrVal(model.ChannelCategoryAttrINT) == model.CategoryStateShieldINT {
			channel.AttrSet(model.ChannelAttrINT, model.ChannelStateShieldINT)
		}
		return s.addChannel(c, channelInfo, channel, tag, uname)
	}
	if category.AttrVal(model.ChannelCategoryAttrINT) == model.CategoryStateShieldINT {
		channel.AttrSet(model.ChannelAttrINT, model.ChannelStateShieldINT)
	}
	return s.editChannel(c, channelInfo, channel, tag, uname)
}

// TODO 代码整合.
func (s *Service) createTag(c context.Context, tag *model.Tag) error {
	var affect int64
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("add tag tran error(%v)", err)
		return err
	}
	if tag.ID, err = s.dao.TxInsertTag(tx, tag); err != nil || tag.ID <= 0 {
		tx.Rollback()
		return err
	}
	if affect, err = s.dao.TxInsertTagCount(tx, tag.ID); err != nil || affect <= 0 {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *Service) addChannel(c context.Context, channel *model.ChannelInfo, chann *model.Channel, tag *model.Tag, uname string) (err error) {
	var (
		updateRules   = make([]*model.ChannelRule, 0)
		originRuleMap = make(map[int64]*model.ChannelRule)
		count         int32
	)
	if originRuleMap, err = s.dao.ChannelRuleMap(c, tag.ID); err != nil {
		return
	}
	if count, err = s.dao.CountChannelByType(c, channel.Type); err != nil {
		return
	}
	for _, rule := range channel.Rules {
		k, ok := originRuleMap[rule.ID]
		if ok && rule.State == k.State {
			r := &model.ChannelRule{
				ID:        k.ID,
				Tid:       k.Tid,
				InRule:    k.InRule,
				NotInRule: k.NotInRule,
				State:     k.State,
				Editor:    k.Editor,
				CTime:     k.CTime,
			}
			updateRules = append(updateRules, r)
			delete(originRuleMap, rule.ID)
			continue
		}
		r := &model.ChannelRule{
			Tid:       tag.ID,
			InRule:    rule.InRule,
			NotInRule: rule.NotInRule,
			State:     rule.State,
			Editor:    uname,
			CTime:     xtime.Time(time.Now().Unix()),
		}
		if ok {
			r.CTime = k.CTime
			delete(originRuleMap, rule.ID)
		}
		updateRules = append(updateRules, r)
	}
	for _, rule := range originRuleMap {
		r := &model.ChannelRule{
			ID:        rule.ID,
			Tid:       rule.Tid,
			InRule:    rule.InRule,
			NotInRule: rule.NotInRule,
			State:     model.StateDel,
			Editor:    uname,
			CTime:     rule.CTime,
		}
		updateRules = append(updateRules, r)
	}
	for index, synonym := range channel.Synonyms {
		if synonym.State != model.StateNormal {
			continue
		}
		synonym.Rank = int32(index)
		synonym.PTid = tag.ID
		synonym.Operator = uname
		synonym.CTime = xtime.Time(time.Now().Unix())
		synonym.MTime = xtime.Time(time.Now().Unix())
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.addChannel() tran error(%v)", err)
		return err
	}
	chann.Rank = count
	chann.AttrSet(model.ChannelAttrCheckBack, channel.CheckBack)
	chann.AttrSet(model.ChannelAttrCheckBack, channel.Activity)
	chann.AttrSet(model.ChannelAttrCheckBack, channel.INTShield)
	tag.Cover = channel.Cover
	tag.Type = model.TypeBiliContent
	tag.Content = channel.Content
	tag.HeadCover = channel.HeadCover
	tag.ShortContent = channel.ShortContent
	if _, err = s.dao.TxInsertChannel(tx, chann); err != nil {
		tx.Rollback()
		return
	}
	if len(updateRules) != 0 {
		if _, err = s.dao.TxUpdateChannelRules(tx, updateRules); err != nil {
			tx.Rollback()
			return
		}
	}
	if _, err = s.dao.TxUpdateTag(tx, tag); err != nil {
		tx.Rollback()
		return
	}
	if _, err = s.dao.TxUpStateChannelSynonym(tx, tag.ID, model.StateDel, uname); err != nil {
		tx.Rollback()
		return
	}
	if len(channel.Synonyms) > 0 {
		if _, err = s.dao.TxUpdateChannelSynonyms(tx, channel.Synonyms); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelTagCache(ctx, tag.ID, tag.Name)
		s.dao.DelChannelGroupCache(ctx, tag.ID)
	})
	return
}

func (s *Service) editChannel(c context.Context, channel *model.ChannelInfo, chann *model.Channel, tag *model.Tag, uname string) (err error) {
	var (
		change         bool
		updateRules    = make([]*model.ChannelRule, 0)
		originRuleMap  = make(map[int64]*model.ChannelRule)
		originGroupMap = make(map[int64]*model.ChannelSynonym)
		updateGroups   = make([]*model.ChannelSynonym, 0)
	)
	if originRuleMap, err = s.dao.ChannelRuleMap(c, tag.ID); err != nil {
		return
	}
	if originGroupMap, _, err = s.dao.ChannelSynonymMap(c, tag.ID); err != nil {
		return
	}
	for _, rule := range channel.Rules {
		k, ok := originRuleMap[rule.ID]
		if ok && rule.State == k.State {
			r := &model.ChannelRule{
				ID:        k.ID,
				Tid:       k.Tid,
				InRule:    k.InRule,
				NotInRule: k.NotInRule,
				State:     k.State,
				Editor:    k.Editor,
				CTime:     k.CTime,
			}
			updateRules = append(updateRules, r)
			delete(originRuleMap, rule.ID)
			continue
		}
		change = true
		r := &model.ChannelRule{
			Tid:       tag.ID,
			InRule:    rule.InRule,
			NotInRule: rule.NotInRule,
			State:     rule.State,
			Editor:    uname,
			CTime:     xtime.Time(time.Now().Unix()),
		}
		if ok {
			r.CTime = k.CTime
			delete(originRuleMap, rule.ID)
		}
		updateRules = append(updateRules, r)
	}
	for _, rule := range originRuleMap {
		r := &model.ChannelRule{
			ID:        rule.ID,
			Tid:       rule.Tid,
			InRule:    rule.InRule,
			NotInRule: rule.NotInRule,
			State:     model.StateDel,
			Editor:    uname,
			CTime:     rule.CTime,
		}
		updateRules = append(updateRules, r)
	}
	for index, group := range channel.Synonyms {
		k, ok := originGroupMap[group.Tid]
		if ok && group.State == k.State && int32(index) == k.Rank && group.Alias == k.Alias {
			cg := &model.ChannelSynonym{
				ID:       k.ID,
				PTid:     k.PTid,
				Tid:      k.Tid,
				Alias:    k.Alias,
				Rank:     k.Rank,
				Operator: k.Operator,
				State:    k.State,
				CTime:    k.CTime,
				MTime:    k.MTime,
			}
			updateGroups = append(updateGroups, cg)
			delete(originGroupMap, group.Tid)
			continue
		}
		change = true
		cg := &model.ChannelSynonym{
			PTid:     tag.ID,
			Tid:      group.Tid,
			Alias:    group.Alias,
			Rank:     int32(index),
			Operator: uname,
			State:    group.State,
			CTime:    xtime.Time(time.Now().Unix()),
			MTime:    xtime.Time(time.Now().Unix()),
		}
		if ok {
			cg.CTime = k.CTime
			delete(originGroupMap, group.Tid)
		}
		updateGroups = append(updateGroups, cg)
	}
	for _, group := range originGroupMap {
		cg := &model.ChannelSynonym{
			PTid:     group.PTid,
			Tid:      group.Tid,
			Alias:    group.Alias,
			Rank:     group.Rank,
			Operator: uname,
			State:    model.StateDel,
			CTime:    group.CTime,
			MTime:    xtime.Time(time.Now().Unix()),
		}
		updateGroups = append(updateGroups, cg)
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.editChannel(%v) tran error(%v)", channel, err)
		return
	}
	if chann.Type != channel.Type {
		change = true
		chann.Type = channel.Type
	}
	if chann.AttrVal(model.ChannelAttrCheckBack) != channel.CheckBack {
		change = true
		chann.AttrSet(model.ChannelAttrCheckBack, channel.CheckBack)
	}
	if chann.AttrVal(model.ChannelAttrActivity) != channel.Activity {
		change = true
		chann.AttrSet(model.ChannelAttrActivity, channel.Activity)
	}
	if chann.AttrVal(model.ChannelAttrINT) != channel.INTShield {
		change = true
		chann.AttrSet(model.ChannelAttrINT, channel.INTShield)
	}
	if chann.State != channel.State {
		change = true
		chann.State = channel.State
	}
	if tag.Cover != channel.Cover {
		change = true
		tag.Cover = channel.Cover
	}
	if tag.HeadCover != channel.HeadCover {
		change = true
		tag.HeadCover = channel.HeadCover
	}
	if tag.ShortContent != channel.ShortContent {
		change = true
		tag.ShortContent = channel.ShortContent
	}
	if tag.Type != model.TypeBiliContent {
		change = true
		tag.Type = model.TypeBiliContent
	}
	if tag.Content != channel.Content {
		change = true
		tag.Content = channel.Content
	}
	if change {
		chann.Editor = uname
	}
	if _, err = s.dao.TxUpChannel(tx, chann); err != nil {
		tx.Rollback()
		return
	}
	if len(updateRules) != 0 {
		if _, err = s.dao.TxUpdateChannelRules(tx, updateRules); err != nil {
			tx.Rollback()
			return
		}
	}
	if _, err = s.dao.TxUpdateTag(tx, tag); err != nil {
		tx.Rollback()
		return
	}
	if len(updateGroups) > 0 {
		if _, err = s.dao.TxUpdateChannelSynonyms(tx, updateGroups); err != nil {
			tx.Rollback()
			return
		}
	}
	if err = tx.Commit(); err != nil {
		return
	}
	// TODO  清缓存、通知search更新
	s.cacheCh.Do(c, func(ctx context.Context) {
		s.dao.DelTagCache(ctx, tag.ID, tag.Name)
		s.dao.DelChannelGroupCache(ctx, tag.ID)
	})
	return
}

// ChannelShieldINT channel shild int.
func (s *Service) ChannelShieldINT(c context.Context, tid int64, state int32, uname string) (err error) {
	var (
		channel *model.Channel
	)
	if channel, err = s.dao.Channel(c, tid); err != nil {
		return
	}
	if channel == nil {
		return ecode.ChannelNotExist
	}
	if channel.AttrVal(model.ChannelAttrINT) == state {
		return ecode.ChannelNoChange
	}
	channel.AttrSet(model.ChannelAttrINT, state)
	channel.Operator = uname
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.ChannelShieldINT(%d,%d,%s) error(%v)", tid, state, uname, err)
		return
	}
	if _, err = s.dao.TxUpChannelAttr(tx, channel); err != nil {
		tx.Rollback()
		return
	}
	return tx.Commit()
}
