package grok

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	g, _ := New()
	if len(g.patterns) == 0 {
		t.Fatal("the Grok object should have some patterns pre loaded")
	}

	g, _ = NewWithConfig(&Config{NamedCapturesOnly: true})
	if len(g.patterns) == 0 {
		t.Fatal("the Grok object should have some patterns pre loaded")
	}
}

func TestParseWithDefaultCaptureMode(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	if captures, err := g.Parse("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["timestamp"] != "23/Apr/2014:22:58:32 +0200" {
			t.Fatalf("%s should be '%s' have '%s'", "timestamp", "23/Apr/2014:22:58:32 +0200", captures["timestamp"])
		}
		if captures["TIME"] != "" {
			t.Fatalf("%s should be '%s' have '%s'", "TIME", "", captures["TIME"])
		}
	}

	g, _ = New()
	if captures, err := g.Parse("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["timestamp"] != "23/Apr/2014:22:58:32 +0200" {
			t.Fatalf("%s should be '%s' have '%s'", "timestamp", "23/Apr/2014:22:58:32 +0200", captures["timestamp"])
		}
		if captures["TIME"] != "22:58:32" {
			t.Fatalf("%s should be '%s' have '%s'", "TIME", "22:58:32", captures["TIME"])
		}
	}
}

func TestMultiParseWithDefaultCaptureMode(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	res, _ := g.ParseToMultiMap("%{COMMONAPACHELOG} %{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:23:58:32 +0200] "GET /index.php HTTP/1.1" 404 207 127.0.0.1 - - [24/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	if len(res["TIME"]) != 0 {
		t.Fatalf("DAY should be an array of 0 elements, but is '%s'", res["TIME"])
	}

	g, _ = New()
	res, _ = g.ParseToMultiMap("%{COMMONAPACHELOG} %{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:23:58:32 +0200] "GET /index.php HTTP/1.1" 404 207 127.0.0.1 - - [24/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	if len(res["TIME"]) != 2 {
		t.Fatalf("TIME should be an array of 2 elements, but is '%s'", res["TIME"])
	}
	if len(res["timestamp"]) != 2 {
		t.Fatalf("timestamp should be an array of 2 elements, but is '%s'", res["timestamp"])
	}
}

func TestNewWithNoDefaultPatterns(t *testing.T) {
	g, _ := NewWithConfig(&Config{SkipDefaultPatterns: true})
	if len(g.patterns) != 0 {
		t.Fatal("Using SkipDefaultPatterns the Grok object should not have any patterns pre loaded")
	}
}

func TestAddPatternErr(t *testing.T) {
	name := "Error"
	pattern := "%{ERR}"

	g, _ := New()
	err := g.addPattern(name, pattern)
	if err == nil {
		t.Fatalf("AddPattern should returns an error when path is invalid")
	}
}

func TestAddPatternsFromPathErr(t *testing.T) {
	g, _ := New()
	err := g.AddPatternsFromPath("./Lorem ipsum Minim qui in.")
	if err == nil {
		t.Fatalf("AddPatternsFromPath should returns an error when path is invalid")
	}
}

func TestConfigPatternsDir(t *testing.T) {
	g, err := NewWithConfig(&Config{PatternsDir: []string{"./patterns"}})
	if err != nil {
		t.Error(err)
	}

	if captures, err := g.Parse("%{SYSLOGLINE}", `Sep 12 23:19:02 docker syslog-ng[25389]: syslog-ng starting up; version='3.5.3'`); err != nil {
		t.Fatalf("error : %s", err.Error())
	} else {
		// pp.Print(captures)
		if captures["program"] != "syslog-ng" {
			t.Fatalf("%s should be '%s' have '%s'", "program", "syslog-ng", captures["program"])
		}
	}

}

func TestAddPatternsFromPathFileOpenErr(t *testing.T) {
	t.Skipped()
}

func TestAddPatternsFromPathFile(t *testing.T) {
	g, _ := New()
	err := g.AddPatternsFromPath("./patterns/grok-patterns")
	if err != nil {
		t.Fatalf("err %#v", err)
	}
}

func TestAddPattern(t *testing.T) {
	name := "DAYO"
	pattern := "(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)"

	g, _ := New()
	cPatterns := len(g.patterns)
	g.AddPattern(name, pattern)
	g.AddPattern(name+"2", pattern)
	if len(g.patterns) != cPatterns+2 {
		t.Fatalf("%d Default patterns should be available, have %d", cPatterns+2, len(g.patterns))
	}

	g, _ = NewWithConfig(&Config{NamedCapturesOnly: true})
	cPatterns = len(g.patterns)
	g.AddPattern(name, pattern)
	g.AddPattern(name+"2", pattern)
	if len(g.patterns) != cPatterns+2 {
		t.Fatalf("%d NamedCapture patterns should be available, have %d", cPatterns+2, len(g.patterns))
	}
}

func TestMatch(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")

	if r, err := g.Match("%{MONTH}", "June"); !r {
		t.Fatalf("June should match %s: err=%s", "%{MONTH}", err.Error())
	}

}
func TestDoesNotMatch(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")
	if r, _ := g.Match("%{MONTH}", "13"); r {
		t.Fatalf("13 should not match %s", "%{MONTH}")
	}
}

func TestErrorMatch(t *testing.T) {
	g, _ := New()
	if _, err := g.Match("(", "13"); err == nil {
		t.Fatal("Error expected")
	}

}

func TestShortName(t *testing.T) {
	g, _ := New()
	g.AddPattern("A", "a")

	m, err := g.Match("%{A}", "a")
	if err != nil {
		t.Fatal("a should match %%{A}: err=%s", err.Error())
	}
	if !m {
		t.Fatal("%%{A} didn't match 'a'")
	}
}

func TestDayCompile(t *testing.T) {
	g, _ := New()
	g.AddPattern("DAY", "(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)")
	pattern := "%{DAY}"
	_, err := g.compile(pattern)
	if err != nil {
		t.Fatal("Error:", err)
	}
}

func TestErrorCompile(t *testing.T) {
	g, _ := New()
	_, err := g.compile("(")
	if err == nil {
		t.Fatal("Error:", err)
	}
}

func TestNamedCaptures(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")

	check := func(key, value, pattern, text string) {
		captures, _ := g.Parse(pattern, text)
		if captures[key] != value {
			t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
		}
	}

	check("jour", "Tue",
		"%{DAY:jour}",
		"Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157",
	)

	check("day-of.week", "Tue",
		"%{DAY:day-of.week}",
		"Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157",
	)
}

func TestErrorCaptureUnknowPattern(t *testing.T) {
	g, _ := New()
	pattern := "%{UNKNOWPATTERN}"
	_, err := g.Parse(pattern, "")
	if err == nil {
		t.Fatal("Expected error not set")
	}
}

func TestParse(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")
	res, _ := g.Parse("%{DAY}", "Tue qds")
	if res["DAY"] != "Tue" {
		t.Fatalf("DAY should be 'Tue' have '%s'", res["DAY"])
	}
}

func TestErrorParseToMultiMap(t *testing.T) {
	g, _ := New()
	pattern := "%{UNKNOWPATTERN}"
	_, err := g.ParseToMultiMap(pattern, "")
	if err == nil {
		t.Fatal("Expected error not set")
	}
}

func TestParseToMultiMap(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")
	res, _ := g.ParseToMultiMap("%{COMMONAPACHELOG} %{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:23:58:32 +0200] "GET /index.php HTTP/1.1" 404 207 127.0.0.1 - - [24/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	if len(res["TIME"]) != 2 {
		t.Fatalf("DAY should be an array of 3 elements, but is '%s'", res["TIME"])
	}
	if res["TIME"][0] != "23:58:32" {
		t.Fatalf("TIME[0] should be '23:58:32' have '%s'", res["TIME"][0])
	}
	if res["TIME"][1] != "22:58:32" {
		t.Fatalf("TIME[1] should be '22:58:32' have '%s'", res["TIME"][1])
	}
}

func TestParseToMultiMapOnlyNamedCaptures(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	g.AddPatternsFromPath("./patterns")
	res, _ := g.ParseToMultiMap("%{COMMONAPACHELOG} %{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207 127.0.0.1 - - [24/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	if len(res["timestamp"]) != 2 {
		t.Fatalf("timestamp should be an array of 2 elements, but is '%s'", res["timestamp"])
	}
	if res["timestamp"][0] != "23/Apr/2014:22:58:32 +0200" {
		t.Fatalf("timestamp[0] should be '23/Apr/2014:22:58:32 +0200' have '%s'", res["DAY"][0])
	}
	if res["timestamp"][1] != "24/Apr/2014:22:58:32 +0200" {
		t.Fatalf("timestamp[1] should be '24/Apr/2014:22:58:32 +0200' have '%s'", res["DAY"][1])
	}
}

func TestCaptureAll(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")

	check := func(key, value, pattern, text string) {

		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if captures[key] != value {
				t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
			}
		}
	}

	check("timestamp", "23/Apr/2014:22:58:32 +0200",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("TIME", "22:58:32",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("SECOND", `17,1599`, "%{TIMESTAMP_ISO8601}", `s d9fq4999s ../ sdf 2013-11-06 04:50:17,1599sd`)
	check("HOSTNAME", `google.com`, "%{HOSTPORT}", `google.com:8080`)
	//HOSTPORT
	check("POSINT", `8080`, "%{HOSTPORT}", `google.com:8080`)
}

func TestNamedCapture(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	g.AddPatternsFromPath("./patterns")

	check := func(key, value, pattern, text string) {
		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if captures[key] != value {
				t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
			}
		}
	}

	check("timestamp", "23/Apr/2014:22:58:32 +0200",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("TIME", "",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("SECOND", ``, "%{TIMESTAMP_ISO8601}", `s d9fq4999s ../ sdf 2013-11-06 04:50:17,1599sd`)
	check("HOSTNAME", ``, "%{HOSTPORT}", `google.com:8080`)
	//HOSTPORT
	check("POSINT", ``, "%{HOSTPORT}", `google.com:8080`)
}

func TestRemoveEmptyValues(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true, RemoveEmptyValues: true})

	capturesExists := func(key, pattern, text string) {
		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if _, ok := captures[key]; ok {
				t.Fatalf("%s should be absent", key)
			}
		}
	}

	capturesExists("rawrequest", "%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)

}

func TestCapturesAndNamedCapture(t *testing.T) {

	check := func(key, value, pattern, text string) {
		g, _ := New()
		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if captures[key] != value {
				t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
			}
		}
	}

	checkNamed := func(key, value, pattern, text string) {
		g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if captures[key] != value {
				t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
			}
		}

	}

	check("DAY", "Tue",
		"%{DAY}",
		"Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157",
	)
	checkNamed("jour", "Tue",
		"%{DAY:jour}",
		"Tue May 15 11:21:42 [conn1047685] moveChunk deleted: 7157",
	)
	check("clientip", "127.0.0.1",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("verb", "GET",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("timestamp", "23/Apr/2014:22:58:32 +0200",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)
	check("bytes", "207",
		"%{COMMONAPACHELOG}",
		`127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`,
	)

	//PATH
	check("WINPATH", `c:\winfows\sdf.txt`, "%{WINPATH}", `s dfqs c:\winfows\sdf.txt`)
	check("WINPATH", `\\sdf\winfows\sdf.txt`, "%{WINPATH}", `s dfqs \\sdf\winfows\sdf.txt`)
	check("UNIXPATH", `/usr/lib/`, "%{UNIXPATH}", `s dfqs /usr/lib/ sqfd`)
	check("UNIXPATH", `/usr/lib`, "%{UNIXPATH}", `s dfqs /usr/lib sqfd`)
	check("UNIXPATH", `/usr/`, "%{UNIXPATH}", `s dfqs /usr/ sqfd`)
	check("UNIXPATH", `/usr`, "%{UNIXPATH}", `s dfqs /usr sqfd`)
	check("UNIXPATH", `/`, "%{UNIXPATH}", `s dfqs / sqfd`)

	//YEAR
	check("YEAR", `4999`, "%{YEAR}", `s d9fq4999s ../ sdf`)
	check("YEAR", `79`, "%{YEAR}", `s d79fq4999s ../ sdf`)
	check("TIMESTAMP_ISO8601", `2013-11-06 04:50:17,1599`, "%{TIMESTAMP_ISO8601}", `s d9fq4999s ../ sdf 2013-11-06 04:50:17,1599sd`)

	//MAC
	check("MAC", `01:02:03:04:ab:cf`, "%{MAC}", `s d9fq4999s ../ sdf 2013- 01:02:03:04:ab:cf  11-06 04:50:17,1599sd`)
	check("MAC", `01-02-03-04-ab-cd`, "%{MAC}", `s d9fq4999s ../ sdf 2013- 01-02-03-04-ab-cd  11-06 04:50:17,1599sd`)

	//QUOTEDSTRING
	check("QUOTEDSTRING", `"lkj"`, "%{QUOTEDSTRING}", `qsdklfjqsd fk"lkj"mkj`)
	check("QUOTEDSTRING", `'lkj'`, "%{QUOTEDSTRING}", `qsdklfjqsd fk'lkj'mkj`)
	check("QUOTEDSTRING", `"fk'lkj'm"`, "%{QUOTEDSTRING}", `qsdklfjqsd "fk'lkj'm"kj`)
	check("QUOTEDSTRING", `'fk"lkj"m'`, "%{QUOTEDSTRING}", `qsdklfjqsd 'fk"lkj"m'kj`)

	//BASE10NUM
	check("BASE10NUM", `1`, "%{BASE10NUM}", `1`) // this is a nice one
	check("BASE10NUM", `8080`, "%{BASE10NUM}", `qsfd8080qsfd`)

}

// Should be run with -race
func TestConcurentParse(t *testing.T) {
	g, _ := New()
	g.AddPatternsFromPath("./patterns")

	check := func(key, value, pattern, text string) {

		if captures, err := g.Parse(pattern, text); err != nil {
			t.Fatalf("error can not capture : %s", err.Error())
		} else {
			if captures[key] != value {
				t.Fatalf("%s should be '%s' have '%s'", key, value, captures[key])
			}
		}
	}

	go check("QUOTEDSTRING", `"lkj"`, "%{QUOTEDSTRING}", `qsdklfjqsd fk"lkj"mkj`)
	go check("QUOTEDSTRING", `'lkj'`, "%{QUOTEDSTRING}", `qsdklfjqsd fk'lkj'mkj`)
	go check("QUOTEDSTRING", `'lkj'`, "%{QUOTEDSTRING}", `qsdklfjqsd fk'lkj'mkj`)
	go check("QUOTEDSTRING", `"fk'lkj'm"`, "%{QUOTEDSTRING}", `qsdklfjqsd "fk'lkj'm"kj`)
	go check("QUOTEDSTRING", `'fk"lkj"m'`, "%{QUOTEDSTRING}", `qsdklfjqsd 'fk"lkj"m'kj`)
}

func TestPatterns(t *testing.T) {
	g, _ := NewWithConfig(&Config{SkipDefaultPatterns: true})
	if len(g.patterns) != 0 {
		t.Fatalf("Patterns should return 0, have '%d'", len(g.patterns))
	}
	name := "DAY0"
	pattern := "(?:Mon(?:day)?|Tue(?:sday)?|Wed(?:nesday)?|Thu(?:rsday)?|Fri(?:day)?|Sat(?:urday)?|Sun(?:day)?)"

	g.AddPattern(name, pattern)
	g.AddPattern(name+"1", pattern)
	if len(g.patterns) != 2 {
		t.Fatalf("Patterns should return 2, have '%d'", len(g.patterns))
	}
}

func TestParseTypedWithDefaultCaptureMode(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	if captures, err := g.ParseTyped("%{IPV4:ip:string} %{NUMBER:status:int} %{NUMBER:duration:float}", `127.0.0.1 200 0.8`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["ip"] != "127.0.0.1" {
			t.Fatalf("%s should be '%s' have '%s'", "ip", "127.0.0.1", captures["ip"])
		} else {
			if captures["status"] != 200 {
				t.Fatalf("%s should be '%d' have '%d'", "status", 200, captures["status"])
			} else {
				if captures["duration"] != 0.8 {
					t.Fatalf("%s should be '%f' have '%f'", "duration", 0.8, captures["duration"])
				}
			}
		}
	}
}

func TestParseTypedWithNoTypeInfo(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	if captures, err := g.ParseTyped("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["timestamp"] != "23/Apr/2014:22:58:32 +0200" {
			t.Fatalf("%s should be '%s' have '%s'", "timestamp", "23/Apr/2014:22:58:32 +0200", captures["timestamp"])
		}
		if captures["TIME"] != nil {
			t.Fatalf("%s should be nil have '%s'", "TIME", captures["TIME"])
		}
	}

	g, _ = New()
	if captures, err := g.ParseTyped("%{COMMONAPACHELOG}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["timestamp"] != "23/Apr/2014:22:58:32 +0200" {
			t.Fatalf("%s should be '%s' have '%s'", "timestamp", "23/Apr/2014:22:58:32 +0200", captures["timestamp"])
		}
		if captures["TIME"] != "22:58:32" {
			t.Fatalf("%s should be '%s' have '%s'", "TIME", "22:58:32", captures["TIME"])
		}
	}
}

func TestParseTypedWithIntegerTypeCoercion(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	if captures, err := g.ParseTyped("%{WORD:coerced:int}", `5.75`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["coerced"] != 5 {
			t.Fatalf("%s should be '%s' have '%s'", "coerced", "5", captures["coerced"])
		}
	}
}

func TestParseTypedWithUnknownType(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	if _, err := g.ParseTyped("%{WORD:word:unknown}", `hello`); err == nil {
		t.Fatalf("parsing an unknown type must result in a conversion error")
	}
}

func TestParseTypedErrorCaptureUnknowPattern(t *testing.T) {
	g, _ := New()
	pattern := "%{UNKNOWPATTERN}"
	_, err := g.ParseTyped(pattern, "")
	if err == nil {
		t.Fatal("Expected error not set")
	}
}

func TestParseTypedWithTypedParents(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	g.AddPattern("TESTCOMMON", `%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes:int}|-)`)
	if captures, err := g.ParseTyped("%{TESTCOMMON}", `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["bytes"] != 207 {
			t.Fatalf("%s should be '%s' have '%s'", "bytes", "207", captures["bytes"])
		}
	}
}

func TestParseTypedWithSemanticHomonyms(t *testing.T) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true, SkipDefaultPatterns: true})

	g.AddPattern("BASE10NUM", `([+-]?(?:[0-9]+(?:\.[0-9]+)?)|\.[0-9]+)`)
	g.AddPattern("NUMBER", `(?:%{BASE10NUM})`)
	g.AddPattern("MYNUM", `%{NUMBER:bytes:int}`)
	g.AddPattern("MYSTR", `%{NUMBER:bytes:string}`)

	if captures, err := g.ParseTyped("%{MYNUM}", `207`); err != nil {
		t.Fatalf("error can not scapture : %s", err.Error())
	} else {
		if captures["bytes"] != 207 {
			t.Fatalf("%s should be %#v have %#v", "bytes", 207, captures["bytes"])
		}
	}
	if captures, err := g.ParseTyped("%{MYSTR}", `207`); err != nil {
		t.Fatalf("error can not capture : %s", err.Error())
	} else {
		if captures["bytes"] != "207" {
			t.Fatalf("%s should be %#v have %#v", "bytes", "207", captures["bytes"])
		}
	}
}

var resultNew *Grok

func BenchmarkNew(b *testing.B) {
	b.ReportAllocs()
	b.ResetTimer()
	var g *Grok
	// run the check function b.N times
	for n := 0; n < b.N; n++ {
		g, _ = NewWithConfig(&Config{NamedCapturesOnly: true})
	}
	resultNew = g
}

func BenchmarkCaptures(b *testing.B) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	b.ReportAllocs()
	b.ResetTimer()
	// run the check function b.N times
	for n := 0; n < b.N; n++ {
		g.Parse(`%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes}|-)`, `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	}
}

func BenchmarkCapturesTypedFake(b *testing.B) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	b.ReportAllocs()
	b.ResetTimer()
	// run the check function b.N times
	for n := 0; n < b.N; n++ {
		g.Parse(`%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes}|-)`, `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	}
}

func BenchmarkCapturesTypedReal(b *testing.B) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	b.ReportAllocs()
	b.ResetTimer()
	// run the check function b.N times
	for n := 0; n < b.N; n++ {
		g.ParseTyped(`%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion:int})?|%{DATA:rawrequest})" %{NUMBER:response:int} (?:%{NUMBER:bytes:int}|-)`, `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
	}
}

func BenchmarkParallelCaptures(b *testing.B) {
	g, _ := NewWithConfig(&Config{NamedCapturesOnly: true})
	b.ReportAllocs()
	b.ResetTimer()

	b.RunParallel(func(b *testing.PB) {
		for b.Next() {
			g.Parse(`%{IPORHOST:clientip} %{USER:ident} %{USER:auth} \[%{HTTPDATE:timestamp}\] "(?:%{WORD:verb} %{NOTSPACE:request}(?: HTTP/%{NUMBER:httpversion})?|%{DATA:rawrequest})" %{NUMBER:response} (?:%{NUMBER:bytes}|-)`, `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207`)
		}
	})
}

func TestGrok_AddPatternsFromMap_not_exist(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddPatternsFromMap panics: %v", r)
		}
	}()
	g, _ := NewWithConfig(&Config{SkipDefaultPatterns: true})
	err := g.AddPatternsFromMap(map[string]string{
		"SOME": "%{NOT_EXIST}",
	})
	if err == nil {
		t.Errorf("AddPatternsFromMap should returns an error")
	}
}

func TestGrok_AddPatternsFromMap_simple(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddPatternsFromMap panics: %v", r)
		}
	}()
	g, _ := NewWithConfig(&Config{SkipDefaultPatterns: true})
	err := g.AddPatternsFromMap(map[string]string{
		"NO3": `\d{3}`,
	})
	if err != nil {
		t.Errorf("AddPatternsFromMap returns an error: %v", err)
	}
	mss, err := g.Parse("%{NO3:match}", "333")
	if err != nil {
		t.Error("parsing error:", err)
		t.FailNow()
	}
	if mss["match"] != "333" {
		t.Errorf("bad match: expected 333, got %s", mss["match"])
	}
}

func TestGrok_AddPatternsFromMap_complex(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("AddPatternsFromMap panics: %v", r)
		}
	}()
	g, _ := NewWithConfig(&Config{
		SkipDefaultPatterns: true,
		NamedCapturesOnly:   true,
	})
	err := g.AddPatternsFromMap(map[string]string{
		"NO3": `\d{3}`,
		"NO6": "%{NO3}%{NO3}",
	})
	if err != nil {
		t.Errorf("AddPatternsFromMap returns an error: %v", err)
	}
	mss, err := g.Parse("%{NO6:number}", "333666")
	if err != nil {
		t.Error("parsing error:", err)
		t.FailNow()
	}
	if mss["number"] != "333666" {
		t.Errorf("bad match: expected 333666, got %s", mss["match"])
	}
}

func TestParseStream(t *testing.T) {
	g, _ := New()
	pTest := func(m map[string]string) error {
		ts, ok := m["timestamp"]
		if !ok {
			t.Error("timestamp not found")
		}
		if len(ts) == 0 {
			t.Error("empty timestamp")
		}
		return nil
	}
	const testLog = `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207
127.0.0.1 - - [23/Apr/2014:22:59:32 +0200] "GET /index.php HTTP/1.1" 404 207
127.0.0.1 - - [23/Apr/2014:23:00:32 +0200] "GET /index.php HTTP/1.1" 404 207
`

	r := bufio.NewReader(strings.NewReader(testLog))
	if err := g.ParseStream(r, "%{COMMONAPACHELOG}", pTest); err != nil {
		t.Fatal(err)
	}
}

func TestParseStreamError(t *testing.T) {
	g, _ := New()
	pTest := func(m map[string]string) error {
		if _, ok := m["timestamp"]; !ok {
			return fmt.Errorf("timestamp not found")
		}
		return nil
	}
	const testLog = `127.0.0.1 - - [23/Apr/2014:22:58:32 +0200] "GET /index.php HTTP/1.1" 404 207
127.0.0.1 - - [xxxxxxxxxxxxxxxxxxxx +0200] "GET /index.php HTTP/1.1" 404 207
127.0.0.1 - - [23/Apr/2014:23:00:32 +0200] "GET /index.php HTTP/1.1" 404 207
`

	r := bufio.NewReader(strings.NewReader(testLog))
	if err := g.ParseStream(r, "%{COMMONAPACHELOG}", pTest); err == nil {
		t.Fatal("Error expected")
	}
}
