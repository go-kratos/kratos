package broker

// Message represent a message
type Message struct {
	Header []Header
	Body   []byte
	Raw    interface{}
}

// Header is a key/value pair type representing headers of a message
type Header struct {
	Key   string
	Value []byte
}
