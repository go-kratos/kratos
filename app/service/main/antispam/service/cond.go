package service

import (
	"time"

	"go-common/app/service/main/antispam/util"
)

// Condition aggregates all query parameters
// if doesn't need pagination, set Condition.Pagination to nil
type Condition struct {
	// set it to nil if you don't want pagination
	// otherwise, page num 1 and page size 20 will be the default
	*util.Pagination

	Offset string
	Limit  string

	Tag  string
	Tags []string

	Content  string
	Contents []string

	Area      string
	State     string
	Search    string
	HitCounts string

	Order, OrderBy        string
	LimitType, LimitScope string
	StartTime, EndTime    *time.Time
	LastModifiedTime      *time.Time
}
