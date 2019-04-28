package regexp

import (
	"go-common/library/log"

	pcre "github.com/GRbit/go-pcre"
)

type Regexp struct {
	regexp pcre.Regexp
	expr   string
}

func Compile(expr string) (reg *Regexp, err error) {
	if reg = cache.Get(expr); reg != nil {
		log.Info("regexp cache hit (%s)", expr)
		return
	}
	reg = &Regexp{expr: expr}
	if reg.regexp, err = pcre.Compile(expr, pcre.UTF8); err != nil {
		reg = nil
		return
	}
	cache.Set(expr, reg)
	return
}

func (r *Regexp) FindAllStringIndex(str string, n int) (pos [][]int) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("FindAllStringIndex error(%v)", x)
		}
	}()
	return r.FindAllIndex([]byte(str))
}

func (r *Regexp) FindAllIndex(b []byte) (pos [][]int) {
	defer func() {
		if x := recover(); x != nil {
			log.Error("FindAllIndex error(%v)", x)
		}
	}()
	return r.regexp.FindAllIndex(b, 0)
}

func (r *Regexp) MatchString(msg string) bool {
	defer func() {
		if x := recover(); x != nil {
			log.Error("MatchString error(%v)", x)
		}
	}()
	return r.regexp.MatcherString(msg, 0).Matches
}

func (r *Regexp) String() string {
	return r.expr
}
