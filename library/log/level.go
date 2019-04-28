package log

// Level of severity.
type Level int

// Verbose is a boolean type that implements Info, Infov (like Printf) etc.
type Verbose bool

// common log level.
const (
	_debugLevel Level = iota
	_infoLevel
	_warnLevel
	_errorLevel
	_fatalLevel
)

var levelNames = [...]string{
	_debugLevel: "DEBUG",
	_infoLevel:  "INFO",
	_warnLevel:  "WARN",
	_errorLevel: "ERROR",
	_fatalLevel: "FATAL",
}

// String implementation.
func (l Level) String() string {
	return levelNames[l]
}
