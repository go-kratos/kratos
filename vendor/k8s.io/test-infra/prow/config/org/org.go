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

package org

import (
	"fmt"
)

// Metadata declares metadata about the GitHub org.
//
// See https://developer.github.com/v3/orgs/#edit-an-organization
type Metadata struct {
	BillingEmail                 *string              `json:"billing_email,omitempty"`
	Company                      *string              `json:"company,omitempty"`
	Email                        *string              `json:"email,omitempty"`
	Name                         *string              `json:"name,omitempty"`
	Description                  *string              `json:"description,omitempty"`
	Location                     *string              `json:"location,omitempty"`
	HasOrganizationProjects      *bool                `json:"has_organization_projects,omitempty"`
	HasRepositoryProjects        *bool                `json:"has_repository_projects,omitempty"`
	DefaultRepositoryPermission  *RepoPermissionLevel `json:"default_repository_permission,omitempty"`
	MembersCanCreateRepositories *bool                `json:"members_can_create_repositories,omitempty"`
}

// Config declares org metadata as well as its people and teams.
type Config struct {
	Metadata
	Teams   map[string]Team `json:"teams,omitempty"`
	Members []string        `json:"members,omitempty"`
	Admins  []string        `json:"admins,omitempty"`
}

// TeamMetadata declares metadata about the github team.
//
// See https://developer.github.com/v3/teams/#edit-team
type TeamMetadata struct {
	Description *string  `json:"description,omitempty"`
	Privacy     *Privacy `json:"privacy,omitempty"`
}

// Team declares metadata as well as its poeple.
type Team struct {
	TeamMetadata
	Members     []string        `json:"members,omitempty"`
	Maintainers []string        `json:"maintainers,omitempty"`
	Children    map[string]Team `json:"teams,omitempty"`

	Previously []string `json:"previously,omitempty"`
}

// RepoPermissionLevel is admin, write, read or none.
//
// See https://developer.github.com/v3/repos/collaborators/#review-a-users-permission-level
type RepoPermissionLevel string

const (
	// Read allows pull but not push
	Read RepoPermissionLevel = "read"
	// Write allows Read plus push
	Write RepoPermissionLevel = "write"
	// Admin allows Write plus change others' rights.
	Admin RepoPermissionLevel = "admin"
	// None disallows everything
	None RepoPermissionLevel = "none"
)

var repoPermissionLevels = map[RepoPermissionLevel]bool{
	Read:  true,
	Write: true,
	Admin: true,
	None:  true,
}

// MarshalText returns the byte representation of the permission
func (l RepoPermissionLevel) MarshalText() ([]byte, error) {
	return []byte(l), nil
}

// UnmarshalText validates the text is a valid string
func (l *RepoPermissionLevel) UnmarshalText(text []byte) error {
	v := RepoPermissionLevel(text)
	if _, ok := repoPermissionLevels[v]; !ok {
		return fmt.Errorf("bad repo permission: %s not in %v", v, repoPermissionLevels)
	}
	*l = v
	return nil
}

// Privacy is secret or closed.
//
// See https://developer.github.com/v3/teams/#edit-team
type Privacy string

const (
	// Closed means it is only visible to org members
	Closed Privacy = "closed"
	// Secret means it is only visible to team members.
	Secret Privacy = "secret"
)

var privacySettings = map[Privacy]bool{
	Closed: true,
	Secret: true,
}

// MarshalText returns bytes that equal secret or closed
func (p Privacy) MarshalText() ([]byte, error) {
	return []byte(p), nil
}

// UnmarshalText returns an error if text != secret or closed
func (p *Privacy) UnmarshalText(text []byte) error {
	v := Privacy(text)
	if _, ok := privacySettings[v]; !ok {
		return fmt.Errorf("bad privacy setting: %s", v)
	}
	*p = v
	return nil
}
