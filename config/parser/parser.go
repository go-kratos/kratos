package parser

// Parser is config json parser.
type Parser interface {
	Format() string
	Marshal(interface{}) ([]byte, error)
	Unmarshal([]byte, interface{}) error
}
