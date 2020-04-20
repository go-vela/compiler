// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/types"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli/v2"
)

func TestNative_EnvironmentStages(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"

	s := yaml.StageSlice{
		&yaml.Stage{
			Name: str,
			Steps: yaml.StepSlice{
				&yaml.Step{
					Image: "alpine",
					Name:  str,
					Pull:  true,
				},
			},
		},
	}

	want := yaml.StageSlice{
		&yaml.Stage{
			Name: str,
			Steps: yaml.StepSlice{
				&yaml.Step{
					Environment: environment(nil, nil, nil, nil),
					Image:       "alpine",
					Name:        str,
					Pull:        true,
				},
			},
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentStages(s)
	if err != nil {
		t.Errorf("EnvironmentStages returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EnvironmentStages is %v, want %v", got, want)
	}
}

func TestNative_EnvironmentSteps(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	str := "foo"
	s := yaml.StepSlice{
		&yaml.Step{
			Image: "alpine",
			Name:  str,
			Pull:  true,
		},
	}

	want := yaml.StepSlice{
		&yaml.Step{
			Environment: environment(nil, nil, nil, nil),
			Image:       "alpine",
			Name:        str,
			Pull:        true,
		},
	}

	// run test
	compiler, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	got, err := compiler.EnvironmentSteps(s)
	if err != nil {
		t.Errorf("EnvironmentSteps returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("EnvironmentSteps is %v, want %v", got, want)
	}
}

func TestNative_environment(t *testing.T) {
	// setup types
	booL := false
	num := 1
	num64 := int64(num)
	str := "foo"
	workspace := "/home/foo/foo"
	// push
	push := "push"
	// tag
	tag := "tag"
	tagref := "refs/tags/1"
	// pull_request
	pull := "pull_request"
	pullref := "refs/pull/1/head"

	tests := []struct {
		w    string
		b    *library.Build
		m    *types.Metadata
		r    *library.Repo
		u    *library.User
		want map[string]string
	}{
		// push
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &push, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &str, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "push", "BUILD_FINISHED": "1", "BUILD_HOST": "TODO", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "foo", "BUILD_STARTED": "1", "BUILD_SOURCE": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/home/foo/foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/home/foo/foo", "CI": "vela", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false"},
		},
		// tag
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &tag, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &tagref, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "tag", "BUILD_FINISHED": "1", "BUILD_HOST": "TODO", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_REF": "refs/tags/1", "BUILD_STARTED": "1", "BUILD_SOURCE": "foo", "BUILD_TAG": "1", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/home/foo/foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/home/foo/foo", "CI": "vela", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false"},
		},
		// pull_request
		{
			w:    workspace,
			b:    &library.Build{ID: &num64, RepoID: &num64, Number: &num, Parent: &num, Event: &pull, Status: &str, Error: &str, Enqueued: &num64, Created: &num64, Started: &num64, Finished: &num64, Deploy: &str, Clone: &str, Source: &str, Title: &str, Message: &str, Commit: &str, Sender: &str, Author: &str, Branch: &str, Ref: &pullref, BaseRef: &str},
			m:    &types.Metadata{Database: &types.Database{Driver: str, Host: str}, Queue: &types.Queue{Channel: str, Driver: str, Host: str}, Source: &types.Source{Driver: str, Host: str}, Vela: &types.Vela{Address: str, WebAddress: str}},
			r:    &library.Repo{ID: &num64, UserID: &num64, Org: &str, Name: &str, FullName: &str, Link: &str, Clone: &str, Branch: &str, Timeout: &num64, Visibility: &str, Private: &booL, Trusted: &booL, Active: &booL, AllowPull: &booL, AllowPush: &booL, AllowDeploy: &booL, AllowTag: &booL, AllowComment: &booL},
			u:    &library.User{ID: &num64, Name: &str, Token: &str, Active: &booL, Admin: &booL},
			want: map[string]string{"BUILD_AUTHOR": "foo", "BUILD_AUTHOR_EMAIL": "", "BUILD_BRANCH": "foo", "BUILD_CHANNEL": "foo", "BUILD_COMMIT": "foo", "BUILD_CREATED": "1", "BUILD_ENQUEUED": "1", "BUILD_EVENT": "pull_request", "BUILD_FINISHED": "1", "BUILD_HOST": "TODO", "BUILD_LINK": "", "BUILD_MESSAGE": "foo", "BUILD_NUMBER": "1", "BUILD_PARENT": "1", "BUILD_PULL_REQUEST_NUMBER": "1", "BUILD_REF": "refs/pull/1/head", "BUILD_STARTED": "1", "BUILD_SOURCE": "foo", "BUILD_TITLE": "foo", "BUILD_WORKSPACE": "/home/foo/foo", "VELA": "true", "VELA_ADDR": "foo", "VELA_CHANNEL": "foo", "VELA_DATABASE": "foo", "VELA_DISTRIBUTION": "TODO", "VELA_HOST": "foo", "VELA_NETRC_MACHINE": "foo", "VELA_NETRC_PASSWORD": "foo", "VELA_NETRC_USERNAME": "x-oauth-basic", "VELA_QUEUE": "foo", "VELA_RUNTIME": "TODO", "VELA_SOURCE": "foo", "VELA_VERSION": "TODO", "VELA_WORKSPACE": "/home/foo/foo", "CI": "vela", "REPOSITORY_BRANCH": "foo", "REPOSITORY_CLONE": "foo", "REPOSITORY_FULL_NAME": "foo", "REPOSITORY_LINK": "foo", "REPOSITORY_NAME": "foo", "REPOSITORY_ORG": "foo", "REPOSITORY_PRIVATE": "false", "REPOSITORY_TIMEOUT": "1", "REPOSITORY_TRUSTED": "false"},
		},
	}

	// run test
	for _, test := range tests {
		got := environment(test.b, test.m, test.r, test.u)

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("environment is %v, want %v", got, test.want)
		}
	}
}
