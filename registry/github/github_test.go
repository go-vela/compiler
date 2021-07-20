// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"testing"

	"github.com/google/go-github/v37/github"
	"golang.org/x/oauth2"
)

func TestGithub_New(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	gitClient := github.NewClient(nil)
	gitClient.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	want := &client{
		Github: gitClient,
		URL:    s.URL,
		API:    s.URL + "/api/v3/",
	}

	// run test
	got, err := New(s.URL, "")

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %v, want %v", got, want)
	}
}

func TestGithub_NewToken(t *testing.T) {
	// setup router
	s := httptest.NewServer(http.NotFoundHandler())
	defer s.Close()

	token := "foobar"
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	gitClient := github.NewClient(tc)
	gitClient.BaseURL, _ = url.Parse(s.URL + "/api/v3/")

	want := &client{
		Github: gitClient,
		URL:    s.URL,
		API:    s.URL + "/api/v3/",
	}

	// run test
	got, err := New(s.URL, token)

	if err != nil {
		t.Errorf("New returned err: %v", err)
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("New is %+v, want %+v", got.Github, want.Github)
	}
}
