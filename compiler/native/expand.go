// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/compiler/template/native"

	"github.com/go-vela/types/yaml"
	"github.com/sirupsen/logrus"
)

// ExpandStages injects the template for each
// templated step in every stage in a yaml configuration.
func (c *client) ExpandStages(s yaml.StageSlice, tmpls map[string]*yaml.Template) (yaml.StageSlice, error) {
	// iterate through all stages
	for _, stage := range s {
		// inject the templates into the steps for the stage
		steps, err := c.ExpandSteps(stage.Steps, tmpls)
		if err != nil {
			return nil, err
		}

		stage.Steps = steps
	}

	return s, nil
}

// ExpandSteps injects the template for each
// templated step in a yaml configuration.
func (c *client) ExpandSteps(s yaml.StepSlice, tmpls map[string]*yaml.Template) (yaml.StepSlice, error) {
	steps := yaml.StepSlice{}

	// iterate through each step
	for _, step := range s {
		bytes := []byte{}

		// lookup step template name
		tmpl, ok := tmpls[step.Template.Name]
		if !ok {
			// add existing step if no template
			steps = append(steps, step)
			continue
		}

		// skip processing template if the type isn't github
		if tmpl.Type != "github" {
			logrus.Errorf("Unsupported template type: %v", tmpl.Type)
			continue
		}

		// parse source from template
		src, err := c.Github.Parse(tmpl.Source, c.repo.GetBranch())
		if err != nil {
			return yaml.StepSlice{}, fmt.Errorf("invalid template source provided for %s: %v", step.Template.Name, err)
		}

		// pull from public github when the host isn't provided or is set to github.com
		if len(src.Host) == 0 || strings.Contains(src.Host, "github.com") {
			bytes, err = c.Github.Template(nil, src)
			if err != nil {
				return yaml.StepSlice{}, err
			}
		}

		// pull from private github installation if the host is not empty
		if len(src.Host) > 0 {
			bytes, err = c.PrivateGithub.Template(c.user, src)
			if err != nil {
				return yaml.StepSlice{}, err
			}
		}

		// render template for steps
		tmplSteps, err := native.Render(string(bytes), step)
		if err != nil {
			return yaml.StepSlice{}, err
		}

		// add templated steps
		steps = append(steps, tmplSteps...)
	}

	return steps, nil
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
