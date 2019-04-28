package egg

import (
	"time"

	"go-common/app/admin/main/feed/conf"
	"go-common/app/admin/main/feed/dao/egg"
	"go-common/app/admin/main/feed/model/common"
	eggModel "go-common/app/admin/main/feed/model/egg"
	Log "go-common/app/admin/main/feed/util"
	"go-common/library/log"
	libTime "go-common/library/time"
)

const (
	//_ActionAdd log action add
	_ActionAdd = "add"
	//_ActionUpdate log action update
	_ActionUpdate = "update"
	//_ActionDel log action delete
	_ActionDel = "delete"
	//_ActionPub log action publish
	_ActionPub = "publish"
)

// Service is egg service
type Service struct {
	dao *egg.Dao
}

// New new a egg service
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: egg.New(c),
	}
	return
}

//EggWithID get egg with ID
func (s *Service) EggWithID(id uint) (egg *eggModel.Egg, err error) {
	egg = &eggModel.Egg{}
	if err = s.dao.DB.Model(&eggModel.Egg{}).Where("id = ?", id).
		Where("`delete` != ?", eggModel.Delete).First(egg).
		Error; err != nil {
		log.Error("eggSrv.DelEgg Update error(%v)", err)
		return
	}
	return
}

//DelEgg update egg
func (s *Service) DelEgg(id uint, person string, uid int64) (err error) {
	tx := s.dao.DB.Begin()
	if err = tx.Model(&eggModel.Egg{}).Where("id = ?", id).Update("delete", eggModel.Delete).Error; err != nil {
		log.Error("eggSrv.DelEgg Update error(%v)", err)
		return
	}
	if err = tx.Model(&eggModel.Query{}).Where("egg_id = ?", id).Update("deleted", eggModel.Delete).Error; err != nil {
		log.Error("eggSrv.DelEgg UpdateQuery error(%v)", err)
		return
	}
	tx.Commit()
	if err = Log.AddLog(eggModel.Business, person, uid, int64(id), _ActionDel, ""); err != nil {
		log.Error("eggSrv.AddEgg AddLog error(%v)", err)
		return
	}
	return
}

//PubEgg publish egg
func (s *Service) PubEgg(id uint, publish uint8, person string, uid int64) (err error) {
	if err = s.dao.DB.Model(&eggModel.Egg{}).Where("id = ?", id).Update("publish", publish).Error; err != nil {
		log.Error("eggSrv.PubEgg Update error(%v)", err)
		return
	}
	if err = Log.AddLog(eggModel.Business, person, uid, int64(id), _ActionPub, ""); err != nil {
		log.Error("eggSrv.PubEgg AddLog error(%v)", err)
		return
	}
	return
}

//UpdateEgg update egg
func (s *Service) UpdateEgg(e *eggModel.Egg, p []eggModel.Plat, w []string) (err error) {
	tx := s.dao.DB.Begin()
	db := tx.Model(&eggModel.Egg{}).Where("id = ?", e.ID).Update(e)
	if err = db.Error; err != nil {
		log.Error("eggSrv.UpdateEgg UpdateEgg egg error(%v)", err)
		return
	}
	db = tx.Model(&eggModel.Plat{}).Where("egg_id = ?", e.ID).Update("deleted", eggModel.Delete)
	if err = db.Error; err != nil {
		log.Error("eggSrv.UpdateEgg UpdatePlat error(%v)", err)
		return
	}
	db = tx.Model(&eggModel.Query{}).Where("egg_id = ?", e.ID).Update("deleted", eggModel.Delete)
	if err = db.Error; err != nil {
		log.Error("eggSrv.UpdateEgg UpdateQuery error(%v)", err)
		return
	}
	for _, v := range p {
		v.EggID = e.ID
		v.Deleted = eggModel.NotDelete
		db := tx.Model(&eggModel.Plat{}).Create(v)
		if err = db.Error; err != nil {
			log.Error("eggSrv.UpdateEgg create egg plat error(%v)", err)
			return
		}
	}
	for _, v := range w {
		q := eggModel.Query{
			Word:    v,
			EggID:   e.ID,
			STime:   e.Stime,
			ETime:   e.Etime,
			Deleted: eggModel.NotDelete,
		}
		db := tx.Model(&eggModel.Query{}).Create(q)
		if err = db.Error; err != nil {
			log.Error("eggSrv.UpdateEgg create egg query error(%v)", err)
			return
		}
	}
	tx.Commit()
	obj := map[string]interface{}{
		"egg":   e,
		"plat":  p,
		"words": w,
	}
	if err = Log.AddLog(eggModel.Business, e.Person, e.UID, int64(e.ID), _ActionUpdate, obj); err != nil {
		log.Error("eggSrv.AddEgg AddLog error(%v)", err)
		return
	}
	return
}

//AddEgg add egg
func (s *Service) AddEgg(e *eggModel.Egg, p []eggModel.Plat, w []string) (err error) {
	tx := s.dao.DB.Begin()
	db := tx.Model(&eggModel.Egg{}).Create(e)
	if err = db.Error; err != nil {
		log.Error("eggSrv.AddEgg create egg error(%v)", err)
		return
	}
	for _, v := range p {
		v.EggID = e.ID
		v.Deleted = eggModel.NotDelete
		db := tx.Model(&eggModel.Plat{}).Create(v)
		if err = db.Error; err != nil {
			log.Error("eggSrv.AddEgg create egg plat error(%v)", err)
			return
		}
	}
	for _, v := range w {
		q := eggModel.Query{
			Word:    v,
			EggID:   e.ID,
			STime:   e.Stime,
			ETime:   e.Etime,
			Deleted: eggModel.NotDelete,
		}
		db := tx.Model(&eggModel.Query{}).Create(q)
		if err = db.Error; err != nil {
			log.Error("eggSrv.AddEgg create egg query error(%v)", err)
			return
		}
	}
	tx.Commit()
	obj := map[string]interface{}{
		"egg":   e,
		"plat":  p,
		"words": w,
	}
	if err = Log.AddLog(eggModel.Business, e.Person, e.UID, int64(e.ID), _ActionAdd, obj); err != nil {
		log.Error("eggSrv.AddEgg AddLog error(%v)", err)
		return
	}
	return
}

//IsWdExist the word will add is exist
func (s *Service) IsWdExist(words []string, sTime, eTime libTime.Time, eggID uint) (exist bool, w string, err error) {
	var (
		c int
	)
	for _, v := range words {
		query := s.dao.DB.Model(&eggModel.Query{}).
			Where("deleted=?", eggModel.NotDelete).
			Where("word = ?", v).
			Where("s_time < ?", eTime).
			Where("e_time > ?", sTime)
		if eggID != 0 {
			query = query.Where("egg_id != ?", eggID)
		}
		if err = query.Count(&c).Error; err != nil {
			log.Error("eggSrv.IsWdExist Query error(%v)", err)
			return
		}
		if c > 0 {
			return true, v, nil
		}
	}
	return false, "", nil
}

func (s *Service) qWord(word string) (ids []uint, err error) {
	q := []eggModel.Query{}
	if err = s.dao.DB.Model(&eggModel.Query{}).Where("deleted=?", eggModel.NotDelete).
		Where("word like ?", "%"+word+"%").Find(&q).Error; err != nil {
		log.Error("eggSrv.IndexEgg Query error(%v)", err)
		return
	}
	for _, v := range q {
		ids = append(ids, v.EggID)
	}
	return
}

//IndexEgg egg list
func (s *Service) IndexEgg(param *eggModel.IndexParam) (values *eggModel.IndexPager, err error) {
	values = &eggModel.IndexPager{
		Page: common.Page{
			Num:  param.Pn,
			Size: param.Ps,
		},
	}
	w := map[string]interface{}{
		"delete": eggModel.NotDelete,
	}
	if param.ID != "" {
		w["id"] = param.ID
	}
	query := s.dao.DB.Model(&eggModel.Index{}).Where(w)
	if param.Stime != "" {
		query = query.Where("stime >= ?", param.Stime)
	}
	if param.Etime != "" {
		query = query.Where("etime <= ?", param.Etime)
	}
	if param.Person != "" {
		query = query.Where("person like ?", "%"+param.Person+"%")
	}
	if param.Word != "" {
		var ids = []uint{}
		if ids, err = s.qWord(param.Word); err != nil {
			return
		}
		if len(ids) != 0 {
			query = query.Where("id in (?)", ids)
		} else {
			query = query.Where("id in (?)", 0)
		}
	}
	if err = query.Order("`id` DESC").Offset((param.Pn - 1) * param.Ps).Limit(param.Ps).Find(&values.Item).Error; err != nil {
		log.Error("eggSrv.IndexEgg Index list error(%v)", err)
		return
	}
	if err = query.Count(&values.Page.Total).Error; err != nil {
		log.Error("eggSrv.IndexEgg Index count error(%v)", err)
		return
	}
	for k, v := range values.Item {
		q := []eggModel.Query{}
		p := []eggModel.Plat{}
		//select egg query words
		w := map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Query{}).Where(w).Find(&q).Error; err != nil {
			log.Error("eggSrv.IndexEgg Query error(%v)", err)
			return
		}
		for _, qV := range q {
			if values.Item[k].Words == "" {
				values.Item[k].Words = qV.Word
			} else {
				values.Item[k].Words = values.Item[k].Words + "," + qV.Word
			}
		}
		//select egg plat
		w = map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Plat{}).Where(w).Find(&p).Error; err != nil {
			log.Error("eggSrv.IndexEgg Plat error(%v)", err)
			return
		}
		values.Item[k].Plat = p
	}
	return
}

//SearchEgg search egg list
func (s *Service) SearchEgg() (values []eggModel.SearchEgg, err error) {
	param := eggModel.IndexParam{}
	w := map[string]interface{}{
		"delete":  eggModel.NotDelete,
		"publish": eggModel.Publish,
	}
	cTime := time.Now().Unix()
	cTime = cTime + 10*60
	tm := time.Unix(cTime, 0)
	param.Stime = tm.Format("2006-01-02 15:04:05")
	param.Etime = param.Stime
	query := s.dao.DB.Model(&eggModel.SearchEgg{}).Where(w)
	query = query.Where("stime <= ?", param.Stime).Where("etime >= ?", param.Etime)
	if err = query.Order("`id` DESC").Find(&values).Error; err != nil {
		log.Error("eggSrv.SearchEgg Index list error(%v)", err)
		return
	}
	for k, v := range values {
		q := []eggModel.Query{}
		p := []eggModel.Plat{}
		Words := []string{}
		//select egg query words
		w := map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Query{}).Where(w).Find(&q).Error; err != nil {
			log.Error("eggSrv.IndexEgg Query error(%v)", err)
			return
		}
		for _, v := range q {
			Words = append(Words, v.Word)
		}
		values[k].Words = Words
		//select egg plat
		w = map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Plat{}).Where(w).Find(&p).Error; err != nil {
			log.Error("eggSrv.IndexEgg Plat error(%v)", err)
			return
		}
		values[k].Plat = make(map[uint8]eggModel.Plat)
		for _, v := range p {
			values[k].Plat[v.Plat] = v
		}
	}
	return
}

//SearchEggWeb search egg list
func (s *Service) SearchEggWeb() (values []eggModel.SearchEggWeb, err error) {
	param := eggModel.IndexParam{}
	w := map[string]interface{}{
		"delete":  eggModel.NotDelete,
		"publish": eggModel.Publish,
	}
	cTime := time.Now().Unix()
	cTime = cTime + 10*60
	tm := time.Unix(cTime, 0)
	param.Stime = tm.Format("2006-01-02 15:04:05")
	param.Etime = param.Stime
	query := s.dao.DB.Model(&eggModel.SearchEgg{}).Where(w)
	query = query.Where("stime <= ?", param.Stime).Where("etime >= ?", param.Etime)
	if err = query.Order("`id` DESC").Find(&values).Error; err != nil {
		log.Error("eggSrv.SearchEgg Index list error(%v)", err)
		return
	}
	for k, v := range values {
		q := []eggModel.Query{}
		p := []eggModel.Plat{}
		Words := []string{}
		//select egg query words
		w := map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Query{}).Where(w).Find(&q).Error; err != nil {
			log.Error("eggSrv.IndexEgg Query error(%v)", err)
			return
		}
		for _, v := range q {
			Words = append(Words, v.Word)
		}
		values[k].Words = Words
		//select egg plat
		w = map[string]interface{}{
			"egg_id":  v.ID,
			"deleted": eggModel.NotDelete,
		}
		if err = s.dao.DB.Model(&eggModel.Plat{}).Where(w).Find(&p).Error; err != nil {
			log.Error("eggSrv.IndexEgg Plat error(%v)", err)
			return
		}
		values[k].Plat = make(map[uint8][]eggModel.Plat)
		for _, v := range p {
			values[k].Plat[v.Plat] = append(values[k].Plat[v.Plat], v)
		}
	}
	return
}
