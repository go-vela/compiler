// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/yaml"
)

const (
	// default ID for pipeline.
	// format: `<org>_<repo>_<build number>`
	pipelineID = "%s_%s_%d"
	// default ID for steps in a stage in a pipeline.
	// format: `<org name>_<repo name>_<build number>_<stage name>_<step name>`
	stageID = "%s_%s_%d_%s_%s"
	// default ID for steps in a pipeline.
	// format: `step_<org name>_<repo name>_<build number>_<step name>`
	stepID = "step_%s_%s_%d_%s"
	// default ID for steps in a pipeline.
	// format: `service_<org name>_<repo name>_<build number>_<service name>`
	serviceID = "service_%s_%s_%d_%s"
)

// TransformStages converts a yaml configuration with stages into an executable pipeline.
func (c *Client) TransformStages(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
	// capture variables for setting the unique ID fields
	org := c.repo.GetOrg()
	name := c.repo.GetName()
	number := c.build.GetNumber()

	// create new executable pipeline
	pipeline := &pipeline.Build{
		Version:  p.Version,
		Metadata: *p.Metadata.ToPipeline(),
		Stages:   *p.Stages.ToPipeline(),
		Secrets:  *p.Secrets.ToPipeline(),
		Services: *p.Services.ToPipeline(),
	}

	// set the unique ID for the executable pipeline
	pipeline.ID = fmt.Sprintf(pipelineID, org, name, number)

	// set the unique ID for every step in every stage of the executable pipeline
	for i, stage := range p.Stages {
		for j, step := range stage.Steps {
			pipeline.Stages[i].Steps[j].ID = fmt.Sprintf(stageID, org, name, number, stage.Name, step.Name)
		}
	}

	// set the unique ID for each service in the pipeline
	for _, service := range pipeline.Services {
		service.ID = fmt.Sprintf(serviceID, org, name, number, service.Name)
	}

	return pipeline.Purge(r), nil
}

// TransformSteps converts a yaml configuration with steps into an executable pipeline.
func (c *Client) TransformSteps(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
	// capture variables for setting the unique ID fields
	org := c.repo.GetOrg()
	name := c.repo.GetName()
	number := c.build.GetNumber()

	// create new executable pipeline
	pipeline := &pipeline.Build{
		Version:  p.Version,
		Metadata: *p.Metadata.ToPipeline(),
		Steps:    *p.Steps.ToPipeline(),
		Secrets:  *p.Secrets.ToPipeline(),
		Services: *p.Services.ToPipeline(),
	}

	// set the unique ID for the executable pipeline
	pipeline.ID = fmt.Sprintf(pipelineID, org, name, number)

	// set the unique ID for every step of the executable pipeline
	for i, step := range p.Steps {
		pipeline.Steps[i].ID = fmt.Sprintf(stepID, org, name, number, step.Name)
	}

	// set the unique ID for each service in the pipeline
	for _, service := range pipeline.Services {
		service.ID = fmt.Sprintf(serviceID, org, name, number, service.Name)
	}

	return pipeline.Purge(r), nil
}
