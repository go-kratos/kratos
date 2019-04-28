package show

import (
	"context"

	"go-common/app/interface/main/app-show/model/show"
	"go-common/library/log"
)

const (
	// real data
	_headSQL = "SELECT s.id,s.plat,s.title,s.type,s.param,s.style,s.rank,s.build,s.conditions,l.name FROM show_head AS s,language AS l WHERE l.id=s.lang_id ORDER BY rank DESC"
	_itemSQL = "SELECT sid,title,random,cover,param FROM show_item"

	// temp preview
	_headTmpSQL = "SELECT s.id,s.plat,s.title,s.type,s.param,s.style,s.rank,s.build,s.conditions,l.name FROM show_head_temp AS s,language AS l WHERE l.id=s.lang_id ORDER BY rank DESC"
	_itemTmpSQL = "SELECT sid,title,random,cover,param FROM show_item_temp"
)

// Heads get show head data.
func (d *Dao) Heads(ctx context.Context) (heads map[int8][]*show.Head, err error) {
	rows, err := d.getHead.Query(ctx)
	if err != nil {
		log.Error("d.getItem.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	heads = make(map[int8][]*show.Head, 20)
	for rows.Next() {
		h := &show.Head{}
		if err = rows.Scan(&h.ID, &h.Plat, &h.Title, &h.Type, &h.Param, &h.Style, &h.Rank, &h.Build, &h.Condition, &h.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		heads[h.Plat] = append(heads[h.Plat], h)
	}
	return
}

// Items get item data.
func (d *Dao) Items(ctx context.Context) (items map[int][]*show.Item, err error) {
	rows, err := d.getItem.Query(ctx)
	if err != nil {
		log.Error("d.getItem.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	items = make(map[int][]*show.Item, 50)
	for rows.Next() {
		i := &show.Item{}
		if err = rows.Scan(&i.Sid, &i.Title, &i.Random, &i.Cover, &i.Param); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		items[i.Sid] = append(items[i.Sid], i)
	}
	return
}

// TempHeads get show temp head data.
func (d *Dao) TempHeads(ctx context.Context) (heads map[int8][]*show.Head, err error) {
	rows, err := d.db.Query(ctx, _headTmpSQL)
	if err != nil {
		log.Error("d.tempHeads.Query() error(%v)", err)
		return
	}
	defer rows.Close()
	heads = make(map[int8][]*show.Head, 20)
	for rows.Next() {
		h := &show.Head{}
		if err = rows.Scan(&h.ID, &h.Plat, &h.Title, &h.Type, &h.Param, &h.Style, &h.Rank, &h.Build, &h.Condition, &h.Language); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		heads[h.Plat] = append(heads[h.Plat], h)
	}
	return
}

// TempItems get temp item data.
func (d *Dao) TempItems(ctx context.Context) (items map[int][]*show.Item, err error) {
	rows, err := d.db.Query(ctx, _itemTmpSQL)
	if err != nil {
		log.Error("d.tempItems.Query error(%v)", err)
		return nil, err
	}
	defer rows.Close()
	items = make(map[int][]*show.Item, 50)
	for rows.Next() {
		i := &show.Item{}
		if err = rows.Scan(&i.Sid, &i.Title, &i.Random, &i.Cover, &i.Param); err != nil {
			log.Error("row.Scan error(%v)", err)
			return nil, err
		}
		items[i.Sid] = append(items[i.Sid], i)
	}
	return
}
