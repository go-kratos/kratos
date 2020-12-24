package parser

// Parser is config parser.
type Parser interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(v interface{}) error
}
