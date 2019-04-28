package manager

import "go-common/app/admin/main/workflow/model"

// Meta business meta
// http://info.bilibili.co/pages/viewpage.action?pageId=9846887
type Meta struct {
	ID        int     `json:"id"`
	PID       int     `json:"pid"`
	Name      string  `json:"name"`
	Flow      int     `json:"flow"`
	FlowState int     `json:"flow_state"`
	State     int     `json:"state"`
	FlowChild []*Flow `json:"flowchild"`
}

// Flow is child flow meta
type Flow struct {
	FlowState int     `json:"flow_state"`
	Child     []*Meta `json:"child"`
}

// ListResponse .
type ListResponse struct {
	*model.CommonResponse
	Data []*Meta `json:"data"`
}

// Role .
type Role struct {
	ID    int    `json:"id"`
	Bid   int8   `json:"bid"`
	Rid   int8   `json:"rid"`
	Name  string `json:"name"`
	Type  int    `json:"type"`
	State int    `json:"state"`
}

// RoleResponse .
type RoleResponse struct {
	*model.CommonResponse
	Data []*Role
}
