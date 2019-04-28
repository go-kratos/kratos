package http

import (
	"fmt"
	"strings"
	"time"

	"go-common/app/service/main/antispam/model"
	"go-common/app/service/main/antispam/service"
	"go-common/app/service/main/antispam/util"
)

// Condition .
type Condition struct {
	*util.Pagination

	Tag  string
	Tags []string

	Content  string
	Contents []string

	Area      string
	Search    string
	State     string
	HitCounts string

	Order, OrderBy        string
	LimitType, LimitScope string
	StartTime, EndTime    *time.Time
}

// ToServiceCond .
func ToServiceCond(cond *Condition) *service.Condition {
	if cond == nil {
		return nil
	}
	res := &service.Condition{
		Pagination: cond.Pagination,

		Area:    cond.Area,
		Order:   cond.Order,
		OrderBy: cond.OrderBy,

		Tags:     cond.Tags,
		Contents: cond.Contents,

		Search:     cond.Search,
		State:      cond.State,
		HitCounts:  cond.HitCounts,
		LimitType:  cond.LimitType,
		LimitScope: cond.LimitScope,

		StartTime: cond.StartTime,
		EndTime:   cond.EndTime,
	}
	// TODO: how to handler it graceful ?
	if cond.Tag != "" {
		res.Tags = []string{cond.Tag}
	}
	if cond.Content != "" {
		res.Contents = []string{cond.Content}
	}
	// history reasons
	if res.OrderBy == "show_up_counts" {
		res.OrderBy = "hit_counts"
	}
	return res
}

// Valid .
func (c *Condition) Valid() error {
	if c.Pagination != nil {
		if c.CurPage == 0 {
			c.CurPage = 1
		}
		if c.PerPage == 0 {
			c.PerPage = 20
		}
	}
	c.Search, c.Order = strings.TrimSpace(c.Search), strings.TrimSpace(c.Order)
	if c.Order == "" {
		c.Order = model.OrderASC
	} else {
		c.Order = strings.ToUpper(c.Order)
	}
	if c.Order != model.OrderASC && c.Order != model.OrderDESC {
		return fmt.Errorf("Order by should be 'ASC' or 'DESC' but got(%s)", c.Order)
	}
	return nil
}
