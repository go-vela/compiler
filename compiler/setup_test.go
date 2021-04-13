// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"reflect"
	"testing"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
)

func TestCompiler_Setup_Native(t *testing.T) {
	// setup types
	_setup := &Setup{
		Driver:   "native",
		Address:  "https://github.com",
		Build:    new(library.Build),
		Metadata: new(types.Metadata),
		Repo:     new(library.Repo),
		User:     new(library.User),
	}

	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
		want    Engine
	}{
		{
			failure: true,
			setup:   _setup,
			want:    nil,
		},
	}

	// run tests
	for _, test := range tests {
		got, err := test.setup.Native()

		if test.failure {
			if err == nil {
				t.Errorf("Native should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Native returned err: %v", err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("Native is %v, want %v", got, test.want)
		}
	}
}

func TestSource_Setup_Validate(t *testing.T) {
	// setup tests
	tests := []struct {
		failure bool
		setup   *Setup
	}{
		{
			failure: false,
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
				Driver:   "",
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
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "https://github.com/",
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
				Address:  "github.com",
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
				Address:  "https://github.com",
				Build:    nil,
				Metadata: new(types.Metadata),
				Repo:     new(library.Repo),
				User:     new(library.User),
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "https://github.com",
				Build:    new(library.Build),
				Metadata: nil,
				Repo:     new(library.Repo),
				User:     new(library.User),
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "https://github.com",
				Build:    new(library.Build),
				Metadata: new(types.Metadata),
				Repo:     nil,
				User:     new(library.User),
			},
		},
		{
			failure: true,
			setup: &Setup{
				Driver:   "native",
				Address:  "https://github.com",
				Build:    new(library.Build),
				Metadata: new(types.Metadata),
				Repo:     new(library.Repo),
				User:     nil,
			},
		},
	}

	// run tests
	for _, test := range tests {
		err := test.setup.Validate()

		if test.failure {
			if err == nil {
				t.Errorf("Validate should have returned err")
			}

			continue
		}

		if err != nil {
			t.Errorf("Validate returned err: %v", err)
		}
	}
}
