package model

import (
	"encoding/json"
	"go-common/library/ecode"
	xtime "go-common/library/time"
	"time"
)

const (
	// MaxInt64 用于最大int64
	MaxInt64 = int64(^uint64(0) >> 1)
	// UserListLen 空间长度
	UserListLen = 20
)

// Cache
const (
	CacheKeyUserBase    = "user_base:%d" //用户基本信息缓存key
	CacheExpireUserBase = 600
)

// CursorValue 用于cursor的定位，这里可以当做通用结构使用，使用者自己根据需求定义cursor_id的含义
type CursorValue struct {
	CursorID   int64      `json:"cursor_id"`
	CursorTime xtime.Time `json:"cursor_time"`
}

// ParseCursor 从cursor_prev和cursor_next，判断请求的方向，以及生成cursor
func ParseCursor(cursorPrev string, cursorNext string) (cursor CursorValue, directionNext bool, err error) {
	// 判断是向前还是向后查询
	directionNext = true
	cursorStr := cursorNext
	if len(cursorNext) == 0 && len(cursorPrev) > 0 {
		directionNext = false
		cursorStr = cursorPrev
	}
	// 解析cursor中的cursor_id
	if len(cursorStr) != 0 {
		var cursorData = []byte(cursorStr)
		err = json.Unmarshal(cursorData, &cursor)
		if err != nil {
			err = ecode.ReqParamErr
			return
		}
	}
	// 第一次请求的时候，携带的svid=0，需要转成max传给dao层
	if directionNext && cursor.CursorID == 0 {
		cursor.CursorID = MaxInt64
		cursor.CursorTime = xtime.Time(time.Now().Unix())
	}
	return
}
