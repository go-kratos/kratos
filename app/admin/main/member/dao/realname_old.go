package dao

import (
	"context"

	"go-common/app/admin/main/member/model"

	"github.com/pkg/errors"
)

// AddOldRealnameIMG is
func (d *Dao) AddOldRealnameIMG(ctx context.Context, img *model.DeDeIdentificationCardApplyImg) error {
	if err := d.account.Create(img).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// AddOldRealnameApply is
func (d *Dao) AddOldRealnameApply(ctx context.Context, apply *model.DeDeIdentificationCardApply) error {
	if err := d.account.Create(apply).Error; err != nil {
		return errors.WithStack(err)
	}
	return nil
}
