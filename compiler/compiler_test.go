// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"testing"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

func TestCompiler_New(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "https://github.com",
				Build:    new(library.Build),
				Metadata: new(types.Metadata),
				Repo:     new(library.Repo),
				User:     new(library.User),
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "",
				Build:    new(library.Build),
				Metadata: new(types.Metadata),
				Repo:     new(library.Repo),
				User:     new(library.User),
			},
		},
	}

	// run tests
	for _, test := range tests {
		_, err := New(test.setup)

		if test.failure {
			if err == nil {
				t.Errorf("New should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("New returned err: %v", err)
		}
	}
}
