package service

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"time"

	"go-common/library/log"
	xtime "go-common/library/time"

	"go-common/app/admin/main/growup/model"
)

// CheatUps get cheat ups
func (s *Service) CheatUps(c context.Context, mid int64, nickname string, from, limit int) (total int, spies []*model.UpSpy, err error) {
	total, err = s.dao.UpSpyCount(c)
	if err != nil {
		return
	}
	var query string
	if mid > 0 {
		query = fmt.Sprintf("WHERE mid=%d", mid)
	}

	if nickname != "" {
		query = fmt.Sprintf("WHERE nickname='%s'", nickname)
	}
	spies, err = s.dao.UpSpies(c, query, from, limit)
	if err != nil {
		return
	}
	if spies == nil {
		spies = make([]*model.UpSpy, 0)
	}
	return
}

// CheatArchives get cheat avs
func (s *Service) CheatArchives(c context.Context, mid, avID int64, nickname string, from, limit int) (total int, spies []*model.ArchiveSpy, err error) {
	var query string
	if mid > 0 {
		query = fmt.Sprintf("WHERE mid=%d", mid)
	}

	if nickname != "" {
		query = fmt.Sprintf("WHERE nickname='%s'", nickname)
	}

	if avID > 0 {
		if len(query) > 0 {
			query += fmt.Sprintf(" AND archive_id=%d", avID)
		} else {
			query = fmt.Sprintf("WHERE archive_id=%d", avID)
		}
	}
	total, err = s.dao.ArchiveSpyCount(c, query)
	if err != nil {
		return
	}
	spies, err = s.dao.ArchiveSpies(c, query, from, limit)
	if err != nil {
		return
	}
	return
}

// ExportCheatUps export up.
func (s *Service) ExportCheatUps(c context.Context, mid int64, nickname string, from, limit int) (res []byte, err error) {
	_, spies, err := s.CheatUps(c, mid, nickname, from, limit)
	if err != nil {
		log.Error("s.ExportSpyUp QuerySpyUp error(%v)", err)
		return
	}
	data := formatSpyUp(spies)
	res, err = FormatCSV(data)
	if err != nil {
		log.Error("s.ExportSpyUp FormatCSV error(%v)", err)
	}
	return
}

func formatSpyUp(infos []*model.UpSpy) (data [][]string) {
	str := []string{"序号", "UID", "昵称", "签约时间", "粉丝数", "作弊粉丝量", "播放量", "作弊播放量"}
	data = append(data, str)
	count := 0
	for _, info := range infos {
		count++
		v := []string{strconv.Itoa(count), strconv.FormatInt(info.MID, 10), info.Nickname, time.Unix(int64(info.SignedAt), 0).Format("2006-01-02 15:04:05"), strconv.Itoa(info.Fans), strconv.Itoa(info.CheatFans), strconv.Itoa(info.PlayCount), strconv.Itoa(info.CheatPlayCount)}
		data = append(data, v)
	}
	return
}

// ExportCheatAvs export cheat avs
func (s *Service) ExportCheatAvs(c context.Context, mid, avID int64, nickname string, from, limit int) (res []byte, err error) {
	_, infos, err := s.CheatArchives(c, mid, avID, nickname, from, limit)
	if err != nil {
		log.Error("s.ExportSpyAV QuerySpyArchive error(%v)", err)
		return
	}
	data := fromatSpyAV(infos)
	res, err = FormatCSV(data)
	if err != nil {
		log.Error("s.ExportSpyAV FormatCSV error(%v)", err)
	}
	return
}

func fromatSpyAV(infos []*model.ArchiveSpy) (data [][]string) {
	str := []string{"序号", "稿件ID", "UP昵称", "投稿时间", "累计收入", "作弊播放量", "作弊收藏量", "作弊硬币量"}
	data = append(data, str)
	count := 0
	for _, info := range infos {
		count++
		v := []string{strconv.Itoa(count), strconv.FormatInt(info.ArchiveID, 10), info.Nickname, time.Unix(int64(info.UploadTime), 0).Format("2006-01-02 15:04:05"), strconv.FormatFloat(float64(info.TotalIncome)/100, 'f', 2, 32), strconv.Itoa(info.CheatPlayCount), strconv.Itoa(info.CheatFavorite), strconv.Itoa(info.CheatCoin)}
		data = append(data, v)
	}
	return
}

// QueryCheatFans query cheat fans.
func (s *Service) QueryCheatFans(c context.Context, from, limit int64) (total int64, fans []*model.CheatFans, err error) {
	total, err = s.dao.CheatFansCount(c)
	if err != nil {
		return
	}
	fans, err = s.dao.CheatFans(c, from, limit)
	return
}

// CheatFans handle checkfans mid.
func (s *Service) CheatFans(c context.Context, mid int64) (err error) {
	cu, err := s.cheatFans(c, mid)
	if err != nil {
		return
	}

	if cu.RealFans >= 1000 {
		return
	}
	_, err = s.dao.InsertCheatFansInfo(c, cheatFansValues(cu))
	return
}

func (s *Service) cheatFans(c context.Context, mid int64) (up *model.CheatFans, err error) {
	up = &model.CheatFans{}
	v, err := s.dao.UpInfo(c, mid, 3)
	if err != nil {
		return
	}
	if v.MID == 0 {
		return
	}
	up.MID = v.MID
	up.Nickname = v.Nickname
	up.SignedAt = v.SignedAt
	up.DeductAt = xtime.Time(time.Now().Unix())

	up.RealFans, err = s.dao.GetUpRealFansCount(c, s.conf.Host.Common, v.MID)
	if err != nil {
		return
	}
	up.CheatFans, err = s.dao.GetUpCheatFansCount(c, s.conf.Host.Common, v.MID)
	if err != nil {
		return
	}
	return
}

func cheatFansValues(cu *model.CheatFans) (values string) {
	var buf bytes.Buffer
	buf.WriteString("(")
	buf.WriteString(strconv.FormatInt(cu.MID, 10))
	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf("'%s'", cu.Nickname))
	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf("'%s'", cu.SignedAt.Time().Format("2006-01-02 15:04:05")))
	buf.WriteByte(',')
	buf.WriteString(strconv.Itoa(cu.RealFans))
	buf.WriteByte(',')
	buf.WriteString(strconv.Itoa(cu.CheatFans))
	buf.WriteByte(',')
	buf.WriteString(fmt.Sprintf("'%s'", cu.DeductAt.Time().Format("2006-01-02 15:04-05")))
	buf.WriteString(")")
	buf.WriteByte(',')
	if buf.Len() > 0 {
		buf.Truncate(buf.Len() - 1)
	}
	values = buf.String()
	buf.Reset()
	return
}
