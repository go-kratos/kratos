package model

import (
	xtime "go-common/library/time"
)

// ArchiveAuditCargoHour is table archive_audit_cargo_hour.
type ArchiveAuditCargoHour struct {
	ID           int64      `json:"id"`
	UID          int64      `json:"uid"`
	StatDate     xtime.Time `json:"stat_date"`
	ReceiveValue int64      `json:"receive_value"`
	AuditValue   int64      `json:"audit_value"`
	Ctime        xtime.Time `json:"ctime"`
	Mtime        xtime.Time `json:"mtime"`
	State        int        `json:"state"`
}

// CargoDetail is archive audit detail.
type CargoDetail struct {
	UID          int64      `json:"uid"`
	StatDate     xtime.Time `json:"stat_date"`
	ReceiveValue int64      `json:"receive_value"`
	AuditValue   int64      `json:"audit_value"`
}

// CargoItem is audit value which is received or done.
type CargoItem struct {
	ReceiveValue int64 `json:"auditing"`
	AuditValue   int64 `json:"audited"`
}

// CargoView is json data compromised contracted with web front.
type CargoView struct {
	Date string             `json:"date"`
	Data map[int]*CargoItem `json:"data"`
}

// CargoViewWrapper is json data for show the archive cargo audit of every auditor.
type CargoViewWrapper struct {
	Username string `json:"username"`
	*CargoView
}
