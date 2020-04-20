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
	"github.com/google/go-cmp/cmp"
)

func TestNative_Render_GoBasic(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "golang",
			Variables: map[string]interface{}{
				"image":       "golang:latest",
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
			Name:     "sample_test",
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
	tmpl, err := ioutil.ReadFile("testdata/go_basic.yml")
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

func TestNative_Render_GoMultiline(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/go/multiline/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/go/multiline/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/go/multiline/tmpl.yml")
	if err != nil {
		t.Errorf("Reading file returned err: %v", err)
	}

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
		}
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_GoConditional_Match(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "golang",
			Variables: map[string]interface{}{
				"image":       "golang:latest",
				"pull_policy": "pull: true",
				"branch":      "master",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"go get ./..."},
			Name:     "sample_install",
			Image:    "golang:latest",
			Pull:     true,
		},
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test",
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
	tmpl, err := ioutil.ReadFile("testdata/go_conditional.yml")
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

func TestNative_Render_GoConditional_NoMatch(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "golang",
			Variables: map[string]interface{}{
				"image":       "golang:latest",
				"pull_policy": "pull: true",
				"branch":      "dev",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"go test ./..."},
			Name:     "sample_test",
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
	tmpl, err := ioutil.ReadFile("testdata/go_conditional.yml")
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

func TestNative_Render_JavaBasic(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "gradle",
			Variables: map[string]interface{}{
				"image":       "openjdk:latest",
				"environment": "{ \"GRADLE_USER_HOME\": \".gradle\", \"GRADLE_OPTS\": \"-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false\"}",
				"pull_policy": "pull: true",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/java_basic.yml")
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

func TestNative_Render_JavaConditional_Match(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "gradle",
			Variables: map[string]interface{}{
				"image":       "openjdk:latest",
				"environment": "{ \"GRADLE_USER_HOME\": \".gradle\", \"GRADLE_OPTS\": \"-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false\"}",
				"pull_policy": "pull: true",
				"branch":      "master",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install",
			Image: "openjdk:latest",
			Pull:  true,
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/java_conditional.yml")
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

func TestNative_Render_JavaConditional_NoMatch(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "gradle",
			Variables: map[string]interface{}{
				"image":       "openjdk:latest",
				"environment": "{ \"GRADLE_USER_HOME\": \".gradle\", \"GRADLE_OPTS\": \"-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false\"}",
				"pull_policy": "pull: true",
				"branch":      "dev",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/java_conditional.yml")
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

func TestNative_Render_JavaLoop_Map(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "gradle",
			Variables: map[string]interface{}{
				"images": map[string]string{
					"12":     "openjdk:12",
					"13":     "openjdk:13",
					"latest": "openjdk:latest",
				},
				"environment": "{ \"GRADLE_USER_HOME\": \".gradle\", \"GRADLE_OPTS\": \"-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false\"}",
				"pull_policy": "pull: true",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_12",
			Image: "openjdk:12",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_13",
			Image: "openjdk:13",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_latest",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/java_loop_map.yml")
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

func TestNative_Render_JavaLoop_Slice(t *testing.T) {
	// setup types
	s := &yaml.Step{
		Name: "sample",
		Template: yaml.StepTemplate{
			Name: "gradle",
			Variables: map[string]interface{}{
				"images":      []string{"openjdk:12", "openjdk:13", "openjdk:latest"},
				"environment": "{ \"GRADLE_USER_HOME\": \".gradle\", \"GRADLE_OPTS\": \"-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false\"}",
				"pull_policy": "pull: true",
			},
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install_0",
			Image: "openjdk:12",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_0",
			Image: "openjdk:12",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build_0",
			Image: "openjdk:12",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install_1",
			Image: "openjdk:13",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_1",
			Image: "openjdk:13",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build_1",
			Image: "openjdk:13",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew downloadDependencies"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_install_2",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew check"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_test_2",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
		&yaml.Step{
			Commands: raw.StringSlice{"./gradlew build distTar"},
			Environment: raw.StringSliceMap{
				"GRADLE_USER_HOME": ".gradle",
				"GRADLE_OPTS":      "-Dorg.gradle.daemon=false -Dorg.gradle.workers.max=1 -Dorg.gradle.parallel=false",
			},
			Name:  "sample_build_2",
			Image: "openjdk:latest",
			Pull:  true,
			Ruleset: yaml.Ruleset{
				If:       yaml.Rules{Event: []string{"push", "pull_request"}},
				Operator: "and",
			},
		},
	}

	// run test
	tmpl, err := ioutil.ReadFile("testdata/java_loop_slice.yml")
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
