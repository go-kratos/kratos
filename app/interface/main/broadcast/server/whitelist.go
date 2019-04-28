package server

import (
	"go-common/app/interface/main/broadcast/conf"
	"log"
	"os"
)

var whitelist *Whitelist

// Whitelist .
type Whitelist struct {
	log  *log.Logger
	list map[int64]struct{} // whitelist for debug
}

// InitWhitelist a whitelist struct.
func InitWhitelist(c *conf.Whitelist) (err error) {
	var (
		mid int64
		f   *os.File
	)
	if f, err = os.OpenFile(c.WhiteLog, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0644); err == nil {
		whitelist = new(Whitelist)
		whitelist.log = log.New(f, "", log.LstdFlags)
		whitelist.list = make(map[int64]struct{})
		for _, mid = range c.Whitelist {
			whitelist.list[mid] = struct{}{}
		}
	}
	return
}

// Contains whitelist contains a mid or not.
func (w *Whitelist) Contains(mid int64) (ok bool) {
	if mid > 0 {
		_, ok = w.list[mid]
	}
	return
}

// Printf calls l.Output to print to the logger.
func (w *Whitelist) Printf(format string, v ...interface{}) {
	w.log.Printf(format, v...)
}
