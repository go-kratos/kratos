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
	"encoding/gob"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

// Cookie holds the secret returned from github that authenticates the user who authorized this app.
type Cookie struct {
	Secret string `json:"secret,omitempty"`
}

// GithubOAuthConfig is a config for requesting users access tokens from Github API. It also has
// a Cookie Store that retains user credentials deriving from Github API.
type GithubOAuthConfig struct {
	ClientID         string   `json:"client_id"`
	ClientSecret     string   `json:"client_secret"`
	RedirectURL      string   `json:"redirect_url"`
	Scopes           []string `json:"scopes,omitempty"`
	FinalRedirectURL string   `json:"final_redirect_url"`

	CookieStore *sessions.CookieStore `json:"-"`
}

// InitGithubOAuthConfig creates an OAuthClient using GithubOAuth config and a Cookie Store
// to retain user credentials.
func (gac *GithubOAuthConfig) InitGithubOAuthConfig(cookie *sessions.CookieStore) {
	gob.Register(&oauth2.Token{})
	gac.CookieStore = cookie
}
