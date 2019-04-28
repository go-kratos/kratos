/*
 *
 * Copyright 2018 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package google defines credentials for google cloud services.
package google

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/alts"
	"google.golang.org/grpc/credentials/oauth"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/internal"
)

const tokenRequestTimeout = 30 * time.Second

// NewDefaultCredentials returns a credentials bundle that is configured to work
// with google services.
//
// This API is experimental.
func NewDefaultCredentials() credentials.Bundle {
	c := &creds{}
	bundle, err := c.NewWithMode(internal.CredsBundleModeFallback)
	if err != nil {
		grpclog.Warningf("google default creds: failed to create new creds: %v", err)
	}
	return bundle
}

// creds implements credentials.Bundle.
type creds struct {
	// Supported modes are defined in internal/internal.go.
	mode string
	// The transport credentials associated with this bundle.
	transportCreds credentials.TransportCredentials
	// The per RPC credentials associated with this bundle.
	perRPCCreds credentials.PerRPCCredentials
}

func (c *creds) TransportCredentials() credentials.TransportCredentials {
	return c.transportCreds
}

func (c *creds) PerRPCCredentials() credentials.PerRPCCredentials {
	if c == nil {
		return nil
	}
	return c.perRPCCreds
}

// NewWithMode should make a copy of Bundle, and switch mode. Modifying the
// existing Bundle may cause races.
func (c *creds) NewWithMode(mode string) (credentials.Bundle, error) {
	newCreds := &creds{mode: mode}

	// Create transport credentials.
	switch mode {
	case internal.CredsBundleModeFallback:
		newCreds.transportCreds = credentials.NewTLS(nil)
	case internal.CredsBundleModeBackendFromBalancer, internal.CredsBundleModeBalancer:
		// Only the clients can use google default credentials, so we only need
		// to create new ALTS client creds here.
		newCreds.transportCreds = alts.NewClientCreds(alts.DefaultClientOptions())
	default:
		return nil, fmt.Errorf("google default creds: unsupported mode: %v", mode)
	}

	if mode == internal.CredsBundleModeFallback || mode == internal.CredsBundleModeBackendFromBalancer {
		// Create per RPC credentials.
		// For the time being, we required per RPC credentials for both TLS and
		// ALTS. In the future, this will only be required for TLS.
		ctx, cancel := context.WithTimeout(context.Background(), tokenRequestTimeout)
		defer cancel()
		var err error
		newCreds.perRPCCreds, err = oauth.NewApplicationDefault(ctx)
		if err != nil {
			grpclog.Warningf("google default creds: failed to create application oauth: %v", err)
		}
	}

	return newCreds, nil
}
