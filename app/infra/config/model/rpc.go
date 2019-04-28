package model

// ArgConf config param.
type ArgConf struct {
	App      string
	BuildVer string
	Ver      int64
	Env      string
	Hosts    map[string]string
	SType    int8
}

// ArgToken token param.
type ArgToken struct {
	App   string
	Token string
	Env   string
}
