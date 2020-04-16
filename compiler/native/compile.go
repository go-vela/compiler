// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/types/pipeline"
)

// Compile produces an executable pipeline from a yaml configuration.
func (c *client) Compile(v interface{}) (*pipeline.Build, error) {
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
		Status:  c.build.GetStatus(),
		Tag:     strings.TrimPrefix(c.build.GetRef(), "refs/tags/"),
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
