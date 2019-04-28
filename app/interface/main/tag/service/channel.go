package service

import (
	"context"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"go-common/app/interface/main/tag/model"
	taGrpcModel "go-common/app/service/main/tag/api"
	rpcModel "go-common/app/service/main/tag/model"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

// get channel caches.
func (s *Service) channelCaches() (err error) {
	var (
		lastID                int64
		batchNum              = int(model.BatchSize)
		channelCategories     []*model.ChannelCategory
		categoryMap           = make(map[int64]*model.ChannelCategory)
		channelByCategoryyMap = make(map[int64][]*model.Channel)
		channelMap            = make(map[int64]*model.Channel)
		recommendChannels     = make([]*model.Channel, 0)
		tidMap                = make(map[int64]struct{})
		ruleClassMap          = make(map[int64]*model.ChannelRuleClassifier)
		ruleMap               = make(map[int64][]*model.ChannelRule)
	)
	// get all channel categories.
	for model.BatchSize == int32(batchNum) {
		var categories []*taGrpcModel.ChannelCategory
		if categories, err = s.dao.ChannelCategories(context.TODO(), lastID, model.BatchSize, model.ChannelCategoryStateOK); err != nil {
			return
		}
		batchNum = len(categories)
		for _, category := range categories {
			if category.Id > lastID {
				lastID = category.Id
			}
			if category.State != model.ChannelCategoryStateOK {
				continue
			}
			t := &model.ChannelCategory{
				ID:        category.Id,
				Name:      category.Name,
				State:     category.State,
				Order:     category.Order,
				CTime:     category.Ctime,
				MTime:     category.Mtime,
				INTShield: category.IntShield,
			}
			categoryMap[category.Id] = t
			channelCategories = append(channelCategories, t)
		}
	}
	// get all channels.
	lastID = 0
	batchNum = int(model.BatchSize)
	for model.BatchSize == int32(batchNum) {
		var channels []*taGrpcModel.Channel
		if channels, err = s.dao.Channels(context.TODO(), lastID, model.BatchSize); err != nil {
			return
		}
		batchNum = len(channels)
		for _, channel := range channels {
			if channel.Tid > lastID {
				lastID = channel.Tid
			}
			if channel.State == model.ChanStateStop || channel.State == model.ChanStateOffline {
				continue
			}
			if _, ok := categoryMap[channel.Type]; !ok {
				continue
			}
			t := &model.Channel{
				ID:      channel.Tid,
				Type:    channel.Type,
				Rank:    channel.Rank,
				Attr:    channel.Attr,
				State:   channel.State,
				CTime:   channel.Ctime,
				MTime:   channel.Mtime,
				TopRank: channel.TopRank,
			}
			tidMap[channel.Tid] = struct{}{}
			channelMap[t.ID] = t
		}
	}
	// get all channel rules.
	lastID = 0
	batchNum = int(model.BatchSize)
	for model.BatchSize == int32(batchNum) {
		var rules []*taGrpcModel.ChannelRule
		if rules, err = s.dao.ChannelRules(context.TODO(), lastID, model.BatchSize, model.ChannelRuleStateOK); err != nil {
			return
		}
		batchNum = len(rules)
		for _, rule := range rules {
			if rule.Id > lastID {
				lastID = rule.Id
			}
			if rule.State != model.ChannelRuleStateOK {
				continue
			}
			if rule.Tid <= 0 {
				continue
			}
			if (rule.InRule == "" && rule.NotinRule == "") || rule.InRule == rule.NotinRule {
				continue
			}
			k, ok := ruleMap[rule.Tid]
			if !ok {
				k = make([]*model.ChannelRule, 0)
			}
			r := &model.ChannelRule{
				Tid: rule.Tid,
			}
			if rule.NotinRule != "" && rule.NotinRule != "0" {
				var (
					tidA int64
					tidB int64
				)
				if tidA, err = strconv.ParseInt(rule.InRule, 10, 64); err != nil || tidA <= 0 {
					log.Error("s.channelCaches() ParseInRule(%s) error: %v", rule, err)
					err = nil
					continue
				}
				if tidB, err = strconv.ParseInt(rule.NotinRule, 10, 64); err != nil || tidB <= 0 {
					log.Error("s.channelCaches() ParseNotInRule(%s) error: %v", rule, err)
					err = nil
					continue
				}
				r.TidA = tidA
				r.TidB = tidB
				r.Flag = model.ChannelRuleFlagMinus
				tidMap[tidA] = struct{}{}
				tidMap[tidB] = struct{}{}
				k = append(k, r)
				ruleMap[rule.Tid] = k
				continue
			}
			switch strings.Count(rule.InRule, ",") {
			case 0:
				var tid int64
				if tid, err = strconv.ParseInt(rule.InRule, 10, 64); err != nil || tid <= 0 {
					log.Error("s.channelCaches() ParseNotInRule(%s) error: %v", rule, err)
					err = nil
					continue
				}
				r.TidA = tid
				r.Flag = model.ChannelRuleFlagSingle
				tidMap[tid] = struct{}{}
			case 1:
				var tids []int64
				if tids, err = xstr.SplitInts(rule.InRule); err != nil {
					log.Error("s.channelCaches() ParseRule(%s) error: %v", rule, err)
					err = nil
					continue
				}
				r.TidA = tids[0]
				r.TidB = tids[1]
				r.Flag = model.ChannelRuleFlagPlus
				tidMap[tids[0]] = struct{}{}
				tidMap[tids[1]] = struct{}{}
			default:
				log.Error("s.channelCaches() ParseRule(%s) error: %v", rule, err)
				continue
			}
			k = append(k, r)
			ruleMap[rule.Tid] = k
		}
	}
	// get rules and channels tag infos.
	var (
		tids   = make([]int64, 0, len(tidMap))
		tagMap map[int64]*taGrpcModel.Tag
	)
	for tid := range tidMap {
		if tid <= 0 {
			continue
		}
		tids = append(tids, tid)
	}
	if tagMap, err = s.dao.TagMap(context.TODO(), tids, model.NoneUserID); err != nil {
		return
	}
	// set channel name,cover,detail_content, flag invalid channel and slite channels by channel's type and channel state.
	var delTids []int64
	for _, t := range channelMap {
		k, ok := tagMap[t.ID]
		if !ok || k.State != model.TagStateNormal {
			delTids = append(delTids, t.ID)
			continue
		}
		t.Name = k.Name
		t.Cover = k.Cover
		t.Content = k.Content
		t.Bind = k.Bind
		t.Sub = k.Sub
		t.HeadCover = k.HeadCover
		t.ShortContent = k.ShortContent
		channelByCategoryyMap[t.Type] = append(channelByCategoryyMap[t.Type], t)
		if t.Recommend() {
			recommendChannels = append(recommendChannels, t)
		}
	}
	// delete status != normal channels.
	for _, tid := range delTids {
		delete(channelMap, tid)
	}
	// set rule name into rules, and splite rule by rule's tid.
	for tid, rules := range ruleMap {
		k, ok := ruleClassMap[tid]
		if !ok {
			k = &model.ChannelRuleClassifier{
				Single: make([]*model.ChannelRule, 0),
				Plus:   make([]*model.ChannelRule, 0),
				Minus:  make([]*model.ChannelRule, 0),
			}
		}
		for _, rule := range rules {
			switch rule.Flag {
			case model.ChannelRuleFlagSingle:
				v, b := tagMap[rule.TidA]
				if !b || v.State != model.TagStateNormal {
					continue
				}
				rule.TidAName = v.Name
				k.Single = append(k.Single, rule)
			case model.ChannelRuleFlagPlus, model.ChannelRuleFlagMinus:
				tA, ba := tagMap[rule.TidA]
				tB, bb := tagMap[rule.TidB]
				if !ba || !bb {
					continue
				}
				if tA.State != model.TagStateNormal || tB.State != model.TagStateNormal {
					continue
				}
				rule.TidAName = tA.Name
				rule.TidBName = tB.Name
				if rule.Flag == model.ChannelRuleFlagPlus {
					k.Plus = append(k.Plus, rule)
				} else {
					k.Minus = append(k.Minus, rule)
				}
			default:
			}
		}
		ruleClassMap[tid] = k
	}
	s.channelLock.Lock()
	s.channelMap = channelMap
	s.channelRecommand = recommendChannels
	s.channelCategories = channelCategories
	s.channelTypeMap = channelByCategoryyMap
	s.channelRule = ruleClassMap
	s.channelLock.Unlock()
	return
}

func (s *Service) channelproc() {
	for {
		time.Sleep(time.Duration(s.c.Tag.ChannelRefreshTime))
		s.channelCaches()
	}
}

// ChannelCategory channel category.
func (s *Service) ChannelCategory(c context.Context) (res []*model.ChannelCategory, err error) {
	categories := s.channelCategories
	res = make([]*model.ChannelCategory, 0, len(categories))
	for _, category := range categories {
		t := &model.ChannelCategory{}
		*t = *category
		res = append(res, t)
	}
	sort.Sort(model.ChannelCategorySort(res))
	return
}

// ChannelCategories channel category.
func (s *Service) ChannelCategories(c context.Context, arg *model.ArgChannelCategories) (res []*model.ChannelCategory, err error) {
	categories := s.channelCategories
	res = make([]*model.ChannelCategory, 0, len(categories))
	for _, category := range categories {
		channels, ok := s.channelTypeMap[category.ID]
		if !ok || len(channels) == 0 {
			continue
		}
		if arg.From == model.ChannelFromINT {
			if category.INTShield == model.ChannelCategoryStateINTField {
				continue
			}
			nodata := true
			// Check that there is at least one unshielded channel under the category.
			for _, channel := range channels {
				if channel.AttrVal(model.ChannelAttrINT) != model.ChannelStateShieldINT {
					nodata = false
					break
				}
			}
			if nodata {
				continue
			}
		}
		t := &model.ChannelCategory{}
		*t = *category
		res = append(res, t)
	}
	sort.Sort(model.ChannelCategorySort(res))
	return
}

// ChannelRule channel rule.
func (s *Service) ChannelRule(c context.Context, tid int64) (res *model.ChannelRuleClassifier, err error) {
	res, ok := s.channelRule[tid]
	if !ok {
		res = &model.ChannelRuleClassifier{
			Single: make([]*model.ChannelRule, 0),
			Plus:   make([]*model.ChannelRule, 0),
			Minus:  make([]*model.ChannelRule, 0),
		}
	}
	return
}

// ChanneList get channel list by channel category id.
func (s *Service) ChanneList(c context.Context, mid, tp int64, from int32) (res []*model.Channel, err error) {
	res = make([]*model.Channel, 0)
	if from == model.ChannelFromINT {
		for _, v := range s.channelCategories {
			if v.ID != tp {
				continue
			}
			if v.INTShield == model.ChannelCategoryStateINTField {
				return
			}
		}
	}
	var (
		tids     []int64
		channels []*model.Channel
	)
	for _, cc := range s.channelTypeMap[tp] {
		if from == model.ChannelFromINT && cc.AttrVal(model.ChannelAttrINT) == model.ChannelStateShieldINT {
			continue
		}
		tids = append(tids, cc.ID)
		channels = append(channels, cc)
	}
	if len(tids) == 0 {
		return
	}
	tagMap, err := s.dao.TagMap(c, tids, mid)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.ChanneList(%d,%d) TagMapByID() tids:%+v,mid:%d,err:%v", mid, tp, tids, mid, err)
		return
	}
	for _, v := range channels {
		k, ok := tagMap[v.ID]
		if !ok || k.State != model.TagStateNormal {
			continue
		}
		t := &model.Channel{}
		*t = *v
		t.Name = k.Name
		t.Cover = k.Cover
		t.HeadCover = k.HeadCover
		t.ShortContent = k.ShortContent
		t.Content = k.Content
		t.Attention = k.Attention
		t.Sub = k.Sub
		t.Bind = k.Bind
		res = append(res, t)
	}
	sort.Sort(model.ChannelSort(res))
	return
}

// RecommandChannel recommand channels list.
func (s *Service) RecommandChannel(c context.Context, mid int64, from int32) (res []*model.Channel, err error) {
	var (
		tids     []int64
		channels []*model.Channel
	)
	for _, cc := range s.channelRecommand {
		if from == model.ChannelFromINT && cc.AttrVal(model.ChannelAttrINT) == model.ChannelStateShieldINT {
			continue
		}
		tids = append(tids, cc.ID)
		channels = append(channels, cc)
	}
	res = make([]*model.Channel, 0, len(channels))
	if len(tids) == 0 {
		return
	}
	tagMap, err := s.dao.TagMap(c, tids, mid)
	if err != nil {
		if ecode.TagNotExist.Equal(err) {
			return
		}
		log.Error("s.RecommandChannel(%d) TagMapByID() tids:%+v,mid:%d,err:%v", mid, tids, mid, err)
		return
	}
	for _, v := range channels {
		k, ok := tagMap[v.ID]
		if !ok || k.State != model.TagStateNormal {
			continue
		}
		cc := &model.Channel{}
		*cc = *v
		cc.Attention = k.Attention
		cc.Sub = k.Sub
		cc.Bind = k.Bind
		res = append(res, cc)
	}
	sort.Sort(model.ChannelRecomendSort(res))
	return
}

// ChannelSquare channel square.
func (s *Service) ChannelSquare(c context.Context, arg *model.ReqChannelSquare) (res *model.ChannelSquare, err error) {
	var channels = make([]*model.Channel, 0, arg.TagNumber)
	if channels, err = s.discoveryChannel(c, arg.Mid, arg.TagNumber, arg.From); err != nil {
		return
	}
	res = &model.ChannelSquare{
		Channels: channels,
		Oids:     make(map[int64][]int64, arg.TagNumber),
	}
	if arg.OidNumber <= 0 {
		return
	}
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	for _, channel := range channels {
		tid := channel.ID
		wg.Add(1)
		go func() {
			defer wg.Done()
			req := &model.ArgChannelResource{
				Tid:        tid,
				Mid:        arg.Mid,
				Plat:       arg.Plat,
				LoginEvent: arg.LoginEvent,
				RequestCNT: arg.OidNumber,
				DisplayID:  arg.DisplayID,
				From:       model.ChannelFromApp,
				Build:      arg.Build,
				Buvid:      arg.Buvid,
				RealIP:     arg.RealIP,
			}
			cr, rerr := s.ChannelResources(c, req)
			if rerr != nil {
				return
			}
			lock.Lock()
			res.Oids[tid] = cr.Oids
			lock.Unlock()
		}()
	}
	wg.Wait()
	return
}

func (s *Service) discoveryChannel(c context.Context, mid int64, tagNumber, from int32) (res []*model.Channel, err error) {
	channels, err := s.RecommandChannel(c, mid, from)
	if err != nil {
		return nil, err
	}
	var (
		topNum     int32
		channelMap = make(map[int64]*model.Channel)
	)
	for _, v := range channels {
		if v.Attent() {
			continue
		}
		if v.Top() {
			if topNum < tagNumber {
				res = append(res, v)
				topNum++
				continue
			}
			return
		}
		channelMap[v.ID] = v
	}
	for _, channel := range channelMap {
		if topNum >= tagNumber {
			return
		}
		res = append(res, channel)
		topNum++
	}
	if topNum >= tagNumber {
		return
	}
	for _, catetory := range s.channelCategories {
		if catetory.INTShield == model.ChannelCategoryStateINTField && catetory.INTShield == from {
			continue
		}
		channeList, _ := s.ChanneList(c, mid, catetory.ID, from)
		for _, v := range channeList {
			if v.Attent() || v.Recommend() {
				continue
			}
			res = append(res, v)
			topNum++
			if topNum >= tagNumber {
				return
			}
		}
	}
	return
}

// DiscoveryChannel discovery channels list.
func (s *Service) DiscoveryChannel(c context.Context, mid int64, from int32) (res []*model.Channel, err error) {
	return s.discoveryChannel(c, mid, model.DiscoveryChannelNum, from)
}

// ResChannelInfos resource channel infos.
func (s *Service) ResChannelInfos(c context.Context, arg *model.ReqChannelResourceInfos) (res map[int64]*model.ChannelInfo, err error) {
	rcMap, err := s.ResChannelCheckBack(c, arg.Oids, rpcModel.ResTypeArchive)
	if err != nil {
		return
	}
	res = make(map[int64]*model.ChannelInfo, len(arg.IDs))
	for index, oid := range arg.Oids {
		tid := arg.Tids[index]
		id := arg.IDs[index]
		ci := &model.ChannelInfo{
			Tid:       tid,
			HitRules:  make([]string, 0),
			HitTNames: make([]string, 0),
		}
		if rc, ok := rcMap[oid]; ok {
			if channelInfo, b := rc.Channels[tid]; b {
				ci = channelInfo
			}
		}
		res[id] = ci
	}
	return
}

// ResChannelCheckBack ResChannelCheckBack.
func (s *Service) ResChannelCheckBack(c context.Context, oids []int64, tp int32) (res map[int64]*model.ResChannelCheckBack, err error) {
	res = make(map[int64]*model.ResChannelCheckBack, len(oids))
	resouceMap, err := s.dao.ResTags(c, oids, tp)
	if err != nil {
		return
	}
	var arcTidsMap = make(map[int64][]int64, len(oids))
	for oid, arcTags := range resouceMap {
		var tids = make([]int64, 0, len(arcTags))
		for _, v := range arcTags {
			if v.State != model.ResTagStateNormal && v.State != model.ResTagStateRegion {
				continue
			}
			tids = append(tids, v.Tid)
		}
		arcTidsMap[oid] = tids
	}
	channelMap, _ := s.resourceChannels(c, model.ManagerYes, arcTidsMap)
	for oid, hitChannelMap := range channelMap {
		var (
			checkBack int32
			hitMap    = make(map[int64]*model.ChannelInfo, len(hitChannelMap))
		)
		for tid, hitRules := range hitChannelMap {
			channel, ok := s.channelMap[tid]
			if !ok || channel == nil {
				continue
			}
			if channel.State != model.ChanStateCommon && channel.State != model.ChanStateRecommend {
				continue
			}
			if checkBack != model.StateResCheckBackYes && (channel.Attr&1 == model.StateResCheckBackYes) {
				checkBack = model.StateResCheckBackYes
			}
			hitTNameMap := make(map[int64]string, len(hitRules))
			rules := make([]string, 0, len(hitRules))
			for _, hitRule := range hitRules {
				switch hitRule.Flag {
				case model.ChannelRuleFlagSingle:
					hitTNameMap[hitRule.TidA] = hitRule.TidAName
					rules = append(rules, hitRule.TidAName)
				case model.ChannelRuleFlagPlus:
					hitTNameMap[hitRule.TidA] = hitRule.TidAName
					hitTNameMap[hitRule.TidB] = hitRule.TidBName
					rules = append(rules, strings.Join([]string{hitRule.TidAName, hitRule.TidBName}, " + "))
				case model.ChannelRuleFlagMinus:
					hitTNameMap[hitRule.TidA] = hitRule.TidAName
					rules = append(rules, strings.Join([]string{hitRule.TidAName, hitRule.TidBName}, " - "))
				default:
					continue
				}
			}
			tnames := make([]string, 0, len(hitTNameMap))
			for _, tname := range hitTNameMap {
				tnames = append(tnames, tname)
			}
			hitMap[tid] = &model.ChannelInfo{
				Tid:       tid,
				TName:     channel.Name,
				HitTNames: tnames,
				HitRules:  rules,
			}
		}
		res[oid] = &model.ResChannelCheckBack{
			Channels:  hitMap,
			CheckBack: checkBack,
		}
	}
	return
}

func (s *Service) resourceChannel(tids []int64, mng int32) (res map[int64][]*model.ChannelRule, channelIDs []int64) {
	res = make(map[int64][]*model.ChannelRule)
	for tid, rcs := range s.channelRule {
		rules := rcs.RuleCALC(tids, mng)
		if len(rules) == 0 {
			continue
		}
		res[tid] = rules
		channelIDs = append(channelIDs, tid)
	}
	return
}

func (s *Service) resourceChannels(c context.Context, mng int32, tidMap map[int64][]int64) (res map[int64]map[int64][]*model.ChannelRule, channelIDs []int64) {
	res = make(map[int64]map[int64][]*model.ChannelRule, len(tidMap))
	wg := sync.WaitGroup{}
	lock := sync.Mutex{}
	for oid, tids := range tidMap {
		tempOid := oid
		tempTids := tids
		wg.Add(1)
		go func() {
			defer wg.Done()
			rules, hitIDs := s.resourceChannel(tempTids, mng)
			lock.Lock()
			res[tempOid] = rules
			channelIDs = append(channelIDs, hitIDs...)
			lock.Unlock()
		}()
	}
	wg.Wait()
	return
}

// ChannelResources channel resources.
func (s *Service) ChannelResources(c context.Context, arg *model.ArgChannelResource) (res *model.ChannelResource, err error) {
	var (
		channel *model.Channel
		tag     *rpcModel.Tag
	)
	if arg.Name != "" {
		if tag, err = s.tagName(c, arg.Mid, arg.Name); err != nil {
			return
		}
		if tag == nil {
			err = ecode.TagNotExist
			return
		}
		arg.Tid = tag.ID
	}
	res = &model.ChannelResource{
		Oids:  make([]int64, 0),
		Pages: new(model.Page),
	}
	channel = s.channelMap[arg.Tid]
	if channel != nil && channel.State > model.ChanStateStop {
		res.IsChannel = true
		arg.Channel = model.TagChannelYes
	} else {
		res.IsChannel = false
		arg.Channel = model.TagChannelNo
	}
	res.Oids, err = s.aiRecommand(c, arg)
	if err != nil && arg.DisplayID == 1 {
		res.Oids, _, err = s.dao.NewArcsCache(c, arg.Tid, 0, 20)
		res.Failover = true
	}
	return
}

// ChannelDetail channel detail.
func (s *Service) ChannelDetail(c context.Context, arg *model.ReqChannelDetail) (res *model.ChannelDetail, err error) {
	var tag *taGrpcModel.Tag
	if arg.Tid > 0 {
		tag, err = s.dao.Tag(c, arg.Tid, arg.Mid)
	} else {
		if arg.TName, err = s.CheckName(arg.TName); err != nil {
			return
		}
		tag, err = s.dao.TagByName(c, arg.Mid, arg.TName)
	}
	if err != nil {
		return
	}
	if tag == nil || tag.State != model.TagStateNormal {
		return nil, ecode.TagNotExist
	}
	channel, ok := s.channelMap[tag.Id]
	if ok && arg.From == model.ChannelFromINT && channel.AttrVal(model.ChannelAttrINT) == model.ChannelStateShieldINT {
		err = ecode.TagAlreadyShield
		return
	}
	res = &model.ChannelDetail{
		Tag: &model.TagInfo{
			ID:           tag.Id,
			Name:         tag.Name,
			Type:         tag.Type,
			Cover:        tag.Cover,
			HeadCover:    tag.HeadCover,
			Content:      tag.Content,
			ShortContent: tag.ShortContent,
			Verify:       tag.Verify,
			Attr:         tag.Attr,
			Attention:    tag.Attention,
			State:        tag.State,
			Bind:         tag.Bind,
			Sub:          tag.Sub,
			CTime:        tag.Ctime,
			MTime:        tag.Mtime,
		},
		Synonym: make([]*model.ChannelSynonym, 0),
	}
	if !ok || channel == nil {
		return
	}
	switch channel.State {
	case model.ChanStateRecommend, model.ChanStateCommon:
		res.Tag.Activity = channel.AttrVal(model.ChannelAttrActivity)
	default:
		return
	}
	synonyms, _ := s.dao.ChannelGroup(c, tag.Id)
	sort.Sort(model.ChannelSynonymSort(synonyms))
	res.Synonym = synonyms
	return
}
