package dao

import "go-common/app/service/main/antispam/util"

// Condition .
type Condition struct {
	*util.Pagination

	Offset string
	Limit  string

	Tags     []string
	Contents []string

	Area      string
	Search    string
	State     string
	HitCounts string

	Order, OrderBy        string
	LimitType, LimitScope string
	StartTime, EndTime    string
	LastModifiedTime      string
}
