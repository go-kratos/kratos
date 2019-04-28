package generator

import (
	"strings"

	"github.com/golang/glog"

	"go-common/app/tool/gengo/namer"
	"go-common/app/tool/gengo/types"
)

// NewImportTracker is
func NewImportTracker(typesToAdd ...*types.Type) namer.ImportTracker {
	tracker := namer.NewDefaultImportTracker(types.Name{})
	tracker.IsInvalidType = func(*types.Type) bool { return false }
	tracker.LocalName = func(name types.Name) string { return golangTrackerLocalName(&tracker, name) }
	tracker.PrintImport = func(path, name string) string { return name + " \"" + path + "\"" }

	tracker.AddTypes(typesToAdd...)
	return &tracker

}

func golangTrackerLocalName(tracker namer.ImportTracker, t types.Name) string {
	path := t.Package

	// Using backslashes in package names causes gengo to produce Go code which
	// will not compile with the gc compiler. See the comment on GoSeperator.
	if strings.ContainsRune(path, '\\') {
		glog.Warningf("Warning: backslash used in import path '%v', this is unsupported.\n", path)
	}

	dirs := strings.Split(path, namer.GoSeperator)
	for n := len(dirs) - 1; n >= 0; n-- {
		// follow kube convention of not having anything between directory names
		name := strings.Join(dirs[n:], "")
		name = strings.Replace(name, "_", "", -1)
		// These characters commonly appear in import paths for go
		// packages, but aren't legal go names. So we'll sanitize.
		name = strings.Replace(name, ".", "", -1)
		name = strings.Replace(name, "-", "", -1)
		if _, found := tracker.PathOf(name); found {
			// This name collides with some other package
			continue
		}
		return name
	}
	panic("can't find import for " + path)
}
