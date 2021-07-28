package host

import "net/url"

type Query struct {
	IsSecure bool
}

func GetQuery(url *url.URL) *Query {
	query := new(Query)
	values, ok := url.Query()["isSecure"]
	if ok && len(values) > 0 && values[0] == "true" {
		query.IsSecure = true
	}
	return query
}
