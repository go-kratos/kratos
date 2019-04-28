package model

// ArgEquip card equip arg.
type ArgEquip struct {
	Mid    int64
	CardID int64 `form:"id" validate:"required,min=1,gte=1"`
}

// ArgMids card mids arg.
type ArgMids struct {
	Mids []int64 `form:"mids,split" validate:"min=1,max=50"`
}
