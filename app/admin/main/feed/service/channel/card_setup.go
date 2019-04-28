package channel

import (
	"encoding/json"
	"fmt"
	"strconv"

	"go-common/app/admin/main/feed/conf"
	accdao "go-common/app/admin/main/feed/dao/account"
	arcdao "go-common/app/admin/main/feed/dao/archive"
	pgcdao "go-common/app/admin/main/feed/dao/pgc"
	showdao "go-common/app/admin/main/feed/dao/show"
	cardmodel "go-common/app/admin/main/feed/model/channel"
	"go-common/app/admin/main/feed/model/common"
	"go-common/app/admin/main/feed/util"
	"go-common/library/log"
)

// Service is search service
type Service struct {
	showDao *showdao.Dao
	pgcDao  *pgcdao.Dao
	accDao  *accdao.Dao
	arcDao  *arcdao.Dao
}

// New new a search service
func New(c *conf.Config) (s *Service) {
	var (
		pgc *pgcdao.Dao
		err error
	)
	if pgc, err = pgcdao.New(c); err != nil {
		log.Error("pgcdao.New error(%v)", err)
		return
	}
	s = &Service{
		showDao: showdao.New(c),
		pgcDao:  pgc,
		accDao:  accdao.New(c),
		arcDao:  arcdao.New(c),
	}
	return
}

//parseConten parse string type id to int type id
func parseConten(content string) (s string, err error) {
	type Content struct {
		ID    string `json:"id"`
		Title string `json:"title"`
	}
	var contents []Content
	if err = json.Unmarshal([]byte(content), &contents); err != nil {
		return
	}
	type ContentTrans struct {
		ID    int64  `json:"id"`
		Title string `json:"title"`
	}
	var cTrans []ContentTrans
	for _, v := range contents {
		var s int64
		if s, err = strconv.ParseInt(v.ID, 10, 64); err != nil {
			return "", err
		}
		v := ContentTrans{
			ID:    s,
			Title: v.Title,
		}
		cTrans = append(cTrans, v)
	}
	var b []byte
	if b, err = json.Marshal(cTrans); err != nil {
		return "", err
	}
	return string(b), nil
}

//AddCardSetup card set up
func (s *Service) AddCardSetup(card *cardmodel.AddCardSetup, person string, uid int64) (err error) {
	var (
		flag bool
		e    error
	)
	flag, e = isDup(card.Content)
	if e != nil {
		return e
	}
	if flag {
		return fmt.Errorf("ID重复")
	}
	card.Person = person
	card.UID = uid
	if card.Content, err = parseConten(card.Content); err != nil {
		return
	}
	if err = s.showDao.DB.Model(&cardmodel.AddCardSetup{}).
		Create(card).Error; err != nil {
		log.Error("chanelSvc.AddCardSetup Create error(%v)", err)
		return
	}
	if card.Type == common.CardPgcsRcmd {
		if err = util.AddLog(cardmodel.LogBusPgcsRcmd, person, uid, 0, cardmodel.ActAddCsPgcRcmd, card); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	} else if card.Type == common.CardUpRcmdNew {
		if err = util.AddLog(cardmodel.LogBusRcmdNew, person, uid, 0, cardmodel.ActAddCsRcmdNew, card); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	}
	return
}

//CardSetupList card set up
func (s *Service) CardSetupList(id int, t string, person string, title string, pn int, ps int) (cPager *cardmodel.SetupPager, err error) {
	cPager = &cardmodel.SetupPager{
		Page: common.Page{
			Num:  pn,
			Size: ps,
		},
	}
	w := map[string]interface{}{
		"deleted": cardmodel.NotDelete,
		"type":    t,
	}
	query := s.showDao.DB.Model(&cardmodel.Setup{})
	if id != 0 {
		w["id"] = id
	}
	if person != "" {
		query = query.Where("person like ?", "%"+person+"%")
	}
	if title != "" {
		if t == "up_rcmd_new" {
			query = query.Where("long_title like ?", "%"+title+"%")
		} else {
			query = query.Where("title like ?", "%"+title+"%")
		}
	}
	if err = query.Where(w).Count(&cPager.Page.Total).Error; err != nil {
		log.Error("chanelSvc.CardSetupList Index count error(%v)", err)
		return
	}
	cards := []*cardmodel.Setup{}
	if err = query.Where(w).Order("`id` DESC").Offset((pn - 1) * ps).Limit(ps).Find(&cards).Error; err != nil {
		log.Error("chanelSvc.CardSetupList First error(%v)", err)
		return
	}
	cPager.Item = cards
	return
}

//DelCardSetup card set up
func (s *Service) DelCardSetup(id int, t string, person string, uid int64) (err error) {
	dbModel := s.showDao.DB.Model(&cardmodel.Setup{})
	dbModel = dbModel.Where("id = ?", id).Where("type = ?", t)
	if err = dbModel.Update("deleted", cardmodel.Delete).Error; err != nil {
		log.Error("chanelSvc.CardSetupList First error(%v)", err)
		return
	}
	if t == common.CardPgcsRcmd {
		if err = util.AddLog(cardmodel.LogBusPgcsRcmd, person, uid, int64(id), cardmodel.ActDelCsPgcRcmd, ""); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	} else if t == common.CardUpRcmdNew {
		if err = util.AddLog(cardmodel.LogBusRcmdNew, person, uid, int64(id), cardmodel.ActDelCsRcmdNew, ""); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	}
	return
}

func isDup(con string) (flag bool, err error) {
	type Content struct {
		ID string `json:"id"`
	}
	value := []Content{}
	if err := json.Unmarshal([]byte(con), &value); err != nil {
		return false, err
	}
	s := make(map[string]bool)
	for _, v := range value {
		if s[v.ID] {
			return true, nil
		}
		s[v.ID] = true
	}
	return false, nil
}

//UpdateCardSetup card set up
func (s *Service) UpdateCardSetup(id int, card *cardmodel.AddCardSetup, person string, uid int64) (err error) {
	var (
		flag bool
		e    error
	)
	flag, e = isDup(card.Content)
	if e != nil {
		return e
	}
	if flag {
		return fmt.Errorf("ID重复")
	}
	if card.Content, err = parseConten(card.Content); err != nil {
		return
	}
	dbModel := s.showDao.DB.Model(&cardmodel.Setup{})
	dbModel = dbModel.Where("id = ?", id).Where("type = ?", card.Type)
	if err = dbModel.Update(card).Error; err != nil {
		log.Error("chanelSvc.CardSetupList First error(%v)", err)
		return
	}
	if card.Type == common.CardPgcsRcmd {
		if err = util.AddLog(cardmodel.LogBusPgcsRcmd, person, uid, int64(id), cardmodel.ActUpCsPgcRcmd, card); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	} else if card.Type == common.CardUpRcmdNew {
		if err = util.AddLog(cardmodel.LogBusRcmdNew, person, uid, int64(id), cardmodel.ActUpCsRcmdNew, card); err != nil {
			log.Error("chanelSvc.UpdateCardSetup AddLog error(%v)", err)
			return
		}
	}
	return
}
