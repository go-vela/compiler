// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types/yaml"
	"github.com/urfave/cli"
)

func TestNative_CloneStage(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "alpine",
						Name:  str,
						Pull:  true,
					},
				},
			},
		},
	}

	want := &yaml.Build{
		Version: "v1",
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: "clone",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "docker.target.com/vela-plugins/git:1",
						Name:  "clone",
						Pull:  true,
					},
				},
			},
			&yaml.Stage{
				Name: str,
				Steps: yaml.StepSlice{
					&yaml.Step{
						Image: "alpine",
						Name:  str,
						Pull:  true,
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.CloneStage(p)
	if err != nil {
		t.Errorf("CloneStage returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CloneStage is %v, want %v", got, want)
	}
}

func TestNative_CloneStep(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	p := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "alpine",
				Name:  str,
				Pull:  true,
			},
		},
	}

	want := &yaml.Build{
		Version: "v1",
		Steps: yaml.StepSlice{
			&yaml.Step{
				Image: "docker.target.com/vela-plugins/git:1",
				Name:  "clone",
				Pull:  true,
			},
			&yaml.Step{
				Image: "alpine",
				Name:  str,
				Pull:  true,
			},
		},
	}
	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.CloneStep(p)
	if err != nil {
		t.Errorf("CloneStep returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("CloneStep is %v, want %v", got, want)
	}
}
