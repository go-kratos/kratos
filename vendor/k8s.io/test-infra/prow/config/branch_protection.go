/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"errors"
	"fmt"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
)

// Policy for the config/org/repo/branch.
// When merging policies, a nil value results in inheriting the parent policy.
type Policy struct {
	deprecatedPolicy
	deprecatedWarning bool // true if a warning message was sent
	// Protect overrides whether branch protection is enabled if set.
	Protect *bool `json:"protect,omitempty"`
	// RequiredStatusChecks configures github contexts
	RequiredStatusChecks *ContextPolicy `json:"required_status_checks,omitempty"`
	// Admins overrides whether protections apply to admins if set.
	Admins *bool `json:"enforce_admins,omitempty"`
	// Restrictions limits who can merge
	Restrictions *Restrictions `json:"restrictions,omitempty"`
	// RequiredPullRequestReviews specifies github approval/review criteria.
	RequiredPullRequestReviews *ReviewPolicy `json:"required_pull_request_reviews,omitempty"`
}

// deprecatedPolicy deserializes fields that are no longer in use
type deprecatedPolicy struct {
	DeprecatedProtect  *bool    `json:"protect-by-default,omitempty"`
	DeprecatedContexts []string `json:"require-contexts,omitempty"`
	DeprecatedPushers  []string `json:"allow-push,omitempty"`
}

func (d deprecatedPolicy) defined() bool {
	return d.DeprecatedProtect != nil || d.DeprecatedContexts != nil || d.DeprecatedPushers != nil
}

func (p Policy) defined() bool {
	return p.Protect != nil || p.RequiredStatusChecks != nil || p.Admins != nil || p.Restrictions != nil || p.RequiredPullRequestReviews != nil
}

// HasProtect returns true if the policy or deprecated policy defines protection
func (p Policy) HasProtect() bool {
	return p.Protect != nil || p.deprecatedPolicy.DeprecatedProtect != nil
}

// ContextPolicy configures required github contexts.
// When merging policies, contexts are appended to context list from parent.
// Strict determines whether merging to the branch invalidates existing contexts.
type ContextPolicy struct {
	// Contexts appends required contexts that must be green to merge
	Contexts []string `json:"contexts,omitempty"`
	// Strict overrides whether new commits in the base branch require updating the PR if set
	Strict *bool `json:"strict,omitempty"`
}

// ReviewPolicy specifies github approval/review criteria.
// Any nil values inherit the policy from the parent, otherwise bool/ints are overridden.
// Non-empty lists are appended to parent lists.
type ReviewPolicy struct {
	// Restrictions appends users/teams that are allowed to merge
	DismissalRestrictions *Restrictions `json:"dismissal_restrictions,omitempty"`
	// DismissStale overrides whether new commits automatically dismiss old reviews if set
	DismissStale *bool `json:"dismiss_stale_reviews,omitempty"`
	// RequireOwners overrides whether CODEOWNERS must approve PRs if set
	RequireOwners *bool `json:"require_code_owner_reviews,omitempty"`
	// Approvals overrides the number of approvals required if set (set to 0 to disable)
	Approvals *int `json:"required_approving_review_count,omitempty"`
}

// Restrictions limits who can merge
// Users and Teams items are appended to parent lists.
type Restrictions struct {
	Users []string `json:"users"`
	Teams []string `json:"teams"`
}

// selectInt returns the child if set, else parent
func selectInt(parent, child *int) *int {
	if child != nil {
		return child
	}
	return parent
}

// selectBool returns the child argument if set, otherwise the parent
func selectBool(parent, child *bool) *bool {
	if child != nil {
		return child
	}
	return parent
}

// unionStrings merges the parent and child items together
func unionStrings(parent, child []string) []string {
	if child == nil {
		return parent
	}
	if parent == nil {
		return child
	}
	s := sets.NewString(parent...)
	s.Insert(child...)
	return s.List()
}

func mergeContextPolicy(parent, child *ContextPolicy) *ContextPolicy {
	if child == nil {
		return parent
	}
	if parent == nil {
		return child
	}
	return &ContextPolicy{
		Contexts: unionStrings(parent.Contexts, child.Contexts),
		Strict:   selectBool(parent.Strict, child.Strict),
	}
}

func mergeReviewPolicy(parent, child *ReviewPolicy) *ReviewPolicy {
	if child == nil {
		return parent
	}
	if parent == nil {
		return child
	}
	return &ReviewPolicy{
		DismissalRestrictions: mergeRestrictions(parent.DismissalRestrictions, child.DismissalRestrictions),
		DismissStale:          selectBool(parent.DismissStale, child.DismissStale),
		RequireOwners:         selectBool(parent.RequireOwners, child.RequireOwners),
		Approvals:             selectInt(parent.Approvals, child.Approvals),
	}
}

func mergeRestrictions(parent, child *Restrictions) *Restrictions {
	if child == nil {
		return parent
	}
	if parent == nil {
		return child
	}
	return &Restrictions{
		Users: unionStrings(parent.Users, child.Users),
		Teams: unionStrings(parent.Teams, child.Teams),
	}
}

// Apply returns a policy that merges the child into the parent
func (p Policy) Apply(child Policy) (Policy, error) {
	if old := child.deprecatedPolicy.defined(); old && child.defined() {
		return p, errors.New("cannot mix Policy and deprecatedPolicy branch protection fields")
	} else if old {
		if !p.deprecatedWarning {
			p.deprecatedWarning = true
			logrus.Warn("WARNING: protect-by-default, require-contexts, allow-push are deprecated. Please replace them before July 2018")
		}
		d := child.deprecatedPolicy
		child = Policy{
			Protect: d.DeprecatedProtect,
		}
		if d.DeprecatedContexts != nil {
			child.RequiredStatusChecks = &ContextPolicy{
				Contexts: d.DeprecatedContexts,
			}
		}
		if d.DeprecatedPushers != nil {
			child.Restrictions = &Restrictions{
				Teams: d.DeprecatedPushers,
			}
		}
	}

	return Policy{
		Protect:                    selectBool(p.Protect, child.Protect),
		RequiredStatusChecks:       mergeContextPolicy(p.RequiredStatusChecks, child.RequiredStatusChecks),
		Admins:                     selectBool(p.Admins, child.Admins),
		Restrictions:               mergeRestrictions(p.Restrictions, child.Restrictions),
		RequiredPullRequestReviews: mergeReviewPolicy(p.RequiredPullRequestReviews, child.RequiredPullRequestReviews),
		deprecatedWarning:          p.deprecatedWarning,
	}, nil
}

// BranchProtection specifies the global branch protection policy
type BranchProtection struct {
	Policy
	ProtectTested         bool           `json:"protect-tested-repos,omitempty"`
	Orgs                  map[string]Org `json:"orgs,omitempty"`
	AllowDisabledPolicies bool           `json:"allow_disabled_policies,omitempty"`

	warned bool // warn if deprecated fields are use
}

// Org holds the default protection policy for an entire org, as well as any repo overrides.
type Org struct {
	Policy
	Repos map[string]Repo `json:"repos,omitempty"`
}

// Repo holds protection policy overrides for all branches in a repo, as well as specific branch overrides.
type Repo struct {
	Policy
	Branches map[string]Branch `json:"branches,omitempty"`
}

// Branch holds protection policy overrides for a particular branch.
type Branch struct {
	Policy
}

// GetBranchProtection returns the policy for a given branch.
//
// Handles merging any policies defined at repo/org/global levels into the branch policy.
func (c *Config) GetBranchProtection(org, repo, branch string) (*Policy, error) {
	bp := c.BranchProtection
	var policy Policy
	policy, err := policy.Apply(bp.Policy)
	if err != nil {
		return nil, err
	}

	if o, ok := bp.Orgs[org]; ok {
		policy, err = policy.Apply(o.Policy)
		if err != nil {
			return nil, err
		}
		if r, ok := o.Repos[repo]; ok {
			policy, err = policy.Apply(r.Policy)
			if err != nil {
				return nil, err
			}
			if b, ok := r.Branches[branch]; ok {
				policy, err = policy.Apply(b.Policy)
				if err != nil {
					return nil, err
				}
				if policy.Protect == nil {
					return nil, errors.New("defined branch policies must set protect")
				}
			}
		}
	} else {
		return nil, nil
	}

	// Automatically require any required prow jobs
	if prowContexts, _ := BranchRequirements(org, repo, branch, c.Presubmits); len(prowContexts) > 0 {
		// Error if protection is disabled
		if policy.Protect != nil && !*policy.Protect {
			return nil, fmt.Errorf("required prow jobs require branch protection")
		}
		ps := Policy{
			RequiredStatusChecks: &ContextPolicy{
				Contexts: prowContexts,
			},
		}
		// Require protection by default if ProtectTested is true
		if bp.ProtectTested {
			yes := true
			ps.Protect = &yes
		}
		policy, err = policy.Apply(ps)
		if err != nil {
			return nil, err
		}
	}

	if policy.Protect != nil && !*policy.Protect {
		// Ensure that protection is false => no protection settings
		var old *bool
		old, policy.Protect = policy.Protect, old
		switch {
		case policy.defined() && bp.AllowDisabledPolicies:
			logrus.Warnf("%s/%s=%s defines a policy but has protect: false", org, repo, branch)
			policy = Policy{
				Protect: policy.Protect,
			}
		case policy.defined():
			return nil, fmt.Errorf("%s/%s=%s defines a policy, which requires protect: true", org, repo, branch)
		}
		policy.Protect = old
	}

	if !policy.defined() {
		return nil, nil
	}
	return &policy, nil
}

func jobRequirements(jobs []Presubmit, branch string, after bool) ([]string, []string) {
	var required, optional []string
	for _, j := range jobs {
		if !j.Brancher.RunsAgainstBranch(branch) {
			continue
		}
		// Does this job require a context or have kids that might need one?
		if !after && !j.AlwaysRun && j.RunIfChanged == "" {
			continue // No
		}
		if j.ContextRequired() { // This job needs a context
			required = append(required, j.Context)
		} else {
			optional = append(optional, j.Context)
		}
		// Check which children require contexts
		r, o := jobRequirements(j.RunAfterSuccess, branch, true)
		required = append(required, r...)
		optional = append(optional, o...)
	}
	return required, optional
}

// BranchRequirements returns required and optional presubmits prow jobs for a given org, repo branch.
func BranchRequirements(org, repo, branch string, presubmits map[string][]Presubmit) ([]string, []string) {
	p, ok := presubmits[org+"/"+repo]
	if !ok {
		return nil, nil
	}
	return jobRequirements(p, branch, false)
}
