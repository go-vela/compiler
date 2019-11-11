// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"

	"github.com/go-vela/types/yaml"
)

// Validate verifies the yaml configuration is accurate.
func (c *client) Validate(p *yaml.Build) error {
	if len(p.Version) == 0 {
		return fmt.Errorf("no version provided")
	}

	if len(p.Stages) == 0 && len(p.Steps) == 0 {
		return fmt.Errorf("no stages or steps provided")
	}

	if len(p.Stages) > 0 && len(p.Steps) > 0 {
		return fmt.Errorf("stages and steps provided")
	}

	for _, service := range p.Services {
		if len(service.Name) == 0 {
			return fmt.Errorf("no name provided for service")
		}

		if len(service.Image) == 0 {
			return fmt.Errorf("no image provided for service %s", service.Name)
		}
	}

	for _, stage := range p.Stages {
		if len(stage.Name) == 0 {
			return fmt.Errorf("no name provided for stage")
		}

		for _, step := range stage.Steps {
			if len(step.Name) == 0 {
				return fmt.Errorf("no name provided for step for stage %s", stage.Name)
			}

			if len(step.Image) == 0 && len(step.Template.Name) == 0 {
				return fmt.Errorf("no image or template provided for step %s for stage %s", step.Name, stage.Name)
			}

			if len(step.Commands) == 0 && len(step.Parameters) == 0 &&
				len(step.Template.Name) == 0 && !step.Detach {
				return fmt.Errorf("no commands or parameters or template provided for step %s for stage %s", step.Name, stage.Name)
			}
		}
	}

	for _, step := range p.Steps {
		if len(step.Name) == 0 {
			return fmt.Errorf("no name provided for step")
		}

		if len(step.Image) == 0 && len(step.Template.Name) == 0 {
			return fmt.Errorf("no image or template provided for step %s", step.Name)
		}

		if len(step.Commands) == 0 && len(step.Parameters) == 0 &&
			len(step.Template.Name) == 0 && !step.Detach {
			return fmt.Errorf("no commands or parameters or template provided for step %s", step.Name)
		}
	}

	return nil
}
