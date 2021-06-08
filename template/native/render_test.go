// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"io/ioutil"
	"testing"

	goyaml "github.com/buildkite/yaml"
	"github.com/google/go-cmp/cmp"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
)

func TestNative_Render(t *testing.T) {
	type args struct {
		velaFile     string
		templateFile string
	}
	tests := []struct {
		name     string
		args     args
		wantFile string
		wantErr  bool
	}{
		{"basic", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/basic/tmpl.yml"}, "testdata/basic/want.yml", false},
		{"multiline", args{velaFile: "testdata/multiline/step.yml", templateFile: "testdata/multiline/tmpl.yml"}, "testdata/multiline/want.yml", false},
		{"conditional match", args{velaFile: "testdata/conditional/step.yml", templateFile: "testdata/conditional/tmpl.yml"}, "testdata/conditional/want.yml", false},
		{"loop map", args{velaFile: "testdata/loop_map/step.yml", templateFile: "testdata/loop_map/tmpl.yml"}, "testdata/loop_map/want.yml", false},
		{"loop slice", args{velaFile: "testdata/loop_slice/step.yml", templateFile: "testdata/loop_slice/tmpl.yml"}, "testdata/loop_slice/want.yml", false},
		{"platform vars", args{velaFile: "testdata/with_vars_plat/step.yml", templateFile: "testdata/with_vars_plat/tmpl.yml"}, "testdata/with_vars_plat/want.yml", false},
		{"invalid template", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/invalid_template.yml"}, "", true},
		{"invalid variable", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/invalid_variables.yml"}, "", true},
		{"invalid yml", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/invalid.yml"}, "", true},
		{"disallowed env func", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/disallowed/tmpl_env.yml"}, "", true},
		{"disallowed expandenv func", args{velaFile: "testdata/basic/step.yml", templateFile: "testdata/disallowed/tmpl_expandenv.yml"}, "", true},
		{"yaml anchor", args{velaFile: "testdata/anchor/step.yml", templateFile: "testdata/anchor/tmpl.yml"}, "", false},
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

			tmpl, err := ioutil.ReadFile(tt.args.templateFile)
			if err != nil {
				t.Error(err)
			}

			fmt.Println(b.Steps[0])

			steps, secrets, services, err := Render(string(tmpl), b.Steps[0])
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
				wantSteps := w.Steps
				wantSecrets := w.Secrets
				wantServices := w.Services

				if diff := cmp.Diff(wantSteps, steps); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantSecrets, secrets); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
				if diff := cmp.Diff(wantServices, services); diff != "" {
					t.Errorf("Render() mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
