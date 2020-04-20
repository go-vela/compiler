package native

import (
	"bytes"
	"fmt"
	"text/template"

	types "github.com/go-vela/types/yaml"

	"github.com/Masterminds/sprig"

	yaml "gopkg.in/yaml.v2"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, s *types.Step) (types.StepSlice, error) {
	buffer := new(bytes.Buffer)
	config := new(types.Build)

	// parse the template with Masterminds/sprig functions
	t, err := template.New(s.Name).Funcs(sprig.TxtFuncMap()).Parse(tmpl)
	if err != nil {
		return types.StepSlice{}, fmt.Errorf("unable to parsing template %s: %v", s.Template.Name, err)
	}

	// apply the variables to the parsed template
	err = t.Execute(buffer, s.Template.Variables)
	if err != nil {
		return types.StepSlice{}, fmt.Errorf("unable to execute template %s: %v", s.Template.Name, err)
	}

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buffer.Bytes(), config)
	if err != nil {
		return types.StepSlice{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	return config.Steps, nil
}
