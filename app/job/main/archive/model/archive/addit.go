package archive

const (
	UpFromPGC       = 1
	UpFromPGCSecret = 5
)

// Addit addit struct
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
