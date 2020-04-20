// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"io/ioutil"
	"reflect"
	"testing"

	goyaml "gopkg.in/yaml.v2"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
)

func TestNative_Render_GoBasic(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/basic/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/basic/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/basic/tmpl.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_GoMultiline(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/multiline/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/multiline/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/multiline/tmpl.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_GoConditional_Match(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/conditional/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/conditional/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/conditional/tmpl.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_GoLoop_Map(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "golang",
			Variables: map[string]interface{}{
				"images": map[string]string{
					"1.11":   "golang:1.11",
					"1.12":   "golang:1.12",
					"latest": "golang:latest",
				},
				"pull_policy": "pull: true",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"go get ./..."},
			Name:     "sample_install",
			Image:    "golang:latest",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_1.11",
			Image:    "golang:1.11",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_1.12",
			Image:    "golang:1.12",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_latest",
			Image:    "golang:latest",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go build"},
			Name:     "sample_build",
			Environment: raw.StringSliceMap{
				"CGO_ENABLED": "0",
				"GOOS":        "linux",
			},
			Image: "golang:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/go_loop_map.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), s)

	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_GoLoop_Slice(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "golang",
			Variables: map[string]interface{}{
				"images":      []string{"golang:1.11", "golang:1.12", "golang:latest"},
				"pull_policy": "pull: true",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"go get ./..."},
			Name:     "sample_install_0",
			Image:    "golang:1.11",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_0",
			Image:    "golang:1.11",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go build"},
			Name:     "sample_build_0",
			Environment: raw.StringSliceMap{
				"CGO_ENABLED": "0",
				"GOOS":        "linux",
			},
			Image: "golang:1.11",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go get ./..."},
			Name:     "sample_install_1",
			Image:    "golang:1.12",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_1",
			Image:    "golang:1.12",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go build"},
			Name:     "sample_build_1",
			Environment: raw.StringSliceMap{
				"CGO_ENABLED": "0",
				"GOOS":        "linux",
			},
			Image: "golang:1.12",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go get ./..."},
			Name:     "sample_install_2",
			Image:    "golang:latest",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test_2",
			Image:    "golang:latest",
			Pull:     true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go build"},
			Name:     "sample_build_2",
			Environment: raw.StringSliceMap{
				"CGO_ENABLED": "0",
				"GOOS":        "linux",
			},
			Image: "golang:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/go_loop_slice.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), s)

	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_InvalidTemplate(t *testing.T) {
	// setup types
	want := yaml.StepSlice{}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/invalid_template.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), &yaml.Step{})

	if err == nil {
		t.Errorf("Render should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_InvalidVariables(t *testing.T) {
	// setup types
	want := yaml.StepSlice{}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/invalid_variables.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), &yaml.Step{})

	if err == nil {
		t.Errorf("Render should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_InvalidYml(t *testing.T) {
	// setup types
	want := yaml.StepSlice{}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/invalid.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), &yaml.Step{})

	if err == nil {
		t.Errorf("Render should have returned err")
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}
