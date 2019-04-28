package dao

import (
	"context"
	"fmt"
	"go-common/app/job/openplatform/open-sug/model"
	"go-common/library/log"
)

// WishCount ...
func (d *Dao) WishCount(c context.Context, item *model.Item) (wishCount int, err error) {
	row := d.ugcDB.QueryRow(c, fmt.Sprintf(_getWishCount, item.ID%16), item.ID)
	if err = row.Scan(&item.WishCount); err != nil {
		item.WishCount = 0
	}
	return item.WishCount, err
}

// CommentCount ...
func (d *Dao) CommentCount(c context.Context, item *model.Item) (commentCount int, err error) {
	row := d.ugcDB.QueryRow(c, fmt.Sprintf(_getCommentCount, item.ID%128), item.ID)
	if err = row.Scan(&item.CommentCount); err != nil {
		item.CommentCount = 0
	}
	return item.CommentCount, err
}

// SalesCount ...
func (d *Dao) SalesCount(c context.Context, item *model.Item) (salesCount int, err error) {
	row := d.ugcDB.QueryRow(c, _getSaleCount, item.ID)
	if err = row.Scan(&item.CommentCount); err != nil {
		item.CommentCount = 0
	}
	return item.CommentCount, err
}

// InsertMatch ...
func (d *Dao) InsertMatch(c context.Context, item *model.Item, season model.Score) (affect int64, err error) {
	score := int(season.Score * 1000)
	res, err := d.ticketDB.Exec(c, _insertMatch, item.ID, season.SeasonID, season.SeasonName, score, item.Name, item.HeadImg, item.SugImg, item.Name, season.SeasonName, score, item.HeadImg, item.SugImg)
	if err != nil {
		log.Error("(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}

// UpdatePic ...
func (d *Dao) UpdatePic(c context.Context, sug *model.Sug) (affect int64, err error) {
	res, err := d.ticketDB.Exec(c, _updatePic, sug.Item.SugImg, sug.Item.Name, sug.Item.HeadImg, sug.Item.ID, sug.SeasonID)
	if err != nil {
		log.Error("(%v)", err)
		return
	}
	affect, err = res.RowsAffected()
	return
}
