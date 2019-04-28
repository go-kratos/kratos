package bytes

// Writer writer.
type Writer struct {
	n   int
	buf []byte
}

// NewWriterSize new a writer with size.
func NewWriterSize(n int) *Writer {
	return &Writer{buf: make([]byte, n)}
}

// Len buff len.
func (w *Writer) Len() int {
	return w.n
}

// Size buff cap.
func (w *Writer) Size() int {
	return len(w.buf)
}

// Reset reset the buff.
func (w *Writer) Reset() {
	w.n = 0
}

// Buffer return buff.
func (w *Writer) Buffer() []byte {
	return w.buf[:w.n]
}

// Peek peek a buf.
func (w *Writer) Peek(n int) []byte {
	var buf []byte
	w.grow(n)
	buf = w.buf[w.n : w.n+n]
	w.n += n
	return buf
}

// Write write a buff.
func (w *Writer) Write(p []byte) {
	w.grow(len(p))
	w.n += copy(w.buf[w.n:], p)
}

func (w *Writer) grow(n int) {
	var buf []byte
	if w.n+n < len(w.buf) {
		return
	}
	buf = make([]byte, 2*len(w.buf)+n)
	copy(buf, w.buf[:w.n])
	w.buf = buf
}
