package up

// .
const (
	CreditBusinessTypeArchive = int8(1)
)

//CreditLog .
type CreditLog struct {
	Type         int8   `json:"type"`
	Optyte       int8   `json:"optype"`
	Reason       int64  `json:"reason"`
	BusinessType int8   `json:"business_type"`
	MID          int64  `json:"mid"`
	OID          int64  `json:"oid"`
	UID          int64  `json:"uid"`
	Content      string `json:"content"`
	Ctime        int64  `json:"ctime"`
	Extra        interface{}
}
