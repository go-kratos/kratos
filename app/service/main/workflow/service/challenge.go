package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	arc "go-common/app/service/main/archive/api"
	"go-common/app/service/main/workflow/model"
	"go-common/library/log"

	"github.com/jinzhu/gorm"
)

const (
	_inGroupSQL  = "INSERT INTO workflow_group(oid,business,rid,tid,lasttime,count,handling) VALUE(?,?,?,?,?,1,1) ON DUPLICATE KEY UPDATE lasttime=?,state=?,count=count+1,handling=handling+1"
	_inGroupSQL3 = "INSERT INTO workflow_group(oid,business,fid,rid,eid,score,tid,lasttime,count,handling) VALUES(?,?,?,?,?,?,?,?,1,1) ON DUPLICATE KEY UPDATE lasttime=?,state=?,count=count+1,handling=handling+1"
	//_inBusiness  = "INSERT INTO workflow_business(oid,typeid,business,title,content,mid,extra,cid,gid) VALUES(?,?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE typeid=values(typeid),title=values(title),content=values(content),extra=values(extra)"
	_inBusiness3 = "INSERT INTO workflow_business(oid,typeid,business,title,content,mid,extra,cid,gid) VALUES(?,?,?,?,?,?,?,?,?)"
	_upBusiness3 = "UPDATE workflow_business SET typeid=?, title=?, content=?, extra=? WHERE gid = ?"
)

// AddChallenge add challenge
func (s *Service) AddChallenge(c context.Context, ap *model.ChallengeParam) (row int64, err error) {
	var (
		lasttime = time.Now()
		apl      model.Group
		chall    *model.Challenge
		attach   *model.Attachment
	)
	tx := s.dao.DB.Begin()
	if tx.Error != nil {
		return
	}
	if tx.Error != nil {
		log.Error("s.dao.AddGroup(%d,%d,%d,%v) error(%v)", ap.Oid, ap.Business, ap.Tid, lasttime, err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("s.dao.AddGroup(%v,%v) rollback error(%v)", ap, lasttime, err1)
			}
		}
	}()
	var rid int64
	if tag, ok := s.tagsCache.TagMap3Tid[int64(ap.Business)][int64(ap.Tid)]; !ok {
		log.Error("not found tag bid(%d) tid(%d)", ap.Business, ap.Tid)
	} else {
		rid = tag.RID
	}
	if err = tx.Exec(_inGroupSQL, ap.Oid, ap.Business, rid, ap.Tid, lasttime, lasttime, model.StateTypePending).Error; err != nil {
		log.Error("tx.Raw(%d,%d,%d,%v) error(%v)", ap.Oid, ap.Business, ap.Tid, lasttime, err)
		return
	}
	if err = tx.Where("oid=? AND business=?", ap.Oid, ap.Business).Find(&apl).Error; err != nil {
		log.Error("tx.Where(%d,%d).Find() error(%v)", ap.Oid, ap.Business, err)
		return
	}
	chall = &model.Challenge{
		Tid:           ap.Tid,
		Gid:           apl.ID,
		Oid:           ap.Oid,
		Mid:           ap.Mid,
		Desc:          ap.Desc,
		MetaData:      ap.MetaData,
		Business:      ap.Business,
		BusinessState: ap.BusinessState,
		Adminid:       ap.AdminID,
		Assignee:      ap.AssigneeID,
	}
	chall.SetState(uint32(0), ap.Role)

	if err = tx.Create(chall).Error; err != nil {
		log.Error("tx.Create(%+v) error(%v)", chall, err)
		return
	}
	row = int64(chall.ID)
	for _, path := range ap.Attachments {
		attach = &model.Attachment{Cid: chall.ID, Path: path}
		if err = tx.Create(&attach).Error; err != nil {
			log.Error("tx.Create(%+v) error(%v)", attach, err)
			return
		}
	}
	// 刷新稿件信息
	if ap.Business == 1 {
		in := &arc.ArcRequest{
			Aid: ap.Oid,
		}
		var res *arc.ArcReply
		if res, err = s.arcClient.Arc(c, in); err != nil {
			log.Error("s.arcClient.Arc(%d) error(%v)", ap.Oid, err)
			return
		}
		ap.BusinessTypeid = res.Arc.TypeID
		ap.BusinessMid = res.Arc.Author.Mid
		ap.BusinessTitle = res.Arc.Title
		ap.BusinessContent = res.Arc.Desc
	}
	if ap.CheckBusiness() {
		bus := &model.Business{}
		if err = tx.Table("workflow_business").Where("gid = ?", apl.ID).Find(bus).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error("s.AddChallenge tx.Where error(%v)", err)
				return
			}
			if err = tx.Exec(_inBusiness3, ap.Oid, ap.BusinessTypeid, ap.Business, ap.BusinessTitle, ap.BusinessContent, ap.BusinessMid, ap.BusinessExtra, chall.ID, apl.ID).Error; err != nil {
				log.Error("tx.Exec(%s) params(oid=%d,typeid=%d,business=%d,title=%s,content=%s,mid=%d,extra=%s,cid=%d,gid=%d) error(%v)", _inBusiness3, ap.Oid, ap.BusinessTypeid, ap.Business, ap.BusinessTitle, ap.BusinessContent, ap.BusinessMid, ap.BusinessExtra, chall.ID, apl.ID, err)
				return
			}
		} else {
			if err = tx.Exec(_upBusiness3, ap.BusinessTypeid, ap.BusinessTitle, ap.BusinessContent, ap.BusinessExtra, apl.ID).Error; err != nil {
				log.Error("tx.Exec(%s) params(typeid=%d,title=%s,content=%s,extra=%s) error(%v)", _upBusiness3, ap.BusinessTypeid, ap.BusinessTitle, ap.BusinessContent, ap.BusinessExtra, err)
			}
		}

	}
	if err = tx.Commit().Error; err != nil {
		log.Error("tx.Commit() err(%v)", err)
		return
	}

	// upsert
	if apl.Count <= 1 {
		s.cache.Do(context.Background(), func(c context.Context) {
			if err = s.dao.UpdateUserTag(c, int64(apl.ID), ap.Tid); err != nil {
				log.Error("s.dao.UpdateUserTag gid:%d error:%v", apl.ID, err)
			}
		})
	}
	return
}

// UpChallengeState update challenge business state by mid && id && business
func (s *Service) UpChallengeState(c context.Context, id int32, mid int64, business, role, state int8) (err error) {
	if role != model.AuditRole && role != model.CustomerServiceRole {
		err = errors.New("UntreatedChallenge Unknown Role")
		log.Error("s.dao.UpChallengeState(%d,%d,%d,%d,%d) error(%v)", id, mid, business, role, state, err)
		return
	}
	chall := &model.Challenge{}
	if err = s.dao.DB.Where("id=? and mid=? and business=?", id, mid, business).Find(chall).Error; err != nil {
		log.Error("s.dao.UpChallengeState(%d,%d,%d,%d,%d) error(%v)", id, mid, business, role, state, err)
		return
	}
	chall.SetState(uint32(state), uint8(role))

	if err = s.dao.DB.Model(&model.Challenge{}).Where("id=? and mid=? and business=?", id, mid, business).Update("dispatch_state", chall.DispatchState).Error; err != nil {
		log.Error("s.dao.DB.UpState(%d,%d,%d,%d,%d) error(%v)", id, mid, business, role, state, err)
		return
	}

	return
}

// CloseChallenge set challenge business state closed by challenge id
func (s *Service) CloseChallenge(c context.Context, id int32, business, role, businessState int8, note string) (err error) {
	if role != model.AuditRole && role != model.CustomerServiceRole {
		err = errors.New("UntreatedChallenge Unknown Role")
		log.Error("s.dao.CloseChallenge(%d,%d) error(%v)", id, role, err)
		return
	}
	chall := &model.Challenge{}
	if err = s.dao.DB.Where("id=? and business=?", id, business).Find(chall).Error; err != nil {
		log.Error("s.dao.Challenge(%d,%d) error(%v)", id, business, err)
		return
	}

	chall.SetState(uint32(model.StateClose), uint8(role))
	if err = s.dao.DB.Model(&model.Challenge{}).Where("id=? and business=?", id, business).Update("dispatch_state", chall.DispatchState).Error; err != nil {
		log.Error("s.dao.UpChallengeState(%d) error(%v)", id, err)
		return
	}

	changeLog := &model.Log{AdminID: chall.Adminid, Oid: chall.Oid, Business: chall.Business, Target: chall.ID, Module: 1, Remark: "close challenge", Note: note}
	if err = s.dao.DB.Create(changeLog).Error; err != nil {
		log.Error("s.dao.CloseChallengeStateLog(%d) error(%v)", id, err)
	}
	return
}

// Challenge get challenge info
func (s *Service) Challenge(c context.Context, ap *model.ChallengeParam) (chall *model.Challenge, err error) {
	chall = &model.Challenge{}
	if err = s.dao.DB.Where("id=? and mid=? and business=?", ap.ID, ap.Mid, ap.Business).Find(chall).Error; err != nil {
		log.Error("s.dao.Challenge(%d,%d,%d) error(%v)", ap.ID, ap.Mid, ap.Business, err)
		return
	}
	// read new state field of challenge
	chall.FromState()
	if err = s.dao.DB.Where("cid=?", ap.ID).Order("id desc").Find(&chall.Attachments).Error; err != nil {
		log.Error("s.dao.Attachments(%d) error(%v)", ap.ID, err)
		return
	}
	if err = s.dao.DB.Where("cid=?", ap.ID).Order("id asc").Find(&chall.Events).Error; err != nil {
		log.Error("s.dao.Events(%d) error(%v)", ap.ID, err)
		return
	}
	return
}

// Challenges get challenge list
func (s *Service) Challenges(c context.Context, ap *model.ChallengeParam) (challs []*model.Challenge, err error) {
	if err = s.dao.DB.Where("mid=? and business=?", ap.Mid, ap.Business).Order("id DESC").Find(&challs).Error; err != nil {
		log.Error("s.dao.Challenges(%d,%d) error(%v)", ap.Mid, ap.Business, err)
		return
	}
	// read new state field of challenge
	for cid := range challs {
		challs[cid].FromState()
	}
	return
}

// UntreatedChallenge get untreated chanllenges by oid
func (s *Service) UntreatedChallenge(c context.Context, oid int64, role int8) (challs []*model.Challenge, err error) {
	if role != model.AuditRole && role != model.CustomerServiceRole {
		err = errors.New("UntreatedChallenge Unknown Role")
		log.Error("s.dao.UntreatedChallenge(%d,%d) error(%v)", oid, role, err)
		return
	}
	allChalls := []*model.Challenge{}
	if err = s.dao.DB.Where("oid=?", oid).Order("id DESC").Find(&allChalls).Error; err != nil {
		log.Error("s.dao.UntreatedChallenge(%d,%d) error(%v)", oid, role, err)
		return
	}
	for _, challenge := range allChalls {
		if value := challenge.GetState(uint8(role)); value == 0 {
			challs = append(challs, challenge)
		}
	}
	for _, chall := range challs {
		if err = s.dao.DB.Where("cid=?", chall.ID).Find(&chall.BusinessInfo).Error; err != nil {
			log.Error("s.dao.UntreatedChallenge(%d,%d) GetBusiness(%d) error(%v)", oid, role, chall.ID, err)
			continue
		}
	}
	err = nil
	return
}

// Callback callback audit event to manage
func (s *Service) Callback(c context.Context, chall *model.Challenge, businessID int8) (err error) {
	switch chall.Business {
	case model.BusinessAudit:
		err = s.dao.Callback(c, chall, businessID)
	default:
	}
	return
}

// CallbackByID callback audit event to manage
func (s *Service) CallbackByID(c context.Context, challengeID int32, businessID int8) (err error) {
	chall := &model.Challenge{}
	if err = s.dao.DB.Where("id=? ", challengeID).Find(chall).Error; err != nil {
		log.Error("s.CallbackByID(%d) error(%v)", challengeID, err)
		return
	}
	return s.Callback(c, chall, businessID)
}

// Challenges3 .
func (s *Service) Challenges3(c context.Context, cp3 *model.ChallengeParam3) (res []*model.Challenge3, err error) {
	if err = s.dao.DB.Where("business=? AND oid=? AND mid=?", cp3.Business, cp3.Oid, cp3.Mid).Order("id DESC").Find(&res).Error; err != nil {
		log.Error("s.Challenges3 error(%v)", err)
	}
	return
}

// AddChallenge3 .
func (s *Service) AddChallenge3(c context.Context, cp3 *model.ChallengeParam3) (row int64, err error) {
	var (
		lasttime  = time.Now()
		group     = &model.Group3{}
		challenge *model.Challenge
	)
	tx := s.dao.DB.Begin()
	if tx.Error != nil {
		log.Error("s.dao.DB.Begin error(%v)", err)
		return
	}
	defer func() {
		if err != nil {
			if err1 := tx.Rollback(); err1 != nil {
				log.Error("s.AddChallenge3 tx.Rollback error(%v)", err)
			}
		}
	}()
	if err = tx.Exec(_inGroupSQL3, cp3.Oid, cp3.Business, cp3.Fid, cp3.Rid, cp3.Eid, cp3.Score, cp3.Tid, lasttime, lasttime, model.StateTypePending).Error; err != nil {
		log.Error("s.AddChallenge3 tx.Exec error(%v)", err)
		return
	}
	if err = tx.Where("business=? AND oid=? AND eid=?", cp3.Business, cp3.Oid, cp3.Eid).Find(&group).Error; err != nil {
		log.Error("s.AddChallenge3 tx.Where error(%v)", err)
		return
	}
	challenge = &model.Challenge{
		Tid:           cp3.Tid,
		Gid:           int32(group.ID),
		Oid:           cp3.Oid,
		Mid:           cp3.Mid,
		Eid:           cp3.Eid,
		Desc:          cp3.Desc,
		MetaData:      cp3.MetaData,
		Business:      cp3.Business,
		BusinessState: cp3.BusinessState,
		Adminid:       cp3.AdminID,
		Assignee:      cp3.AssigneeID,
	}
	challenge.SetState(uint32(0), cp3.Role)
	if err = tx.Create(challenge).Error; err != nil {
		log.Error("s.AddChallenge3 tx.Create error(%v)", err)
		return
	}
	row = int64(challenge.ID)
	if len(cp3.Attachments) > 0 {
		values := []string{}
		valueArgs := []interface{}{}
		for _, a := range cp3.Attachments {
			values = append(values, "(?,?)")
			valueArgs = append(valueArgs, challenge.ID, a)
		}
		stmt := fmt.Sprintf("INSERT INTO workflow_attachment(cid,path) VALUES %s", strings.Join(values, ","))
		if err = tx.Exec(stmt, valueArgs...).Error; err != nil {
			return
		}
	}
	if cp3.CheckBusiness() {
		// todo workflow_business.oid save cp3.aid if exist
		var aid = cp3.Oid
		if cp3.Aid > 0 {
			aid = cp3.Aid
		}
		// todo insert or update business record
		bus := &model.Business{}
		if err = tx.Table("workflow_business").Where("gid = ?", group.ID).Find(bus).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				log.Error("s.AddChallenge3 tx.Where error(%v)", err)
				return
			}
			// insert business table
			if err = tx.Exec(_inBusiness3, aid, cp3.BusinessTypeid, cp3.Business, cp3.BusinessTitle, cp3.BusinessContent, cp3.BusinessMid, cp3.BusinessExtra, challenge.ID, group.ID).Error; err != nil {
				log.Error("tx.Exec(%s) params(oid=%d,typeid=%d,business=%d,title=%s,content=%s,mid=%d,extra=%s,cid=%d,gid=%d) error(%v)", _inBusiness3, cp3.Oid, cp3.BusinessTypeid, cp3.Business, cp3.BusinessTitle, cp3.BusinessContent, cp3.BusinessMid, cp3.BusinessExtra, challenge.ID, group.ID, err)
				return
			}
		} else {
			// update business table
			if err = tx.Exec(_upBusiness3, cp3.BusinessTypeid, cp3.BusinessTitle, cp3.BusinessContent, cp3.BusinessExtra, group.ID).Error; err != nil {
				log.Error("tx.Exec(%s) params(typeid=%d,title=%s,content=%s,extra=%s) error(%v)", _upBusiness3, cp3.BusinessTypeid, cp3.BusinessTitle, cp3.BusinessContent, cp3.BusinessExtra, err)
			}
		}
	}
	if err = tx.Commit().Error; err != nil {
		log.Error("s.AddChallenge3 tx.Commit error(%v)", err)
	}

	// upsert
	if group.Count <= 1 {
		s.cache.Do(context.Background(), func(c context.Context) {
			if err = s.dao.UpdateUserTag(c, group.ID, cp3.Tid); err != nil {
				log.Error("s.dao.UpdateUserTag gid:%d error:%v", group.ID, err)
			}
		})
	}
	return
}

// GroupState3 .
func (s *Service) GroupState3(c context.Context, cp3 *model.ChallengeParam3) (state int, err error) {
	err = s.dao.DB.Table("workflow_group").Select("state").Where("business=? AND oid=? AND eid=?", cp3.Business, cp3.Oid, cp3.Eid).Row().Scan(&state)
	if err == sql.ErrNoRows {
		err = nil
		state = -1
	}
	return
}
