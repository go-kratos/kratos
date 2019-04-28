package http

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"go-common/app/admin/main/creative/model/academy"
	"go-common/app/interface/openplatform/article/model"
	"go-common/app/service/main/archive/api"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"go-common/library/net/metadata"

	"github.com/davecgh/go-spew/spew"
	"github.com/jinzhu/gorm"
)

func addArc(c *bm.Context) {
	var err error
	v := new(struct {
		OID          int64  `form:"oid" validate:"required"`
		Business     int8   `form:"business"`
		Comment      string `form:"comment"`
		CourseTID    int64  `form:"course_tid" validate:"required"`
		OperTID      int64  `form:"oper_tid" validate:"required"`
		ClassTID     string `form:"class_tid" validate:"required"`
		ArticleTID   int64  `form:"article_tid"`
		RecommendTID int64  `form:"recommend_tid" validate:"required"`
	})
	ip := metadata.String(c, metadata.RemoteIP)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if arc, err := checkExist(v.OID); err == nil && arc.OID != 0 {
		c.JSON(nil, ecode.CreativeAcademyOIDExistErr)
		return
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "添加单个视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "添加单个专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts([]int64{v.OID})})
	c.JSON(nil, bulkInsertArcs(c, []int64{v.OID}, v.Business, v.CourseTID, v.OperTID, v.ArticleTID, v.RecommendTID, v.ClassTID, v.Comment, ip))
}

func upArcTag(c *bm.Context) {
	var err error
	v := new(struct {
		OID          int64  `form:"oid" validate:"required"`
		Business     int8   `form:"business"`
		Comment      string `form:"comment"`
		CourseTID    int64  `form:"course_tid" validate:"required"`
		OperTID      int64  `form:"oper_tid" validate:"required"`
		ClassTID     string `form:"class_tid" validate:"required"`
		ArticleTID   int64  `form:"article_tid"`
		RecommendTID int64  `form:"recommend_tid" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = checkExist(v.OID); err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "更新单个视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "更新单个专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts([]int64{v.OID})})
	c.JSON(nil, bulkUpdateArcs([]int64{v.OID}, v.Business, v.CourseTID, v.OperTID, v.ArticleTID, v.RecommendTID, v.ClassTID, v.Comment))
}

func removeArcTag(c *bm.Context) {
	var err error
	v := new(struct {
		OID      int64 `form:"oid" validate:"required"`
		Business int8  `form:"business"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = checkExist(v.OID); err != nil {
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "移除单个视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "移除单个专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts([]int64{v.OID})})
	c.JSON(nil, bulkRemoveArcs([]int64{v.OID}, v.Business))
}

func viewArc(c *bm.Context) {
	var err error
	v := new(struct {
		OID      int64 `form:"oid"`
		Business int8  `form:"business"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if _, err = checkExist(v.OID); err != nil {
		c.JSON(nil, err)
		return
	}
	ap := &academy.EsParam{
		OID:      v.OID,
		Business: v.Business,
		State:    academy.DefaultState,
		Pn:       1,
		Ps:       1,
		IP:       metadata.String(c, metadata.RemoteIP),
	}
	res, err := search(c, ap)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res == nil || len(res.Items) == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(res.Items, nil)
}

func listArc(c *bm.Context) {
	var (
		err     error
		tids    []int64
		tidsMap map[int][]int64
	)
	v := new(struct {
		OID       int64  `form:"oid"`
		Keyword   string `form:"keyword"`
		Uname     string `form:"uname"`
		Business  int8   `form:"business"`
		TID       string `form:"tids"`
		State     int    `form:"state" default:"2018"`
		Copyright int    `form:"copyright"`
		Pn        int    `form:"pn" validate:"required,min=1"`
		Ps        int    `form:"ps" validate:"required,min=1"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn == 0 {
		v.Pn = 1
	}
	if v.Ps > 20 {
		v.Ps = 20
	}

	if v.TID != "" {
		if tids, err = xstr.SplitInts(v.TID); err != nil {
			log.Error("strconv.ParseInt(%s) error(%v)", v.TID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
		if len(tids) >= 0 {
			tags, getMapErr := getTagParentChildMap()
			if getMapErr != nil {
				c.JSON(nil, getMapErr)
				return
			}
			tidsMap = filterTIDs(tids, tags)
		}
	}

	ap := &academy.EsParam{
		OID:       v.OID,
		Keyword:   v.Keyword,
		Business:  v.Business,
		Uname:     v.Uname,
		TID:       tids,
		Copyright: v.Copyright,
		State:     v.State,
		Pn:        v.Pn,
		Ps:        v.Ps,
		IP:        metadata.String(c, metadata.RemoteIP),
		TidsMap:   tidsMap,
	}
	res, err := search(c, ap)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    res,
	}))
}

func bindArcInfo(c context.Context, oids []int64, bs int8, ip string) (res map[int64]*academy.ArchiveMeta, err error) {
	res = make(map[int64]*academy.ArchiveMeta)
	if bs == academy.BusinessForArchvie {
		arcs, err := svc.Archives(c, oids)
		if err != nil {
			return nil, err
		}
		stat, err := svc.Stats(c, oids, ip)
		if err != nil {
			log.Error("s.arc.Stats oids(%+v)|business(%d)|error(%v)", oids, bs, err)
			return nil, err
		}
		for _, oid := range oids {
			a := &academy.ArchiveMeta{}
			if v, ok := arcs[oid]; ok && v != nil {
				a.OID = oid
				a.Title = v.Title
				a.State = v.State
				a.Type = v.TypeName
				a.Cover = v.Pic
				a.UName = v.Author.Name
				if t, ok := stat[oid]; ok && t != nil {
					a.Hot = countArcHot(t, int64(v.PubDate))
				}
				res[oid] = a
			}
		}
	} else if bs == academy.BusinessForArticle {
		arts, err := svc.Articles(c, oids)
		if err != nil {
			return nil, err
		}
		for _, oid := range oids {
			if v, ok := arts[oid]; ok && v != nil {
				a := &academy.ArchiveMeta{
					OID:   oid,
					Title: v.Title,
					State: v.State,
				}
				if v.Category != nil {
					a.Type = v.Category.Name
				}
				if len(v.ImageURLs) > 0 {
					a.Cover = v.ImageURLs[0]
				}
				if v.Author != nil {
					a.UName = v.Author.Name
				}
				a.Hot = countArtHot(v)
				res[oid] = a
			}
		}
	}
	return
}

func bindTags(oidTIDsMap map[int64][]int64) (res map[int64]map[int][]*academy.TagMeta, err error) {
	res = make(map[int64]map[int][]*academy.TagMeta)
	for oid, tids := range oidTIDsMap {
		var tags []*academy.Tag
		if err = svc.DB.Model(&academy.Tag{}).Where("id in (?)", tids).Find(&tags).Error; err != nil {
			log.Error("creative-admin bindTags error(%v)", err)
			return
		}
		oidTIDs := make(map[int][]*academy.TagMeta)
		ctgs := make(map[int64][]*academy.TagMeta)
		for _, v := range tags {
			tg := renderTag(v)
			switch tg.Type {
			case academy.Course:
				oidTIDs[academy.Course] = append(oidTIDs[academy.Course], tg)
			case academy.Operation:
				oidTIDs[academy.Operation] = append(oidTIDs[academy.Operation], tg)
			case academy.Classify:
				ctgs[tg.ParentID] = append(ctgs[tg.ParentID], tg)
			case academy.ArticleClass:
				oidTIDs[academy.ArticleClass] = append(oidTIDs[academy.ArticleClass], tg)
			case academy.Recommend:
				oidTIDs[academy.Recommend] = append(oidTIDs[academy.Recommend], tg)
			}
		}
		for pid, tgs := range ctgs {
			parent := &academy.Tag{}
			if err = svc.DB.Model(&academy.Tag{}).Where("id = ?", pid).Find(parent).Error; err != nil {
				log.Error("creative-admin bindTags get parent tag error(%v)", err)
				continue
			}
			p := renderTag(parent)
			p.Children = tgs
			oidTIDs[academy.Classify] = append(oidTIDs[academy.Classify], p)
		}
		res[oid] = oidTIDs
	}
	return
}

func batchAddArc(c *bm.Context) {
	var (
		err                    error
		oids, newOIDs, oldOIDs []int64
	)
	v := new(struct {
		OIDs         string `form:"oids" validate:"required"`
		Business     int8   `form:"business"`
		Comment      string `form:"comment"`
		CourseTID    int64  `form:"course_tid" validate:"required"`
		OperTID      int64  `form:"oper_tid" validate:"required"`
		ClassTID     string `form:"class_tid" validate:"required"`
		ArticleTID   int64  `form:"article_tid"`
		RecommendTID int64  `form:"recommend_tid" validate:"required"`
	})
	ip := metadata.String(c, metadata.RemoteIP)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OIDs != "" {
		if oids, err = xstr.SplitInts(v.OIDs); err != nil {
			log.Error("batchAddArc xstr.SplitInts OIDs(%+v)|error(%v)", v.OIDs, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	for _, oid := range oids {
		if arc, err := checkExist(oid); err == nil && arc.OID != 0 { //表里已存在该稿件
			oldOIDs = append(oldOIDs, oid)
		} else {
			newOIDs = append(newOIDs, oid)
		}
	}
	if len(newOIDs) == 0 {
		log.Error("batchAddArc oldOIDs(%+v)", oldOIDs)
		c.JSON(nil, ecode.CreativeAcademyOIDExistErr)
		return
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "批量添加视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "批量添加专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts(newOIDs)})
	c.JSON(nil, bulkInsertArcs(c, newOIDs, v.Business, v.CourseTID, v.OperTID, v.ArticleTID, v.RecommendTID, v.ClassTID, v.Comment, ip))
}

//批量插入稿件表
func bulkInsertArcs(c *bm.Context, oids []int64, bs int8, courseTID, operTID, articleTID, recommendTID int64, classTID, comment, ip string) (err error) {
	arcInfo, err := bindArcInfo(c, oids, bs, ip)
	if err != nil {
		return err
	}
	arcs := make([]*academy.Archive, 0, len(oids))
	for _, oid := range oids {
		a, ok := arcInfo[oid]
		if !ok || a == nil { //校验oid是否有效
			if bs == academy.BusinessForArchvie {
				err = ecode.CreativeArcServiceErr
			} else if bs == academy.BusinessForArticle {
				err = ecode.CreativeArticleRPCErr
			}
			log.Error("bulkInsertArcs add archive with invalid oid(%d)|business(%d)", oid, bs)
			return
		}
		arcs = append(arcs, setArcParam(oid, a.Hot, bs, comment))
	}
	valArcs := make([]string, 0, len(arcs))
	valArcArgs := make([]interface{}, 0)
	for _, v := range arcs {
		valArcs = append(valArcs, "(?, ?, ?, ?, ?, ?, ?)")
		valArcArgs = append(valArcArgs, v.OID, v.Hot, v.Business, v.State, v.Comment, v.CTime, v.MTime)
	}
	//批量插入稿件标签关联表
	ctids, err := xstr.SplitInts(classTID)
	if err != nil {
		return
	}
	tags := make([]*academy.ArchiveTag, 0)
	for _, oid := range oids {
		tags = append(tags, setTagParam(oid, courseTID, bs), setTagParam(oid, operTID, bs), setTagParam(oid, recommendTID, bs))
		for _, cid := range ctids { //分类标签支持绑定多个二级标签
			tags = append(tags, setTagParam(oid, cid, bs))
		}
		if bs == academy.BusinessForArticle && articleTID > 0 { //专栏特殊分类
			tags = append(tags, setTagParam(oid, articleTID, bs))
		}
	}
	valTags := make([]string, 0)
	valTagArgs := make([]interface{}, 0)
	for _, v := range tags {
		valTags = append(valTags, "(?, ?, ?, ?, ?, ?)")
		valTagArgs = append(valTagArgs, v.OID, v.TID, v.State, v.CTime, v.MTime, v.Business)
	}
	sqlArcStr := fmt.Sprintf("INSERT INTO academy_archive (oid, hot, business, state, comment, ctime, mtime) VALUES %s ON DUPLICATE KEY UPDATE state=0, comment=VALUES(comment), mtime=VALUES(mtime)", strings.Join(valArcs, ","))
	sqlTagStr := fmt.Sprintf("INSERT INTO academy_archive_tag (oid, tid, state, ctime, mtime, business) VALUES %s ON DUPLICATE KEY UPDATE state=0, mtime=VALUES(mtime)", strings.Join(valTags, ","))
	tx := svc.DB.Begin()
	if err = tx.Exec(sqlArcStr, valArcArgs...).Error; err != nil {
		log.Error("creative-admin bulkInsertArcs error(%v)", err)
		tx.Rollback()
		return
	}
	if err = tx.Exec(sqlTagStr, valTagArgs...).Error; err != nil {
		log.Error("creative-admin bulkInsertArcs error(%v)", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return nil
}

func setArcParam(oid, hot int64, bs int8, comment string) (arc *academy.Archive) {
	now := time.Now().Format("2006-01-02 15:04:05")
	arc = &academy.Archive{
		OID:      oid,
		Business: bs,
		CTime:    now,
		MTime:    now,
		Comment:  comment,
		Hot:      hot,
		State:    academy.StateNormal,
	}
	return
}

func setTagParam(oid, tid int64, bs int8) (res *academy.ArchiveTag) {
	now := time.Now().Format("2006-01-02 15:04:05")
	res = &academy.ArchiveTag{
		OID:      oid,
		TID:      tid,
		CTime:    now,
		MTime:    now,
		Business: bs,
		State:    academy.StateNormal,
	}
	return
}

func batchUpArc(c *bm.Context) {
	var (
		err  error
		oids []int64
	)
	v := new(struct {
		OIDs         string `form:"oids" validate:"required"`
		Business     int8   `form:"business"`
		Comment      string `form:"comment"`
		CourseTID    int64  `form:"course_tid" validate:"required"`
		OperTID      int64  `form:"oper_tid" validate:"required"`
		ClassTID     string `form:"class_tid" validate:"required"`
		ArticleTID   int64  `form:"article_tid"`
		RecommendTID int64  `form:"recommend_tid" validate:"required"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OIDs != "" {
		if oids, err = xstr.SplitInts(v.OIDs); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	for _, oid := range oids {
		if _, err = checkExist(oid); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "批量更新视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "批量更新专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts(oids)})
	c.JSON(nil, bulkUpdateArcs(oids, v.Business, v.CourseTID, v.OperTID, v.ArticleTID, v.RecommendTID, v.ClassTID, v.Comment))
}

func bulkUpdateArcs(oids []int64, bs int8, courseTID, operTID, articleTID, recommendTID int64, classTID, comment string) (err error) {
	var (
		tids                         []int64
		arcTags                      []*academy.ArchiveTags
		arcMapTids                   map[int64][]int64
		delOIDTidsMap, newOIDTidsMap map[int64][]int64
	)
	//更新其他类别标签
	var (
		ats []*academy.ArchiveTags
		ids []int64
	)
	if err = svc.DB.Raw("SELECT b.*, c.type FROM (SELECT t.id,t.tid,a.oid FROM academy_archive AS a LEFT JOIN academy_archive_tag AS t ON t.oid = a.oid WHERE  a.oid IN (?) AND a.business=? AND a.state=0 AND t.state=0) b LEFT JOIN academy_tag AS c  ON c.id=b.tid WHERE c.type!=?", oids, bs, academy.Classify).Find(&ats).Error; err != nil {
		log.Error("creative-admin bulkUpdateArcs get all archive tags error(%v)", err)
		return
	}
	if len(ats) == 0 {
		return
	}
	ids = make([]int64, 0, len(ats))
	for _, a := range ats {
		if a == nil {
			log.Error("creative-admin bulkUpdateArcs update other tags get nil a(%+v)", a)
			return
		}
		ids = append(ids, a.ID)
	}
	tx := svc.DB.Begin()
	// 对于提交上来的其他标签先统一删除再统一插入
	if err = tx.Model(&academy.ArchiveTag{}).Where("id IN (?)", ids).
		Updates(map[string]interface{}{
			"state": academy.StateRemove,
		}).Error; err != nil {
		log.Error("bulkUpdateArcs first delete other tags by ids(%v)|error(%v)", ids, err)
		tx.Rollback()
		return
	}
	tags := make([]*academy.ArchiveTag, 0)
	for _, oid := range oids {
		tags = append(tags, setTagParam(oid, courseTID, bs), setTagParam(oid, operTID, bs), setTagParam(oid, recommendTID, bs))
		if bs == academy.BusinessForArticle && articleTID > 0 { //专栏特殊分类
			tags = append(tags, setTagParam(oid, articleTID, bs))
		}
	}
	valTags := make([]string, 0)
	valTagArgs := make([]interface{}, 0)
	for _, v := range tags {
		valTags = append(valTags, "(?, ?, ?, ?, ?, ?)")
		valTagArgs = append(valTagArgs, v.OID, v.TID, v.State, v.CTime, v.MTime, v.Business)
	}
	sqlTagStr := fmt.Sprintf("INSERT INTO academy_archive_tag (oid, tid, state, ctime, mtime, business) VALUES %s ON DUPLICATE KEY UPDATE state=0,oid=VALUES(oid),tid=VALUES(tid),business=VALUES(business)", strings.Join(valTags, ","))
	if err = tx.Exec(sqlTagStr, valTagArgs...).Error; err != nil {
		log.Error("bulkUpdateArcs update other tags valTagArgs(%+v)", valTagArgs...)
		tx.Rollback()
		return
	}
	//处理分类标签
	tids, err = xstr.SplitInts(classTID)
	if err != nil {
		return err
	}
	if err = svc.DB.Raw("SELECT b.*, c.type FROM (SELECT t.id,t.tid,a.oid FROM academy_archive AS a LEFT JOIN academy_archive_tag AS t ON t.oid = a.oid WHERE  a.oid IN (?) AND a.business=? AND a.state=0 AND t.state=0) b LEFT JOIN  academy_tag AS c  ON c.id=b.tid WHERE c.type=?", oids, bs, academy.Classify).Find(&arcTags).Error; err != nil {
		log.Error("creative-admin bulkUpdateArcs archive class tags  error(%v)", err)
		return
	}
	arcMapTids = make(map[int64][]int64)
	for _, v := range arcTags {
		arcMapTids[v.OID] = append(arcMapTids[v.OID], v.TID)
	}
	newOIDTidsMap = make(map[int64][]int64)
	delOIDTidsMap = make(map[int64][]int64)
	for oid, ids := range arcMapTids {
		isInDB := make(map[int64]int64)
		for _, id := range ids {
			isInDB[id] = id
		}
		log.Info("creative-admin bulkUpdateArcs oid(%d)|isInDB(%+v)", oid, isInDB)
		isInSub := make(map[int64]int64)
		for _, tid := range tids {
			isInSub[tid] = tid
			if _, ok := isInDB[tid]; !ok {
				newOIDTidsMap[oid] = append(newOIDTidsMap[oid], tid)
			}
		}
		log.Info("creative-admin bulkUpdateArcs oid(%d)|isInSubmit(%+v)", oid, isInSub)
		for _, id := range ids {
			if _, ok := isInSub[id]; !ok {
				delOIDTidsMap[oid] = append(delOIDTidsMap[oid], id)
			}
		}
	}
	if len(newOIDTidsMap) > 0 { //insert on update
		tags := make([]*academy.ArchiveTag, 0)
		for oid, tgs := range newOIDTidsMap { //分类标签支持绑定多个二级标签
			for _, cid := range tgs {
				tags = append(tags, setTagParam(oid, cid, bs))
			}
		}
		valTags := make([]string, 0)
		valTagArgs := make([]interface{}, 0)
		for _, v := range tags {
			valTags = append(valTags, "(?, ?, ?, ?, ?, ?)")
			valTagArgs = append(valTagArgs, v.OID, v.TID, v.State, v.CTime, v.MTime, v.Business)
		}
		sqlTagStr := fmt.Sprintf("INSERT INTO academy_archive_tag (oid, tid, state, ctime, mtime, business) VALUES %s ON DUPLICATE KEY UPDATE state=0, oid=VALUES(oid), tid=VALUES(tid), business=VALUES(business), mtime=VALUES(mtime)", strings.Join(valTags, ","))
		if err = tx.Exec(sqlTagStr, valTagArgs...).Error; err != nil {
			log.Error("creative-admin bulkUpdateArcs insert new class tags error(%v)", err)
			tx.Rollback()
			return
		}
	}
	if len(delOIDTidsMap) > 0 { //delete
		delMapID := make(map[int64]int64)
		for _, tgs := range delOIDTidsMap {
			for _, tid := range tgs {
				delMapID[tid] = tid
			}
		}
		delIDs := make([]int64, 0)
		for _, a := range arcTags {
			if _, ok := delMapID[a.TID]; ok {
				delIDs = append(delIDs, a.ID)
			}
		}
		if err = tx.Model(&academy.ArchiveTag{}).Where("id IN (?)", delIDs).
			Updates(map[string]interface{}{
				"state": academy.StateRemove,
				"mtime": time.Now().Format("2006-01-02 15:04:05"),
			}).Error; err != nil {
			log.Error("creative-admin bulkUpdateArcs delete class tags by ids(%v)|error(%v)", delIDs, err)
			tx.Rollback()
			return
		}
	}
	// 统一更新comment
	if err = tx.Model(&academy.Archive{}).Where("business = ?", bs).Where("oid IN (?)", oids).
		Updates(map[string]interface{}{
			"comment": comment,
			"mtime":   time.Now().Format("2006-01-02 15:04:05"),
		}).Error; err != nil {
		log.Error("creative-admin bulkUpdateArcs update all comment error(%v)", err)
		tx.Rollback()
		return
	}
	tx.Commit()
	return
}

func batchRemoveArc(c *bm.Context) {
	var (
		err  error
		oids []int64
	)
	v := new(struct {
		OIDs     string `form:"oids" validate:"required"`
		Business int8   `form:"business"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.OIDs != "" {
		if oids, err = xstr.SplitInts(v.OIDs); err != nil {
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	for _, oid := range oids {
		if _, err = checkExist(oid); err != nil {
			c.JSON(nil, err)
			return
		}
	}
	uid, uname := getUIDName(c)
	action := ""
	if v.Business == academy.BusinessForArchvie {
		action = "批量移除视频稿件"
	} else if v.Business == academy.BusinessForArticle {
		action = "批量移除专栏稿件"
	}
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: action, OIDs: xstr.JoinInts(oids)})
	c.JSON(nil, bulkRemoveArcs(oids, v.Business))
}

func bulkRemoveArcs(oids []int64, bs int8) (err error) {
	var (
		db  *gorm.DB
		now = time.Now().Format("2006-01-02 15:04:05")
		ats []*academy.ArchiveTags
		IDs []int64
	)
	db = svc.DB.Raw("SELECT t.id,t.tid,a.oid,a.business FROM academy_archive AS a LEFT JOIN academy_archive_tag AS t ON t.oid = a.oid WHERE  a.oid IN (?) AND a.business=?", oids, bs)
	if err = db.Find(&ats).Error; err != nil {
		log.Error("creative-admin bulkRemoveArcs error(%v)", err)
		return
	}
	if len(ats) == 0 {
		return
	}
	for _, a := range ats {
		IDs = append(IDs, a.ID)
	}
	tx := svc.DB.Begin()
	if err = tx.Model(&academy.Archive{}).Where("oid IN (?) AND business=?", oids, bs).Updates(map[string]interface{}{
		"mtime": now,
		"state": academy.StateRemove,
	}).Error; err != nil {
		log.Error("creative-admin bulkRemoveArcs error(%v)", err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&academy.ArchiveTag{}).Where("id IN (?)", IDs).Updates(map[string]interface{}{
		"mtime": now,
		"state": academy.StateRemove,
	}).Error; err != nil {
		tx.Rollback()
		log.Error("creative-admin bulkRemoveArcs error(%v)", err)
		return
	}
	tx.Commit()
	return
}

func arcCountByTids(tids []int64) (res map[int64]int, err error) {
	var (
		countSQL = "SELECT tid, count(DISTINCT oid) AS count  FROM academy_archive_tag WHERE state=0 AND tid IN (?) GROUP BY tid"
		ats      []*academy.ArchiveCount
	)
	if err = svc.DB.Raw(countSQL, tids).Find(&ats).Error; err != nil {
		log.Error("creative-admin get arcCountByTids error(%v)", err)
		return
	}
	if len(ats) == 0 {
		return
	}
	res = make(map[int64]int)
	for _, a := range ats {
		res[a.TID] = a.Count
	}
	return
}

func checkExist(oid int64) (arc *academy.Archive, err error) {
	arc = &academy.Archive{}
	err = svc.DB.Model(&academy.Archive{}).Where("state=?", academy.StateNormal).Where("oid=?", oid).Find(arc).Error
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			log.Error("creative-admin checkExist oid(%d)|error(%v)", oid, err)
		}
	}
	return
}

func fixArchive(c *bm.Context) {
	var err error
	v := new(struct {
		ID       int64 `form:"id"`
		OID      int64 `form:"oid"`
		Business int8  `form:"business"`
		State    int8  `form:"state"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if err = svc.DB.Model(&academy.Archive{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"oid":      v.OID,
		"business": v.Business,
		"state":    v.State,
		"mtime":    time.Now().Format("2006-01-02 15:04:05"),
	}).Error; err != nil {
		log.Error("creative-admin fixArchive error(%v)", err)
	}
	c.JSON(nil, err)
}

func search(c *bm.Context, ap *academy.EsParam) (res *academy.Archives, err error) {
	var (
		oids        []int64
		tempMapTIDs map[int64][]int64
		arcs        []*academy.ArchiveOrigin
		items       []*academy.ArchiveMeta
		com         = make(map[int64]string)
	)
	res = &academy.Archives{
		Items: []*academy.ArchiveMeta{},
		Pager: &academy.Pager{},
	}
	sear, err := svc.ArchivesWithES(c, ap)
	if err != nil {
		log.Error("search svc.ArchivesWithES error(%v)", err)
		return
	}
	arcs = make([]*academy.ArchiveOrigin, 0)
	if sear == nil || len(sear.Result) == 0 {
		return
	}
	tempMapTIDs = make(map[int64][]int64)
	for _, v := range sear.Result {
		oids = append(oids, v.OID)
		tempMapTIDs[v.OID] = v.TID
	}
	res.Pager.Total = sear.Page.Total
	res.Pager.Num = sear.Page.Num
	res.Pager.Size = sear.Page.Size
	g, _ := errgroup.WithContext(c)
	var (
		arcInfo     map[int64]*academy.ArchiveMeta
		tagInfo     map[int64]map[int][]*academy.TagMeta
		bindMapTIDs map[int64][]int64
	)
	bindMapTIDs = make(map[int64][]int64)
	for _, oid := range oids {
		if v, ok := tempMapTIDs[oid]; ok {
			bindMapTIDs[oid] = v
		}
	}
	g.Go(func() error {
		arcInfo, err = bindArcInfo(c, oids, ap.Business, metadata.String(c, metadata.RemoteIP))
		return err
	})
	g.Go(func() error {
		com, err = bindArcComment(oids, ap.Business)
		return err
	})
	g.Go(func() error {
		tagInfo, err = bindTags(bindMapTIDs)
		return err
	})
	if err = g.Wait(); err != nil {
		return
	}
	log.Info("search arcInfo(%s)", spew.Sdump(arcInfo))
	items = make([]*academy.ArchiveMeta, 0, len(arcs))
	for _, oid := range oids {
		a, ok := arcInfo[oid]
		if !ok || a == nil {
			log.Error("search get archive info error by oid(%d)", oid)
			return
		}
		if v, ok := tagInfo[oid]; ok {
			a.Tags = v
		}
		if co, ok := com[oid]; ok {
			a.Comment = co
		}
		items = append(items, a)
	}
	res.Items = items
	return
}

func bindArcComment(oids []int64, bs int8) (res map[int64]string, err error) {
	var arcs []*academy.Archive
	if err = svc.DB.Model(&academy.Archive{}).Where("oid in(?)", oids).Where("business=?", bs).Group("oid").Find(&arcs).Error; err != nil {
		log.Error("bindArcComment d.DB.Model oids(%+v)|business(%d)|error(%v)", oids, bs, err)
		return
	}
	res = make(map[int64]string)
	for _, v := range arcs {
		res[v.OID] = v.Comment
	}
	return
}

//countArcHot 视频=硬币*0.4+收藏*0.3+弹幕*0.4+评论*0.4+播放*0.25+点赞*0.4+分享*0.6 最新视频（一天内发布）提权[总值*1.5]
func countArcHot(t *api.Stat, ptime int64) int64 {
	if t == nil {
		return 0
	}
	hot := float64(t.Coin)*0.4 +
		float64(t.Fav)*0.3 +
		float64(t.Danmaku)*0.4 +
		float64(t.Reply)*0.4 +
		float64(t.View)*0.25 +
		float64(t.Like)*0.4 +
		float64(t.Share)*0.6
	if ptime >= time.Now().AddDate(0, 0, -1).Unix() && ptime <= time.Now().Unix() {
		hot *= 1.5
	}
	return int64(math.Floor(hot))
}

// countArtHot 专栏=硬币*0.4+收藏*0.3+评论*0.4+阅读*0.25+点赞*0.4+分享*0.6 最新专栏（一天内发布）提权[总值*1.5]
func countArtHot(t *model.Meta) int64 {
	if t.Stats == nil {
		return 0
	}
	hot := float64(t.Stats.Coin)*0.4 +
		float64(t.Stats.Favorite)*0.3 +
		float64(t.Stats.Reply)*0.4 +
		float64(t.Stats.View)*0.25 +
		float64(t.Stats.Like)*0.4 +
		float64(t.Stats.Share)*0.6
	if int64(t.PublishTime) >= time.Now().AddDate(0, 0, -1).Unix() && int64(t.PublishTime) <= time.Now().Unix() {
		hot *= 1.5
	}
	return int64(math.Floor(hot))
}

// getParentChildMap
func getTagParentChildMap() (res map[int64]*academy.Tag, err error) {
	var (
		db   *gorm.DB
		tags []*academy.Tag
	)

	db = svc.DB.Order("rank ASC").Find(&tags)
	if err = db.Error; err != nil {
		log.Error("creative-admin getTagParentChildMap error(%v)", err)
		return
	}
	res = make(map[int64]*academy.Tag)
	for _, t := range tags {
		res[t.ID] = t
	}
	for _, v := range res {
		if v == nil {
			continue
		}
		if v.ParentID == 0 {
			for _, t := range tags {
				if t == nil {
					continue
				}
				if t.ParentID == v.ID {
					v.Children = append(v.Children, t)
				}
			}
		}
	}

	return
}

// filterTIDs
func filterTIDs(tids []int64, parentChildMap map[int64]*academy.Tag) (res map[int][]int64) {
	if len(tids) == 0 {
		return
	}
	log.Info("s.filterTIDs origin tids(%+v)", tids)
	res = make(map[int][]int64)
	ochs := make([]int64, 0) //原始提交的二级标签
	ops := make([]int64, 0)  //原始提交的一级标签
	qchs := make([]int64, 0) //通过一级标签查询出来的二级标签
	for _, id := range tids {
		t, ok := parentChildMap[id]
		if !ok || t == nil {
			continue
		}
		if t.Type == academy.Classify {
			if t.ParentID != 0 { //原始提交的二级标签
				ochs = append(ochs, id)
			} else if t.ParentID == 0 && len(t.Children) > 0 { //通过一级标签查询出来的二级标签
				for _, v := range t.Children {
					qchs = append(qchs, v.ID)
				}
			} else if t.ParentID == 0 && len(t.Children) == 0 {
				ops = append(ops, id)
			}
		} else {
			res[int(t.Type)] = append(res[int(t.Type)], id)
		}
	}
	if len(ochs) > 0 { //如果分类标签中提交了原始的二级标签则认为按该二级标签进行筛选，如果可以查询到二级标签认为筛选全部二级，否则一级参与查询.
		res[academy.Classify] = ochs
	} else if len(qchs) > 0 {
		res[academy.Classify] = qchs
	} else if len(ops) > 0 {
		res[academy.Classify] = ops
	}
	log.Info("s.filterTIDs res(%s)", spew.Sdump(res))
	return
}
