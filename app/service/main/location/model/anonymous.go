package model

// AnonymousIP IP database.
type AnonymousIP struct {
	IsAnonymous       bool `json:"is_anonymous" maxminddb:"is_anonymous"`
	IsAnonymousVPN    bool `json:"is_anonymous_vpn" maxminddb:"is_anonymous_vpn"`
	IsHostingProvider bool `json:"is_hosting_provider" maxminddb:"is_hosting_provider"`
	IsPublicProxy     bool `json:"is_public_proxy" maxminddb:"is_public_proxy"`
	IsTorExitNode     bool `json:"is_tor_exit_node" maxminddb:"is_tor_exit_node"`
}
