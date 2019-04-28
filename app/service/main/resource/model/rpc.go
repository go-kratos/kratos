package model

// ArgBanner for banner func
type ArgBanner struct {
	Plat      int8
	Build     int
	AID       int64
	MID       int64
	ResIDs    string
	Channel   string
	IP        string
	Buvid     string
	Network   string
	MobiApp   string
	Device    string
	IsAd      bool
	OpenEvent string
	AdExtra   string
	Version   string
}

// ArgRes for resource func
type ArgRes struct {
	ResID int
}

// ArgRess for resources func
type ArgRess struct {
	ResIDs []int
}

// ArgPaster for paster func
type ArgPaster struct {
	Platform int8
	AdType   int8
	Aid      string
	TypeId   string
	TypeID   string
	Buvid    string
}

// ArgCmtbox for ctmbox
type ArgCmtbox struct {
	ID int64
}

// ArgAbTest for abTest
type ArgAbTest struct {
	Groups string
	IP     string
}
