package mcnmodel

//EmptyCmd empty cmd
type EmptyCmd struct {
}

//CmdReloadRank reload rank
type CmdReloadRank struct {
	SignID int64 `form:"sign_id"`
}
