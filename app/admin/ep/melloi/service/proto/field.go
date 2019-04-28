// Copyright (c) 2017 Ernest Micklei
//
// MIT License
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package proto

import (
	"text/scanner"
)

// Field is an abstract message field.
type Field struct {
	Position      scanner.Position
	Comment       *Comment
	Name          string
	Type          string
	Sequence      int
	Options       []*Option
	InlineComment *Comment
	Parent        Visitee
}

// inlineComment is part of commentInliner.
func (f *Field) inlineComment(c *Comment) {
	f.InlineComment = c
}

// NormalField represents a field in a Message.
type NormalField struct {
	*Field
	Repeated bool
	Optional bool // proto2
	Required bool // proto2
}

func newNormalField() *NormalField { return &NormalField{Field: new(Field)} }

// Accept dispatches the call to the visitor.
func (f *NormalField) Accept(v Visitor) {
	v.VisitNormalField(f)
}

// Doc is part of Documented
func (f *NormalField) Doc() *Comment {
	return f.Comment
}

// parse expects:
// [ "repeated" | "optional" ] type fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
func (f *NormalField) parse(p *Parser) error {
	for {
		_, tok, lit := p.nextTypeName()
		switch tok {
		case tREPEATED:
			f.Repeated = true
			return f.parse(p)
		case tOPTIONAL: // proto2
			f.Optional = true
			return f.parse(p)
		case tIDENT:
			f.Type = lit
			return parseFieldAfterType(f.Field, p)
		default:
			goto done
		}
	}
done:
	return nil
}

// parseFieldAfterType expects:
// fieldName "=" fieldNumber [ "[" fieldOptions "]" ] ";
func parseFieldAfterType(f *Field, p *Parser) error {
	pos, tok, lit := p.next()
	if tok != tIDENT {
		if !isKeyword(tok) {
			return p.unexpected(lit, "field identifier", f)
		}
	}
	f.Name = lit
	pos, tok, lit = p.next()
	if tok != tEQUALS {
		return p.unexpected(lit, "field =", f)
	}
	i, err := p.nextInteger()
	if err != nil {
		return p.unexpected(lit, "field sequence number", f)
	}
	f.Sequence = i
	// see if there are options
	pos, tok, _ = p.next()
	if tLEFTSQUARE != tok {
		p.nextPut(pos, tok, lit)
		return nil
	}
	// consume options
	for {
		o := new(Option)
		o.Position = pos
		o.IsEmbedded = true
		err := o.parse(p)
		if err != nil {
			return err
		}
		f.Options = append(f.Options, o)

		pos, tok, lit = p.next()
		if tRIGHTSQUARE == tok {
			break
		}
		if tCOMMA != tok {
			return p.unexpected(lit, "option ,", o)
		}
	}
	return nil
}

// MapField represents a map entry in a message.
type MapField struct {
	*Field
	KeyType string
}

func newMapField() *MapField { return &MapField{Field: new(Field)} }

// Accept dispatches the call to the visitor.
func (f *MapField) Accept(v Visitor) {
	v.VisitMapField(f)
}

// parse expects:
// mapField = "map" "<" keyType "," type ">" mapName "=" fieldNumber [ "[" fieldOptions "]" ] ";"
// keyType = "int32" | "int64" | "uint32" | "uint64" | "sint32" | "sint64" |
//           "fixed32" | "fixed64" | "sfixed32" | "sfixed64" | "bool" | "string"
func (f *MapField) parse(p *Parser) error {
	_, tok, lit := p.next()
	if tLESS != tok {
		return p.unexpected(lit, "map keyType <", f)
	}
	_, tok, lit = p.nextTypeName()
	if tIDENT != tok {
		return p.unexpected(lit, "map identifier", f)
	}
	f.KeyType = lit
	_, tok, lit = p.next()
	if tCOMMA != tok {
		return p.unexpected(lit, "map type separator ,", f)
	}
	_, tok, lit = p.nextTypeName()
	if tIDENT != tok {
		return p.unexpected(lit, "map valueType identifier", f)
	}
	f.Type = lit
	_, tok, lit = p.next()
	if tGREATER != tok {
		return p.unexpected(lit, "map valueType >", f)
	}
	return parseFieldAfterType(f.Field, p)
}

func (f *Field) parent(v Visitee) { f.Parent = v }
