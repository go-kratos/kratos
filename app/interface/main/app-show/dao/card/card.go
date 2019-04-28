package card

import (
	"context"
	"strconv"
	"time"

	"go-common/app/interface/main/app-card/model/card/operate"
	"go-common/app/interface/main/app-show/conf"
	"go-common/app/interface/main/app-show/model"
	"go-common/app/interface/main/app-show/model/card"
	"go-common/library/cache/memcache"
	"go-common/library/database/sql"
	"go-common/library/log"
)

const (
	// daily_selection
	_appColumnSQL     = "SELECT id,tab,resource_id,tpl,name,plat_ver FROM app_column WHERE state=1"
	_appPosRecSQL     = "SELECT p.id,p.tab,p.resource_id,p.type,p.title,p.cover,p.re_type,p.re_value,p.plat_ver,p.desc,p.tag_id FROM app_pos_rec AS p WHERE p.stime<? AND p.etime>? AND p.state=1 ORDER BY p.weight ASC"
	_appContentRSQL   = "SELECT c.id,c.module,c.rec_id,c.ctype,c.cvalue,c.ctitle,c.tag_id FROM app_content AS c, app_pos_rec AS r WHERE c.rec_id=r.id AND r.state=1 AND r.stime<? AND r.etime>? AND c.module=1"
	_appColumnNperSQL = "SELECT n.id,n.column_id,n.name,n.desc,n.nper,n.nper_time,n.cover,n.plat_ver,n.title,n.re_type,n.re_value FROM app_column_nper AS n WHERE n.cron_time<? AND n.state=1 ORDER BY n.nper DESC"
	_appContentNSQL   = "SELECT c.id,c.module,c.rec_id,c.ctype,c.cvalue,c.ctitle,c.tag_id FROM app_content AS c, app_column_nper AS n WHERE c.rec_id=n.id AND n.state=1 AND n.cron_time<? AND c.module=2"
	_appColumnList    = "SELECT c.id,c.name,cn.id,cn.title,cn.plat_ver FROM app_column AS c,app_column_nper AS cn WHERE c.id=cn.column_id AND c.state=1 AND cn.state=1 AND cn.cron_time<? ORDER BY cn.nper DESC"
	// hot card
	_cardSQL = `SELECT c.id,c.title,c.card_type,c.card_value,c.recommand_reason,c.recommand_state,c.priority FROM popular_card AS c
	WHERE c.stime<? AND c.etime>? AND c.check=2 AND c.is_delete=0 ORDER BY c.priority ASC`
	_cardPlatSQL   = `SELECT card_id,plat,conditions,build FROM popular_card_plat WHERE is_delete=0`
	_cardSetSQL    = `SELECT c.id,c.type,c.value,c.title,c.long_title,c.content FROM card_set AS c WHERE c.deleted=0`
	_eventTopicSQL = `SELECT c.id,c.title,c.desc,c.cover,c.re_type,c.re_value,c.corner FROM event_topic AS c WHERE c.deleted=0`
)

// Dao is card dao.
type Dao struct {
	db          *sql.DB
	column      *sql.Stmt
	posRec      *sql.Stmt
	recContent  *sql.Stmt
	nperContent *sql.Stmt
	columnNper  *sql.Stmt
	columnList  *sql.Stmt
	// memcache
	mc     *memcache.Pool
	expire int32
}

func New(c *conf.Config) *Dao {
	d := &Dao{
		db: sql.NewMySQL(c.MySQL.Show),
		// memcache
		mc:     memcache.NewPool(c.Memcache.Cards.Config),
		expire: int32(time.Duration(c.Memcache.Cards.Expire) / time.Second),
	}
	d.column = d.db.Prepared(_appColumnSQL)
	d.posRec = d.db.Prepared(_appPosRecSQL)
	d.recContent = d.db.Prepared(_appContentRSQL)
	d.nperContent = d.db.Prepared(_appContentNSQL)
	d.columnNper = d.db.Prepared(_appColumnNperSQL)
	d.columnList = d.db.Prepared(_appColumnList)
	return d
}

// Columns
func (d *Dao) Columns(ctx context.Context) (res map[int8][]*card.Column, err error) {
	res = map[int8][]*card.Column{}
	rows, err := d.column.Query(ctx)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.Column{}
		if err = rows.Scan(&c.ID, &c.Tab, &c.RegionID, &c.Tpl, &c.Name, &c.PlatVer); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		for _, limit := range c.ColumnPlatChange() {
			tmpc := &card.Column{}
			*tmpc = *c
			tmpc.Plat = limit.Plat
			tmpc.Build = limit.Build
			tmpc.Condition = limit.Condition
			tmpc.PlatVer = ""
			tmpc.ColumnGotoChannge()
			res[tmpc.Plat] = append(res[tmpc.Plat], tmpc)
		}
	}
	return
}

// PosRecs
func (d *Dao) PosRecs(ctx context.Context, now time.Time) (res map[int8]map[int][]*card.Card, err error) {
	res = map[int8]map[int][]*card.Card{}
	rows, err := d.posRec.Query(ctx, now, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.Card{}
		if err = rows.Scan(&c.ID, &c.Tab, &c.RegionID, &c.Type, &c.Title, &c.Cover, &c.Rtype, &c.Rvalue, &c.PlatVer, &c.Desc, &c.TagID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		for _, limit := range c.CardPlatChange() {
			tmpc := &card.Card{}
			*tmpc = *c
			tmpc.Plat = limit.Plat
			tmpc.Build = limit.Build
			tmpc.Condition = limit.Condition
			tmpc.PlatVer = ""
			tmpc.CardGotoChannge()
			if cards, ok := res[tmpc.Plat]; ok {
				cards[tmpc.RegionID] = append(cards[tmpc.RegionID], tmpc)
			} else {
				res[tmpc.Plat] = map[int][]*card.Card{
					tmpc.RegionID: []*card.Card{tmpc},
				}
			}
		}
	}
	return
}

// RecContents
func (d *Dao) RecContents(ctx context.Context, now time.Time) (res map[int][]*card.Content, aids map[int][]int64, err error) {
	res = map[int][]*card.Content{}
	aids = map[int][]int64{}
	rows, err := d.recContent.Query(ctx, now, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.Content{}
		if err = rows.Scan(&c.ID, &c.Module, &c.RecID, &c.Type, &c.Value, &c.Title, &c.TagID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		res[c.RecID] = append(res[c.RecID], c)
		if c.Type == model.CardGotoAv {
			aidInt, _ := strconv.ParseInt(c.Value, 10, 64)
			aids[c.RecID] = append(aids[c.RecID], aidInt)
		}
	}
	return
}

// NperContents
func (d *Dao) NperContents(ctx context.Context, now time.Time) (res map[int][]*card.Content, aids map[int][]int64, err error) {
	res = map[int][]*card.Content{}
	aids = map[int][]int64{}
	rows, err := d.nperContent.Query(ctx, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.Content{}
		if err = rows.Scan(&c.ID, &c.Module, &c.RecID, &c.Type, &c.Value, &c.Title, &c.TagID); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		res[c.RecID] = append(res[c.RecID], c)
		if c.Type == model.CardGotoAv {
			aidInt, _ := strconv.ParseInt(c.Value, 10, 64)
			aids[c.RecID] = append(aids[c.RecID], aidInt)
		}
	}
	return
}

// ColumnNpers
func (d *Dao) ColumnNpers(ctx context.Context, now time.Time) (res map[int8][]*card.ColumnNper, err error) {
	res = map[int8][]*card.ColumnNper{}
	rows, err := d.columnNper.Query(ctx, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.ColumnNper{}
		if err = rows.Scan(&c.ID, &c.ColumnID, &c.Name, &c.Desc, &c.Nper, &c.NperTime, &c.Cover, &c.PlatVer, &c.Title, &c.Rtype, &c.Rvalue); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		for _, limit := range c.ColumnNperPlatChange() {
			tmpc := &card.ColumnNper{}
			*tmpc = *c
			tmpc.Plat = limit.Plat
			tmpc.Build = limit.Build
			tmpc.Condition = limit.Condition
			tmpc.PlatVer = ""
			tmpc.ColumnNperGotoChange()
			res[tmpc.Plat] = append(res[tmpc.Plat], tmpc)
		}
	}
	return
}

// ColumnList
func (d *Dao) ColumnPlatList(ctx context.Context, now time.Time) (res map[int8][]*card.ColumnList, err error) {
	res = map[int8][]*card.ColumnList{}
	rows, err := d.columnList.Query(ctx, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	for rows.Next() {
		c := &card.ColumnList{}
		if err = rows.Scan(&c.Ceid, &c.Cname, &c.Cid, &c.Name, &c.PlatVer); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		for _, limit := range c.ColumnListPlatChange() {
			tmpc := &card.ColumnList{}
			*tmpc = *c
			tmpc.Plat = limit.Plat
			tmpc.Build = limit.Build
			tmpc.Condition = limit.Condition
			tmpc.PlatVer = ""
			res[tmpc.Plat] = append(res[tmpc.Plat], tmpc)
		}
	}
	return
}

// ColumnList
func (d *Dao) ColumnList(ctx context.Context, now time.Time) (res []*card.ColumnList, err error) {
	res = []*card.ColumnList{}
	rows, err := d.columnList.Query(ctx, now)
	if err != nil {
		log.Error("query error(%v)", err)
		return
	}
	for rows.Next() {
		c := &card.ColumnList{}
		if err = rows.Scan(&c.Ceid, &c.Cname, &c.Cid, &c.Name, &c.PlatVer); err != nil {
			log.Error("rows.Scan error(%v)", err)
			res = nil
			return
		}
		res = append(res, c)
	}
	return
}

// Card channel card
func (d *Dao) Card(ctx context.Context, now time.Time) (res []*card.PopularCard, err error) {
	rows, err := d.db.Query(ctx, _cardSQL, now, now)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.PopularCard{}
		var valueStr string
		if err = rows.Scan(&c.ID, &c.Title, &c.Type, &valueStr, &c.Reason, &c.ReasonType, &c.Pos); err != nil {
			return
		}
		c.Value, _ = strconv.ParseInt(valueStr, 10, 64)
		res = append(res, c)
	}
	return
}

// CardPlat channel card  plat
func (d *Dao) CardPlat(ctx context.Context) (res map[int64]map[int8][]*card.PopularCardPlat, err error) {
	res = map[int64]map[int8][]*card.PopularCardPlat{}
	rows, err := d.db.Query(ctx, _cardPlatSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		c := &card.PopularCardPlat{}
		if err = rows.Scan(&c.CardID, &c.Plat, &c.Condition, &c.Build); err != nil {
			return
		}
		if r, ok := res[c.CardID]; !ok {
			res[c.CardID] = map[int8][]*card.PopularCardPlat{
				c.Plat: []*card.PopularCardPlat{c},
			}
		} else {
			r[c.Plat] = append(r[c.Plat], c)
		}
	}
	return
}

// CardSet card set
func (d *Dao) CardSet(ctx context.Context) (res map[int64]*operate.CardSet, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _cardSetSQL); err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64]*operate.CardSet)
	for rows.Next() {
		var (
			c     = &operate.CardSet{}
			value string
		)
		if err = rows.Scan(&c.ID, &c.Type, &value, &c.Title, &c.LongTitle, &c.Content); err != nil {
			return
		}
		c.Value, _ = strconv.ParseInt(value, 10, 64)
		res[c.ID] = c
	}
	return
}

// EventTopic event_topic all
func (d *Dao) EventTopic(ctx context.Context) (res map[int64]*operate.EventTopic, err error) {
	var rows *sql.Rows
	if rows, err = d.db.Query(ctx, _eventTopicSQL); err != nil {
		return
	}
	defer rows.Close()
	res = make(map[int64]*operate.EventTopic)
	for rows.Next() {
		c := &operate.EventTopic{}
		if err = rows.Scan(&c.ID, &c.Title, &c.Desc, &c.Cover, &c.ReType, &c.ReValue, &c.Corner); err != nil {
			return
		}
		res[c.ID] = c
	}
	return
}
