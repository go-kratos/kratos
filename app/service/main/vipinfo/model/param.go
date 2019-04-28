package model

// ArgMid arg mid.
type ArgMid struct {
	Mid int64 `form:"mid" validate:"required,min=1,gte=1"`
}

// ArgMids card mids arg.
type ArgMids struct {
	Mids []int64 `form:"mids,split" validate:"min=1,max=100"`
}
