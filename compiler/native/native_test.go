// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"flag"
	"reflect"
	"testing"

	"github.com/go-vela/lexi/registry/github"

	"github.com/go-vela/types/library"

	"github.com/urfave/cli"
)

func TestNative_New(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New("", "")
	want := &Client{
		Github: public,
	}

	// run test
	got, err := New(c)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_New_PrivateGithub(t *testing.T) {
	// setup types
	url := "http://foo.example.com"
	token := "someToken"
	set := flag.NewFlagSet("test", 0)
	set.Bool("github-driver", true, "doc")
	set.String("github-url", url, "doc")
	set.String("github-token", token, "doc")
	c := cli.NewContext(nil, set, nil)
	public, _ := github.New("", "")
	private, _ := github.New(url, token)
	want := &Client{
		Github:        public,
		PrivateGithub: private,
	}

	// run test
	got, err := New(c)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestNative_WithBuild(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	b := &library.Build{ID: &id}

	want, _ := New(c)
	want.build = b

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithBuild(b), want) {
		t.Errorf("WithBuild is %v, want %v", got, want)
	}
}

func TestNative_WithFiles(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	f := []string{"foo"}

	want, _ := New(c)
	want.files = f

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithFiles(f), want) {
		t.Errorf("WithFiles is %v, want %v", got, want)
	}
}

func TestNative_WithRepo(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	r := &library.Repo{ID: &id}

	want, _ := New(c)
	want.repo = r

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithRepo(r), want) {
		t.Errorf("WithRepo is %v, want %v", got, want)
	}
}

func TestNative_WithUser(t *testing.T) {
	// setup types
	set := flag.NewFlagSet("test", 0)
	c := cli.NewContext(nil, set, nil)

	id := int64(1)
	u := &library.User{ID: &id}

	want, _ := New(c)
	want.user = u

	// run test
	got, err := New(c)
	if err != nil {
		t.Errorf("Unable to create new compiler: %v", err)
	}

	if !reflect.DeepEqual(got.WithUser(u), want) {
		t.Errorf("WithUser is %v, want %v", got, want)
	}
}
