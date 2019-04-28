package survey

import (
	"gopkg.in/AlecAivazis/survey.v1/core"
)

/*
Input is a regular text input that prints each character the user types on the screen
and accepts the input with the enter key. Response type is a string.

	name := ""
	prompt := &survey.Input{ Message: "What is your name?" }
	survey.AskOne(prompt, &name, nil)
*/
type Input struct {
	core.Renderer
	Message string
	Default string
	Help    string
}

// data available to the templates when processing
type InputTemplateData struct {
	Input
	Answer     string
	ShowAnswer bool
	ShowHelp   bool
}

// Templates with Color formatting. See Documentation: https://github.com/mgutz/ansi#style-format
var InputQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }} {{color "reset"}}
{{- if .ShowAnswer}}
  {{- color "cyan"}}{{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
  {{- if and .Help (not .ShowHelp)}}{{color "cyan"}}[{{ HelpInputRune }} for help]{{color "reset"}} {{end}}
  {{- if .Default}}{{color "white"}}({{.Default}}) {{color "reset"}}{{end}}
{{- end}}`

func (i *Input) Prompt() (interface{}, error) {
	// render the template
	err := i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i},
	)
	if err != nil {
		return "", err
	}

	// start reading runes from the standard in
	rr := i.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	cursor := i.NewCursor()

	line := []rune{}
	// get the next line
	for {
		line, err = rr.ReadLine(0)
		if err != nil {
			return string(line), err
		}
		// terminal will echo the \n so we need to jump back up one row
		cursor.PreviousLine(1)

		if string(line) == string(core.HelpInputRune) && i.Help != "" {
			err = i.Render(
				InputQuestionTemplate,
				InputTemplateData{Input: *i, ShowHelp: true},
			)
			if err != nil {
				return "", err
			}
			continue
		}
		break
	}

	// if the line is empty
	if line == nil || len(line) == 0 {
		// use the default value
		return i.Default, err
	}

	// we're done
	return string(line), err
}

func (i *Input) Cleanup(val interface{}) error {
	return i.Render(
		InputQuestionTemplate,
		InputTemplateData{Input: *i, Answer: val.(string), ShowAnswer: true},
	)
}
