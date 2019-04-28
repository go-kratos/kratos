package reply

import (
	"context"
	"go-common/app/interface/main/reply/model/reply"
	"go-common/library/database/sql"
)

const (
	_selEmojiSQL        = "select id,package_id,name,url,state,remark from emoji where state=0 order by sort"
	_selEmoByPidSQL     = "select id,package_id,name,url,state,remark from emoji where state=0 and package_id=? order by sort"
	_selEmojiPackageSQL = "select id,name,url,state from emoji_package where state=0 order by sort"
)

// EmoDao emoji dao
type EmoDao struct {
	db *sql.DB
}

// NewEmojiDao NewEmojiDao
func NewEmojiDao(db *sql.DB) (dao *EmoDao) {
	dao = &EmoDao{
		db: db,
	}
	return
}

// EmojiList get all emoji
func (dao *EmoDao) EmojiList(c context.Context) (emo []*reply.Emoji, err error) {
	rows, err := dao.db.Query(c, _selEmojiSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		emoji := &reply.Emoji{}
		if err = rows.Scan(&emoji.ID, &emoji.PackageID, &emoji.Name, &emoji.URL, &emoji.State, &emoji.Remark); err != nil {
			return
		}
		emo = append(emo, emoji)
	}
	err = rows.Err()
	return
}

// EmojiListByPid get emoji by package_id
func (dao *EmoDao) EmojiListByPid(c context.Context, pid int64) (emo []*reply.Emoji, err error) {
	rows, err := dao.db.Query(c, _selEmoByPidSQL, pid)
	if err != nil {
		return
	}
	defer rows.Close()
	emo = make([]*reply.Emoji, 0)
	for rows.Next() {
		emoji := &reply.Emoji{}
		if err = rows.Scan(&emoji.ID, &emoji.PackageID, &emoji.Name, &emoji.URL, &emoji.State, &emoji.Remark); err != nil {
			return
		}
		emo = append(emo, emoji)
	}
	err = rows.Err()
	return
}

// ListEmojiPack get all emojipack
func (dao *EmoDao) ListEmojiPack(c context.Context) (packs []*reply.EmojiPackage, err error) {
	rows, err := dao.db.Query(c, _selEmojiPackageSQL)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		p := &reply.EmojiPackage{}
		if err = rows.Scan(&p.ID, &p.Name, &p.URL, &p.State); err != nil {
			return
		}
		packs = append(packs, p)
	}
	err = rows.Err()
	return
}
