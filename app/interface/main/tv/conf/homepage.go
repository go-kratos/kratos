package conf

// TVApp defines the configuration on behalf of the tv app
type TVApp struct {
	MobiApp  string
	Build    string
	Platform string
}

// PageConf defines the configuration for the home page info
type PageConf struct {
	FollowSize    int
	HideIndexShow []string // the zones that we need to hide the index show part
}
