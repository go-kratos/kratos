package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"go-common/app/job/main/dm2/model"
	filterMdl "go-common/app/service/main/filter/api/grpc/v1"
	"go-common/library/database/sql"
	"go-common/library/ecode"
	"go-common/library/log"
)

const (
	_filterAreaSubtitle = "subtitle"
	_filterCapacity     = 5
	_contentSizeLimit   = 300
)

// SubtitleFilter .
// 1、只处理状态为审核待检测，发布待检测的数据
// 2、如果是审核待检测的数据，检测完毕，状态改为带审核，删除草稿缓存，删除字幕缓存
// 3、如果检测失败，状态改为审核驳回，并且更新驳回理由 删除缓存
// 4、如果是发布待检测的状态，检测完毕，状态改为发布，更新数据库逻辑，发布表更新，删除字幕缓存
// 5、如果发布检测失败。状态改为审核驳回，并且更新驳回理由。删除缓存
// 6、如果消费失败 数据丢失 容错
func (s *Service) SubtitleFilter(c context.Context, oid int64, subtitleID int64) (err error) {
	var (
		subtitle *model.Subtitle
	)
	if subtitle, err = s.dao.GetSubtitle(c, oid, subtitleID); err != nil {
		log.Error("params(oid::%v,subtitleID:%v)", oid, subtitleID)
		return
	}
	if subtitle == nil {
		log.Error("params(oid:%v,subtitleID:%v not found)", oid, subtitleID)
		return
	}
	switch subtitle.Status {
	case model.SubtitleStatusCheckToAudit:
		if err = s.checkToAudit(c, subtitle); err != nil {
			log.Error("checkToAudit.params(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
	case model.SubtitleStatusCheckPublish:
		if err = s.checkToPublish(c, subtitle); err != nil {
			log.Error("checkToPublish.params(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
	default:
		return
	}
	return
}

func (s *Service) checkToAudit(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		status = model.SubtitleStatusToAudit
		hits   []string
	)
	if hits, err = s.checkBfsData(c, subtitle); err != nil {
		log.Error("checkBfsData(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	if len(hits) > 0 {
		subtitle.RejectComment = "敏感词:" + strings.Join(hits, ",")
		status = model.SubtitleStatusAuditBack
		subtitle.PubTime = time.Now().Unix()
	}
	subtitle.Status = status
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("UpdateSubtitleStatus(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	s.dao.DelSubtitleDraftCache(c, subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	s.dao.DelSubtitleCache(c, subtitle.Oid, subtitle.ID)
	return
}

func (s *Service) checkToPublish(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		status = model.SubtitleStatusPublish
		hits   []string
	)
	if hits, err = s.checkBfsData(c, subtitle); err != nil && err != ecode.SubtitleSizeLimit {
		log.Error("checkBfsData(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	if err == ecode.SubtitleSizeLimit {
		subtitle.RejectComment = "单条字幕数超过限制"
		status = model.SubtitleStatusAuditBack
	}
	if len(hits) > 0 {
		subtitle.RejectComment = "敏感词:" + strings.Join(hits, ",")
		status = model.SubtitleStatusAuditBack
	}
	subtitle.Status = status
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("UpdateSubtitleStatus(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	if status == model.SubtitleStatusPublish {
		if err = s.auditPass(c, subtitle); err != nil {
			log.Error("auditPass(subtitle:%+v),error(%v)", subtitle, err)
			return
		}
		return
	}
	if err = s.auditReject(c, subtitle); err != nil {
		log.Error("auditReject(subtitle:%+v),error(%v)", subtitle, err)
		return
	}
	return
}

// checkBfsData .
func (s *Service) checkBfsData(c context.Context, subtitle *model.Subtitle) (hits []string, err error) {
	var (
		body *model.SubtitleBody
		bs   []byte
	)
	if bs, err = s.dao.BfsData(c, subtitle.SubtitleURL); err != nil {
		log.Error("BfsData.params(SubtitleURL:%v),error(%v)", subtitle.SubtitleURL, err)
		return
	}
	body = &model.SubtitleBody{}
	if err = json.Unmarshal(bs, body); err != nil {
		log.Error("checkToAudit.Unmarshal,error(%v)", err)
		return
	}
	if hits, err = s.checkFilter(c, body); err != nil {
		log.Error("checkFilter(body:%+v),error(%v)", body, err)
		return
	}
	return
}

// checkFilter .
func (s *Service) checkFilter(c context.Context, body *model.SubtitleBody) (hits []string, err error) {
	var (
		msgMap  map[string]string
		msgMaps []map[string]string
		reply   *filterMdl.MHitReply
		hitMap  map[string]struct{}
	)
	msgMap = make(map[string]string)
	for idx, item := range body.Bodys {
		if len(item.Content) > _contentSizeLimit {
			err = ecode.SubtitleSizeLimit
			return
		}
		msgMap[fmt.Sprint(idx)] = item.Content
		if (idx+1)%_filterCapacity == 0 {
			msgMaps = append(msgMaps, msgMap)
			msgMap = make(map[string]string)
		}
	}
	if len(msgMap) > 0 {
		msgMaps = append(msgMaps, msgMap)
	}
	hitMap = make(map[string]struct{})
	for _, msgMap = range msgMaps {
		if reply, err = s.filterRPC.MHit(c, &filterMdl.MHitReq{
			Area:   _filterAreaSubtitle,
			MsgMap: msgMap,
		}); err != nil {
			log.Error("checkFilter(msgMap:%+v),error(%v)", msgMap, err)
			return
		}
		for _, rl := range reply.GetRMap() {
			for _, hit := range rl.GetHits() {
				hitMap[hit] = struct{}{}
			}
		}
	}
	for k := range hitMap {
		hits = append(hits, k)
	}
	return
}

// auditReject subtitle reject
func (s *Service) auditReject(c context.Context, subtitle *model.Subtitle) (err error) {
	subtitle.Status = model.SubtitleStatusAuditBack
	if err = s.dao.UpdateSubtitle(c, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	s.dao.DelSubtitleDraftCache(context.Background(), subtitle.Oid, subtitle.Type, subtitle.Mid, subtitle.Lan)
	s.dao.DelSubtitleCache(context.Background(), subtitle.Oid, subtitle.ID)
	return
}

// auditPass .
func (s *Service) auditPass(c context.Context, subtitle *model.Subtitle) (err error) {
	var (
		tx          *sql.Tx
		subtitlePub *model.SubtitlePub
	)
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Error("params(subtitle:%+v).err(%v)", subtitle, err)
			return
		}
		if err = tx.Commit(); err != nil {
			log.Error("params(subtitle:%+v).err(%v)", subtitle, err)
			return
		}
	}()
	subtitle.RejectComment = ""
	if tx, err = s.dao.BeginBiliDMTran(c); err != nil {
		log.Error("error(%v)", err)
		return
	}
	if err = s.dao.TxUpdateSubtitle(tx, subtitle); err != nil {
		log.Error("params(%+v).error(%v)", subtitle, err)
		return
	}
	subtitlePub = &model.SubtitlePub{
		Oid:        subtitle.Oid,
		Type:       subtitle.Type,
		Lan:        subtitle.Lan,
		SubtitleID: subtitle.ID,
	}
	if err = s.dao.TxAddSubtitlePub(tx, subtitlePub); err != nil {
		log.Error("params(%+v).error(%v)", subtitlePub, err)
		return
	}
	if err = s.dao.DelSubtitleCache(c, subtitle.Oid, subtitle.ID); err != nil {
		log.Error("DelSubtitleCache.params(subtitle:%+v).err(%v)", subtitle, err)
		return
	}
	if err = s.dao.DelVideoSubtitleCache(c, subtitle.Oid, subtitle.Type); err != nil {
		log.Error("DelVideoSubtitleCache.params(subtitle:%+v).err(%v)", subtitle, err)
		return
	}
	return
}
