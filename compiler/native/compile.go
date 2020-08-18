// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/yaml"
)

// ModifyRequest contains the payload passed to the modification endpoint
type ModifyRequest struct {
	Pipeline *yaml.Build `json:"pipeline,omitempty"`
	Build    int         `json:"build,omitempty"`
	Repo     string      `json:"repo,omitempty"`
	User     string      `json:"user,omitempty"`
}

// modifyConfig sends the configuration to external http endpoint for modification.
func (c *client) modifyConfig(build *yaml.Build, libraryBuild *library.Build, repo *library.Repo) (*yaml.Build, error) {
	// create request to send to endpoint
	modReq := &ModifyRequest{
		Pipeline: build,
		Build:    libraryBuild.GetNumber(),
		Repo:     repo.GetName(),
		User:     libraryBuild.GetAuthor(),
	}

	// marshal json to send in request
	b, err := json.Marshal(modReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal modify payload")
	}

	// create POST request
	req, err := http.NewRequest("POST", c.ModificationService.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	// create the client with timeouts
	client := &http.Client{
		Timeout: c.ModificationService.Timeout,
	}

	// add content-type and auth headers
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.ModificationService.Secret))

	// send the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// fail if the response code was not 200
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("modification endpoint returned status code %v", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read payload")
	}

	newBuild := new(yaml.Build)
	// unmarshal the response into the yaml.Build struct
	err = json.Unmarshal(body, newBuild)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal modification payload")
	}

	return newBuild, nil
}

// Compile produces an executable pipeline from a yaml configuration.
func (c *client) Compile(v interface{}, build *library.Build, repo *library.Repo) (*pipeline.Build, error) {
	// parse the object into a yaml configuration
	p, err := c.Parse(v)
	if err != nil {
		return nil, err
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline: %w", err)
	}

	// create map of templates for easy lookup
	tmpls := mapFromTemplates(p.Templates)

	// create the ruledata to purge steps
	r := &pipeline.RuleData{
		Branch:  c.build.GetBranch(),
		Comment: c.comment,
		Event:   c.build.GetEvent(),
		Path:    c.files,
		Repo:    c.repo.GetFullName(),
		Tag:     strings.TrimPrefix(c.build.GetRef(), "refs/tags/"),
		Target:  c.build.GetDeploy(),
	}

	// inject the environment variables into the services
	p.Services, err = c.EnvironmentServices(p.Services)
	if err != nil {
		return nil, err
	}

	// inject the environment variables into the secrets
	p.Secrets, err = c.EnvironmentSecrets(p.Secrets)
	if err != nil {
		return nil, err
	}

	if len(p.Stages) > 0 {
		// inject the clone stage
		p, err = c.CloneStage(p)
		if err != nil {
			return nil, err
		}

		// inject the init stage
		p, err = c.InitStage(p)
		if err != nil {
			return nil, err
		}

		// inject the templates into the stages
		p.Stages, err = c.ExpandStages(p.Stages, tmpls)
		if err != nil {
			return nil, err
		}

		if c.ModificationService.Endpoint != "" {
			// send config to external endpoint for modification
			p, err = c.modifyConfig(p, build, repo)
			if err != nil {
				return nil, err
			}
		}

		// validate the yaml configuration
		err = c.Validate(p)
		if err != nil {
			return nil, fmt.Errorf("invalid pipeline: %w", err)
		}

		// inject the environment variables into the stages
		p.Stages, err = c.EnvironmentStages(p.Stages)
		if err != nil {
			return nil, err
		}

		// inject the substituted environment variables into the stages
		p.Stages, err = c.SubstituteStages(p.Stages)
		if err != nil {
			return nil, err
		}

		// inject the scripts into the stages
		p.Stages, err = c.ScriptStages(p.Stages)
		if err != nil {
			return nil, err
		}

		// return executable representation
		return c.TransformStages(r, p)
	}

	// inject the clone step
	p, err = c.CloneStep(p)
	if err != nil {
		return nil, err
	}

	// inject the init step
	p, err = c.InitStep(p)
	if err != nil {
		return nil, err
	}

	// inject the templates into the steps
	p.Steps, err = c.ExpandSteps(p.Steps, tmpls)
	if err != nil {
		return nil, err
	}

	if c.ModificationService.Endpoint != "" {
		// send config to external endpoint for modification
		p, err = c.modifyConfig(p, build, repo)
		if err != nil {
			return nil, err
		}
	}

	// validate the yaml configuration
	err = c.Validate(p)
	if err != nil {
		return nil, fmt.Errorf("invalid pipeline: %w", err)
	}

	// inject the environment variables into the steps
	p.Steps, err = c.EnvironmentSteps(p.Steps)
	if err != nil {
		return nil, err
	}

	// inject the substituted environment variables into the steps
	p.Steps, err = c.SubstituteSteps(p.Steps)
	if err != nil {
		return nil, err
	}

	// inject the scripts into the steps
	p.Steps, err = c.ScriptSteps(p.Steps)
	if err != nil {
		return nil, err
	}

	// return executable representation
	return c.TransformSteps(r, p)
}
