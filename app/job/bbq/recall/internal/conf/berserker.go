package conf

// BerserkerConfig .
type BerserkerConfig struct {
	Keys []*BerserkerKey
	API  []*BerserkerAPI
}

// BerserkerAPI .
type BerserkerAPI struct {
	Name string
	URL  string
}

// BerserkerKey .
type BerserkerKey struct {
	Owner  string
	AppKey string
	Secret string
}
