package model

import (
	"encoding/json"
)

//VideoDetail qa video task detail
type VideoDetail struct {
	CID          int64   `json:"cid" validate:"required,gt=0"`
	AID          int64   `json:"aid" validate:"required,gt=0"`
	TaskID       int64   `json:"task_id" validate:"omitempty,gt=0"`
	TaskUTime    int64   `json:"task_utime"`
	Attribute    int32   `json:"attribute"`
	TagID        int64   `json:"tag_id" validate:"omitempty,gt=0"`
	ArcTitle     string  `json:"arc_title"`
	ArcTypeID    int64   `json:"arc_typeid" validate:"required,gt=0"`
	AuditStatus  int16   `json:"audit_status" `
	AuditSubmit  string  `json:"audit_submit" validate:"required"`
	AuditDetails string  `json:"audit_details" validate:"required"`
	MID          int64   `json:"mid" validate:"required,gt=0"`
	UPGroups     []int64 `json:"up_groups"`
	Fans         int64   `json:"fans"`
}

//QATaskVideo qa video
type QATaskVideo struct {
	QATask
	VideoDetail
	AttributeList map[string]int32 `json:"attribute_list"`
}

//VideoTaskInfo  info
type VideoTaskInfo struct {
	QATaskVideo
	UPGroupList []*UPGroup    `json:"up_group_list"`
	Warnings    []interface{} `json:"warnings"`
}

//TaskVideoDetail detail
type TaskVideoDetail struct {
	Task         *VideoTaskInfo   `json:"task"`
	Video        *Video           `json:"review"`
	VideoHistory []*VideoOperInfo `json:"video_history"`
}

//GetNote get task submit note
func (qv *QATaskVideo) GetNote() (note string, err error) {
	submit := map[string]string{}
	if err = json.Unmarshal([]byte(qv.AuditSubmit), &submit); err != nil {
		return
	}

	note = submit["note"]
	return
}

//GetAttributeList get task submit attribute
func (qv *QATaskVideo) GetAttributeList() {
	qv.AttributeList = AttributeList(qv.Attribute)
}

//GetWarnings get task warnings
func (qv *VideoTaskInfo) GetWarnings() {
	qv.Warnings = []interface{}{}
	for _, g := range qv.UPGroupList {
		qv.Warnings = append(qv.Warnings, g)
	}
}
