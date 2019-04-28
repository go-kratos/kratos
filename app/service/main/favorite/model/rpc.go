package model

const (
	Article       = int8(1)
	TypeVideo     = int8(2)
	TypeMusic     = int8(3)
	TypeTopic     = int8(4)
	TypePlayVideo = int8(5)
	TypePlayList  = int8(6)
	TypeBangumi   = int8(7)
	TypeMoe       = int8(8)
	TypeComic     = int8(9)
	TypeEsports   = int8(10)
	TypeMediaList = int8(11)
	TypeMusicNew  = int8(12)
)

type ArgAllFolders struct {
	Type   int8
	Mid    int64
	Vmid   int64
	Oid    int64
	RealIP string
}

type ArgFolder struct {
	Type   int8
	Fid    int64
	Mid    int64
	Vmid   int64
	RealIP string
}

type ArgFVmid struct {
	Fid  int64
	Vmid int64
}

func (f *ArgFVmid) MediaID() int64 {
	return f.Fid*100 + f.Vmid%100
}

type ArgFolders struct {
	Type   int8
	Mid    int64
	FVmids []*ArgFVmid
	RealIP string
}

type ArgAddFolder struct {
	Type        int8
	Mid         int64
	Name        string
	Description string
	Cover       string
	Public      int8
	Cookie      string
	AccessKey   string
	RealIP      string
}

type ArgUpdateFolder struct {
	Type        int8
	Fid         int64
	Mid         int64
	Name        string
	Description string
	Cover       string
	Public      int8
	Cookie      string
	AccessKey   string
	RealIP      string
}

type ArgDelFolder struct {
	Type   int8
	Mid    int64
	Fid    int64
	RealIP string
}

type ArgFavs struct {
	Type    int8
	Mid     int64
	Vmid    int64
	Fid     int64
	Tv      int
	Tid     int
	Pn      int
	Ps      int
	Keyword string
	Order   string
	RealIP  string
}

type ArgAdd struct {
	Type   int8
	Mid    int64
	Oid    int64
	Fid    int64
	RealIP string
}

type ArgDel struct {
	Type   int8
	Mid    int64
	Oid    int64
	Fid    int64
	RealIP string
}

type ArgAdds struct {
	Type   int8
	Mid    int64
	Oid    int64
	Fids   []int64
	RealIP string
}

type ArgDels struct {
	Type   int8
	Mid    int64
	Oid    int64
	Fids   []int64
	RealIP string
}

type ArgMultiAdd struct {
	Type   int8
	Mid    int64
	Oids   []int64
	Fid    int64
	RealIP string
}

type ArgMultiDel struct {
	Type   int8
	Mid    int64
	Oids   []int64
	Fid    int64
	RealIP string
}

type ArgIsFav struct {
	Type   int8
	Mid    int64
	Oid    int64
	RealIP string
}
type ArgIsFavs struct {
	Type   int8
	Mid    int64
	Oids   []int64
	RealIP string
}
type ArgInDefaultFolder struct {
	Type   int8
	Mid    int64
	Oid    int64
	RealIP string
}

type ArgIsFavedByFid struct {
	Type   int8
	Mid    int64
	Oid    int64
	Fid    int64
	RealIP string
}

type ArgCntUserFolders struct {
	Type   int8
	Mid    int64
	Vmid   int64
	RealIP string
}

type ArgAddVideo struct {
	Mid       int64
	Fids      []int64
	Aid       int64
	Cookie    string
	AccessKey string
	RealIP    string
}

type ArgFavoredVideos struct {
	Mid    int64
	Aids   []int64
	RealIP string
}

type ArgUsers struct {
	Type   int8
	Oid    int64
	Pn     int
	Ps     int
	RealIP string
}

type ArgTlists struct {
	Type   int8
	Mid    int64
	Vmid   int64
	Fid    int64
	RealIP string
}

type ArgRecents struct {
	Type   int8
	Mid    int64
	Size   int
	RealIP string
}
