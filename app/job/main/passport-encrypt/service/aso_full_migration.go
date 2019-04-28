package service

import (
	"context"

	"go-common/app/job/main/passport-encrypt/model"
	"go-common/library/log"
)

// fullMigration full migration
func (s *Service) fullMigration(start, end, inc, limit int64, group string) {
	var count int64

	for {
		log.Info(" group is %s,start %d, end %d,inc %d,limit %d", group, start, end, inc, limit)
		var (
			res      []*model.OriginAccount
			err      error
			effected int64
		)
		if res, err = s.d.AsoAccounts(context.TODO(), start, end); err != nil {
			log.Error("error is (%+v)\n", err)
			break
		}

		for _, acc := range res {
			enAcc := EncryptAccount(acc)
			if effected, err = s.d.AddAsoAccount(context.TODO(), enAcc); err != nil || effected == 0 {
				log.Error("insert error (%+v)", err)
			}
		}
		start = end
		count = end
		end = end + inc

		if count >= limit {
			break
		}
	}
}
