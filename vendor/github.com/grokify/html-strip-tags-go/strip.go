// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package strip

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"text/template"
	"text/template/parse"
	"unicode"
	"unicode/utf8"
)

// htmlNospaceEscaper escapes for inclusion in unquoted attribute values.
func htmlNospaceEscaper(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeHTML {
		return htmlReplacer(StripTags(s), htmlNospaceNormReplacementTable, false)
	}
	return htmlReplacer(s, htmlNospaceReplacementTable, false)
}

// attrEscaper escapes for inclusion in quoted attribute values.
func attrEscaper(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeHTML {
		return htmlReplacer(StripTags(s), htmlNormReplacementTable, true)
	}
	return htmlReplacer(s, htmlReplacementTable, true)
}

// rcdataEscaper escapes for inclusion in an RCDATA element body.
func rcdataEscaper(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeHTML {
		return htmlReplacer(s, htmlNormReplacementTable, true)
	}
	return htmlReplacer(s, htmlReplacementTable, true)
}

// htmlEscaper escapes for inclusion in HTML text.
func htmlEscaper(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeHTML {
		return s
	}
	return htmlReplacer(s, htmlReplacementTable, true)
}

// htmlReplacementTable contains the runes that need to be escaped
// inside a quoted attribute value or in a text node.
var htmlReplacementTable = []string{
	// http://www.w3.org/TR/html5/syntax.html#attribute-value-(unquoted)-state
	// U+0000 NULL Parse error. Append a U+FFFD REPLACEMENT
	// CHARACTER character to the current attribute's value.
	// "
	// and similarly
	// http://www.w3.org/TR/html5/syntax.html#before-attribute-value-state
	0:    "\uFFFD",
	'"':  "&#34;",
	'&':  "&amp;",
	'\'': "&#39;",
	'+':  "&#43;",
	'<':  "&lt;",
	'>':  "&gt;",
}

// htmlNormReplacementTable is like htmlReplacementTable but without '&' to
// avoid over-encoding existing entities.
var htmlNormReplacementTable = []string{
	0:    "\uFFFD",
	'"':  "&#34;",
	'\'': "&#39;",
	'+':  "&#43;",
	'<':  "&lt;",
	'>':  "&gt;",
}

// htmlNospaceReplacementTable contains the runes that need to be escaped
// inside an unquoted attribute value.
// The set of runes escaped is the union of the HTML specials and
// those determined by running the JS below in browsers:
// <div id=d></div>
// <script>(function () {
// var a = [], d = document.getElementById("d"), i, c, s;
// for (i = 0; i < 0x10000; ++i) {
//   c = String.fromCharCode(i);
//   d.innerHTML = "<span title=" + c + "lt" + c + "></span>"
//   s = d.getElementsByTagName("SPAN")[0];
//   if (!s || s.title !== c + "lt" + c) { a.push(i.toString(16)); }
// }
// document.write(a.join(", "));
// })()</script>
var htmlNospaceReplacementTable = []string{
	0:    "&#xfffd;",
	'\t': "&#9;",
	'\n': "&#10;",
	'\v': "&#11;",
	'\f': "&#12;",
	'\r': "&#13;",
	' ':  "&#32;",
	'"':  "&#34;",
	'&':  "&amp;",
	'\'': "&#39;",
	'+':  "&#43;",
	'<':  "&lt;",
	'=':  "&#61;",
	'>':  "&gt;",
	// A parse error in the attribute value (unquoted) and
	// before attribute value states.
	// Treated as a quoting character by IE.
	'`': "&#96;",
}

// htmlNospaceNormReplacementTable is like htmlNospaceReplacementTable but
// without '&' to avoid over-encoding existing entities.
var htmlNospaceNormReplacementTable = []string{
	0:    "&#xfffd;",
	'\t': "&#9;",
	'\n': "&#10;",
	'\v': "&#11;",
	'\f': "&#12;",
	'\r': "&#13;",
	' ':  "&#32;",
	'"':  "&#34;",
	'\'': "&#39;",
	'+':  "&#43;",
	'<':  "&lt;",
	'=':  "&#61;",
	'>':  "&gt;",
	// A parse error in the attribute value (unquoted) and
	// before attribute value states.
	// Treated as a quoting character by IE.
	'`': "&#96;",
}

// htmlReplacer returns s with runes replaced according to replacementTable
// and when badRunes is true, certain bad runes are allowed through unescaped.
func htmlReplacer(s string, replacementTable []string, badRunes bool) string {
	written, b := 0, new(bytes.Buffer)
	for i, r := range s {
		if int(r) < len(replacementTable) {
			if repl := replacementTable[r]; len(repl) != 0 {
				b.WriteString(s[written:i])
				b.WriteString(repl)
				// Valid as long as replacementTable doesn't
				// include anything above 0x7f.
				written = i + utf8.RuneLen(r)
			}
		} else if badRunes {
			// No-op.
			// IE does not allow these ranges in unquoted attrs.
		} else if 0xfdd0 <= r && r <= 0xfdef || 0xfff0 <= r && r <= 0xffff {
			fmt.Fprintf(b, "%s&#x%x;", s[written:i], r)
			written = i + utf8.RuneLen(r)
		}
	}
	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}

// stripTags takes a snippet of HTML and returns only the text content.
// For example, `<b>&iexcl;Hi!</b> <script>...</script>` -> `&iexcl;Hi! `.
func StripTags(html string) string {
	var b bytes.Buffer
	s, c, i, allText := []byte(html), context{}, 0, true
	// Using the transition funcs helps us avoid mangling
	// `<div title="1>2">` or `I <3 Ponies!`.
	for i != len(s) {
		if c.delim == delimNone {
			st := c.state
			// Use RCDATA instead of parsing into JS or CSS styles.
			if c.element != elementNone && !isInTag(st) {
				st = stateRCDATA
			}
			d, nread := transitionFunc[st](c, s[i:])
			i1 := i + nread
			if c.state == stateText || c.state == stateRCDATA {
				// Emit text up to the start of the tag or comment.
				j := i1
				if d.state != c.state {
					for j1 := j - 1; j1 >= i; j1-- {
						if s[j1] == '<' {
							j = j1
							break
						}
					}
				}
				b.Write(s[i:j])
			} else {
				allText = false
			}
			c, i = d, i1
			continue
		}
		i1 := i + bytes.IndexAny(s[i:], delimEnds[c.delim])
		if i1 < i {
			break
		}
		if c.delim != delimSpaceOrTagEnd {
			// Consume any quote.
			i1++
		}
		c, i = context{state: stateTag, element: c.element}, i1
	}
	if allText {
		return html
	} else if c.state == stateText || c.state == stateRCDATA {
		b.Write(s[i:])
	}
	return b.String()
}

// htmlNameFilter accepts valid parts of an HTML attribute or tag name or
// a known-safe HTML attribute.
func htmlNameFilter(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeHTMLAttr {
		return s
	}
	if len(s) == 0 {
		// Avoid violation of structure preservation.
		// <input checked {{.K}}={{.V}}>.
		// Without this, if .K is empty then .V is the value of
		// checked, but otherwise .V is the value of the attribute
		// named .K.
		return filterFailsafe
	}
	s = strings.ToLower(s)
	if t := attrType(s); t != contentTypePlain {
		// TODO: Split attr and element name part filters so we can whitelist
		// attributes.
		return filterFailsafe
	}
	for _, r := range s {
		switch {
		case '0' <= r && r <= '9':
		case 'a' <= r && r <= 'z':
		default:
			return filterFailsafe
		}
	}
	return s
}

// commentEscaper returns the empty string regardless of input.
// Comment content does not correspond to any parsed structure or
// human-readable content, so the simplest and most secure policy is to drop
// content interpolated into comments.
// This approach is equally valid whether or not static comment content is
// removed from the template.
func commentEscaper(args ...interface{}) string {
	return ""
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// context describes the state an HTML parser must be in when it reaches the
// portion of HTML produced by evaluating a particular template node.
//
// The zero value of type context is the start context for a template that
// produces an HTML fragment as defined at
// http://www.w3.org/TR/html5/syntax.html#the-end
// where the context element is null.
type context struct {
	state   state
	delim   delim
	urlPart urlPart
	jsCtx   jsCtx
	attr    attr
	element element
	err     *Error
}

func (c context) String() string {
	return fmt.Sprintf("{%v %v %v %v %v %v %v}", c.state, c.delim, c.urlPart, c.jsCtx, c.attr, c.element, c.err)
}

// eq reports whether two contexts are equal.
func (c context) eq(d context) bool {
	return c.state == d.state &&
		c.delim == d.delim &&
		c.urlPart == d.urlPart &&
		c.jsCtx == d.jsCtx &&
		c.attr == d.attr &&
		c.element == d.element &&
		c.err == d.err
}

// mangle produces an identifier that includes a suffix that distinguishes it
// from template names mangled with different contexts.
func (c context) mangle(templateName string) string {
	// The mangled name for the default context is the input templateName.
	if c.state == stateText {
		return templateName
	}
	s := templateName + "$htmltemplate_" + c.state.String()
	if c.delim != 0 {
		s += "_" + c.delim.String()
	}
	if c.urlPart != 0 {
		s += "_" + c.urlPart.String()
	}
	if c.jsCtx != 0 {
		s += "_" + c.jsCtx.String()
	}
	if c.attr != 0 {
		s += "_" + c.attr.String()
	}
	if c.element != 0 {
		s += "_" + c.element.String()
	}
	return s
}

// state describes a high-level HTML parser state.
//
// It bounds the top of the element stack, and by extension the HTML insertion
// mode, but also contains state that does not correspond to anything in the
// HTML5 parsing algorithm because a single token production in the HTML
// grammar may contain embedded actions in a template. For instance, the quoted
// HTML attribute produced by
//     <div title="Hello {{.World}}">
// is a single token in HTML's grammar but in a template spans several nodes.
type state uint8

const (
	// stateText is parsed character data. An HTML parser is in
	// this state when its parse position is outside an HTML tag,
	// directive, comment, and special element body.
	stateText state = iota
	// stateTag occurs before an HTML attribute or the end of a tag.
	stateTag
	// stateAttrName occurs inside an attribute name.
	// It occurs between the ^'s in ` ^name^ = value`.
	stateAttrName
	// stateAfterName occurs after an attr name has ended but before any
	// equals sign. It occurs between the ^'s in ` name^ ^= value`.
	stateAfterName
	// stateBeforeValue occurs after the equals sign but before the value.
	// It occurs between the ^'s in ` name =^ ^value`.
	stateBeforeValue
	// stateHTMLCmt occurs inside an <!-- HTML comment -->.
	stateHTMLCmt
	// stateRCDATA occurs inside an RCDATA element (<textarea> or <title>)
	// as described at http://www.w3.org/TR/html5/syntax.html#elements-0
	stateRCDATA
	// stateAttr occurs inside an HTML attribute whose content is text.
	stateAttr
	// stateURL occurs inside an HTML attribute whose content is a URL.
	stateURL
	// stateJS occurs inside an event handler or script element.
	stateJS
	// stateJSDqStr occurs inside a JavaScript double quoted string.
	stateJSDqStr
	// stateJSSqStr occurs inside a JavaScript single quoted string.
	stateJSSqStr
	// stateJSRegexp occurs inside a JavaScript regexp literal.
	stateJSRegexp
	// stateJSBlockCmt occurs inside a JavaScript /* block comment */.
	stateJSBlockCmt
	// stateJSLineCmt occurs inside a JavaScript // line comment.
	stateJSLineCmt
	// stateCSS occurs inside a <style> element or style attribute.
	stateCSS
	// stateCSSDqStr occurs inside a CSS double quoted string.
	stateCSSDqStr
	// stateCSSSqStr occurs inside a CSS single quoted string.
	stateCSSSqStr
	// stateCSSDqURL occurs inside a CSS double quoted url("...").
	stateCSSDqURL
	// stateCSSSqURL occurs inside a CSS single quoted url('...').
	stateCSSSqURL
	// stateCSSURL occurs inside a CSS unquoted url(...).
	stateCSSURL
	// stateCSSBlockCmt occurs inside a CSS /* block comment */.
	stateCSSBlockCmt
	// stateCSSLineCmt occurs inside a CSS // line comment.
	stateCSSLineCmt
	// stateError is an infectious error state outside any valid
	// HTML/CSS/JS construct.
	stateError
)

var stateNames = [...]string{
	stateText:        "stateText",
	stateTag:         "stateTag",
	stateAttrName:    "stateAttrName",
	stateAfterName:   "stateAfterName",
	stateBeforeValue: "stateBeforeValue",
	stateHTMLCmt:     "stateHTMLCmt",
	stateRCDATA:      "stateRCDATA",
	stateAttr:        "stateAttr",
	stateURL:         "stateURL",
	stateJS:          "stateJS",
	stateJSDqStr:     "stateJSDqStr",
	stateJSSqStr:     "stateJSSqStr",
	stateJSRegexp:    "stateJSRegexp",
	stateJSBlockCmt:  "stateJSBlockCmt",
	stateJSLineCmt:   "stateJSLineCmt",
	stateCSS:         "stateCSS",
	stateCSSDqStr:    "stateCSSDqStr",
	stateCSSSqStr:    "stateCSSSqStr",
	stateCSSDqURL:    "stateCSSDqURL",
	stateCSSSqURL:    "stateCSSSqURL",
	stateCSSURL:      "stateCSSURL",
	stateCSSBlockCmt: "stateCSSBlockCmt",
	stateCSSLineCmt:  "stateCSSLineCmt",
	stateError:       "stateError",
}

func (s state) String() string {
	if int(s) < len(stateNames) {
		return stateNames[s]
	}
	return fmt.Sprintf("illegal state %d", int(s))
}

// isComment is true for any state that contains content meant for template
// authors & maintainers, not for end-users or machines.
func isComment(s state) bool {
	switch s {
	case stateHTMLCmt, stateJSBlockCmt, stateJSLineCmt, stateCSSBlockCmt, stateCSSLineCmt:
		return true
	}
	return false
}

// isInTag return whether s occurs solely inside an HTML tag.
func isInTag(s state) bool {
	switch s {
	case stateTag, stateAttrName, stateAfterName, stateBeforeValue, stateAttr:
		return true
	}
	return false
}

// delim is the delimiter that will end the current HTML attribute.
type delim uint8

const (
	// delimNone occurs outside any attribute.
	delimNone delim = iota
	// delimDoubleQuote occurs when a double quote (") closes the attribute.
	delimDoubleQuote
	// delimSingleQuote occurs when a single quote (') closes the attribute.
	delimSingleQuote
	// delimSpaceOrTagEnd occurs when a space or right angle bracket (>)
	// closes the attribute.
	delimSpaceOrTagEnd
)

var delimNames = [...]string{
	delimNone:          "delimNone",
	delimDoubleQuote:   "delimDoubleQuote",
	delimSingleQuote:   "delimSingleQuote",
	delimSpaceOrTagEnd: "delimSpaceOrTagEnd",
}

func (d delim) String() string {
	if int(d) < len(delimNames) {
		return delimNames[d]
	}
	return fmt.Sprintf("illegal delim %d", int(d))
}

// urlPart identifies a part in an RFC 3986 hierarchical URL to allow different
// encoding strategies.
type urlPart uint8

const (
	// urlPartNone occurs when not in a URL, or possibly at the start:
	// ^ in "^http://auth/path?k=v#frag".
	urlPartNone urlPart = iota
	// urlPartPreQuery occurs in the scheme, authority, or path; between the
	// ^s in "h^ttp://auth/path^?k=v#frag".
	urlPartPreQuery
	// urlPartQueryOrFrag occurs in the query portion between the ^s in
	// "http://auth/path?^k=v#frag^".
	urlPartQueryOrFrag
	// urlPartUnknown occurs due to joining of contexts both before and
	// after the query separator.
	urlPartUnknown
)

var urlPartNames = [...]string{
	urlPartNone:        "urlPartNone",
	urlPartPreQuery:    "urlPartPreQuery",
	urlPartQueryOrFrag: "urlPartQueryOrFrag",
	urlPartUnknown:     "urlPartUnknown",
}

func (u urlPart) String() string {
	if int(u) < len(urlPartNames) {
		return urlPartNames[u]
	}
	return fmt.Sprintf("illegal urlPart %d", int(u))
}

// jsCtx determines whether a '/' starts a regular expression literal or a
// division operator.
type jsCtx uint8

const (
	// jsCtxRegexp occurs where a '/' would start a regexp literal.
	jsCtxRegexp jsCtx = iota
	// jsCtxDivOp occurs where a '/' would start a division operator.
	jsCtxDivOp
	// jsCtxUnknown occurs where a '/' is ambiguous due to context joining.
	jsCtxUnknown
)

func (c jsCtx) String() string {
	switch c {
	case jsCtxRegexp:
		return "jsCtxRegexp"
	case jsCtxDivOp:
		return "jsCtxDivOp"
	case jsCtxUnknown:
		return "jsCtxUnknown"
	}
	return fmt.Sprintf("illegal jsCtx %d", int(c))
}

// element identifies the HTML element when inside a start tag or special body.
// Certain HTML element (for example <script> and <style>) have bodies that are
// treated differently from stateText so the element type is necessary to
// transition into the correct context at the end of a tag and to identify the
// end delimiter for the body.
type element uint8

const (
	// elementNone occurs outside a special tag or special element body.
	elementNone element = iota
	// elementScript corresponds to the raw text <script> element.
	elementScript
	// elementStyle corresponds to the raw text <style> element.
	elementStyle
	// elementTextarea corresponds to the RCDATA <textarea> element.
	elementTextarea
	// elementTitle corresponds to the RCDATA <title> element.
	elementTitle
)

var elementNames = [...]string{
	elementNone:     "elementNone",
	elementScript:   "elementScript",
	elementStyle:    "elementStyle",
	elementTextarea: "elementTextarea",
	elementTitle:    "elementTitle",
}

func (e element) String() string {
	if int(e) < len(elementNames) {
		return elementNames[e]
	}
	return fmt.Sprintf("illegal element %d", int(e))
}

// attr identifies the most recent HTML attribute when inside a start tag.
type attr uint8

const (
	// attrNone corresponds to a normal attribute or no attribute.
	attrNone attr = iota
	// attrScript corresponds to an event handler attribute.
	attrScript
	// attrStyle corresponds to the style attribute whose value is CSS.
	attrStyle
	// attrURL corresponds to an attribute whose value is a URL.
	attrURL
)

var attrNames = [...]string{
	attrNone:   "attrNone",
	attrScript: "attrScript",
	attrStyle:  "attrStyle",
	attrURL:    "attrURL",
}

func (a attr) String() string {
	if int(a) < len(attrNames) {
		return attrNames[a]
	}
	return fmt.Sprintf("illegal attr %d", int(a))
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// escapeTemplates rewrites the named templates, which must be
// associated with t, to guarantee that the output of any of the named
// templates is properly escaped.  Names should include the names of
// all templates that might be Executed but need not include helper
// templates.  If no error is returned, then the named templates have
// been modified.  Otherwise the named templates have been rendered
// unusable.
func escapeTemplates(tmpl *Template, names ...string) error {
	e := newEscaper(tmpl)
	for _, name := range names {
		c, _ := e.escapeTree(context{}, name, 0)
		var err error
		if c.err != nil {
			err, c.err.Name = c.err, name
		} else if c.state != stateText {
			err = &Error{ErrEndContext, name, 0, fmt.Sprintf("ends in a non-text context: %v", c)}
		}
		if err != nil {
			// Prevent execution of unsafe templates.
			for _, name := range names {
				if t := tmpl.set[name]; t != nil {
					t.text.Tree = nil
					t.Tree = nil
				}
			}
			return err
		}
	}
	e.commit()
	for _, name := range names {
		if t := tmpl.set[name]; t != nil {
			t.escaped = true
			t.Tree = t.text.Tree
		}
	}
	return nil
}

// funcMap maps command names to functions that render their inputs safe.
var funcMap = template.FuncMap{
	"html_template_attrescaper":     attrEscaper,
	"html_template_commentescaper":  commentEscaper,
	"html_template_cssescaper":      cssEscaper,
	"html_template_cssvaluefilter":  cssValueFilter,
	"html_template_htmlnamefilter":  htmlNameFilter,
	"html_template_htmlescaper":     htmlEscaper,
	"html_template_jsregexpescaper": jsRegexpEscaper,
	"html_template_jsstrescaper":    jsStrEscaper,
	"html_template_jsvalescaper":    jsValEscaper,
	"html_template_nospaceescaper":  htmlNospaceEscaper,
	"html_template_rcdataescaper":   rcdataEscaper,
	"html_template_urlescaper":      urlEscaper,
	"html_template_urlfilter":       urlFilter,
	"html_template_urlnormalizer":   urlNormalizer,
}

// equivEscapers matches contextual escapers to equivalent template builtins.
var equivEscapers = map[string]string{
	"html_template_attrescaper":    "html",
	"html_template_htmlescaper":    "html",
	"html_template_nospaceescaper": "html",
	"html_template_rcdataescaper":  "html",
	"html_template_urlescaper":     "urlquery",
	"html_template_urlnormalizer":  "urlquery",
}

// escaper collects type inferences about templates and changes needed to make
// templates injection safe.
type escaper struct {
	tmpl *Template
	// output[templateName] is the output context for a templateName that
	// has been mangled to include its input context.
	output map[string]context
	// derived[c.mangle(name)] maps to a template derived from the template
	// named name templateName for the start context c.
	derived map[string]*template.Template
	// called[templateName] is a set of called mangled template names.
	called map[string]bool
	// xxxNodeEdits are the accumulated edits to apply during commit.
	// Such edits are not applied immediately in case a template set
	// executes a given template in different escaping contexts.
	actionNodeEdits   map[*parse.ActionNode][]string
	templateNodeEdits map[*parse.TemplateNode]string
	textNodeEdits     map[*parse.TextNode][]byte
}

// newEscaper creates a blank escaper for the given set.
func newEscaper(t *Template) *escaper {
	return &escaper{
		t,
		map[string]context{},
		map[string]*template.Template{},
		map[string]bool{},
		map[*parse.ActionNode][]string{},
		map[*parse.TemplateNode]string{},
		map[*parse.TextNode][]byte{},
	}
}

// filterFailsafe is an innocuous word that is emitted in place of unsafe values
// by sanitizer functions. It is not a keyword in any programming language,
// contains no special characters, is not empty, and when it appears in output
// it is distinct enough that a developer can find the source of the problem
// via a search engine.
const filterFailsafe = "ZgotmplZ"

// escape escapes a template node.
func (e *escaper) escape(c context, n parse.Node) context {
	switch n := n.(type) {
	case *parse.ActionNode:
		return e.escapeAction(c, n)
	case *parse.IfNode:
		return e.escapeBranch(c, &n.BranchNode, "if")
	case *parse.ListNode:
		return e.escapeList(c, n)
	case *parse.RangeNode:
		return e.escapeBranch(c, &n.BranchNode, "range")
	case *parse.TemplateNode:
		return e.escapeTemplate(c, n)
	case *parse.TextNode:
		return e.escapeText(c, n)
	case *parse.WithNode:
		return e.escapeBranch(c, &n.BranchNode, "with")
	}
	panic("escaping " + n.String() + " is unimplemented")
}

// escapeAction escapes an action template node.
func (e *escaper) escapeAction(c context, n *parse.ActionNode) context {
	if len(n.Pipe.Decl) != 0 {
		// A local variable assignment, not an interpolation.
		return c
	}
	c = nudge(c)
	s := make([]string, 0, 3)
	switch c.state {
	case stateError:
		return c
	case stateURL, stateCSSDqStr, stateCSSSqStr, stateCSSDqURL, stateCSSSqURL, stateCSSURL:
		switch c.urlPart {
		case urlPartNone:
			s = append(s, "html_template_urlfilter")
			fallthrough
		case urlPartPreQuery:
			switch c.state {
			case stateCSSDqStr, stateCSSSqStr:
				s = append(s, "html_template_cssescaper")
			default:
				s = append(s, "html_template_urlnormalizer")
			}
		case urlPartQueryOrFrag:
			s = append(s, "html_template_urlescaper")
		case urlPartUnknown:
			return context{
				state: stateError,
				err:   errorf(ErrAmbigContext, n.Line, "%s appears in an ambiguous URL context", n),
			}
		default:
			panic(c.urlPart.String())
		}
	case stateJS:
		s = append(s, "html_template_jsvalescaper")
		// A slash after a value starts a div operator.
		c.jsCtx = jsCtxDivOp
	case stateJSDqStr, stateJSSqStr:
		s = append(s, "html_template_jsstrescaper")
	case stateJSRegexp:
		s = append(s, "html_template_jsregexpescaper")
	case stateCSS:
		s = append(s, "html_template_cssvaluefilter")
	case stateText:
		s = append(s, "html_template_htmlescaper")
	case stateRCDATA:
		s = append(s, "html_template_rcdataescaper")
	case stateAttr:
		// Handled below in delim check.
	case stateAttrName, stateTag:
		c.state = stateAttrName
		s = append(s, "html_template_htmlnamefilter")
	default:
		if isComment(c.state) {
			s = append(s, "html_template_commentescaper")
		} else {
			panic("unexpected state " + c.state.String())
		}
	}
	switch c.delim {
	case delimNone:
		// No extra-escaping needed for raw text content.
	case delimSpaceOrTagEnd:
		s = append(s, "html_template_nospaceescaper")
	default:
		s = append(s, "html_template_attrescaper")
	}
	e.editActionNode(n, s)
	return c
}

// allIdents returns the names of the identifiers under the Ident field of the node,
// which might be a singleton (Identifier) or a slice (Field).
func allIdents(node parse.Node) []string {
	switch node := node.(type) {
	case *parse.IdentifierNode:
		return []string{node.Ident}
	case *parse.FieldNode:
		return node.Ident
	}
	panic("unidentified node type in allIdents")
}

// ensurePipelineContains ensures that the pipeline has commands with
// the identifiers in s in order.
// If the pipeline already has some of the sanitizers, do not interfere.
// For example, if p is (.X | html) and s is ["escapeJSVal", "html"] then it
// has one matching, "html", and one to insert, "escapeJSVal", to produce
// (.X | escapeJSVal | html).
func ensurePipelineContains(p *parse.PipeNode, s []string) {
	if len(s) == 0 {
		return
	}
	n := len(p.Cmds)
	// Find the identifiers at the end of the command chain.
	idents := p.Cmds
	for i := n - 1; i >= 0; i-- {
		if cmd := p.Cmds[i]; len(cmd.Args) != 0 {
			if _, ok := cmd.Args[0].(*parse.IdentifierNode); ok {
				continue
			}
		}
		idents = p.Cmds[i+1:]
	}
	dups := 0
	for _, idNode := range idents {
		for _, ident := range allIdents(idNode.Args[0]) {
			if escFnsEq(s[dups], ident) {
				dups++
				if dups == len(s) {
					return
				}
			}
		}
	}
	newCmds := make([]*parse.CommandNode, n-len(idents), n+len(s)-dups)
	copy(newCmds, p.Cmds)
	// Merge existing identifier commands with the sanitizers needed.
	for _, idNode := range idents {
		pos := idNode.Args[0].Position()
		for _, ident := range allIdents(idNode.Args[0]) {
			i := indexOfStr(ident, s, escFnsEq)
			if i != -1 {
				for _, name := range s[:i] {
					newCmds = appendCmd(newCmds, newIdentCmd(name, pos))
				}
				s = s[i+1:]
			}
		}
		newCmds = appendCmd(newCmds, idNode)
	}
	// Create any remaining sanitizers.
	for _, name := range s {
		newCmds = appendCmd(newCmds, newIdentCmd(name, p.Position()))
	}
	p.Cmds = newCmds
}

// redundantFuncs[a][b] implies that funcMap[b](funcMap[a](x)) == funcMap[a](x)
// for all x.
var redundantFuncs = map[string]map[string]bool{
	"html_template_commentescaper": {
		"html_template_attrescaper":    true,
		"html_template_nospaceescaper": true,
		"html_template_htmlescaper":    true,
	},
	"html_template_cssescaper": {
		"html_template_attrescaper": true,
	},
	"html_template_jsregexpescaper": {
		"html_template_attrescaper": true,
	},
	"html_template_jsstrescaper": {
		"html_template_attrescaper": true,
	},
	"html_template_urlescaper": {
		"html_template_urlnormalizer": true,
	},
}

// appendCmd appends the given command to the end of the command pipeline
// unless it is redundant with the last command.
func appendCmd(cmds []*parse.CommandNode, cmd *parse.CommandNode) []*parse.CommandNode {
	if n := len(cmds); n != 0 {
		last, ok := cmds[n-1].Args[0].(*parse.IdentifierNode)
		next, _ := cmd.Args[0].(*parse.IdentifierNode)
		if ok && redundantFuncs[last.Ident][next.Ident] {
			return cmds
		}
	}
	return append(cmds, cmd)
}

// indexOfStr is the first i such that eq(s, strs[i]) or -1 if s was not found.
func indexOfStr(s string, strs []string, eq func(a, b string) bool) int {
	for i, t := range strs {
		if eq(s, t) {
			return i
		}
	}
	return -1
}

// escFnsEq reports whether the two escaping functions are equivalent.
func escFnsEq(a, b string) bool {
	if e := equivEscapers[a]; e != "" {
		a = e
	}
	if e := equivEscapers[b]; e != "" {
		b = e
	}
	return a == b
}

// newIdentCmd produces a command containing a single identifier node.
func newIdentCmd(identifier string, pos parse.Pos) *parse.CommandNode {
	return &parse.CommandNode{
		NodeType: parse.NodeCommand,
		Args:     []parse.Node{parse.NewIdentifier(identifier).SetPos(pos)},
	}
}

// nudge returns the context that would result from following empty string
// transitions from the input context.
// For example, parsing:
//     `<a href=`
// will end in context{stateBeforeValue, attrURL}, but parsing one extra rune:
//     `<a href=x`
// will end in context{stateURL, delimSpaceOrTagEnd, ...}.
// There are two transitions that happen when the 'x' is seen:
// (1) Transition from a before-value state to a start-of-value state without
//     consuming any character.
// (2) Consume 'x' and transition past the first value character.
// In this case, nudging produces the context after (1) happens.
func nudge(c context) context {
	switch c.state {
	case stateTag:
		// In `<foo {{.}}`, the action should emit an attribute.
		c.state = stateAttrName
	case stateBeforeValue:
		// In `<foo bar={{.}}`, the action is an undelimited value.
		c.state, c.delim, c.attr = attrStartStates[c.attr], delimSpaceOrTagEnd, attrNone
	case stateAfterName:
		// In `<foo bar {{.}}`, the action is an attribute name.
		c.state, c.attr = stateAttrName, attrNone
	}
	return c
}

// join joins the two contexts of a branch template node. The result is an
// error context if either of the input contexts are error contexts, or if the
// the input contexts differ.
func join(a, b context, line int, nodeName string) context {
	if a.state == stateError {
		return a
	}
	if b.state == stateError {
		return b
	}
	if a.eq(b) {
		return a
	}

	c := a
	c.urlPart = b.urlPart
	if c.eq(b) {
		// The contexts differ only by urlPart.
		c.urlPart = urlPartUnknown
		return c
	}

	c = a
	c.jsCtx = b.jsCtx
	if c.eq(b) {
		// The contexts differ only by jsCtx.
		c.jsCtx = jsCtxUnknown
		return c
	}

	// Allow a nudged context to join with an unnudged one.
	// This means that
	//   <p title={{if .C}}{{.}}{{end}}
	// ends in an unquoted value state even though the else branch
	// ends in stateBeforeValue.
	if c, d := nudge(a), nudge(b); !(c.eq(a) && d.eq(b)) {
		if e := join(c, d, line, nodeName); e.state != stateError {
			return e
		}
	}

	return context{
		state: stateError,
		err:   errorf(ErrBranchEnd, line, "{{%s}} branches end in different contexts: %v, %v", nodeName, a, b),
	}
}

// escapeBranch escapes a branch template node: "if", "range" and "with".
func (e *escaper) escapeBranch(c context, n *parse.BranchNode, nodeName string) context {
	c0 := e.escapeList(c, n.List)
	if nodeName == "range" && c0.state != stateError {
		// The "true" branch of a "range" node can execute multiple times.
		// We check that executing n.List once results in the same context
		// as executing n.List twice.
		c1, _ := e.escapeListConditionally(c0, n.List, nil)
		c0 = join(c0, c1, n.Line, nodeName)
		if c0.state == stateError {
			// Make clear that this is a problem on loop re-entry
			// since developers tend to overlook that branch when
			// debugging templates.
			c0.err.Line = n.Line
			c0.err.Description = "on range loop re-entry: " + c0.err.Description
			return c0
		}
	}
	c1 := e.escapeList(c, n.ElseList)
	return join(c0, c1, n.Line, nodeName)
}

// escapeList escapes a list template node.
func (e *escaper) escapeList(c context, n *parse.ListNode) context {
	if n == nil {
		return c
	}
	for _, m := range n.Nodes {
		c = e.escape(c, m)
	}
	return c
}

// escapeListConditionally escapes a list node but only preserves edits and
// inferences in e if the inferences and output context satisfy filter.
// It returns the best guess at an output context, and the result of the filter
// which is the same as whether e was updated.
func (e *escaper) escapeListConditionally(c context, n *parse.ListNode, filter func(*escaper, context) bool) (context, bool) {
	e1 := newEscaper(e.tmpl)
	// Make type inferences available to f.
	for k, v := range e.output {
		e1.output[k] = v
	}
	c = e1.escapeList(c, n)
	ok := filter != nil && filter(e1, c)
	if ok {
		// Copy inferences and edits from e1 back into e.
		for k, v := range e1.output {
			e.output[k] = v
		}
		for k, v := range e1.derived {
			e.derived[k] = v
		}
		for k, v := range e1.called {
			e.called[k] = v
		}
		for k, v := range e1.actionNodeEdits {
			e.editActionNode(k, v)
		}
		for k, v := range e1.templateNodeEdits {
			e.editTemplateNode(k, v)
		}
		for k, v := range e1.textNodeEdits {
			e.editTextNode(k, v)
		}
	}
	return c, ok
}

// escapeTemplate escapes a {{template}} call node.
func (e *escaper) escapeTemplate(c context, n *parse.TemplateNode) context {
	c, name := e.escapeTree(c, n.Name, n.Line)
	if name != n.Name {
		e.editTemplateNode(n, name)
	}
	return c
}

// escapeTree escapes the named template starting in the given context as
// necessary and returns its output context.
func (e *escaper) escapeTree(c context, name string, line int) (context, string) {
	// Mangle the template name with the input context to produce a reliable
	// identifier.
	dname := c.mangle(name)
	e.called[dname] = true
	if out, ok := e.output[dname]; ok {
		// Already escaped.
		return out, dname
	}
	t := e.template(name)
	if t == nil {
		// Two cases: The template exists but is empty, or has never been mentioned at
		// all. Distinguish the cases in the error messages.
		if e.tmpl.set[name] != nil {
			return context{
				state: stateError,
				err:   errorf(ErrNoSuchTemplate, line, "%q is an incomplete or empty template", name),
			}, dname
		}
		return context{
			state: stateError,
			err:   errorf(ErrNoSuchTemplate, line, "no such template %q", name),
		}, dname
	}
	if dname != name {
		// Use any template derived during an earlier call to escapeTemplate
		// with different top level templates, or clone if necessary.
		dt := e.template(dname)
		if dt == nil {
			dt = template.New(dname)
			dt.Tree = &parse.Tree{Name: dname, Root: t.Root.CopyList()}
			e.derived[dname] = dt
		}
		t = dt
	}
	return e.computeOutCtx(c, t), dname
}

// computeOutCtx takes a template and its start context and computes the output
// context while storing any inferences in e.
func (e *escaper) computeOutCtx(c context, t *template.Template) context {
	// Propagate context over the body.
	c1, ok := e.escapeTemplateBody(c, t)
	if !ok {
		// Look for a fixed point by assuming c1 as the output context.
		if c2, ok2 := e.escapeTemplateBody(c1, t); ok2 {
			c1, ok = c2, true
		}
		// Use c1 as the error context if neither assumption worked.
	}
	if !ok && c1.state != stateError {
		return context{
			state: stateError,
			// TODO: Find the first node with a line in t.text.Tree.Root
			err: errorf(ErrOutputContext, 0, "cannot compute output context for template %s", t.Name()),
		}
	}
	return c1
}

// escapeTemplateBody escapes the given template assuming the given output
// context, and returns the best guess at the output context and whether the
// assumption was correct.
func (e *escaper) escapeTemplateBody(c context, t *template.Template) (context, bool) {
	filter := func(e1 *escaper, c1 context) bool {
		if c1.state == stateError {
			// Do not update the input escaper, e.
			return false
		}
		if !e1.called[t.Name()] {
			// If t is not recursively called, then c1 is an
			// accurate output context.
			return true
		}
		// c1 is accurate if it matches our assumed output context.
		return c.eq(c1)
	}
	// We need to assume an output context so that recursive template calls
	// take the fast path out of escapeTree instead of infinitely recursing.
	// Naively assuming that the input context is the same as the output
	// works >90% of the time.
	e.output[t.Name()] = c
	return e.escapeListConditionally(c, t.Tree.Root, filter)
}

// delimEnds maps each delim to a string of characters that terminate it.
var delimEnds = [...]string{
	delimDoubleQuote: `"`,
	delimSingleQuote: "'",
	// Determined empirically by running the below in various browsers.
	// var div = document.createElement("DIV");
	// for (var i = 0; i < 0x10000; ++i) {
	//   div.innerHTML = "<span title=x" + String.fromCharCode(i) + "-bar>";
	//   if (div.getElementsByTagName("SPAN")[0].title.indexOf("bar") < 0)
	//     document.write("<p>U+" + i.toString(16));
	// }
	delimSpaceOrTagEnd: " \t\n\f\r>",
}

var doctypeBytes = []byte("<!DOCTYPE")

// escapeText escapes a text template node.
func (e *escaper) escapeText(c context, n *parse.TextNode) context {
	s, written, i, b := n.Text, 0, 0, new(bytes.Buffer)
	for i != len(s) {
		c1, nread := contextAfterText(c, s[i:])
		i1 := i + nread
		if c.state == stateText || c.state == stateRCDATA {
			end := i1
			if c1.state != c.state {
				for j := end - 1; j >= i; j-- {
					if s[j] == '<' {
						end = j
						break
					}
				}
			}
			for j := i; j < end; j++ {
				if s[j] == '<' && !bytes.HasPrefix(bytes.ToUpper(s[j:]), doctypeBytes) {
					b.Write(s[written:j])
					b.WriteString("&lt;")
					written = j + 1
				}
			}
		} else if isComment(c.state) && c.delim == delimNone {
			switch c.state {
			case stateJSBlockCmt:
				// http://es5.github.com/#x7.4:
				// "Comments behave like white space and are
				// discarded except that, if a MultiLineComment
				// contains a line terminator character, then
				// the entire comment is considered to be a
				// LineTerminator for purposes of parsing by
				// the syntactic grammar."
				if bytes.IndexAny(s[written:i1], "\n\r\u2028\u2029") != -1 {
					b.WriteByte('\n')
				} else {
					b.WriteByte(' ')
				}
			case stateCSSBlockCmt:
				b.WriteByte(' ')
			}
			written = i1
		}
		if c.state != c1.state && isComment(c1.state) && c1.delim == delimNone {
			// Preserve the portion between written and the comment start.
			cs := i1 - 2
			if c1.state == stateHTMLCmt {
				// "<!--" instead of "/*" or "//"
				cs -= 2
			}
			b.Write(s[written:cs])
			written = i1
		}
		if i == i1 && c.state == c1.state {
			panic(fmt.Sprintf("infinite loop from %v to %v on %q..%q", c, c1, s[:i], s[i:]))
		}
		c, i = c1, i1
	}

	if written != 0 && c.state != stateError {
		if !isComment(c.state) || c.delim != delimNone {
			b.Write(n.Text[written:])
		}
		e.editTextNode(n, b.Bytes())
	}
	return c
}

// contextAfterText starts in context c, consumes some tokens from the front of
// s, then returns the context after those tokens and the unprocessed suffix.
func contextAfterText(c context, s []byte) (context, int) {
	if c.delim == delimNone {
		c1, i := tSpecialTagEnd(c, s)
		if i == 0 {
			// A special end tag (`</script>`) has been seen and
			// all content preceding it has been consumed.
			return c1, 0
		}
		// Consider all content up to any end tag.
		return transitionFunc[c.state](c, s[:i])
	}

	i := bytes.IndexAny(s, delimEnds[c.delim])
	if i == -1 {
		i = len(s)
	}
	if c.delim == delimSpaceOrTagEnd {
		// http://www.w3.org/TR/html5/syntax.html#attribute-value-(unquoted)-state
		// lists the runes below as error characters.
		// Error out because HTML parsers may differ on whether
		// "<a id= onclick=f("     ends inside id's or onclick's value,
		// "<a class=`foo "        ends inside a value,
		// "<a style=font:'Arial'" needs open-quote fixup.
		// IE treats '`' as a quotation character.
		if j := bytes.IndexAny(s[:i], "\"'<=`"); j >= 0 {
			return context{
				state: stateError,
				err:   errorf(ErrBadHTML, 0, "%q in unquoted attr: %q", s[j:j+1], s[:i]),
			}, len(s)
		}
	}
	if i == len(s) {
		// Remain inside the attribute.
		// Decode the value so non-HTML rules can easily handle
		//     <button onclick="alert(&quot;Hi!&quot;)">
		// without having to entity decode token boundaries.
		for u := []byte(html.UnescapeString(string(s))); len(u) != 0; {
			c1, i1 := transitionFunc[c.state](c, u)
			c, u = c1, u[i1:]
		}
		return c, len(s)
	}
	if c.delim != delimSpaceOrTagEnd {
		// Consume any quote.
		i++
	}
	// On exiting an attribute, we discard all state information
	// except the state and element.
	return context{state: stateTag, element: c.element}, i
}

// editActionNode records a change to an action pipeline for later commit.
func (e *escaper) editActionNode(n *parse.ActionNode, cmds []string) {
	if _, ok := e.actionNodeEdits[n]; ok {
		panic(fmt.Sprintf("node %s shared between templates", n))
	}
	e.actionNodeEdits[n] = cmds
}

// editTemplateNode records a change to a {{template}} callee for later commit.
func (e *escaper) editTemplateNode(n *parse.TemplateNode, callee string) {
	if _, ok := e.templateNodeEdits[n]; ok {
		panic(fmt.Sprintf("node %s shared between templates", n))
	}
	e.templateNodeEdits[n] = callee
}

// editTextNode records a change to a text node for later commit.
func (e *escaper) editTextNode(n *parse.TextNode, text []byte) {
	if _, ok := e.textNodeEdits[n]; ok {
		panic(fmt.Sprintf("node %s shared between templates", n))
	}
	e.textNodeEdits[n] = text
}

// commit applies changes to actions and template calls needed to contextually
// autoescape content and adds any derived templates to the set.
func (e *escaper) commit() {
	for name := range e.output {
		e.template(name).Funcs(funcMap)
	}
	for _, t := range e.derived {
		if _, err := e.tmpl.text.AddParseTree(t.Name(), t.Tree); err != nil {
			panic("error adding derived template")
		}
	}
	for n, s := range e.actionNodeEdits {
		ensurePipelineContains(n.Pipe, s)
	}
	for n, name := range e.templateNodeEdits {
		n.Name = name
	}
	for n, s := range e.textNodeEdits {
		n.Text = s
	}
}

// template returns the named template given a mangled template name.
func (e *escaper) template(name string) *template.Template {
	t := e.tmpl.text.Lookup(name)
	if t == nil {
		t = e.derived[name]
	}
	return t
}

// Forwarding functions so that clients need only import this package
// to reach the general escaping functions of text/template.

// HTMLEscape writes to w the escaped HTML equivalent of the plain text data b.
func HTMLEscape(w io.Writer, b []byte) {
	template.HTMLEscape(w, b)
}

// HTMLEscapeString returns the escaped HTML equivalent of the plain text data s.
func HTMLEscapeString(s string) string {
	return template.HTMLEscapeString(s)
}

// HTMLEscaper returns the escaped HTML equivalent of the textual
// representation of its arguments.
func HTMLEscaper(args ...interface{}) string {
	return template.HTMLEscaper(args...)
}

// JSEscape writes to w the escaped JavaScript equivalent of the plain text data b.
func JSEscape(w io.Writer, b []byte) {
	template.JSEscape(w, b)
}

// JSEscapeString returns the escaped JavaScript equivalent of the plain text data s.
func JSEscapeString(s string) string {
	return template.JSEscapeString(s)
}

// JSEscaper returns the escaped JavaScript equivalent of the textual
// representation of its arguments.
func JSEscaper(args ...interface{}) string {
	return template.JSEscaper(args...)
}

// URLQueryEscaper returns the escaped value of the textual representation of
// its arguments in a form suitable for embedding in a URL query.
func URLQueryEscaper(args ...interface{}) string {
	return template.URLQueryEscaper(args...)
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// transitionFunc is the array of context transition functions for text nodes.
// A transition function takes a context and template text input, and returns
// the updated context and the number of bytes consumed from the front of the
// input.
var transitionFunc = [...]func(context, []byte) (context, int){
	stateText:        tText,
	stateTag:         tTag,
	stateAttrName:    tAttrName,
	stateAfterName:   tAfterName,
	stateBeforeValue: tBeforeValue,
	stateHTMLCmt:     tHTMLCmt,
	stateRCDATA:      tSpecialTagEnd,
	stateAttr:        tAttr,
	stateURL:         tURL,
	stateJS:          tJS,
	stateJSDqStr:     tJSDelimited,
	stateJSSqStr:     tJSDelimited,
	stateJSRegexp:    tJSDelimited,
	stateJSBlockCmt:  tBlockCmt,
	stateJSLineCmt:   tLineCmt,
	stateCSS:         tCSS,
	stateCSSDqStr:    tCSSStr,
	stateCSSSqStr:    tCSSStr,
	stateCSSDqURL:    tCSSStr,
	stateCSSSqURL:    tCSSStr,
	stateCSSURL:      tCSSStr,
	stateCSSBlockCmt: tBlockCmt,
	stateCSSLineCmt:  tLineCmt,
	stateError:       tError,
}

var commentStart = []byte("<!--")
var commentEnd = []byte("-->")

// tText is the context transition function for the text state.
func tText(c context, s []byte) (context, int) {
	k := 0
	for {
		i := k + bytes.IndexByte(s[k:], '<')
		if i < k || i+1 == len(s) {
			return c, len(s)
		} else if i+4 <= len(s) && bytes.Equal(commentStart, s[i:i+4]) {
			return context{state: stateHTMLCmt}, i + 4
		}
		i++
		end := false
		if s[i] == '/' {
			if i+1 == len(s) {
				return c, len(s)
			}
			end, i = true, i+1
		}
		j, e := eatTagName(s, i)
		if j != i {
			if end {
				e = elementNone
			}
			// We've found an HTML tag.
			return context{state: stateTag, element: e}, j
		}
		k = j
	}
}

var elementContentType = [...]state{
	elementNone:     stateText,
	elementScript:   stateJS,
	elementStyle:    stateCSS,
	elementTextarea: stateRCDATA,
	elementTitle:    stateRCDATA,
}

// tTag is the context transition function for the tag state.
func tTag(c context, s []byte) (context, int) {
	// Find the attribute name.
	i := eatWhiteSpace(s, 0)
	if i == len(s) {
		return c, len(s)
	}
	if s[i] == '>' {
		return context{
			state:   elementContentType[c.element],
			element: c.element,
		}, i + 1
	}
	j, err := eatAttrName(s, i)
	if err != nil {
		return context{state: stateError, err: err}, len(s)
	}
	state, attr := stateTag, attrNone
	if i == j {
		return context{
			state: stateError,
			err:   errorf(ErrBadHTML, 0, "expected space, attr name, or end of tag, but got %q", s[i:]),
		}, len(s)
	}
	switch attrType(string(s[i:j])) {
	case contentTypeURL:
		attr = attrURL
	case contentTypeCSS:
		attr = attrStyle
	case contentTypeJS:
		attr = attrScript
	}
	if j == len(s) {
		state = stateAttrName
	} else {
		state = stateAfterName
	}
	return context{state: state, element: c.element, attr: attr}, j
}

// tAttrName is the context transition function for stateAttrName.
func tAttrName(c context, s []byte) (context, int) {
	i, err := eatAttrName(s, 0)
	if err != nil {
		return context{state: stateError, err: err}, len(s)
	} else if i != len(s) {
		c.state = stateAfterName
	}
	return c, i
}

// tAfterName is the context transition function for stateAfterName.
func tAfterName(c context, s []byte) (context, int) {
	// Look for the start of the value.
	i := eatWhiteSpace(s, 0)
	if i == len(s) {
		return c, len(s)
	} else if s[i] != '=' {
		// Occurs due to tag ending '>', and valueless attribute.
		c.state = stateTag
		return c, i
	}
	c.state = stateBeforeValue
	// Consume the "=".
	return c, i + 1
}

var attrStartStates = [...]state{
	attrNone:   stateAttr,
	attrScript: stateJS,
	attrStyle:  stateCSS,
	attrURL:    stateURL,
}

// tBeforeValue is the context transition function for stateBeforeValue.
func tBeforeValue(c context, s []byte) (context, int) {
	i := eatWhiteSpace(s, 0)
	if i == len(s) {
		return c, len(s)
	}
	// Find the attribute delimiter.
	delim := delimSpaceOrTagEnd
	switch s[i] {
	case '\'':
		delim, i = delimSingleQuote, i+1
	case '"':
		delim, i = delimDoubleQuote, i+1
	}
	c.state, c.delim, c.attr = attrStartStates[c.attr], delim, attrNone
	return c, i
}

// tHTMLCmt is the context transition function for stateHTMLCmt.
func tHTMLCmt(c context, s []byte) (context, int) {
	if i := bytes.Index(s, commentEnd); i != -1 {
		return context{}, i + 3
	}
	return c, len(s)
}

// specialTagEndMarkers maps element types to the character sequence that
// case-insensitively signals the end of the special tag body.
var specialTagEndMarkers = [...]string{
	elementScript:   "</script",
	elementStyle:    "</style",
	elementTextarea: "</textarea",
	elementTitle:    "</title",
}

// tSpecialTagEnd is the context transition function for raw text and RCDATA
// element states.
func tSpecialTagEnd(c context, s []byte) (context, int) {
	if c.element != elementNone {
		if i := strings.Index(strings.ToLower(string(s)), specialTagEndMarkers[c.element]); i != -1 {
			return context{}, i
		}
	}
	return c, len(s)
}

// tAttr is the context transition function for the attribute state.
func tAttr(c context, s []byte) (context, int) {
	return c, len(s)
}

// tURL is the context transition function for the URL state.
func tURL(c context, s []byte) (context, int) {
	if bytes.IndexAny(s, "#?") >= 0 {
		c.urlPart = urlPartQueryOrFrag
	} else if len(s) != eatWhiteSpace(s, 0) && c.urlPart == urlPartNone {
		// HTML5 uses "Valid URL potentially surrounded by spaces" for
		// attrs: http://www.w3.org/TR/html5/index.html#attributes-1
		c.urlPart = urlPartPreQuery
	}
	return c, len(s)
}

// tJS is the context transition function for the JS state.
func tJS(c context, s []byte) (context, int) {
	i := bytes.IndexAny(s, `"'/`)
	if i == -1 {
		// Entire input is non string, comment, regexp tokens.
		c.jsCtx = nextJSCtx(s, c.jsCtx)
		return c, len(s)
	}
	c.jsCtx = nextJSCtx(s[:i], c.jsCtx)
	switch s[i] {
	case '"':
		c.state, c.jsCtx = stateJSDqStr, jsCtxRegexp
	case '\'':
		c.state, c.jsCtx = stateJSSqStr, jsCtxRegexp
	case '/':
		switch {
		case i+1 < len(s) && s[i+1] == '/':
			c.state, i = stateJSLineCmt, i+1
		case i+1 < len(s) && s[i+1] == '*':
			c.state, i = stateJSBlockCmt, i+1
		case c.jsCtx == jsCtxRegexp:
			c.state = stateJSRegexp
		case c.jsCtx == jsCtxDivOp:
			c.jsCtx = jsCtxRegexp
		default:
			return context{
				state: stateError,
				err:   errorf(ErrSlashAmbig, 0, "'/' could start a division or regexp: %.32q", s[i:]),
			}, len(s)
		}
	default:
		panic("unreachable")
	}
	return c, i + 1
}

// tJSDelimited is the context transition function for the JS string and regexp
// states.
func tJSDelimited(c context, s []byte) (context, int) {
	specials := `\"`
	switch c.state {
	case stateJSSqStr:
		specials = `\'`
	case stateJSRegexp:
		specials = `\/[]`
	}

	k, inCharset := 0, false
	for {
		i := k + bytes.IndexAny(s[k:], specials)
		if i < k {
			break
		}
		switch s[i] {
		case '\\':
			i++
			if i == len(s) {
				return context{
					state: stateError,
					err:   errorf(ErrPartialEscape, 0, "unfinished escape sequence in JS string: %q", s),
				}, len(s)
			}
		case '[':
			inCharset = true
		case ']':
			inCharset = false
		default:
			// end delimiter
			if !inCharset {
				c.state, c.jsCtx = stateJS, jsCtxDivOp
				return c, i + 1
			}
		}
		k = i + 1
	}

	if inCharset {
		// This can be fixed by making context richer if interpolation
		// into charsets is desired.
		return context{
			state: stateError,
			err:   errorf(ErrPartialCharset, 0, "unfinished JS regexp charset: %q", s),
		}, len(s)
	}

	return c, len(s)
}

var blockCommentEnd = []byte("*/")

// tBlockCmt is the context transition function for /*comment*/ states.
func tBlockCmt(c context, s []byte) (context, int) {
	i := bytes.Index(s, blockCommentEnd)
	if i == -1 {
		return c, len(s)
	}
	switch c.state {
	case stateJSBlockCmt:
		c.state = stateJS
	case stateCSSBlockCmt:
		c.state = stateCSS
	default:
		panic(c.state.String())
	}
	return c, i + 2
}

// tLineCmt is the context transition function for //comment states.
func tLineCmt(c context, s []byte) (context, int) {
	var lineTerminators string
	var endState state
	switch c.state {
	case stateJSLineCmt:
		lineTerminators, endState = "\n\r\u2028\u2029", stateJS
	case stateCSSLineCmt:
		lineTerminators, endState = "\n\f\r", stateCSS
		// Line comments are not part of any published CSS standard but
		// are supported by the 4 major browsers.
		// This defines line comments as
		//     LINECOMMENT ::= "//" [^\n\f\d]*
		// since http://www.w3.org/TR/css3-syntax/#SUBTOK-nl defines
		// newlines:
		//     nl ::= #xA | #xD #xA | #xD | #xC
	default:
		panic(c.state.String())
	}

	i := bytes.IndexAny(s, lineTerminators)
	if i == -1 {
		return c, len(s)
	}
	c.state = endState
	// Per section 7.4 of EcmaScript 5 : http://es5.github.com/#x7.4
	// "However, the LineTerminator at the end of the line is not
	// considered to be part of the single-line comment; it is
	// recognized separately by the lexical grammar and becomes part
	// of the stream of input elements for the syntactic grammar."
	return c, i
}

// tCSS is the context transition function for the CSS state.
func tCSS(c context, s []byte) (context, int) {
	// CSS quoted strings are almost never used except for:
	// (1) URLs as in background: "/foo.png"
	// (2) Multiword font-names as in font-family: "Times New Roman"
	// (3) List separators in content values as in inline-lists:
	//    <style>
	//    ul.inlineList { list-style: none; padding:0 }
	//    ul.inlineList > li { display: inline }
	//    ul.inlineList > li:before { content: ", " }
	//    ul.inlineList > li:first-child:before { content: "" }
	//    </style>
	//    <ul class=inlineList><li>One<li>Two<li>Three</ul>
	// (4) Attribute value selectors as in a[href="http://example.com/"]
	//
	// We conservatively treat all strings as URLs, but make some
	// allowances to avoid confusion.
	//
	// In (1), our conservative assumption is justified.
	// In (2), valid font names do not contain ':', '?', or '#', so our
	// conservative assumption is fine since we will never transition past
	// urlPartPreQuery.
	// In (3), our protocol heuristic should not be tripped, and there
	// should not be non-space content after a '?' or '#', so as long as
	// we only %-encode RFC 3986 reserved characters we are ok.
	// In (4), we should URL escape for URL attributes, and for others we
	// have the attribute name available if our conservative assumption
	// proves problematic for real code.

	k := 0
	for {
		i := k + bytes.IndexAny(s[k:], `("'/`)
		if i < k {
			return c, len(s)
		}
		switch s[i] {
		case '(':
			// Look for url to the left.
			p := bytes.TrimRight(s[:i], "\t\n\f\r ")
			if endsWithCSSKeyword(p, "url") {
				j := len(s) - len(bytes.TrimLeft(s[i+1:], "\t\n\f\r "))
				switch {
				case j != len(s) && s[j] == '"':
					c.state, j = stateCSSDqURL, j+1
				case j != len(s) && s[j] == '\'':
					c.state, j = stateCSSSqURL, j+1
				default:
					c.state = stateCSSURL
				}
				return c, j
			}
		case '/':
			if i+1 < len(s) {
				switch s[i+1] {
				case '/':
					c.state = stateCSSLineCmt
					return c, i + 2
				case '*':
					c.state = stateCSSBlockCmt
					return c, i + 2
				}
			}
		case '"':
			c.state = stateCSSDqStr
			return c, i + 1
		case '\'':
			c.state = stateCSSSqStr
			return c, i + 1
		}
		k = i + 1
	}
}

// tCSSStr is the context transition function for the CSS string and URL states.
func tCSSStr(c context, s []byte) (context, int) {
	var endAndEsc string
	switch c.state {
	case stateCSSDqStr, stateCSSDqURL:
		endAndEsc = `\"`
	case stateCSSSqStr, stateCSSSqURL:
		endAndEsc = `\'`
	case stateCSSURL:
		// Unquoted URLs end with a newline or close parenthesis.
		// The below includes the wc (whitespace character) and nl.
		endAndEsc = "\\\t\n\f\r )"
	default:
		panic(c.state.String())
	}

	k := 0
	for {
		i := k + bytes.IndexAny(s[k:], endAndEsc)
		if i < k {
			c, nread := tURL(c, decodeCSS(s[k:]))
			return c, k + nread
		}
		if s[i] == '\\' {
			i++
			if i == len(s) {
				return context{
					state: stateError,
					err:   errorf(ErrPartialEscape, 0, "unfinished escape sequence in CSS string: %q", s),
				}, len(s)
			}
		} else {
			c.state = stateCSS
			return c, i + 1
		}
		c, _ = tURL(c, decodeCSS(s[:i+1]))
		k = i + 1
	}
}

// tError is the context transition function for the error state.
func tError(c context, s []byte) (context, int) {
	return c, len(s)
}

// eatAttrName returns the largest j such that s[i:j] is an attribute name.
// It returns an error if s[i:] does not look like it begins with an
// attribute name, such as encountering a quote mark without a preceding
// equals sign.
func eatAttrName(s []byte, i int) (int, *Error) {
	for j := i; j < len(s); j++ {
		switch s[j] {
		case ' ', '\t', '\n', '\f', '\r', '=', '>':
			return j, nil
		case '\'', '"', '<':
			// These result in a parse warning in HTML5 and are
			// indicative of serious problems if seen in an attr
			// name in a template.
			return -1, errorf(ErrBadHTML, 0, "%q in attribute name: %.32q", s[j:j+1], s)
		default:
			// No-op.
		}
	}
	return len(s), nil
}

var elementNameMap = map[string]element{
	"script":   elementScript,
	"style":    elementStyle,
	"textarea": elementTextarea,
	"title":    elementTitle,
}

// asciiAlpha reports whether c is an ASCII letter.
func asciiAlpha(c byte) bool {
	return 'A' <= c && c <= 'Z' || 'a' <= c && c <= 'z'
}

// asciiAlphaNum reports whether c is an ASCII letter or digit.
func asciiAlphaNum(c byte) bool {
	return asciiAlpha(c) || '0' <= c && c <= '9'
}

// eatTagName returns the largest j such that s[i:j] is a tag name and the tag type.
func eatTagName(s []byte, i int) (int, element) {
	if i == len(s) || !asciiAlpha(s[i]) {
		return i, elementNone
	}
	j := i + 1
	for j < len(s) {
		x := s[j]
		if asciiAlphaNum(x) {
			j++
			continue
		}
		// Allow "x-y" or "x:y" but not "x-", "-y", or "x--y".
		if (x == ':' || x == '-') && j+1 < len(s) && asciiAlphaNum(s[j+1]) {
			j += 2
			continue
		}
		break
	}
	return j, elementNameMap[strings.ToLower(string(s[i:j]))]
}

// eatWhiteSpace returns the largest j such that s[i:j] is white space.
func eatWhiteSpace(s []byte, i int) int {
	for j := i; j < len(s); j++ {
		switch s[j] {
		case ' ', '\t', '\n', '\f', '\r':
			// No-op.
		default:
			return j
		}
	}
	return len(s)
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Error describes a problem encountered during template Escaping.
type Error struct {
	// ErrorCode describes the kind of error.
	ErrorCode ErrorCode
	// Name is the name of the template in which the error was encountered.
	Name string
	// Line is the line number of the error in the template source or 0.
	Line int
	// Description is a human-readable description of the problem.
	Description string
}

// ErrorCode is a code for a kind of error.
type ErrorCode int

// We define codes for each error that manifests while escaping templates, but
// escaped templates may also fail at runtime.
//
// Output: "ZgotmplZ"
// Example:
//   <img src="{{.X}}">
//   where {{.X}} evaluates to `javascript:...`
// Discussion:
//   "ZgotmplZ" is a special value that indicates that unsafe content reached a
//   CSS or URL context at runtime. The output of the example will be
//     <img src="#ZgotmplZ">
//   If the data comes from a trusted source, use content types to exempt it
//   from filtering: URL(`javascript:...`).
const (
	// OK indicates the lack of an error.
	OK ErrorCode = iota

	// ErrAmbigContext: "... appears in an ambiguous URL context"
	// Example:
	//   <a href="
	//      {{if .C}}
	//        /path/
	//      {{else}}
	//        /search?q=
	//      {{end}}
	//      {{.X}}
	//   ">
	// Discussion:
	//   {{.X}} is in an ambiguous URL context since, depending on {{.C}},
	//  it may be either a URL suffix or a query parameter.
	//   Moving {{.X}} into the condition removes the ambiguity:
	//   <a href="{{if .C}}/path/{{.X}}{{else}}/search?q={{.X}}">
	ErrAmbigContext

	// ErrBadHTML: "expected space, attr name, or end of tag, but got ...",
	//   "... in unquoted attr", "... in attribute name"
	// Example:
	//   <a href = /search?q=foo>
	//   <href=foo>
	//   <form na<e=...>
	//   <option selected<
	// Discussion:
	//   This is often due to a typo in an HTML element, but some runes
	//   are banned in tag names, attribute names, and unquoted attribute
	//   values because they can tickle parser ambiguities.
	//   Quoting all attributes is the best policy.
	ErrBadHTML

	// ErrBranchEnd: "{{if}} branches end in different contexts"
	// Example:
	//   {{if .C}}<a href="{{end}}{{.X}}
	// Discussion:
	//   Package html/template statically examines each path through an
	//   {{if}}, {{range}}, or {{with}} to escape any following pipelines.
	//   The example is ambiguous since {{.X}} might be an HTML text node,
	//   or a URL prefix in an HTML attribute. The context of {{.X}} is
	//   used to figure out how to escape it, but that context depends on
	//   the run-time value of {{.C}} which is not statically known.
	//
	//   The problem is usually something like missing quotes or angle
	//   brackets, or can be avoided by refactoring to put the two contexts
	//   into different branches of an if, range or with. If the problem
	//   is in a {{range}} over a collection that should never be empty,
	//   adding a dummy {{else}} can help.
	ErrBranchEnd

	// ErrEndContext: "... ends in a non-text context: ..."
	// Examples:
	//   <div
	//   <div title="no close quote>
	//   <script>f()
	// Discussion:
	//   Executed templates should produce a DocumentFragment of HTML.
	//   Templates that end without closing tags will trigger this error.
	//   Templates that should not be used in an HTML context or that
	//   produce incomplete Fragments should not be executed directly.
	//
	//   {{define "main"}} <script>{{template "helper"}}</script> {{end}}
	//   {{define "helper"}} document.write(' <div title=" ') {{end}}
	//
	//   "helper" does not produce a valid document fragment, so should
	//   not be Executed directly.
	ErrEndContext

	// ErrNoSuchTemplate: "no such template ..."
	// Examples:
	//   {{define "main"}}<div {{template "attrs"}}>{{end}}
	//   {{define "attrs"}}href="{{.URL}}"{{end}}
	// Discussion:
	//   Package html/template looks through template calls to compute the
	//   context.
	//   Here the {{.URL}} in "attrs" must be treated as a URL when called
	//   from "main", but you will get this error if "attrs" is not defined
	//   when "main" is parsed.
	ErrNoSuchTemplate

	// ErrOutputContext: "cannot compute output context for template ..."
	// Examples:
	//   {{define "t"}}{{if .T}}{{template "t" .T}}{{end}}{{.H}}",{{end}}
	// Discussion:
	//   A recursive template does not end in the same context in which it
	//   starts, and a reliable output context cannot be computed.
	//   Look for typos in the named template.
	//   If the template should not be called in the named start context,
	//   look for calls to that template in unexpected contexts.
	//   Maybe refactor recursive templates to not be recursive.
	ErrOutputContext

	// ErrPartialCharset: "unfinished JS regexp charset in ..."
	// Example:
	//     <script>var pattern = /foo[{{.Chars}}]/</script>
	// Discussion:
	//   Package html/template does not support interpolation into regular
	//   expression literal character sets.
	ErrPartialCharset

	// ErrPartialEscape: "unfinished escape sequence in ..."
	// Example:
	//   <script>alert("\{{.X}}")</script>
	// Discussion:
	//   Package html/template does not support actions following a
	//   backslash.
	//   This is usually an error and there are better solutions; for
	//   example
	//     <script>alert("{{.X}}")</script>
	//   should work, and if {{.X}} is a partial escape sequence such as
	//   "xA0", mark the whole sequence as safe content: JSStr(`\xA0`)
	ErrPartialEscape

	// ErrRangeLoopReentry: "on range loop re-entry: ..."
	// Example:
	//   <script>var x = [{{range .}}'{{.}},{{end}}]</script>
	// Discussion:
	//   If an iteration through a range would cause it to end in a
	//   different context than an earlier pass, there is no single context.
	//   In the example, there is missing a quote, so it is not clear
	//   whether {{.}} is meant to be inside a JS string or in a JS value
	//   context.  The second iteration would produce something like
	//
	//     <script>var x = ['firstValue,'secondValue]</script>
	ErrRangeLoopReentry

	// ErrSlashAmbig: '/' could start a division or regexp.
	// Example:
	//   <script>
	//     {{if .C}}var x = 1{{end}}
	//     /-{{.N}}/i.test(x) ? doThis : doThat();
	//   </script>
	// Discussion:
	//   The example above could produce `var x = 1/-2/i.test(s)...`
	//   in which the first '/' is a mathematical division operator or it
	//   could produce `/-2/i.test(s)` in which the first '/' starts a
	//   regexp literal.
	//   Look for missing semicolons inside branches, and maybe add
	//   parentheses to make it clear which interpretation you intend.
	ErrSlashAmbig
)

func (e *Error) Error() string {
	if e.Line != 0 {
		return fmt.Sprintf("html/template:%s:%d: %s", e.Name, e.Line, e.Description)
	} else if e.Name != "" {
		return fmt.Sprintf("html/template:%s: %s", e.Name, e.Description)
	}
	return "html/template: " + e.Description
}

// errorf creates an error given a format string f and args.
// The template Name still needs to be supplied.
func errorf(k ErrorCode, line int, f string, args ...interface{}) *Error {
	return &Error{k, "", line, fmt.Sprintf(f, args...)}
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// attrTypeMap[n] describes the value of the given attribute.
// If an attribute affects (or can mask) the encoding or interpretation of
// other content, or affects the contents, idempotency, or credentials of a
// network message, then the value in this map is contentTypeUnsafe.
// This map is derived from HTML5, specifically
// http://www.w3.org/TR/html5/Overview.html#attributes-1
// as well as "%URI"-typed attributes from
// http://www.w3.org/TR/html4/index/attributes.html
var attrTypeMap = map[string]contentType{
	"accept":          contentTypePlain,
	"accept-charset":  contentTypeUnsafe,
	"action":          contentTypeURL,
	"alt":             contentTypePlain,
	"archive":         contentTypeURL,
	"async":           contentTypeUnsafe,
	"autocomplete":    contentTypePlain,
	"autofocus":       contentTypePlain,
	"autoplay":        contentTypePlain,
	"background":      contentTypeURL,
	"border":          contentTypePlain,
	"checked":         contentTypePlain,
	"cite":            contentTypeURL,
	"challenge":       contentTypeUnsafe,
	"charset":         contentTypeUnsafe,
	"class":           contentTypePlain,
	"classid":         contentTypeURL,
	"codebase":        contentTypeURL,
	"cols":            contentTypePlain,
	"colspan":         contentTypePlain,
	"content":         contentTypeUnsafe,
	"contenteditable": contentTypePlain,
	"contextmenu":     contentTypePlain,
	"controls":        contentTypePlain,
	"coords":          contentTypePlain,
	"crossorigin":     contentTypeUnsafe,
	"data":            contentTypeURL,
	"datetime":        contentTypePlain,
	"default":         contentTypePlain,
	"defer":           contentTypeUnsafe,
	"dir":             contentTypePlain,
	"dirname":         contentTypePlain,
	"disabled":        contentTypePlain,
	"draggable":       contentTypePlain,
	"dropzone":        contentTypePlain,
	"enctype":         contentTypeUnsafe,
	"for":             contentTypePlain,
	"form":            contentTypeUnsafe,
	"formaction":      contentTypeURL,
	"formenctype":     contentTypeUnsafe,
	"formmethod":      contentTypeUnsafe,
	"formnovalidate":  contentTypeUnsafe,
	"formtarget":      contentTypePlain,
	"headers":         contentTypePlain,
	"height":          contentTypePlain,
	"hidden":          contentTypePlain,
	"high":            contentTypePlain,
	"href":            contentTypeURL,
	"hreflang":        contentTypePlain,
	"http-equiv":      contentTypeUnsafe,
	"icon":            contentTypeURL,
	"id":              contentTypePlain,
	"ismap":           contentTypePlain,
	"keytype":         contentTypeUnsafe,
	"kind":            contentTypePlain,
	"label":           contentTypePlain,
	"lang":            contentTypePlain,
	"language":        contentTypeUnsafe,
	"list":            contentTypePlain,
	"longdesc":        contentTypeURL,
	"loop":            contentTypePlain,
	"low":             contentTypePlain,
	"manifest":        contentTypeURL,
	"max":             contentTypePlain,
	"maxlength":       contentTypePlain,
	"media":           contentTypePlain,
	"mediagroup":      contentTypePlain,
	"method":          contentTypeUnsafe,
	"min":             contentTypePlain,
	"multiple":        contentTypePlain,
	"name":            contentTypePlain,
	"novalidate":      contentTypeUnsafe,
	// Skip handler names from
	// http://www.w3.org/TR/html5/webappapis.html#event-handlers-on-elements,-document-objects,-and-window-objects
	// since we have special handling in attrType.
	"open":        contentTypePlain,
	"optimum":     contentTypePlain,
	"pattern":     contentTypeUnsafe,
	"placeholder": contentTypePlain,
	"poster":      contentTypeURL,
	"profile":     contentTypeURL,
	"preload":     contentTypePlain,
	"pubdate":     contentTypePlain,
	"radiogroup":  contentTypePlain,
	"readonly":    contentTypePlain,
	"rel":         contentTypeUnsafe,
	"required":    contentTypePlain,
	"reversed":    contentTypePlain,
	"rows":        contentTypePlain,
	"rowspan":     contentTypePlain,
	"sandbox":     contentTypeUnsafe,
	"spellcheck":  contentTypePlain,
	"scope":       contentTypePlain,
	"scoped":      contentTypePlain,
	"seamless":    contentTypePlain,
	"selected":    contentTypePlain,
	"shape":       contentTypePlain,
	"size":        contentTypePlain,
	"sizes":       contentTypePlain,
	"span":        contentTypePlain,
	"src":         contentTypeURL,
	"srcdoc":      contentTypeHTML,
	"srclang":     contentTypePlain,
	"start":       contentTypePlain,
	"step":        contentTypePlain,
	"style":       contentTypeCSS,
	"tabindex":    contentTypePlain,
	"target":      contentTypePlain,
	"title":       contentTypePlain,
	"type":        contentTypeUnsafe,
	"usemap":      contentTypeURL,
	"value":       contentTypeUnsafe,
	"width":       contentTypePlain,
	"wrap":        contentTypePlain,
	"xmlns":       contentTypeURL,
}

// attrType returns a conservative (upper-bound on authority) guess at the
// type of the named attribute.
func attrType(name string) contentType {
	name = strings.ToLower(name)
	if strings.HasPrefix(name, "data-") {
		// Strip data- so that custom attribute heuristics below are
		// widely applied.
		// Treat data-action as URL below.
		name = name[5:]
	} else if colon := strings.IndexRune(name, ':'); colon != -1 {
		if name[:colon] == "xmlns" {
			return contentTypeURL
		}
		// Treat svg:href and xlink:href as href below.
		name = name[colon+1:]
	}
	if t, ok := attrTypeMap[name]; ok {
		return t
	}
	// Treat partial event handler names as script.
	if strings.HasPrefix(name, "on") {
		return contentTypeJS
	}

	// Heuristics to prevent "javascript:..." injection in custom
	// data attributes and custom attributes like g:tweetUrl.
	// http://www.w3.org/TR/html5/dom.html#embedding-custom-non-visible-data-with-the-data-*-attributes
	// "Custom data attributes are intended to store custom data
	//  private to the page or application, for which there are no
	//  more appropriate attributes or elements."
	// Developers seem to store URL content in data URLs that start
	// or end with "URI" or "URL".
	if strings.Contains(name, "src") ||
		strings.Contains(name, "uri") ||
		strings.Contains(name, "url") {
		return contentTypeURL
	}
	return contentTypePlain
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Strings of content from a trusted source.
type (
	// CSS encapsulates known safe content that matches any of:
	//   1. The CSS3 stylesheet production, such as `p { color: purple }`.
	//   2. The CSS3 rule production, such as `a[href=~"https:"].foo#bar`.
	//   3. CSS3 declaration productions, such as `color: red; margin: 2px`.
	//   4. The CSS3 value production, such as `rgba(0, 0, 255, 127)`.
	// See http://www.w3.org/TR/css3-syntax/#parsing and
	// https://web.archive.org/web/20090211114933/http://w3.org/TR/css3-syntax#style
	CSS string

	// HTML encapsulates a known safe HTML document fragment.
	// It should not be used for HTML from a third-party, or HTML with
	// unclosed tags or comments. The outputs of a sound HTML sanitizer
	// and a template escaped by this package are fine for use with HTML.
	HTML string

	// HTMLAttr encapsulates an HTML attribute from a trusted source,
	// for example, ` dir="ltr"`.
	HTMLAttr string

	// JS encapsulates a known safe EcmaScript5 Expression, for example,
	// `(x + y * z())`.
	// Template authors are responsible for ensuring that typed expressions
	// do not break the intended precedence and that there is no
	// statement/expression ambiguity as when passing an expression like
	// "{ foo: bar() }\n['foo']()", which is both a valid Expression and a
	// valid Program with a very different meaning.
	JS string

	// JSStr encapsulates a sequence of characters meant to be embedded
	// between quotes in a JavaScript expression.
	// The string must match a series of StringCharacters:
	//   StringCharacter :: SourceCharacter but not `\` or LineTerminator
	//                    | EscapeSequence
	// Note that LineContinuations are not allowed.
	// JSStr("foo\\nbar") is fine, but JSStr("foo\\\nbar") is not.
	JSStr string

	// URL encapsulates a known safe URL or URL substring (see RFC 3986).
	// A URL like `javascript:checkThatFormNotEditedBeforeLeavingPage()`
	// from a trusted source should go in the page, but by default dynamic
	// `javascript:` URLs are filtered out since they are a frequently
	// exploited injection vector.
	URL string
)

type contentType uint8

const (
	contentTypePlain contentType = iota
	contentTypeCSS
	contentTypeHTML
	contentTypeHTMLAttr
	contentTypeJS
	contentTypeJSStr
	contentTypeURL
	// contentTypeUnsafe is used in attr.go for values that affect how
	// embedded content and network messages are formed, vetted,
	// or interpreted; or which credentials network messages carry.
	contentTypeUnsafe
)

// indirect returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil).
func indirect(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	if t := reflect.TypeOf(a); t.Kind() != reflect.Ptr {
		// Avoid creating a reflect.Value if it's not a pointer.
		return a
	}
	v := reflect.ValueOf(a)
	for v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

var (
	errorType       = reflect.TypeOf((*error)(nil)).Elem()
	fmtStringerType = reflect.TypeOf((*fmt.Stringer)(nil)).Elem()
)

// indirectToStringerOrError returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of fmt.Stringer
// or error,
func indirectToStringerOrError(a interface{}) interface{} {
	if a == nil {
		return nil
	}
	v := reflect.ValueOf(a)
	for !v.Type().Implements(fmtStringerType) && !v.Type().Implements(errorType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// stringify converts its arguments to a string and the type of the content.
// All pointers are dereferenced, as in the text/template package.
func stringify(args ...interface{}) (string, contentType) {
	if len(args) == 1 {
		switch s := indirect(args[0]).(type) {
		case string:
			return s, contentTypePlain
		case CSS:
			return string(s), contentTypeCSS
		case HTML:
			return string(s), contentTypeHTML
		case HTMLAttr:
			return string(s), contentTypeHTMLAttr
		case JS:
			return string(s), contentTypeJS
		case JSStr:
			return string(s), contentTypeJSStr
		case URL:
			return string(s), contentTypeURL
		}
	}
	for i, arg := range args {
		args[i] = indirectToStringerOrError(arg)
	}
	return fmt.Sprint(args...), contentTypePlain
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// endsWithCSSKeyword reports whether b ends with an ident that
// case-insensitively matches the lower-case kw.
func endsWithCSSKeyword(b []byte, kw string) bool {
	i := len(b) - len(kw)
	if i < 0 {
		// Too short.
		return false
	}
	if i != 0 {
		r, _ := utf8.DecodeLastRune(b[:i])
		if isCSSNmchar(r) {
			// Too long.
			return false
		}
	}
	// Many CSS keywords, such as "!important" can have characters encoded,
	// but the URI production does not allow that according to
	// http://www.w3.org/TR/css3-syntax/#TOK-URI
	// This does not attempt to recognize encoded keywords. For example,
	// given "\75\72\6c" and "url" this return false.
	return string(bytes.ToLower(b[i:])) == kw
}

// isCSSNmchar reports whether rune is allowed anywhere in a CSS identifier.
func isCSSNmchar(r rune) bool {
	// Based on the CSS3 nmchar production but ignores multi-rune escape
	// sequences.
	// http://www.w3.org/TR/css3-syntax/#SUBTOK-nmchar
	return 'a' <= r && r <= 'z' ||
		'A' <= r && r <= 'Z' ||
		'0' <= r && r <= '9' ||
		r == '-' ||
		r == '_' ||
		// Non-ASCII cases below.
		0x80 <= r && r <= 0xd7ff ||
		0xe000 <= r && r <= 0xfffd ||
		0x10000 <= r && r <= 0x10ffff
}

// decodeCSS decodes CSS3 escapes given a sequence of stringchars.
// If there is no change, it returns the input, otherwise it returns a slice
// backed by a new array.
// http://www.w3.org/TR/css3-syntax/#SUBTOK-stringchar defines stringchar.
func decodeCSS(s []byte) []byte {
	i := bytes.IndexByte(s, '\\')
	if i == -1 {
		return s
	}
	// The UTF-8 sequence for a codepoint is never longer than 1 + the
	// number hex digits need to represent that codepoint, so len(s) is an
	// upper bound on the output length.
	b := make([]byte, 0, len(s))
	for len(s) != 0 {
		i := bytes.IndexByte(s, '\\')
		if i == -1 {
			i = len(s)
		}
		b, s = append(b, s[:i]...), s[i:]
		if len(s) < 2 {
			break
		}
		// http://www.w3.org/TR/css3-syntax/#SUBTOK-escape
		// escape ::= unicode | '\' [#x20-#x7E#x80-#xD7FF#xE000-#xFFFD#x10000-#x10FFFF]
		if isHex(s[1]) {
			// http://www.w3.org/TR/css3-syntax/#SUBTOK-unicode
			//   unicode ::= '\' [0-9a-fA-F]{1,6} wc?
			j := 2
			for j < len(s) && j < 7 && isHex(s[j]) {
				j++
			}
			r := hexDecode(s[1:j])
			if r > unicode.MaxRune {
				r, j = r/16, j-1
			}
			n := utf8.EncodeRune(b[len(b):cap(b)], r)
			// The optional space at the end allows a hex
			// sequence to be followed by a literal hex.
			// string(decodeCSS([]byte(`\A B`))) == "\nB"
			b, s = b[:len(b)+n], skipCSSSpace(s[j:])
		} else {
			// `\\` decodes to `\` and `\"` to `"`.
			_, n := utf8.DecodeRune(s[1:])
			b, s = append(b, s[1:1+n]...), s[1+n:]
		}
	}
	return b
}

// isHex reports whether the given character is a hex digit.
func isHex(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// hexDecode decodes a short hex digit sequence: "10" -> 16.
func hexDecode(s []byte) rune {
	n := '\x00'
	for _, c := range s {
		n <<= 4
		switch {
		case '0' <= c && c <= '9':
			n |= rune(c - '0')
		case 'a' <= c && c <= 'f':
			n |= rune(c-'a') + 10
		case 'A' <= c && c <= 'F':
			n |= rune(c-'A') + 10
		default:
			panic(fmt.Sprintf("Bad hex digit in %q", s))
		}
	}
	return n
}

// skipCSSSpace returns a suffix of c, skipping over a single space.
func skipCSSSpace(c []byte) []byte {
	if len(c) == 0 {
		return c
	}
	// wc ::= #x9 | #xA | #xC | #xD | #x20
	switch c[0] {
	case '\t', '\n', '\f', ' ':
		return c[1:]
	case '\r':
		// This differs from CSS3's wc production because it contains a
		// probable spec error whereby wc contains all the single byte
		// sequences in nl (newline) but not CRLF.
		if len(c) >= 2 && c[1] == '\n' {
			return c[2:]
		}
		return c[1:]
	}
	return c
}

// isCSSSpace reports whether b is a CSS space char as defined in wc.
func isCSSSpace(b byte) bool {
	switch b {
	case '\t', '\n', '\f', '\r', ' ':
		return true
	}
	return false
}

// cssEscaper escapes HTML and CSS special characters using \<hex>+ escapes.
func cssEscaper(args ...interface{}) string {
	s, _ := stringify(args...)
	var b bytes.Buffer
	written := 0
	for i, r := range s {
		var repl string
		switch r {
		case 0:
			repl = `\0`
		case '\t':
			repl = `\9`
		case '\n':
			repl = `\a`
		case '\f':
			repl = `\c`
		case '\r':
			repl = `\d`
		// Encode HTML specials as hex so the output can be embedded
		// in HTML attributes without further encoding.
		case '"':
			repl = `\22`
		case '&':
			repl = `\26`
		case '\'':
			repl = `\27`
		case '(':
			repl = `\28`
		case ')':
			repl = `\29`
		case '+':
			repl = `\2b`
		case '/':
			repl = `\2f`
		case ':':
			repl = `\3a`
		case ';':
			repl = `\3b`
		case '<':
			repl = `\3c`
		case '>':
			repl = `\3e`
		case '\\':
			repl = `\\`
		case '{':
			repl = `\7b`
		case '}':
			repl = `\7d`
		default:
			continue
		}
		b.WriteString(s[written:i])
		b.WriteString(repl)
		written = i + utf8.RuneLen(r)
		if repl != `\\` && (written == len(s) || isHex(s[written]) || isCSSSpace(s[written])) {
			b.WriteByte(' ')
		}
	}
	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}

var expressionBytes = []byte("expression")
var mozBindingBytes = []byte("mozbinding")

// cssValueFilter allows innocuous CSS values in the output including CSS
// quantities (10px or 25%), ID or class literals (#foo, .bar), keyword values
// (inherit, blue), and colors (#888).
// It filters out unsafe values, such as those that affect token boundaries,
// and anything that might execute scripts.
func cssValueFilter(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeCSS {
		return s
	}
	b, id := decodeCSS([]byte(s)), make([]byte, 0, 64)

	// CSS3 error handling is specified as honoring string boundaries per
	// http://www.w3.org/TR/css3-syntax/#error-handling :
	//     Malformed declarations. User agents must handle unexpected
	//     tokens encountered while parsing a declaration by reading until
	//     the end of the declaration, while observing the rules for
	//     matching pairs of (), [], {}, "", and '', and correctly handling
	//     escapes. For example, a malformed declaration may be missing a
	//     property, colon (:) or value.
	// So we need to make sure that values do not have mismatched bracket
	// or quote characters to prevent the browser from restarting parsing
	// inside a string that might embed JavaScript source.
	for i, c := range b {
		switch c {
		case 0, '"', '\'', '(', ')', '/', ';', '@', '[', '\\', ']', '`', '{', '}':
			return filterFailsafe
		case '-':
			// Disallow <!-- or -->.
			// -- should not appear in valid identifiers.
			if i != 0 && b[i-1] == '-' {
				return filterFailsafe
			}
		default:
			if c < 0x80 && isCSSNmchar(rune(c)) {
				id = append(id, c)
			}
		}
	}
	id = bytes.ToLower(id)
	if bytes.Index(id, expressionBytes) != -1 || bytes.Index(id, mozBindingBytes) != -1 {
		return filterFailsafe
	}
	return string(b)
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package template (html/template) implements data-driven templates for
generating HTML output safe against code injection. It provides the
same interface as package text/template and should be used instead of
text/template whenever the output is HTML.

The documentation here focuses on the security features of the package.
For information about how to program the templates themselves, see the
documentation for text/template.

Introduction

This package wraps package text/template so you can share its template API
to parse and execute HTML templates safely.

  tmpl, err := template.New("name").Parse(...)
  // Error checking elided
  err = tmpl.Execute(out, data)

If successful, tmpl will now be injection-safe. Otherwise, err is an error
defined in the docs for ErrorCode.

HTML templates treat data values as plain text which should be encoded so they
can be safely embedded in an HTML document. The escaping is contextual, so
actions can appear within JavaScript, CSS, and URI contexts.

The security model used by this package assumes that template authors are
trusted, while Execute's data parameter is not. More details are
provided below.

Example

  import "text/template"
  ...
  t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
  err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")

produces

  Hello, <script>alert('you have been pwned')</script>!

but the contextual autoescaping in html/template

  import "html/template"
  ...
  t, err := template.New("foo").Parse(`{{define "T"}}Hello, {{.}}!{{end}}`)
  err = t.ExecuteTemplate(out, "T", "<script>alert('you have been pwned')</script>")

produces safe, escaped HTML output

  Hello, &lt;script&gt;alert(&#39;you have been pwned&#39;)&lt;/script&gt;!


Contexts

This package understands HTML, CSS, JavaScript, and URIs. It adds sanitizing
functions to each simple action pipeline, so given the excerpt

  <a href="/search?q={{.}}">{{.}}</a>

At parse time each {{.}} is overwritten to add escaping functions as necessary.
In this case it becomes

  <a href="/search?q={{. | urlquery}}">{{. | html}}</a>


Errors

See the documentation of ErrorCode for details.


A fuller picture

The rest of this package comment may be skipped on first reading; it includes
details necessary to understand escaping contexts and error messages. Most users
will not need to understand these details.


Contexts

Assuming {{.}} is `O'Reilly: How are <i>you</i>?`, the table below shows
how {{.}} appears when used in the context to the left.

  Context                          {{.}} After
  {{.}}                            O'Reilly: How are &lt;i&gt;you&lt;/i&gt;?
  <a title='{{.}}'>                O&#39;Reilly: How are you?
  <a href="/{{.}}">                O&#39;Reilly: How are %3ci%3eyou%3c/i%3e?
  <a href="?q={{.}}">              O&#39;Reilly%3a%20How%20are%3ci%3e...%3f
  <a onx='f("{{.}}")'>             O\x27Reilly: How are \x3ci\x3eyou...?
  <a onx='f({{.}})'>               "O\x27Reilly: How are \x3ci\x3eyou...?"
  <a onx='pattern = /{{.}}/;'>     O\x27Reilly: How are \x3ci\x3eyou...\x3f

If used in an unsafe context, then the value might be filtered out:

  Context                          {{.}} After
  <a href="{{.}}">                 #ZgotmplZ

since "O'Reilly:" is not an allowed protocol like "http:".


If {{.}} is the innocuous word, `left`, then it can appear more widely,

  Context                              {{.}} After
  {{.}}                                left
  <a title='{{.}}'>                    left
  <a href='{{.}}'>                     left
  <a href='/{{.}}'>                    left
  <a href='?dir={{.}}'>                left
  <a style="border-{{.}}: 4px">        left
  <a style="align: {{.}}">             left
  <a style="background: '{{.}}'>       left
  <a style="background: url('{{.}}')>  left
  <style>p.{{.}} {color:red}</style>   left

Non-string values can be used in JavaScript contexts.
If {{.}} is

  struct{A,B string}{ "foo", "bar" }

in the escaped template

  <script>var pair = {{.}};</script>

then the template output is

  <script>var pair = {"A": "foo", "B": "bar"};</script>

See package json to understand how non-string content is marshalled for
embedding in JavaScript contexts.


Typed Strings

By default, this package assumes that all pipelines produce a plain text string.
It adds escaping pipeline stages necessary to correctly and safely embed that
plain text string in the appropriate context.

When a data value is not plain text, you can make sure it is not over-escaped
by marking it with its type.

Types HTML, JS, URL, and others from content.go can carry safe content that is
exempted from escaping.

The template

  Hello, {{.}}!

can be invoked with

  tmpl.Execute(out, HTML(`<b>World</b>`))

to produce

  Hello, <b>World</b>!

instead of the

  Hello, &lt;b&gt;World&lt;b&gt;!

that would have been produced if {{.}} was a regular string.


Security Model

http://js-quasis-libraries-and-repl.googlecode.com/svn/trunk/safetemplate.html#problem_definition defines "safe" as used by this package.

This package assumes that template authors are trusted, that Execute's data
parameter is not, and seeks to preserve the properties below in the face
of untrusted data:

Structure Preservation Property:
"... when a template author writes an HTML tag in a safe templating language,
the browser will interpret the corresponding portion of the output as a tag
regardless of the values of untrusted data, and similarly for other structures
such as attribute boundaries and JS and CSS string boundaries."

Code Effect Property:
"... only code specified by the template author should run as a result of
injecting the template output into a page and all code specified by the
template author should run as a result of the same."

Least Surprise Property:
"A developer (or code reviewer) familiar with HTML, CSS, and JavaScript, who
knows that contextual autoescaping happens should be able to look at a {{.}}
and correctly infer what sanitization happens."
*/
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nextJSCtx returns the context that determines whether a slash after the
// given run of tokens starts a regular expression instead of a division
// operator: / or /=.
//
// This assumes that the token run does not include any string tokens, comment
// tokens, regular expression literal tokens, or division operators.
//
// This fails on some valid but nonsensical JavaScript programs like
// "x = ++/foo/i" which is quite different than "x++/foo/i", but is not known to
// fail on any known useful programs. It is based on the draft
// JavaScript 2.0 lexical grammar and requires one token of lookbehind:
// http://www.mozilla.org/js/language/js20-2000-07/rationale/syntax.html
func nextJSCtx(s []byte, preceding jsCtx) jsCtx {
	s = bytes.TrimRight(s, "\t\n\f\r \u2028\u2029")
	if len(s) == 0 {
		return preceding
	}

	// All cases below are in the single-byte UTF-8 group.
	switch c, n := s[len(s)-1], len(s); c {
	case '+', '-':
		// ++ and -- are not regexp preceders, but + and - are whether
		// they are used as infix or prefix operators.
		start := n - 1
		// Count the number of adjacent dashes or pluses.
		for start > 0 && s[start-1] == c {
			start--
		}
		if (n-start)&1 == 1 {
			// Reached for trailing minus signs since "---" is the
			// same as "-- -".
			return jsCtxRegexp
		}
		return jsCtxDivOp
	case '.':
		// Handle "42."
		if n != 1 && '0' <= s[n-2] && s[n-2] <= '9' {
			return jsCtxDivOp
		}
		return jsCtxRegexp
	// Suffixes for all punctuators from section 7.7 of the language spec
	// that only end binary operators not handled above.
	case ',', '<', '>', '=', '*', '%', '&', '|', '^', '?':
		return jsCtxRegexp
	// Suffixes for all punctuators from section 7.7 of the language spec
	// that are prefix operators not handled above.
	case '!', '~':
		return jsCtxRegexp
	// Matches all the punctuators from section 7.7 of the language spec
	// that are open brackets not handled above.
	case '(', '[':
		return jsCtxRegexp
	// Matches all the punctuators from section 7.7 of the language spec
	// that precede expression starts.
	case ':', ';', '{':
		return jsCtxRegexp
	// CAVEAT: the close punctuators ('}', ']', ')') precede div ops and
	// are handled in the default except for '}' which can precede a
	// division op as in
	//    ({ valueOf: function () { return 42 } } / 2
	// which is valid, but, in practice, developers don't divide object
	// literals, so our heuristic works well for code like
	//    function () { ... }  /foo/.test(x) && sideEffect();
	// The ')' punctuator can precede a regular expression as in
	//     if (b) /foo/.test(x) && ...
	// but this is much less likely than
	//     (a + b) / c
	case '}':
		return jsCtxRegexp
	default:
		// Look for an IdentifierName and see if it is a keyword that
		// can precede a regular expression.
		j := n
		for j > 0 && isJSIdentPart(rune(s[j-1])) {
			j--
		}
		if regexpPrecederKeywords[string(s[j:])] {
			return jsCtxRegexp
		}
	}
	// Otherwise is a punctuator not listed above, or
	// a string which precedes a div op, or an identifier
	// which precedes a div op.
	return jsCtxDivOp
}

// regexpPrecederKeywords is a set of reserved JS keywords that can precede a
// regular expression in JS source.
var regexpPrecederKeywords = map[string]bool{
	"break":      true,
	"case":       true,
	"continue":   true,
	"delete":     true,
	"do":         true,
	"else":       true,
	"finally":    true,
	"in":         true,
	"instanceof": true,
	"return":     true,
	"throw":      true,
	"try":        true,
	"typeof":     true,
	"void":       true,
}

var jsonMarshalType = reflect.TypeOf((*json.Marshaler)(nil)).Elem()

// indirectToJSONMarshaler returns the value, after dereferencing as many times
// as necessary to reach the base type (or nil) or an implementation of json.Marshal.
func indirectToJSONMarshaler(a interface{}) interface{} {
	v := reflect.ValueOf(a)
	for !v.Type().Implements(jsonMarshalType) && v.Kind() == reflect.Ptr && !v.IsNil() {
		v = v.Elem()
	}
	return v.Interface()
}

// jsValEscaper escapes its inputs to a JS Expression (section 11.14) that has
// neither side-effects nor free variables outside (NaN, Infinity).
func jsValEscaper(args ...interface{}) string {
	var a interface{}
	if len(args) == 1 {
		a = indirectToJSONMarshaler(args[0])
		switch t := a.(type) {
		case JS:
			return string(t)
		case JSStr:
			// TODO: normalize quotes.
			return `"` + string(t) + `"`
		case json.Marshaler:
			// Do not treat as a Stringer.
		case fmt.Stringer:
			a = t.String()
		}
	} else {
		for i, arg := range args {
			args[i] = indirectToJSONMarshaler(arg)
		}
		a = fmt.Sprint(args...)
	}
	// TODO: detect cycles before calling Marshal which loops infinitely on
	// cyclic data. This may be an unacceptable DoS risk.

	b, err := json.Marshal(a)
	if err != nil {
		// Put a space before comment so that if it is flush against
		// a division operator it is not turned into a line comment:
		//     x/{{y}}
		// turning into
		//     x//* error marshalling y:
		//          second line of error message */null
		return fmt.Sprintf(" /* %s */null ", strings.Replace(err.Error(), "*/", "* /", -1))
	}

	// TODO: maybe post-process output to prevent it from containing
	// "<!--", "-->", "<![CDATA[", "]]>", or "</script"
	// in case custom marshallers produce output containing those.

	// TODO: Maybe abbreviate \u00ab to \xab to produce more compact output.
	if len(b) == 0 {
		// In, `x=y/{{.}}*z` a json.Marshaler that produces "" should
		// not cause the output `x=y/*z`.
		return " null "
	}
	first, _ := utf8.DecodeRune(b)
	last, _ := utf8.DecodeLastRune(b)
	var buf bytes.Buffer
	// Prevent IdentifierNames and NumericLiterals from running into
	// keywords: in, instanceof, typeof, void
	pad := isJSIdentPart(first) || isJSIdentPart(last)
	if pad {
		buf.WriteByte(' ')
	}
	written := 0
	// Make sure that json.Marshal escapes codepoints U+2028 & U+2029
	// so it falls within the subset of JSON which is valid JS.
	for i := 0; i < len(b); {
		rune, n := utf8.DecodeRune(b[i:])
		repl := ""
		if rune == 0x2028 {
			repl = `\u2028`
		} else if rune == 0x2029 {
			repl = `\u2029`
		}
		if repl != "" {
			buf.Write(b[written:i])
			buf.WriteString(repl)
			written = i + n
		}
		i += n
	}
	if buf.Len() != 0 {
		buf.Write(b[written:])
		if pad {
			buf.WriteByte(' ')
		}
		b = buf.Bytes()
	}
	return string(b)
}

// jsStrEscaper produces a string that can be included between quotes in
// JavaScript source, in JavaScript embedded in an HTML5 <script> element,
// or in an HTML5 event handler attribute such as onclick.
func jsStrEscaper(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeJSStr {
		return replace(s, jsStrNormReplacementTable)
	}
	return replace(s, jsStrReplacementTable)
}

// jsRegexpEscaper behaves like jsStrEscaper but escapes regular expression
// specials so the result is treated literally when included in a regular
// expression literal. /foo{{.X}}bar/ matches the string "foo" followed by
// the literal text of {{.X}} followed by the string "bar".
func jsRegexpEscaper(args ...interface{}) string {
	s, _ := stringify(args...)
	s = replace(s, jsRegexpReplacementTable)
	if s == "" {
		// /{{.X}}/ should not produce a line comment when .X == "".
		return "(?:)"
	}
	return s
}

// replace replaces each rune r of s with replacementTable[r], provided that
// r < len(replacementTable). If replacementTable[r] is the empty string then
// no replacement is made.
// It also replaces runes U+2028 and U+2029 with the raw strings `\u2028` and
// `\u2029`.
func replace(s string, replacementTable []string) string {
	var b bytes.Buffer
	written := 0
	for i, r := range s {
		var repl string
		switch {
		case int(r) < len(replacementTable) && replacementTable[r] != "":
			repl = replacementTable[r]
		case r == '\u2028':
			repl = `\u2028`
		case r == '\u2029':
			repl = `\u2029`
		default:
			continue
		}
		b.WriteString(s[written:i])
		b.WriteString(repl)
		written = i + utf8.RuneLen(r)
	}
	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}

var jsStrReplacementTable = []string{
	0:    `\0`,
	'\t': `\t`,
	'\n': `\n`,
	'\v': `\x0b`, // "\v" == "v" on IE 6.
	'\f': `\f`,
	'\r': `\r`,
	// Encode HTML specials as hex so the output can be embedded
	// in HTML attributes without further encoding.
	'"':  `\x22`,
	'&':  `\x26`,
	'\'': `\x27`,
	'+':  `\x2b`,
	'/':  `\/`,
	'<':  `\x3c`,
	'>':  `\x3e`,
	'\\': `\\`,
}

// jsStrNormReplacementTable is like jsStrReplacementTable but does not
// overencode existing escapes since this table has no entry for `\`.
var jsStrNormReplacementTable = []string{
	0:    `\0`,
	'\t': `\t`,
	'\n': `\n`,
	'\v': `\x0b`, // "\v" == "v" on IE 6.
	'\f': `\f`,
	'\r': `\r`,
	// Encode HTML specials as hex so the output can be embedded
	// in HTML attributes without further encoding.
	'"':  `\x22`,
	'&':  `\x26`,
	'\'': `\x27`,
	'+':  `\x2b`,
	'/':  `\/`,
	'<':  `\x3c`,
	'>':  `\x3e`,
}

var jsRegexpReplacementTable = []string{
	0:    `\0`,
	'\t': `\t`,
	'\n': `\n`,
	'\v': `\x0b`, // "\v" == "v" on IE 6.
	'\f': `\f`,
	'\r': `\r`,
	// Encode HTML specials as hex so the output can be embedded
	// in HTML attributes without further encoding.
	'"':  `\x22`,
	'$':  `\$`,
	'&':  `\x26`,
	'\'': `\x27`,
	'(':  `\(`,
	')':  `\)`,
	'*':  `\*`,
	'+':  `\x2b`,
	'-':  `\-`,
	'.':  `\.`,
	'/':  `\/`,
	'<':  `\x3c`,
	'>':  `\x3e`,
	'?':  `\?`,
	'[':  `\[`,
	'\\': `\\`,
	']':  `\]`,
	'^':  `\^`,
	'{':  `\{`,
	'|':  `\|`,
	'}':  `\}`,
}

// isJSIdentPart reports whether the given rune is a JS identifier part.
// It does not handle all the non-Latin letters, joiners, and combining marks,
// but it does handle every codepoint that can occur in a numeric literal or
// a keyword.
func isJSIdentPart(r rune) bool {
	switch {
	case r == '$':
		return true
	case '0' <= r && r <= '9':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case r == '_':
		return true
	case 'a' <= r && r <= 'z':
		return true
	}
	return false
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Template is a specialized Template from "text/template" that produces a safe
// HTML document fragment.
type Template struct {
	escaped bool
	// We could embed the text/template field, but it's safer not to because
	// we need to keep our version of the name space and the underlying
	// template's in sync.
	text *template.Template
	// The underlying template's parse tree, updated to be HTML-safe.
	Tree       *parse.Tree
	*nameSpace // common to all associated templates
}

// nameSpace is the data structure shared by all templates in an association.
type nameSpace struct {
	mu  sync.Mutex
	set map[string]*Template
}

// Templates returns a slice of the templates associated with t, including t
// itself.
func (t *Template) Templates() []*Template {
	ns := t.nameSpace
	ns.mu.Lock()
	defer ns.mu.Unlock()
	// Return a slice so we don't expose the map.
	m := make([]*Template, 0, len(ns.set))
	for _, v := range ns.set {
		m = append(m, v)
	}
	return m
}

// escape escapes all associated templates.
func (t *Template) escape() error {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	if !t.escaped {
		if err := escapeTemplates(t, t.Name()); err != nil {
			return err
		}
		t.escaped = true
	}
	return nil
}

// Execute applies a parsed template to the specified data object,
// writing the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel.
func (t *Template) Execute(wr io.Writer, data interface{}) error {
	if err := t.escape(); err != nil {
		return err
	}
	return t.text.Execute(wr, data)
}

// ExecuteTemplate applies the template associated with t that has the given
// name to the specified data object and writes the output to wr.
// If an error occurs executing the template or writing its output,
// execution stops, but partial results may already have been written to
// the output writer.
// A template may be executed safely in parallel.
func (t *Template) ExecuteTemplate(wr io.Writer, name string, data interface{}) error {
	tmpl, err := t.lookupAndEscapeTemplate(name)
	if err != nil {
		return err
	}
	return tmpl.text.Execute(wr, data)
}

// lookupAndEscapeTemplate guarantees that the template with the given name
// is escaped, or returns an error if it cannot be. It returns the named
// template.
func (t *Template) lookupAndEscapeTemplate(name string) (tmpl *Template, err error) {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	tmpl = t.set[name]
	if tmpl == nil {
		return nil, fmt.Errorf("html/template: %q is undefined", name)
	}
	if tmpl.text.Tree == nil || tmpl.text.Root == nil {
		return nil, fmt.Errorf("html/template: %q is an incomplete template", name)
	}
	if t.text.Lookup(name) == nil {
		panic("html/template internal error: template escaping out of sync")
	}
	if tmpl != nil && !tmpl.escaped {
		err = escapeTemplates(tmpl, name)
	}
	return tmpl, err
}

// Parse parses a string into a template. Nested template definitions
// will be associated with the top-level template t. Parse may be
// called multiple times to parse definitions of templates to associate
// with t. It is an error if a resulting template is non-empty (contains
// content other than template definitions) and would replace a
// non-empty template with the same name.  (In multiple calls to Parse
// with the same receiver template, only one call can contain text
// other than space, comments, and template definitions.)
func (t *Template) Parse(src string) (*Template, error) {
	t.nameSpace.mu.Lock()
	t.escaped = false
	t.nameSpace.mu.Unlock()
	ret, err := t.text.Parse(src)
	if err != nil {
		return nil, err
	}
	// In general, all the named templates might have changed underfoot.
	// Regardless, some new ones may have been defined.
	// The template.Template set has been updated; update ours.
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	for _, v := range ret.Templates() {
		name := v.Name()
		tmpl := t.set[name]
		if tmpl == nil {
			tmpl = t.new(name)
		}
		// Restore our record of this text/template to its unescaped original state.
		tmpl.escaped = false
		tmpl.text = v
		tmpl.Tree = v.Tree
	}
	return t, nil
}

// AddParseTree creates a new template with the name and parse tree
// and associates it with t.
//
// It returns an error if t has already been executed.
func (t *Template) AddParseTree(name string, tree *parse.Tree) (*Template, error) {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	if t.escaped {
		return nil, fmt.Errorf("html/template: cannot AddParseTree to %q after it has executed", t.Name())
	}
	text, err := t.text.AddParseTree(name, tree)
	if err != nil {
		return nil, err
	}
	ret := &Template{
		false,
		text,
		text.Tree,
		t.nameSpace,
	}
	t.set[name] = ret
	return ret, nil
}

// Clone returns a duplicate of the template, including all associated
// templates. The actual representation is not copied, but the name space of
// associated templates is, so further calls to Parse in the copy will add
// templates to the copy but not to the original. Clone can be used to prepare
// common templates and use them with variant definitions for other templates
// by adding the variants after the clone is made.
//
// It returns an error if t has already been executed.
func (t *Template) Clone() (*Template, error) {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	if t.escaped {
		return nil, fmt.Errorf("html/template: cannot Clone %q after it has executed", t.Name())
	}
	textClone, err := t.text.Clone()
	if err != nil {
		return nil, err
	}
	ret := &Template{
		false,
		textClone,
		textClone.Tree,
		&nameSpace{
			set: make(map[string]*Template),
		},
	}
	for _, x := range textClone.Templates() {
		name := x.Name()
		src := t.set[name]
		if src == nil || src.escaped {
			return nil, fmt.Errorf("html/template: cannot Clone %q after it has executed", t.Name())
		}
		x.Tree = x.Tree.Copy()
		ret.set[name] = &Template{
			false,
			x,
			x.Tree,
			ret.nameSpace,
		}
	}
	return ret, nil
}

// New allocates a new HTML template with the given name.
func New(name string) *Template {
	tmpl := &Template{
		false,
		template.New(name),
		nil,
		&nameSpace{
			set: make(map[string]*Template),
		},
	}
	tmpl.set[name] = tmpl
	return tmpl
}

// New allocates a new HTML template associated with the given one
// and with the same delimiters. The association, which is transitive,
// allows one template to invoke another with a {{template}} action.
func (t *Template) New(name string) *Template {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	return t.new(name)
}

// new is the implementation of New, without the lock.
func (t *Template) new(name string) *Template {
	tmpl := &Template{
		false,
		t.text.New(name),
		nil,
		t.nameSpace,
	}
	tmpl.set[name] = tmpl
	return tmpl
}

// Name returns the name of the template.
func (t *Template) Name() string {
	return t.text.Name()
}

// FuncMap is the type of the map defining the mapping from names to
// functions. Each function must have either a single return value, or two
// return values of which the second has type error. In that case, if the
// second (error) argument evaluates to non-nil during execution, execution
// terminates and Execute returns that error. FuncMap has the same base type
// as FuncMap in "text/template", copied here so clients need not import
// "text/template".
type FuncMap map[string]interface{}

// Funcs adds the elements of the argument map to the template's function map.
// It panics if a value in the map is not a function with appropriate return
// type. However, it is legal to overwrite elements of the map. The return
// value is the template, so calls can be chained.
func (t *Template) Funcs(funcMap FuncMap) *Template {
	t.text.Funcs(template.FuncMap(funcMap))
	return t
}

// Delims sets the action delimiters to the specified strings, to be used in
// subsequent calls to Parse, ParseFiles, or ParseGlob. Nested template
// definitions will inherit the settings. An empty delimiter stands for the
// corresponding default: {{ or }}.
// The return value is the template, so calls can be chained.
func (t *Template) Delims(left, right string) *Template {
	t.text.Delims(left, right)
	return t
}

// Lookup returns the template with the given name that is associated with t,
// or nil if there is no such template.
func (t *Template) Lookup(name string) *Template {
	t.nameSpace.mu.Lock()
	defer t.nameSpace.mu.Unlock()
	return t.set[name]
}

// Must is a helper that wraps a call to a function returning (*Template, error)
// and panics if the error is non-nil. It is intended for use in variable initializations
// such as
//	var t = template.Must(template.New("name").Parse("html"))
func Must(t *Template, err error) *Template {
	if err != nil {
		panic(err)
	}
	return t
}

// ParseFiles creates a new Template and parses the template definitions from
// the named files. The returned template's name will have the (base) name and
// (parsed) contents of the first file. There must be at least one file.
// If an error occurs, parsing stops and the returned *Template is nil.
func ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(nil, filenames...)
}

// ParseFiles parses the named files and associates the resulting templates with
// t. If an error occurs, parsing stops and the returned template is nil;
// otherwise it is t. There must be at least one file.
func (t *Template) ParseFiles(filenames ...string) (*Template, error) {
	return parseFiles(t, filenames...)
}

// parseFiles is the helper for the method and function. If the argument
// template is nil, it is created from the first file.
func parseFiles(t *Template, filenames ...string) (*Template, error) {
	if len(filenames) == 0 {
		// Not really a problem, but be consistent.
		return nil, fmt.Errorf("html/template: no files named in call to ParseFiles")
	}
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}
		s := string(b)
		name := filepath.Base(filename)
		// First template becomes return value if not already defined,
		// and we use that one for subsequent New calls to associate
		// all the templates together. Also, if this file has the same name
		// as t, this file becomes the contents of t, so
		//  t, err := New(name).Funcs(xxx).ParseFiles(name)
		// works. Otherwise we create a new template associated with t.
		var tmpl *Template
		if t == nil {
			t = New(name)
		}
		if name == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(name)
		}
		_, err = tmpl.Parse(s)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

// ParseGlob creates a new Template and parses the template definitions from the
// files identified by the pattern, which must match at least one file. The
// returned template will have the (base) name and (parsed) contents of the
// first file matched by the pattern. ParseGlob is equivalent to calling
// ParseFiles with the list of files matched by the pattern.
func ParseGlob(pattern string) (*Template, error) {
	return parseGlob(nil, pattern)
}

// ParseGlob parses the template definitions in the files identified by the
// pattern and associates the resulting templates with t. The pattern is
// processed by filepath.Glob and must match at least one file. ParseGlob is
// equivalent to calling t.ParseFiles with the list of files matched by the
// pattern.
func (t *Template) ParseGlob(pattern string) (*Template, error) {
	return parseGlob(t, pattern)
}

// parseGlob is the implementation of the function and method ParseGlob.
func parseGlob(t *Template, pattern string) (*Template, error) {
	filenames, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	if len(filenames) == 0 {
		return nil, fmt.Errorf("html/template: pattern matches no files: %#q", pattern)
	}
	return parseFiles(t, filenames...)
}

// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// urlFilter returns its input unless it contains an unsafe protocol in which
// case it defangs the entire URL.
func urlFilter(args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeURL {
		return s
	}
	if i := strings.IndexRune(s, ':'); i >= 0 && strings.IndexRune(s[:i], '/') < 0 {
		protocol := strings.ToLower(s[:i])
		if protocol != "http" && protocol != "https" && protocol != "mailto" {
			return "#" + filterFailsafe
		}
	}
	return s
}

// urlEscaper produces an output that can be embedded in a URL query.
// The output can be embedded in an HTML attribute without further escaping.
func urlEscaper(args ...interface{}) string {
	return urlProcessor(false, args...)
}

// urlEscaper normalizes URL content so it can be embedded in a quote-delimited
// string or parenthesis delimited url(...).
// The normalizer does not encode all HTML specials. Specifically, it does not
// encode '&' so correct embedding in an HTML attribute requires escaping of
// '&' to '&amp;'.
func urlNormalizer(args ...interface{}) string {
	return urlProcessor(true, args...)
}

// urlProcessor normalizes (when norm is true) or escapes its input to produce
// a valid hierarchical or opaque URL part.
func urlProcessor(norm bool, args ...interface{}) string {
	s, t := stringify(args...)
	if t == contentTypeURL {
		norm = true
	}
	var b bytes.Buffer
	written := 0
	// The byte loop below assumes that all URLs use UTF-8 as the
	// content-encoding. This is similar to the URI to IRI encoding scheme
	// defined in section 3.1 of  RFC 3987, and behaves the same as the
	// EcmaScript builtin encodeURIComponent.
	// It should not cause any misencoding of URLs in pages with
	// Content-type: text/html;charset=UTF-8.
	for i, n := 0, len(s); i < n; i++ {
		c := s[i]
		switch c {
		// Single quote and parens are sub-delims in RFC 3986, but we
		// escape them so the output can be embedded in single
		// quoted attributes and unquoted CSS url(...) constructs.
		// Single quotes are reserved in URLs, but are only used in
		// the obsolete "mark" rule in an appendix in RFC 3986
		// so can be safely encoded.
		case '!', '#', '$', '&', '*', '+', ',', '/', ':', ';', '=', '?', '@', '[', ']':
			if norm {
				continue
			}
		// Unreserved according to RFC 3986 sec 2.3
		// "For consistency, percent-encoded octets in the ranges of
		// ALPHA (%41-%5A and %61-%7A), DIGIT (%30-%39), hyphen (%2D),
		// period (%2E), underscore (%5F), or tilde (%7E) should not be
		// created by URI producers
		case '-', '.', '_', '~':
			continue
		case '%':
			// When normalizing do not re-encode valid escapes.
			if norm && i+2 < len(s) && isHex(s[i+1]) && isHex(s[i+2]) {
				continue
			}
		default:
			// Unreserved according to RFC 3986 sec 2.3
			if 'a' <= c && c <= 'z' {
				continue
			}
			if 'A' <= c && c <= 'Z' {
				continue
			}
			if '0' <= c && c <= '9' {
				continue
			}
		}
		b.WriteString(s[written:i])
		fmt.Fprintf(&b, "%%%02x", c)
		written = i + 1
	}
	if written == 0 {
		return s
	}
	b.WriteString(s[written:])
	return b.String()
}
