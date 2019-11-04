// Copyright (c) 2019 Target Brands, Inc. All rights reserved.
//
// Use of this source code is governed by the LICENSE file in this repository.

package compiler

import (
	"github.com/go-vela/types/library"
	"github.com/go-vela/types/pipeline"
	"github.com/go-vela/types/yaml"
)

// Engine represents an interface for converting a yaml
// configuration to an executable pipeline for Vela.
type Engine interface {
	// Compiler Interface Functions

	// Compile defines a function that produces an executable
	// representation of a pipeline from an object. This calls
	// Parse internally to convert the object to a yaml configuration.
	Compile(interface{}) (*pipeline.Build, error)

	// Parse defines a function that converts
	// an object to a yaml configuration.
	Parse(interface{}) (*yaml.Build, error)

	// Validate defines a function that verifies
	// the yaml configuration is accurate.
	Validate(*yaml.Build) error

	// Clone Compiler Interface Functions

	// CloneStage defines a function that injects the
	// stage clone process into a yaml configuration.
	CloneStage(*yaml.Build) (*yaml.Build, error)
	// CloneStep defines a function that injects the
	// step clone process into a yaml configuration.
	CloneStep(*yaml.Build) (*yaml.Build, error)

	// Environment Compiler Interface Functions

	// EnvironmentStages defines a function that injects the environment
	// variables for each step in every stage into a yaml configuration.
	EnvironmentStages(yaml.StageSlice) (yaml.StageSlice, error)
	// EnvironmentSteps defines a function that injects the environment
	// variables for each step into a yaml configuration.
	EnvironmentSteps(yaml.StepSlice) (yaml.StepSlice, error)

	// Expand Compiler Interface Functions

	// ExpandStages defines a function that injects the template
	// for each templated step in every stage in a yaml configuration.
	ExpandStages(yaml.StageSlice, map[string]*yaml.Template) (yaml.StageSlice, error)
	// ExpandSteps defines a function that injects the template
	// for each templated step in a yaml configuration.
	ExpandSteps(yaml.StepSlice, map[string]*yaml.Template) (yaml.StepSlice, error)

	// Script Compiler Interface Functions

	// ScriptStages defines a function that injects the script
	// for each step in every stage in a yaml configuration.
	ScriptStages(yaml.StageSlice) (yaml.StageSlice, error)
	// ScriptSteps defines a function that injects the script
	// for each step in a yaml configuration.
	ScriptSteps(yaml.StepSlice) (yaml.StepSlice, error)

	// Transform Compiler Interface Functions

	// TransformStages defines a function that converts a yaml
	// configuration with stages into an executable pipeline.
	TransformStages(*pipeline.RuleData, *yaml.Build) (*pipeline.Build, error)
	// TransformSteps defines a function that converts a yaml
	// configuration with steps into an executable pipeline.
	TransformSteps(*pipeline.RuleData, *yaml.Build) (*pipeline.Build, error)

	// With Compiler Interface Functions

	// WithBuild defines a function that sets
	// the library build type in the Engine.
	WithBuild(*library.Build) Engine
	// WithFiles defines a function that sets
	// the changeset files in the Engine.
	WithFiles([]string) Engine
	// WithRepo defines a function that sets
	// the library repo type in the Engine.
	WithRepo(*library.Repo) Engine
	// WithUser defines a function that sets
	// the library user type in the Engine.
	WithUser(*library.User) Engine
}
