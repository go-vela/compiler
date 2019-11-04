// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/compiler/compiler"

	"github.com/go-vela/compiler/registry"
	"github.com/go-vela/compiler/registry/github"

	"github.com/go-vela/types/library"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

type Client struct {
	Github        registry.Service
	PrivateGithub registry.Service

	build *library.Build
	files []string
	repo  *library.Repo
	user  *library.User
}

// New returns a Pipeline implementation that integrates with the supported registries.
func New(ctx *cli.Context) (*Client, error) {
	logrus.Debug("Creating registry clients from CLI configuration")
	c := Client{}

	// setup github template service
	github, err := setupGithub()
	if err != nil {
		return nil, err
	}
	c.Github = github

	if ctx.Bool("github-driver") {
		// setup private github service
		privGithub, err := setupPrivateGithub(ctx.String("github-url"), ctx.String("github-token"))
		if err != nil {
			return nil, err
		}
		c.PrivateGithub = privGithub
	}

	return &c, nil
}

// setupGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupGithub() (registry.Service, error) {
	logrus.Tracef("Creating %s registry client from CLI configuration", "github")
	return github.New("", "")
}

// setupPrivateGithub is a helper function to setup the
// Github registry service from the CLI arguments.
func setupPrivateGithub(addr, token string) (registry.Service, error) {
	logrus.Tracef("Creating private %s registry client from CLI configuration", "github")
	return github.New(addr, token)
}

// WithBuild sets the library build type in the Engine.
func (c *Client) WithBuild(b *library.Build) compiler.Engine {
	if b != nil {
		c.build = b
	}

	return c
}

// WithFiles sets the changeset files in the Engine.
func (c *Client) WithFiles(f []string) compiler.Engine {
	if f != nil {
		c.files = f
	}

	return c
}

// WithRepo sets the library repo type in the Engine.
func (c *Client) WithRepo(r *library.Repo) compiler.Engine {
	if r != nil {
		c.repo = r
	}

	return c
}

// WithUser sets the library user type in the Engine.
func (c *Client) WithUser(u *library.User) compiler.Engine {
	if u != nil {
		c.user = u
	}

	return c
}
