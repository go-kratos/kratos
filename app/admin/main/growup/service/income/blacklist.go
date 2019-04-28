package income

import (
	"bytes"
	"context"
	"strconv"

	model "go-common/app/admin/main/growup/model/income"
	"go-common/library/database/sql"
	"go-common/library/log"
)

// ArchiveBlack stop archives income, add archive into av_black_list
func (s *Service) ArchiveBlack(c context.Context, typ int, aIDs []int64, mid int64) (err error) {
	if len(aIDs) == 0 {
		return
	}
	tx, err := s.dao.BeginTran(c)
	if err != nil {
		log.Error("s.dao.BeginTran error(%v)", err)
		return
	}

	if err = s.TxInsertAvBlacklist(c, tx, typ, aIDs, mid, _avBlack, len(aIDs)); err != nil {
		log.Error("s.InsertAvBlacklist error(%v)", err)
		return
	}

	if err = tx.Commit(); err != nil {
		log.Error("tx.Commit error")
	}
	return
}

// GetAvBlackListByAvIds get av_black_list by av_id and ctype
func (s *Service) GetAvBlackListByAvIds(c context.Context, avs []*model.ArchiveIncome, ctype int) (avBMap map[int64]struct{}, err error) {
	avIDMap := make(map[int64]struct{})
	for _, av := range avs {
		avIDMap[av.AvID] = struct{}{}
	}
	avIDList := []int64{}
	for avID := range avIDMap {
		avIDList = append(avIDList, avID)
	}
	avBMap = make(map[int64]struct{})
	if len(avIDList) > 0 {
		avBMap, err = s.dao.ListAvBlackList(c, avIDList, ctype)
		if err != nil {
			log.Error("s.dao.ListAvBlackList error(%v)", err)
			return
		}
	}
	return
}

// TxInsertAvBlacklist insert av_black_list
func (s *Service) TxInsertAvBlacklist(c context.Context, tx *sql.Tx, ctype int, aIDs []int64, mid int64, reason int, count int) (err error) {
	nickname, err := s.dao.GetUpInfoNicknameByMID(c, mid, getUpInfoTable(ctype))
	if err != nil {
		log.Error("s.dao.GetUpInfoNicknameByMID error(%v)", err)
		return
	}
	isDeleted, hasSigned := 0, 0
	if nickname != "" {
		hasSigned = 1
	}

	var buf bytes.Buffer
	for _, id := range aIDs {
		buf.WriteString("(")
		buf.WriteString(strconv.FormatInt(id, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.FormatInt(mid, 10))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(ctype))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(reason))
		buf.WriteByte(',')
		buf.WriteString("\"" + nickname + "\"")
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(hasSigned))
		buf.WriteByte(',')
		buf.WriteString(strconv.Itoa(isDeleted))
		buf.WriteString(")")
		buf.WriteByte(',')
	}
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	vals := buf.String()
	buf.Reset()

	rows, err := s.dao.TxInsertAvBlackList(tx, vals)
	if err != nil {
		tx.Rollback()
		return
	}
	if rows < int64(count) {
		log.Info("TxInsertAvBlackList(%v) rows(%d) < count(%d) error", vals, rows, count)
	}
	return
}
