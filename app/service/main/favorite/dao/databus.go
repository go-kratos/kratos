package dao

import (
	"context"
	"strconv"

	"go-common/app/service/main/favorite/model"
)

func (d *Dao) send(c context.Context, mid int64, msg *model.Message) error {
	key := strconv.FormatInt(mid, 10)
	return d.jobDatabus.Send(c, key, msg)
}

// PubAddFav push the add resource event into databus.
func (d *Dao) PubSortFavs(c context.Context, tp int8, mid, fid int64, sorts []model.SortFav) {
	msg := &model.Message{
		Field:    model.FieldResource,
		Action:   model.ActionSortFavs,
		Type:     tp,
		Mid:      mid,
		Fid:      fid,
		SortFavs: sorts,
	}
	d.send(c, mid, msg)
}

// PubAddFav push the add resource event into databus.
func (d *Dao) PubAddFav(c context.Context, tp int8, mid, fid, oid int64, attr int32, ts int64, otype int8) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionAdd,
		Type:       tp,
		Mid:        mid,
		Fid:        fid,
		Oid:        oid,
		FolderAttr: attr,
		FTime:      ts,
		Otype:      otype,
	}
	d.send(c, mid, msg)
}

// PubDelFav push the delete favorite event into databus.
func (d *Dao) PubDelFav(c context.Context, tp int8, mid, fid, oid int64, attr int32, ts int64, otype int8) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionDel,
		Type:       tp,
		Mid:        mid,
		Fid:        fid,
		Oid:        oid,
		FolderAttr: attr,
		FTime:      ts,
		Otype:      otype,
	}
	d.send(c, mid, msg)
}

// PubInitRelationFids push the relationfids cache event into databus.
func (d *Dao) PubInitRelationFids(c context.Context, tp int8, mid int64) {
	msg := &model.Message{
		Field:  model.FieldResource,
		Action: model.ActionInitRelationFids,
		Type:   tp,
		Mid:    mid,
	}
	d.send(c, mid, msg)
}

// PubInitFolderRelations push the folder relations cache event into databus.
func (d *Dao) PubInitFolderRelations(c context.Context, tp int8, mid, fid int64) {
	msg := &model.Message{
		Field:  model.FieldResource,
		Action: model.ActionInitFolderRelations,
		Type:   tp,
		Mid:    mid,
		Fid:    fid,
	}
	d.send(c, mid, msg)
}

// PubInitAllFolderRelations push the folder relations cache event into databus.
func (d *Dao) PubInitAllFolderRelations(c context.Context, tp int8, mid, fid int64) {
	msg := &model.Message{
		Field:  model.FieldResource,
		Action: model.ActionInitAllFolderRelations,
		Type:   tp,
		Mid:    mid,
		Fid:    fid,
	}
	d.send(c, mid, msg)
}

// PubAddFolder push the add folder action event into databus.
func (d *Dao) PubAddFolder(c context.Context, typ int8, mid, fid int64, attr int32) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionFolderAdd,
		Type:       typ,
		Mid:        mid,
		Fid:        fid,
		FolderAttr: attr,
	}
	d.send(c, mid, msg)
}

// PubDelFolder push the del folder action event into databus.
func (d *Dao) PubDelFolder(c context.Context, typ int8, mid, fid int64, attr int32, ts int64) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionFolderDel,
		Type:       typ,
		Mid:        mid,
		Fid:        fid,
		FolderAttr: attr,
		FTime:      ts,
	}
	d.send(c, mid, msg)
}

// PubMultiDelFavs push the multi del fav relations event into databus.
func (d *Dao) PubMultiDelFavs(c context.Context, typ int8, mid, fid, rows int64, attr int32, oids []int64, ts int64) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionMultiDel,
		Type:       typ,
		Mid:        mid,
		Fid:        fid,
		Affected:   rows,
		FolderAttr: attr,
		Oids:       oids,
		FTime:      ts,
	}
	d.send(c, mid, msg)
}

// PubMultiAddFavs push the multi add fav relations event into databus.
func (d *Dao) PubMultiAddFavs(c context.Context, typ int8, mid, fid, rows int64, attr int32, oids []int64, ts int64) {
	msg := &model.Message{
		Field:      model.FieldResource,
		Action:     model.ActionMultiAdd,
		Type:       typ,
		Mid:        mid,
		Fid:        fid,
		Affected:   rows,
		FolderAttr: attr,
		Oids:       oids,
		FTime:      ts,
	}
	d.send(c, mid, msg)
}

// PubMoveFavs push the move resources event into databus.
func (d *Dao) PubMoveFavs(c context.Context, typ int8, mid, ofid, nfid, rows int64, oids []int64, ts int64) {
	msg := &model.Message{
		Field:    model.FieldResource,
		Action:   model.ActionMove,
		Type:     typ,
		Mid:      mid,
		OldFid:   ofid,
		NewFid:   nfid,
		Affected: rows,
		Oids:     oids,
		FTime:    ts,
	}
	d.send(c, mid, msg)
}

// PubCopyFavs push the copy resources event into databus.
func (d *Dao) PubCopyFavs(c context.Context, typ int8, mid, ofid, nfid, rows int64, oids []int64, ts int64) {
	msg := &model.Message{
		Field:    model.FieldResource,
		Action:   model.ActionCopy,
		Type:     typ,
		Mid:      mid,
		OldFid:   ofid,
		NewFid:   nfid,
		Affected: rows,
		Oids:     oids,
		FTime:    ts,
	}
	d.send(c, mid, msg)
}

// PubClean push the clean video event into databus.
func (d *Dao) PubClean(c context.Context, typ int8, mid, fid, ftime int64) {
	msg := &model.Message{
		Field:  model.FieldResource,
		Action: model.ActionClean,
		Type:   typ,
		Mid:    mid,
		Fid:    fid,
		FTime:  ftime,
	}
	d.send(c, mid, msg)
}
