// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli"
)

func TestNative_TransformStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Ports: []string{"5432:5432"},
				Name:  "postgres backend",
				Image: "postgres:latest",
			},
		},
		Worker: yaml.Worker{
			Name:    "worker_1",
			Flavor:  "16cpu8gb",
			Runtime: "dockers",
		},
		Stages: yaml.StageSlice{
			&yaml.Stage{
				Name: "install deps",
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew downloadDependencies"},
						Image:    "openjdk:latest",
						Name:     "install",
						Pull:     true,
					},
				},
			},
			&yaml.Stage{
				Name:  "test",
				Needs: []string{"install"},
				Steps: yaml.StepSlice{
					&yaml.Step{
						Commands: []string{"./gradlew check"},
						Image:    "openjdk:latest",
						Name:     "test",
						Pull:     true,
						Ruleset: yaml.Ruleset{
							If: yaml.Rules{
								Event: []string{"push"},
							},
							Operator: "and",
						},
					},
				},
			},
		},
	}

	want := &pipeline.Build{
		ID:      "__0",
		Version: "v1",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:     "service___0_postgres-backend",
				Ports:  []string{"5432:5432"},
				Name:   "postgres backend",
				Image:  "postgres:latest",
				Number: 1,
			},
		},
		Worker: pipeline.Worker{
			Name:    "worker_1",
			Flavor:  "16cpu8gb",
			Runtime: "dockers",
		},
		Stages: pipeline.StageSlice{
			&pipeline.Stage{
				Name: "install deps",
				Steps: pipeline.ContainerSlice{
					&pipeline.Container{
						ID:       "__0_install-deps_install",
						Commands: []string{"./gradlew downloadDependencies"},
						Image:    "openjdk:latest",
						Name:     "install",
						Number:   1,
						Pull:     true,
					},
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.TransformStages(new(pipeline.RuleData), p)
	if err != nil {
		t.Errorf("TransformStages returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TransformStages is %v, want %v", got, want)
	}
}

func TestNative_TransformSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	p := &yaml.Build{
		Version: "v1",
		Services: yaml.ServiceSlice{
			&yaml.Service{
				Ports: []string{"5432:5432"},
				Name:  "postgres backend",
				Image: "postgres:latest",
			},
		},
		Worker: yaml.Worker{
			Name:    "worker_1",
			Flavor:  "16cpu8gb",
			Runtime: "dockers",
		},
		Steps: yaml.StepSlice{
			&yaml.Step{
				Commands: []string{"./gradlew downloadDependencies"},
				Image:    "openjdk:latest",
				Name:     "install deps",
				Pull:     true,
			},
			&yaml.Step{
				Commands: []string{"./gradlew check"},
				Image:    "openjdk:latest",
				Name:     "test",
				Pull:     true,
				Ruleset: yaml.Ruleset{
					If: yaml.Rules{
						Event: []string{"push"},
					},
					Operator: "and",
				},
			},
		},
	}

	want := &pipeline.Build{
		ID:      "__0",
		Version: "v1",
		Services: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:     "service___0_postgres-backend",
				Ports:  []string{"5432:5432"},
				Name:   "postgres backend",
				Image:  "postgres:latest",
				Number: 1,
			},
		},
		Worker: pipeline.Worker{
			Name:    "worker_1",
			Flavor:  "16cpu8gb",
			Runtime: "dockers",
		},
		Steps: pipeline.ContainerSlice{
			&pipeline.Container{
				ID:       "step___0_install-deps",
				Commands: []string{"./gradlew downloadDependencies"},
				Image:    "openjdk:latest",
				Name:     "install deps",
				Number:   1,
				Pull:     true,
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Creating compiler returned err: %v", err)
	}

	got, err := compiler.TransformSteps(new(pipeline.RuleData), p)
	if err != nil {
		t.Errorf("TransformSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("TransformSteps is %v, want %v", got, want)
	}
}
