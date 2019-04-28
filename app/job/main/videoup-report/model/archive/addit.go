package archive

const (
	//UpFromPGC pgc
	UpFromPGC = 1
	//UpFromPGCSecret pgc secret
	UpFromPGCSecret = 5
	//UpFromCoopera pgc cooperate
	UpFromCoopera = 6
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

//IsPGC is archive from pgc
func (addit *Addit) IsPGC() bool {
	return addit.UpFrom == UpFromPGC || addit.UpFrom == UpFromPGCSecret || addit.UpFrom == UpFromCoopera
}
