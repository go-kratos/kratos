package archive

//UpFromPGC 生产组, UpFromPGCSecret 机密生产组, UpFromCoopera 企业
const (
	UpFromPGC       = 1
	UpFromPGCSecret = 5
	UpFromCoopera   = 6
)

//Addit id,aid,source,redirect_url,dynamic
type Addit struct {
	ID          int64
	Aid         int64
	Source      string
	RedirectURL string
	MissionID   int64
	UpFrom      int8
	OrderID     int
	Dynamic     string
}
