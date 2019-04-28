package http

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"go-common/app/admin/main/creative/model/academy"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

func addTag(c *bm.Context) {
	var (
		err     error
		linkIDs []int64
		now     = time.Now().Format("2006-01-02 15:04:05")
	)
	v := new(struct {
		ParentID int64  `form:"parent_id"`
		Business int8   `form:"business"`
		Type     int8   `form:"type"`
		Desc     string `form:"desc"`
		Name     string `form:"name"`
		LinkID   string `form:"link_id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.LinkID != "" {
		if linkIDs, err = xstr.SplitInts(v.LinkID); err != nil {
			log.Error("addTag xstr.SplitInts h5 linkid(%+v)|error(%v)", v.LinkID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}
	m := &academy.Tag{
		ParentID: v.ParentID,
		State:    academy.StateUnBlock,
		Type:     v.Type,
		Business: v.Business,
		Name:     v.Name,
		Desc:     v.Desc,
		CTime:    now,
		MTime:    now,
	}
	tx := svc.DB.Begin()
	if err = tx.Create(m).Error; err != nil {
		log.Error("creative-admin academy addTag error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&academy.Tag{}).Where("id=?", m.ID).Updates(map[string]interface{}{
		"rank": m.ID,
	}).Error; err != nil {
		log.Error("creative-admin academy addTag update rank error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}

	if len(linkIDs) > 0 { //插入h5标签关联的web标签
		valLink := make([]string, 0, len(linkIDs))
		valLinkArgs := make([]interface{}, 0)
		for _, lid := range linkIDs {
			valLink = append(valLink, "(?, ?, ?, ?)")
			valLinkArgs = append(valLinkArgs, m.ID, lid, now, now)
		}
		sqlLinkStr := fmt.Sprintf("INSERT INTO academy_tag_link (tid, link_id, ctime, mtime) VALUES %s", strings.Join(valLink, ","))
		if err = tx.Exec(sqlLinkStr, valLinkArgs...).Error; err != nil {
			log.Error("academy bulk insert h5 link tag error(%v)", err)
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	tid := m.ID
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "添加标签", TID: tid, OName: m.Name})
	c.JSON(map[string]interface{}{
		"id": tid,
	}, nil)
}

func upTag(c *bm.Context) {
	var (
		tg      = &academy.Tag{}
		linkIDs []int64
		err     error
		now     = time.Now().Format("2006-01-02 15:04:05")
		tid     int64
	)

	v := new(struct {
		ID     int64  `form:"id"`
		Desc   string `form:"desc"`
		Name   string `form:"name"`
		LinkID string `form:"link_id"`
	})
	if err = c.Bind(v); err != nil {
		return
	}

	if v.LinkID != "" {
		if linkIDs, err = xstr.SplitInts(v.LinkID); err != nil {
			log.Error("upTag xstr.SplitInts h5 linkid(%+v)|error(%v)", v.LinkID, err)
			c.JSON(nil, ecode.RequestErr)
			return
		}
	}

	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	tid = v.ID

	if err = svc.DB.Model(&academy.Tag{}).Where("id=?", tid).Find(tg).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if tg == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}

	tx := svc.DB.Begin()
	if err = tx.Model(&academy.Tag{}).Where("id=?", tid).Updates(map[string]interface{}{
		"name": v.Name,
		"desc": v.Desc,
	}).Error; err != nil {
		log.Error("academy upTag error(%v)", err)
		tx.Rollback()
		return
	}

	if len(linkIDs) > 0 { //插入h5标签关联的web标签
		// 对于提交上来的其他标签先统一删除再统一插入
		var tl academy.TagLink
		if err = tx.Where("tid =?", tid).Delete(&tl).Error; err != nil {
			log.Error("upTag TagLink delete by tid(%v)|error(%v)", tid, err)
			tx.Rollback()
			return
		}

		valLink := make([]string, 0, len(linkIDs))
		valLinkArgs := make([]interface{}, 0)
		for _, lid := range linkIDs {
			valLink = append(valLink, "(?, ?, ?, ?)")
			valLinkArgs = append(valLinkArgs, tid, lid, now, now)
		}

		sqlLinkStr := fmt.Sprintf("INSERT INTO academy_tag_link (tid, link_id, ctime, mtime) VALUES %s ON DUPLICATE KEY UPDATE tid=VALUES(tid),link_id=VALUES(link_id)", strings.Join(valLink, ","))
		if err = tx.Exec(sqlLinkStr, valLinkArgs...).Error; err != nil {
			log.Error("upTag TagLink valLinkArgs(%+v)", valLinkArgs...)
			tx.Rollback()
			return
		}
	}
	tx.Commit()

	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新标签", TID: tid, OName: v.Name})
	c.JSON(nil, err)
}

func bindTag(c *bm.Context) {
	var (
		tg  = &academy.Tag{}
		err error
	)
	v := new(struct {
		ID    int64 `form:"id"`
		State int8  `form:"state"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Model(&academy.Tag{}).Where("id=?", v.ID).Find(tg, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if tg == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Tag{ID: v.ID}).Updates(map[string]interface{}{
		"state": v.State,
		"mtime": time.Now().Format("2006-01-02 15:04:05"),
	}).Error; err != nil {
		log.Error("creative-admin academy delTag error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新标签状态", TID: v.ID, OState: v.State})
	c.JSON(nil, err)
}

func viewTag(c *bm.Context) {
	var (
		err error
		tg  = &academy.Tag{}
	)
	v := new(struct {
		ID       int64 `form:"id"`
		Type     int8  `form:"type"`
		Business int8  `form:"business"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("type = ?", v.Type).Where("business = ?", v.Business).Find(tg, v.ID).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if tg == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	c.JSON(renderTag(tg), nil)
}

func listTag(c *bm.Context) {
	var (
		err  error
		db   *gorm.DB
		tags []*academy.Tag
	)
	v := new(struct {
		Type int8 `form:"type"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Type != 0 {
		db = svc.DB.Where("type = ?", v.Type).Order("rank ASC").Find(&tags)
	} else {
		db = svc.DB.Order("rank ASC").Find(&tags)
	}
	if err = db.Error; err != nil {
		log.Error("creative-admin academy listTag  error(%v)", err)
		c.JSON(nil, err)
		return
	}
	tgs, _ := tagTree(tags, v.Type)
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data": map[string]interface{}{
			"classify": academy.TagClass(),
			"list":     tgs,
		},
	}))
}

func tagTree(tags []*academy.Tag, ty int8) (res map[int8][]*academy.TagMeta, err error) {
	all := make(map[int8]map[int64]*academy.TagMeta)
	top := make(map[int64]*academy.TagMeta)
	res = make(map[int8][]*academy.TagMeta)
	tids := make([]int64, 0, len(tags))

	linkMap, _ := h5RelatedWebTIDsMap() //获取h5标签关联的web标签
	if ty == academy.H5 {
		for _, lids := range linkMap {
			tids = append(tids, lids...) //取所有关联的web标签id
		}
	} else {
		for _, v := range tags {
			if v.Type != academy.H5 {
				tids = append(tids, v.ID) //取所有web标签id
			}
		}
	}

	arcCountMap, err := arcCountByTids(tids) //获取web标签下面的稿件数量
	if err != nil {
		return
	}

	h5ArcCountMap := make(map[int64]int)
	for id, lids := range linkMap {
		count := 0
		for _, lid := range lids {
			if c, ok := arcCountMap[lid]; ok {
				count += c
			}
		}
		h5ArcCountMap[id] = count
	}

	for _, v := range tags { //获取父级节点
		if v == nil {
			continue
		}
		t := renderTag(v)

		if t.ParentID == 0 {
			if c, ok := arcCountMap[t.ID]; ok { //获取除分类标签和h5标签之外的标签关联稿件数量
				t.Count = c
			}

			top[t.ID] = t
			all[t.Type] = top                    //存储一级标签对象id map
			res[t.Type] = append(res[t.Type], t) //存储一级标签对象
		}
	}

	for _, v := range tags { //获取子节点并获取关联稿件数量
		if v == nil {
			continue
		}
		t := renderTag(v)

		p, ok := all[t.Type][t.ParentID] //获取一级标签对象
		if ok && p != nil && p.Type == t.Type {

			if c, ok := arcCountMap[t.ID]; ok { //获取分类标签的二级关联稿件数量
				t.Count = c
			}

			if c, ok := h5ArcCountMap[t.ID]; ok { //获取h5标签二级关联稿件数量
				t.Count = c
			}

			if l, ok := linkMap[v.ID]; ok && v.Type == academy.H5 {
				t.LinkID = l
			}

			p.Children = append(p.Children, t)
		}
	}

	for _, v := range res {
		for _, t := range v { //计算分类标签或者h5标签一级关联稿件数量
			if t != nil && t.ParentID == 0 && (t.Type == academy.H5 || t.Type == academy.Classify) {
				for _, v := range t.Children {
					t.Count += v.Count
				}
			}
		}
	}

	return
}

func renderTag(v *academy.Tag) (tg *academy.TagMeta) {
	tg = &academy.TagMeta{
		ID:       v.ID,
		Type:     v.Type,
		State:    v.State,
		Business: v.Business,
		ParentID: v.ParentID,
		Name:     v.Name,
		Desc:     v.Desc,
		Rank:     v.Rank,
	}
	return
}

func fixTag(c *bm.Context) {
	var err error
	v := new(struct {
		ID       int64  `form:"id"`
		ParentID int64  `form:"parent_id"`
		Business int8   `form:"business"`
		State    int8   `form:"state"`
		Type     int8   `form:"type"`
		Name     string `form:"name"`
		Desc     string `form:"desc"`
		Rank     int64  `form:"rank"`
	})
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Model(&academy.Tag{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"parent_id": v.ParentID,
		"business":  v.Business,
		"type":      v.Type,
		"state":     v.State,
		"name":      v.Name,
		"desc":      v.Desc,
		"mtime":     time.Now().Format("2006-01-02 15:04:05"),
		"rank":      v.Rank,
	}).Error; err != nil {
		log.Error("creative-admin academy upTag error(%v)", err)
	}
	c.JSON(nil, err)
}

func orderTag(c *bm.Context) {
	var err error
	v := new(struct {
		ID         int64 `form:"id"  validate:"required"`
		Rank       int64 `form:"rank" validate:"required"`
		SwitchID   int64 `form:"switch_id"  validate:"required"`
		SwitchRank int64 `form:"switch_rank" validate:"required"`
	})
	if err = c.BindWith(v, binding.Form); err != nil {
		return
	}
	tx := svc.DB.Begin()
	tg := &academy.Tag{}
	if err = tx.Where("id=?", v.ID).First(&tg).Error; err != nil {
		log.Error("creative-admin academy orderTag error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	stg := &academy.Tag{}
	if err = tx.Where("id=?", v.SwitchID).First(&stg).Error; err != nil {
		log.Error("creative-admin academy orderTag error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	if err = tx.Table("academy_tag").Where("id=?", v.ID).Updates(
		map[string]interface{}{
			"rank":  v.SwitchRank,
			"mtime": time.Now().Format("2006-01-02 15:04:05")},
	).Error; err != nil {
		log.Error("creative-admin academy orderTag error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	if err = tx.Table("academy_tag").Where("id=?", v.SwitchID).Updates(
		map[string]interface{}{
			"rank":  v.Rank,
			"mtime": time.Now().Format("2006-01-02 15:04:05")},
	).Error; err != nil {
		log.Error("creative-admin academy orderTag error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新标签", TID: v.ID, OName: ""})
	c.JSON(nil, err)
}

func h5RelatedWebTIDsMap() (res map[int64][]int64, err error) {
	var tls []*academy.TagLink
	if err = svc.DB.Find(&tls).Error; err != nil {
		log.Error("academy h5TIDsMap error(%v)", err)
		return
	}
	if len(tls) == 0 {
		return
	}

	res = make(map[int64][]int64)
	for _, t := range tls {
		res[t.TID] = append(res[t.TID], t.LinkID)
	}
	return
}
