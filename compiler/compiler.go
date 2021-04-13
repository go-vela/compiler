// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"fmt"

	"github.com/go-vela/types/constants"

	"github.com/sirupsen/logrus"
)

// nolint: godot // top level comment ends in a list
//
// New creates and returns a Vela compiler engine capable
// of parsing and producing pipelines while integrating
// with the configured template registries.
//
// Currently the following compiler engines are supported:
//
// * Native
func New(s *Setup) (Engine, error) {
	// validate the setup being provided
	//
	// https://pkg.go.dev/github.com/go-vela/compiler/compiler?tab=doc#Setup.Validate
	err := s.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Debug("creating compiler engine from setup")
	// process the compiler driver being provided
	switch s.Driver {
	case constants.DriverNative:
		// handle the native compiler driver being provided
		//
		// https://pkg.go.dev/github.com/go-vela/compiler/compiler?tab=doc#Setup.Native
		return s.Native()
	default:
		// handle an invalid compiler driver being provided
		return nil, fmt.Errorf("invalid compiler driver provided: %s", s.Driver)
	}
}
