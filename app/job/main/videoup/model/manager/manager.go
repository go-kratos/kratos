package manager

// upper type.
const (
	UpperTypeWhite      int8 = 1
	UpperTypeBlack      int8 = 2
	UpperTypePGC        int8 = 3
	UpperTypeUGCX       int8 = 3
	UpperTypePolitices  int8 = 5
	UpperTypeEnterprise int8 = 7
	UpperTypeSigned     int8 = 15
)

// first round audit result for video.
const (
	FirstRoundLock       int16 = -4    //锁定
	FirstRoundRejectBack int16 = -2    //打回
	FirstRoundWait       int16 = -1    //待审
	FirstRoundOpen       int16 = 0     //开放
	FirstRoundLoginOpen  int16 = 10000 //会员开放(需登录，非VIP)
)

//User user info
type User struct {
	ID         int    `json:"uid"`
	Username   string `json:"username"`
	Department string `json:"department"`
}
