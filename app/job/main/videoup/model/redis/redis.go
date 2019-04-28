package redis

//redis action list
const (
	// ActionForDispatchFinish      = "dispatchFinish"
	// ActionForModifyArchive       = "modifyArchive"
	ActionForBvcCapable = "bvcCapable"
	// ActionForSecondRound         = "secondRound"
	ActionForSendOpenMsg    = "sendOpenMsg"
	ActionForSendBblog      = "sendBblog"
	ActionForVideoshot      = "addVideoshot"
	ActionForVideocovers    = "addVideocovers"
	ActionForPostFirstRound = "postFirstRound"
)

//Retry struct
type Retry struct {
	Action string `json:"action"`
}

//RetryJSON struct
type RetryJSON struct {
	Action string `json:"action"`
	Data   struct {
		Aid     int64  `json:"aid"`
		Route   string `json:"route"`
		Mid     int64  `json:"mid"`
		Content string `json:"content"`
	} `json:"data"`
}
