package service

import (
	"fmt"
	"strings"

	"context"
	"go-common/app/admin/main/laser/model"
)

// QueryTaskLog is query the finished task
func (s *Service) QueryTaskLog(ctx context.Context, mid int64, taskID int64, platform int, taskState int, sortBy string, pageNo int, pageSize int) (logs []*model.TaskLog, count int64, err error) {
	var wherePairs []string
	if mid > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("mid = %d", mid))
	}
	if taskID > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("task_id = %d", taskID))
	}
	if platform > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("platform = %d", platform))
	}
	if taskState > 0 {
		wherePairs = append(wherePairs, fmt.Sprintf("task_state = %d", taskState))
	}

	var queryStmt string
	if len(wherePairs) > 0 {
		queryStmt = "WHERE " + strings.Join(wherePairs, " AND ")
	}
	sort := buildSortStmt(sortBy)
	var limit, offset int
	if pageNo >= 0 && pageSize > 0 {
		offset = (pageNo - 1) * pageSize
		limit = pageSize
	} else {
		offset = 0
		limit = 20
	}
	logs, count, err = s.dao.QueryTaskLog(ctx, queryStmt, sort, offset, limit)
	return
}
