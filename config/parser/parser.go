package parser

// Parser is config parser.
type Parser interface {
	Marshal(interface{}) error
	Unmarshal([]byte, interface{}) error
}
