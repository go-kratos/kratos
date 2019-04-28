package service

import (
	"fmt"
	"strings"
	"time"

	"context"
	"go-common/app/admin/main/laser/model"
	"go-common/library/log"
	xtime "go-common/library/time"
)

// AddTask is add a log task
func (s *Service) AddTask(ctx context.Context, mid int64, username string, adminID int64, logDate int64, contactEmail string, platform int, sourceType int) (err error) {
	t, err := s.dao.FindTask(ctx, mid, 0)
	if err != nil {
		return
	}
	if t != nil {
		err = fmt.Errorf("存在该 mid:%d 未完成的任务,请先删除后再添加", mid)
		log.Error("s.AddTask() error(%v)", err)
		return
	}
	_, err = s.dao.AddTask(ctx, mid, username, adminID, time.Unix(logDate, 0).Format("2006-01-02 15:03:04"), contactEmail, platform, sourceType)
	go s.dao.AddTaskInfoCache(ctx, mid, &model.TaskInfo{
		MID:        mid,
		LogDate:    xtime.Time(logDate),
		SourceType: sourceType,
		Platform:   platform,
		Empty:      false,
	})
	return
}

// DeleteTask is delete a log task
func (s *Service) DeleteTask(ctx context.Context, taskID int64, username string, adminID int64) (err error) {
	t, err := s.dao.QueryTaskInfoByIDSQL(ctx, taskID)
	if err != nil {
		return
	}
	go s.dao.RemoveTaskInfoCache(ctx, t.MID)
	return s.dao.DeleteTask(ctx, taskID, username, adminID)
}

// QueryTask is query Task by query params
func (s *Service) QueryTask(ctx context.Context, mid int64, logDateStart int64, logDateEnd int64, sourceType int, platform int, state int, sortBy string, pageNo int, pageSize int) (tasks []*model.Task, count int64, err error) {
	var wherePairs []string
	wherePairs = append(wherePairs, "is_deleted = 0")
	wherePairs = append(wherePairs, fmt.Sprintf("state = %d", state))
	if sourceType > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("source_type = %d", sourceType))
	}
	if platform > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("platform = %d", platform))
	}
	if mid > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("mid = %d", mid))
	}
	layout := "2006-01-02 15:03:04"
	if logDateStart > 0 && logDateEnd > 0 {
		start := time.Unix(logDateStart, 0).Format(layout)
		end := time.Unix(logDateEnd, 0).Format(layout)
		wherePairs = append(wherePairs, fmt.Sprintf(" log_date between '%s' and '%s' ", start, end))
	} else if logDateStart > 0 {
		start := time.Unix(logDateStart, 0).Format(layout)
		wherePairs = append(wherePairs, fmt.Sprintf("log_date >= '%s' ", start))
	}

	var queryStmt string
	if len(wherePairs) > 0 {
		queryStmt = strings.Join(wherePairs, " AND ")
	}
	sort := buildSortStmt(sortBy)
	var limit, offset int
	if pageNo > 0 && pageSize > 0 {
		offset = (pageNo - 1) * pageSize
		limit = pageSize
	} else {
		offset = 0
		limit = 20
	}
	return s.dao.QueryTask(ctx, queryStmt, sort, offset, limit)
}

func buildSortStmt(sortBy string) (sort string) {
	if sortBy == "" {
		sort = "mtime Desc"
		return
	}
	if strings.HasPrefix(sortBy, "-") {
		sort = strings.TrimPrefix(sortBy, "-") + "  " + "Desc"
	}
	return sort
}

// UpdateTask is update undone task where state = 0.
func (s *Service) UpdateTask(ctx context.Context, username string, adminID int64, taskID int64, mid int64, logDate int64, contactEmail string, sourceType int, platform int) (err error) {
	var updatePairs []string
	updatePairs = append(updatePairs, fmt.Sprintf("source_type = %d", sourceType))
	updatePairs = append(updatePairs, fmt.Sprintf("platform = %d", platform))
	updatePairs = append(updatePairs, fmt.Sprintf("mid = %d", mid))
	updatePairs = append(updatePairs, fmt.Sprintf("log_date = '%s' ", time.Unix(logDate, 0).Format("2006-01-02 15:03:04")))
	updatePairs = append(updatePairs, fmt.Sprintf("contact_email = '%s'", contactEmail))

	if adminID > 0 {
		updatePairs = append(updatePairs, fmt.Sprintf("admin_id = %d", adminID))
	}
	if len(username) != 0 {
		updatePairs = append(updatePairs, fmt.Sprintf("username = '%s' ", username))
	}
	updateStmt := strings.Join(updatePairs, ", ")
	// 0 undone
	state := 0
	// 下列参数前端必传.
	go s.dao.AddTaskInfoCache(ctx, mid, &model.TaskInfo{
		MID:        mid,
		LogDate:    xtime.Time(logDate),
		SourceType: sourceType,
		Platform:   platform,
		Empty:      false,
	})
	return s.dao.UpdateTask(ctx, taskID, state, updateStmt)
}
