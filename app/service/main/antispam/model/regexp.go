package model

import (
	"fmt"
	"regexp"
	"strings"

	"go-common/app/service/main/antispam/util"
)

const (
	// OperationLimit .
	OperationLimit = "limit"
	// OperationRestrictLimit .
	OperationRestrictLimit = "restrict"
	// OperationPutToWhiteList .
	OperationPutToWhiteList = "white"
	// OperationIgnore .
	OperationIgnore = "ignore"
)

// Regexp .
type Regexp struct {
	ID        int64          `json:"id"`
	Area      string         `json:"area"`
	AdminID   int64          `json:"admin_id"`
	AdminName string         `json:"-"`
	Reg       *regexp.Regexp `json:"-"`
	Name      string         `json:"name"`
	Operation string         `json:"op"`
	Content   string         `json:"content"`
	State     string         `json:"state"`
	CTime     util.JSONTime  `json:"-"`
	MTime     util.JSONTime  `json:"mtime"`
}

// FindString .
func (r *Regexp) FindString(content string) string {
	if hits := r.Reg.FindStringSubmatch(content); len(hits) >= 2 {
		return strings.TrimSpace(hits[1])
	}
	return ""
}

func (r *Regexp) String() string {
	return fmt.Sprintf("name:%s, operation:%s, regexp:%s\n",
		r.Name, r.Operation, r.Content)
}
