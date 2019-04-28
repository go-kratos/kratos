package v1

// Tuple .
type Tuple struct {
	Key string
	Val string
}

// ShareRequest .
type ShareRequest struct {
	Svid    int64 `form:"svid"`
	Channel int32 `form:"share_channel"`
}

// ShareResponse .
type ShareResponse struct {
	URL    []*Tuple `json:"url"`
	Params []*Tuple `json:"params"`
}

// ShareCallbackRequest .
type ShareCallbackRequest struct {
	Svid    int64  `form:"svid"`
	URL     string `form:"url"`
	Type    string `form:"type"`
	Ctime   int64  `form:"ctime"`
	Channel int32  `form:"share_channel"`
}

// ShareCallbackResponse struct
type ShareCallbackResponse struct {
	ShareCount int32 `json:"share"`
}
