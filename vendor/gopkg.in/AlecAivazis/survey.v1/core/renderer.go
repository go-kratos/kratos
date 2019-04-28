package core

import (
	"fmt"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

type Renderer struct {
	stdio          terminal.Stdio
	lineCount      int
	errorLineCount int
}

var ErrorTemplate = `{{color "red"}}{{ ErrorIcon }} Sorry, your reply was invalid: {{.Error}}{{color "reset"}}
`

func (r *Renderer) WithStdio(stdio terminal.Stdio) {
	r.stdio = stdio
}

func (r *Renderer) Stdio() terminal.Stdio {
	return r.stdio
}

func (r *Renderer) NewRuneReader() *terminal.RuneReader {
	return terminal.NewRuneReader(r.stdio)
}

func (r *Renderer) NewCursor() *terminal.Cursor {
	return &terminal.Cursor{
		In:  r.stdio.In,
		Out: r.stdio.Out,
	}
}

func (r *Renderer) Error(invalid error) error {
	// since errors are printed on top we need to reset the prompt
	// as well as any previous error print
	r.resetPrompt(r.lineCount + r.errorLineCount)
	// we just cleared the prompt lines
	r.lineCount = 0
	out, err := RunTemplate(ErrorTemplate, invalid)
	if err != nil {
		return err
	}
	// keep track of how many lines are printed so we can clean up later
	r.errorLineCount = strings.Count(out, "\n")

	// send the message to the user
	fmt.Fprint(terminal.NewAnsiStdout(r.stdio.Out), out)
	return nil
}

func (r *Renderer) resetPrompt(lines int) {
	// clean out current line in case tmpl didnt end in newline
	cursor := r.NewCursor()
	cursor.HorizontalAbsolute(0)
	terminal.EraseLine(r.stdio.Out, terminal.ERASE_LINE_ALL)
	// clean up what we left behind last time
	for i := 0; i < lines; i++ {
		cursor.PreviousLine(1)
		terminal.EraseLine(r.stdio.Out, terminal.ERASE_LINE_ALL)
	}
}

func (r *Renderer) Render(tmpl string, data interface{}) error {
	r.resetPrompt(r.lineCount)
	// render the template summarizing the current state
	out, err := RunTemplate(tmpl, data)
	if err != nil {
		return err
	}

	// keep track of how many lines are printed so we can clean up later
	r.lineCount = strings.Count(out, "\n")

	// print the summary
	fmt.Fprint(terminal.NewAnsiStdout(r.stdio.Out), out)

	// nothing went wrong
	return nil
}
