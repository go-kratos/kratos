package survey

import (
	"errors"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

/*
Select is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a string.

	color := ""
	prompt := &survey.Select{
		Message: "Choose a color:",
		Options: []string{"red", "blue", "green"},
	}
	survey.AskOne(prompt, &color, nil)
*/
type Select struct {
	core.Renderer
	Message       string
	Options       []string
	Default       string
	Help          string
	PageSize      int
	VimMode       bool
	FilterMessage string
	filter        string
	selectedIndex int
	useDefault    bool
	showingHelp   bool
}

// the data available to the templates when processing
type SelectTemplateData struct {
	Select
	PageEntries   []string
	SelectedIndex int
	Answer        string
	ShowAnswer    bool
	ShowHelp      bool
}

var SelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else}}
  {{- "  "}}{{- color "cyan"}}[Use arrows to move, type to filter{{- if and .Help (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $choice := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan+b"}}{{ SelectFocusIcon }} {{else}}{{color "default+hb"}}  {{end}}
    {{- $choice}}
    {{- color "reset"}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (s *Select) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	options := s.filterOptions()
	oldFilter := s.filter

	// if the user pressed the enter key
	if key == terminal.KeyEnter {
		if s.selectedIndex < len(options) {
			return []rune(options[s.selectedIndex]), 0, true
		}
		// if the user pressed the up arrow or 'k' to emulate vim
	} else if key == terminal.KeyArrowUp || (s.VimMode && key == 'k') {
		s.useDefault = false

		// if we are at the top of the list
		if s.selectedIndex == 0 {
			// start from the button
			s.selectedIndex = len(options) - 1
		} else {
			// otherwise we are not at the top of the list so decrement the selected index
			s.selectedIndex--
		}
		// if the user pressed down or 'j' to emulate vim
	} else if key == terminal.KeyArrowDown || (s.VimMode && key == 'j') {
		s.useDefault = false
		// if we are at the bottom of the list
		if s.selectedIndex == len(options)-1 {
			// start from the top
			s.selectedIndex = 0
		} else {
			// increment the selected index
			s.selectedIndex++
		}
		// only show the help message if we have one
	} else if key == core.HelpInputRune && s.Help != "" {
		s.showingHelp = true
		// if the user wants to toggle vim mode on/off
	} else if key == terminal.KeyEscape {
		s.VimMode = !s.VimMode
		// if the user hits any of the keys that clear the filter
	} else if key == terminal.KeyDeleteWord || key == terminal.KeyDeleteLine {
		s.filter = ""
		// if the user is deleting a character in the filter
	} else if key == terminal.KeyDelete || key == terminal.KeyBackspace {
		// if there is content in the filter to delete
		if s.filter != "" {
			// subtract a line from the current filter
			s.filter = s.filter[0 : len(s.filter)-1]
			// we removed the last value in the filter
		}
	} else if key >= terminal.KeySpace {
		s.filter += string(key)
		// make sure vim mode is disabled
		s.VimMode = false
		// make sure that we use the current value in the filtered list
		s.useDefault = false
	}

	s.FilterMessage = ""
	if s.filter != "" {
		s.FilterMessage = " " + s.filter
	}
	if oldFilter != s.filter {
		// filter changed
		options = s.filterOptions()
		if len(options) > 0 && len(options) <= s.selectedIndex {
			s.selectedIndex = len(options) - 1
		}
	}

	// figure out the options and index to render

	// TODO if we have started filtering and were looking at the end of a list
	// and we have modified the filter then we should move the page back!
	opts, idx := paginate(s.PageSize, options, s.selectedIndex)

	// render the options
	s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			SelectedIndex: idx,
			ShowHelp:      s.showingHelp,
			PageEntries:   opts,
		},
	)

	// if we are not pressing ent
	if len(options) <= s.selectedIndex {
		return []rune{}, 0, false
	}
	return []rune(options[s.selectedIndex]), 0, true
}

func (s *Select) filterOptions() []string {
	filter := strings.ToLower(s.filter)
	if filter == "" {
		return s.Options
	}
	answer := []string{}
	for _, o := range s.Options {
		if strings.Contains(strings.ToLower(o), filter) {
			answer = append(answer, o)
		}
	}
	return answer
}

func (s *Select) Prompt() (interface{}, error) {
	// if there are no options to render
	if len(s.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// start off with the first option selected
	sel := 0
	// if there is a default
	if s.Default != "" {
		// find the choice
		for i, opt := range s.Options {
			// if the option correponds to the default
			if opt == s.Default {
				// we found our initial value
				sel = i
				// stop looking
				break
			}
		}
	}
	// save the selected index
	s.selectedIndex = sel

	// figure out the options and index to render
	opts, idx := paginate(s.PageSize, s.Options, sel)

	// ask the question
	err := s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:        *s,
			PageEntries:   opts,
			SelectedIndex: idx,
		},
	)
	if err != nil {
		return "", err
	}

	// by default, use the default value
	s.useDefault = true

	rr := s.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	cursor := s.NewCursor()
	cursor.Hide()       // hide the cursor
	defer cursor.Show() // show the cursor when we're done

	// start waiting for input
	for {
		r, _, err := rr.ReadRune()
		if err != nil {
			return "", err
		}
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		s.OnChange(nil, 0, r)
	}
	options := s.filterOptions()
	s.filter = ""
	s.FilterMessage = ""

	var val string
	// if we are supposed to use the default value
	if s.useDefault || s.selectedIndex >= len(options) {
		// if there is a default value
		if s.Default != "" {
			// use the default value
			val = s.Default
		} else if len(options) > 0 {
			// there is no default value so use the first
			val = options[0]
		}
		// otherwise the selected index points to the value
	} else if s.selectedIndex < len(options) {
		// the
		val = options[s.selectedIndex]
	}
	return val, err
}

func (s *Select) Cleanup(val interface{}) error {
	return s.Render(
		SelectQuestionTemplate,
		SelectTemplateData{
			Select:     *s,
			Answer:     val.(string),
			ShowAnswer: true,
		},
	)
}
