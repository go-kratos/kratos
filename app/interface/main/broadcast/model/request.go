package model

// ChangeRoomReq .
type ChangeRoomReq struct {
	RoomID string `json:"room_id"`
}

// RegisterOpReq .
type RegisterOpReq struct {
	Operation  int32   `json:"operation"`
	Operations []int32 `json:"operations"`
}

// UnregisterOpReq .
type UnregisterOpReq struct {
	Operation  int32   `json:"operation"`
	Operations []int32 `json:"operations"`
}
