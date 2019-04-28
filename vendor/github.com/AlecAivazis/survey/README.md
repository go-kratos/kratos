# Survey

[![Build Status](https://travis-ci.org/AlecAivazis/survey.svg?branch=feature%2Fpretty)](https://travis-ci.org/AlecAivazis/survey)
[![GoDoc](http://img.shields.io/badge/godoc-reference-5272B4.svg)](https://godoc.org/gopkg.in/AlecAivazis/survey.v1)

A library for building interactive prompts. 

<img width="550" src="https://thumbs.gfycat.com/VillainousGraciousKouprey-size_restricted.gif"/>

```go
package main

import (
    "fmt"
    "gopkg.in/AlecAivazis/survey.v1"
)

// the questions to ask
var qs = []*survey.Question{
    {
        Name:     "name",
        Prompt:   &survey.Input{Message: "What is your name?"},
        Validate: survey.Required,
        Transform: survey.Title,
    },
    {
        Name: "color",
        Prompt: &survey.Select{
            Message: "Choose a color:",
            Options: []string{"red", "blue", "green"},
            Default: "red",
        },
    },
    {
        Name: "age",
        Prompt:   &survey.Input{Message: "How old are you?"},
    },
}

func main() {
    // the answers will be written to this struct
    answers := struct {
        Name          string                  // survey will match the question and field names
        FavoriteColor string `survey:"color"` // or you can tag fields to match a specific name
        Age           int                     // if the types don't match exactly, survey will try to convert for you
    }{}

    // perform the questions
    err := survey.Ask(qs, &answers)
    if err != nil {
        fmt.Println(err.Error())
        return
    }

    fmt.Printf("%s chose %s.", answers.Name, answers.FavoriteColor)
}
```

## Table of Contents

1. [Examples](#examples)
1. [Prompts](#prompts)
   1. [Input](#input)
   1. [Multiline](#multiline)
   1. [Password](#password)
   1. [Confirm](#confirm)
   1. [Select](#select)
   1. [MultiSelect](#multiselect)
   1. [Editor](#editor)
1. [Validation](#validation)
   1. [Built-in Validators](#built-in-validators)
1. [Help Text](#help-text)
   1. [Changing the input rune](#changing-the-input-run)
1. [Custom Types](#custom-types)
1. [Customizing Output](#customizing-output)
1. [Versioning](#versioning)
1. [Testing](#testing)

## Examples

Examples can be found in the `examples/` directory. Run them
to see basic behavior:

```bash
go get gopkg.in/AlecAivazis/survey.v1

cd $GOPATH/src/gopkg.in/AlecAivazis/survey.v1

go run examples/simple.go
go run examples/validation.go
```

## Prompts

### Input

<img src="https://thumbs.gfycat.com/LankyBlindAmericanpainthorse-size_restricted.gif" width="400px"/>

```golang
name := ""
prompt := &survey.Input{
    Message: "ping",
}
survey.AskOne(prompt, &name, nil)
```

### Multiline

<img src="https://thumbs.gfycat.com/ImperfectShimmeringBeagle-size_restricted.gif" width="400px"/>

```golang
text := ""
prompt := &survey.Multiline{
    Message: "ping",
}
survey.AskOne(prompt, &text, nil)
```

### Password

<img src="https://thumbs.gfycat.com/CompassionateSevereHypacrosaurus-size_restricted.gif" width="400px" />

```golang
password := ""
prompt := &survey.Password{
    Message: "Please type your password",
}
survey.AskOne(prompt, &password, nil)
```

### Confirm

<img src="https://thumbs.gfycat.com/UnkemptCarefulGermanpinscher-size_restricted.gif" width="400px"/>

```golang
name := false
prompt := &survey.Confirm{
    Message: "Do you like pie?",
}
survey.AskOne(prompt, &name, nil)
```

### Select

<img src="https://thumbs.gfycat.com/GrimFilthyAmazonparrot-size_restricted.gif" width="450px"/>

```golang
color := ""
prompt := &survey.Select{
    Message: "Choose a color:",
    Options: []string{"red", "blue", "green"},
}
survey.AskOne(prompt, &color, nil)
```

The user can filter for options by typing while the prompt is active. The user can also press `esc` to toggle 
the ability cycle through the options with the j and k keys to do down and up respectively.

By default, the select prompt is limited to showing 7 options at a time
and will paginate lists of options longer than that. To increase, you can either
change the global `survey.PageSize`, or set the `PageSize` field on the prompt:

```golang
prompt := &survey.Select{..., PageSize: 10}
```

### MultiSelect

<img src="https://thumbs.gfycat.com/SharpTameAntelope-size_restricted.gif" width="450px"/>

```golang
days := []string{}
prompt := &survey.MultiSelect{
    Message: "What days do you prefer:",
    Options: []string{"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"},
}
survey.AskOne(prompt, &days, nil)
```

The user can filter for options by typing while the prompt is active. The user can also press `esc` to toggle 
the ability cycle through the options with the j and k keys to do down and up respectively.

By default, the MultiSelect prompt is limited to showing 7 options at a time
and will paginate lists of options longer than that. To increase, you can either
change the global `survey.PageSize`, or set the `PageSize` field on the prompt:

```golang
prompt := &survey.MultiSelect{..., PageSize: 10}
```

### Editor

Launches the user's preferred editor (defined by the $EDITOR environment variable) on a
temporary file. Once the user exits their editor, the contents of the temporary file are read in as
the result. If neither of those are present, notepad (on Windows) or vim (Linux or Mac) is used.

## Validation

Validating individual responses for a particular question can be done by defining a
`Validate` field on the `survey.Question` to be validated. This function takes an
`interface{}` type and returns an error to show to the user, prompting them for another
response:

```golang
q := &survey.Question{
    Prompt: &survey.Input{Message: "Hello world validation"},
    Validate: func (val interface{}) error {
        // since we are validating an Input, the assertion will always succeed
        if str, ok := val.(string) ; !ok || len(str) > 10 {
            return errors.New("This response cannot be longer than 10 characters.")
        }
    }
}
```

### Built-in Validators

`survey` comes prepackaged with a few validators to fit common situations. Currently these
validators include:

| name         | valid types | description                                                 | notes                                                                                 |
| ------------ | ----------- | ----------------------------------------------------------- | ------------------------------------------------------------------------------------- |
| Required     | any         | Rejects zero values of the response type                    | Boolean values pass straight through since the zero value (false) is a valid response |
| MinLength(n) | string      | Enforces that a response is at least the given length       |                                                                                       |
| MaxLength(n) | string      | Enforces that a response is no longer than the given length |                                                                                       |

## Help Text

All of the prompts have a `Help` field which can be defined to provide more information to your users:

<img src="https://thumbs.gfycat.com/CloudyRemorsefulFossa-size_restricted.gif" width="400px" style="margin-top: 8px"/>

```golang
&survey.Input{
    Message: "What is your phone number:",
    Help:    "Phone number should include the area code",
}
```

### Changing the input rune

In some situations, `?` is a perfectly valid response. To handle this, you can change the rune that survey
looks for by setting the `HelpInputRune` variable in `survey/core`:

```golang
import (
    "gopkg.in/AlecAivazis/survey.v1"
    surveyCore "gopkg.in/AlecAivazis/survey.v1/core"
)

number := ""
prompt := &survey.Input{
    Message: "If you have this need, please give me a reasonable message.",
    Help:    "I couldn't come up with one.",
}

surveyCore.HelpIcon = '^'

survey.AskOne(prompt, &number, nil)
```

## Custom Types

survey will assign prompt answers to your custom types if they implement this interface:

```golang
type settable interface {
    WriteAnswer(field string, value interface{}) error
}
```

Here is an example how to use them:

```golang
type MyValue struct {
    value string
}
func (my *MyValue) WriteAnswer(name string, value interface{}) error {
     my.value = value.(string)
}

myval := MyValue{}
survey.AskOne(
    &survey.Input{
        Message: "Enter something:",
    },
    &myval,
    nil,
)
```

## Customizing Output

Customizing the icons and various parts of survey can easily be done by setting the following variables
in `survey/core`:

| name               | default | description                                                   |
| ------------------ | ------- | ------------------------------------------------------------- |
| ErrorIcon          | X       | Before an error                                               |
| HelpIcon           | i       | Before help text                                              |
| QuestionIcon       | ?       | Before the message of a prompt                                |
| SelectFocusIcon    | >       | Marks the current focus in `Select` and `MultiSelect` prompts |
| UnmarkedOptionIcon | [ ]     | Marks an unselected option in a `MultiSelect` prompt          |
| MarkedOptionIcon   | [x]     | Marks a chosen selection in a `MultiSelect` prompt            |

## Versioning

This project tries to maintain semantic GitHub releases as closely as possible and relies on [gopkg.in](http://labix.org/gopkg.in)
to maintain those releases. Importing version 1 of survey would look like:

```golang
package main

import "gopkg.in/AlecAivazis/survey.v1"
```

## Testing

You can test your program's interactive prompts using [go-expect](https://github.com/Netflix/go-expect). The library
can be used to expect a match on stdout and respond on stdin. Since `os.Stdout` in a `go test` process is not a TTY,
if you are manipulating the cursor or using `survey`, you will need a way to interpret terminal / ANSI escape sequences
for things like `CursorLocation`. `vt10x.NewVT10XConsole` will create a `go-expect` console that also multiplexes
stdio to an in-memory [virtual terminal](https://github.com/hinshun/vt10x).

For example, you can test a binary utilizing `survey` by connecting the Console's tty to a subprocess's stdio. 

```go
func TestCLI(t *testing.T) {
 	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
 	// sequences (i.e. cursor position report).
 	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
	defer c.Close()

	donec := make(chan struct{})
	go func() {
		defer close(donec)
    		c.ExpectString("What is your name?")
    		c.SendLine("Johnny Appleseed")
    		c.ExpectEOF()
  	}()

	cmd := exec.Command("your-cli")
  	cmd.Stdin = c.Tty()
  	cmd.Stdout = c.Tty()
  	cmd.Stderr = c.Tty()

  	err = cmd.Run()
  	require.Nil(t, err)

  	// Close the slave end of the pty, and read the remaining bytes from the master end.
  	c.Tty().Close()
  	<-donec

  	// Dump the terminal's screen.
  	t.Log(expect.StripTrailingEmptyLines(state.String()))
}
```

If your application is decoupled from `os.Stdout` and `os.Stdin`, you can even test through the tty alone.
`survey` itself is tested in this manner.

```go
func TestCLI(t *testing.T) {
  	// Multiplex stdin/stdout to a virtual terminal to respond to ANSI escape
  	// sequences (i.e. cursor position report).
	c, state, err := vt10x.NewVT10XConsole()
	require.Nil(t, err)
  	defer c.Close()

  	donec := make(chan struct{})
	go func() {
    		defer close(donec)
    		c.ExpectString("What is your name?")
    		c.SendLine("Johnny Appleseed")
    		c.ExpectEOF()
	}()

  	prompt := &Input{
    		Message: "What is your name?",
  	}
  	prompt.WithStdio(Stdio(c))

  	answer, err := prompt.Prompt()
  	require.Nil(t, err)
  	require.Equal(t, "Johnny Appleseed", answer)

  	// Close the slave end of the pty, and read the remaining bytes from the master end.
  	c.Tty().Close()
  	<-donec

  	// Dump the terminal's screen.
  	t.Log(expect.StripTrailingEmptyLines(state.String()))
}
```
