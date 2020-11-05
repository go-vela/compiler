// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"reflect"
	"testing"

	"github.com/go-vela/types/raw"
	"github.com/stretchr/testify/assert"
	"go.starlark.net/starlark"
)

func TestNative_Render_convertTemplateVars(t *testing.T) {
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

	strWant := starlark.NewDict(0)
	strWant.SetKey(starlark.String("pull"), starlark.String("always"))

	arrayWant := starlark.NewDict(0)
	arrayWant.SetKey(starlark.String("tags"), tags)

	mapWant := starlark.NewDict(0)
	mapWant.SetKey(starlark.String("commands"), commands)

	tests := []struct {
		name string
		args map[string]interface{}
		want *starlark.Dict
	}{
		{
			name: "test for a user passed string",
			args: map[string]interface{}{"pull": "always"},
			want: strWant,
		},
		{
			name: "test for a user passed array",
			args: map[string]interface{}{"tags": []string{"latest", "1.14", "1.15"}},
			want: arrayWant,
		},
		{
			name: "test for a user passed map",
			// nolint // ignore line length
			args: map[string]interface{}{"commands": map[string]string{"test": "go test ./...", "build": "go build"}},
			want: mapWant,
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertTemplateVars(tt.args)
			assert.NoError(t, err)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertTemplateVars() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNative_Render_velaEnvironmentData(t *testing.T) {
	// setup types
	build := starlark.NewDict(1)
	err := build.SetKey(starlark.String("author"), starlark.String("octocat"))
	assert.NoError(t, err)

	repo := starlark.NewDict(1)
	err = repo.SetKey(starlark.String("full_name"), starlark.String("go-vela/hello-world"))
	assert.NoError(t, err)

	user := starlark.NewDict(1)
	err = user.SetKey(starlark.String("admin"), starlark.String("true"))
	assert.NoError(t, err)

	system := starlark.NewDict(1)
	err = system.SetKey(starlark.String("workspace"), starlark.String("/vela/src/github.com/go-vela/hello-world"))
	assert.NoError(t, err)

	withAllPre := starlark.NewDict(0)
	withAllPre.SetKey(starlark.String("build"), build)
	withAllPre.SetKey(starlark.String("repo"), repo)
	withAllPre.SetKey(starlark.String("user"), user)
	withAllPre.SetKey(starlark.String("system"), system)

	tests := []struct {
		name    string
		args    raw.StringSliceMap
		want    *starlark.Dict
		wantErr bool
	}{
		{
			name: "with all vela prefixed var",
			args: raw.StringSliceMap{
				"VELA_BUILD_AUTHOR":   "octocat",
				"VELA_REPO_FULL_NAME": "go-vela/hello-world",
				"VELA_USER_ADMIN":     "true",
				"VELA_WORKSPACE":      "/vela/src/github.com/go-vela/hello-world",
			},
			want: withAllPre,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := convertPlatformVars(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertPlatformVars() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertPlatformVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
