package aids

import (
	"context"
	"strconv"
	"strings"

	"go-common/app/admin/main/app/conf"
	aidsdao "go-common/app/admin/main/app/dao/aids"
	"go-common/app/admin/main/app/model/aids"
	"go-common/library/ecode"
)

// Service aids
type Service struct {
	dao *aidsdao.Dao
}

// New new dao
func New(c *conf.Config) (s *Service) {
	s = &Service{
		dao: aidsdao.New(c),
	}
	return
}

// Save save
func (s *Service) Save(ctx context.Context, aidsStr string) (err error) {
	aidsArr := strings.Split(aidsStr, ",")
	for _, v := range aidsArr {
		tmp, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return ecode.RequestErr
		}
		a := &aids.Param{
			Aid: tmp,
		}
		s.dao.Insert(ctx, a)
	}
	return
}
