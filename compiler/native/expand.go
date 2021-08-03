// Copyright (c) 2021 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package native

import (
	"fmt"
	"strings"

	"github.com/go-vela/compiler/template/native"
	"github.com/go-vela/compiler/template/starlark"
	"github.com/spf13/afero"

	"github.com/go-vela/types/raw"
	"github.com/go-vela/types/yaml"
	"github.com/sirupsen/logrus"
)

// ExpandStages injects the template for each
// templated step in every stage in a yaml configuration.
//
// nolint: lll // ignore long line length due to variable names
func (c *client) ExpandStages(s *yaml.Build, tmpls map[string]*yaml.Template) (yaml.StageSlice, yaml.SecretSlice, yaml.ServiceSlice, error) {
	// iterate through all stages
	for _, stage := range s.Stages {
		// inject the templates into the steps for the stage
		steps, secrets, services, err := c.ExpandSteps(&yaml.Build{Steps: stage.Steps, Secrets: s.Secrets, Services: s.Services}, tmpls)
		if err != nil {
			return nil, nil, nil, err
		}

		stage.Steps = steps
		s.Secrets = secrets
		s.Services = services
	}

	return s.Stages, s.Secrets, s.Services, nil
}

// ExpandSteps injects the template for each
// templated step in a yaml configuration.
//
// nolint: lll,funlen // ignore long line length due to variable names
func (c *client) ExpandSteps(s *yaml.Build, tmpls map[string]*yaml.Template) (yaml.StepSlice, yaml.SecretSlice, yaml.ServiceSlice, error) {
	steps := yaml.StepSlice{}
	secrets := s.Secrets
	services := s.Services

	// iterate through each step
	for _, step := range s.Steps {
		bytes := []byte{}

		// skip if no template is provided for the step
		if len(step.Template.Name) == 0 {
			// add existing step if no template
			steps = append(steps, step)
			continue
		}

		// lookup step template name
		tmpl, ok := tmpls[step.Template.Name]
		if !ok {
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, fmt.Errorf("missing template source for template %s in pipeline for step %s", step.Template.Name, step.Name)
		}

		// Create some default global environment inject vars
		// these are used below to overwrite to an empty
		// map if they should not be injected into a container
		envGlobalSteps := s.Environment

		if !s.Metadata.HasEnvironment("steps") {
			envGlobalSteps = make(raw.StringSliceMap)
		}

		// inject environment information for template
		step, err := c.EnvironmentStep(step, envGlobalSteps)
		if err != nil {
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
		}

		// nolint:ineffassign // ignoring linter on overwriting the "bytes" var
		switch {
		case c.local:
			a := &afero.Afero{
				Fs: afero.NewOsFs(),
			}

			bytes, err = a.ReadFile(tmpl.Source)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
			}

		case strings.EqualFold(tmpl.Type, "github"):
			// parse source from template
			src, err := c.Github.Parse(tmpl.Source)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, fmt.Errorf("invalid template source provided for %s: %v", step.Template.Name, err)
			}

			// pull from public github when the host isn't provided or is set to github.com
			if len(src.Host) == 0 || strings.Contains(src.Host, "github.com") {
				bytes, err = c.Github.Template(nil, src)
				if err != nil {
					return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
				}
			}

			// pull from private github installation if the host is not empty
			if len(src.Host) > 0 {
				bytes, err = c.PrivateGithub.Template(c.user, src)
				if err != nil {
					return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
				}
			}

		default:
			logrus.Errorf("Unsupported template type: %v", tmpl.Type)
			continue
		}

		var tmplSteps yaml.StepSlice
		var tmplSecrets yaml.SecretSlice
		var tmplServices yaml.ServiceSlice

		// TODO: provide friendlier error messages with file type mismatches
		switch tmpl.Format {
		case "go", "golang", "":
			// render template for steps
			tmplSteps, tmplSecrets, tmplServices, err = native.RenderStep(string(bytes), step)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
			}
		case "starlark":
			// render template for steps
			tmplSteps, tmplSecrets, tmplServices, err = starlark.RenderStep(string(bytes), step)
			if err != nil {
				return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, err
			}
		default:
			return yaml.StepSlice{}, yaml.SecretSlice{}, yaml.ServiceSlice{}, fmt.Errorf("format of %s is unsupported", tmpl.Format)
		}

		// loop over secrets within template
		for _, secret := range tmplSecrets {
			found := false
			// loop over secrets within base configuration
			for _, sec := range secrets {
				// check if the template secret and base secret name match
				if sec.Name == secret.Name {
					found = true
				}
			}

			// only append template secret if it does not exist within base configuration
			if !secret.Origin.Empty() || !found {
				secrets = append(secrets, secret)
			}
		}

		// loop over services within template
		for _, service := range tmplServices {
			found := false
			for _, serv := range services {
				if serv.Name == service.Name {
					found = true
				}
			}

			// only append template service if it does not exist within base configuration
			if !found {
				services = append(services, service)
			}
		}

		// add templated steps
		steps = append(steps, tmplSteps...)
	}

	return steps, secrets, services, nil
}

// helper function that creates a map of templates from a yaml configuration.
func mapFromTemplates(templates []*yaml.Template) map[string]*yaml.Template {
	m := make(map[string]*yaml.Template)

	for _, tmpl := range templates {
		m[tmpl.Name] = tmpl
	}

	return m
}
