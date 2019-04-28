package pgc

import "fmt"

const (
	_definition = "SD"
)

// License Owner Request message
// License represents the data that we need to send to the license owner for auditing
type License struct {
	TId       string
	InputTime string
	Sign      string
	XMLData   *XMLData
}

// XMLData reprensents the main body of xml data sent to license owner
type XMLData struct {
	Service *Service `xml:"Service"`
}

// Service body+head
type Service struct {
	ID   string `xml:"id,attr"`
	Head *Head
	Body *Body
}

// Head some header info
type Head struct {
	TradeID string `xml:"TradeId"`
	Date    string
	Count   int
}

// Body Media list
type Body struct {
	ProgramSetList *PSList `xml:"programSetList"`
}

// PSList is short for programSetList
type PSList struct {
	ProgramSet []*PS `xml:"programSet"`
}

// PS is short for ProgramSet
type PS struct {
	ProgramSetID     string       `xml:"programSetId"`
	ProgramSetName   string       `xml:"programSetName"`
	ProgramSetClass  string       `xml:"programSetClass"`
	ProgramSetType   string       `xml:"programSetType"`
	ProgramSetPoster string       `xml:"programSetPoster"`
	Portrait         string       `xml:"portrait"` // upper's portrait
	Producer         string       `xml:"producer"` // upper's name
	PublishDate      string       `xml:"publishDate"`
	Copyright        string       `xml:"copyright"`
	ProgramCount     int          `xml:"programCount"`
	CREndData        string       `xml:"cREndDate"`
	DefinitionType   string       `xml:"definitionType"`
	CpCode           string       `xml:"cpCode"`
	PayStatus        int          `xml:"payStatus"`
	PrimitiveName    string       `xml:"primitiveName"`
	Alias            string       `xml:"alias"`
	Zone             string       `xml:"zone"`
	LeadingRole      string       `xml:"leadingRole"`
	ProgramSetDesc   string       `xml:"programSetDesc"`
	Staff            string       `xml:"Staff"`
	SubGenre         string       `xml:"subGenre"`
	ProgramList      *ProgramList `xml:"programList,omitempty"`
}

// ProgramList contains different EP
type ProgramList struct {
	Program []*Program `xml:"program"`
}

// Program represents one EP data
type Program struct {
	ProgramID        string  `xml:"programId"`
	ProgramName      string  `xml:"programName"`
	ProgramPoster    string  `xml:"programPoster"`
	ProgramLength    int     `xml:"programLength"`
	PublishDate      string  `xml:"publishDate"`
	IfPreview        int     `xml:"ifPreview"`
	Number           string  `xml:"number"`
	DefinitionType   string  `xml:"definitionType"`
	PlayCount        int     `xml:"playCount"`
	Drm              int     `xml:"drm"`
	ProgramMediaList *PMList `xml:"programMediaList"`
	ProgramDesc      string  `xml:"programDesc"`
}

// PMList is short for programMediaList
type PMList struct {
	ProgramMedia []*PMedia `xml:"programMedia"`
}

// PMedia is short for ProgramMedia
type PMedia struct {
	MediaID    string `xml:"mediaId"`
	PlayURL    string `xml:"playUrl"`
	Definition string `xml:"definition"`
	HTMLURL    string `xml:"htmlUrl"`
}

// MakePMedia is used to construct PMedia structure
func MakePMedia(prefix, playurl string, cid int64) *PMedia {
	return &PMedia{
		MediaID:    fmt.Sprintf("%s%d", prefix, cid),
		PlayURL:    playurl,
		Definition: _definition,
		HTMLURL:    playurl,
	}
}

// Document is the result structure of license owner's response
type Document struct {
	Response *Response
}

// Response is the main content of response
type Response struct {
	TradeID      string `xml:"TradeId"`
	ResponseCode string
	ResponseInfo string
	ResponseTime string `xml:"responseTime"`
	ErrorList    *ErrorList
}

// ErrorList is the list of error returned by the license owner
type ErrorList struct {
	Error *Error
}

// Error one error body
type Error struct {
	ID      string `xml:"Id"`
	Message string
}

// DelBody is the bodu message of deletion
type DelBody struct {
	ProgramList *ProgramList `xml:"programList"`
}

// CreatePMedia creates PMedia struct
func CreatePMedia(prefix string, epid int, url string) *PMedia {
	return &PMedia{
		MediaID:    prefix + fmt.Sprintf("%d", epid),
		PlayURL:    url,
		Definition: "SD",
		HTMLURL:    url,
	}
}

// CreateProgram creates program
func CreateProgram(prefix string, ep *TVEpContent) *Program {
	r := &Program{
		ProgramID:      prefix + fmt.Sprintf("%d", ep.ID),
		ProgramName:    ep.LongTitle,
		ProgramPoster:  ep.Cover,
		ProgramLength:  int(ep.Length * 60),
		PublishDate:    "1970-01-01",
		IfPreview:      0,
		Number:         ep.Title,
		DefinitionType: "SD",
		PlayCount:      0,
		Drm:            ep.PayStatus,
	}
	r.isPay()
	return r
}

// ReqEpLicCall is the request struct for epLicCall function
type ReqEpLicCall struct {
	EpLic *License
	SID   int64
	Conts []*Content
}

// isPay .
func (p *Program) isPay() {
	if p.Drm == 2 {
		p.Drm = 0
	} else {
		p.Drm = 1
	}
}
