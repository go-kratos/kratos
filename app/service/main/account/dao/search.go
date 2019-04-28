package dao

import (
	"context"

	"go-common/app/service/main/account/model"

	"github.com/pkg/errors"
)

// MidsByName is.
func (d *Dao) MidsByName(ctx context.Context, names []string) ([]int64, error) {
	r := d.es.NewRequest("member_user").
		Fields("mid").
		Index("user_base").
		WhereIn("kwname", names).
		Ps(len(names)).
		Pn(1)
	result := &model.SearchMemberResult{}
	if err := r.Scan(ctx, &result); err != nil {
		return nil, errors.WithStack(err)
	}
	return result.Mids(), nil
}
