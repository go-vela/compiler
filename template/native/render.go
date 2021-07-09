package native

import (
	"bytes"
	"fmt"
	"text/template"

	types "github.com/go-vela/types/yaml"

	"github.com/Masterminds/sprig/v3"

	yaml "github.com/buildkite/yaml"
)

//TODO: remove this when done.
func Render(template string, step *types.Step) (types.StepSlice, error) {
	return nil, nil
}

// RenderStep combines the template with the step in the yaml pipeline.
// nolint: lll // ignore long line length due to return args
func RenderStep(tmpl string, s *types.Step) (types.StepSlice, types.SecretSlice, types.ServiceSlice, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(s.Environment, s.Name)}
	templateFuncMap := map[string]interface{}{
		"vela": velaFuncs.returnPlatformVar,
	}
	// modify Masterminds/sprig functions
	// to remove OS functions
	//
	// https://masterminds.github.io/sprig/os.html
	sf := sprig.TxtFuncMap()
	delete(sf, "env")
	delete(sf, "expandenv")

	// parse the template with Masterminds/sprig functions
	//
	// https://pkg.go.dev/github.com/Masterminds/sprig?tab=doc#TxtFuncMap
	t, err := template.New(s.Name).Funcs(sf).Funcs(templateFuncMap).Parse(tmpl)
	if err != nil {
		// nolint: lll // ignore long line length due to return arguments
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, fmt.Errorf("unable to parse template %s: %v", s.Template.Name, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, s.Template.Variables)
	if err != nil {
		// nolint: lll // ignore long line length due to return arguments
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, fmt.Errorf("unable to execute template %s: %v", s.Template.Name, err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		// nolint: lll // ignore long line length due to return args
		return types.StepSlice{}, types.SecretSlice{}, types.ServiceSlice{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	return config.Steps, config.Secrets, config.Services, nil
}

// RenderBuild renders the templated build
func RenderBuild(b string, envs map[string]string) (*types.Build, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(envs)}
	templateFuncMap := map[string]interface{}{
		"vela": velaFuncs.returnPlatformVar,
	}
	// modify Masterminds/sprig functions
	// to remove OS functions
	//
	// https://masterminds.github.io/sprig/os.html
	sf := sprig.TxtFuncMap()
	delete(sf, "env")
	delete(sf, "expandenv")

	// parse the template with Masterminds/sprig functions
	//
	// https://pkg.go.dev/github.com/Masterminds/sprig?tab=doc#TxtFuncMap
	t, err := template.New("build").Funcs(sf).Funcs(templateFuncMap).Parse(b)
	if err != nil {
		return nil, err
	}

	// execute the template
	err = t.Execute(buffer, "")
	if err != nil {
		return nil, fmt.Errorf("unable to execute template: %w", err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal yaml: %w", err)
	}

	return config, nil
}
