package model

import (
	"encoding/hex"
)

// User .
type User struct {
	Mid    int64  `json:"mid"`
	UserID string `json:"userid"`
	Pwd    []byte `json:"pwd"`
	Salt   string `json:"salt"`
	Status int8   `json:"status"`
	Tel    []byte `json:"tel"`
	Cid    string `json:"cid"`
	Email  []byte `json:"email"`
}

// DecodeUser .
type DecodeUser struct {
	Mid    int64  `json:"mid"`
	UserID string `json:"userid"`
	Pwd    string `json:"pwd"`
	Salt   string `json:"salt"`
	Status int8   `json:"status"`
	Tel    string `json:"tel"`
	Cid    string `json:"cid"`
	Email  string `json:"email"`
}

// Decode decode user
func (d *User) Decode() *DecodeUser {
	return &DecodeUser{
		Mid:    d.Mid,
		UserID: d.UserID,
		Pwd:    hex.EncodeToString(d.Pwd),
		Salt:   d.Salt,
		Status: d.Status,
		Tel:    hex.EncodeToString(d.Tel),
		Cid:    d.Cid,
		Email:  hex.EncodeToString(d.Email),
	}
}
