package archive

// Addit id,aid,source,redirect_url,description
type Addit struct {
	ID          int64
	Aid         int64
	Source      string
	RedirectURL string
	MissionID   int64
	UpFrom      int8
	OrderID     int
	Description string
}
