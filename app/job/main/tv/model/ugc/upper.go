package ugc

import "strings"

// Upper reprensents the uppers
type Upper struct {
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

// EasyUp is the simple version of upper
type EasyUp struct {
	MID  int64
	Name string
	Face string
}

// ToUpper transform an EasyUp to Upper
func (u *EasyUp) ToUpper(ori *Upper) *Upper {
	if ori == nil { // if no original data is given
		ori = &Upper{
			Valid: 1,
		}
	}
	return &Upper{
		MID:     u.MID,
		OriFace: u.Face,
		CMSFace: u.Face,
		OriName: u.Name,
		CMSName: u.Name,
		Toinit:  ori.Toinit,
		Submit:  ori.Submit,
		Valid:   ori.Valid,
		Deleted: ori.Deleted,
	}
}

// IsSame returns whether the upper is the same
func (u *Upper) IsSame(name, face string) (f bool, n bool) {
	n = u.OriName == name
	if strings.Contains(u.OriFace, "bfs") &&
		strings.Contains(face, "bfs") {
		f = bfsFName(u.OriFace) == bfsFName(face)
	} else {
		f = u.OriFace == face
	}
	return
}

// bfsFName picks the file name from bfs url
func bfsFName(bfsurl string) (fileName string) {
	var index = strings.LastIndex(bfsurl, "/")
	if index >= 0 && index+1 < len(bfsurl) {
		fileName = bfsurl[index+1:]
	}
	return
}

// ReqSetUp is the structure of request to function in Dao, set upper value
type ReqSetUp struct {
	Value  string
	MID    int64
	UpType int
}
