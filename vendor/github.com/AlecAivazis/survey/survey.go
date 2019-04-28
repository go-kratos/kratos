package survey

import (
	"errors"
	"io"
	"os"

	"gopkg.in/AlecAivazis/survey.v1/core"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

// PageSize is the default maximum number of items to show in select/multiselect prompts
var PageSize = 7

// DefaultAskOptions is the default options on ask, using the OS stdio.
var DefaultAskOptions = AskOptions{
	Stdio: terminal.Stdio{
		In:  os.Stdin,
		Out: os.Stdout,
		Err: os.Stderr,
	},
}

// Validator is a function passed to a Question after a user has provided a response.
// If the function returns an error, then the user will be prompted again for another
// response.
type Validator func(ans interface{}) error

// Transformer is a function passed to a Question after a user has provided a response.
// The function can be used to implement a custom logic that will result to return
// a different representation of the given answer.
//
// Look `TransformString`, `ToLower` `Title` and `ComposeTransformers` for more.
type Transformer func(ans interface{}) (newAns interface{})

// Question is the core data structure for a survey questionnaire.
type Question struct {
	Name      string
	Prompt    Prompt
	Validate  Validator
	Transform Transformer
}

// Prompt is the primary interface for the objects that can take user input
// and return a response.
type Prompt interface {
	Prompt() (interface{}, error)
	Cleanup(interface{}) error
	Error(error) error
}

// AskOpt allows setting optional ask options.
type AskOpt func(options *AskOptions) error

// AskOptions provides additional options on ask.
type AskOptions struct {
	Stdio terminal.Stdio
}

// WithStdio specifies the standard input, output and error files survey
// interacts with. By default, these are os.Stdin, os.Stdout, and os.Stderr.
func WithStdio(in terminal.FileReader, out terminal.FileWriter, err io.Writer) AskOpt {
	return func(options *AskOptions) error {
		options.Stdio.In = in
		options.Stdio.Out = out
		options.Stdio.Err = err
		return nil
	}
}

type wantsStdio interface {
	WithStdio(terminal.Stdio)
}

/*
AskOne performs the prompt for a single prompt and asks for validation if required.
Response types should be something that can be casted from the response type designated
in the documentation. For example:

	name := ""
	prompt := &survey.Input{
		Message: "name",
	}

	survey.AskOne(prompt, &name, nil)

*/
func AskOne(p Prompt, response interface{}, v Validator, opts ...AskOpt) error {
	err := Ask([]*Question{{Prompt: p, Validate: v}}, response, opts...)
	if err != nil {
		return err
	}

	return nil
}

/*
Ask performs the prompt loop, asking for validation when appropriate. The response
type can be one of two options. If a struct is passed, the answer will be written to
the field whose name matches the Name field on the corresponding question. Field types
should be something that can be casted from the response type designated in the
documentation. Note, a survey tag can also be used to identify a Otherwise, a
map[string]interface{} can be passed, responses will be written to the key with the
matching name. For example:

	qs := []*survey.Question{
		{
			Name:     "name",
			Prompt:   &survey.Input{Message: "What is your name?"},
			Validate: survey.Required,
			Transform: survey.Title,
		},
	}

	answers := struct{ Name string }{}


	err := survey.Ask(qs, &answers)
*/
func Ask(qs []*Question, response interface{}, opts ...AskOpt) error {

	options := DefaultAskOptions
	for _, opt := range opts {
		if err := opt(&options); err != nil {
			return err
		}
	}

	// if we weren't passed a place to record the answers
	if response == nil {
		// we can't go any further
		return errors.New("cannot call Ask() with a nil reference to record the answers")
	}

	// go over every question
	for _, q := range qs {
		// If Prompt implements controllable stdio, pass in specified stdio.
		if p, ok := q.Prompt.(wantsStdio); ok {
			p.WithStdio(options.Stdio)
		}

		// grab the user input and save it
		ans, err := q.Prompt.Prompt()
		// if there was a problem
		if err != nil {
			return err
		}

		// if there is a validate handler for this question
		if q.Validate != nil {
			// wait for a valid response
			for invalid := q.Validate(ans); invalid != nil; invalid = q.Validate(ans) {
				err := q.Prompt.Error(invalid)
				// if there was a problem
				if err != nil {
					return err
				}

				// ask for more input
				ans, err = q.Prompt.Prompt()
				// if there was a problem
				if err != nil {
					return err
				}
			}
		}

		if q.Transform != nil {
			// check if we have a transformer available, if so
			// then try to acquire the new representation of the
			// answer, if the resulting answer is not nil.
			if newAns := q.Transform(ans); newAns != nil {
				ans = newAns
			}
		}

		// tell the prompt to cleanup with the validated value
		q.Prompt.Cleanup(ans)

		// if something went wrong
		if err != nil {
			// stop listening
			return err
		}

		// add it to the map
		err = core.WriteAnswer(response, q.Name, ans)
		// if something went wrong
		if err != nil {
			return err
		}

	}

	// return the response
	return nil
}

// paginate returns a single page of choices given the page size, the total list of
// possible choices, and the current selected index in the total list.
func paginate(page int, choices []string, sel int) ([]string, int) {
	// the number of elements to show in a single page
	var pageSize int
	// if the select has a specific page size
	if page != 0 {
		// use the specified one
		pageSize = page
		// otherwise the select does not have a page size
	} else {
		// use the package default
		pageSize = PageSize
	}

	var start, end, cursor int

	if len(choices) < pageSize {
		// if we dont have enough options to fill a page
		start = 0
		end = len(choices)
		cursor = sel

	} else if sel < pageSize/2 {
		// if we are in the first half page
		start = 0
		end = pageSize
		cursor = sel

	} else if len(choices)-sel-1 < pageSize/2 {
		// if we are in the last half page
		start = len(choices) - pageSize
		end = len(choices)
		cursor = sel - start

	} else {
		// somewhere in the middle
		above := pageSize / 2
		below := pageSize - above

		cursor = pageSize / 2
		start = sel - above
		end = sel + below
	}

	// return the subset we care about and the index
	return choices[start:end], cursor
}
