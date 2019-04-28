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
	"strings"
)

// token represents a lexical token.
type token int

const (
	// Special tokens
	tILLEGAL token = iota
	tEOF
	tWS

	// Literals
	tIDENT

	// Misc characters
	tSEMICOLON   // ;
	tCOLON       // :
	tEQUALS      // =
	tQUOTE       // "
	tSINGLEQUOTE // '
	tLEFTPAREN   // (
	tRIGHTPAREN  // )
	tLEFTCURLY   // {
	tRIGHTCURLY  // }
	tLEFTSQUARE  // [
	tRIGHTSQUARE // ]
	tCOMMENT     // /
	tLESS        // <
	tGREATER     // >
	tCOMMA       // ,
	tDOT         // .

	// Keywords
	keywordsStart
	tSYNTAX
	tSERVICE
	tRPC
	tRETURNS
	tMESSAGE
	tIMPORT
	tPACKAGE
	tOPTION
	tREPEATED
	tWEAK
	tPUBLIC

	// special fields
	tONEOF
	tMAP
	tRESERVED
	tENUM
	tSTREAM

	// BEGIN proto2
	tOPTIONAL
	tGROUP
	tEXTENSIONS
	tEXTEND
	tREQUIRED
	// END proto2
	keywordsEnd
)

// typeTokens exists for future validation
// const typeTokens = "double float int32 int64 uint32 uint64 sint32 sint64 fixed32 sfixed32 sfixed64 bool string bytes"

// isKeyword returns if tok is in the keywords range
func isKeyword(tok token) bool {
	return keywordsStart < tok && tok < keywordsEnd
}

// isWhitespace checks for space,tab and newline
func isWhitespace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\n'
}

// isDigit returns true if the rune is a digit.
func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') }

// isString checks if the literal is quoted (single or double).
func isString(lit string) bool {
	return (strings.HasPrefix(lit, "\"") &&
		strings.HasSuffix(lit, "\"")) ||
		(strings.HasPrefix(lit, "'") &&
			strings.HasSuffix(lit, "'"))
}

func isComment(lit string) bool {
	return strings.HasPrefix(lit, "//") || strings.HasPrefix(lit, "/*")
}

func unQuote(lit string) string {
	return strings.Trim(lit, "\"'")
}

func asToken(literal string) token {
	switch literal {
	// delimiters
	case ";":
		return tSEMICOLON
	case ":":
		return tCOLON
	case "=":
		return tEQUALS
	case "\"":
		return tQUOTE
	case "'":
		return tSINGLEQUOTE
	case "(":
		return tLEFTPAREN
	case ")":
		return tRIGHTPAREN
	case "{":
		return tLEFTCURLY
	case "}":
		return tRIGHTCURLY
	case "[":
		return tLEFTSQUARE
	case "]":
		return tRIGHTSQUARE
	case "<":
		return tLESS
	case ">":
		return tGREATER
	case ",":
		return tCOMMA
	case ".":
		return tDOT
	// words
	case "syntax":
		return tSYNTAX
	case "service":
		return tSERVICE
	case "rpc":
		return tRPC
	case "returns":
		return tRETURNS
	case "option":
		return tOPTION
	case "message":
		return tMESSAGE
	case "import":
		return tIMPORT
	case "package":
		return tPACKAGE
	case "oneof":
		return tONEOF
	// special fields
	case "map":
		return tMAP
	case "reserved":
		return tRESERVED
	case "enum":
		return tENUM
	case "repeated":
		return tREPEATED
	case "weak":
		return tWEAK
	case "public":
		return tPUBLIC
	case "stream":
		return tSTREAM
	// proto2
	case "optional":
		return tOPTIONAL
	case "group":
		return tGROUP
	case "extensions":
		return tEXTENSIONS
	case "extend":
		return tEXTEND
	case "required":
		return tREQUIRED
	case "ws":
		return tWS
	case "ill":
		return tILLEGAL
	default:
		// special cases
		if isComment(literal) {
			return tCOMMENT
		}
		return tIDENT
	}
}
