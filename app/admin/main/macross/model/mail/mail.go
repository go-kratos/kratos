package mail

import "mime/multipart"

// Type for mail
type Type uint8

// Mail types
const (
	TypeTextPlain Type = iota
	TypeTextHTML
)

// Attach def.
type Attach struct {
	Name        string
	File        multipart.File
	ShouldUnzip bool
}

// Mail def.
type Mail struct {
	ToAddresses  []*Address `json:"to_addresses"`
	CcAddresses  []*Address `json:"cc_addresses"`
	BccAddresses []*Address `json:"bcc_addresses"`
	Subject      string     `json:"subject"`
	Body         string     `json:"body"`
	Type         Type       `json:"type"`
}

// Address def.
type Address struct {
	Address string `json:"address"`
	Name    string `json:"name"`
}
