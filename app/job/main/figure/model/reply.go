package model

import (
	repmol "go-common/app/job/main/reply/model/reply"
)

// ReplyEvent reply event
type ReplyEvent struct {
	Action  string          `json:"action"`
	Mid     int64           `json:"mid"`
	Subject *repmol.Subject `json:"subject"`
	Reply   *repmol.Reply   `json:"reply"`
	Report  *repmol.Report  `json:"report"`
}

const (
	// EventAdd add reply
	EventAdd = "reply"
	// EventLike add like
	EventLike = "like"
	// EventLikeCancel like cal
	EventLikeCancel = "like_cancel"
	// EventHate hate
	EventHate = "hate"
	// EventHateCancel hate cal
	EventHateCancel = "hate_cancel"
	// EventReportDel user report approved
	EventReportDel = "report_del"
	// EventReportRecover user report recover
	EventReportRecover = "report_recover"
)
