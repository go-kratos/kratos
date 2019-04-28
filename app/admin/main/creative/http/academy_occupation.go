package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"go-common/app/admin/main/creative/model/academy"
	"go-common/library/ecode"
	"go-common/library/log"
	bm "go-common/library/net/http/blademaster"
	"go-common/library/net/http/blademaster/binding"
	"go-common/library/net/http/blademaster/render"
	"go-common/library/sync/errgroup"
	"go-common/library/xstr"

	"github.com/jinzhu/gorm"
)

func addOccupation(c *bm.Context) {
	var (
		err error
		now = time.Now().Format("2006-01-02 15:04:05")
	)
	v := new(academy.Occupation)
	if err = c.Bind(v); err != nil {
		log.Error("addOccupation c.Bind error(%v)", err)
		return
	}
	m := &academy.Occupation{
		State:        academy.StateNormal,
		Name:         v.Name,
		Desc:         v.Desc,
		Logo:         v.Logo,
		MainStep:     v.MainStep,
		MainSoftware: v.MainSoftware,
		CTime:        now,
		MTime:        now,
	}
	tx := svc.DB.Begin()
	if err = tx.Create(m).Error; err != nil {
		log.Error("academy addOccupation error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&academy.Occupation{}).Where("id=?", m.ID).Updates(map[string]interface{}{
		"rank": m.ID,
	}).Error; err != nil {
		log.Error("academy addOccupation error(%v)", err)
		tx.Rollback()
	}
	tx.Commit()
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "添加职业", TID: m.ID, OName: m.Name})
	c.JSON(map[string]interface{}{
		"id": m.ID,
	}, nil)
}

func upOccupation(c *bm.Context) {
	var (
		oc  = &academy.Occupation{}
		err error
	)
	v := new(academy.Occupation)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(oc).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if oc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Occupation{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"name":          v.Name,
		"desc":          v.Desc,
		"logo":          v.Logo,
		"main_step":     v.MainStep,
		"main_software": v.MainSoftware,
	}).Error; err != nil {
		log.Error("academy bindOccupation error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新职业", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func bindOccupation(c *bm.Context) {
	var (
		oc  = &academy.Occupation{}
		err error
	)
	v := new(academy.Occupation)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(oc).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if oc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Occupation{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"state": v.State,
	}).Error; err != nil {
		log.Error("academy upOccupation error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新职业", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func addSkill(c *bm.Context) {
	var (
		err error
		now = time.Now().Format("2006-01-02 15:04:05")
	)
	v := new(academy.Skill)
	if err = c.Bind(v); err != nil {
		log.Error("addSkill c.Bind error(%v)", err)
		return
	}
	m := &academy.Skill{
		State: academy.StateNormal,
		OID:   v.OID,
		Name:  v.Name,
		Desc:  v.Desc,
		CTime: now,
		MTime: now,
	}
	if err = svc.DB.Create(m).Error; err != nil {
		log.Error("academy addSkill error(%v)", err)
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "添加步骤", TID: m.ID, OName: m.Name})
	c.JSON(map[string]interface{}{
		"id": m.ID,
	}, nil)
}

func upSkill(c *bm.Context) {
	var (
		sk  = &academy.Skill{}
		err error
	)
	v := new(academy.Skill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(sk).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if sk == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Skill{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"name": v.Name,
		"desc": v.Desc,
	}).Error; err != nil {
		log.Error("academy upSkill error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新步骤", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func bindSkill(c *bm.Context) {
	var (
		oc  = &academy.Skill{}
		err error
	)
	v := new(academy.Skill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(oc).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if oc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Skill{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"state": v.State,
	}).Error; err != nil {
		log.Error("academy bindSkill error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新职业", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func addSoftware(c *bm.Context) {
	var (
		err error
		now = time.Now().Format("2006-01-02 15:04:05")
	)
	v := new(academy.Software)
	if err = c.Bind(v); err != nil {
		return
	}
	m := &academy.Software{
		State: academy.StateNormal,
		SkID:  v.SkID,
		Name:  v.Name,
		Desc:  v.Desc,
		CTime: now,
		MTime: now,
	}
	if err = svc.DB.Create(m).Error; err != nil {
		log.Error("academy addSoftware error(%v)", err)
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "添加步骤", TID: m.ID, OName: m.Name})
	c.JSON(map[string]interface{}{
		"id": m.ID,
	}, nil)
}

func upSoftware(c *bm.Context) {
	var (
		oc  = &academy.Software{}
		err error
	)
	v := new(academy.Software)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(oc).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if oc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Software{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"name": v.Name,
		"desc": v.Desc,
	}).Error; err != nil {
		log.Error("academy upSoftware error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新步骤", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func bindSoftware(c *bm.Context) {
	var (
		oc  = &academy.Software{}
		err error
	)
	v := new(academy.Software)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	if err = svc.DB.Where("id=?", v.ID).Find(oc).Error; err != nil {
		c.JSON(nil, err)
		return
	}
	if oc == nil {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.Software{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"state": v.State,
	}).Error; err != nil {
		log.Error("academy bindSoftware error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新职业", TID: v.ID, OName: v.Name})
	c.JSON(nil, err)
}

func listOccupation(c *bm.Context) {
	var (
		ocs []*academy.Occupation
		sks []*academy.Skill
		sfs []*academy.Software

		sfMap map[int64][]*academy.Software
		skMap map[int64][]*academy.Skill

		sfIDs, skIDs, oIDs        []int64
		softMap, skillMap, occMap map[int64]int

		g, _ = errgroup.WithContext(c)
		err  error
	)

	g.Go(func() error { //软件
		if err = svc.DB.Find(&sfs).Error; err != nil {
			log.Error("listSoftware error(%v)", err)
			return err
		}

		if len(sfs) == 0 {
			return nil
		}

		sfMap = make(map[int64][]*academy.Software)
		sfIDs = make([]int64, 0, len(sfs))
		for _, v := range sfs {
			sfMap[v.SkID] = append(sfMap[v.SkID], v) //按技能聚合软件
			sfIDs = append(sfIDs, v.ID)
		}
		softMap, _ = arcCountBySids(sfIDs) //获取软件对应稿件数量

		if len(softMap) == 0 {
			return nil
		}
		for _, v := range sfs {
			if n, ok := softMap[v.ID]; ok { //映射软件对应稿件数量
				v.Count = n
			}
		}

		return nil
	})

	g.Go(func() error { //技能
		if err = svc.DB.Find(&sks).Error; err != nil {
			log.Error("listSkill error(%v)", err)
			return err
		}

		if len(sks) == 0 {
			return nil
		}

		skMap = make(map[int64][]*academy.Skill)
		skIDs = make([]int64, 0, len(sks))
		for _, v := range sks {
			skMap[v.OID] = append(skMap[v.OID], v) //按职业聚合技能
			skIDs = append(skIDs, v.ID)
		}
		skillMap, _ = arcCountBySkids(skIDs) //获取技能对应稿件数量
		return nil
	})

	g.Go(func() error { //职业
		if err = svc.DB.Order("rank ASC").Find(&ocs).Error; err != nil {
			log.Error("listOccupation error(%v)", err)
			return err
		}

		if len(ocs) == 0 {
			return nil
		}

		oIDs = make([]int64, 0, len(ocs))
		for _, v := range ocs {
			oIDs = append(oIDs, v.ID)
		}
		occMap, _ = arcCountByPids(oIDs) //获取职业对应稿件数量

		return nil
	})

	if err = g.Wait(); err != nil {
		c.JSON(nil, err)
		return
	}

	for _, v := range sks {
		if sf, ok := sfMap[v.ID]; ok { //添加软件节点
			v.Software = sf
		}
		if n, ok := skillMap[v.ID]; ok { //映射技能对应稿件数量
			v.Count = n
		}
	}

	for _, v := range ocs {
		if sk, ok := skMap[v.ID]; ok { //添加技能节点
			v.Skill = sk
		}
		if n, ok := occMap[v.ID]; ok { //映射职业对应稿件数量
			v.Count = n
		}
	}

	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    ocs,
	}))
}

func arcCountByPids(tids []int64) (res map[int64]int, err error) {
	var (
		countSQL = "SELECT pid AS tid, count(DISTINCT aid) AS count  FROM academy_arc_skill WHERE state=0 AND pid IN (?) GROUP BY pid"
		ats      []*academy.ArchiveCount
	)
	if err = svc.DB.Raw(countSQL, tids).Find(&ats).Error; err != nil {
		log.Error("occupation arcCountByPids error(%v)", err)
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

func arcCountBySkids(tids []int64) (res map[int64]int, err error) {
	var (
		countSQL = "SELECT skid AS tid, count(DISTINCT aid) AS count  FROM academy_arc_skill WHERE state=0 AND skid IN (?) GROUP BY skid"
		ats      []*academy.ArchiveCount
	)
	if err = svc.DB.Raw(countSQL, tids).Find(&ats).Error; err != nil {
		log.Error("skill arcCountByPids error(%v)", err)
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

func arcCountBySids(tids []int64) (res map[int64]int, err error) {
	var (
		countSQL = "SELECT sid AS tid, count(DISTINCT aid) AS count  FROM academy_arc_skill WHERE state=0 AND sid IN (?) GROUP BY sid"
		ats      []*academy.ArchiveCount
	)
	if err = svc.DB.Raw(countSQL, tids).Find(&ats).Error; err != nil {
		log.Error("software arcCountBySids error(%v)", err)
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

func orderOccupation(c *bm.Context) {
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
	var oc academy.Occupation
	if err = tx.Where("id=?", v.ID).First(&oc).Error; err != nil {
		log.Error("academy orderOccupation error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	var soc academy.Occupation
	if err = tx.Where("id=?", v.SwitchID).First(&soc).Error; err != nil {
		log.Error("academy orderOccupation error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}

	if err = tx.Model(&academy.Occupation{}).Where("id=?", v.ID).Updates(
		map[string]interface{}{
			"rank": v.SwitchRank,
		},
	).Error; err != nil {
		log.Error("academy orderOccupation error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	if err = tx.Model(&academy.Occupation{}).Where("id=?", v.SwitchID).Updates(
		map[string]interface{}{
			"rank": v.Rank,
		},
	).Error; err != nil {
		log.Error("academy orderOccupation error(%v)", err)
		c.JSON(nil, err)
		tx.Rollback()
		return
	}
	tx.Commit()
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新职业排序", TID: v.ID, OName: ""})
	c.JSON(nil, err)
}

func checkSkillArcExist(aid, sid int64) (res *academy.ArcSkill, err error) {
	var as academy.ArcSkill
	err = svc.DB.Where("aid=?", aid).Where("sid=?", sid).Find(&as).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		}
		log.Error("academy checkSkillArcExist aid(%d)|sid(%d)|error(%v)", aid, sid, err)
	}
	res = &as
	return
}

func viewSkillArc(c *bm.Context) {
	var (
		err error
		as  = &academy.ArcSkill{}
	)
	v := new(academy.ArcSkill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.AID == 0 {
		c.JSON(nil, ecode.RequestErr)
		return
	}
	arc, err := svc.Archive(c, v.AID)
	if err != nil {
		c.JSON(nil, err)
		return
	}

	if arc != nil {
		as.AID = v.AID
		as.Title = arc.Title
		as.Pic = arc.Pic
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    as,
	}))
}

func addSkillArc(c *bm.Context) {
	var (
		err error
		now = time.Now().Format("2006-01-02 15:04:05")
	)
	v := new(academy.ArcSkill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.AID == 0 || v.PID == 0 || v.SkID == 0 || v.SID == 0 || v.Type == 0 {
		log.Error("academy addSkillArc v(%+v)", v)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	res, err := checkSkillArcExist(v.AID, v.SID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res != nil && res.ID != 0 {
		c.JSON(nil, ecode.CreativeAcademyDuplicateSoftIDErr)
		return
	}

	m := &academy.ArcSkill{
		State: academy.StateNormal,
		AID:   v.AID,
		Type:  v.Type,
		PID:   v.PID,
		SkID:  v.SkID,
		SID:   v.SID,
		CTime: now,
		MTime: now,
	}
	if err = svc.DB.Create(m).Error; err != nil {
		log.Error("academy addSoftware error(%v)", err)
		c.JSON(nil, err)
		return
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "添加技能稿件", OIDs: xstr.JoinInts([]int64{v.AID})})
	c.JSON(nil, err)
}

func upSkillArc(c *bm.Context) {
	var err error
	v := new(academy.ArcSkill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.ID == 0 || v.AID == 0 || v.PID == 0 || v.SkID == 0 || v.SID == 0 || v.Type == 0 {
		log.Error("academy upSkillArc v(%+v)", v)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if err = svc.DB.Model(&academy.ArcSkill{}).Where("id=?", v.ID).Updates(map[string]interface{}{
		"aid":  v.AID,
		"type": v.Type,
		"pid":  v.PID,
		"skid": v.SkID,
		"sid":  v.SID,
	}).Error; err != nil {
		log.Error("academy upSkillArc error(%v)", err)
	}

	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新技能稿件", OIDs: xstr.JoinInts([]int64{v.AID})})
	c.JSON(nil, err)
}

func bindSkillArc(c *bm.Context) {
	var err error
	v := new(academy.ArcSkill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.AID == 0 || v.SID == 0 {
		log.Error("academy bindSkillArc v(%+v)", v)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	res, err := checkSkillArcExist(v.AID, v.SID)
	if err != nil {
		c.JSON(nil, err)
		return
	}
	if res != nil && res.ID == 0 {
		c.JSON(nil, ecode.NothingFound)
		return
	}
	if err = svc.DB.Model(&academy.ArcSkill{}).Where("aid=?", v.AID).Where("sid=?", v.SID).Updates(map[string]interface{}{
		"state": v.State,
	}).Error; err != nil {
		log.Error("academy bindSkillArc error(%v)", err)
	}
	uid, uname := getUIDName(c)
	svc.SendAcademyLog(c, &academy.LogParam{UID: uid, UName: uname, Action: "更新技能稿件状态", TID: v.ID, OName: v.Title})
	c.JSON(nil, err)
}

func skillArcList(c *bm.Context) {
	var (
		items []*academy.ArcSkill
		aids  []int64
		total int
		err   error
	)
	v := new(academy.ArcSkill)
	if err = c.Bind(v); err != nil {
		return
	}
	if v.Pn < 1 {
		v.Pn = 1
	}
	if v.Ps > 20 {
		v.Ps = 20
	}
	db := svc.DB.Model(&academy.ArcSkill{})
	if v.PID != 0 {
		db = db.Where("pid=?", v.PID)
	}
	if v.SkID != 0 {
		db = db.Where("skid=?", v.SkID)
	}
	db.Count(&total)
	if err = db.Order("ctime DESC").Offset((v.Pn - 1) * v.Ps).Limit(v.Ps).Find(&items).Error; err != nil {
		log.Error("academy skillArcList error(%v)", err)
		c.JSON(nil, err)
		return
	}
	for _, v := range items {
		aids = append(aids, v.AID)
	}
	arcs, err := svc.Archives(c, aids)
	if err != nil {
		log.Error("academy skillArcList error(%v)", err)
	}
	for _, v := range items {
		if a, ok := arcs[v.AID]; ok {
			v.Pic = a.Pic
			v.Title = a.Title
		}
	}
	data := &academy.ArcSkills{
		Items: items,
		Pager: &academy.Pager{
			Num:   v.Pn,
			Size:  v.Ps,
			Total: total,
		},
	}
	c.Render(http.StatusOK, render.MapJSON(map[string]interface{}{
		"code":    0,
		"message": "0",
		"data":    data,
	}))
}

func subSearchKeywords(c *bm.Context) {
	var err error
	bs, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Error("subSearchKeywords ioutil.ReadAll error(%v)", err)
		c.JSON(nil, ecode.RequestErr)
		return
	}
	c.Request.Body.Close()

	var v []*academy.SearchKeywords
	err = json.Unmarshal(bs, &v)
	if err != nil {
		log.Error("subSearchKeywords json.Unmarshal v(%+v) error(%v)", v, err)
		c.JSON(nil, ecode.RequestErr)
		return
	}

	if err = svc.SubSearchKeywords(v); err != nil {
		c.JSON(nil, err)
		return
	}

	c.JSON(nil, nil)
}

func searchKeywords(c *bm.Context) {
	res, err := svc.SearchKeywords()
	if err != nil {
		c.JSON(nil, err)
		return
	}
	c.JSON(res, nil)
}
