package service

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

// uat-i0 i0 ...
var bfsRegexp = regexp.MustCompile(`//.{1,6}\.hdslb+\.com/.+(?:jpg|gif|png|webp|jpeg)$`)

func xssFilter(content string) string {
	p := bluemonday.NewPolicy()
	p.AllowElements("b", "br", "del")
	p.AllowAttrs("target", "href").OnElements("a")
	p.AllowAttrs("class").OnElements("caption", "dl", "dd", "dt", "h2", "h3", "h4", "h5", "h6", "li", "ol", "strong", "ul")
	p.AllowAttrs("class", "style").OnElements("h1", "p", "span")
	p.AllowAttrs("class", "cite").OnElements("blockquote")
	p.AllowAttrs("class", "contenteditable").OnElements("figure", "figcaption", "code")
	p.AllowAttrs("class", "contenteditable", "aid", "style").OnElements("div")
	p.AllowAttrs("color", "size", "face").OnElements("font")
	p.AllowAttrs("class", "contenteditable", "data-lang").OnElements("pre")
	p.AllowAttrs("src", "alt", "title", "width", "aid", "class", "height", "id", "_src", "type", "data-size", "data-vote-id").OnElements("img")
	p.RequireParseableURLs(true)
	p.AllowRelativeURLs(true) // support //i0.hdslb.com
	p.AllowURLSchemes("http", "https", "bilibili")
	p.AllowAttrs("src").Matching(bfsRegexp).OnElements("img")
	return p.Sanitize(content)
}
