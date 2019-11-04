// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"

	"github.com/urfave/cli"
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
					Environment: environment(nil, nil, nil),
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
			Environment: environment(nil, nil, nil),
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
	b := &library.Build{
		ID:       &num64,
		RepoID:   &num64,
		Number:   &num,
		Parent:   &num,
		Event:    &str,
		Status:   &str,
		Error:    &str,
		Enqueued: &num64,
		Created:  &num64,
		Started:  &num64,
		Finished: &num64,
		Deploy:   &str,
		Clone:    &str,
		Source:   &str,
		Title:    &str,
		Message:  &str,
		Commit:   &str,
		Sender:   &str,
		Author:   &str,
		Branch:   &str,
		Ref:      &str,
		BaseRef:  &str,
	}
	r := &library.Repo{
		ID:          &num64,
		UserID:      &num64,
		Org:         &str,
		Name:        &str,
		FullName:    &str,
		Link:        &str,
		Clone:       &str,
		Branch:      &str,
		Timeout:     &num64,
		Visibility:  &str,
		Private:     &booL,
		Trusted:     &booL,
		Active:      &booL,
		AllowPull:   &booL,
		AllowPush:   &booL,
		AllowDeploy: &booL,
		AllowTag:    &booL,
	}
	u := &library.User{
		ID:     &num64,
		Name:   &str,
		Token:  &str,
		Active: &booL,
		Admin:  &booL,
	}

	workspace := fmt.Sprintf("/home/%s_%s_%d", r.GetOrg(), r.GetName(), b.GetNumber())

	want := map[string]string{
		"BUILD_BRANCH":         b.GetBranch(),
		"BUILD_CHANNEL":        "vela",
		"BUILD_COMMIT":         b.GetCommit(),
		"BUILD_CREATED":        unmarshal(b.GetCreated()),
		"BUILD_ENQUEUED":       unmarshal(b.GetEnqueued()),
		"BUILD_EVENT":          b.GetEvent(),
		"BUILD_FINISHED":       unmarshal(b.GetFinished()),
		"BUILD_HOST":           "TODO",
		"BUILD_MESSAGE":        b.GetMessage(),
		"BUILD_NUMBER":         unmarshal(b.GetNumber()),
		"BUILD_PARENT":         unmarshal(b.GetParent()),
		"BUILD_REF":            b.GetRef(),
		"BUILD_STARTED":        unmarshal(b.GetStarted()),
		"BUILD_SOURCE":         b.GetSource(),
		"BUILD_TITLE":          b.GetTitle(),
		"BUILD_WORKSPACE":      workspace,
		"VELA":                 unmarshal(true),
		"VELA_ADDR":            "TODO",
		"VELA_CHANNEL":         "vela",
		"VELA_DATABASE":        "postgres",
		"VELA_DISTRIBUTION":    "linux",
		"VELA_HOST":            "TODO",
		"VELA_NETRC_MACHINE":   "github.com",
		"VELA_NETRC_PASSWORD":  u.GetToken(),
		"VELA_NETRC_USERNAME":  "x-oauth-basic",
		"VELA_QUEUE":           "redis",
		"VELA_RUNTIME":         "docker",
		"VELA_SOURCE":          "https://github.com",
		"VELA_VERSION":         "TODO",
		"VELA_WORKSPACE":       workspace,
		"CI":                   "vela",
		"REPOSITORY_BRANCH":    r.GetBranch(),
		"REPOSITORY_CLONE":     r.GetClone(),
		"REPOSITORY_FULL_NAME": r.GetFullName(),
		"REPOSITORY_LINK":      r.GetLink(),
		"REPOSITORY_NAME":      r.GetName(),
		"REPOSITORY_ORG":       r.GetOrg(),
		"REPOSITORY_PRIVATE":   unmarshal(r.GetPrivate()),
		"REPOSITORY_TIMEOUT":   unmarshal(r.GetTimeout()),
		"REPOSITORY_TRUSTED":   unmarshal(r.GetTrusted()),
	}

	// run test
	got := environment(b, r, u)

	if !reflect.DeepEqual(got, want) {
		t.Errorf("environment is %v, want %v", got, want)
	}
}
