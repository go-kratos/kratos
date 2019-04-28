package log

import (
	"context"
	"io"
	"path/filepath"
	"time"

	"go-common/library/log/internal/filewriter"
)

// level idx
const (
	_infoIdx = iota
	_warnIdx
	_errorIdx
	_totalIdx
)

var _fileNames = map[int]string{
	_infoIdx:  "info.log",
	_warnIdx:  "warning.log",
	_errorIdx: "error.log",
}

// FileHandler .
type FileHandler struct {
	render Render
	fws    [_totalIdx]*filewriter.FileWriter
}

// NewFile crete a file logger.
func NewFile(dir string, bufferSize, rotateSize int64, maxLogFile int) *FileHandler {
	// new info writer
	newWriter := func(name string) *filewriter.FileWriter {
		var options []filewriter.Option
		if rotateSize > 0 {
			options = append(options, filewriter.MaxSize(rotateSize))
		}
		if maxLogFile > 0 {
			options = append(options, filewriter.MaxFile(maxLogFile))
		}
		w, err := filewriter.New(filepath.Join(dir, name), options...)
		if err != nil {
			panic(err)
		}
		return w
	}
	handler := &FileHandler{
		render: newPatternRender("[%D %T] [%L] [%S] %M"),
	}
	for idx, name := range _fileNames {
		handler.fws[idx] = newWriter(name)
	}
	return handler
}

// Log loggint to file .
func (h *FileHandler) Log(ctx context.Context, lv Level, args ...D) {
	d := make(map[string]interface{}, 10+len(args))
	for _, arg := range args {
		d[arg.Key] = arg.Value
	}
	// add extra fields
	addExtraField(ctx, d)
	d[_time] = time.Now().Format(_timeFormat)
	var w io.Writer
	switch lv {
	case _warnLevel:
		w = h.fws[_warnIdx]
	case _errorLevel:
		w = h.fws[_errorIdx]
	default:
		w = h.fws[_infoIdx]
	}
	h.render.Render(w, d)
	w.Write([]byte("\n"))
}

// Close log handler
func (h *FileHandler) Close() error {
	for _, fw := range h.fws {
		// ignored error
		fw.Close()
	}
	return nil
}

// SetFormat set log format
func (h *FileHandler) SetFormat(format string) {
	h.render = newPatternRender(format)
}
