package filewriter

import (
	"bytes"
	"container/list"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// FileWriter create file log writer
type FileWriter struct {
	opt    option
	dir    string
	fname  string
	ch     chan *bytes.Buffer
	stdlog *log.Logger
	pool   *sync.Pool

	lastRotateFormat string
	lastSplitNum     int

	current *wrapFile
	files   *list.List

	closed int32
	wg     sync.WaitGroup
}

type rotateItem struct {
	rotateTime int64
	rotateNum  int
	fname      string
}

func parseRotateItem(dir, fname, rotateFormat string) (*list.List, error) {
	fis, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// parse exists log file filename
	parse := func(s string) (rt rotateItem, err error) {
		// remove filename and left "." error.log.2018-09-12.001 -> 2018-09-12.001
		rt.fname = s
		s = strings.TrimLeft(s[len(fname):], ".")
		seqs := strings.Split(s, ".")
		var t time.Time
		switch len(seqs) {
		case 2:
			if rt.rotateNum, err = strconv.Atoi(seqs[1]); err != nil {
				return
			}
			fallthrough
		case 1:
			if t, err = time.Parse(rotateFormat, seqs[0]); err != nil {
				return
			}
			rt.rotateTime = t.Unix()
		}
		return
	}

	var items []rotateItem
	for _, fi := range fis {
		if strings.HasPrefix(fi.Name(), fname) && fi.Name() != fname {
			rt, err := parse(fi.Name())
			if err != nil {
				// TODO deal with error
				continue
			}
			items = append(items, rt)
		}
	}
	sort.Slice(items, func(i, j int) bool {
		if items[i].rotateTime == items[j].rotateTime {
			return items[i].rotateNum > items[j].rotateNum
		}
		return items[i].rotateTime > items[j].rotateTime
	})
	l := list.New()

	for _, item := range items {
		l.PushBack(item)
	}
	return l, nil
}

type wrapFile struct {
	fsize int64
	fp    *os.File
}

func (w *wrapFile) size() int64 {
	return w.fsize
}

func (w *wrapFile) write(p []byte) (n int, err error) {
	n, err = w.fp.Write(p)
	w.fsize += int64(n)
	return
}

func newWrapFile(fpath string) (*wrapFile, error) {
	fp, err := os.OpenFile(fpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	fi, err := fp.Stat()
	if err != nil {
		return nil, err
	}
	return &wrapFile{fp: fp, fsize: fi.Size()}, nil
}

// New FileWriter A FileWriter is safe for use by multiple goroutines simultaneously.
func New(fpath string, fns ...Option) (*FileWriter, error) {
	opt := defaultOption
	for _, fn := range fns {
		fn(&opt)
	}

	fname := filepath.Base(fpath)
	if fname == "" {
		return nil, fmt.Errorf("filename can't empty")
	}
	dir := filepath.Dir(fpath)
	fi, err := os.Stat(dir)
	if err == nil && !fi.IsDir() {
		return nil, fmt.Errorf("%s already exists and not a directory", dir)
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create dir %s error: %s", dir, err.Error())
		}
	}

	current, err := newWrapFile(fpath)
	if err != nil {
		return nil, err
	}

	stdlog := log.New(os.Stderr, "flog ", log.LstdFlags)
	ch := make(chan *bytes.Buffer, opt.ChanSize)

	files, err := parseRotateItem(dir, fname, opt.RotateFormat)
	if err != nil {
		// set files a empty list
		files = list.New()
		stdlog.Printf("parseRotateItem error: %s", err)
	}

	lastRotateFormat := time.Now().Format(opt.RotateFormat)
	var lastSplitNum int
	if files.Len() > 0 {
		rt := files.Front().Value.(rotateItem)
		//  check contains is mush esay than compared with timestamp
		if strings.Contains(rt.fname, lastRotateFormat) {
			lastSplitNum = rt.rotateNum
		}
	}

	fw := &FileWriter{
		opt:    opt,
		dir:    dir,
		fname:  fname,
		stdlog: stdlog,
		ch:     ch,
		pool:   &sync.Pool{New: func() interface{} { return new(bytes.Buffer) }},

		lastSplitNum:     lastSplitNum,
		lastRotateFormat: lastRotateFormat,

		files:   files,
		current: current,
	}

	fw.wg.Add(1)
	go fw.daemon()

	return fw, nil
}

// Write write data to log file, return write bytes is pseudo just for implement io.Writer.
func (f *FileWriter) Write(p []byte) (int, error) {
	// atomic is not necessary
	if atomic.LoadInt32(&f.closed) == 1 {
		f.stdlog.Printf("%s", p)
		return 0, fmt.Errorf("filewriter already closed")
	}
	// because write to file is asynchronousc,
	// copy p to internal buf prevent p be change on outside
	buf := f.getBuf()
	buf.Write(p)

	if f.opt.WriteTimeout == 0 {
		select {
		case f.ch <- buf:
			return len(p), nil
		default:
			// TODO: write discard log to to stdout?
			return 0, fmt.Errorf("log channel is full, discard log")
		}
	}

	// write log with timeout
	timeout := time.NewTimer(f.opt.WriteTimeout)
	select {
	case f.ch <- buf:
		return len(p), nil
	case <-timeout.C:
		// TODO: write discard log to to stdout?
		return 0, fmt.Errorf("log channel is full, discard log")
	}
}

func (f *FileWriter) daemon() {
	// TODO: check aggsbuf size prevent it too big
	aggsbuf := &bytes.Buffer{}
	tk := time.NewTicker(f.opt.RotateInterval)
	// TODO: make it configrable
	aggstk := time.NewTicker(10 * time.Millisecond)
	var err error
	for {
		select {
		case t := <-tk.C:
			f.checkRotate(t)
		case buf, ok := <-f.ch:
			if ok {
				aggsbuf.Write(buf.Bytes())
				f.putBuf(buf)
			}
		case <-aggstk.C:
			if aggsbuf.Len() > 0 {
				if err = f.write(aggsbuf.Bytes()); err != nil {
					f.stdlog.Printf("write log error: %s", err)
				}
				aggsbuf.Reset()
			}
		}
		if atomic.LoadInt32(&f.closed) != 1 {
			continue
		}
		// read all buf from channel and break loop
		if err = f.write(aggsbuf.Bytes()); err != nil {
			f.stdlog.Printf("write log error: %s", err)
		}
		for buf := range f.ch {
			if err = f.write(buf.Bytes()); err != nil {
				f.stdlog.Printf("write log error: %s", err)
			}
			f.putBuf(buf)
		}
		break
	}
	f.wg.Done()
}

// Close close file writer
func (f *FileWriter) Close() error {
	atomic.StoreInt32(&f.closed, 1)
	close(f.ch)
	f.wg.Wait()
	return nil
}

func (f *FileWriter) checkRotate(t time.Time) {
	formatFname := func(format string, num int) string {
		if num == 0 {
			return fmt.Sprintf("%s.%s", f.fname, format)
		}
		return fmt.Sprintf("%s.%s.%03d", f.fname, format, num)
	}
	format := t.Format(f.opt.RotateFormat)

	if f.opt.MaxFile != 0 {
		for f.files.Len() > f.opt.MaxFile {
			rt := f.files.Remove(f.files.Front()).(rotateItem)
			fpath := filepath.Join(f.dir, rt.fname)
			if err := os.Remove(fpath); err != nil {
				f.stdlog.Printf("remove file %s error: %s", fpath, err)
			}
		}
	}

	if format != f.lastRotateFormat || (f.opt.MaxSize != 0 && f.current.size() > f.opt.MaxSize) {
		var err error
		// close current file first
		if err = f.current.fp.Close(); err != nil {
			f.stdlog.Printf("close current file error: %s", err)
		}

		// rename file
		fname := formatFname(f.lastRotateFormat, f.lastSplitNum)
		oldpath := filepath.Join(f.dir, f.fname)
		newpath := filepath.Join(f.dir, fname)
		if err = os.Rename(oldpath, newpath); err != nil {
			f.stdlog.Printf("rename file %s to %s error: %s", oldpath, newpath, err)
			return
		}

		f.files.PushBack(rotateItem{fname: fname /*rotateNum: f.lastSplitNum, rotateTime: t.Unix() unnecessary*/})

		if format != f.lastRotateFormat {
			f.lastRotateFormat = format
			f.lastSplitNum = 0
		} else {
			f.lastSplitNum++
		}

		// recreate current file
		f.current, err = newWrapFile(filepath.Join(f.dir, f.fname))
		if err != nil {
			f.stdlog.Printf("create log file error: %s", err)
		}
	}
}

func (f *FileWriter) write(p []byte) error {
	// f.current may be nil, if newWrapFile return err in checkRotate, redirect log to stderr
	if f.current == nil {
		f.stdlog.Printf("can't write log to file, please check stderr log for detail")
		f.stdlog.Printf("%s", p)
	}
	_, err := f.current.write(p)
	return err
}

func (f *FileWriter) putBuf(buf *bytes.Buffer) {
	buf.Reset()
	f.pool.Put(buf)
}

func (f *FileWriter) getBuf() *bytes.Buffer {
	return f.pool.Get().(*bytes.Buffer)
}
