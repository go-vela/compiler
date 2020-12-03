// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"io/ioutil"
	"reflect"
	"testing"

	goyaml "github.com/goccy/go-yaml"

	"github.com/go-vela/types/yaml"
)

func TestNative_Render_Basic(t *testing.T) {
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

func TestNative_Render_Multiline(t *testing.T) {
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

func TestNative_Render_Conditional_Match(t *testing.T) {
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

func TestNative_Render_Loop_Map(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/loop_map/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/loop_map/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/loop_map/tmpl.yml")
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

func TestNative_Render_Loop_Slice(t *testing.T) {
	// setup types
	sFile, _ := ioutil.ReadFile("testdata/loop_slice/step.yml")
	b := &yaml.Build{}
	_ = goyaml.Unmarshal(sFile, b)

	wFile, _ := ioutil.ReadFile("testdata/loop_slice/want.yml")
	w := &yaml.Build{}
	_ = goyaml.Unmarshal(wFile, w)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/loop_slice/tmpl.yml")
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
