package survey

import (
	"errors"
	"strings"

	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

/*
MultiSelect is a prompt that presents a list of various options to the user
for them to select using the arrow keys and enter. Response type is a slice of strings.

	days := []string{}
	prompt := &survey.MultiSelect{
		Message: "What days do you prefer:",
		Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
	}
	survey.AskOne(prompt, &days, nil)
*/
type MultiSelect struct {
	core.Renderer
	Message       string
	Options       []string
	Default       []string
	Help          string
	PageSize      int
	VimMode       bool
	FilterMessage string
	filter        string
	selectedIndex int
	checked       map[string]bool
	showingHelp   bool
}

// data available to the templates when processing
type MultiSelectTemplateData struct {
	MultiSelect
	Answer        string
	ShowAnswer    bool
	Checked       map[string]bool
	SelectedIndex int
	ShowHelp      bool
	PageEntries   []string
}

var MultiSelectQuestionTemplate = `
{{- if .ShowHelp }}{{- color "cyan"}}{{ HelpIcon }} {{ .Help }}{{color "reset"}}{{"\n"}}{{end}}
{{- color "green+hb"}}{{ QuestionIcon }} {{color "reset"}}
{{- color "default+hb"}}{{ .Message }}{{ .FilterMessage }}{{color "reset"}}
{{- if .ShowAnswer}}{{color "cyan"}} {{.Answer}}{{color "reset"}}{{"\n"}}
{{- else }}
	{{- "  "}}{{- color "cyan"}}[Use arrows to move, type to filter{{- if and .Help (not .ShowHelp)}}, {{ HelpInputRune }} for more help{{end}}]{{color "reset"}}
  {{- "\n"}}
  {{- range $ix, $option := .PageEntries}}
    {{- if eq $ix $.SelectedIndex}}{{color "cyan"}}{{ SelectFocusIcon }}{{color "reset"}}{{else}} {{end}}
    {{- if index $.Checked $option}}{{color "green"}} {{ MarkedOptionIcon }} {{else}}{{color "default+hb"}} {{ UnmarkedOptionIcon }} {{end}}
    {{- color "reset"}}
    {{- " "}}{{$option}}{{"\n"}}
  {{- end}}
{{- end}}`

// OnChange is called on every keypress.
func (m *MultiSelect) OnChange(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
	options := m.filterOptions()
	oldFilter := m.filter

	if key == terminal.KeyArrowUp || (m.VimMode && key == 'k') {
		// if we are at the top of the list
		if m.selectedIndex == 0 {
			// go to the bottom
			m.selectedIndex = len(options) - 1
		} else {
			// decrement the selected index
			m.selectedIndex--
		}
	} else if key == terminal.KeyArrowDown || (m.VimMode && key == 'j') {
		// if we are at the bottom of the list
		if m.selectedIndex == len(options)-1 {
			// start at the top
			m.selectedIndex = 0
		} else {
			// increment the selected index
			m.selectedIndex++
		}
		// if the user pressed down and there is room to move
	} else if key == terminal.KeySpace {
		if m.selectedIndex < len(options) {
			if old, ok := m.checked[options[m.selectedIndex]]; !ok {
				// otherwise just invert the current value
				m.checked[options[m.selectedIndex]] = true
			} else {
				// otherwise just invert the current value
				m.checked[options[m.selectedIndex]] = !old
			}
			m.filter = ""
		}
		// only show the help message if we have one to show
	} else if key == core.HelpInputRune && m.Help != "" {
		m.showingHelp = true
	} else if key == terminal.KeyEscape {
		m.VimMode = !m.VimMode
	} else if key == terminal.KeyDeleteWord || key == terminal.KeyDeleteLine {
		m.filter = ""
	} else if key == terminal.KeyDelete || key == terminal.KeyBackspace {
		if m.filter != "" {
			m.filter = m.filter[0 : len(m.filter)-1]
		}
	} else if key >= terminal.KeySpace {
		m.filter += string(key)
		m.VimMode = false
	}

	m.FilterMessage = ""
	if m.filter != "" {
		m.FilterMessage = " " + m.filter
	}
	if oldFilter != m.filter {
		// filter changed
		options = m.filterOptions()
		if len(options) > 0 && len(options) <= m.selectedIndex {
			m.selectedIndex = len(options) - 1
		}
	}
	// paginate the options

	// TODO if we have started filtering and were looking at the end of a list
	// and we have modified the filter then we should move the page back!
	opts, idx := paginate(m.PageSize, options, m.selectedIndex)

	// render the options
	m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: idx,
			Checked:       m.checked,
			ShowHelp:      m.showingHelp,
			PageEntries:   opts,
		},
	)

	// if we are not pressing ent
	return line, 0, true
}

func (m *MultiSelect) filterOptions() []string {
	filter := strings.ToLower(m.filter)
	if filter == "" {
		return m.Options
	}
	answer := []string{}
	for _, o := range m.Options {
		if strings.Contains(strings.ToLower(o), filter) {
			answer = append(answer, o)
		}
	}
	return answer
}

func (m *MultiSelect) Prompt() (interface{}, error) {
	// compute the default state
	m.checked = make(map[string]bool)
	// if there is a default
	if len(m.Default) > 0 {
		for _, dflt := range m.Default {
			for _, opt := range m.Options {
				// if the option correponds to the default
				if opt == dflt {
					// we found our initial value
					m.checked[opt] = true
					// stop looking
					break
				}
			}
		}
	}

	// if there are no options to render
	if len(m.Options) == 0 {
		// we failed
		return "", errors.New("please provide options to select from")
	}

	// paginate the options
	opts, idx := paginate(m.PageSize, m.Options, m.selectedIndex)

	cursor := m.NewCursor()
	cursor.Hide()       // hide the cursor
	defer cursor.Show() // show the cursor when we're done

	// ask the question
	err := m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: idx,
			Checked:       m.checked,
			PageEntries:   opts,
		},
	)
	if err != nil {
		return "", err
	}

	rr := m.NewRuneReader()
	rr.SetTermMode()
	defer rr.RestoreTermMode()

	// start waiting for input
	for {
		r, _, _ := rr.ReadRune()
		if r == '\r' || r == '\n' {
			break
		}
		if r == terminal.KeyInterrupt {
			return "", terminal.InterruptErr
		}
		if r == terminal.KeyEndTransmission {
			break
		}
		m.OnChange(nil, 0, r)
	}
	m.filter = ""
	m.FilterMessage = ""

	answers := []string{}
	for _, option := range m.Options {
		if val, ok := m.checked[option]; ok && val {
			answers = append(answers, option)
		}
	}

	return answers, nil
}

// Cleanup removes the options section, and renders the ask like a normal question.
func (m *MultiSelect) Cleanup(val interface{}) error {
	// execute the output summary template with the answer
	return m.Render(
		MultiSelectQuestionTemplate,
		MultiSelectTemplateData{
			MultiSelect:   *m,
			SelectedIndex: m.selectedIndex,
			Checked:       m.checked,
			Answer:        strings.Join(val.([]string), ", "),
			ShowAnswer:    true,
		},
	)
}
