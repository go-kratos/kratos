package dao

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"go-common/app/admin/main/member/model"
	"go-common/library/cache/memcache"
	"go-common/library/log"

	"github.com/pkg/errors"
)

var (
	_reasonKey = "reason"
)

func realnameImageKey(IMGData string) string {
	return fmt.Sprintf("realname_image_%s", IMGData)
}

// RealnameInfo is.
func (d *Dao) RealnameInfo(ctx context.Context, mid int64) (info *model.DBRealnameInfo, err error) {
	info = &model.DBRealnameInfo{}
	if err = d.memberRead.Table(info.TableName()).Where("mid = ?", mid).First(&info).Error; err != nil {
		info = nil
		err = errors.WithStack(err)
	}
	return
}

// BatchRealnameInfo is
func (d *Dao) BatchRealnameInfo(ctx context.Context, mids []int64) (map[int64]*model.DBRealnameInfo, error) {
	find := make([]*model.DBRealnameInfo, 0, len(mids))
	db := d.memberRead.Table("realname_info").Where("mid IN (?)", mids).Where("status=?", model.RealnameApplyStatePassed.DBStatus()).Find(&find)
	if err := db.Error; err != nil {
		return nil, err
	}
	result := make(map[int64]*model.DBRealnameInfo, len(find))
	for _, ri := range find {
		result[ri.MID] = ri
	}
	return result, nil
}

// RealnameInfoByCardMD5 is.
func (d *Dao) RealnameInfoByCardMD5(ctx context.Context, cardMD5 string, state int, channel uint8) (infos []*model.DBRealnameInfo, err error) {
	infos = make([]*model.DBRealnameInfo, 0)
	var info *model.DBRealnameInfo
	db := d.memberRead.Table(info.TableName()).Where("card_md5 = ?", cardMD5).Where("channel = ?", channel)
	if state >= 0 {
		db = db.Where("status = ?", state)
	}

	if err = db.Find(&infos).Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "RealnameInfoByCardMD5")
			return
		}
		err = nil
		return
	}
	return
}

// UpdateRealnameAlipayApply is.
func (d *Dao) UpdateRealnameAlipayApply(ctx context.Context, id int64, adminID int64, adminName string, state int, reason string) (err error) {
	var (
		mld *model.DBRealnameAlipayApply
		ups = map[string]interface{}{
			"operator":      adminName,
			"operator_id":   adminID,
			"operator_time": time.Now(),
			"status":        state,
			"reason":        reason,
		}
	)
	if err = d.member.Table(mld.TableName()).Where("id = ?", id).Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateRealnameAlipayApply")
	}
	return
}

// UpdateRealnameInfo is.
func (d *Dao) UpdateRealnameInfo(ctx context.Context, mid int64, state int, reason string) (err error) {
	var (
		mld *model.DBRealnameInfo
		ups = map[string]interface{}{
			"status": state,
			"reason": reason,
		}
	)
	if err = d.member.Table(mld.TableName()).Where("mid = ?", mid).Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateRealnameInfo")
	}
	return
}

// RealnameMainApply is.
func (d *Dao) RealnameMainApply(ctx context.Context, id int64) (apply *model.DBRealnameApply, err error) {
	apply = &model.DBRealnameApply{
		Status: model.RealnameApplyStateNone.DBStatus(),
	}
	if err = d.memberRead.Table(apply.TableName()).Where("id = ?", id).First(&apply).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

// RealnameAlipayApply is.
func (d *Dao) RealnameAlipayApply(ctx context.Context, id int64) (apply *model.DBRealnameAlipayApply, err error) {
	apply = &model.DBRealnameAlipayApply{
		Status: model.RealnameApplyStateNone.DBStatus(),
	}
	if err = d.memberRead.Table(apply.TableName()).Where("id = ?", id).First(&apply).Error; err != nil {
		err = errors.WithStack(err)
	}
	return
}

// RealnameMainList is.
func (d *Dao) RealnameMainList(ctx context.Context, mids []int64, cardType int, country int, opName string, tsFrom, tsTo int64, state int, pn, ps int, isDesc bool) (list []*model.DBRealnameApply, total int, err error) {
	var (
		mdl *model.DBRealnameApply
	)
	db := d.memberRead.Table(mdl.TableName())
	if state >= 0 {
		db = db.Where("status = ?", state)
	}
	if len(mids) > 0 {
		db = db.Where("mid in (?)", mids)
	}
	if cardType >= 0 {
		db = db.Where("card_type = ?", cardType)
	}
	if country >= 0 {
		db = db.Where("country = ?", country)
	}
	if opName != "" {
		db = db.Where("operator = ?", opName)
	}
	if tsTo > 0 {
		timeTo := time.Unix(tsTo, 0)
		db = db.Where("mtime <= ?", timeTo)
	}
	if tsFrom > 0 {
		timeFrom := time.Unix(tsFrom, 0)
		db = db.Where("mtime >= ?", timeFrom)
	}
	if err = db.Count(&total).Error; err != nil {
		err = errors.Wrapf(err, "realname apply list count")
		return
	}
	mtimeSort := "mtime ASC"
	if isDesc {
		mtimeSort = "mtime DESC"
	}

	db = db.Order(mtimeSort).Offset((pn - 1) * ps).Limit(ps)
	if err = db.Find(&list).Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "realname apply list")
			return
		}
		err = nil
		return
	}
	return
}

// RealnameAlipayList is.
func (d *Dao) RealnameAlipayList(ctx context.Context, mids []int64, tsFrom, tsTo int64, state int, pn, ps int, isDesc bool) (list []*model.DBRealnameAlipayApply, total int, err error) {
	var (
		mdl *model.DBRealnameAlipayApply
	)
	db := d.memberRead.Table(mdl.TableName())
	if state >= 0 {
		db = db.Where("status = ?", state)
	}
	if len(mids) > 0 {
		db = db.Where("mid in (?)", mids)
	}
	if tsTo > 0 {
		timeTo := time.Unix(tsTo, 0)
		db = db.Where("mtime <= ?", timeTo)
	}
	if tsFrom > 0 {
		timeFrom := time.Unix(tsFrom, 0)
		db = db.Where("mtime >= ?", timeFrom)
	}
	if err = db.Count(&total).Error; err != nil {
		err = errors.Wrapf(err, "realname apply list count")
		return
	}
	mtimeSort := "mtime ASC"
	if isDesc {
		mtimeSort = "mtime DESC"
	}

	db = db.Order(mtimeSort).Offset((pn - 1) * ps).Limit(ps)

	if err = db.Find(&list).Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "realname apply list")
			return
		}
		err = nil
		return
	}
	return
}

// RealnameSearchCards is.
func (d *Dao) RealnameSearchCards(ctx context.Context, cardMD5s []string) (list []*model.DBRealnameInfo, err error) {
	var (
		mdl *model.DBRealnameInfo
	)
	db := d.memberRead.Table(mdl.TableName()).Where("card_md5 in (?)", cardMD5s)
	if err = db.Find(&list).Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "realname apply list")
			return
		}
		err = nil
		return
	}
	return
}

// UpdateRealnameMainApply .
func (d *Dao) UpdateRealnameMainApply(ctx context.Context, id int, state int, opname string, opid int64, optime time.Time, remark string) (err error) {
	var (
		mld *model.DBRealnameApply
		ups = map[string]interface{}{
			"operator":      opname,
			"operator_id":   opid,
			"operator_time": optime,
			"remark":        remark,
			"remark_status": 1,
			"status":        state,
		}
	)
	if err = d.member.Table(mld.TableName()).Where("id = ?", id).Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateRealnameApply")
	}
	return
}

// UpdateOldRealnameApply .
func (d *Dao) UpdateOldRealnameApply(ctx context.Context, id int64, state int, opname string, opid int64, optime time.Time, remark string) (err error) {
	var (
		ups = map[string]interface{}{
			"operater": opname,
			// "operator_id":   opid,
			"operater_time": optime.Unix(),
			"remark":        remark,
			"remark_status": 1,
			"status":        state,
		}
	)
	if err = d.account.Table("dede_identification_card_apply").Where("id = ?", id).Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateOldRealnameApply")
	}
	return
}

// RealnameApplyIMG .
func (d *Dao) RealnameApplyIMG(ctx context.Context, ids []int64) (imgMap map[int64]*model.DBRealnameApplyIMG, err error) {
	var (
		db   = d.memberRead.Where("id in (?)", ids)
		list []*model.DBRealnameApplyIMG
	)
	imgMap = make(map[int64]*model.DBRealnameApplyIMG)
	if err = db.Find(&list).Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "RealnameApplyIMG")
			return
		}
		err = nil
	}
	for _, l := range list {
		imgMap[l.ID] = l
	}
	return
}

// RealnameApplyCount .
func (d *Dao) RealnameApplyCount(ctx context.Context, mid int64) (count int, err error) {
	var (
		ml        *model.DBRealnameApply
		aml       *model.DBRealnameAlipayApply
		tempCount int
		db        = d.memberRead.Table(ml.TableName()).Where("mid = ?", mid)
		db2       = d.memberRead.Table(aml.TableName()).Where("mid = ?", mid)
	)
	if err = db.Count(&tempCount).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	count += tempCount
	if err = db2.Count(&tempCount).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	count += tempCount
	return
}

// RealnameReasonList .
func (d *Dao) RealnameReasonList(ctx context.Context) (list []string, total int, err error) {
	var (
		conf = &model.DBRealnameConfig{}
		db   = d.memberRead.Where("`key` = ?", _reasonKey).First(&conf)
	)
	if err = db.Error; err != nil {
		if !db.RecordNotFound() {
			err = errors.Wrapf(err, "RealnameReasonList")
			return
		}
		err = nil
	}
	list = decodeReason(conf.Data)
	total = len(list)
	return
}

// UpdateRealnameReason .
func (d *Dao) UpdateRealnameReason(ctx context.Context, list []string) (err error) {
	var (
		conf *model.DBRealnameConfig
		ups  = map[string]interface{}{
			"data": encodeReason(list),
		}
	)
	if err = d.member.Table(conf.TableName()).Where("`key` = ?", _reasonKey).Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateRealnameReason")
	}
	return
}

func encodeReason(list []string) (data string) {
	raw := []byte(strings.Join(list, "(#=_=#)"))
	return base64.StdEncoding.EncodeToString(raw)
}

func decodeReason(data string) (list []string) {
	var (
		raw []byte
		err error
	)
	if raw, err = base64.StdEncoding.DecodeString(data); err != nil {
		err = errors.WithStack(err)
		log.Error("%+v", err)
		return
	}
	list = strings.Split(string(raw), "(#=_=#)")
	return
}

// AddRealnameIMG is
func (d *Dao) AddRealnameIMG(ctx context.Context, img *model.DBRealnameApplyIMG) error {
	if err := d.member.Create(img).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// AddRealnameApply is
func (d *Dao) AddRealnameApply(ctx context.Context, apply *model.DBRealnameApply) error {
	if err := d.member.Create(apply).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// SubmitRealnameInfo is
func (d *Dao) SubmitRealnameInfo(ctx context.Context, info *model.DBRealnameInfo) error {
	ups := map[string]interface{}{
		"channel":   info.Channel,
		"realname":  info.Realname,
		"country":   info.Country,
		"card_type": info.CardType,
		"card":      info.Card,
		"card_md5":  info.CardMD5,
		"status":    info.Status,
		"reason":    info.Reason,
	}
	if err := d.member.Table(info.TableName()).Where("mid=?", info.MID).Assign(ups).FirstOrCreate(info).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// GetRealnameImageCache is
func (d *Dao) GetRealnameImageCache(ctx context.Context, IMGData string) ([]byte, error) {
	key := realnameImageKey(IMGData)
	conn := d.memcache.Get(ctx)
	defer conn.Close()
	item, err := conn.Get(key)
	if err != nil {
		return nil, err
	}
	out := []byte{}
	if err := conn.Scan(item, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// SetRealnameImageCache is
func (d *Dao) SetRealnameImageCache(ctx context.Context, IMGData string, data []byte) error {
	key := realnameImageKey(IMGData)
	conn := d.memcache.Get(ctx)
	defer conn.Close()
	return conn.Set(&memcache.Item{
		Key:        key,
		Value:      data,
		Flags:      memcache.FlagRAW | memcache.FlagGzip,
		Expiration: 6 * 3600, // cache for 6 hours
	})
}

// RecentRealnameApplyImg is
func (d *Dao) RecentRealnameApplyImg(ctx context.Context, duration time.Duration) ([]*model.DBRealnameApplyIMG, error) {
	from := time.Now().Add(-duration)
	result := []*model.DBRealnameApplyIMG{}
	db := d.memberRead.Where("mtime>=?", from).Order("id desc").Find(&result)
	if err := db.Error; err != nil {
		return nil, err
	}
	return result, nil
}

// RejectRealnameMainApply is
func (d *Dao) RejectRealnameMainApply(ctx context.Context, mid int64, opname string, opid int64, remark string) (err error) {
	var (
		mld *model.DBRealnameApply
		ups = map[string]interface{}{
			"operator":      opname,
			"operator_id":   opid,
			"operator_time": time.Now(),
			"remark":        remark,
			"remark_status": 1,
			"status":        model.RealnameApplyStateRejective.DBStatus(),
		}
	)
	if err = d.member.Table(mld.TableName()).
		Where("mid = ?", mid).
		Where("status = ?", model.RealnameApplyStatePassed.DBStatus()).
		Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "RejectRealnameMainApply")
	}
	return
}

// RejectRealnameAlipayApply is
func (d *Dao) RejectRealnameAlipayApply(ctx context.Context, mid int64, opname string, opid int64, reason string) (err error) {
	var (
		mld *model.DBRealnameAlipayApply
		ups = map[string]interface{}{
			"operator":      opname,
			"operator_id":   opid,
			"operator_time": time.Now(),
			"status":        model.RealnameApplyStateRejective.DBStatus(),
			"reason":        reason,
		}
	)
	if err = d.member.Table(mld.TableName()).
		Where("mid = ?", mid).
		Where("status = ?", model.RealnameApplyStatePassed.DBStatus()).
		Updates(ups).Error; err != nil {
		err = errors.Wrapf(err, "UpdateRealnameAlipayApply")
	}
	return
}

// LastPassedRealnameMainApply is
func (d *Dao) LastPassedRealnameMainApply(ctx context.Context, mid int64) (*model.DBRealnameApply, error) {
	apply := &model.DBRealnameApply{}
	if err := d.member.Table("realname_apply").Where("mid=?", mid).
		Where("status=?", model.RealnameApplyStatePassed.DBStatus()).
		Order("id DESC").Limit(1).Last(apply).Error; err != nil {
		return nil, err
	}
	return apply, nil
}

// LastPassedRealnameAlipayApply is
func (d *Dao) LastPassedRealnameAlipayApply(ctx context.Context, mid int64) (*model.DBRealnameAlipayApply, error) {
	apply := &model.DBRealnameAlipayApply{}
	if err := d.member.Table("realname_alipay_apply").Where("mid=?", mid).
		Where("status=?", model.RealnameApplyStatePassed.DBStatus()).
		Order("id DESC").Limit(1).Last(apply).Error; err != nil {
		return nil, err
	}
	return apply, nil
}
