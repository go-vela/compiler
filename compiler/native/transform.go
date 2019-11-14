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
func (c *client) TransformStages(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
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

	// set the unique ID for each step in each stage of the executable pipeline
	for _, stage := range pipeline.Stages {
		for _, step := range stage.Steps {
			// create pattern for steps
			pattern := fmt.Sprintf(stageID, org, name, number, stage.Name, step.Name)

			// set id to the pattern
			step.ID = pattern

			// sanitize the step to ensure ID field is valid
			step.Sanitize()
		}
	}

	// set the unique ID for each service in the executable pipeline
	for _, service := range pipeline.Services {
		// create pattern for services
		pattern := fmt.Sprintf(serviceID, org, name, number, service.Name)

		// set id to the pattern
		service.ID = pattern

		// sanitize the service to ensure ID field is valid
		service.Sanitize()
	}

	return pipeline.Purge(r), nil
}

// TransformSteps converts a yaml configuration with steps into an executable pipeline.
func (c *client) TransformSteps(r *pipeline.RuleData, p *yaml.Build) (*pipeline.Build, error) {
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

	// set the unique ID for each step in the executable pipeline
	for _, step := range pipeline.Steps {
		// create pattern for steps
		pattern := fmt.Sprintf(stepID, org, name, number, step.Name)

		// set id to the pattern
		step.ID = pattern

		// sanitize the step to ensure ID field is valid
		step.Sanitize()
	}

	// set the unique ID for each service in the executable pipeline
	for _, service := range pipeline.Services {
		// create pattern for services
		pattern := fmt.Sprintf(serviceID, org, name, number, service.Name)

		// set id to the pattern
		service.ID = pattern

		// sanitize the service to ensure ID field is valid
		service.Sanitize()
	}

	return pipeline.Purge(r), nil
}
