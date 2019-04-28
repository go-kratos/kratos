package dao

import (
	"context"
	"fmt"
	"strconv"

	"go-common/app/admin/main/member/model"

	"github.com/pkg/errors"
)

var (
	_sharding = int64(100)
)

func baseTable(mid int64) string {
	return fmt.Sprintf("user_base_%02d", mid%_sharding)
}

func expTable(mid int64) string {
	return fmt.Sprintf("user_exp_%02d", mid%_sharding)
}

// Base is.
func (d *Dao) Base(ctx context.Context, mid int64) (*model.Base, error) {
	b := &model.Base{}
	if err := d.member.Table(baseTable(mid)).Where("mid=?", mid).Find(b).Error; err != nil {
		return nil, err
	}
	b.RandFaceURL()
	return b, nil
}

// Bases is.
func (d *Dao) Bases(ctx context.Context, mids []int64) (map[int64]*model.Base, error) {
	bs := make(map[int64]*model.Base, len(mids))
	for _, mid := range mids {
		b, err := d.Base(ctx, mid)
		if err != nil {
			continue
		}
		bs[b.Mid] = b
	}
	return bs, nil
}

// Exp is.
func (d *Dao) Exp(ctx context.Context, mid int64) (*model.Exp, error) {
	e := &model.Exp{}
	if err := d.member.Table(expTable(mid)).Where("mid=?", mid).Find(e).Error; err != nil {
		return nil, err
	}
	return e, nil
}

// Moral is.
func (d *Dao) Moral(ctx context.Context, mid int64) (*model.Moral, error) {
	m := &model.Moral{
		Mid:   mid,
		Moral: 7000,
	}
	if err := d.member.Table("user_moral").
		Where(&model.Moral{Mid: mid}).
		FirstOrCreate(m).Error; err != nil {
		return nil, err
	}
	return m, nil
}

// UpName is.
func (d *Dao) UpName(ctx context.Context, mid int64, name string) error {
	ups := map[string]string{
		"name": name,
	}
	if err := d.member.Table(baseTable(mid)).Where("mid=?", mid).Updates(ups).Error; err != nil {
		err = errors.Wrap(err, "dao update name")
		return err
	}
	return nil
}

// UpSign is.
func (d *Dao) UpSign(ctx context.Context, mid int64, sign string) error {
	ups := map[string]string{
		"sign": sign,
	}
	if err := d.member.Table(baseTable(mid)).Where("mid=?", mid).Updates(ups).Error; err != nil {
		err = errors.Wrap(err, "dao update sign")
		return err
	}
	return nil
}

// UpFace is.
func (d *Dao) UpFace(ctx context.Context, mid int64, face string) error {
	ups := map[string]string{
		"face": face,
	}
	if err := d.member.Table(baseTable(mid)).Where("mid=?", mid).Updates(ups).Error; err != nil {
		err = errors.Wrap(err, "dao update face")
		return err
	}
	return nil
}

// BatchUserAddit is.
func (d *Dao) BatchUserAddit(ctx context.Context, mids []int64) (map[int64]*model.UserAddit, error) {
	uas := []*model.UserAddit{}
	if err := d.member.Table("user_addit").Where("mid in (?)", mids).Find(&uas).Error; err != nil {
		return nil, err
	}
	uasMap := make(map[int64]*model.UserAddit, len(uas))
	for _, ua := range uas {
		uasMap[ua.Mid] = ua
	}
	return uasMap, nil
}

// PubExpMsg is.
func (d *Dao) PubExpMsg(ctx context.Context, msg *model.AddExpMsg) error {
	return d.expMsgDatabus.Send(ctx, strconv.FormatInt(msg.Mid, 10), msg)
}

// UserAddit is.
func (d *Dao) UserAddit(ctx context.Context, mid int64) (*model.UserAddit, error) {
	addit := &model.UserAddit{}
	if err := d.member.Table("user_addit").Where("mid=?", mid).Find(addit).Error; err != nil {
		return nil, err
	}
	return addit, nil
}

// UpAdditRemark update remark.
func (d *Dao) UpAdditRemark(ctx context.Context, mid int64, remark string) error {
	addit := &model.UserAddit{
		Mid: mid,
	}
	upRemark := map[string]string{
		"remark": remark,
	}
	if err := d.member.Table("user_addit").Where("mid=?", mid).Assign(upRemark).FirstOrCreate(addit).Error; err != nil {
		err = errors.Wrap(err, "dao insert or update addit")
		return err
	}
	return nil
}
