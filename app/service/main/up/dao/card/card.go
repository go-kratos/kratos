package card

import (
	"context"
	"fmt"
	"go-common/app/service/main/up/model"
	"go-common/library/database/sql"
	"go-common/library/log"
	"go-common/library/xstr"
)

const (
	_UpInfoBaseColumn = `mid, name_cn, name_en, name_alias, signature, content, nationality, 
nation, gender, blood_type, constellation, height, weight, birth_place, birth_date, occupation, 
tags, masterpieces, school, location, interests, platform, platform_account`
	_countUpSQL             = "SELECT count(distinct mid) FROM card_up"
	_listUpMidSQL           = "SELECT mid FROM card_up order BY mtime DESC"
	_listUpInfoSQL          = "SELECT " + _UpInfoBaseColumn + " FROM card_up limit ? offset ?"
	_getUpInfoByMidSQL      = "SELECT " + _UpInfoBaseColumn + " FROM card_up WHERE mid=?"
	_listUpVideoIDSQL       = "SELECT avid FROM card_up_video WHERE mid=? ORDER BY id DESC"
	_listUpImageSQL         = "SELECT url, height, width FROM card_up_image WHERE mid=? ORDER BY id DESC"
	_listUpAccountSQL       = "SELECT url, title, picture, abstract FROM card_up_account WHERE mid = ?"
	_listUpInfoByMidsSQL    = "SELECT " + _UpInfoBaseColumn + " FROM card_up WHERE mid IN (%s)"
	_listUpVideoIDByMidsSQL = "SELECT mid, avid FROM card_up_video WHERE mid IN (%s) ORDER BY id DESC"
	_listUpImageByMidsSQL   = "SELECT mid, url, height, width FROM card_up_image WHERE mid IN (%s) ORDER BY id DESC"
	_listUpAccountByMidsSQL = "SELECT mid, url, title, picture, abstract FROM card_up_account WHERE mid IN (%s)"
)

// CountUpCard count up num
func (d *Dao) CountUpCard(ctx context.Context) (total int, err error) {
	row := d.db.QueryRow(ctx, _countUpSQL)
	err = row.Scan(&total)
	if err != nil {
		log.Error("CountUpCard row.Scan error(%v)", err)
	}
	return
}

// ListUpInfo page list up mids
func (d *Dao) ListUpInfo(ctx context.Context, offset uint, size uint) (infos []*model.UpCardInfo, err error) {
	rows, err := d.db.Query(ctx, _listUpInfoSQL, size, offset)
	if err != nil {
		log.Error("ListUpInfo d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		info := &model.UpCardInfo{}
		err = rows.Scan(&info.MID,
			&info.NameCN, &info.NameEN, &info.NameAlias,
			&info.Signature,
			&info.Content, &info.Nationality, &info.Nation,
			&info.Gender, &info.BloodType, &info.Constellation,
			&info.Height, &info.Weight, &info.BirthPlace,
			&info.BirthDate, &info.Occupation, &info.Tags,
			&info.Masterpieces, &info.School, &info.Location,
			&info.Interests, &info.Platform, &info.PlatformAccount)
		if err != nil {
			log.Error("ListUpInfo rows.Scan error(%v)", err)
			return
		}
		infos = append(infos, info)
	}

	return
}

// MidUpInfoMap get <mid, UpInfo> map by mids
func (d *Dao) MidUpInfoMap(ctx context.Context, mids []int64) (midUpInfoMap map[int64]*model.UpCardInfo, err error) {
	midUpInfoMap = make(map[int64]*model.UpCardInfo)
	rows, err := d.db.Query(ctx, fmt.Sprintf(_listUpInfoByMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("MidUpInfoMap d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		info := &model.UpCardInfo{}
		err = rows.Scan(&info.MID,
			&info.NameCN, &info.NameEN, &info.NameAlias,
			&info.Signature,
			&info.Content, &info.Nationality, &info.Nation,
			&info.Gender, &info.BloodType, &info.Constellation,
			&info.Height, &info.Weight, &info.BirthPlace,
			&info.BirthDate, &info.Occupation, &info.Tags,
			&info.Masterpieces, &info.School, &info.Location,
			&info.Interests, &info.Platform, &info.PlatformAccount)
		if err != nil {
			log.Error("MidUpInfoMap rows.Scan error(%v)", err)
			return
		}
		midUpInfoMap[info.MID] = info
	}

	return
}

// MidAccountsMap get <mid, Accounts> map by mids
func (d *Dao) MidAccountsMap(ctx context.Context, mids []int64) (midAccountsMap map[int64][]*model.UpCardAccount, err error) {
	midAccountsMap = make(map[int64][]*model.UpCardAccount)
	rows, err := d.db.Query(ctx, fmt.Sprintf(_listUpAccountByMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("MidAccountsMap d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		account := new(model.UpCardAccount)
		err = rows.Scan(&mid, &account.URL, &account.Title, &account.Picture, &account.Desc)
		if err != nil {
			log.Error("MidAccountsMap row.Scan error(%v)", err)
			return
		}
		midAccountsMap[mid] = append(midAccountsMap[mid], account)
	}

	return
}

// MidImagesMap get <mid, Images> map by mids
func (d *Dao) MidImagesMap(ctx context.Context, mids []int64) (midImagesMap map[int64][]*model.UpCardImage, err error) {
	midImagesMap = make(map[int64][]*model.UpCardImage)
	rows, err := d.db.Query(ctx, fmt.Sprintf(_listUpImageByMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("MidImagesMap d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		image := new(model.UpCardImage)
		err = rows.Scan(&mid, &image.URL, &image.Height, &image.Width)
		if err != nil {
			log.Error("MidImagesMap row.Scan error(%v)", err)
			return
		}
		midImagesMap[mid] = append(midImagesMap[mid], image)
	}

	return
}

// MidAvidsMap get <mid, Avids> map by mids
func (d *Dao) MidAvidsMap(ctx context.Context, mids []int64) (midAvidsMap map[int64][]int64, err error) {
	midAvidsMap = make(map[int64][]int64)
	rows, err := d.db.Query(ctx, fmt.Sprintf(_listUpVideoIDByMidsSQL, xstr.JoinInts(mids)))
	if err != nil {
		log.Error("MidAvidsMap d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		var avid int64
		err = rows.Scan(&mid, &avid)
		if err != nil {
			log.Error("MidAvidsMap row.Scan error(%v)", err)
			return
		}
		midAvidsMap[mid] = append(midAvidsMap[mid], avid)
	}

	return
}

// ListUpMID list up mids
func (d *Dao) ListUpMID(ctx context.Context) (mids []int64, err error) {
	rows, err := d.db.Query(ctx, _listUpMidSQL)
	if err != nil {
		log.Error("ListCardBase d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var mid int64
		err = rows.Scan(&mid)
		if err != nil {
			log.Error("ListCardBase row.Scan error(%v)", err)
			return
		}
		mids = append(mids, mid)
	}

	return
}

// GetUpInfo get up info by mid
func (d *Dao) GetUpInfo(ctx context.Context, mid int64) (card *model.UpCardInfo, err error) {
	row := d.db.QueryRow(ctx, _getUpInfoByMidSQL, mid)
	card = &model.UpCardInfo{}
	if err = row.Scan(&card.MID,
		&card.NameCN, &card.NameEN, &card.NameAlias,
		&card.Signature,
		&card.Content, &card.Nationality, &card.Nation,
		&card.Gender, &card.BloodType, &card.Constellation,
		&card.Height, &card.Weight, &card.BirthPlace,
		&card.BirthDate, &card.Occupation, &card.Tags,
		&card.Masterpieces, &card.School, &card.Location,
		&card.Interests, &card.Platform, &card.PlatformAccount); err != nil {
		if err == sql.ErrNoRows {
			card = nil
			err = nil
		} else {
			log.Error("GetUpCard row.Scan error(%v)", err)
			return
		}
	}
	return
}

// ListUpAccount list up accounts by mid
func (d *Dao) ListUpAccount(ctx context.Context, mid int64) (accounts []*model.UpCardAccount, err error) {
	rows, err := d.db.Query(ctx, _listUpAccountSQL, mid)
	if err != nil {
		log.Error("listUpAccount d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		account := new(model.UpCardAccount)
		err = rows.Scan(&account.URL, &account.Title, &account.Picture, &account.Desc)
		if err != nil {
			log.Error("listUpAccount row.Scan error(%v)", err)
			return
		}
		accounts = append(accounts, account)
	}

	return
}

// ListUpImage list up images by mid
func (d *Dao) ListUpImage(ctx context.Context, mid int64) (images []*model.UpCardImage, err error) {
	rows, err := d.db.Query(ctx, _listUpImageSQL, mid)
	if err != nil {
		log.Error("listUpImage d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		image := new(model.UpCardImage)
		err = rows.Scan(&image.URL, &image.Height, &image.Width)
		if err != nil {
			log.Error("listUpImage row.Scan error(%v)", err)
			return
		}
		images = append(images, image)
	}

	return
}

// ListAVID list avids by mid
func (d *Dao) ListAVID(ctx context.Context, mid int64) (avids []int64, err error) {
	rows, err := d.db.Query(ctx, _listUpVideoIDSQL, mid)
	if err != nil {
		log.Error("listUpVideo d.db.Query error(%v)", err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var avid int64
		err = rows.Scan(&avid)
		if err != nil {
			log.Error("listUpVideo row.Scan error(%v)", err)
			return
		}
		avids = append(avids, avid)
	}

	return
}
