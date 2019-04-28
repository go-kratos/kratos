// Package http detect localized from http header
// and set localized info to metadata as key 'locale'
// The locale used following the specification defined at
// http://www.rfc-editor.org/rfc/bcp/bcp47.txt.
// Examples are: "en-US", "fr-CH", "es-MX"
package http

import (
	"net/url"
	"sort"
	"strconv"
	"strings"

	bm "go-common/library/net/http/blademaster"
)

type language struct {
	Name      string
	Weighting float64
}

// detectLocalizedWeb detect locale from  HTTP Accept-Language header
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Accept-Language
func detectLocalizedWeb(c *bm.Context) (langs []string) {
	parseLang := func(s string) (language, error) {
		seps := strings.SplitN(s, ";", 2)
		lang := language{Name: seps[0]}
		if len(seps) == 1 {
			lang.Weighting = 1
			return lang, nil
		}
		params, err := url.ParseQuery(seps[1])
		if err != nil {
			return lang, err
		}
		lang.Weighting, err = strconv.ParseFloat(params.Get("q"), 32)
		return lang, err
	}
	items := strings.Split(c.Request.Header.Get("Accept-Language"), ",")
	if items[0] == "" {
		return
	}
	// three language is most common accept language send by browser
	languages := make([]language, 0, len(items))
	for _, s := range items {
		l, err := parseLang(s)
		if err != nil {
			//TODO(weicheng): deal with error
			continue
		}
		languages = append(languages, l)
	}
	sort.Slice(languages, func(i, j int) bool {
		return languages[i].Weighting > languages[j].Weighting
	})
	for i := range languages {
		langs = append(langs, languages[i].Name)
	}
	return
}
