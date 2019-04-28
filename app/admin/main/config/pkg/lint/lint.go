package lint

import (
	"errors"
	"fmt"
	"io"
	"strings"
)

var linterMap map[string]Linter

// ErrLintNotExists .
var ErrLintNotExists = errors.New("linter not exists")

func init() {
	linterMap = make(map[string]Linter)
}

// RegisterLinter register linter for a kind of file
func RegisterLinter(filetype string, linter Linter) {
	if _, ok := linterMap[filetype]; ok {
		panic("linter for filetype " + filetype + " already exists")
	}
	linterMap[filetype] = linter
}

// LineErr error contains line number
type LineErr struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// Error lint error
type Error []LineErr

func (errs Error) Error() string {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		messages = append(messages, fmt.Sprintf("%d:%s", err.Line, err.Message))
	}
	return strings.Join(messages, "\n")
}

func (errs Error) String() string {
	return errs.Error()
}

// Linter lint config file
type Linter interface {
	Lint(r io.Reader) Error
}

// Lint config file,
func Lint(filetype string, r io.Reader) error {
	lint, ok := linterMap[filetype]
	if !ok {
		return ErrLintNotExists
	}
	if lintErr := lint.Lint(r); lintErr != nil {
		return lintErr
	}
	return nil
}
