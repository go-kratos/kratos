package file

import (
	"io"
	"bufio"
	"bytes"
)

// Message represents a reader event with timestamp, content and actual number
// of bytes read from input before decoding.
//type Message struct {
//	Ts      time.Time // timestamp the content was read
//	Content []byte    // actual content read
//	Bytes   int       // total number of bytes read to generate the message
//	//Fields  common.MapStr // optional fields that can be added by reader
//}

type Reader interface {
	Next() ([]byte, int, error)
}

type LineReader struct {
	reader     io.Reader
	rb         *bufio.Reader
	bufferSize int
	nl         []byte
	nlSize     int
	scan       *bufio.Scanner
}

// New creates a new reader object
func NewLineReader(input io.Reader, bufferSize int) (*LineReader, error) {
	nl := []byte{'\n'}

	r := &LineReader{
		reader:     input,
		bufferSize: bufferSize,
		nl:         nl,
		nlSize:     len(nl),
	}
	r.rb = bufio.NewReaderSize(input, r.bufferSize)
	r.scan = bufio.NewScanner(r.rb)
	r.scan.Split(ScanLines)

	return r, nil
}

func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, data[0:i], nil
	}

	// Request more data.
	return 0, nil, nil
}

// Next reads the next line until the new line character
func (r *LineReader) Next() ([]byte, int, error) {
	body, err := r.rb.ReadBytes('\n')
	advance := len(body)
	//if err == io.EOF && advance > 0 {
	//	return body, advance, err
	//}
	// remove '\n'
	if len(body) > 0 && body[len(body)-1] == '\n' {
		body = body[0:len(body)-1]
	}

	// remove '\r'
	if len(body) > 0 && body[len(body)-1] == '\r' {
		body = body[0: len(body)-1]
	}

	return body, advance, err
}
