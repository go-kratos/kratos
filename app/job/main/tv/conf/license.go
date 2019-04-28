package conf

import xtime "go-common/library/time"

// Sync struct defines the parameters for the data sync to license owner
type Sync struct {
	HTTPTimeout xtime.Duration
	DialTimeout xtime.Duration
	LogSize     int
	LConf       LicenseConf // conf for the sync with License Owner
	PlayURL     PlayURL     // playurl config
	API         LicenseURL  // license owner url
	Frequency   Duration
	AuditPrefix string // the prefix for audit pgc data
	UGCPrefix   string // the prefix for audit ugc data
	Sign        string
}

// LicenseConf defubes the configuration about the comm with the license owner
type LicenseConf struct {
	// how many programs can be contained in one message
	SizeMsg int
	// cpcode recognized by License owner
	CPCode string
	// number of modified season to sync in one time
	NbSeason int
}

// PlayURL defines the conf to have the temp play URL
type PlayURL struct {
	Upsigsecret string // key of playurl
	Deadline    string // deadline of playurl
	PlayPath    string // path of playurl
	API         string // the api to get the playurl with CID
	Qn          string // quality of the video
	Deadcodes   []int  // playurl response codes, for them we think the video is dead and delete it
}

// Duration defines the frequencies of the data sync/wait
type Duration struct {
	// Modified Season sync frequency
	FreModSeason xtime.Duration
	// how much time wait if error
	ErrorWait xtime.Duration
	// unit: seconds. if it's 3600, that means when we found season is delayed ( not in DB yet ), we postpone all its eps auditing one hour
	AuditDelay int64
	// unit: seconds. used for rejected season case, we re-audit its content in one day
	RejectWait int
	//  one minute for the data to sync ( avoid selecting the same data )
	WaitCall int
}

// LicenseURL defines the API address of the license owner
type LicenseURL struct {
	AddURL       string
	DelSeasonURL string
	DelEPURL     string
	UpdateURL    string
}
