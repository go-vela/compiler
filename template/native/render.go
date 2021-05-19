package native

import (
	"bytes"
	"fmt"
	"text/template"

	types "github.com/go-vela/types/yaml"

	"github.com/Masterminds/sprig/v3"

	yaml "gopkg.in/yaml.v2"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, s *types.Step) (types.StepSlice, types.SecretSlice, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	velaFuncs := funcHandler{envs: convertPlatformVars(s.Environment)}
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
		return types.StepSlice{}, types.SecretSlice{}, fmt.Errorf("unable to parse template %s: %v", s.Template.Name, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, s.Template.Variables)
	if err != nil {
		// nolint: lll // ignore long line length due to return arguments
		return types.StepSlice{}, types.SecretSlice{}, fmt.Errorf("unable to execute template %s: %v", s.Template.Name, err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return types.StepSlice{}, types.SecretSlice{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	return config.Steps, config.Secrets, nil
}
