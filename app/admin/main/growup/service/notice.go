package service

import (
	"context"
	"fmt"
	"strings"

	"go-common/app/admin/main/growup/model"
)

// InsertNotice insert notice
func (s *Service) InsertNotice(c context.Context, title string, typ int, platform int, link string, status int) (err error) {
	notice := &model.Notice{
		Title:    title,
		Type:     typ,
		Platform: platform,
		Link:     link,
		Status:   status,
	}
	_, err = s.dao.InsertNotice(c, notice)
	return
}

// Notices notices
func (s *Service) Notices(c context.Context, typ int, status int, platform int, from int, limit int) (total int, notices []*model.Notice, err error) {
	query := queryStr(typ, status, platform)
	total, err = s.dao.NoticeCount(c, query)
	if err != nil {
		return
	}
	notices, err = s.dao.Notices(c, query, from, limit)
	if notices == nil {
		notices = make([]*model.Notice, 0)
	}
	return
}

func queryStr(typ int, status int, platform int) (query string) {
	if typ != 0 {
		query += " AND "
		query += fmt.Sprintf("type=%d", typ)
	}
	if status != 0 {
		query += " AND "
		query += fmt.Sprintf("status=%d", status)
	}

	if platform != 0 {
		query += " AND "
		query += fmt.Sprintf("platform=%d", platform)
	}
	query += " AND is_deleted = 0"
	return
}

// UpdateNotice update notice
func (s *Service) UpdateNotice(c context.Context, typ int, platform int, title string, link string, id int64, status int) (err error) {
	var kv string
	if typ != 0 {
		kv += fmt.Sprintf("type=%d,", typ)
	}

	if platform != 0 {
		kv += fmt.Sprintf("platform=%d,", platform)
	}

	if len(title) != 0 {
		kv += fmt.Sprintf("title='%s',", title)
	}

	if len(link) != 0 {
		kv += fmt.Sprintf("link='%s',", link)
	}

	if status != 0 {
		kv += fmt.Sprintf("status=%d,", status)
	}

	if len(kv) == 0 {
		return
	}
	kv = strings.TrimRight(kv, ",")
	_, err = s.dao.UpdateNotice(c, kv, id)
	return
}
