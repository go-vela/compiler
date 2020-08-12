// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/types"
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
			env[k] = library.ToString(v)
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

// helper function to merge two maps together.
func mergeMap(combinedMap, loopMap map[string]string) map[string]string {
	for key, value := range loopMap {
		combinedMap[key] = value
	}

	return combinedMap
}

// helper function that creates the standard set of environment variables for a pipeline.
func environment(b *library.Build, m *types.Metadata, r *library.Repo, u *library.User) map[string]string {
	workspace := "/vela"

	env := make(map[string]string)

	// vela specific environment variables
	env["VELA"] = library.ToString(true)
	env["VELA_ADDR"] = "TODO"
	env["VELA_CHANNEL"] = "TODO"
	env["VELA_DATABASE"] = "TODO"
	env["VELA_DISTRIBUTION"] = "TODO"
	env["VELA_HOST"] = "TODO"
	env["VELA_NETRC_MACHINE"] = "TODO"
	env["VELA_NETRC_PASSWORD"] = u.GetToken()
	env["VELA_NETRC_USERNAME"] = "x-oauth-basic"
	env["VELA_QUEUE"] = "TODO"
	env["VELA_RUNTIME"] = "TODO"
	env["VELA_SOURCE"] = "TODO"
	env["VELA_VERSION"] = "TODO"
	env["CI"] = "vela"

	// populate environment variables from metadata
	if m != nil {
		env["VELA_ADDR"] = m.Vela.WebAddress
		env["VELA_CHANNEL"] = m.Queue.Channel
		env["VELA_DATABASE"] = m.Database.Driver
		env["VELA_HOST"] = m.Vela.Address
		env["VELA_NETRC_MACHINE"] = m.Source.Host
		env["VELA_QUEUE"] = m.Queue.Driver
		env["VELA_SOURCE"] = m.Source.Driver
		workspace = fmt.Sprintf("/vela/src/%s/%s/%s", m.Source.Host, r.GetOrg(), r.GetName())
	}

	env["VELA_WORKSPACE"] = workspace

	// populate environment variables from repo library
	mergeMap(env, r.Environment())
	// populate environment variables from build library
	mergeMap(env, b.Environment(workspace))
	// populate environment variables from user library
	mergeMap(env, u.Environment())

	// TODO: add this to types

	if m != nil {
		env["BUILD_CHANNEL"] = m.Queue.Channel
	}

	return env
}
