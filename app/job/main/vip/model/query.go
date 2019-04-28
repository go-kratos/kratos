package model

import "go-common/library/time"

//QueryBcoinSalary .
type QueryBcoinSalary struct {
	StartID       int64     `json:"start_id"`
	EndID         int64     `json:"end_id"`
	StartMonth    time.Time `json:"start_month"`
	EndMonth      time.Time `json:"end_month"`
	GiveNowStatus int8      `json:"give_now_status"`
	Status        int8      `json:"status"`
}
