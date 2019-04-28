package model

//EmptyReq nothing
type EmptyReq struct {
}

//WarmUpReq warm up
type WarmUpReq struct {
	Mid    int64 `form:"mid"`
	LastID int   `form:"last_id"`
	Size   int   `form:"size"`
}

//AddStaffReq .
type AddStaffReq struct {
	StaffMid int64 `form:"staff_mid"`
	Aid      int64 `form:"aid"`
}
