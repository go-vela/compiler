// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"io/ioutil"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	goyaml "github.com/goccy/go-yaml"
	"github.com/google/go-cmp/cmp"
)

func TestStarlark_Render(t *testing.T) {
	type args struct {
		velaFile     string
		starlarkFile string
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{"basic", args{velaFile: "testdata/basic/step.yml", starlarkFile: "testdata/basic/template.py"}, "testdata/basic/want.yml", false},
		{"with method", args{velaFile: "testdata/with_method/step.yml", starlarkFile: "testdata/with_method/template.star"}, "testdata/with_method/want.yml", false},
		{"user vars", args{velaFile: "testdata/with_vars/step.yml", starlarkFile: "testdata/with_vars/template.star"}, "testdata/with_vars/want.yml", false},
		{"platform vars", args{velaFile: "testdata/with_vars_plat/step.yml", starlarkFile: "testdata/with_vars_plat/template.star"}, "testdata/with_vars_plat/want.yml", false},
		{"cancel due to complexity", args{velaFile: "testdata/cancel/step.yml", starlarkFile: "testdata/cancel/template.star"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sFile, err := ioutil.ReadFile(tt.args.velaFile)
			if err != nil {
				t.Error(err)
			}
			b := &yaml.Build{}
			err = goyaml.Unmarshal(sFile, b)
			if err != nil {
				t.Error(err)
			}
			b.Steps[0].Environment = raw.StringSliceMap{
				"VELA_REPO_FULL_NAME": "octocat/hello-world",
			}

			tmpl, err := ioutil.ReadFile(tt.args.starlarkFile)
			if err != nil {
				t.Error(err)
			}

			got, err := Render(string(tmpl), b.Steps[0])
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr != true {
				wFile, err := ioutil.ReadFile(tt.wantFile)
				if err != nil {
					t.Error(err)
				}
				w := &yaml.Build{}
				err = goyaml.Unmarshal(wFile, w)
				if err != nil {
					t.Error(err)
				}
				want := w.Steps

				if diff := cmp.Diff(want, got); diff != "" {
					t.Errorf("ExpandSteps() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
