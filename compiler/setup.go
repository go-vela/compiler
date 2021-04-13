// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"fmt"
	"strings"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
)

// Setup represents the configuration necessary for
// creating a Vela compiler engine capable of
// integrating with the configured template
// registries.
type Setup struct {
	// Compiler Configuration

	// specifies the driver to use for the compiler client
	Driver string
	// specifies the address to use for the compiler client
	Address string
	// specifies the token to use for the compiler client
	Token string

	// specifies the build to use for the compiler client
	Build *library.Build
	// specifies the comment to use for the compiler client
	Comment string
	// specifies the files to use for the compiler client
	Files []string
	// enables running local mode for the compiler client
	Local bool
	// specifies the metadata to use for the compiler client
	Metadata *types.Metadata
	// specifies the repo to use for the compiler client
	Repo *library.Repo
	// specifies the user to use for the compiler client
	User *library.User
}

// Native creates and returns a Vela compiler engine capable
// of integrating with Github as a template registry.
func (s *Setup) Native() (Engine, error) {
	logrus.Trace("creating native compiler client from setup")

	return nil, fmt.Errorf("unsupported compiler driver: %s", constants.DriverNative)
}

// Validate verifies the necessary fields for the
// provided configuration are populated correctly.
func (s *Setup) Validate() error {
	logrus.Trace("validating compiler setup for client")

	// verify a compiler driver was provided
	if len(s.Driver) == 0 {
		return fmt.Errorf("no compiler driver provided")
	}

	// verify a compiler address was provided
	if len(s.Address) == 0 {
		return fmt.Errorf("no compiler address provided")
	}

	// check if the compiler address has a scheme
	if !strings.Contains(s.Address, "://") {
		return fmt.Errorf("compiler address must be fully qualified (<scheme>://<host>)")
	}

	// check if the compiler address has a trailing slash
	if strings.HasSuffix(s.Address, "/") {
		return fmt.Errorf("compiler address must not have trailing slash")
	}

	// verify a compiler build was provided
	if s.Build == nil {
		return fmt.Errorf("no compiler build provided")
	}

	// verify a compiler metadata was provided
	if s.Metadata == nil {
		return fmt.Errorf("no compiler metadata provided")
	}

	// verify a compiler repo was provided
	if s.Repo == nil {
		return fmt.Errorf("no compiler repo provided")
	}

	// verify a compiler user was provided
	if s.User == nil {
		return fmt.Errorf("no compiler user provided")
	}

	// setup is valid
	return nil
}
