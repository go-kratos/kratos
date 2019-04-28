package service

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	"go-common/app/admin/main/aegis/model"
	"go-common/app/admin/main/aegis/model/net"
	"go-common/library/ecode"
	"go-common/library/log"
	"go-common/library/xstr"
)

func (s *Service) checkBindToken(c context.Context, tokenID []int64) (tokenMap map[string]*net.Token, err error, msg string) {
	var (
		tokens []*net.Token
	)

	tokenMap = map[string]*net.Token{}
	if len(tokenID) == 0 {
		return
	}

	if tokens, err = s.gorm.Tokens(c, tokenID); err != nil {
		log.Error("checkBindToken s.gorm.Tokens error(%v)", err)
		return
	}
	for _, tk := range tokens {
		if !tk.IsAssign() {
			err = ecode.AegisTokenNotAssign
			msg = fmt.Sprintf(ecode.AegisTokenNotAssign.Message(), tk.ChName)
			return
		}
		id := strconv.FormatInt(tk.ID, 10)
		tokenMap[id] = tk
	}
	if len(tokenID) == len(tokens) {
		return
	}
	for _, id := range tokenID {
		if _, exist := tokenMap[strconv.FormatInt(id, 10)]; !exist {
			err = ecode.AegisTokenNotFound
			msg = fmt.Sprintf("id=%d的令牌 %s", id, ecode.AegisTokenNotAssign.Message())
			return
		}
	}
	return
}

func (s *Service) fetchOldBindAndLog(c context.Context, elementID int64, tp []int8) (all map[int64]*net.TokenBind, availableLog []string, err error) {
	var (
		oldBindMap   map[int64][]*net.TokenBind
		oldBindToken map[int64][]*net.Token
	)

	all = map[int64]*net.TokenBind{}
	availableLog = []string{}
	if oldBindMap, err = s.gorm.TokenBindByElement(c, []int64{elementID}, tp, false); err != nil {
		log.Error("fetchOldBindAndLog s.gorm.TokenBindByElement error(%v)", err)
		return
	}
	oldBindAvailable := []*net.TokenBind{}
	allBindMap := map[int64]*net.TokenBind{}
	for _, item := range oldBindMap[elementID] {
		allBindMap[item.ID] = item
		oldBindAvailable = append(oldBindAvailable, item)
	}
	if oldBindToken, err = s.bindTokens(c, oldBindAvailable); err != nil {
		log.Error("fetchOldBindAndLog s.bindTokens error(%v)", err)
		return
	}
	for bindID, items := range oldBindToken {
		sub := []string{}
		for _, item := range items {
			sub = append(sub, item.FormatLog())
		}
		chname := ""
		tp := int8(0)
		if allBindMap[bindID] != nil {
			chname = allBindMap[bindID].ChName
			tp = allBindMap[bindID].Type
		}
		availableLog = append(availableLog, fmt.Sprintf(net.BindLogTemp, chname, tp, strings.Join(sub, ",")))
	}
	all = allBindMap
	return
}

func (s *Service) compareFlowBind(c context.Context, tx *gorm.DB, flowID int64, tokenID []int64, isUpdate bool) (diff string, changed []int64, err error, msg string) {
	var (
		tokenMap     map[string]*net.Token
		oldBindAll   = map[int64]*net.TokenBind{}
		newFormatLog = []string{}
		oldFormatLog = []string{}
		disable      = time.Now()
		recovered    = net.Recovered
		existTokenID = map[string]int{}
	)

	log.Info("compareFlowBind start flow(%d) tokenid(%+v) isUpdate(%v)", flowID, tokenID, isUpdate)
	//新绑定查询，检查是否为赋值语句token & token是否存在
	if len(tokenID) > 0 {
		if tokenMap, err, msg = s.checkBindToken(c, tokenID); err != nil {
			log.Error("compareFlowBind s.checkBindToken error(%v)", err)
			return
		}
		for _, item := range tokenMap {
			newFormatLog = append(newFormatLog, fmt.Sprintf(net.BindLogTemp, item.ChName, item.Type, item.FormatLog()))
		}
	}

	//获取现有所有绑定和可用绑定
	if isUpdate {
		if oldBindAll, oldFormatLog, err = s.fetchOldBindAndLog(c, flowID, []int8{net.BindTypeFlow}); err != nil {
			log.Error("compareFlowBind s.fetchOldBindAndLog error(%v)", err)
			return
		}
	}
	if len(tokenMap) == 0 && len(oldBindAll) == 0 {
		return
	}

	//从旧的中过滤新的
	for _, item := range oldBindAll {
		if _, exist := tokenMap[item.TokenID]; exist {
			existTokenID[item.TokenID] = 1
			if item.IsAvailable() {
				continue
			}

			item.DisableTime = recovered
		} else {
			item.DisableTime = disable
		}
		if err = s.gorm.UpdateFields(c, tx, net.TableTokenBind, item.ID, map[string]interface{}{"disable_time": item.DisableTime}); err != nil {
			log.Error("compareFlowBind s.gorm.UpdateFields error(%v)", err)
			return
		}
		changed = append(changed, item.ID)
	}

	//从新的中过滤旧的
	for tokenID, tk := range tokenMap {
		if _, exist := existTokenID[tokenID]; exist {
			continue
		}
		nw := &net.TokenBind{TokenID: tokenID, ChName: tk.ChName, ElementID: flowID, Type: net.BindTypeFlow}
		if err = s.gorm.AddItem(c, tx, nw); err != nil {
			log.Error("compareFlowBind s.gorm.AddItem error(%v)", err)
			return
		}
		changed = append(changed, nw.ID)
	}

	if len(changed) > 0 && (len(newFormatLog) > 0 || len(oldFormatLog) > 0) {
		diff = model.LogFieldTemp(model.LogFieldTokenID, strings.Join(newFormatLog, ";"), strings.Join(oldFormatLog, ";"), isUpdate)
	}
	log.Info("compareFlowBind end flow(%d) tokenid(%+v) isupdate(%v) diff(%s)", flowID, tokenID, isUpdate, diff)
	return
}

func (s *Service) compareTranBind(c context.Context, tx *gorm.DB, tranID int64, binds []*net.TokenBindParam, isUpdate bool) (diff string, changed []int64, err error, msg string) {
	var (
		relatedTokenID []int64
		tokenID        = []int64{}
		existBindParam = map[int64]*net.TokenBindParam{}
		tokenMap       map[string]*net.Token
		oldBindAll     = map[int64]*net.TokenBind{}
		newFormatLog   = []string{}
		oldFormatLog   []string
		disable        = time.Now()
		recovered      = net.Recovered
	)

	log.Info("compareTranBind start transition(%d) binds(%+v) isUpdate(%v)", tranID, binds, isUpdate)

	//现有全部绑定关系和可用绑定
	if isUpdate {
		if oldBindAll, oldFormatLog, err = s.fetchOldBindAndLog(c, tranID, net.BindTranType); err != nil {
			log.Error("compareTranBind s.fetchOldBindAndLog error(%v)", err)
			return
		}
	}

	//新绑定处理
	for _, item := range binds {
		//是否已存在
		if item.ID > 0 && oldBindAll[item.ID] == nil {
			log.Error("compareTranBind binds(%+v) not found", item)
			err = ecode.RequestErr
			return
		}
		if item.ID > 0 {
			existBindParam[item.ID] = item
		}

		//tokenid排序重组
		if relatedTokenID, err = xstr.SplitInts(item.TokenID); err != nil {
			log.Error("compareTranBind xstr.SplitInts(%+v) error(%v)", item, err)
			return
		}
		if len(relatedTokenID) == 0 {
			log.Error("compareTranBind bind(%+v) tokenid empty ", item)
			err = ecode.RequestErr
			return
		}
		sort.Sort(net.Int64Slice(relatedTokenID))
		item.TokenID = xstr.JoinInts(relatedTokenID)
		tokenID = append(tokenID, relatedTokenID...)
	}
	//检查是否为赋值语句token & token是否存在
	if len(tokenID) > 0 {
		if tokenMap, err, msg = s.checkBindToken(c, tokenID); err != nil {
			log.Error("compareTranBind s.checkBindToken error(%v)", err)
			return
		}
	}

	for _, item := range binds {
		//前端没传中文名则自行组合
		nwLog := []string{}
		chname := item.ChName
		tkname := map[string]int64{}
		for _, tid := range strings.Split(item.TokenID, ",") {
			tk := tokenMap[tid]
			if tk == nil {
				continue
			}

			//token_name级别的过滤
			if tkname[tk.Name] > 0 {
				log.Error("compareTranBind bind(%+v) duplicated with token name(%s)", item, tk.Name)
				err = ecode.RequestErr
				return
			}
			tkname[tk.Name] = tk.ID

			if item.ChName == "" {
				chname = chname + tokenMap[tid].ChName
			}
			nwLog = append(nwLog, tokenMap[tid].FormatLog())
		}
		chnameMerge := []rune(chname)
		last := len(chnameMerge)
		if last > 16 {
			last = 16
		}
		item.ChName = string(chnameMerge[:last])
		newFormatLog = append(newFormatLog, fmt.Sprintf(net.BindLogTemp, item.ChName, item.Type, strings.Join(nwLog, ",")))
	}

	//从旧的中过滤新的
	for id, item := range oldBindAll {
		updateField := map[string]interface{}{}
		if newParam, exist := existBindParam[id]; exist {
			if !item.IsAvailable() {
				item.DisableTime = recovered
				updateField["disable_time"] = recovered
			}
			if item.TokenID != newParam.TokenID {
				item.TokenID = newParam.TokenID
				updateField["token_id"] = newParam.TokenID
			}
			if newParam.ChName != "" && newParam.ChName != item.ChName {
				item.ChName = newParam.ChName
				updateField["ch_name"] = newParam.ChName
			}
			if newParam.Type != item.Type {
				item.Type = newParam.Type
				updateField["type"] = newParam.Type
			}
		} else if item.IsAvailable() {
			item.DisableTime = disable
			updateField["disable_time"] = disable
		} else {
			continue
		}
		if len(updateField) == 0 {
			continue
		}

		if err = s.gorm.UpdateFields(c, tx, net.TableTokenBind, item.ID, updateField); err != nil {
			log.Error("compareTranBind s.gorm.UpdateFields error(%v)", err)
			return
		}
		changed = append(changed, item.ID)
	}
	//从新的中过滤旧的
	for _, newParam := range binds {
		if newParam.ID > 0 && oldBindAll[newParam.ID] != nil {
			continue
		}

		nw := &net.TokenBind{TokenID: newParam.TokenID, ChName: newParam.ChName, ElementID: tranID, Type: newParam.Type}
		if err = s.gorm.AddItem(c, tx, nw); err != nil {
			log.Error("compareTranBind s.gorm.AddItem error(%v)", err)
			return
		}
		changed = append(changed, nw.ID)
	}

	if len(changed) > 0 && (len(newFormatLog) > 0 || len(oldFormatLog) > 0) {
		diff = model.LogFieldTemp(model.LogFieldTokenID, strings.Join(newFormatLog, ";"), strings.Join(oldFormatLog, ";"), isUpdate)
	}
	log.Info("compareTranBind end transition(%d) binds(%+v) isupdate(%v) diff(%s)", tranID, binds, isUpdate, diff)
	return
}

func (s *Service) bindTokens(c context.Context, binds []*net.TokenBind) (result map[int64][]*net.Token, err error) {
	var (
		tokenIDSlice []int64
		tokens       []*net.Token
	)

	result = map[int64][]*net.Token{}
	if len(binds) == 0 {
		return
	}

	tids := []int64{}
	bindTokenMap := map[int64][]int64{}
	for _, item := range binds {
		if tokenIDSlice, err = xstr.SplitInts(item.TokenID); err != nil {
			log.Error("tokenBindDetail xstr.SplitInts(%s) error(%v)", item.TokenID, err)
			return
		}
		bindTokenMap[item.ID] = tokenIDSlice
		if len(tokenIDSlice) == 0 {
			continue
		}

		tids = append(tids, tokenIDSlice...)
	}
	if len(tids) == 0 {
		return
	}

	if tokens, err = s.gorm.Tokens(c, tids); err != nil {
		log.Error("tokenBindDetail s.gorm.Tokens(%v) error(%v)", tids, err)
		return
	}
	if len(tokens) == 0 {
		return
	}
	tokenMap := map[int64]*net.Token{}
	for _, item := range tokens {
		tokenMap[item.ID] = item
	}
	for bindID, tidList := range bindTokenMap {
		result[bindID] = []*net.Token{}
		for _, id := range tidList {
			if tokenMap[id] == nil {
				continue
			}

			result[bindID] = append(result[bindID], tokenMap[id])
		}
	}
	return
}

//tokenBindDetail 一条绑定的详细信息，获取了token详情
func (s *Service) tokenBindDetail(c context.Context, binds []*net.TokenBind) (result []*net.TokenBindDetail, err error) {
	var (
		tokenIDSlice []int64
		tokens       []*net.Token
	)

	if len(binds) == 0 {
		return
	}
	details := make([]*net.TokenBindDetail, len(binds))

	//get tokens
	tids := []int64{}
	bindTokenMap := map[string][]int64{}
	for k, item := range binds {
		details[k] = &net.TokenBindDetail{
			ID:          item.ID,
			Type:        item.Type,
			ElementID:   item.ElementID,
			TokenID:     item.TokenID,
			ChName:      item.ChName,
			DisableTime: item.DisableTime,
			Tokens:      []*net.Token{},
		}
		if tokenIDSlice, err = xstr.SplitInts(item.TokenID); err != nil {
			log.Error("tokenBindDetail xstr.SplitInts(%s) error(%v)", item.TokenID, err)
			return
		} else if len(tokenIDSlice) == 0 {
			continue
		} else {
			tids = append(tids, tokenIDSlice...)
		}
		bindTokenMap[item.TokenID] = tokenIDSlice
	}
	if len(tids) == 0 {
		result = details
		return
	}
	if tokens, err = s.gorm.Tokens(c, tids); err == nil && len(tokens) == 0 {
		log.Error("tokenBindDetail s.gorm.Tokens(%v) not found", tids)
		result = details
		err = nil
		return
	}
	if err != nil {
		log.Error("TokenBindDetail find token error(%v) tids(%v)", err, tids)
		return
	}
	tokenMap := make(map[int64]*net.Token, len(tokens))
	for _, item := range tokens {
		tokenMap[item.ID] = item
	}

	//dispatch tokens to detail
	for _, item := range details {
		for _, id := range bindTokenMap[item.TokenID] {
			if tokenMap[id] != nil {
				item.Tokens = append(item.Tokens, tokenMap[id])
			}
		}
	}
	result = details
	return
}
