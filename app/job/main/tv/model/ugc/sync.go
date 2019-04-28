package ugc

// SimpleArc provides the fields for license owner sync
type SimpleArc struct {
	AID      int64
	MID      int64
	TypeID   int32
	Videos   int64
	Title    string
	Cover    string
	Content  string
	Duration int64
	Pubtime  string
}

// SimpleVideo provides the fields for license owner sync
type SimpleVideo struct {
	ID          int
	CID         int64
	IndexOrder  int64
	Eptitle     string
	Duration    int64
	Description string
}

// LicSke represents the skeleton of a license audit message
type LicSke struct {
	Arc    *SimpleArc
	Videos []*SimpleVideo
}
