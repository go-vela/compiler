// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
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

func TestNative_Render_userData(t *testing.T) {
	// setup types
	tags := starlark.Tuple(nil)
	tags = append(tags, starlark.String("latest"))
	tags = append(tags, starlark.String("1.14"))
	tags = append(tags, starlark.String("1.15"))

	commands := starlark.NewDict(16)
	err := commands.SetKey(starlark.String("test"), starlark.String("go test ./..."))
	assert.NoError(t, err)
	err = commands.SetKey(starlark.String("build"), starlark.String("go build"))
	assert.NoError(t, err)

	tests := []struct {
		name string
		args map[string]interface{}
		want starlark.StringDict
	}{
		{
			name: "test for a user passed string",
			args: map[string]interface{}{"pull": "always"},
			want: starlark.StringDict{"pull": starlark.String("always")},
		},
		{
			name: "test for a user passed array",
			args: map[string]interface{}{"tags": []string{"latest", "1.14", "1.15"}},
			want: starlark.StringDict{"tags": tags},
		},
		{
			name: "test for a user passed map",
			args: map[string]interface{}{"commands": map[string]string{"test": "go test ./...", "build": "go build"}},
			want: starlark.StringDict{"commands": commands},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := userData(tt.args)
			assert.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userData() = %v, want %v", got, tt.want)
			}
		})
	}
}
