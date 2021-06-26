// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package github

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-vela/compiler/registry"

	"github.com/go-vela/types/library"

	"github.com/google/go-github/v36/github"
)

// Template captures the templated pipeline configuration from the GitHub repo.
func (c *client) Template(u *library.User, s *registry.Source) ([]byte, error) {
	// use default GitHub OAuth client we provide
	cli := c.Github
	if u != nil {
		// create GitHub OAuth client with user's token
		cli = c.newClientToken(u.GetToken())
	}

	// create the options to pass
	opts := &github.RepositoryContentGetOptions{}

	// set the reference for the options to capture the templated pipeline
	// configuration. if no ref is set, it will pull from the default
	// branch on the targeted repo, see:
	// https://docs.github.com/en/rest/reference/repos#get-repository-content--parameters
	if len(s.Ref) > 0 {
		opts.Ref = s.Ref
	}

	// send API call to capture the templated pipeline configuration
	//
	// nolint: lll // ignore long line length due to variable names
	data, _, resp, err := cli.Repositories.GetContents(context.Background(), s.Org, s.Repo, s.Name, opts)
	if err != nil {
		if resp.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	// data is not nil if template exists
	if data != nil {
		strData, err := data.GetContent()
		if err != nil {
			return nil, err
		}

		return []byte(strData), nil
	}

	return nil, fmt.Errorf("no valid template found at %s/%s/%s", s.Org, s.Repo, s.Name)
}
