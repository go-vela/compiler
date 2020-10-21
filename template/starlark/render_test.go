// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"io/ioutil"
	"reflect"
	"testing"

	goyaml "gopkg.in/yaml.v2"

	"github.com/go-vela/types/yaml"
	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"
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
		if diff := cmp.Diff(want, got); diff != "" {
			t.Errorf("MakeGatewayInfo() mismatch (-want +got):\n%s", diff)
		}
		t.Errorf("Render is %v, want %v", got, want)
	}
}
