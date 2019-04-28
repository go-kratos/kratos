package model

import (
	"go-common/library/time"
)

//MemberVerify is.
type MemberVerify struct {
	Mid  int64  `json:"mid"`
	Type int64  `json:"type"`
	Desc string `json:"desc"`
}

//AccountInfo is.
type AccountInfo struct {
	Mid      int64     `json:"mid"`
	Name     string    `json:"name"`
	Cert     int64     `json:"cert"`
	CertDesc string    `json:"certdesc"`
	Ts       time.Time `json:"ts"`
}
