package conf

// PageCfg cfg
type PageCfg struct {
	Name    string // page name
	Top     int    // recom positions
	Middle  int    // last updated
	Bottom  int    // bottom list
	TopM    int    // top interv module id
	MiddleM int    // middle interv module id
}

// ZonesInfo loads all the zones' ID and name
type ZonesInfo struct {
	PGCZonesID    []int          // all the zones' ID that need to be loaded
	UGCZonesID    []int16        // all the ugc zones' ID
	PageIDs       []int          // all the page ID's
	ZonesName     []string       // all the zones' name that need to be loaded
	TargetTypes   []int32        // ugc types that we could pick source data
	UgcTypes      []int32        // ugc archive type order, for search types listing
	OldIdxJump    int            // when the module is UGC and we can't find the old index page corresponding, we tell the client to jump to this category's index page
	OldIdxMapping map[string]int // ugc idx ID old and new mapping
}
