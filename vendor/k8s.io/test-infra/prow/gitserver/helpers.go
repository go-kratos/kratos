package gitserver

import (
	"strings"
)

// HasLabel checks if label is in the label set "issueLabels".
func HasLabel(label string, issueLabels []Label) bool {
	for _, l := range issueLabels {
		if strings.ToLower(l.Name) == strings.ToLower(label) {
			return true
		}
	}
	return false
}

// ChangedLabels describe a gitlab PR changed labels
func ChangedLabels(action PullRequestEventAction, previous, current []Label) []Label {
	labels := make([]Label, 0)
	if action == PullRequestActionLabeled {
		for _, l := range current {
			if !HasLabel(l.Name, previous) {
				labels = append(labels, l)
			}
		}
	} else if action == PullRequestActionUnlabeled {
		for _, l := range previous {
			if !HasLabel(l.Name, current) {
				labels = append(labels, l)
			}
		}
	}
	return labels
}
