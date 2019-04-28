package income

import (
	"bytes"
	"context"
	"strconv"

	incomeD "go-common/app/job/main/growup/dao/income"
	model "go-common/app/job/main/growup/model/income"
)

// UpAccountSvr up account service
type UpAccountSvr struct {
	batchSize int
	dao       *incomeD.Dao
}

// NewUpAccountSvr new up account service
func NewUpAccountSvr(dao *incomeD.Dao, batchSize int) (svr *UpAccountSvr) {
	return &UpAccountSvr{
		batchSize: batchSize,
		dao:       dao,
	}
}

// UpAccount get up account
func (s *UpAccountSvr) UpAccount(c context.Context, limit int64) (m map[int64]*model.UpAccount, err error) {
	var id int64
	m = make(map[int64]*model.UpAccount)
	for {
		var um map[int64]*model.UpAccount
		um, id, err = s.dao.UpAccounts(c, id, limit)
		if err != nil {
			return
		}
		if len(um) == 0 {
			break
		}
		for mid, acc := range um {
			m[mid] = acc
		}
	}
	return
}

// BatchInsertUpAccount batch insert up account
func (s *UpAccountSvr) BatchInsertUpAccount(c context.Context, us map[int64]*model.UpAccount) (err error) {
	var (
		buff    = make([]*model.UpAccount, s.batchSize)
		buffEnd = 0
	)
	for _, u := range us {
		if u.DataState != 1 {
			continue
		}
		buff[buffEnd] = u
		buffEnd++
		if buffEnd >= s.batchSize {
			values := upAccountValues(buff[:buffEnd])
			buffEnd = 0
			_, err = s.dao.InsertUpAccount(c, values)
			if err != nil {
				return
			}
		}
	}
	if buffEnd > 0 {
		values := upAccountValues(buff[:buffEnd])
		buffEnd = 0
		_, err = s.dao.InsertUpAccount(c, values)
	}
	return
}

// UpdateUpAccount update up account
func (s *UpAccountSvr) UpdateUpAccount(c context.Context, us map[int64]*model.UpAccount) (err error) {
	for _, u := range us {
		if u.DataState != 2 {
			continue
		}
		var time int
		for {
			var rows int64
			rows, err = s.dao.UpdateUpAccount(c, u.MID, u.Version, u.TotalIncome, u.TotalUnwithdrawIncome)
			if err != nil {
				return
			}
			time++
			if rows > 0 {
				break
			}
			if time >= 10 {
				break
			}
			s.reload(c, u)
		}
	}
	return
}

func (s *UpAccountSvr) reload(c context.Context, upAccount *model.UpAccount) (err error) {
	result, err := s.dao.UpAccount(c, upAccount.MID)
	if err != nil {
		return
	}
	upAccount.TotalIncome = result.TotalIncome
	upAccount.TotalUnwithdrawIncome = result.TotalUnwithdrawIncome
	upAccount.Version = result.Version
	upAccount.WithdrawDateVersion = result.WithdrawDateVersion
	return
}

func upAccountValues(us []*model.UpAccount) (values string) {
	var buf bytes.Buffer
	for _, u := range us {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(u.MID, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(u.HasSignContract))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TotalIncome, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.TotalUnwithdrawIncome, 10))
		buf.WriteByte(',')
		buf.WriteString("'" + u.WithdrawDateVersion + "'")
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(u.Version, 10))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}
