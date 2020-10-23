// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	goyaml "gopkg.in/yaml.v2"
)

func TestNative_Render_StarlarkBasic(t *testing.T) {
	// setup types
	sFile, err := ioutil.ReadFile("testdata/basic/step.yml")
	assert.NoError(t, err)
	b := &yaml.Build{}
	err = goyaml.Unmarshal(sFile, b)
	assert.NoError(t, err)

	wFile, err := ioutil.ReadFile("testdata/basic/want.yml")
	assert.NoError(t, err)
	w := &yaml.Build{}
	err = goyaml.Unmarshal(wFile, w)
	assert.NoError(t, err)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/basic/template.py")
	assert.NoError(t, err)

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

// ctx.vela.build_number
// ctx.vela.build.number
// ctx.vela.build.foo_bar
// ctx.vela.host
// ctx.vela.system.<thing>

// ctx.vela.host (value)
// ctx.vela.build (dict)

// if prexix == VELA_
// remove VELA_ add to vela dict
// if prefix == BUILD_
// do y

func TestNative_Render_StarlarkWithMethod(t *testing.T) {
	// setup types
	sFile, err := ioutil.ReadFile("testdata/with_method/step.yml")
	assert.NoError(t, err)
	b := &yaml.Build{}
	err = goyaml.Unmarshal(sFile, b)
	assert.NoError(t, err)

	wFile, err := ioutil.ReadFile("testdata/with_method/want.yml")
	assert.NoError(t, err)
	w := &yaml.Build{}
	err = goyaml.Unmarshal(wFile, w)
	assert.NoError(t, err)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/with_method/template.star")
	assert.NoError(t, err)

	got, err := Render(string(tmpl), b.Steps[0])
	if err != nil {
		t.Errorf("Render returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Render is %v, want %v", got, want)
	}
}

func TestNative_Render_StarlarkWithVars(t *testing.T) {
	// setup types
	sFile, err := ioutil.ReadFile("testdata/with_vars/step.yml")
	assert.NoError(t, err)
	b := &yaml.Build{}
	err = goyaml.Unmarshal(sFile, b)
	assert.NoError(t, err)

	wFile, err := ioutil.ReadFile("testdata/with_vars/want.yml")
	assert.NoError(t, err)
	w := &yaml.Build{}
	err = goyaml.Unmarshal(wFile, w)
	assert.NoError(t, err)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/with_vars/template.star")
	assert.NoError(t, err)

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

func TestNative_Render_StarlarkWithVarsPlat(t *testing.T) {
	// setup types
	sFile, err := ioutil.ReadFile("testdata/with_vars_plat/step.yml")
	assert.NoError(t, err)
	b := &yaml.Build{}
	err = goyaml.Unmarshal(sFile, b)
	assert.NoError(t, err)

	b.Steps[0].Environment = raw.StringSliceMap{
		"VELA_REPO_FULL_NAME": "octocat/hello-world",
	}

	wFile, err := ioutil.ReadFile("testdata/with_vars_plat/want.yml")
	assert.NoError(t, err)
	w := &yaml.Build{}
	err = goyaml.Unmarshal(wFile, w)
	assert.NoError(t, err)

	want := w.Steps

	// run test
	tmpl, err := ioutil.ReadFile("testdata/with_vars_plat/template.star")
	assert.NoError(t, err)

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
