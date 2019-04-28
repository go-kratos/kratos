package upper

// Upper reprensents the uppers
type Upper struct {
	ID      int
	MID     int64
	Toinit  int
	Submit  int    // 1=need report
	OriName string // original name
	CMSName string // cms intervened name
	OriFace string // original face
	CMSFace string // cms intervened face
	Valid   int    // auth info: 1=online,0=hidden
	Deleted int
}

// CanShow tells whether the upper can be displayed or not
func (u *Upper) CanShow() bool {
	return u.Valid == 1 && u.Deleted == 0
}
