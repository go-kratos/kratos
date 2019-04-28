package bws

import "go-common/library/time"

// UserAchieve .
type UserAchieve struct {
	ID    int64     `json:"id"`
	Aid   int64     `json:"aid"`
	Award int64     `json:"award"`
	Ctime time.Time `json:"ctime"`
}

// UserAchieveDetail .
type UserAchieveDetail struct {
	*UserAchieve
	Name          string `json:"name"`
	Icon          string `json:"icon"`
	Dic           string `json:"dic"`
	LockType      int64  `json:"lockType"`
	Unlock        int64  `json:"unlock"`
	Bid           int64  `json:"bid"`
	IconBig       string `json:"icon_big"`
	IconActive    string `json:"icon_active"`
	IconActiveBig string `json:"icon_active_big"`
	SuitID        int64  `json:"suit_id"`
}

// CountAchieves count achieve
type CountAchieves struct {
	Aid   int64 `json:"aid"`
	Count int64 `json:"count"`
}
