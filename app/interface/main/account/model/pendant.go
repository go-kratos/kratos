package model

import (
	"strconv"

	usmdl "go-common/app/service/main/usersuit/model"
)

// EquipPHP struct.
type EquipPHP struct {
	Pid        int64   `json:"pid"`
	Coins      float64 `json:"coins"`
	Image      string  `json:"image"`
	ImageModel string  `json:"image_model"`
	FaceURL    string  `json:"face_url"`
}

// GroupPHP php group result
type GroupPHP struct {
	Name    string           `json:"group_name"`
	Count   int64            `json:"group_count"`
	Pendant []*usmdl.Pendant `json:"pendant_info"`
}

// PendantPHP php pendant result
type PendantPHP struct {
	Pid        int64  `json:"pid"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageModel string `json:"image_model"`
}

// GroupEntryPHP php vip pendant result
type GroupEntryPHP struct {
	Pid        int64  `json:"pid"`
	Money      int64  `json:"money"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageModel string `json:"image_model"`
}

// GroupVipPHP php vip pendant result
type GroupVipPHP struct {
	Pid        int64  `json:"pid"`
	Money      int64  `json:"money"`
	MoneyType  int8   `json:"money_type"`
	Expire     int64  `json:"display_expire"`
	Name       string `json:"name"`
	Image      string `json:"image"`
	ImageModel string `json:"image_model"`
}

// MyPHP struct.
type MyPHP struct {
	Pid         int64  `json:"pid"`
	Name        string `json:"name"`
	MoneyType   int8   `json:"money_type"`
	Image       string `json:"image"`
	ImageModel  string `json:"image_model"`
	Expire      int64  `json:"expire"`
	IsActivated int8   `json:"is_activated"`
	IsOnline    int8   `json:"is_online"`
	IsVip       int8   `json:"is_vip"`
}

// MyHistoryPHP struct.
type MyHistoryPHP struct {
	Pid        int64  `json:"pid"`
	Image      string `json:"image"`
	Name       string `json:"name"`
	BuyTime    int64  `json:"buy_time"`
	PayID      string `json:"pay_id"`
	Cost       string `json:"cost"`
	TimeLength int64  `json:"time_length"`
}

// FormatImgURL format images url
func FormatImgURL(mid int64, img string) (url string) {
	if len(img) > 0 {
		return "http://i" + strconv.FormatInt(mid%3, 10) + ".hdslb.com" + img
	}
	return
}
