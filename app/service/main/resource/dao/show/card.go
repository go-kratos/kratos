package show

import (
	"context"
	"strconv"
	"time"

	"go-common/app/service/main/resource/model"
	"go-common/library/log"
)

var (
	_appPosRecSQL   = "SELECT r.id,r.tab,r.resource_id,r.type,r.title,r.cover,r.re_type,r.re_value,r.plat_ver,r.desc,r.tag_id FROM app_pos_rec AS r WHERE r.stime<? AND r.etime>? AND r.state=1 AND r.resource_id=3 ORDER BY r.weight ASC"
	_appContentRSQL = "SELECT c.id,c.module,c.rec_id,c.ctype,c.cvalue,c.ctitle,c.tag_id FROM app_content AS c, app_pos_rec AS r WHERE c.rec_id=r.id AND r.state=1 AND r.stime<? AND r.etime>? AND c.module=1"
)

// PosRecs get pos resrouce
func (d *Dao) PosRecs(c context.Context, now time.Time) (res map[int8][]*model.Card, err error) {
	res = map[int8][]*model.Card{}
	rows, err := d.db.Query(c, _appPosRecSQL, now, now)
	if err != nil {
		log.Error("d.PosRecs error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		card := &model.Card{}
		if err = rows.Scan(&card.ID, &card.Tab, &card.RegionID, &card.Type, &card.Title, &card.Cover, &card.Rtype, &card.Rvalue, &card.PlatVer, &card.Desc, &card.TagID); err != nil {
			log.Error("d.PosRecs rows.Scan error(%v)", err)
			res = nil
			return
		}
		for _, limit := range card.CardPlatChange() {
			tmpc := &model.Card{}
			*tmpc = *card
			tmpc.Plat = limit.Plat
			tmpc.Build = limit.Build
			tmpc.Condition = limit.Condition
			tmpc.PlatVer = ""
			tmpc.TypeStr = model.GotoDaily
			tmpc.Goto = model.GotoDaily
			tmpc.Param = tmpc.Rvalue
			tmpc.URI = model.FillURI(tmpc.Goto, tmpc.Param)
			res[tmpc.Plat] = append(res[tmpc.Plat], tmpc)
		}
	}
	err = rows.Err()
	return
}

// RecContents get resource contents
func (d *Dao) RecContents(c context.Context, now time.Time) (res map[int][]*model.Content, aids map[int][]int64, err error) {
	res = map[int][]*model.Content{}
	aids = map[int][]int64{}
	rows, err := d.db.Query(c, _appContentRSQL, now, now)
	if err != nil {
		log.Error("d.RecContents error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		card := &model.Content{}
		if err = rows.Scan(&card.ID, &card.Module, &card.RecID, &card.Type, &card.Value, &card.Title, &card.TagID); err != nil {
			log.Error("d.RecContents rows.Scan error(%v)", err)
			res = nil
			return
		}
		res[card.RecID] = append(res[card.RecID], card)
		if card.Type == model.CardGotoAv {
			aidInt, _ := strconv.ParseInt(card.Value, 10, 64)
			aids[card.RecID] = append(aids[card.RecID], aidInt)
		}
	}
	err = rows.Err()
	return
}
