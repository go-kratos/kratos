package parser

// Parser is config parser.
type Parser interface {
	Format() string
	Marshal(interface{}) error
	Unmarshal([]byte, interface{}) error
}
