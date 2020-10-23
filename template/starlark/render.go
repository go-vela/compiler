// Copyright (c) 2020 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package starlark

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-vela/types/raw"
	types "github.com/go-vela/types/yaml"
	"go.starlark.net/starlark"
	"gopkg.in/yaml.v2"

	"go.starlark.net/starlarkstruct"
)

// Render combines the template with the step in the yaml pipeline.
func Render(tmpl string, s *types.Step) (types.StepSlice, error) {
	// TODO: investigate way to not use tmp filesystem

	// setup filesystem
	file, err := ioutil.TempFile("/tmp", "sample")
	if err != nil {
		return nil, err
	}

	defer os.Remove(file.Name())

	_, err = file.Write([]byte(tmpl))
	if err != nil {
		return nil, err
	}

	config := new(types.Build)

	thread := &starlark.Thread{Name: "my thread"}
	globals, err := starlark.ExecFile(thread, file.Name(), nil, nil)
	if err != nil {
		return nil, err
	}

	mainVal, ok := globals["main"]
	if !ok {
		return nil, fmt.Errorf("no main function found")
	}
	main, ok := mainVal.(starlark.Callable)
	if !ok {
		return nil, fmt.Errorf("main must be a function")
	}

	userVars, err := userData(s.Template.Variables)
	if err != nil {
		return nil, err
	}

	velaVars, err := velaEnvironmentData(s.Environment)
	if err != nil {
		return nil, err
	}

	args := starlark.Tuple([]starlark.Value{
		starlarkstruct.FromStringDict(
			starlark.String("context"),
			starlark.StringDict{
				"vars": starlarkstruct.FromStringDict(starlark.String("vars"), userVars),
				"vela": starlarkstruct.FromStringDict(starlark.String("vela"), velaVars),
			},
		),
	})

	// Call Starlark function from Go.
	mainVal, err = starlark.Call(thread, main, args, nil)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	switch v := mainVal.(type) {
	case *starlark.List:
		for i := 0; i < v.Len(); i++ {
			item := v.Index(i)
			buf.WriteString("---\n")
			err = writeJSON(buf, item)
			if err != nil {
				return nil, err
			}
			buf.WriteString("\n")
		}
	case *starlark.Dict:
		buf.WriteString("---\n")
		err = writeJSON(buf, v)
		if err != nil {
			return nil, err
		}
	default:
		return nil, fmt.Errorf("invalid return type (got a %s)", mainVal.Type())
	}

	fmt.Println(buf.String())

	// unmarshal the template to the pipeline
	err = yaml.Unmarshal(buf.Bytes(), config)
	if err != nil {
		return types.StepSlice{}, fmt.Errorf("unable to unmarshal yaml: %v", err)
	}

	// ensure all templated steps have template prefix
	for index, newStep := range config.Steps {
		config.Steps[index].Name = fmt.Sprintf("%s_%s", s.Name, newStep.Name)
	}

	fmt.Println(buf.String())

	return config.Steps, nil
}

func userData(m map[string]interface{}) (starlark.StringDict, error) {
	dict := make(starlark.StringDict)

	for key, value := range m {
		val, err := toStarlark(value)
		if err != nil {
			return nil, err
		}
		dict[key] = val
	}

	return dict, nil
}

func velaEnvironmentData(slice raw.StringSliceMap) (starlark.StringDict, error) {
	build := starlark.NewDict(0)
	repo := starlark.NewDict(0)
	user := starlark.NewDict(0)
	system := starlark.NewDict(0)

	// TODO look at fixing access to variables:
	// ctx.vela.repo["full_name"] vs ctx.vela.repo.full_name
	dict := starlark.StringDict{
		"build":  build,
		"repo":   repo,
		"user":   user,
		"system": system,
	}

	for key, value := range slice {
		key = strings.ToLower(key)
		if strings.HasPrefix(key, "vela_") {
			key = strings.TrimPrefix(key, "vela_")

			switch {
			case strings.HasPrefix(key, "build_"):
				err := build.SetKey(starlark.String(strings.TrimPrefix(key, "build_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(key, "repo_"):
				err := repo.SetKey(starlark.String(strings.TrimPrefix(key, "repo_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(key, "user_"):
				err := user.SetKey(starlark.String(strings.TrimPrefix(key, "user_")), starlark.String(value))
				if err != nil {
					return nil, err
				}
			default:
				err := system.SetKey(starlark.String(key), starlark.String(value))
				if err != nil {
					return nil, err
				}
			}
		}
	}

	return dict, nil
}
