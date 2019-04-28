package model

// BcoinSalaryResp salary dto.
type BcoinSalaryResp struct {
	BcoinList    []*VipBcoinSalary `json:"bcoin_list"`
	DaysNextGive int32             `json:"days_next_give"`
}
