// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/types"
	"github.com/go-vela/types/constants"
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/yaml"
)

// EnvironmentStages injects environment variables
// for each step in every stage in a yaml configuration.
func (c *client) EnvironmentStages(s yaml.StageSlice) (yaml.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		// inject the environment variables into the steps for the stage
		steps, err := c.EnvironmentSteps(stage.Steps)
		if err != nil {
			return nil, err
		}

		stage.Steps = steps
	}

	return s, nil
}

// EnvironmentSteps injects environment variables
// for each step in a yaml configuration.
func (c *client) EnvironmentSteps(s yaml.StepSlice) (yaml.StepSlice, error) {
	// iterate through all steps
	for _, step := range s {
		// make empty map of environment variables
		env := make(map[string]string)
		// gather set of default environment variables
		defaultEnv := environment(c.build, c.metadata, c.repo, c.user)

		// inject the declared environment
		// variables to the build step
		for k, v := range step.Environment {
			env[k] = v
		}

		// inject the default environment
		// variables to the build step
		// we do this after injecting the declared environment
		// to ensure the default env overrides any conflicts
		for k, v := range defaultEnv {
			env[k] = v
		}

		// inject the declared parameter
		// variables to the build step
		for k, v := range step.Parameters {
			if v == nil {
				continue
			}

			// parameter keys are passed to the image
			// as PARAMETER_ environment variables
			k = "PARAMETER_" + strings.ToUpper(k)

			// parameter values are passed to the image
			// as string environment variables
			env[k] = unmarshal(v)
		}

		// overwrite existing build step environment
		step.Environment = env
	}

	return s, nil
}

// EnvironmentServices injects environment variables
// for each service in a yaml configuration.
func (c *client) EnvironmentServices(s yaml.ServiceSlice) (yaml.ServiceSlice, error) {
	// iterate through all services
	for _, service := range s {
		// make empty map of environment variables
		env := make(map[string]string)
		// gather set of default environment variables
		defaultEnv := environment(c.build, c.metadata, c.repo, c.user)

		// inject the declared environment
		// variables to the build service
		for k, v := range service.Environment {
			env[k] = v
		}

		// inject the default environment
		// variables to the build service
		// we do this after injecting the declared environment
		// to ensure the default env overrides any conflicts
		for k, v := range defaultEnv {
			env[k] = v
		}

		// overwrite existing build service environment
		service.Environment = env
	}

	return s, nil
}

// EnvironmentSecrets injects environment variables
// for each secret plugin in a yaml configuration.
func (c *client) EnvironmentSecrets(s yaml.SecretSlice) (yaml.SecretSlice, error) {
	// iterate through all secrets
	for _, secret := range s {
		// skip non plugin secrets
		if secret.Origin.Empty() {
			continue
		}

		// make empty map of environment variables
		env := make(map[string]string)
		// gather set of default environment variables
		defaultEnv := environment(c.build, c.metadata, c.repo, c.user)

		// inject the declared environment
		// variables to the build secret
		for k, v := range secret.Origin.Environment {
			env[k] = v
		}

		// inject the default environment
		// variables to the build secret
		// we do this after injecting the declared environment
		// to ensure the default env overrides any conflicts
		for k, v := range defaultEnv {
			env[k] = v
		}

		// inject the declared parameter
		// variables to the build secret
		for k, v := range secret.Origin.Parameters {
			if v == nil {
				continue
			}

			// parameter keys are passed to the image
			// as PARAMETER_ environment variables
			k = "PARAMETER_" + strings.ToUpper(k)

			// parameter values are passed to the image
			// as string environment variables
			env[k] = unmarshal(v)
		}

		// overwrite existing build secret environment
		secret.Origin.Environment = env
	}

	return s, nil
}

// helper function that creates the standard set of environment variables for a pipeline.
func environment(b *library.Build, m *types.Metadata, r *library.Repo, u *library.User) map[string]string {
	workspace := "/vela"

	if m != nil {
		workspace = fmt.Sprintf("/vela/src/%s/%s/%s", m.Source.Host, r.GetOrg(), r.GetName())
	}

	env := map[string]string{
		// build specific environment variables
		"BUILD_AUTHOR":       b.GetAuthor(),
		"BUILD_AUTHOR_EMAIL": b.GetEmail(),
		"BUILD_BRANCH":       b.GetBranch(),
		"BUILD_CHANNEL":      "TODO",
		"BUILD_COMMIT":       b.GetCommit(),
		"BUILD_CREATED":      unmarshal(b.GetCreated()),
		"BUILD_ENQUEUED":     unmarshal(b.GetEnqueued()),
		"BUILD_EVENT":        b.GetEvent(),
		"BUILD_FINISHED":     unmarshal(b.GetFinished()),
		"BUILD_HOST":         "TODO",
		"BUILD_LINK":         b.GetLink(),
		"BUILD_MESSAGE":      b.GetMessage(),
		"BUILD_NUMBER":       unmarshal(b.GetNumber()),
		"BUILD_PARENT":       unmarshal(b.GetParent()),
		"BUILD_REF":          b.GetRef(),
		"BUILD_STARTED":      unmarshal(b.GetStarted()),
		"BUILD_SOURCE":       b.GetSource(),
		"BUILD_TITLE":        b.GetTitle(),
		"BUILD_WORKSPACE":    workspace,

		// vela specific environment variables
		"VELA":                unmarshal(true),
		"VELA_ADDR":           "TODO",
		"VELA_CHANNEL":        "TODO",
		"VELA_DATABASE":       "TODO",
		"VELA_DISTRIBUTION":   "TODO",
		"VELA_HOST":           "TODO",
		"VELA_NETRC_MACHINE":  "TODO",
		"VELA_NETRC_PASSWORD": u.GetToken(),
		"VELA_NETRC_USERNAME": "x-oauth-basic",
		"VELA_QUEUE":          "TODO",
		"VELA_RUNTIME":        "TODO",
		"VELA_SOURCE":         "TODO",
		"VELA_VERSION":        "TODO",
		"VELA_WORKSPACE":      workspace,
		"CI":                  "vela",

		// repo specific environment variables
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

	// set tag environment variable if proper build event
	if b.GetEvent() == constants.EventTag {
		env["BUILD_TAG"] = strings.SplitN(b.GetRef(), "refs/tags/", 2)[1]
	}

	// set pull request number variable if proper build event
	if b.GetEvent() == constants.EventPull {
		env["BUILD_PULL_REQUEST_NUMBER"] = strings.SplitN(b.GetRef(), "/", 4)[2]
	}

	// set deployment target if proper build event
	if b.GetEvent() == constants.EventDeploy {
		env["BUILD_TARGET"] = b.GetDeploy()
	}

	// populate environment variables from metadata
	if m != nil {
		env["BUILD_CHANNEL"] = m.Queue.Channel
		env["VELA_ADDR"] = m.Vela.WebAddress
		env["VELA_CHANNEL"] = m.Queue.Channel
		env["VELA_DATABASE"] = m.Database.Driver
		env["VELA_HOST"] = m.Vela.Address
		env["VELA_NETRC_MACHINE"] = m.Source.Host
		env["VELA_QUEUE"] = m.Queue.Driver
		env["VELA_SOURCE"] = m.Source.Driver
	}

	return env
}
