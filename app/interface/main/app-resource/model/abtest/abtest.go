package abtest

import (
	resource "go-common/app/service/main/resource/model"
)

type List struct {
	ID   int64  `json:"group_id,lomitempty"`
	Name string `json:"group_name,omitempty"`
}

func (l *List) ListChange(r *resource.AbTest) {
	l.ID = r.ID
	l.Name = r.Name
}
