package dao

import (
	"context"
	"fmt"

	"go-common/app/admin/main/member/model"
	"go-common/library/database/sql"
	xtime "go-common/library/time"

	"github.com/pkg/errors"
)

const (
	_reviewName = "user_property_review"
)

// Reviews is. todo delete
func (d *Dao) Reviews(ctx context.Context, mid int64, property []int8, state []int8, isMonitor, isDesc bool, operator string, stime, etime xtime.Time, pn, ps int) ([]*model.UserPropertyReview, int, error) {
	from, to, err := d.prepareReviewRange(ctx, stime, etime)
	if err != nil {
		return nil, 0, err
	}
	order := "asc"
	if isDesc {
		order = "desc"
	}
	rws := []*model.UserPropertyReview{}
	total := 0
	query := d.member.Table(_reviewName).
		Order("id "+order).
		Where("id>=?", from).
		Where("id<=?", to).
		Where("is_monitor=?", isMonitor)

	if mid > 0 {
		query = query.Where("mid=?", mid)
	}
	if len(property) > 0 {
		query = query.Where("property IN (?)", property)
	}
	if operator != "" {
		query = query.Where("operator=?", operator)
	}
	if len(state) > 0 {
		query = query.Where("state IN (?)", state)
	}
	query = query.Count(&total)
	query = query.Offset((pn - 1) * ps).Limit(ps)
	if err := query.Find(&rws).Error; err != nil {
		if err == sql.ErrNoRows {
			return []*model.UserPropertyReview{}, 0, nil
		}
		err = errors.Wrap(err, "reviews")
		return nil, 0, err
	}
	for _, rw := range rws {
		if rw.Property == model.ReviewPropertyFace {
			rw.BuildFaceURL()
		}
	}
	return rws, total, nil
}

func (d *Dao) prepareReviewRange(ctx context.Context, stime, etime xtime.Time) (int64, int64, error) {
	// from id
	rw := &model.UserPropertyReview{}
	if err := d.member.Table(_reviewName).
		Select("id").
		Where("ctime>?", stime).
		Order("ctime asc").
		Limit(1).
		Find(rw).Error; err != nil {
		return 0, 0, err
	}
	from := rw.ID

	// to id
	rw = &model.UserPropertyReview{}
	if err := d.member.Table(_reviewName).
		Select("id").
		Where("ctime<?", etime).
		Order("ctime desc").
		Limit(1).
		Find(rw).Error; err != nil {
		return 0, 0, err
	}
	to := rw.ID

	return from, to, nil
}

// ReviewAudit is.
func (d *Dao) ReviewAudit(ctx context.Context, id []int64, state int8, remark, operator string) error {
	ups := map[string]interface{}{
		"state":    state,
		"remark":   remark,
		"operator": operator,
	}
	if err := d.member.Table(_reviewName).Where("id IN (?)", id).Where("state=?", model.ReviewStateWait).UpdateColumns(ups).Error; err != nil {
		return errors.Wrap(err, "review audit")
	}
	return nil
}

// Review is.
func (d *Dao) Review(ctx context.Context, id int64) (*model.UserPropertyReview, error) {
	r := &model.UserPropertyReview{}
	if err := d.member.Table(_reviewName).Where("id=?", id).First(r).Error; err != nil {
		return nil, err
	}
	return r, nil
}

// ReviewByIDs is.
func (d *Dao) ReviewByIDs(ctx context.Context, ids []int64, state []int8) ([]*model.UserPropertyReview, error) {
	rws := []*model.UserPropertyReview{}
	query := d.member.Table(_reviewName).Where("id IN (?)", ids)
	if len(state) > 0 {
		query = query.Where("state in (?)", state)
	}
	if err := query.Find(&rws).Error; err != nil {
		if err == sql.ErrNoRows {
			return []*model.UserPropertyReview{}, nil
		}
		err = errors.Wrap(err, "review by ids")
		return nil, err
	}
	return rws, nil
}

// UpdateReviewFace is.
func (d *Dao) UpdateReviewFace(ctx context.Context, id int64, face string) error {
	ups := map[string]interface{}{
		"new": face,
	}
	if err := d.member.Table(_reviewName).
		Where("id=?", id).
		Where("property=?", model.ReviewPropertyFace).
		Update(ups).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// MvArchivedFaceToPriv mvArchivedFaceToPriv.
func (d *Dao) MvArchivedFaceToPriv(ctx context.Context, face, privFace, operator, remark string) error {
	ups := map[string]interface{}{
		"new":    privFace,
		"remark": fmt.Sprintf("将归档图片移动到新bucket, face: %s,remark: %s,operator: %s", face, remark, operator),
	}
	if err := d.member.Table(_reviewName).
		Where("new=?", face).
		Where("property=?", model.ReviewPropertyFace).
		Where("state=?", model.ReviewStateArchived).
		Where("is_monitor=?", false).
		Update(ups).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// IncrFaceReject incrFaceReject.
func (d *Dao) IncrFaceReject(ctx context.Context, mid int64) error {
	if err := d.member.Exec("INSERT INTO user_addit (mid,face_reject)VALUES(?,1) ON DUPLICATE KEY UPDATE face_reject=face_reject+1", mid).Error; err != nil {
		return errors.Wrapf(err, "mid: %d", mid)
	}
	return nil
}

// IncrViolationCount is.
func (d *Dao) IncrViolationCount(ctx context.Context, mid int64) error {
	if err := d.member.Exec("INSERT INTO user_addit (mid,violation_count)VALUES(?,1) ON DUPLICATE KEY UPDATE violation_count=violation_count+1", mid).Error; err != nil {
		return errors.Wrapf(err, "mid: %d", mid)
	}
	return nil
}

// FaceAutoPass is.
func (d *Dao) FaceAutoPass(ctx context.Context, ids []int64, etime xtime.Time) (err error) {
	ups := map[string]interface{}{
		"operator": "system/auto",
		"remark":   "48 小时未处理自动通过",
		"state":    model.ReviewStatePass,
	}
	if err = d.member.Table(_reviewName).
		Where("state in (?)", []int8{model.ReviewStateWait, model.ReviewStateQueuing}).
		Where("id in (?)", ids).
		Where("mtime<?", etime).
		Update(ups).Error; err != nil {
		err = errors.WithStack(err)
		return
	}
	return
}

// UpdateRemark is.
func (d *Dao) UpdateRemark(ctx context.Context, id int64, remark string) (err error) {
	ups := map[string]interface{}{
		"remark": remark,
	}
	if err := d.member.Table(_reviewName).
		Where("id=?", id).
		Update(ups).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

//QueuingFaceReviewsByTime is.
func (d *Dao) QueuingFaceReviewsByTime(c context.Context, stime, etime xtime.Time) ([]*model.UserPropertyReview, error) {
	rws := []*model.UserPropertyReview{}
	if err := d.member.Table(_reviewName).
		Where("ctime>? ", stime).
		Where("ctime<? ", etime).
		Where("property=?", model.ReviewPropertyFace).
		Where("state=?", model.ReviewStateQueuing).
		Find(&rws).Error; err != nil {
		return nil, err
	}
	return rws, nil
}

//AuditQueuingFace is
func (d *Dao) AuditQueuingFace(c context.Context, id int64, remark string, state int8) error {
	ups := map[string]interface{}{
		"remark":   remark,
		"state":    state,
		"operator": "system",
	}
	if err := d.member.Table(_reviewName).
		Where("id=?", id).
		Where("state=?", model.ReviewStateQueuing).
		Update(ups).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
