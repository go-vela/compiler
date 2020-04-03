// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli"

	"github.com/kr/pretty"
)

func TestNative_SubstituteStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	s := yaml.StageSlice{
		{
			Name: "simple",
			Steps: yaml.StepSlice{
				{
					Commands:    []string{"echo ${FOO}", "echo $${BAR}"},
					Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
					Image:       "alpine:latest",
					Name:        "simple",
					Pull:        true,
				},
			},
		},
		{
			Name: "advanced",
			Steps: yaml.StepSlice{
				{
					Commands:    []string{"echo ${COMPLEX}"},
					Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
					Image:       "alpine:latest",
					Name:        "advanced",
					Pull:        true,
				},
			},
		},
	}

	want := yaml.StageSlice{
		{
			Name: "simple",
			Steps: yaml.StepSlice{
				{
					Commands:    []string{"echo baz", "echo ${BAR}"},
					Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
					Image:       "alpine:latest",
					Name:        "simple",
					Pull:        true,
				},
			},
		},
		{
			Name: "advanced",
			Steps: yaml.StepSlice{
				{
					Commands:    []string{"echo \"{\\\"hello\\\":\\n  \\\"world\\\"}\""},
					Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
					Image:       "alpine:latest",
					Name:        "advanced",
					Pull:        true,
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.SubstituteStages(s)
	if err != nil {
		t.Errorf("SubstituteStages returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("SubstituteStages is %v, want %v", got, want)
	}
}

func TestNative_SubstituteSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := yaml.StepSlice{
		{
			Commands:    []string{"echo ${FOO}", "echo $${BAR}"},
			Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
			Image:       "alpine:latest",
			Name:        "simple",
			Pull:        true,
		},
		{
			Commands:    []string{"echo ${COMPLEX}"},
			Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
			Image:       "alpine:latest",
			Name:        "advanced",
			Pull:        true,
		},
	}

	want := yaml.StepSlice{
		{
			Commands:    []string{"echo baz", "echo ${BAR}"},
			Environment: map[string]string{"FOO": "baz", "BAR": "baz"},
			Image:       "alpine:latest",
			Name:        "simple",
			Pull:        true,
		},
		{
			Commands:    []string{"echo \"{\\\"hello\\\":\\n  \\\"world\\\"}\""},
			Environment: map[string]string{"COMPLEX": "{\"hello\":\n  \"world\"}"},
			Image:       "alpine:latest",
			Name:        "advanced",
			Pull:        true,
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.SubstituteSteps(p)
	if err != nil {
		t.Errorf("SubstituteSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		pretty.Ldiff(t, got, want)
		t.Errorf("SubstituteSteps is %v, want %v", got, want)
	}
}
