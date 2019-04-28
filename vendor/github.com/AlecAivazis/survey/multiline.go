package survey

import (
	"fmt"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

type Multiline struct {
	core.Renderer
	Message string
	Default string
	Help    string
}

// data available to the templates when processing
type MultilineTemplateData struct {
	Multiline
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var MultilineQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- "\n"}}{{color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
  {{- color "cyan"}}[Enter 2 empty lines to finish]{{color "reset"}}
{{- end}}`

func (i *Multiline) Prompt() (interface{}, error) {
	// render the template
	err := i.Render(
		MultilineQuestionTemplate,
		MultilineTemplateData{Multiline: *i},
	)
	if err != nil {
		return "", err
	}
	fmt.Println()

	// start reading runes from the standard in
	rr := i.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	cursor := i.NewCursor()

	multiline := make([]string, 0)

	emptyOnce := false
	// get the next line
	for {
		line := []rune{}
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}

		if string(line) == "" {
			if emptyOnce {
				numLines := len(multiline) + 2
				cursor.PreviousLine(numLines)
				for j := 0; j < numLines; j++ {
					terminal.EraseLine(i.Stdio().Out, terminal.ERASE_LINE_ALL)
					cursor.NextLine(1)
				}
				cursor.PreviousLine(numLines)
				break
			}
			emptyOnce = true
		} else {
			emptyOnce = false
		}
		multiline = append(multiline, string(line))
	}

	val := strings.Join(multiline, "\n")
	val = strings.TrimSpace(val)

	// if the line is empty
	if len(val) == 0 {
		// use the default value
		return i.Default, err
	}

	// we're done
	return val, err
}

func (i *Multiline) Cleanup(val interface{}) error {
	return i.Render(
		MultilineQuestionTemplate,
		MultilineTemplateData{Multiline: *i, Answer: val.(string), ShowAnswer: true},
	)
}
