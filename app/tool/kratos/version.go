package main

import (
	"bytes"
	"runtime"
	"text/template"
)

var (
	// GitCommit is  git commit
	GitCommit = "library-import"

	// Version is version
	Version = "library-import"

	// BuildTime is BuildTime
	BuildTime = "library-import"

	// Channel is Channel
	Channel = "library-import"
)

// VersionOptions include version
type VersionOptions struct {
	GitCommit string
	Version   string
	BuildTime string
	GoVersion string
	Os        string
	Arch      string
	Channel   string
}

var versionTemplate = ` Version:      {{.Version}}
 Go version:   {{.GoVersion}}
 Git commit:   {{.GitCommit}}
 Built:        {{.BuildTime}}
 OS/Arch:      {{.Os}}/{{.Arch}}
 Channel:      {{.Channel}}
 `

func getVersion() string {
	var doc bytes.Buffer
	vo := VersionOptions{
		GitCommit: GitCommit,
		Version:   Version,
		BuildTime: BuildTime,
		GoVersion: runtime.Version(),
		Os:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Channel:   Channel,
	}
	tmpl, _ := template.New("version").Parse(versionTemplate)
	tmpl.Execute(&doc, vo)
	return doc.String()
}
