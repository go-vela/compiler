// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"github.com/go-vela/types/yaml"
)

const (
	// default image for clone process.
	cloneImage = "docker.target.com/vela-plugins/git:1"
	// default name for clone stage.
	cloneStageName = "clone"
	// default name for clone step.
	cloneStepName = "clone"
)

// CloneStage injects the stage clone process into a yaml configuration.
func (c *Client) CloneStage(p *yaml.Build) (*yaml.Build, error) {
	stages := yaml.StageSlice{}

	// create new clone stage
	clone := &yaml.Stage{
		Name: cloneStageName,
		Steps: yaml.StepSlice{
			&yaml.Step{
				Detach:     false,
				Image:      cloneImage,
				Name:       cloneStepName,
				Privileged: false,
				Pull:       true,
			},
		},
	}

	// add clone stage as first stage
	stages = append(stages, clone)

	// add existing stages after clone stage
	for _, stage := range p.Stages {
		stages = append(stages, stage)
	}

	// overwrite existing stages
	p.Stages = stages

	return p, nil
}

// CloneStep injects the step clone process into a yaml configuration.
func (c *Client) CloneStep(p *yaml.Build) (*yaml.Build, error) {
	steps := yaml.StepSlice{}

	// create new clone step
	clone := &yaml.Step{
		Detach:     false,
		Image:      cloneImage,
		Name:       cloneStepName,
		Privileged: false,
		Pull:       true,
	}

	// add clone step as first step
	steps = append(steps, clone)

	// add existing steps after clone step
	for _, step := range p.Steps {
		steps = append(steps, step)
	}

	// overwrite existing steps
	p.Steps = steps

	return p, nil
}
