package model

// const common
const (
	TimeFormatSec = "2006-01-02 15:04:05"
	TimeFormatDay = "2006-01-02"
)

// MangerInfo struct
type MangerInfo struct {
	OID   int64  `json:"id"`
	Uname string `json:"username"`
}

// EmptyList struct
type EmptyList struct {
	Count int      `json:"total_count"`
	Order string   `json:"order"`
	Sort  string   `json:"sort"`
	PN    int      `json:"pn"`
	PS    int      `json:"ps"`
	List  struct{} `json:"list"`
}

// ArrayUnique struct
func ArrayUnique(ids []int64) (res []int64) {
	length := len(ids)
	for i := 0; i < length; i++ {
		if (i > 0 && ids[i-1] == ids[i]) || ids[i] == 0 {
			continue
		}
		res = append(res, ids[i])
	}
	return
}
