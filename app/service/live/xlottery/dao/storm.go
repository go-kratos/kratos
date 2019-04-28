package dao

import (
	"context"

	"go-common/app/service/live/xlottery/model"

	"github.com/pkg/errors"
)

// InsertSpecialGift 插入SpecialGift
func (d *Dao) InsertSpecialGift(sg *model.SpecialGift) (int64, error) {

	stmt, err := d.db.Prepare("insert into ap_special_gift (uid,room_id,gift_id,gift_num,create_time, custom_field) values (?,?,?,?,?,?) ")
	if err != nil {
		return 0, errors.WithStack(err)
	}

	result, err := stmt.Exec(context.TODO(), sg.UID, sg.RoomID, sg.GiftID, sg.GiftNum, sg.CreateTime, sg.CustomField)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return result.LastInsertId()
}

// FindBeatByBeatIDAndUID 根据beatid 和 uid 查询 beat
func (d *Dao) FindBeatByBeatIDAndUID(id, uid int64) (*model.Beat, error) {
	row := d.db.QueryRow(context.TODO(), "select id ,uid ,content,status, operator,update_time,ctime ,mtime from ap_user_beats_info where  id = ? and uid = ?", id, uid)
	var b model.Beat
	err := row.Scan(&b.ID, &b.UID, &b.Content, &b.Status, &b.Operator, &b.UpdateTime, &b.Ctime, &b.Mtime)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &b, nil
}

// FindShieldKeyWorkByUID 根据uid 查找屏蔽词
func (d *Dao) FindShieldKeyWorkByUID(uid int64) ([]*model.ShieldKeyWork, error) {
	row, err := d.db.Query(context.TODO(), "select id, uid ,original_keyword, keyword, ctime from ap_shield_keywork where uid = ?", uid)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	skArray := make([]*model.ShieldKeyWork, 0)
	for row.Next() {
		var b model.ShieldKeyWork
		err := row.Scan(&b.ID, &b.UID, &b.OriginalKeyword, &b.KeyWord, &b.Ctime)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		skArray = append(skArray, &b)
	}

	return skArray, nil
}
